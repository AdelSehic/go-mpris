package mpris

import (
	"math"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/rs/zerolog"
)

var log zerolog.Logger
var conn *dbus.Conn

func InitLogger(l zerolog.Logger) {
	log = l
}

func StartListening() (chan *dbus.Signal, error) {
	var err error
	conn, err = dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	if err := conn.AddMatchSignal(
		dbus.WithMatchInterface(SERVICE_DBUS),
		dbus.WithMatchMember(PROPERTY_OWNER_CHANGE),
		dbus.WithMatchArg0Namespace(SERVICE_MPRIS),
	); err != nil {
		return nil, err
	}

	// Listen for property changes (track changes, play/pause, etc.)
	if err := conn.AddMatchSignal(
		dbus.WithMatchInterface(DBUS_PROPERTIES),
		dbus.WithMatchMember(PROPERTIES_CHANGED),
		dbus.WithMatchPathNamespace(OBJECT_MPRIS),
	); err != nil {
		return nil, err
	}

	signals := make(chan *dbus.Signal, 10)
	conn.Signal(signals)

	if err := SetActivePlayer(); err != nil {
		return nil, err
	}

	return signals, err
}

func SetActivePlayer() error {
	names, err := GetActivePlayers()
	if err != nil {
		return err
	}

	players := make([]string, 0, 5)
	log.Debug().Msg("Searching for a player ...")
	for _, n := range names {
		if strings.HasSuffix(n, "playerctld") {
			continue
		}
		players = append(players, n)
	}

	if len(players) == 0 {
		log.Debug().Msg("No active players found")
		ActivePlayer = NilPlayer
		return nil
	}

	playerIndex := findPlayerByPriority(players)
	SetPlayer(&Signal{Value: players[playerIndex]})

	return nil
}

func findPlayerByPriority(players []string) int {
	firefox := math.MaxInt
	for k, name := range players {
		if strings.Contains(name, "spotify") {
			return k
		} else if strings.Contains(name, "firefox") {
			firefox = k
		}
	}
	if firefox != math.MaxInt {
		return firefox
	}
	return 0
}

func GetActivePlayers() ([]string, error) {
	names := make([]string, 0)
	if err := conn.BusObject().Call(MPRIS_LIST_PLAYERS, 0).Store(&names); err != nil {
		return nil, err
	}
	playerNames := make([]string, 0)
	for _, name := range names {
		if strings.HasPrefix(name, SERVICE_MPRIS) {
			playerNames = append(playerNames, name)
		}
	}
	return playerNames, nil
}

func Close() {
	conn.Close()
}

package mpris

import (
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/rs/zerolog"
)

const (
	SERVICE_DBUS  = "org.freedesktop.DBus"
	SERVICE_MPRIS = "org.mpris.MediaPlayer2"

	PLAYER_PLAY  = SERVICE_MPRIS + ".Player.Play"
	PLAYER_PAUSE = SERVICE_MPRIS + ".Player.Pause"

	MPRIS_LIST_PLAYERS = SERVICE_DBUS + ".ListNames"

	PROPERTY_OWNER_CHANGE = "NameOwnerChanged"
	OBJECT_MPRIS          = "/org/mpris/MediaPlayer2"
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

	log.Debug().Msg("Searching for a player a player ...")
	for _, n := range names {
		AddPlayer(&Signal{Name: n})
		return nil
	}
	log.Debug().Msg("No active players found")
	return nil
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

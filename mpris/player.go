package mpris

import (
	"sync"

	"github.com/godbus/dbus/v5"
)

var Players sync.Map
var ActivePlayer *Player

type Player struct {
	Name   string
	ID     string
	State  bool
	Object dbus.BusObject
}

func AddPlayer(data *Signal) {
	state := false
	if data.NewOwner != "" {
		state = true
	}
	playerObj := conn.Object(data.Name, OBJECT_MPRIS)
	if playerObj == nil {
		log.Error().Msg("Failed to create player object")
		return
	}
	newPlayer := &Player{
		Name:   data.Name,
		ID:     data.NewOwner,
		State:  state,
		Object: playerObj,
	}
	Players.Store(data.Name, newPlayer)
	ActivePlayer = newPlayer
	log.Debug().Msgf("%s set as the active player", newPlayer.Name)
}

func GetPlayer(name string) (*Player, bool) {
	player, err := Players.Load(name)
	return player.(*Player), err
}

func DeletePlayer(name string) {
	Players.Delete(name)
}

func PlayerExists(name string) bool {
	_, exists := Players.Load(name)
	return exists
}

func PlayerState(name string) bool {
	player, ok := GetPlayer(name)
	if !ok {
		return false
	}
	return player.State
}

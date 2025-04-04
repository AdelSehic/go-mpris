package mpris

import (
	"github.com/godbus/dbus/v5"
)

var ActivePlayer *Player

type Player struct {
	Name     string
	ID       string
	State    bool
	Object   dbus.BusObject
	Metadata map[string]any
}

func SetPlayer(data *Signal) {
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
	ActivePlayer = newPlayer
	log.Debug().Msgf("%s set as the active player", newPlayer.Name)

	ActivePlayer.UpdatePlayerState()
	ActivePlayer.UpdatePlayerMetadata()
}

func (p *Player) PlayerState() bool {
	return p.State
}

func (p *Player) UpdatePlayerState() {
	variant, err := p.Object.GetProperty(PLAYER_STATUS)
	if err != nil {
		log.Error().Err(err).Str("Player Name", p.Name).Msg("Failed to get player state")
		return
	}

	status, ok := variant.Value().(string)
	if !ok {
		log.Error().Msg("Failed to cast player variant into string")

	}
	log.Debug().Str("status", status).Str("player", p.Name).Msg("Player status retrieved")

	p.State = status == STATUS_PLAYING
}

func (p *Player) UpdatePlayerMetadata() {
	variant, err := p.Object.GetProperty(PLAYER_METADATA)
	if err != nil {
		log.Error().Err(err).Str("Player Name", p.Name).Msg("Failed to get player metadata")
		return
	}

	rawData, ok := variant.Value().(map[string]dbus.Variant)
	if !ok {
		log.Error().Msg("Failed to cast player metadata into map[string]any")

	}
	log.Debug().Str("player", p.Name).Msg("Player metadata retrieved")

	data := make(map[string]any)
	for k, v := range rawData {
		data[k] = v.Value()
	}

	p.Metadata = data
}

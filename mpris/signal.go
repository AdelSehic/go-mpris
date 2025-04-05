package mpris

import (
	"errors"
	"fmt"

	"github.com/godbus/dbus/v5"
)

const (
	TYPE_PLAYER_CHANGE = "player_change"
	TYPE_TRACK_CHANGE  = "track_change"
	TYPE_STATUS_CHANGE = "status_change"
	TYPE_UNHANDLED     = "unhandled"

	FIELD_PLAYBACK_STATUS = "PlaybackStatus"
	FIELD_METADATA        = "Metadata"
)

type Signal struct {
	Value    any
	OldOwner string
	NewOwner string
	Type     string
}

func ParseSignal(body []any) (*Signal, error) {
	name, ok1 := body[0].(string)
	oldOwner, ok2 := body[1].(string)
	newOwner, ok3 := body[2].(string)

	if ok1 && ok2 && ok3 {
		return &Signal{
			Type:     TYPE_PLAYER_CHANGE,
			Value:    name,
			OldOwner: oldOwner,
			NewOwner: newOwner,
		}, nil
	}

	dbusData, ok2 := body[1].(map[string]dbus.Variant)
	if !ok2 {
		return nil, errors.New(fmt.Sprintf("Unexpected signal format: %+v", body))
	}

	data := make(map[string]any, 0)
	for k, v := range dbusData {
		data[k] = v
	}

	if status, ok := data[FIELD_PLAYBACK_STATUS]; ok {
		state := status.(dbus.Variant).String()
		return &Signal{
			Type:  TYPE_STATUS_CHANGE,
			Value: state == STATUS_PLAYING,
		}, nil
	}

	if dbusData, ok := data[FIELD_METADATA]; ok {
		dataMap := dbusData.(dbus.Variant).Value().(map[string]dbus.Variant)
		return &Signal{
			Type:  TYPE_TRACK_CHANGE,
			Value: parseMetadata(dataMap),
		}, nil
	}

	return &Signal{
		Type: TYPE_UNHANDLED,
	}, nil

}

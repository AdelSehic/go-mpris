package mpris

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (p *Player) PlayPause() error {
	call := p.Object.Call(PLAYER_PLAYPAUSE, 0)
	if call.Err != nil {
		return call.Err
	}
	return nil
}

func (p *Player) Play() error {
	call := p.Object.Call(PLAYER_PLAY, 0)
	if call.Err != nil {
		return call.Err
	}
	return nil
}

func (p *Player) Pause() error {
	call := p.Object.Call(PLAYER_PAUSE, 0)
	if call.Err != nil {
		return call.Err
	}
	return nil
}

func (p *Player) Next() error {
	call := p.Object.Call(PLAYER_NEXT, 0)
	if call.Err != nil {
		return call.Err
	}
	return nil
}

func (p *Player) Previous() error {
	call := p.Object.Call(PLAYER_PREVIOUS, 0)
	if call.Err != nil {
		return call.Err
	}
	return nil
}

func (p *Player) FormattedMetadata() string {
	return fmt.Sprintf("%s - %s", strings.Join(p.Meta.Artist, ", "), p.Meta.Title)
}

func (p *Player) JSONMetadata() ([]byte, error) {
	data, err := json.Marshal(p.Meta)
	if err != nil {
		return nil, err
	}
	return data, nil
}

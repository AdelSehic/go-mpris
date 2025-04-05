package mpris

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (p *Player) PlayPause() {
	call := p.Object.Call(PLAYER_PLAYPAUSE, 0)
	if call.Err != nil {
		panic(call.Err)
	}
}

func (p *Player) Play() {
	call := p.Object.Call(PLAYER_PLAY, 0)
	if call.Err != nil {
		panic(call.Err)
	}
}

func (p *Player) Pause() {
	call := p.Object.Call(PLAYER_PAUSE, 0)
	if call.Err != nil {
		panic(call.Err)
	}
}

func (p *Player) Next() {
	call := p.Object.Call(PLAYER_NEXT, 0)
	if call.Err != nil {
		panic(call.Err)
	}
}

func (p *Player) Previous() {
	call := p.Object.Call(PLAYER_PREVIOUS, 0)
	if call.Err != nil {
		panic(call.Err)
	}
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

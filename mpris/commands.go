package mpris

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

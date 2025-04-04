package mpris

func (p *Player) Play() {
	call := p.Object.Call(PLAYER_PLAY, 0)
	p.State = true
	if call.Err != nil {
		panic(call.Err)
	}
}

func (p *Player) Pause() {
	call := p.Object.Call(PLAYER_PAUSE, 0)
	p.State = false
	if call.Err != nil {
		panic(call.Err)
	}
}

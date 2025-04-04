package mpris

type Signal struct {
	Name     string
	OldOwner string
	NewOwner string
}

func ParseSignal(body []any) *Signal {
	name, ok1 := body[0].(string)
	oldOwner, ok2 := body[1].(string)
	newOwner, ok3 := body[2].(string)

	if ok1 && ok2 && ok3 {
		return &Signal{
			Name:     name,
			OldOwner: oldOwner,
			NewOwner: newOwner,
		}
	} else {
		log.Error().Msgf("Unexpected signal format: %+v", body)
		return nil
	}
}

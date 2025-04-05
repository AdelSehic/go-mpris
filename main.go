package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/AdelSehic/mpris-go/logger"
	"github.com/AdelSehic/mpris-go/mpris"
	"github.com/rs/zerolog"
)

func main() {
	logger.InitLogger(zerolog.DebugLevel)
	logger := logger.GetLogger()
	mpris.InitLogger(logger)
	logger.Info().Msg("Starting the application")

	if err := StartPipes(); err != nil {
		logger.Fatal().Err(err).Msg("Couldn't start the pipes")
	}
	logger.Debug().Msg("Initialized the pipes")

	signals, err := mpris.StartListening()
	if err != nil {
		logger.Fatal().Err(err).Msg("Issue starting dbus listening")
	}
	defer mpris.Close()
	logger.Debug().Msg("Dbus session initated")

	go scanInput(&logger)
	metaChan := StartTextScroller(pipes[PIPE_TITLE], 25, 250, 3000)
	if mpris.ActivePlayer != nil {
		metaChan <- mpris.ActivePlayer.Meta
		WriteStatusIcon()
	}

	logger.Debug().Msg("App set up")
	var data *mpris.Signal
	for signal := range signals {
		logger.Debug().Msgf("New message: %+v", signal)
		if data, err = mpris.ParseSignal(signal.Body); err != nil {
			logger.Error().Err(err).Msg("Couldn't parse signal body")
			continue
		}
		switch data.Type {
		case mpris.TYPE_PLAYER_CHANGE:
			mpris.SetActivePlayer()
			metaChan <- mpris.ActivePlayer.Meta
			pipes[PIPE_STATUS_ICON].Write([]byte(StatusIcon()))
		case mpris.TYPE_STATUS_CHANGE:
			mpris.ActivePlayer.State = data.Value.(bool)
			WriteStatusIcon()
			logger.Info().Msg("Player status changed")
		case mpris.TYPE_TRACK_CHANGE:
			mpris.ActivePlayer.Meta = data.Value.(*mpris.Metadata)
			metaChan <- data.Value.(*mpris.Metadata)
			logger.Info().Msg("Track changed")
		default:
			logger.Info().Str("Message type", data.Type).Msg("Recieved an unkown message type")
		}
	}
}

func scanInput(logger *zerolog.Logger) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		input := scanner.Text()
		if mpris.ActivePlayer == mpris.NilPlayer || mpris.ActivePlayer.Name == mpris.SERVICE_MPRIS+".playerctld" {
			continue
		}
		switch input {
		case "p":
			mpris.ActivePlayer.PlayPause()
			logger.Info().Msg("Toggle")
		case "d":
			mpris.ActivePlayer.Play()
			logger.Info().Msg("Playing")
		case "s":
			mpris.ActivePlayer.Pause()
			logger.Info().Msg("Pausing")
		case "n":
			mpris.ActivePlayer.Next()
			logger.Info().Msg("Next")
		case "b":
			mpris.ActivePlayer.Previous()
			logger.Info().Msg("Previous")
		case "l":
			players, _ := mpris.GetActivePlayers()
			logger.WithLevel(zerolog.NoLevel).Str("players", strings.Join(players, ", ")).Msg("Active players")
		case "m":
			mpris.ActivePlayer.UpdatePlayerMetadata()
			logger.WithLevel(zerolog.NoLevel).Any("Metadata", mpris.ActivePlayer.Meta).Msg("Player metadata")
		case "t":
			fmt.Println(mpris.ActivePlayer.FormattedMetadata())
		case "j":
			data, err := mpris.ActivePlayer.JSONMetadata()
			if err != nil {
				logger.Error().Err(err).Msg("Failed to marshall Metadata into JSON")
			} else {
				fmt.Println(string(data))
			}
		case "a":
			fmt.Printf("Active Player: %+v\r\n", mpris.ActivePlayer)
		default:
			logger.Info().Msg("Unkown command")
		}
	}
}

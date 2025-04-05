package main

import (
	"bufio"
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

	signals, err := mpris.StartListening()
	if err != nil {
		logger.Fatal().Err(err).Msg("Issue starting dbus listening")
	}
	defer mpris.Close()

	go scanInput(&logger)

	for signal := range signals {
		logger.Debug().Msgf("New message: %+v", signal.Body...)
		data := mpris.ParseSignal(signal.Body)
		if data.NewOwner == "" {
			mpris.SetActivePlayer()
			continue
		}
		mpris.SetPlayer(data)
	}
}

func scanInput(logger *zerolog.Logger) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		input := scanner.Text()
		if mpris.ActivePlayer == nil || mpris.ActivePlayer.Name == mpris.SERVICE_MPRIS+".playerctld" {
			continue
		}
		switch input {
		case "p":
			mpris.ActivePlayer.Play()
			logger.Info().Msg("Playing")
		case "s":
			mpris.ActivePlayer.Pause()
			logger.Info().Msg("Pausing")
		case "l":
			players, _ := mpris.GetActivePlayers()
			logger.Info().Str("players", strings.Join(players, ", ")).Msg("Active players")
		case "m":
			mpris.ActivePlayer.UpdatePlayerMetadata()
			logger.Info().Any("Metadata", mpris.ActivePlayer.Metadata).Msg("Player metadata")
		default:
			logger.Info().Msg("Unkown command")
		}
	}
}

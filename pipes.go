package main

import (
	"os"
	"syscall"
)

const (
	DIR              = "/tmp/go-mpris"
	PIPE_TITLE       = "/tmp/go-mpris/title"
	PIPE_STATUS_ICON = "/tmp/go-mpris/status-icon"
)

var pipesToCreate = []string{PIPE_TITLE, PIPE_STATUS_ICON}

var pipes map[string]*os.File

func StartPipes() error {
	pipes = make(map[string]*os.File, 0)

	if err := os.MkdirAll(DIR, 0755); err != nil {
		return err
	}

	for _, pipeName := range pipesToCreate {
		if _, err := os.Stat(pipeName); os.IsNotExist(err) {
			if err := syscall.Mknod(pipeName, syscall.S_IFIFO|0666, 0); err != nil {
				return err
			}
		}

		pipeFile, err := os.OpenFile(pipeName, os.O_WRONLY|os.O_CREATE, os.ModeNamedPipe)
		if err != nil {
			return err
		}
		pipes[pipeName] = pipeFile
	}

	return nil
}

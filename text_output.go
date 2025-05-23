package main

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/AdelSehic/mpris-go/mpris"
)

func FormatMetadata(meta *mpris.Metadata) string {
	if len(meta.Artist) == 0 || meta.Title == "" {
		return "nothing playing . . ."
	}
	return fmt.Sprintf("%s - %s", strings.Join(meta.Artist, ", "), meta.Title)
}

func StartTextScroller(w io.Writer, maxLength, interval, pauseLen int) chan *mpris.Metadata {
	dataChan := make(chan *mpris.Metadata, 10)
	startIndex, endIndex, strLen := 0, 0, 0
	output := ""
	paused := true
	stopped := false
	pauseUntil := time.Now().Add(time.Duration(pauseLen) * time.Millisecond)

	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
		for {
			select {
			case data := <-dataChan:
				output = FormatMetadata(data)
				strLen = len(output)
				endIndex = strLen + 2*maxLength
				output = output + strings.Repeat(" ", maxLength) + output + strings.Repeat(" ", maxLength)
				end := min(startIndex+maxLength, endIndex)
				w.Write([]byte(output[startIndex:end] + "\n"))
				startIndex = 0
				stopped = strLen < maxLength
			case <-ticker.C:
				if mpris.ActivePlayer == mpris.NilPlayer {
					w.Write([]byte(" \n"))
					startIndex = 0
					WriteStatusIcon()
					continue
				}
				if paused && time.Now().Before(pauseUntil) {
					continue
				}
				if startIndex >= strLen+maxLength {
					startIndex = 0
					paused = true
					pauseUntil = time.Now().Add(time.Duration(pauseLen) * time.Millisecond)
				}
				end := min(startIndex+maxLength, endIndex)
				w.Write([]byte(output[startIndex:end] + "\n"))
				startIndex++
				if !mpris.ActivePlayer.State || stopped {
					startIndex = 0
				}
			}
		}
	}()

	return dataChan
}

func WriteStatusIcon() {
	pipes[PIPE_STATUS_ICON].Write([]byte(StatusIcon()))
}

func StatusIcon() string {
	statusIcon := "󰐊"
	if mpris.ActivePlayer == mpris.NilPlayer {
		statusIcon = ""
	} else if mpris.ActivePlayer.State {
		statusIcon = "󰏤"
	}
	return statusIcon + "\n"
}

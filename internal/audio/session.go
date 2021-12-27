package audio

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/diamondburned/oggreader"
)

const (
	stoppedState = iota
	changingState
	pausedState
	playingState
)

type Session struct {
	Connection discord.VoiceConnection

	source string
	isOpus bool

	Position time.Duration

	state   int
	channel chan int

	context context.Context
	cancel  context.CancelFunc
}

func New(connection discord.VoiceConnection) *Session {
	return &Session{Connection: connection}
}

func (s *Session) PlayURL(source string, isOpus bool) error {
	if s.state != stoppedState && s.state != changingState {
		s.Stop()
	}

	s.context, s.cancel = context.WithCancel(context.Background())
	s.source, s.isOpus = source, isOpus

	codec := "libopus"
	if s.isOpus {
		codec = "copy"
	}

	ffmpeg := exec.CommandContext(
		s.context, "ffmpeg",
		"-loglevel", "error", "-reconnect", "1", "-reconnect_streamed", "1",
		"-reconnect_delay_max", "5", /*"-ss", utils.FormatTime(s.Position),*/
		"-i", source, "-vn", "-codec", codec, "-vbr", "off",
		"-frame_duration", "20", "-f", "opus", "-",
	)

	stdout, err := ffmpeg.StdoutPipe()
	if err != nil {
		s.stop()
		return fmt.Errorf("failed to get ffmpeg stdout: %w", err)
	}

	var stderr bytes.Buffer
	ffmpeg.Stderr = &stderr

	if err := ffmpeg.Start(); err != nil {
		s.stop()
		return fmt.Errorf("failed to start ffmpeg process: %w", err)
	}

	if err := s.SendSpeaking(); err != nil {
		s.stop()
		return fmt.Errorf("failed to send speaking packet to discord: %w", err)
	}

	s.setState(playingState)

	if err := oggreader.DecodeBuffered(s, stdout); err != nil && s.state != changingState {
		s.stop()
		return err
	}

	if err, std := ffmpeg.Wait(), stderr.String(); err != nil && s.state != changingState && std != "" {
		s.stop()
		return fmt.Errorf("ffmpeg returned error")
	}

	if s.state == changingState {
		return s.PlayURL(s.source, s.isOpus)
	}

	s.stop()
	return nil
}

func (s *Session) Destroy() {
	s.Stop()
	_ = s.Connection.Disconnect()
}

func (s *Session) Seek(position time.Duration) {
	if s.state == stoppedState {
		return
	}

	s.Position = position
	s.setState(changingState)
	s.Stop()
}

func (s *Session) Resume() {
	if s.state == pausedState {
		s.setState(playingState)
		_ = s.SendSpeaking()
	}
}

func (s *Session) Pause() {
	if s.state != stoppedState && s.state != changingState {
		s.setState(pausedState)
	}
}

func (s *Session) TogglePause() {
	if s.state == stoppedState || s.state == changingState {
		return
	}
	if s.IsPaused() {
		s.Resume()
		return
	}
	s.Pause()
}

func (s *Session) IsPaused() bool { return s.state == pausedState }

func (s *Session) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *Session) SendSpeaking() error {
	return s.Connection.Speaking(true)
}

func (s *Session) Write(data []byte) (int, error) {
	if s.state == stoppedState || s.state == changingState {
		return 0, context.Canceled
	}

	if s.state == pausedState {
		s.channel = make(chan int)

		for {
			if newState := <-s.channel; newState != pausedState {
				close(s.channel)
				s.channel = nil
				break
			}
		}
	}

	s.Position = s.Position + (20 * time.Millisecond)
	return len(data), s.Connection.WriteOpus(data) // FIXME
}

func (s *Session) setState(state int) {
	s.state = state

	if s.channel != nil {
		s.channel <- state
	}
}

func (s *Session) stop() {
	s.cancel()
	s.setState(stoppedState)
	s.Position = 0
}

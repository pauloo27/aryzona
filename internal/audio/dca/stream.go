package dca

import (
	"io"
	"sync"
	"time"

	"github.com/pauloo27/aryzona/internal/core/routine"
	"github.com/pauloo27/aryzona/internal/discord/model"
)

// based on https://git.notagovernment.agency/ItsClairton/Anny/ which é...
// Isso é baseado no https://github.com/jonas747/dca porém com algumas correções e mais básico

type OpusReader interface {
	OpusFrame() (frame []byte, err error)
	FrameDuration() time.Duration
}

type StreamingSession struct {
	connection model.VoiceConnection
	err        error
	framesSent int

	source *EncodeSession

	callback chan error
	sync.Mutex
	running  bool
	paused   bool
	finished bool
}

func NewStream(source *EncodeSession, vc model.VoiceConnection, callback chan error) *StreamingSession {

	session := &StreamingSession{
		source:     source,
		connection: vc,
		callback:   callback,
	}

	routine.GoAndRecover(session.stream)

	return session
}

func (s *StreamingSession) stream() {
	s.Lock()

	if s.running {
		s.Unlock()
		return
	}

	s.running = true
	s.Unlock()

	defer func() {
		s.Lock()
		s.running = false
		s.Unlock()
	}()

	for {
		s.Lock()

		if s.paused {
			s.Unlock()
			return
		}
		s.Unlock()
		err := s.readNext()

		if err != nil {
			s.Lock()
			s.finished = true

			if err != io.EOF {
				s.err = err
			}

			if s.callback != nil {
				s.callback <- err
			}

			s.Unlock()
			break
		}

	}

}

func (s *StreamingSession) readNext() error {
	opus, err := s.source.OpusFrame()
	if err != nil {
		return err
	}

	_, err = s.connection.WriteOpus(opus)
	if err != nil {
		return err
	}
	s.Lock()
	s.framesSent++
	s.Unlock()

	return nil
}

func (s *StreamingSession) PlaybackPosition() int {
	s.Lock()
	time := s.framesSent * int(s.source.FrameDuration())
	s.Unlock()
	return time
}

func (s *StreamingSession) Pause() {
	s.paused = true
}

func (s *StreamingSession) Resume() {
	s.Lock()
	defer s.Unlock()
	if !s.paused {
		return
	}
	s.paused = false
	routine.GoAndRecover(s.stream)
}

func (s *StreamingSession) TogglePause() {
	if s.finished {
		return
	}

	if s.paused {
		s.Resume()
	} else {
		s.Pause()
	}
}

func (s *StreamingSession) Finished() (bool, error) {
	s.Lock()

	err := s.err
	state := s.finished

	s.Unlock()
	return state, err
}

func (s *StreamingSession) Paused() bool {
	s.Lock()

	state := s.paused

	s.Unlock()
	return state
}

func (s *StreamingSession) Source() *EncodeSession {
	s.Lock()

	source := s.source

	s.Unlock()
	return source
}

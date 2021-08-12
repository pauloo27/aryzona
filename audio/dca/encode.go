package dca

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/Pauloo27/aryzona/logger"
	"github.com/Pauloo27/aryzona/utils"
	"github.com/jonas747/ogg"
)

// based on https://git.notagovernment.agency/ItsClairton/Anny/ which é...
// Isso é baseado no https://github.com/jonas747/dca porém com algumas correções e mais básico

var remoteFfmpegArgs = []string{
	"-reconnect", "1",
	"-reconnect_streamed", "1",
}

type EncodeSession struct {
	sync.Mutex
	path   string
	reader io.Reader

	ruining bool
	started time.Time
	channel chan []byte
	process *os.Process

	lastFrame int
	err       error
	isOpus    bool
	isLocal   bool

	buff bytes.Buffer
}

func EncodeData(path string, isOpus, isLocal bool) *EncodeSession {

	session := &EncodeSession{
		path:    path,
		channel: make(chan []byte, 100),
		isOpus:  isOpus,
		isLocal: isLocal,
	}

	go session.run()
	return session
}

/* #nosec G204 */
func (e *EncodeSession) run() error {
	defer func() {
		e.Lock()
		e.ruining = false
		e.Unlock()
	}()

	e.Lock()
	e.ruining = true

	commonArgs := []string{
		"-i", e.path,
		"-acodec", utils.ConditionalString(e.isOpus, "copy", "libopus"),
		"-analyzeduration", "0",
		"-loglevel", "0",
		"-map", "0:a",
		"-f", "ogg",
		"-ar", "48000",
		"-ac", "2",
		"-application", "lowdelay",
		"-frame_duration", "20", "pipe:1",
	}

	var ffmpegArgs []string

	if !e.isLocal {
		ffmpegArgs = append(remoteFfmpegArgs, commonArgs...)
	} else {
		ffmpegArgs = commonArgs
	}

	ffmpeg := exec.Command("ffmpeg", ffmpegArgs...)

	if e.reader != nil {
		ffmpeg.Stdin = e.reader
	}

	stdout, err := ffmpeg.StdoutPipe()
	if err != nil {
		e.Unlock()
		close(e.channel)
		return err
	}

	ffmpeg.Stderr = os.Stdout

	err = ffmpeg.Start()
	if err != nil {
		e.Unlock()
		close(e.channel)
		return err
	}

	e.started = time.Now()

	e.process = ffmpeg.Process
	e.Unlock()

	var wg sync.WaitGroup
	defer close(e.channel)
	e.readStdout(stdout)
	wg.Wait()
	err = ffmpeg.Wait()
	if err != nil {
		if err.Error() != "signal: killed" {
			logger.Errorf("%v", err)
			e.Lock()
			e.err = err
			e.Unlock()
		}
	}
	return nil
}

func (e *EncodeSession) readStdout(stdout io.ReadCloser) error {
	decoder := ogg.NewPacketDecoder(ogg.NewDecoder(stdout))

	skipPackets := 2
	for {
		packet, _, err := decoder.Decode()

		if skipPackets > 0 {
			skipPackets--
			continue
		}
		if err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				return err
			}
			break
		}

		err = e.writeOpusFrame(packet)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *EncodeSession) writeOpusFrame(frame []byte) error {
	var buffer bytes.Buffer

	err := binary.Write(&buffer, binary.LittleEndian, int16(len(frame)))

	if err != nil {
		return err
	}

	_, err = buffer.Write(frame)
	if err != nil {
		return err
	}

	e.channel <- buffer.Bytes()
	e.Lock()
	e.lastFrame++
	e.Unlock()

	return nil
}

func (e *EncodeSession) Stop() error {
	e.Lock()
	defer e.Unlock()
	if !e.ruining || e.process == nil {
		return errors.New("not ruining")
	}

	err := e.process.Kill()
	return err
}

func (e *EncodeSession) Cleanup() {
	e.Stop()
	for range e.channel {
		// empty till closed
		// Cats can be right-pawed or left-pawed.
	}
}

func (e *EncodeSession) FrameDuration() time.Duration {
	return time.Duration(20) * time.Millisecond
}

func (e *EncodeSession) OpusFrame() ([]byte, error) {
	f := <-e.channel
	if f == nil {
		return nil, io.EOF
	}

	if len(f) < 2 {
		return nil, errors.New("bad frame")
	}

	return f[2:], nil
}

func (e *EncodeSession) ReadFrame() ([]byte, error) {
	f := <-e.channel

	if f == nil {
		return nil, io.EOF
	}
	return f, nil
}

func (e *EncodeSession) Read(p []byte) (n int, err error) {
	if e.buff.Len() >= len(p) {
		return e.buff.Read(p)
	}

	for e.buff.Len() < len(p) {
		f, err := e.ReadFrame()
		if err != nil {
			break
		}
		e.buff.Write(f)
	}

	return e.buff.Read(p)
}

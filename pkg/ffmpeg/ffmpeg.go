package ffmpeg

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// EncoderConfig contains configuration parameters for FFmpeg encoding.
type EncoderConfig struct {
	Width        int
	Height       int
	FPS          int
	Output       string
	FFMPEG_Flags []string
}

// NewEncoderConfig creates a new EncoderConfig with default values.
// Height = 1920, Width = 1080, FPS = 60, Output = "output.mp4"
func NewEncoderConfig() *EncoderConfig {
	return &EncoderConfig{
		Width:        1920,
		Height:       1080,
		FPS:          60,
		Output:       "output.mp4",
		FFMPEG_Flags: nil,
	}
}

// Encoder represents an FFmpeg encoder.
type Encoder struct {
	Config *EncoderConfig
	cmd    *exec.Cmd
	stdin  io.WriteCloser
}

// NewEncoder creates a new Encoder with the given configuration.
func NewEncoder(config *EncoderConfig) (*Encoder, error) {
	cmd := exec.Command("ffmpeg")

	args := append(
		config.FFMPEG_Flags,
		"-y",
		"-f", "rawvideo",
		"-pixel_format", "rgba",
		"-video_size", fmt.Sprintf("%dx%d", config.Width, config.Height),
		"-framerate", fmt.Sprintf("%d", config.FPS),
		"-i", "pipe:0",
		config.Output,
	)

	cmd.Args = append(cmd.Args, args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("error creating FFmpeg stdin pipe: %v", err)
	}

	return &Encoder{
		Config: config,
		cmd:    cmd,
		stdin:  stdin,
	}, nil
}

// StartEncoding begins the FFmpeg encoding process.
func (e *Encoder) StartEncoding() error {
	e.cmd.Stderr = os.Stderr
	return e.cmd.Start()
}

// WriteFrame sends a video frame to FFmpeg for encoding.
func (e *Encoder) WriteFrame(frame []byte) error {
	_, err := e.stdin.Write(frame)
	return err
}

// CloseInputPipe closes the input pipe to signal the end of input.
func (e *Encoder) CloseInputPipe() error {
	return e.stdin.Close()
}

// Wait waits for the FFmpeg process to finish.
func (e *Encoder) Wait() error {
	return e.cmd.Wait()
}

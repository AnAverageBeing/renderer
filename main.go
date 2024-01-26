// a simple shit to test code i wrote
package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/AnAverageBeing/renderer/pkg/ffmpeg"
	"github.com/fogleman/gg"
)

const (
	frameWidth  = 1920
	frameHeight = 1080
	ballRadius  = 50
	fps         = 120
	duration    = 30 //sec
)

func main() {
	config := ffmpeg.NewEncoderConfig()
	config.Width = frameWidth
	config.Height = frameHeight
	config.Output = "out.mov"
	config.FPS = fps

	encoder, err := ffmpeg.NewEncoder(config)
	if err != nil {
		log.Fatalf("Error creating FFmpeg encoder: %v", err)
		return
	}
	defer encoder.Wait()

	if err := encoder.StartEncoding(); err != nil {
		log.Fatalf("Error starting FFmpeg encoding: %v", err)
		return
	}

	ctx := gg.NewContext(frameWidth, frameHeight)
	ballX, ballY := rand.Float64()*(frameWidth-2*ballRadius)+ballRadius, rand.Float64()*(frameHeight-2*ballRadius)+ballRadius
	ballSpeedX, ballSpeedY := rand.Float64()*4, rand.Float64()*4

	frameSize := frameWidth * frameHeight * 4
	frame := make([]byte, frameSize)

	log.Println("Initialization complete. Rendering frames...")

	renderFrames(ctx, encoder, frame, ballX, ballY, ballSpeedX, ballSpeedY)

	log.Println("Frame rendering complete. Cleaning up...")
	cleanup(encoder)
}

func renderFrames(ctx *gg.Context, encoder *ffmpeg.Encoder, frame []byte, ballX, ballY, ballSpeedX, ballSpeedY float64) {
	for i := 0; i < fps*duration; i++ {
		ctx.SetRGB(1, 1, 1)
		ctx.Clear()

		ballX += ballSpeedX
		ballY += ballSpeedY

		// Bounce when hitting the boundaries
		if ballX < ballRadius || ballX > frameWidth-ballRadius {
			ballSpeedX = -ballSpeedX
		}
		if ballY < ballRadius || ballY > frameHeight-ballRadius {
			ballSpeedY = -ballSpeedY
		}

		ctx.SetRGB(math.Sin(float64(time.Now().Unix())), math.Cos(float64(time.Now().Unix())), math.Sin(float64(time.Now().Unix())))
		ctx.DrawCircle(ballX, ballY, ballRadius)
		ctx.Fill()

		index := 0

		for y := 0; y < frameHeight; y++ {
			for x := 0; x < frameWidth; x++ {
				r, g, b, a := ctx.Image().At(x, y).RGBA()
				frame[index] = uint8(r & 0xFF)
				frame[index+1] = uint8(g & 0xFF)
				frame[index+2] = uint8(b & 0xFF)
				frame[index+3] = uint8(a & 0xFF)
				index += 4
			}
		}

		if err := encoder.WriteFrame(frame); err != nil {
			log.Printf("Error writing frame to FFmpeg: %v", err)
			break
		}

		time.Sleep(time.Second / time.Duration(fps))
	}
}

func cleanup(encoder *ffmpeg.Encoder) {
	if err := encoder.CloseInputPipe(); err != nil {
		log.Printf("Error closing FFmpeg input pipe: %v", err)
	}
	log.Println("Cleanup complete. Waiting for FFmpeg process to finish...")
	encoder.Wait()
	log.Println("FFmpeg process finished.")
}

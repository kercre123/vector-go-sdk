package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
	"image/color"
	"math/rand"
	"time"
)

func main() {
	var serial = flag.String("serial", "", "Vector's Serial Number")
	flag.Parse()

	sdk_wrapper.InitSDK(*serial)

	ctx := context.Background()
	start := make(chan bool)
	stop := make(chan bool)

	go func() {
		_ = sdk_wrapper.Robot.BehaviorControl(ctx, start, stop)
	}()

	for {
		select {
		case <-start:
			sdk_wrapper.MoveHead(3.0)
			sdk_wrapper.SetBackgroundColor(color.RGBA{0, 0, 0, 0})

			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)
			dieImage := fmt.Sprintf("data/images/dice/%d.png", r1.Intn(6)+1)

			sdk_wrapper.DisplayAnimatedGif("data/images/dice/roll-the-dice.gif", sdk_wrapper.ANIMATED_GIF_SPEED_FAST, 3, false)
			sdk_wrapper.DisplayImageWithTransition(dieImage, 1000, sdk_wrapper.IMAGE_TRANSITION_FADE_IN, 10)
			sdk_wrapper.PlaySound(sdk_wrapper.SYSTEMSOUND_WIN, 100)
			sdk_wrapper.DisplayImage(dieImage, 5000, true)
			stop <- true
			return
		}
	}
}
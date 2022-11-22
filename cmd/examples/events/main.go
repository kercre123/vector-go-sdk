package main

import (
	"context"
	"flag"
	sdk_wrapper "github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
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
			go func() {
				println("Listening for events...")
				isBusy := false
				for {
					evt := sdk_wrapper.WaitForEvent()
					if evt != nil {
						evtRobotState := evt.GetRobotState()
						if evtRobotState != nil {
							if evtRobotState.TouchData.IsBeingTouched == true && !isBusy {
								isBusy = true
								println("I am being touched.")
								sdk_wrapper.PlayAnimation("anim_eyepose_angry", 1, false, false, false)
								_ = sdk_wrapper.PlaySound("data/audio/roar.wav", 100)
								time.Sleep(time.Duration(1000) * time.Millisecond)
								isBusy = false
							}
						}
					}
				}
			}()
			for {
				time.Sleep(time.Duration(100) * time.Millisecond)
			}
			stop <- true
			return
		}
	}
}
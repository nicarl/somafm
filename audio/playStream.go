package audio

import (
	"fmt"
	"net/http"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func PlayRemoteFile(streamUrl string) error {
	resp, err := http.Get(streamUrl)

	if err != nil {
		return err
	}

	streamer, format, err := mp3.Decode(resp.Body)
	if err != nil {
		return err
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	ctrlStreamer := &beep.Ctrl{Streamer: streamer, Paused: false}

	done := make(chan bool)
	speaker.Play(beep.Seq(ctrlStreamer, beep.Callback(func() {
		done <- true
	})))

	for {
		select {
		case <-done:
			return nil
		case <-time.After(time.Second):
			speaker.Lock()
			fmt.Printf("time elapsed: %d of %d\n", format.SampleRate.D(streamer.Position()), format.SampleRate.D(streamer.Len()))
			speaker.Unlock()
		}
	}

}

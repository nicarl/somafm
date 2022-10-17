package audio

import (
	"net/http"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func PlayMusic(streamUrl string, done <-chan bool, errs chan<- error) {
	resp, err := http.Get(streamUrl)

	if err != nil {
		errs <- err
		return
	}

	streamer, _, err := mp3.Decode(resp.Body)
	if err != nil {
		errs <- err
		return
	}
	close(errs)
	defer streamer.Close()

	ctrlStreamer := &beep.Ctrl{Streamer: streamer, Paused: false}
	speaker.Play(ctrlStreamer)

	for {
		<-done
		return
	}
}

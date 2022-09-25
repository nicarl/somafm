package audio

import (
	"net/http"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

// TODO error handling! this is used as goroutine
func PlayMusic(streamUrl string, done <-chan bool, setVolume <-chan float32, errs chan<- error) {
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
	defer streamer.Close()

	ctrlStreamer := &beep.Ctrl{Streamer: streamer, Paused: false}
	volume := &effects.Volume{
		Streamer: ctrlStreamer,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}
	speaker.Play(volume)

	for {
		select {
		case <-done:
			return
		case volumeChange := <-setVolume:
			speaker.Lock()
			volume.Volume += float64(volumeChange)
			speaker.Unlock()
		}
	}
}

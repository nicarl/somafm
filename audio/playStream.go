package audio

import (
	"log"
	"net/http"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func PlayMusic(streamURL string, done chan bool, setVolume <-chan float32) {
	err := playRemoteFile(streamURL, done, setVolume)
	if err != nil {
		log.Fatal(err)
	}
}

func playRemoteFile(streamUrl string, done chan bool, setVolume <-chan float32) error {
	resp, err := http.Get(streamUrl)

	if err != nil {
		return err
	}

	streamer, _, err := mp3.Decode(resp.Body)
	if err != nil {
		return err
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
			return nil
		case volumeChange := <-setVolume:
			speaker.Lock()
			volume.Volume += float64(volumeChange)
			speaker.Unlock()
		}
	}

}

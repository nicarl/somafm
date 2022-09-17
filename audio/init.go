package audio

import (
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

func InitSpeaker() error {
	sampleRate := beep.SampleRate(44100)
	if err := speaker.Init(sampleRate, sampleRate.N(time.Second/10)); err != nil {
		return err
	}
	return nil
}

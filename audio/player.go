package audio

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
)

const targetSampleRate = beep.SampleRate(44100)

var streamClient = &http.Client{
	Transport: &http.Transport{
		TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
	},
}

type Player struct {
	mu            sync.Mutex
	ctrl          *beep.Ctrl
	volumeCtrl    *effects.Volume
	volumeLevel   float64
	playing       bool
	paused        bool
	speakerInited bool
	currentResp   *http.Response
	currentDec    beep.StreamSeekCloser
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) Play(streamURL string) error {
	p.Stop()

	req, err := http.NewRequest("GET", streamURL, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "somafm-cli/1.0")

	resp, err := streamClient.Do(req)
	if err != nil {
		return fmt.Errorf("fetching stream: %w", err)
	}

	streamer, format, err := mp3.Decode(resp.Body)
	if err != nil {
		resp.Body.Close()
		return fmt.Errorf("decoding mp3: %w", err)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.speakerInited {
		bufSize := targetSampleRate.N(200 * time.Millisecond)
		if err := speaker.Init(targetSampleRate, bufSize); err != nil {
			streamer.Close()
			resp.Body.Close()
			return fmt.Errorf("initializing speaker: %w", err)
		}
		p.speakerInited = true
	}

	var audioStream beep.Streamer = streamer
	if format.SampleRate != targetSampleRate {
		audioStream = beep.Resample(4, format.SampleRate, targetSampleRate, streamer)
	}

	p.ctrl = &beep.Ctrl{Streamer: audioStream, Paused: false}
	p.volumeCtrl = &effects.Volume{
		Streamer: p.ctrl,
		Base:     2,
		Volume:   p.volumeLevel,
		Silent:   false,
	}
	p.currentDec = streamer
	p.currentResp = resp
	p.playing = true
	p.paused = false

	speaker.Play(p.volumeCtrl)
	return nil
}

func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cleanupLocked()
}

func (p *Player) cleanupLocked() {
	if !p.playing {
		return
	}

	if p.ctrl != nil {
		speaker.Lock()
		p.ctrl.Paused = true
		speaker.Unlock()
	}
	if p.currentDec != nil {
		p.currentDec.Close()
		p.currentDec = nil
	}
	if p.currentResp != nil {
		p.currentResp.Body.Close()
		p.currentResp = nil
	}
	p.playing = false
	p.paused = false
	p.ctrl = nil
	p.volumeCtrl = nil
}

func (p *Player) TogglePause() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.playing || p.ctrl == nil {
		return
	}

	speaker.Lock()
	p.ctrl.Paused = !p.ctrl.Paused
	p.paused = p.ctrl.Paused
	speaker.Unlock()
}

func (p *Player) VolumeUp() {
	p.adjustVolume(0.5)
}

func (p *Player) VolumeDown() {
	p.adjustVolume(-0.5)
}

func (p *Player) adjustVolume(delta float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.volumeLevel += delta
	if p.volumeLevel < -5 {
		p.volumeLevel = -5
	}
	if p.volumeLevel > 2 {
		p.volumeLevel = 2
	}

	if p.volumeCtrl != nil {
		speaker.Lock()
		p.volumeCtrl.Volume = p.volumeLevel
		speaker.Unlock()
	}
}

func (p *Player) IsPlaying() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.playing
}

func (p *Player) IsPaused() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.paused
}

func (p *Player) GetVolume() float64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.volumeLevel
}

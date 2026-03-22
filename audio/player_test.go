package audio

import (
	"testing"
)

func TestNewPlayer(t *testing.T) {
	p := NewPlayer()
	if p == nil {
		t.Fatal("NewPlayer() returned nil")
	}
	if p.IsPlaying() {
		t.Error("new player should not be playing")
	}
	if p.IsPaused() {
		t.Error("new player should not be paused")
	}
	if p.GetVolume() != 0 {
		t.Errorf("new player volume = %v, want 0", p.GetVolume())
	}
}

func TestVolumeControl(t *testing.T) {
	p := NewPlayer()

	p.VolumeUp()
	if p.GetVolume() != 0.5 {
		t.Errorf("after VolumeUp() volume = %v, want 0.5", p.GetVolume())
	}

	p.VolumeDown()
	if p.GetVolume() != 0 {
		t.Errorf("after VolumeDown() volume = %v, want 0", p.GetVolume())
	}
}

func TestVolumeBounds(t *testing.T) {
	p := NewPlayer()

	// Volume up to max
	for i := 0; i < 20; i++ {
		p.VolumeUp()
	}
	if p.GetVolume() != 2 {
		t.Errorf("max volume = %v, want 2", p.GetVolume())
	}

	// Volume down to min
	for i := 0; i < 30; i++ {
		p.VolumeDown()
	}
	if p.GetVolume() != -5 {
		t.Errorf("min volume = %v, want -5", p.GetVolume())
	}
}

func TestStopWhenNotPlaying(t *testing.T) {
	p := NewPlayer()
	// Should not panic
	p.Stop()
	if p.IsPlaying() {
		t.Error("player should not be playing after Stop()")
	}
}

func TestTogglePauseWhenNotPlaying(t *testing.T) {
	p := NewPlayer()
	// Should not panic
	p.TogglePause()
	if p.IsPaused() {
		t.Error("player should not be paused when not playing")
	}
}

func TestPlayInvalidURL(t *testing.T) {
	p := NewPlayer()
	err := p.Play("http://invalid.invalid.invalid:99999/stream")
	if err == nil {
		t.Error("Play() with invalid URL should return error")
	}
	if p.IsPlaying() {
		t.Error("player should not be playing after failed Play()")
	}
}

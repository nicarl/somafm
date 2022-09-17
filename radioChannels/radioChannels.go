package radioChannels

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type somafmResponse struct {
	Channels []rawRadioChan
}

type rawRadioChan struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Dj          string `json:"dj"`
	Genre       string `json:"genre"`
	Playlists   []playlist
}
type playlist struct {
	Url     string `json:"url"`
	Format  string `json:"format"`
	Quality string `json:"quality"`
}

type RadioChan struct {
	Id          string
	Title       string
	Description string
	Dj          string
	Genre       string
	StreamURL   string
}

func findMP3Playlist(radioCh rawRadioChan) (string, error) {
	var mp3Playlist string
	for i := range radioCh.Playlists {
		if radioCh.Playlists[i].Format == "mp3" {
			mp3Playlist = radioCh.Playlists[i].Url
			break
		}
	}
	if mp3Playlist == "" {
		return mp3Playlist, fmt.Errorf("Could not find mp3 playlist for channel")
	}

	return mp3Playlist, nil
}

func getStreamURL(radioCh rawRadioChan) (string, error) {
	var streamUrl string
	playlist, err := findMP3Playlist(radioCh)
	if err != nil {
		return streamUrl, err
	}
	resp, err := http.Get(playlist)
	if err != nil {
		return streamUrl, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return streamUrl, err
	}
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "File") {
			split := strings.Split(line, "=")
			if len(split) == 2 {
				return split[1], nil
			}
		}
	}

	return streamUrl, fmt.Errorf("Could not find stream url")
}

func convertRawChannels(channels []rawRadioChan) ([]RadioChan, error) {
	type empty struct{}

	n := len(channels)

	radioChannels := make([]RadioChan, n)
	sem := make(chan empty, n)
	errs := make(chan error, n)

	for i, ch := range channels {
		go func(i int, ch rawRadioChan) {
			streamUrl, err := getStreamURL(ch)
			if err != nil {
				errs <- err
				return
			}
			radioChannels[i] = RadioChan{
				Id:          ch.Id,
				Title:       ch.Title,
				Description: ch.Description,
				Dj:          ch.Dj,
				Genre:       ch.Genre,
				StreamURL:   streamUrl,
			}
			sem <- empty{}
		}(i, ch)
	}
	for i := 0; i < n; i++ {
		<-sem
	}
	select {
	case err := <-errs:
		return nil, err
	default:
		return radioChannels, nil
	}
}

func getRawChannels() ([]rawRadioChan, error) {
	resp, err := http.Get("https://somafm.com/channels.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)

	var data somafmResponse

	if err = d.Decode(&data); err != nil {
		return nil, err
	}
	if len(data.Channels) == 0 {
		return nil, fmt.Errorf("Did not find channels")
	}
	return data.Channels, nil
}

func GetChannels() ([]RadioChan, error) {
	rawChannels, err := getRawChannels()
	if err != nil {
		return nil, err
	}

	return convertRawChannels(rawChannels)
}

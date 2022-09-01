package radioChannels

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type somafmResponse struct {
	Channels []RadioCh
}

type RadioCh struct {
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

func findMP3Playlist(radioCh RadioCh) (string, error) {
	var mp3Playlist string
	for i := range radioCh.Playlists {
		if radioCh.Playlists[i].Format == "mp3" {
			mp3Playlist = radioCh.Playlists[i].Url
			break
		}
	}
	if &mp3Playlist == nil {
		return mp3Playlist, fmt.Errorf("Could not find mp3 playlist for channel")
	}

	return mp3Playlist, nil
}

func GetStreamUrl(radioCh RadioCh) (string, error) {
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

func GetChannels() ([]RadioCh, error) {
	resp, err := http.Get("https://somafm.com/channels.json")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)

	var data somafmResponse

	if err = d.Decode(&data); err != nil {
		return nil, err
	}

	return data.Channels, nil
}

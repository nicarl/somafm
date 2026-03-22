package radiochannels

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type somafmResponse struct {
	Channels []rawRadioChan
}

type rawRadioChan struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Dj          string     `json:"dj"`
	Genre       string     `json:"genre"`
	LastPlaying string     `json:"lastPlaying"`
	Playlists   []playlist `json:"playlists"`
}

type playlist struct {
	URL     string `json:"url"`
	Format  string `json:"format"`
	Quality string `json:"quality"`
}

type Song struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
}

type songsResponse struct {
	Songs []Song `json:"songs"`
}

type RadioChan struct {
	ID          string
	Title       string
	Description string
	Dj          string
	Genre       string
	StreamURL   string
	Songs       []Song
}

func (radioChan RadioChan) GetDetails() string {
	var details string
	if radioChan.Dj != "" {
		details = fmt.Sprintf("%s\n\nDJ: %s\nGenre: %s", radioChan.Description, radioChan.Dj, radioChan.Genre)
	} else {
		details = fmt.Sprintf("%s\n\nGenre: %s", radioChan.Description, radioChan.Genre)
	}

	if len(radioChan.Songs) > 0 {
		details += "\n\nRecent Tracks:"
		limit := 5
		if len(radioChan.Songs) < limit {
			limit = len(radioChan.Songs)
		}
		for _, s := range radioChan.Songs[:limit] {
			details += fmt.Sprintf("\n  %s - %s", s.Artist, s.Title)
		}
	}

	return details
}

func findMP3Playlist(radioCh rawRadioChan) (string, error) {
	for _, pl := range radioCh.Playlists {
		if pl.Format == "mp3" {
			return pl.URL, nil
		}
	}
	return "", fmt.Errorf("could not find mp3 playlist for channel %s", radioCh.ID)
}

func ParsePLS(body string) (string, error) {
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "File") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}
	return "", fmt.Errorf("no stream URL found in PLS data")
}

func getStreamURL(ctx context.Context, radioCh rawRadioChan) (string, error) {
	playlistURL, err := findMP3Playlist(radioCh)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", playlistURL, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetching playlist: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading playlist: %w", err)
	}

	return ParsePLS(string(body))
}

func resolveChannels(ctx context.Context, channels []rawRadioChan) []RadioChan {
	var mu sync.Mutex
	var wg sync.WaitGroup
	result := make([]RadioChan, 0, len(channels))

	for _, ch := range channels {
		wg.Add(1)
		go func(ch rawRadioChan) {
			defer wg.Done()

			streamURL, err := getStreamURL(ctx, ch)
			if err != nil {
				return
			}

			mu.Lock()
			result = append(result, RadioChan{
				ID:          ch.ID,
				Title:       ch.Title,
				Description: ch.Description,
				Dj:          ch.Dj,
				Genre:       ch.Genre,
				StreamURL:   streamURL,
			})
			mu.Unlock()
		}(ch)
	}

	wg.Wait()

	sort.Slice(result, func(i, j int) bool {
		return result[i].Title < result[j].Title
	})

	return result
}

func getRawChannels(ctx context.Context) ([]rawRadioChan, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://somafm.com/channels.json", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data somafmResponse
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if len(data.Channels) == 0 {
		return nil, fmt.Errorf("no channels found")
	}
	return data.Channels, nil
}

func FetchSongs(channelID string) ([]Song, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("https://somafm.com/songs/%s.json", channelID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data songsResponse
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Songs, nil
}

func GetChannels() ([]RadioChan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rawChannels, err := getRawChannels(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching channels: %w", err)
	}

	channels := resolveChannels(ctx, rawChannels)
	if len(channels) == 0 {
		return nil, fmt.Errorf("failed to resolve any channel stream URLs")
	}

	return channels, nil
}

package radiochannels

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFindMP3Playlist(t *testing.T) {
	tests := []struct {
		name     string
		channel  rawRadioChan
		wantURL  string
		wantErr  bool
	}{
		{
			name: "finds mp3 playlist",
			channel: rawRadioChan{
				ID: "test",
				Playlists: []playlist{
					{URL: "https://example.com/aac.pls", Format: "aac", Quality: "high"},
					{URL: "https://example.com/mp3.pls", Format: "mp3", Quality: "high"},
				},
			},
			wantURL: "https://example.com/mp3.pls",
			wantErr: false,
		},
		{
			name: "no mp3 playlist",
			channel: rawRadioChan{
				ID: "test",
				Playlists: []playlist{
					{URL: "https://example.com/aac.pls", Format: "aac", Quality: "high"},
				},
			},
			wantURL: "",
			wantErr: true,
		},
		{
			name: "empty playlists",
			channel: rawRadioChan{
				ID:        "test",
				Playlists: []playlist{},
			},
			wantURL: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findMP3Playlist(tt.channel)
			if (err != nil) != tt.wantErr {
				t.Errorf("findMP3Playlist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantURL {
				t.Errorf("findMP3Playlist() = %v, want %v", got, tt.wantURL)
			}
		})
	}
}

func TestParsePLS(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    string
		wantErr bool
	}{
		{
			name: "standard PLS with numbered entries",
			body: `[playlist]
numberofentries=3
File1=https://ice1.somafm.com/groovesalad-128-mp3
Title1=SomaFM: Groove Salad (#1)
Length1=-1
File2=https://ice2.somafm.com/groovesalad-128-mp3
Title2=SomaFM: Groove Salad (#2)
Length2=-1
Version=2`,
			want:    "https://ice1.somafm.com/groovesalad-128-mp3",
			wantErr: false,
		},
		{
			name: "PLS with unnumbered File entry",
			body: `[playlist]
File=https://stream.example.com/radio.mp3
Title=Example Radio
Length=-1`,
			want:    "https://stream.example.com/radio.mp3",
			wantErr: false,
		},
		{
			name:    "empty PLS",
			body:    "[playlist]\nnumberofentries=0\nVersion=2",
			want:    "",
			wantErr: true,
		},
		{
			name:    "PLS with no File entries",
			body:    "[playlist]\nTitle1=Something\nLength1=-1",
			want:    "",
			wantErr: true,
		},
		{
			name: "PLS with URL containing equals sign",
			body: `[playlist]
File1=https://stream.example.com/radio?quality=high
Title1=Radio
Length1=-1`,
			want:    "https://stream.example.com/radio?quality=high",
			wantErr: false,
		},
		{
			name: "PLS with whitespace",
			body: `[playlist]
  File1=https://stream.example.com/radio.mp3
Title1=Radio`,
			want:    "https://stream.example.com/radio.mp3",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePLS(tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePLS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParsePLS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDetails(t *testing.T) {
	t.Run("with DJ", func(t *testing.T) {
		ch := RadioChan{
			Description: "Chill vibes",
			Dj:          "DJ Cool",
			Genre:       "ambient",
		}
		got := ch.GetDetails()
		want := "Chill vibes\n\nDJ: DJ Cool\nGenre: ambient"
		if got != want {
			t.Errorf("GetDetails() = %v, want %v", got, want)
		}
	})

	t.Run("without DJ", func(t *testing.T) {
		ch := RadioChan{
			Description: "Chill vibes",
			Genre:       "ambient",
		}
		got := ch.GetDetails()
		want := "Chill vibes\n\nGenre: ambient"
		if got != want {
			t.Errorf("GetDetails() = %v, want %v", got, want)
		}
	})

	t.Run("with songs", func(t *testing.T) {
		ch := RadioChan{
			Description: "Chill vibes",
			Genre:       "ambient",
			Songs: []Song{
				{Artist: "Artist1", Title: "Song1"},
				{Artist: "Artist2", Title: "Song2"},
			},
		}
		got := ch.GetDetails()
		want := "Chill vibes\n\nGenre: ambient\n\nRecent Tracks:\n  Artist1 - Song1\n  Artist2 - Song2"
		if got != want {
			t.Errorf("GetDetails() = %q, want %q", got, want)
		}
	})

	t.Run("with more than 5 songs shows only 5", func(t *testing.T) {
		songs := make([]Song, 7)
		for i := range songs {
			songs[i] = Song{Artist: fmt.Sprintf("Artist%d", i+1), Title: fmt.Sprintf("Song%d", i+1)}
		}
		ch := RadioChan{
			Description: "Chill vibes",
			Genre:       "ambient",
			Songs:       songs,
		}
		got := ch.GetDetails()
		if strings.Contains(got, "Artist6") {
			t.Errorf("GetDetails() should only show 5 songs, got: %s", got)
		}
		if !strings.Contains(got, "Artist5") {
			t.Errorf("GetDetails() should show 5th song, got: %s", got)
		}
	})
}

func TestFetchSongs(t *testing.T) {
	songsJSON := `{"id":"test","songs":[{"title":"Track1","artist":"Band1","album":"Album1"},{"title":"Track2","artist":"Band2","album":"Album2"}]}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, songsJSON)
	}))
	defer server.Close()

	// We can't easily test FetchSongs with a custom URL since it's hardcoded,
	// but we test the JSON parsing logic directly.
	t.Run("json parsing", func(t *testing.T) {
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatalf("GET error = %v", err)
		}
		defer resp.Body.Close()

		var data songsResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			t.Fatalf("Decode error = %v", err)
		}

		if len(data.Songs) != 2 {
			t.Errorf("got %d songs, want 2", len(data.Songs))
		}
		if data.Songs[0].Artist != "Band1" {
			t.Errorf("first song artist = %v, want Band1", data.Songs[0].Artist)
		}
		if data.Songs[0].Title != "Track1" {
			t.Errorf("first song title = %v, want Track1", data.Songs[0].Title)
		}
	})
}

func TestGetStreamURL(t *testing.T) {
	plsContent := `[playlist]
numberofentries=2
File1=https://ice1.example.com/stream-128-mp3
Title1=Test Stream (#1)
Length1=-1
File2=https://ice2.example.com/stream-128-mp3
Title2=Test Stream (#2)
Length2=-1
Version=2`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, plsContent)
	}))
	defer server.Close()

	ch := rawRadioChan{
		ID: "test",
		Playlists: []playlist{
			{URL: server.URL + "/test.pls", Format: "mp3", Quality: "high"},
		},
	}

	ctx := context.Background()
	got, err := getStreamURL(ctx, ch)
	if err != nil {
		t.Fatalf("getStreamURL() error = %v", err)
	}
	want := "https://ice1.example.com/stream-128-mp3"
	if got != want {
		t.Errorf("getStreamURL() = %v, want %v", got, want)
	}
}

func TestGetRawChannels(t *testing.T) {
	channels := somafmResponse{
		Channels: []rawRadioChan{
			{
				ID:    "groove",
				Title: "Groove Salad",
				Playlists: []playlist{
					{URL: "https://example.com/groove.pls", Format: "mp3", Quality: "high"},
				},
			},
			{
				ID:    "drone",
				Title: "Drone Zone",
				Playlists: []playlist{
					{URL: "https://example.com/drone.pls", Format: "mp3", Quality: "high"},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(channels); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	// We can't easily test getRawChannels with a custom URL since it's hardcoded,
	// but we test the JSON parsing logic by calling the decoder directly.
	t.Run("json parsing", func(t *testing.T) {
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatalf("GET error = %v", err)
		}
		defer resp.Body.Close()

		var data somafmResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			t.Fatalf("Decode error = %v", err)
		}

		if len(data.Channels) != 2 {
			t.Errorf("got %d channels, want 2", len(data.Channels))
		}
		if data.Channels[0].ID != "groove" {
			t.Errorf("first channel id = %v, want groove", data.Channels[0].ID)
		}
	})
}

func TestResolveChannels(t *testing.T) {
	plsContent := `[playlist]
File1=https://ice1.example.com/stream-128-mp3
Title1=Test
Length1=-1
Version=2`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, plsContent)
	}))
	defer server.Close()

	rawChannels := []rawRadioChan{
		{
			ID:    "beta",
			Title: "Beta Channel",
			Playlists: []playlist{
				{URL: server.URL + "/beta.pls", Format: "mp3", Quality: "high"},
			},
		},
		{
			ID:    "alpha",
			Title: "Alpha Channel",
			Playlists: []playlist{
				{URL: server.URL + "/alpha.pls", Format: "mp3", Quality: "high"},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	channels := resolveChannels(ctx, rawChannels)
	if len(channels) != 2 {
		t.Fatalf("got %d channels, want 2", len(channels))
	}

	// Should be sorted alphabetically by title
	if channels[0].Title != "Alpha Channel" {
		t.Errorf("first channel = %v, want Alpha Channel", channels[0].Title)
	}
	if channels[1].Title != "Beta Channel" {
		t.Errorf("second channel = %v, want Beta Channel", channels[1].Title)
	}
}

func TestResolveChannelsPartialFailure(t *testing.T) {
	plsContent := `[playlist]
File1=https://ice1.example.com/stream-128-mp3
Title1=Test
Length1=-1
Version=2`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/fail.pls" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, plsContent)
	}))
	defer server.Close()

	rawChannels := []rawRadioChan{
		{
			ID:    "good",
			Title: "Good Channel",
			Playlists: []playlist{
				{URL: server.URL + "/good.pls", Format: "mp3", Quality: "high"},
			},
		},
		{
			ID:    "bad",
			Title: "Bad Channel",
			Playlists: []playlist{
				{URL: server.URL + "/fail.pls", Format: "mp3", Quality: "high"},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	channels := resolveChannels(ctx, rawChannels)

	// Should still return the successful channel
	if len(channels) != 1 {
		t.Fatalf("got %d channels, want 1", len(channels))
	}
	if channels[0].Title != "Good Channel" {
		t.Errorf("channel = %v, want Good Channel", channels[0].Title)
	}
}

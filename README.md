# somafm

> TUI application to listen to SomaFM stations

If you enjoy listening to [SomaFM](https://somafm.com/) please support SomaFM by donating.

## Features

- Browse and play all SomaFM channels
- Pause/resume playback
- Volume control
- Filter channels by name or genre
- Now-playing indicator
- Cross-platform (Linux, macOS, Windows)

## Installation

### Homebrew (macOS & Linux)

```sh
brew tap nicarl/somafm
brew install somafm
```

### Install manually from release

[Download the latest release for your platform](https://github.com/nicarl/somafm/releases) and place the binary in your `$PATH`.

### Install manually from source

Requires Go 1.24+ and platform audio headers (e.g. `libasound2-dev` on Linux).

```sh
git clone https://github.com/nicarl/somafm.git
cd somafm
go build ./cmd/somafm.go
```

## Usage

```sh
somafm
```

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `enter` | Play selected channel |
| `space` | Pause / resume |
| `+` / `-` | Volume up / down |
| `/` | Filter channels |
| `esc` | Clear filter |
| `j` / `down` | Move cursor down |
| `k` / `up` | Move cursor up |
| `q` | Quit |

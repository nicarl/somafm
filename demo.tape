# VHS tape file for recording somafm demo GIF
# https://github.com/charmbracelet/vhs
#
# Prerequisites:
#   - VHS: go install github.com/charmbracelet/vhs@latest
#   - ffmpeg and ttyd (>= 1.7.2) must be installed
#
# Usage:
#   1. Build the binary:  go build ./cmd/somafm.go
#   2. Record the GIF:    PATH="$PWD:$PATH" vhs demo.tape
#
# The output will be written to docs/demo.gif.

Output docs/demo.gif

Set Theme "Catppuccin Mocha"
Set FontSize 16
Set Width 1200
Set Height 700
Set Padding 20
Set Shell bash

# Hide the initial loading phase
Hide
Type "somafm"
Enter
Sleep 10s
Show

# App is now loaded - show the channel list
Sleep 2s

# Browse channels - scroll down through list
Down
Sleep 600ms
Down
Sleep 600ms
Down
Sleep 600ms
Down
Sleep 600ms
Down
Sleep 600ms
Down
Sleep 600ms

# Scroll back up a bit
Up
Sleep 600ms
Up
Sleep 600ms

# Select and play a channel
Enter
Sleep 3s

# Adjust volume up
Type "+"
Sleep 500ms
Type "+"
Sleep 500ms
Type "+"
Sleep 500ms

# Adjust volume down
Type "-"
Sleep 500ms
Type "-"
Sleep 500ms

# Browse while playing
Down
Sleep 600ms
Down
Sleep 600ms
Down
Sleep 600ms

# Switch to details panel
Tab
Sleep 1.5s

# Scroll details
Down
Sleep 500ms
Down
Sleep 500ms

# Switch back to channel list
Tab
Sleep 1s

# Use the filter
Type "/"
Sleep 500ms
Type "groove"
Sleep 2s

# Clear filter
Escape
Sleep 1.5s

# Pause playback
Space
Sleep 2s

# Resume playback
Space
Sleep 2s

# Quit
Type "q"
Sleep 500ms

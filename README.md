# LUDO-GO

A TUI music player which can stream music as per commands. The searching and streaming is done from Piped Api. The VLC player is used for playback.

## Features

- Search and play audio
- Build a song queue
- Remove or skip songs in queue
- Start radio of related songs for the given query

## Installation

1. Make sure VLC Player *>=3.0.18* is available on the PATH
2. Download binary, place in a location and add the location to PATH
3. Open terminal and type in `ludo`. Enter help for list of commands.

### PreRequisite

1. VLC Player
The player is build on libvlc bindings to [libvlc v3.0.18](https://www.nuget.org/packages/VideoLAN.LibVLC.Windows/3.0.18). <br>
VLC Player *>=3.0.1* should be installed and be added to the path.

## Build

1. Follow instructions on [libvlc-go](https://github.com/adrg/libvlc-go) to setup libvlc
2. Add the following Go Env variables

	`CGO_CFLAGS="-I$PATH_TO_LIBVLC/include"` <br>
	`CGO_LDFLAGS="-L$PATH_TO_LIBVLC"` <br>
	`CGO_ENABLED=1`
3. Build the binaries

	`go build -o bin/ludo.exe main.go`

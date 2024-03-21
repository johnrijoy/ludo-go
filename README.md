# LUDO-GO

A VLC-based music player which can stream music as per commands.

## Dev Testing
1. main.go file contains main logic. Navigate to folder and enter command.
    `go run .`

## PreRequisite
The player is build on bindings to [libvlc v3.0.18](https://www.nuget.org/packages/VideoLAN.LibVLC.Windows/3.0.18). This will be required on the path to play audio.

## Build

`go build -o bin/ludo.exe main.go`

# Video Transcoder
Wrapper around [ffmpeg](https://ffmpeg.org/) to transcode media to Apple's ProRes codec.

This will prepare your videos to be used in [DaVinci Resolve](https://www.richardlackey.com/avc-hevc-transcode-davinci-resolve/) or other editing programs that don't like H.264 / H.265 encoded media. 

[![Go Report Card](https://goreportcard.com/badge/github.com/TrafeX/videotranscoder)](https://goreportcard.com/report/github.com/TrafeX/videotranscoder)

## Installation

You need to have ffmpeg installed on your system.

Then, install the binary with `go get github.com/TrafeX/videotranscoder`.

## Usage

```
videotranscode -source=/path/to/folder/with/source/videos -target=/path/to/place/transcoded/folder
```

Note: In the target folder a new folder will be created that has the same name as the source folder with a `-transcoded` postfix.


### Overwrite existing files
By default existing transcoded video files in the target folder are not overridden. You can change this by passing the `-overwrite` flag.
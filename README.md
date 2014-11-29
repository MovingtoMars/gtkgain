GtkGain
=======

GtkGain is a project written in Go which provides a GTK+3 frontend to vorbisgain, mp3gain, metaflac and libvorbis in order to show, add, and remove replaygain tack and album tags.

## Installation

GtkGain has the following depends on:
`flac libvorbis vorbisgain mp3gain taglib gtk>=3.10`

To install GtkGain the Go way (provided you have Go set up and $GOPATH/bin added to your $PATH):
```bash
$ go get github.com/MovingtoMars/gtkgain
$ go install github.com/MovingtoMars/gtkgain
```

## License

GtkGain is licenses under the GNU GPLv3.0

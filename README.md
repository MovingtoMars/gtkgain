GtkGain
=======

GtkGain is a project written in Go which provides a GTK+3 frontend to vorbisgain, mp3gain, metaflac and libvorbis in order to show, add, and remove replaygain track and album tags.

## Installation

GtkGain depends on:
```
flac libvorbis vorbisgain mp3gain taglib gtk>=3.10
```

Installing the normal way:
-----------------------

1. [Download](https://github.com/MovingtoMars/gtkgain/archive/master.zip) and extract the source
2. Run `$ make`in the directory you extracted
3. Run the binary outputted to `bin/`

Binaries
-------

Prebuilt versions of GtkGain can be found at https://bintray.com/mvtm/generic/GtkGain/view

Installing the Go way
-------------------

To install GtkGain the Go way (provided you have Go set up and $GOPATH/bin added to your $PATH):
```bash
$ go get github.com/MovingtoMars/gtkgain
$ go install github.com/MovingtoMars/gtkgain
```

## License

GtkGain is licensed under the GNU GPLv3.0

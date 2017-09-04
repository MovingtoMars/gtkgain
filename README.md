GtkGain
=======

GtkGain is a project written in Go which provides a GTK+3 frontend to vorbisgain, mp3gain, metaflac and libvorbis in order to show, add, and remove replaygain track and album tags.

## Dependencies

Runtime:
```
libvorbis taglib gtk>=3.10 vorbisgain mp3gain flac
```

Build-only dependencies:
```
git golang
```

Note that development versions of the first three runtime dependencies are required for building.

### Ubuntu

For building:

```bash
$ sudo apt-get install git golang golang-go.tools libgtk-3-dev libtagc0-dev libvorbis-dev mp3gain flac vorbisgain
```

**Note: Ubuntu 14.04 Trusty is the most recent Ubuntu release to include `mp3gain`. Now it is unavailable.**

Make sure the `universe` repo is enabled in Software and Updates.

## Installation

### Building from source (recommended)

1. Install the build dependencies (command for Ubuntu shown above)
2. [Download](https://github.com/MovingtoMars/gtkgain/archive/master.zip) and extract the source
3. `$ cd /path/to/gtkgain/`
4. `$ make`
5. `$ sudo make install`

### Installing the Go way

To install GtkGain the Go way (provided you have Go set up and $GOPATH/bin added to your $PATH):
```bash
$ go get github.com/MovingtoMars/gtkgain/src/gtkgain
$ go install github.com/MovingtoMars/gtkgain/src/gtkgain
```

You can then run GtkGain from your terminal by typing `gtkgain`.

## Screenshots

### Arch Linux

![Screenshot on Arch Linux](http://i.imgur.com/GMFJvmU.png "Arch Linux")

### Ubuntu

![Screenshot on Ubuntu](http://i.imgur.com/vXEbuws.png "Ubuntu")
![Screenshot on Ubuntu](http://i.imgur.com/77MIanW.png "Ubuntu")

## License

GtkGain is licensed under the GNU GPLv3.0

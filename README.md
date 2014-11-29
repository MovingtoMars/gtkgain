GtkGain
=======

GtkGain is a project written in Go which provides a GTK+3 frontend to vorbisgain, mp3gain, metaflac and libvorbis in order to show, add, and remove replaygain track and album tags.

## Dependencies

GtkGain depends on:
```
flac libvorbis vorbisgain mp3gain taglib gtk>=3.10
```

Note that some earlier versions of mp3gain are known to corrupt ID3v2 tags.

### Ubuntu

For building:

```bash
$ sudo apt-get install git golang golang-go.tools libgtk-3-dev libtagc0-dev libvorbis-dev mp3gain flac vorbisgain
```

Make sure the `universe` repo is enabled in Software and Updates.

## Installation

### Building from source (recommended)

1. Install the build dependencies (command for Ubuntu shown above)
2. [Download](https://github.com/MovingtoMars/gtkgain/archive/master.zip) and extract the source
3. Use the `make`command in the directory you extracted
4. Run the binary outputted to `bin/`

### Binaries

Prebuilt versions of GtkGain can be found at https://bintray.com/mvtm/generic/GtkGain/view

### Installing the Go way

To install GtkGain the Go way (provided you have Go set up and $GOPATH/bin added to your $PATH):
```bash
$ go get github.com/MovingtoMars/gtkgain
$ go install github.com/MovingtoMars/gtkgain
```
## Screenshots

### Arch Linux

![Screenshot on Arch Linux](http://i.imgur.com/GMFJvmU.png "Arch Linux")

### Ubuntu

![Screenshot on Ubuntu](http://i.imgur.com/vXEbuws.png "Ubuntu")
![Screenshot on Ubuntu](http://i.imgur.com/77MIanW.png "Ubuntu")

## License

GtkGain is licensed under the GNU GPLv3.0

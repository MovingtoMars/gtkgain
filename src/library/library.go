package library

import (
	"github.com/wtolson/go-taglib"
	
	"io/ioutil"
	"mime"
	"path/filepath"
	"strings"
)

type Library struct {
	albums []*Album
	songLoadReceiver func(s *Song)
	loadFinishReceiver func()
}

// New() returns an initialised instance of Library.
func New() *Library {
	return &Library{albums: make([]*Album, 0)}
}

func (l *Library) SetSongLoadReceiver(rec func(s *Song)) {
	l.songLoadReceiver = rec
}

func (l *Library) SetLoadFinishReceiver(rec func()) {
	l.loadFinishReceiver = rec
}

// ImportFromDir() imports all songs in a directory, recursively.
func (l *Library) ImportFromDir(dir string) {
	l.importFromDir(dir)
	if l.loadFinishReceiver != nil {
		l.loadFinishReceiver()
	}
}

func (l *Library) importFromDir(dir string) {
	ls, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, file := range ls {
		if file.IsDir() {
			l.importFromDir(dir + "/" + file.Name())
		} else {
			l.ImportSong(dir + "/" + file.Name())
		}
	}
}

// GetPathAudioFormat() returns the audio format of the audio file located at path.
func GetPathAudioFormat(path string) AudioFormat {
	mt := mime.TypeByExtension(filepath.Ext(path))
	
	if strings.Contains(mt, "flac") {
		return FLAC
	} else if strings.Contains(mt, "ogg") {
		return OGG_VORBIS // TODO make sure audio is vorbis
	} else if strings.Contains(mt, "mpeg") || strings.Contains(mt, "mp3") {
		return MP3
	}
    
	return UNKNOWN
}

func (l *Library) openAlbum(name, artist string) *Album {
	for _, a := range l.albums {
		if a.name == name {
			return a
		}
	}
	a := &Album {name: name, artist: artist, songs: make([]*Song, 0)}
	l.albums = append(l.albums, a)
	return a
}

// ImportSong() imports a single song into the library.
func (l *Library) ImportSong(path string) {	
	form := GetPathAudioFormat(path)
	if form == UNKNOWN {
		return
	}
	
	tags, err := taglib.Read(path)
	if err != nil {
		return
	}
	
	album := l.openAlbum(tags.Album(), tags.Artist())
	s := &Song {format: form, path: path, title: tags.Title(), track: tags.Track(), album: album}
	s.tgain, _ = s.LoadGain(GAIN_TRACK)
	s.again, _ = s.LoadGain(GAIN_ALBUM)
	album.addSong(s)
	
	tags.Close()
	
	if l.songLoadReceiver != nil {
		l.songLoadReceiver(s)
	}
}

func (l *Library) GetAlbum(name string) *Album {
	for _, a := range l.albums {
		if a.name == name {
			return a
		}
	}
	return nil
}

func (l *Library) Albums() []*Album {
	return l.albums
}

func (l *Library) UntaggedAlbums() []*Album {
	ret := make([]*Album, 0)
	for _, a := range l.albums {
		if !a.tagged {
			ret = append(ret, a)
		}
	}
	return ret
}

func (l *Library) TaggedSongs() []*Song {
	ret := make([]*Song, 0)
	for _, a := range l.albums {
		for _, s := range a.songs {
			if s.Gain(GAIN_TRACK) != "?" || s.Gain(GAIN_ALBUM) != "?" {
				ret = append(ret, s)
			}
		}
	}
	return ret
}

func (l *Library) String() string {
	ret := ""
	for _, a := range l.albums {
		ret += a.String() + "\n"
	}
	return ret[:len(ret) - 1]
}

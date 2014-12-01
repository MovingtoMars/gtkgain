package library

import (
	"errors"
	"sync"
)

type GainType int

const (
	GAIN_TRACK GainType = iota
	GAIN_ALBUM
)

var ErrUnknownGainType = errors.New("unknown gain type")

// TODO locks for gain
type AudioFormat int

const (
	UNKNOWN AudioFormat = iota
	MP3
	OGG_VORBIS
	FLAC
)

var ErrUnknownFormat = errors.New("unknown audio format")

func (v AudioFormat) String() string {
	switch v {
	case UNKNOWN:
		return "UNKNOWN"
	case MP3:
		return "MP3"
	case OGG_VORBIS:
		return "OGG_VORBIS"
	case FLAC:
		return "FLAC"
	default:
		return "ERROR_FORMAT"
	}
}

type Song struct {
	path, title string
	format AudioFormat
	track int
	album *Album
	tgain, again string
	gainLock sync.Mutex
}

func (s *Song) String() string {
	return s.title
}

func (s *Song) Title() string {
	return s.title
}

func (s *Song) Track() int {
	return s.track
}

func (s *Song) AlbumName() string {
	return s.album.name
}

func (s *Song) Path() string {
	return s.path
}

// MT safe
func (s *Song) Gain(t GainType) (ret string) {
	s.gainLock.Lock()
	switch t {
	case GAIN_ALBUM:
		ret = s.again
	case GAIN_TRACK:
		ret = s.tgain
	default:
		ret = ""
	}
	s.gainLock.Unlock()
	return
}

// MT safe
func (s *Song) SetGain(g string, t GainType) {
	s.gainLock.Lock()
	switch t {
	case GAIN_ALBUM:
		s.again = g
	case GAIN_TRACK:
		s.tgain = g
	}
	s.gainLock.Unlock()
}

func (s *Song) LoadGain(t GainType) (string, error) {
	ret := ""
	var err error
	switch s.format {
	case FLAC:
		ret, err = flacGetTagGain(s.path, t)
	case OGG_VORBIS:
		ret, err = vorbisGetTagGain(s.path, t)
	case MP3:
		ret, err = mp3GetTagGain(s.path, t)
	}
	
	if ret == "" {
		return "?", err
	}
	return ret, err
}

func (s *Song) UntagGain(songUpdateReceiver func(*Song)) error {
	var err error
	switch s.format {
	case FLAC:
		err = flacUntagGain([]string {s.path})
	case OGG_VORBIS:
		err = vorbisUntagGain([]string {s.path})
	case UNKNOWN:
		return ErrUnknownFormat
	default:
		return ErrUnknownFormat
	}
	
	s.gainLock.Lock()
	s.tgain, _ = s.LoadGain(GAIN_TRACK)
	s.again, _ = s.LoadGain(GAIN_ALBUM)
	s.gainLock.Unlock()
	songUpdateReceiver(s)
	return err
}

func SongsUntagGain(list []*Song, songUpdateReceiver func(*Song)) error {
	var err error
	flac := make([]string, 0)
	vorbis := make([]string, 0)
	mp3 := make([]string, 0)
	for _, s := range list {
		switch s.format {
		case FLAC:
			flac = append(flac, s.path)
		case OGG_VORBIS:
			vorbis = append(vorbis, s.path)
		case MP3:
			mp3 = append(mp3, s.path)
		default:
			err = ErrUnknownFormat
		}
	}
	
	flacUntagGain(flac)
	vorbisUntagGain(vorbis)
	mp3UntagGain(mp3)
	
	for _, s := range list {
		s.SetGain("?", GAIN_TRACK)
		s.SetGain("?", GAIN_ALBUM)
		s.album.tagged = false
		songUpdateReceiver(s)
	}
	
	return err
}

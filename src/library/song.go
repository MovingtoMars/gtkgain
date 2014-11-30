package library

import (
	"errors"
)

type GainType int

const (
	GAIN_TRACK GainType = iota
	GAIN_ALBUM
)

// TODO locks for gain
type AudioFormat int

const (
	UNKNOWN AudioFormat = iota
	MP3
	OGG_VORBIS
	FLAC
)

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

func (s *Song) TrackGain() string {
	/*if s.tgain == "" {
		return "?"
	}*/
	return s.tgain
}

func (s *Song) AlbumGain() string {
	/*if s.again == "" {
		return "?"
	}*/
	return s.again
}

// TODO merge this and album function into one, with each type being optional
func (s *Song) LoadTrackGain() string {
	ret := ""
	switch s.format {
	case FLAC:
		ret = flacGetTagGain(s.path, "REPLAYGAIN_TRACK_GAIN")
	case OGG_VORBIS:
		ret = vorbisGetTagGain(s.path, "REPLAYGAIN_TRACK_GAIN")
	case MP3:
		ret = mp3GetTagGain(s.path, GAIN_TRACK)
	}
	
	if ret == "" {
		return "?" 
	}
	return ret
}

func (s *Song) LoadAlbumGain() string {
	ret := ""
	switch s.format {
	case FLAC:
		ret = flacGetTagGain(s.path, "REPLAYGAIN_ALBUM_GAIN")
	case OGG_VORBIS:
		ret = vorbisGetTagGain(s.path, "REPLAYGAIN_ALBUM_GAIN")
	case MP3:
		ret = mp3GetTagGain(s.path, GAIN_ALBUM)
	}
	
	if ret == "" {
		return "?" 
	}
	return ret
}

func (s *Song) UntagGain(songUpdateReceiver func(*Song)) error {
	var err error
	switch s.format {
	case FLAC:
		err = flacUntagGain([]string {s.path})
	case OGG_VORBIS:
		err = vorbisUntagGain([]string {s.path})
	case UNKNOWN:
		return errors.New("can't tag unknown/inconsistent formatted album")
	default:
		return errors.New("unknown format type")
	}
	
	s.tgain = s.LoadTrackGain()
	s.again = s.LoadAlbumGain()
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
			err = errors.New("unknown format type")
		}
	}
	
	flacUntagGain(flac)
	vorbisUntagGain(vorbis)
	mp3UntagGain(mp3)
	
	for _, s := range list {
		s.tgain = "?"
		s.again = "?"
		s.album.tagged = false
		songUpdateReceiver(s)
	}
	
	return err
}

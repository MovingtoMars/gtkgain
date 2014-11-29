package library

import (
	"strconv"
	"errors"
)

type Album struct {
	songs []*Song
	name, artist string
	format AudioFormat // UNKNOWN if formats of all songs in the album aren't the same
	tagged bool
}

func (v *Album) addSong(s *Song) {
	if len(v.songs) == 0 {
		v.format = s.format
		v.tagged = true
	} else if v.format != s.format {
		v.format = UNKNOWN
	}
	v.tagged = v.tagged && (s.AlbumGain() != "?")
	v.songs = append(v.songs, s)
}

func (v *Album) GetTrack(trackNo int) *Song {
	for _, s := range v.songs {
		if s.track == trackNo {
			return s
		}
	}
	return nil
}

func (v *Album) GetSongs() []*Song {
	return v.songs
}

func (v *Album) String() string {
	return v.name + " - " + strconv.FormatBool(v.tagged)
}

func (v *Album) TagGain(songUpdateReceiver func(*Song)) error {
	paths := make([]string, 0)
	for _, s := range v.songs {
		paths = append(paths, s.path)
	}
	
	var err error
	
	switch v.format {
	case FLAC:
		err = flacTagGainAlbum(paths)
	case OGG_VORBIS:
		err = vorbisTagGainAlbum(paths)
	case MP3:
		err = mp3TagGainAlbum(paths)
	case UNKNOWN:
		return errors.New("can't tag unknown/inconsistently formatted album")
	default:
		return errors.New("unknown format type")
	}
	
	for _, s := range v.songs {
		s.tgain = s.LoadTrackGain()
		s.again = s.LoadAlbumGain()
		v.tagged = true
		songUpdateReceiver(s)
	}
	return err
}

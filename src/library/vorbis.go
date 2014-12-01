package library

import (
	"github.com/vchimishuk/chub/src/ogg/libvorbis"
	
	"strings"
	"os/exec"
	"log"
)

func vorbisGetTagGain(path string, t GainType) (string, error) {
	var tag string
	switch t {
	case GAIN_ALBUM:
		tag = "REPLAYGAIN_ALBUM_GAIN"
	case GAIN_TRACK:
		tag = "REPLAYGAIN_TRACK_GAIN"
	default:
		log.Println("unknown GainType in call to vorbisGetTagGain")
		return "", ErrUnknownGainType
	}
	
	f, err := libvorbis.New(path)
	if err != nil {
		return "", err
	}
	
	defer f.Close()
	for _, c := range f.Comment().UserComments {
		if len(c) < 500 {
			if strings.HasPrefix(c, tag) {
				return c[len(tag) + 1:], nil
			}
		}
	}
	return "", nil
}

func vorbisTagGain(path []string, t GainType) error {
	switch t {
	case GAIN_ALBUM:
		cmd := exec.Command("vorbisgain", append([]string {"-aq"}, path...)...)
		_, err := cmd.CombinedOutput()
		return err
	case GAIN_TRACK:
		cmd := exec.Command("vorbisgain", append([]string {"-q"}, path...)...)
		_, err := cmd.CombinedOutput()
		return err
	default:
		log.Println("unknown GainType in call to vorbisTagGain")
		return ErrUnknownGainType
	}
}

func vorbisUntagGain(path []string) error {
	cmd := exec.Command("vorbisgain", append([]string {"-cq"}, path...)...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	
	return nil
}

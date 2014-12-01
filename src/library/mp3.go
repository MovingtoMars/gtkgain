package library

import (
	"os/exec"
	"strings"
	"log"
)

const (
	MP3_ASTR = "Recommended \"Album\" dB change: "
	MP3_TSTR = "Recommended \"Track\" dB change: "
)

func mp3GetTagGain(path string, t GainType) (string, error) {
	cmd := exec.Command("mp3gain", "-s", "i", "-s", "c", path)
	b, err := cmd.CombinedOutput()
	if err != nil || string(b) == "" {
		return "", err
	}
	
	result := string(b)
	var identifier string
	switch t {
	case GAIN_ALBUM:
		identifier = MP3_ASTR
	case GAIN_TRACK:
		identifier = MP3_TSTR
	default:
		log.Println("unknown GainType in call to mp3TagGain")
		return "", ErrUnknownGainType
	}
	
	for _, c := range strings.Split(result, "\n") {
		if strings.HasPrefix(c, identifier) {
			return c[len(identifier):], nil
		}
	}
	
	return "", nil
}

func mp3TagGain(path []string, t GainType) error {
	switch t {
	case GAIN_ALBUM:
		cmd := exec.Command("mp3gain", append([]string {"-a", "-q", "-c"}, path...)...)
		_, err := cmd.CombinedOutput()
		return err
	case GAIN_TRACK:
		cmd := exec.Command("mp3gain", append([]string {"-e", "-q", "-c"}, path...)...)
		_, err := cmd.CombinedOutput()
		return err
	default:
		log.Println("unknown GainType in call to mp3TagGain")
		return ErrUnknownGainType
	}
}

func mp3UntagGain(path []string) error {
	cmd := exec.Command("mp3gain", append([]string {"-s", "i", "-s", "d"}, path...)...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	
	return nil
}

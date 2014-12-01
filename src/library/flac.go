package library

import (
	"os/exec"
	"log"
)

func flacGetTagGain(path string, t GainType) (string, error) {
	var tag string
	switch t {
	case GAIN_ALBUM:
		tag = "REPLAYGAIN_ALBUM_GAIN"
	case GAIN_TRACK:
		tag = "REPLAYGAIN_TRACK_GAIN"
	default:
		log.Println("unknown GainType in call to flacGetTagGain")
		return "", ErrUnknownGainType
	}
	
	cmd := exec.Command("metaflac", "--show-tag=" + tag, path)
	b, err := cmd.CombinedOutput()
	if err != nil || string(b) == "" {
		return "", err
	}
	
	return string(b)[len(tag) + 1:len(b) - 1], nil
}

func flacTagGain(path []string, t GainType) error {
	switch t {
	case GAIN_ALBUM:
		cmd := exec.Command("metaflac", append([]string {"--add-replay-gain"}, path...)...)
		_, err := cmd.CombinedOutput()
		return err
	case GAIN_TRACK:
		var aerr error
		for _, p := range path {
			cmd := exec.Command("metaflac", "--add-replay-gain", p)
			_, err := cmd.CombinedOutput()
			if err != nil {
				aerr = err
			}
		}
		return aerr
	default:
		log.Println("unknown GainType in call to flacTagGain")
		return ErrUnknownGainType
	}
}

func flacUntagGain(path []string) error {
	cmd := exec.Command("metaflac", append([]string {"--remove-replay-gain"}, path...)...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	
	return nil
}

package library

import (
	"os/exec"
)

func flacGetTagGain(path, tag string) string {
	cmd := exec.Command("metaflac", "--show-tag=" + tag, path)
	b, err := cmd.CombinedOutput()
	if err != nil || string(b) == "" {
		return ""
	}
	
	return string(b)[len(tag) + 1:len(b) - 1]
}

func flacTagGainAlbum(path []string) error {
	cmd := exec.Command("metaflac", append([]string {"--add-replay-gain"}, path...)...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	
	return nil
}

func flacUntagGain(path []string) error {
	cmd := exec.Command("metaflac", append([]string {"--remove-replay-gain"}, path...)...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	
	return nil
}

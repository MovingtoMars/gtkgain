package library

import (
	"github.com/vchimishuk/chub/src/ogg/libvorbis"
	
	"strings"
	"os/exec"
)

func vorbisGetTagGain(path, tag string) string {
	f, _ := libvorbis.New(path)
	defer f.Close()
	for _, c := range f.Comment().UserComments {
		if len(c) < 500 {
			if strings.HasPrefix(c, tag) {
				return c[len(tag) + 1:]
			}
		}
	}
	return ""
}

func vorbisTagGainAlbum(path []string) error {
	cmd := exec.Command("vorbisgain", append([]string {"-aq"}, path...)...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	
	return nil
}

func vorbisUntagGain(path []string) error {
	cmd := exec.Command("vorbisgain", append([]string {"-cq"}, path...)...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	
	return nil
}

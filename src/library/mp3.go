package library

import (
	"os/exec"
	"strings"
)

const (
	MP3_ASTR = "Recommended \"Album\" dB change: "
	MP3_TSTR = "Recommended \"Track\" dB change: "
)

func mp3GetTagGain(path string, t GainType) string {
	cmd := exec.Command("mp3gain", "-s", "i", "-s", "c", path)
	b, err := cmd.CombinedOutput()
	if err != nil || string(b) == "" {
		return ""
	}
	
	str := string(b)
	
	if t == GAIN_ALBUM {
		for _, c := range strings.Split(str, "\n") {
			if strings.HasPrefix(c, MP3_ASTR) {
				return c[len(MP3_ASTR):]
			}
		}
	} else if t == GAIN_TRACK {
		for _, c := range strings.Split(str, "\n") {
			if strings.HasPrefix(c, MP3_TSTR) {
				return c[len(MP3_TSTR):]
			}
		}
	}
	
	return ""
}

/*func mp3LoadTagGainMultiple(songs []*Song) {
	fmt.Println("called")
	paths := make([]string, len(songs))
	for i, s := range songs {
		paths[i] = s.path
	}
	
	cmd := exec.Command("mp3gain", append([]string {"-s", "i", "-s", "c"}, paths...)...)
	b, err := cmd.CombinedOutput()
	if err != nil || string(b) == "" {
		return
	}
	
	str := string(b)
	fmt.Println(str)
	
	i := 0
	for _, line := range strings.Split(str, "\n") {
		if line == "" {
			i++
			if i > len(songs) {
				return
			}
		} else if strings.HasPrefix(line, MP3_ASTR) {
			songs[i].again = line[len(MP3_ASTR):]
		} else if strings.HasPrefix(line, MP3_TSTR) {
			songs[i].tgain = line[len(MP3_TSTR):]
		}
	}
}*/

func mp3TagGainAlbum(path []string) error {
	cmd := exec.Command("mp3gain", append([]string {"-aqc"}, path...)...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	
	return nil
}

func mp3UntagGain(path []string) error {
	cmd := exec.Command("mp3gain", append([]string {"-s", "i", "-s", "d"}, path...)...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	
	return nil
}

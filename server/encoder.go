package server

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/google/uuid"
)

func EncodeEpisode(episode *Episode) (string, error) {
	return glueAudioFiles(episode.FilePaths)
}

func glueAudioFiles(filepaths []string) (string, error) {
	output := path.Join(os.TempDir(), uuid.New().String()+".mp3")
	args := []string{"-i", "concat:" + strings.Join(filepaths, "|"), "-acodec", "copy", output}
	cmd := exec.Command("ffmpeg", args...)
	fmt.Println(cmd)
	if err := exec.Command("ffmpeg", args...).Run(); err != nil {
		return "", err
	}
	return output, nil
}

package undercast

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

type FFMpegMediaConverter struct {
	workDir string
}

func NewFFMpegMediaConverter(workDir string) (*FFMpegMediaConverter, error) {
	stat, err := os.Stat(workDir)
	if err != nil {
		err = os.MkdirAll(workDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	} else if !stat.IsDir() {
		return nil, fmt.Errorf("%s already exists and is not a directory", workDir)
	}
	return &FFMpegMediaConverter{workDir: workDir}, nil
}

func (conv *FFMpegMediaConverter) Concatenate(filepaths []string, filename string, format string) (string, error) {
	resultFilepath := path.Join(conv.workDir, filename)
	args := []string{"-y", "-i", "concat:" + strings.Join(filepaths, "|"), "-acodec", format, resultFilepath}
	cmd := exec.Command("ffmpeg", args...)
	log.Printf("Running %s\n", cmd)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return resultFilepath, nil
}

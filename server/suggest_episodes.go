package server

import (
	"path"
	"sort"
	"strings"
)

func suggestEpisodes(name string, filepaths []string) (episodes []Episode) {
	dirToFilesMap := filepathsToMap(filepaths)
	if len(dirToFilesMap) == 1 {
		sort.Strings(filepaths)
		episodes = append(episodes, Episode{Name: name, FilePaths: filepaths})
		return
	}
	for dirname, filepaths := range dirToFilesMap {
		sort.Strings(filepaths)
		episodes = append(episodes, Episode{Name: dirname, FilePaths: filepaths})
	}
	return
}

func filepathsToMap(filepaths []string) map[string][]string {
	result := map[string][]string{}
	for _, f := range filepaths {
		dirname, _ := path.Split(f)
		dirname = strings.Trim(dirname, "/")
		bits := strings.Split(dirname, "/")
		dirname = bits[len(bits)-1]
		if result[dirname] == nil {
			result[dirname] = []string{}
		}
		result[dirname] = append(result[dirname], f)
	}
	return result
}

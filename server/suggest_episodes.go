package server

import (
	"path"
	"sort"
	"strings"
)


func suggestEpisodes(name string, filenames []string) (episodes []Episode) {
	dirToFilesMap := filenamesToMap(filenames)
	if len(dirToFilesMap) == 1 {
		sort.Strings(filenames)
		episodes = append(episodes, Episode{Name: name, FileNames: filenames})
		return
	}
	for dirname, filenames := range dirToFilesMap {
		sort.Strings(filenames)
		episodes = append(episodes, Episode{Name: dirname, FileNames: filenames})
	}
	return
}

func filenamesToMap(filenames []string) map[string][]string {
	result := map[string][]string{}
	for _, f := range filenames {
		dirname, _ := path.Split(f)
		dirname = strings.Trim(dirname, "/")
		if result[dirname] == nil {
			result[dirname] = []string{}
		}
		result[dirname] = append(result[dirname], f)
	}
	return result
}

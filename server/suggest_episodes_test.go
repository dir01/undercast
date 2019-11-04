package server

import (
	"reflect"
	"testing"
)

func TestFlatList(t *testing.T) {
	actual := suggestEpisodes("Around world in 80 days", []string{
		"around_world_in_80_days_03_verne_64kb.mp3",
		"around_world_in_80_days_01_verne_64kb.mp3",
		"around_world_in_80_days_02_verne_64kb.mp3",
	})
	expected := []Episode{
		Episode{
			Name: "Around world in 80 days",
			FileNames: []string{
				"around_world_in_80_days_01_verne_64kb.mp3",
				"around_world_in_80_days_02_verne_64kb.mp3",
				"around_world_in_80_days_03_verne_64kb.mp3",
			},
		},
	}
	assertDeepEquals(t, expected, actual)
}

func TestSubDirs(t *testing.T) {
	actual := suggestEpisodes("Oscar Wilde", []string{
		"The Picture of Dorian Gray/chapter2.mp3",
		"The Importance of Being Earnest/02 - act 1.mp3",
		"The Importance of Being Earnest/01 - act 1.mp3",
		"The Picture of Dorian Gray/chapter1.mp3",
	})
	expected := []Episode{
		Episode{Name: "The Picture of Dorian Gray", FileNames: []string{
			"The Picture of Dorian Gray/chapter1.mp3",
			"The Picture of Dorian Gray/chapter2.mp3",
		}},
		Episode{Name: "The Importance of Being Earnest", FileNames: []string{
			"The Importance of Being Earnest/01 - act 1.mp3",
			"The Importance of Being Earnest/02 - act 1.mp3",
		}},
	}
	assertDeepEquals(t, expected, actual)
}

func assertDeepEquals(t *testing.T, expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Deep equality test failed.\nGot %#v\nWant %#v", actual, expected)
	}
}

package backup

import (
	"strings"
	"testing"

	"go.mlcdf.fr/sc-backup/internal/backend/mock"
	"go.mlcdf.fr/sc-backup/internal/domain"
)

func TestValidateUser(t *testing.T) {
	username := "username-that-does-not-exists"
	err := validateUser(username)
	if err == nil {
		t.Errorf("username %s should not exist", username)
	}

	username = "mlcdf"
	err = validateUser(username)
	if err != nil {
		t.Errorf("username %s should exist", username)
	}
}

func TestMakeListURL(t *testing.T) {
	testCases := []struct {
		url      string
		index    int8
		expected string
	}{
		{
			url:      "https://www.senscritique.com/liste/Vu_au_cinema/363578",
			index:    1,
			expected: "https://www.senscritique.com/liste/Vu_au_cinema/363578/page-1",
		},
		{
			url:      "https://www.senscritique.com/liste/Vu_au_cinema/363578#page-1/",
			index:    1,
			expected: "https://www.senscritique.com/liste/Vu_au_cinema/363578#page-1",
		},
		{
			url:      "https://www.senscritique.com/liste/Vu_au_cinema/363578",
			index:    3,
			expected: "https://www.senscritique.com/liste/Vu_au_cinema/363578/page-3",
		},
		{
			url:      "https://www.senscritique.com/liste/Vu_au_cinema/363578#page-1/",
			index:    3,
			expected: "https://www.senscritique.com/liste/Vu_au_cinema/363578#page-3",
		},
		{
			url:      "https://www.senscritique.com/liste/Vu_au_cinema/363578#page-1",
			index:    4,
			expected: "https://www.senscritique.com/liste/Vu_au_cinema/363578#page-4",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.url, func(t *testing.T) {
			if result := makeListURL(tC.url, int(tC.index)); result != tC.expected {
				t.Errorf("expected %s, got %s", tC.expected, result)
			}
		})
	}
}

func TestBackupList(t *testing.T) {
	back := mock.NewBackend()
	List("https://www.senscritique.com/liste/Vu_au_cinema/363578", back)

	stuff := back.Data["vu-au-cinema"]
	if stuff == nil {
		t.Errorf("slug vu-au-cinema not found")
	}

	list, ok := stuff.(*domain.List)
	if ok == false {
		t.Errorf("cast back.Data into domain.List")
	}

	if expectedDescription := "Depuis le 1er janvier 2014."; list.Description != expectedDescription {
		t.Errorf("expected description '%s', got '%s'", expectedDescription, list.Description)
	}

	if l := len(list.Entries); l < 100 {
		t.Errorf("too few entries: %d", l)
	}

	entry := list.Entries[0]
	if entry.ID == "" {
		t.Errorf("entry.ID cannot be empty %v", entry)
	}

	if expected := true; entry.Favorite != expected {
		t.Errorf("expected: %t, got: %t for %s", expected, entry.Favorite, entry.Title)
	}

	if len(entry.Authors) == 0 {
		t.Errorf("entry.Authors cannot be empty %v", entry)
	}

	if len(entry.Genre) > 4 {
		t.Errorf("entry.Genre is too long %v", entry.Genre)
	}

	if entry.Genre[0] != "Aventure" {
		t.Errorf("entry.Genre[0] is not Aventure %v", entry.Genre)
	}

	if entry.Genre[1] != "Comédie" {
		t.Errorf("entry.Genre[1] is not Comédie %v", entry.Genre)
	}

	if strings.Contains(entry.Genre[len(entry.Genre)-1], ".") {
		t.Errorf("last item in entry.Genre contains a dot %v", entry.Genre)
	}

	if expected := false; list.Entries[1].Favorite != expected {
		t.Errorf("expected: %t, got: %t for %s", expected, list.Entries[1].Favorite, list.Entries[1].Title)
	}
}

func TestBackupCollection(t *testing.T) {
	back := mock.NewBackend()
	Collection("mlcdf", back)

	stuff := back.Data["films-done"]
	if stuff == nil {
		t.Errorf("slug films-done not found")
	}
	collection, ok := stuff.(*domain.Collection)
	if ok == false {
		t.Errorf("cast back.Data into domain.Collection")
	}

	if expectedFilter := "done"; collection.Filter != expectedFilter {
		t.Errorf("expected filter '%s', got '%s'", expectedFilter, collection.Filter)
	}

	if expectedCategory := "films"; collection.Category != expectedCategory {
		t.Errorf("expected category '%s', got '%s'", expectedCategory, collection.Category)
	}

	if l := len(collection.Entries); l < 800 {
		t.Errorf("too few entries: %d", l)
	}

	entry := collection.Entries[0]
	if entry.ID == "" {
		t.Errorf("entry.ID cannot be empty %v", entry)
	}

	if len(entry.Authors) == 0 {
		t.Errorf("entry.Authors cannot be empty %v", entry)
	}
}

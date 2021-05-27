package backup

import (
	"testing"

	"go.mlcdf.fr/sc-backup/internal/backend"
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
	back := backend.NewMemory()
	List("https://www.senscritique.com/liste/Vu_au_cinema/363578", back)

	list, ok := back.Data.(*domain.List)
	if ok == false {
		t.Errorf("cast back.Data into domain.List")
	}

	if expectedSlug := "vu-au-cinema"; list.Slug() != expectedSlug {
		t.Errorf("expected slug '%s', got '%s'", expectedSlug, list.Slug())
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

	if len(entry.Authors) == 0 {
		t.Errorf("entry.Authors cannot be empty %v", entry)
	}
}

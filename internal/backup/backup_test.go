package backup

import (
	"testing"

	"github.com/mlcdf/sc-backup/internal/backend"
	"github.com/mlcdf/sc-backup/internal/sc"
)

func TestCheckForValidUser(t *testing.T) {
	username := "username-that-does-not-exists"
	err := checkForValidUser(username)
	if err == nil {
		t.Errorf("username %s should not exist", username)
	}

	username = "mlcdf"
	err = checkForValidUser(username)
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
	if back.Slug != "vu-au-cinema" {
		t.Errorf("expected %s, got %s", "vu-au-cinema", back.Slug)
	}

	entries := back.Stuff.([]*sc.Entry)
	if l := len(entries); l < 100 {
		t.Errorf("too few entries: %d", l)
	}

	entry := entries[0]
	if entry.ID == "" {
		t.Errorf("entry.ID cannot be empty %v", entry)
	}

	if len(entry.Authors) == 0 {
		t.Errorf("entry.Authors cannot be empty %v", entry)
	}
}

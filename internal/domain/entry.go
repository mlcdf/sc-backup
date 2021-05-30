package domain

import (
	"fmt"

	"github.com/metal3d/go-slugify"
)

// Entry represents an entry in a collection or list : a movie, series, books, etc...
type Entry struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	OriginalTitle string   `json:"original_title,omitempty"`
	Year          int      `json:"year,omitempty"`
	Authors       []string `json:"authors"`
	Rating        int      `json:"rating,omitempty"`
	DoneDate      string   `json:"done_date,omitempty"`
	Comment       string   `json:"comment,omitempty"`
	Favorite      bool     `json:"favorite"`
}

var _ Serializable = (*Collection)(nil)

type Collection struct {
	Entries  []*Entry `json:"entries"`
	Category string   `json:"category"`
	Filter   string   `json:"filter"`
	Username string   `json:"username"`
}

func NewCollection(entries []*Entry, Category, Filter, Username string) *Collection {
	return &Collection{
		Entries:  entries,
		Category: Category,
		Filter:   Filter,
		Username: Username,
	}
}

func (c *Collection) Slug() string {
	return fmt.Sprintf("%s-%s", c.Category, c.Filter)
}

func (c *Collection) CSV() []*Entry {
	return c.Entries
}

func (c *Collection) JSON() interface{} {
	return c
}

var _ Serializable = (*List)(nil)

type List struct {
	Entries     []*Entry `json:"entries"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
}

func NewList(entries []*Entry, Title, Description string) *List {
	return &List{
		Entries:     entries,
		Title:       Title,
		Description: Description,
	}
}

func (l *List) Slug() string {
	return slugify.Marshal(l.Title, true)
}

func (l *List) CSV() []*Entry {
	return l.Entries
}

func (l *List) JSON() interface{} {
	return l
}

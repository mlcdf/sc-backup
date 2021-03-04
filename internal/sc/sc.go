package sc

const URL = "https://www.senscritique.com"

var Categories = []string{"films", "series", "bd", "livres", "albums", "morceaux"}
var Filters = []string{"done", "wish"}

// Entry represents an entry in a collection or list : a movie, series, books, etc...
type Entry struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	FrenchTitle string   `json:"french_title"`
	Year        int      `json:"year,omitempty"`
	Authors     []string `json:"authors"`
	Rating      int      `json:"rating,omitempty"`
}

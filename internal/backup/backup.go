package backup

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/metal3d/go-slugify"
	"github.com/mlcdf/sc-backup/internal/backend"
	"github.com/mlcdf/sc-backup/internal/logx"
	"github.com/mlcdf/sc-backup/internal/pool"
	"github.com/mlcdf/sc-backup/internal/sc"
	"github.com/pkg/errors"
)

type ParseFunc func(document *goquery.Document) ([]*sc.Entry, error)

var client = &http.Client{
	Timeout: time.Second * 20,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func request(url string) (*http.Response, error) {
	logx.Debug("GET %s", url)
	res, err := client.Get(url)

	// check for response error
	if err != nil {
		return nil, errors.Wrapf(err, "failed to GET %s", url)
	}

	if res.StatusCode > 400 {
		return nil, fmt.Errorf("error: http %d for url %s", res.StatusCode, res.Request.URL)
	}

	return res, nil
}

func makeURL(username string, category string, filter string) string {
	return fmt.Sprintf("%s/%s/collection/%s/%s/all/all/all/all/all/all/all/page-", sc.URL, username, filter, category)
}

func makeListURL(url string, index int) string {
	if strings.Contains(url, "page-") {
		re := regexp.MustCompile(`page-(.*)`)
		url = re.ReplaceAllString(url, "page-"+strconv.Itoa(index))
	} else {
		if i := strings.LastIndex(url, "/"); i != -1 {
			url = url + "/"
		}
		url = url + "page-" + strconv.Itoa(index)
	}
	return url
}

func validateUser(username string) error {
	res, err := request(sc.URL + "/" + username)

	if err != nil {
		return errors.Wrap(err, "failed to validate user")
	}

	if res.StatusCode == 301 {
		return fmt.Errorf("username %s does not exist or has a limited profil", username)
	}
	return nil
}

func parseDocument(document *goquery.Document) ([]*sc.Entry, error) {
	entries := make([]*sc.Entry, 0)
	document.Find(".elco-collection-item, .elli-item").Each(func(i int, s *goquery.Selection) {
		id, _ := s.Find(".elco-collection-content > .elco-collection-poster, .elli-media figure").Attr("data-sc-product-id")
		title := strings.TrimSpace(s.Find(".elco-title a").Text())

		var entry = &sc.Entry{
			ID:          id,
			Title:       title,
			FrenchTitle: title,
		}

		entry.Authors = make([]string, 0, 5)
		s.Find(".elco-product-detail a.elco-baseline-a, .elli-content a.elco-baseline-a").Each(func(i int, s *goquery.Selection) {
			author := strings.TrimSpace(s.Text())
			entry.Authors = append(entry.Authors, author)
		})

		parsedDate := strings.TrimSpace(s.Find(".elco-date").Text())
		// some works don't have year, for example Œdipe Roi
		// https://www.senscritique.com/mlcdf/collection/done/livres/all/all/all/all/all/all/list/page-1
		if parsedDate != "" {
			year, err := strconv.Atoi(parsedDate[1 : len(parsedDate)-1])
			if err != nil {
				log.Fatal(err)
			}
			entry.Year = year
		}

		ratingString := strings.TrimSpace(s.Find(".elco-collection-rating.user > a > div > span").Text())

		if ratingString != "" {
			rating, err := strconv.Atoi(ratingString)
			if err != nil {
				log.Fatal(err)
			}
			entry.Rating = rating
		}

		entries = append(entries, entry)
	})

	return entries, nil
}

func collectionSize(document *goquery.Document, filter string) (int, error) {
	_nbOfEntries := strings.TrimSpace(document.Find(fmt.Sprintf("[data-sc-collection-filter=%s] span span", filter)).Text())

	if _nbOfEntries == "" {
		if document.Find(".elco-collection-item-empty").Length() > 0 {
			return 0, nil
		}

		return 0, fmt.Errorf("error: failed to parsed nbOfEntries")
	}
	nbOfEntries, err := strconv.Atoi(_nbOfEntries[1 : len(_nbOfEntries)-1])
	if err != nil {
		return 0, err
	}
	return nbOfEntries, nil
}

func listSize(document *goquery.Document) (int, error) {
	sizeString := strings.TrimSpace(document.Find("[data-rel=list-products-count]").Text())
	if sizeString == "" {
		return 0, nil
	}

	size, err := strconv.Atoi(sizeString)
	if err != nil {
		return 0, err
	}
	return size, nil
}

func listTitle(document *goquery.Document) (string, error) {
	title := strings.TrimSpace(document.Find(".d-heading1.elme-listTitle").Text())

	if title == "" {
		return "", fmt.Errorf("title cannot be empty")
	}
	return title, nil
}

func extractPage(url string, parseFunc ParseFunc) ([]*sc.Entry, error) {
	res, err := request(url)
	if err != nil {
		return nil, err
	}

	document, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}

	entries, err := parseFunc(document)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

// List backs up a list
func List(url string, back backend.Backend) error {
	res, err := request(url)
	if err != nil {
		return nil
	}

	err = back.Create()
	if err != nil {
		return err
	}

	document, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil
	}

	size, err := listSize(document)
	if err != nil {
		return errors.Wrapf(err, "%s", url)
	}

	title, err := listTitle(document)
	if err != nil {
		return errors.Wrapf(err, "%s", url)
	}

	entries, err := parseDocument(document)
	if err != nil {
		return err
	}

	// extract list comments too

	nbOfPages := math.Ceil(float64(size) / 30)

	if nbOfPages > 1 {
		tasks := []*pool.Task{}

		for i := 2; i <= int(nbOfPages); i++ {
			i := i
			tasks = append(tasks, pool.NewTask(func() (interface{}, error) {
				entries, err := extractPage(makeListURL(url, i), parseDocument)
				if err != nil {
					return nil, err
				}
				return entries, nil
			}))
		}

		p := pool.NewPool(tasks, 20)
		p.Run()

		entries, err = p.Merge(entries)
		if err != nil {
			return err
		}
	}

	err = back.SaveList(entries, slugify.Marshal(title, true))
	if err != nil {
		return err
	}

	return nil
}

// Collection backs up a user collection
func Collection(username string, back backend.Backend) error {
	err := validateUser(username)
	if err != nil {
		return err
	}

	logx.Info("Backing up collection for user %s", username)
	back.Create()

	dates, err := parseJournal(username)
	if err != nil {
		return nil
	}

	for _, category := range sc.Categories {
		for _, filter := range sc.Filters {

			url := makeURL(username, category, filter)
			res, err := request(url)
			if err != nil {
				return err
			}

			document, err := goquery.NewDocumentFromResponse(res)
			if err != nil {
				return err
			}

			size, err := collectionSize(document, filter)
			if err != nil {
				return errors.Wrapf(err, "%s", url)
			}

			entries, err := parseDocument(document)
			if err != nil {
				return err
			}

			nbOfPages := math.Ceil(float64(size) / 18)
			if nbOfPages > 1 {
				tasks := []*pool.Task{}

				for i := 2; i <= int(nbOfPages); i++ {
					i := i
					tasks = append(tasks, pool.NewTask(func() (interface{}, error) {
						entries, err := extractPage(url+strconv.Itoa(i), parseDocument)
						if err != nil {
							return nil, err
						}
						return entries, nil
					}))
				}

				p := pool.NewPool(tasks, 20)
				p.Run()

				entries, err = p.Merge(entries)
				if err != nil {
					return err
				}
			}

			if filter == "done" {
				for _, entry := range entries {
					for _, d := range dates {
						if entry.ID == d.ID {
							entry.DoneDate = d.DoneDate
						}
					}
				}
			}

			err = back.SaveCollection(entries, fmt.Sprintf("%s-%s", category, filter))
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func parseJournal(username string) ([]*sc.Entry, error) {
	url := sc.URL + "/" + username + "/journal/all/all"
	res, err := request(url)
	if err != nil {
		return nil, err
	}

	document, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}

	size, err := journalSize(document)
	if err != nil {
		return nil, err
	}

	entries, err := extractDoneDate(document)
	if err != nil {
		return nil, err
	}

	nbOfPages := math.Ceil(float64(size) / 20)
	if nbOfPages > 1 {
		tasks := []*pool.Task{}

		for i := 2; i <= int(nbOfPages); i++ {
			i := i
			tasks = append(tasks, pool.NewTask(func() (interface{}, error) {
				entries, err := extractPage(sc.URL+"/"+username+"/journal/all/all/all/page-"+strconv.Itoa(i)+".ajax", extractDoneDate)
				if err != nil {
					return nil, err
				}
				return entries, nil
			}))
		}

		p := pool.NewPool(tasks, 20)
		p.Run()

		entries, err = p.Merge(entries)
		if err != nil {
			return nil, err
		}
	}

	return entries, nil
}

func extractDoneDate(document *goquery.Document) ([]*sc.Entry, error) {
	entries := make([]*sc.Entry, 0)

	document.Find(".eldi-list-item").Each(func(i int, s *goquery.Selection) {
		date, exists := s.Attr("data-sc-datedone")
		if !exists {
			// ce n'est pas une oeuvre, mais un titre année ou mois
			// on les ignore
			return
		}

		s.Find(".eldi-collection-container").Each(func(i int, s *goquery.Selection) {
			parsedId, exists := s.Find(".eldi-collection-poster").Attr("data-sc-product-id")
			if !exists {
				// pour les épisodes de série, on arrive ici par exemple.
				// on les ignore
				return
			}
			id := strings.TrimSpace(parsedId)
			e := &sc.Entry{
				ID:       id,
				DoneDate: date,
			}
			entries = append(entries, e)
		})
	})
	return entries, nil
}

func journalSize(document *goquery.Document) (int, error) {
	size := 0
	document.Find(".elco-collection-count").Each(func(i int, s *goquery.Selection) {
		parsedValue := strings.TrimSpace(s.Text())
		if parsedValue != "" {
			nb, err := strconv.Atoi(parsedValue[1 : len(parsedValue)-1])
			if err != nil {
				log.Fatal(err)
			}
			size += nb
		}
	})
	return size, nil
}

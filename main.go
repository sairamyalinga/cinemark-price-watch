package main

import (
	"fmt"
	"encoding/json"
	"strings"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	//Get request to HTML page
	url := "https://www.cinemark.com/theatres/wa-bellevue/cinemark-lincoln-square-cinemas-and-imax"
	res, err := http.Get(url)

	if err != nil {
		log.Println("Failed fetching data from website:", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("Status code error: %s", res.Status)
	}

	//Load the HTML doc
	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal("Error loading the document:", err)
	}

	movies := make(map[string]map[string][]string)

	doc.Find("div.showtimeMovieBlock").Each(func(i int, s *goquery.Selection) {
		title := s.Find("a.movieLink > h3").Text()

		var formats []string
		s.Find(".print-type-list").Each(func(j int, f *goquery.Selection) {
			formats = append(formats, strings.TrimSpace(f.Text()))
		})

		movies[title] = make(map[string][]string)
		s.Find(".showtimeMovieTimes").Each(func(j int, f *goquery.Selection) {
			f.Find("a.showtime-link").Each(func(k int, t *goquery.Selection) {
				movies[title][formats[j]] = append(movies[title][formats[j]], t.Text())
			})
		})
	})

	b, err := json.MarshalIndent(movies, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))

}

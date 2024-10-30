package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func getPriceDetails(url string) string {
	details, err := http.Get(url)
	if err != nil {
		log.Println("Failed to fetch price and seat arrangement:", err)
	}
	defer details.Body.Close()
	if details.StatusCode != 200 {
		log.Printf("Status code error: %s", details.Status)
	}
	doc, err := goquery.NewDocumentFromReader(details.Body)
	if err != nil {
		log.Fatal("Error loading the Seating details document:", err)
	}
	price := doc.Find("div.ticketSelectorHeader > h3").Text()
	fmt.Println(price);
	return price

}

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
		//get movie title
		title := s.Find("a.movieLink > h3").Text()
		var formats []string

		//extract movie formats
		s.Find(".attribute-list__item").Each(func(j int, f *goquery.Selection) {
			aTag := f.Find("a")
			if aTag.Length() > 0 {
				dataId, exists := aTag.Attr("data-print-type")
				if exists{
					formats = append(formats, dataId)
				}
			} else {
				formats = append(formats, f.Text())
			}

		})
		fmt.Println(formats)
		movies[title] = make(map[string][]string)
		x := 0
		s.Find(".showtimeMovieTimes").Each(func(j int, f *goquery.Selection) {
			movieTag := f.Find("a.showtime-link")
			if movieTag.Length() > 0{
			f.Find("a.showtime-link").Each(func(k int, t *goquery.Selection) {
				fmt.Println("Movie and format:", title, formats[x])
				seatingURL, exists := t.Attr("href")
				if exists {
					movies[title][formats[x]] = append(movies[title][formats[x]], t.Text(), seatingURL)
					fmt.Println(x, t.Text())
				}
			})
			x = x + 1;
		}     
			
		})
	})

	b, err := json.MarshalIndent(movies, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))

}

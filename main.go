package main

import (
	_ "encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	_"time"

	"github.com/PuerkitoBio/goquery"
)

type Showtime struct {
	Time string
	URL string
}

var movies = make(map[string]map[string][]Showtime)

// Parse the HTML of website

func parseHTML(url string) {

	//Get request to HTML page
	url = "https://www.cinemark.com/theatres/wa-bellevue/cinemark-lincoln-square-cinemas-and-imax"
	root := "https://www.cinemark.com"
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


	doc.Find("div.showtimeMovieBlock").Each(func(i int, s *goquery.Selection) {
		//get movie title
		title := s.Find("a.movieLink > h3").Text()
		movies[title] = make(map[string][]Showtime)
		var currentFormat string
		var times []Showtime

		//Extracting formats and times
		s.Find("div.col-xs-12.col-sm-10").Each(func(j int, f *goquery.Selection){
			f.Children().Each(func(k int, child *goquery.Selection){
				// check if child is format
				if child.Is("ul.attribute-list"){
					child.Find("li.attribute-list__item").Each(func(l int, formatItem *goquery.Selection) {
						aTag := formatItem.Find("a")
						if aTag.Length() > 0 {
							dataId, exists := aTag.Attr("data-print-type")
							if exists {
								currentFormat = dataId // Update current format
							}
						} else {
							currentFormat = formatItem.Text() //Update current format (get text if its not <a> tag)
						}
					})
				}
				//check if child is times list
				if child.Is("div.showtimeMovieTimes"){
					child.Find("a.showtime-link").Each(func(l int, timeLink *goquery.Selection) {
						seatingURL, exists := timeLink.Attr("href") //extract seating URls
						//append show time and seating url
						if exists {
							times = append(times, Showtime{
								Time: timeLink.Text(), 
								URL: root+seatingURL,
							})
						}
					})
					// add the showtimes to format
					if currentFormat != "" {
						movies[title][currentFormat] = times
						times = nil
						currentFormat = ""
					}
				}
			})
			if currentFormat == "" && len(times) > 0 {
				movies[title]["General"] = times // default key if no format is found
			}
		})	
	})
	
}

//function to get price for the show
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
	price := doc.Find("div#ticketSelectorHeader > h3").Text()
	price = strings.Split(price," ")[1]
	return price

}

// function to retrieve seating url from records
func getURL (movie string, time string) (string, error) {
	formats, exists := movies[movie]
	if !exists{
		return "", fmt.Errorf("Movie doesn't exist")
	}
	for _, showtimes := range(formats){
		for _, st := range(showtimes){
			if st.Time == time {
				return st.URL, nil
			}
		}
	}
	return "", fmt.Errorf("No show at that time")
}


func main() {

	//TODO: In an infinite for loop parseHTML and get the details and then check the price came down to desired
	//time.Sleep(1 * time.Hour)

	

}

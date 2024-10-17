package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"

	"github.com/PuerkitoBio/goquery"
)

func main(){
	//Get request to HTML page
	url := "https://www.cinemark.com/theatres/wa-bellevue/cinemark-lincoln-square-cinemas-and-imax"
	res, err := http.Get(url)

	if err != nil {
		log.Println("Failed fetching data from website:",err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("Status code error: %s", res.Status)
	}

	// Read the response body
    // body, err := ioutil.ReadAll(res.Body)
    // if err != nil {
    //     fmt.Println("Error reading response body:", err)
    //     return
    // }

    // htmlContent := string(body)
    // fmt.Println(htmlContent)

	//Load the HTML doc
	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal("Error loading the document:", err)
	}
	

}
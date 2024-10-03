// TODO: handle dates and guids

package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel `xml:"channel"`
}

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Title string `xml:"title"`
	Link string `xml:"link"`
	Description string `xml:"description"`
	Language string `xml:"language"`
	Image Image `xml:"image"`
	Items []Item `xml:"item"`
}

type Item struct {
	XMLName xml.Name `xml:"item"`
	Title string `xml:"title"`
	Link string `xml:"link"`
	Guid string `xml:"guid"`
	PubDate string `xml:"pubDate"`
	Description string `xml:"description"`
	Category string `xml:"category"`
}

type Image struct {
	XMLName xml.Name `xml:"image"`
	Title string `xml:"title"`
	Link string `xml:"link"`
	Url string `xml:"url"`
}

const RSSFile = "rss.xml"
const OutputFile = "site/rss.xml"
const rssAddress = "http://localhost/rss.xml"

func ParseRSS(posts []Post, rssPath string){

	fileBytes := getRSS()
	var rss RSS
	xml.Unmarshal(fileBytes, &rss)

	for _, post := range posts {
		var newItem Item
		newItem.Title = post.Title
		newItem.Link = SitePath + post.Link
		newItem.Guid = SitePath + post.Link
		newItem.Description = post.Description
		newItem.Category = "Technology"
		currentTime := time.Now()
		newItem.PubDate = fmt.Sprintf("%d-%d-%d",
						   			  currentTime.Day(),
						   			  currentTime.Month(),
					   	   			  currentTime.Year())
		if !checkPost(newItem, rss){
			rss.Channel.Items = append(rss.Channel.Items, newItem)
		}
	}

	byteData, err := xml.MarshalIndent(rss, " ", "  ")
	if err != nil {
		log.Fatal("PROBLEM MARSHALING RSS: ", err)
	}
	os.WriteFile(OutputFile, byteData, 0755)
}

func compareItems(in1 Item, in2 Item) bool {
	matched := true
	if in1.Title != in2.Title { matched = false}
	if in1.Link != in2.Link { matched = false}
	return matched
}

func getRSS() []byte {
	resp, err := http.Get(rssAddress)
	if err != nil {
		log.Fatal("Problem with getting url: ", )
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("problem parsing body")
	}
	return body
}

func checkPost(toCheck Item, rss RSS) bool {
	found := false
	for _, item := range rss.Channel.Items {
		if compareItems(item, toCheck) {
			fmt.Println("FOUND ITEM: ", toCheck.Title)
			found = true
		}
	}
	return found
}

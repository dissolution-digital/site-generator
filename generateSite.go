package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"text/template"
)

const outputDir = "./site/"
const postPathBase = "/posts/"
const SitePath = "https://dissolution.digital"

type Post struct {
	Title string
	Link string
	Description string
	FilePath string
}

type templatePost struct {
	title, content string
}

func main(){
	if err := os.Mkdir(outputDir, os.ModePerm); err != nil {
		fmt.Println("Output dir exists, chilling")
	}
	thePosts := generatePosts()
	insertPostNav(thePosts)
	generateIndex(thePosts)
	ParseRSS(thePosts, SitePath)
	generateErrors()
}

func insertPostNav( posts []Post) {
	for ind, file := range posts {
		data, err := os.ReadFile(file.FilePath)
		if err != nil {
			fmt.Println("PROBLEM READING FILE DURING INSERTING POSTNAV: " +
						file.FilePath)
		}
		var navbar string
		workData := string(data)
		if ind == 0 {
			if len(posts) > 1 {
				nextPost := posts[ind+1].Link
				navbar = "<h2><a href=/ > HOME | </a><a href=" +
						  nextPost + "> NEXT > </a></h2>"
			} else {
				navbar = "<h2><a href=/ > HOME </a>"
			}
		} else if ind == len(posts) - 1{
			prevPost := posts[ind-1].Link
			navbar = "<h2><a href=" + prevPost + 
					  "> < PREV </a><a href=/ > | HOME | </a></h2>"
		} else {
			nextPost := posts[ind+1].Link
			prevPost := posts[ind-1].Link
			navbar = "<h2><a href=" + prevPost + 
					  "> < PREV </a><a href=/ > | HOME | </a><a href=" + 
					  nextPost + "> NEXT > </a></h2>"
		}
		workData = strings.Replace(workData, ".png>", ".png>" + navbar, 1)
		if err := os.WriteFile(file.FilePath, []byte(workData), 0755); err != nil {
			log.Fatal("CANT WRITE FILE: " + file.FilePath)
		}
	}
}

func generatePostLink(filename string) string {
	htmlExt := strings.ReplaceAll(filename, ".md", ".html")
	returnLink := postPathBase + htmlExt
	return returnLink
}

func processPost( filePath string ) ProcessedMD {
	postTemplate := "./templates/post.html"
	postPage := template.Must(template.ParseFiles(postTemplate))
	
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("PROBLEM OPENING " + filePath + " FOR READING")
	}
	defer file.Close()

	var fileLines []string
	byLine := bufio.NewScanner(file)
	for byLine.Scan() {
		fileLines = append(fileLines, byLine.Text())
	}
	parsedHTML := ProcessMdtoHTML(fileLines)
	var buff bytes.Buffer
	if err := postPage.Execute(&buff, parsedHTML); err != nil {
		log.Fatal("failure to execute post template")
	}
	parsedHTML.Content = buff.String()
	return parsedHTML
}

func generatePosts() []Post {
	var returnPosts []Post
	dataDir := "./data"
	postDir := outputDir + postPathBase
	if err := os.Mkdir(postDir, os.ModePerm); err != nil {
		fmt.Println("Posts dir exists, chilling")
	}
	files, err := os.ReadDir(dataDir)
	if err != nil {
		log.Fatal("Problem reading data dir in Post Parser")
	}
	for _, f := range files {
		var thisItem Post
		thisItem.Link = generatePostLink(f.Name())
		filePath := dataDir + "/" + f.Name()
		tarFile := outputDir + thisItem.Link
		thisItem.FilePath = tarFile
		processed := processPost(filePath)
		htmlOut := processed.Content
		thisItem.Title = processed.Title
		thisItem.Description = processed.Description
		if err := os.WriteFile(tarFile, []byte(htmlOut), 0755); err != nil {
			log.Fatal("CANT WRITE FILE: " + tarFile)
		}
		returnPosts = append(returnPosts, thisItem)
	}
	return returnPosts
}

func generateIndex( articles []Post ){
	indexFile := "./site/index.html"
	postsString := "<h2>"
	slices.Reverse(articles)
	for _, f := range articles {
		postLink := "<a href=" + f.Link + "> " + f.Title + " </a></br>"
		postsString += postLink
	}
	postsString += "</h2>"
	postTemplate := "./templates/index.html"
	postPage := template.Must(template.ParseFiles(postTemplate))
	var buff bytes.Buffer
	if err := postPage.Execute(&buff, postsString); err != nil {
		log.Fatal("failure to execute post template")
	}
	if err := os.WriteFile(indexFile, []byte(buff.String()), 0755); err != nil {
		log.Fatal("CANT WRITE FILE: " + indexFile)
	}
}

func generateErrors(){
	missingString := "<h1>Not Found.</br>Get Outta Here.</h1>"
	badString := "<h1>Not Allowed.</br>Get Outta Here.</h1>"
	errorTemplate := "./templates/error.html"
	errorPage := template.Must(template.ParseFiles(errorTemplate))
	var missBuff bytes.Buffer
	var badBuff bytes.Buffer
	if err := errorPage.Execute(&missBuff, missingString); err != nil {
		log.Fatal("Problem templating 404 page")
	}
	if err := errorPage.Execute(&badBuff, badString); err != nil {
		log.Fatal("Problem templating 403 page")
	}
	if err := os.WriteFile("./site/404.html", []byte(missBuff.String()), 0755); err != nil {
		log.Fatal("CANT WRITE 404 FILE")
	}
	if err := os.WriteFile("./site/403.html", []byte(badBuff.String()), 0755); err != nil {
		log.Fatal("CANT WRITE 403 FILE")
	}
}

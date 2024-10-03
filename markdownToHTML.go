package main

import (
	"regexp"
	"strings"
)

type ProcessedMD struct {
	Title string
	Description string
	Content string
}

func ProcessMdtoHTML(lines []string) ProcessedMD {
	var outputString string
	body := false
	var returnObj ProcessedMD
	for _, line := range lines {
		lineText := line
		var workText string
		if strings.HasPrefix(lineText, "# ") {
			workText = strings.ReplaceAll(lineText, "# ", "")
			returnObj.Title = workText
			workText = "<h1>" + workText
			workText += "</h1>"
		} else if strings.HasPrefix(lineText, "## ") {
			workText = strings.ReplaceAll(lineText, "## ", "<h2>")
			workText += "</h2>"
		} else if strings.HasPrefix(lineText, "### ") {
			workText = strings.ReplaceAll(lineText, "### ", "<h3>")
			workText += "</h3>"
		} else if strings.HasPrefix(lineText, "#### ") {
			desc := strings.ReplaceAll(lineText, "#### ", "")
			returnObj.Description = desc
		} else if lineText == "" {
			workText = "</br></br>\n"
		} else {
			if body == false {
				workText += "<p>"
				body = true
			}
			re := regexp.MustCompile(`\[.*?\]\(.*?\)`)
			linkSlice := re.FindAllString(lineText, -1)
			for _, link := range linkSlice{
				split := strings.Split(link ,"](")
				linkString := split[0][1:]
				linkLink := split[1][:len(split[1])-1]
				linkHTML := "<a href=" + linkLink + ">" + linkString + "</a>"
				lineText = strings.Replace(lineText, link, linkHTML, 1)
			}
			workText += lineText
		}
		outputString += workText
	}
	outputString += "</p>"
	returnObj.Content = outputString
	return returnObj
}

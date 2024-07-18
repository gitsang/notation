package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func f1() {
	html := `
<!DOCTYPE html>
<html lang="en">
<body class="header-link-page">
<header class="header header-app">
	<div class="show-editor-off">
        <div id="header-second" class="header-second">
            <slice-practice-lists id="title-practice-lists" user="sangc" slice="lPrzc"></slice-practice-lists>
        </div>
	</div>
</header>
	`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		panic(err)
	}

	sliceValue, exists := doc.Find("slice-practice-lists[id='title-practice-lists']").Attr("slice")
	if exists {
		fmt.Println("Slice value:", sliceValue)
	} else {
		fmt.Println("slice attribute not found")
	}
}

func f2() {
	f, err := os.Open("./create_notation.response.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		panic(err)
	}

	sliceValue, exists := doc.Find("slice-practice-lists[id='title-practice-lists']").Attr("slice")
	if exists {
		fmt.Println("Slice value:", sliceValue)
	} else {
		fmt.Println("slice attribute not found")
	}
}

func f3() {
	f, err := os.Open("./overview.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		panic(err)
	}

	doc.Find("div.slice-item.slice-item-full").Each(func(i int, s *goquery.Selection) {
		scoreId, _ := s.Attr("data-score")
		fmt.Println("ScoreId:", scoreId)
		name := s.Find("a.slice-item-title").Text()
		words := strings.Fields(name)
		name = strings.Join(words, " ")
		fmt.Println("Name:", name)
		s.Find("div.slice-item-info").Each(func(j int, info *goquery.Selection) {
			info.Find("span.only10col").Each(func(k int, span *goquery.Selection) {
				fmt.Println("Value:", span.Find("span.text-muted").Text())
			})
		})
	})
}

func main() {
	f3()
}

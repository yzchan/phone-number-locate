package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

const ChromeUA string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36"

func main() {
	client := http.Client{}
	url := "https://ip138.com/mobile.asp?mobile=1360508&action=mobile"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", ChromeUA)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	//fmt.Println(resp.StatusCode)
	dom, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	s := dom.Find("div.table table tbody")
	address := s.Find("tr ").Eq(1).Find("td").Eq(1).Text()
	fmt.Println(address)
	card := s.Find("tr ").Eq(2).Find("td").Eq(1).Text()
	fmt.Println(card)
}

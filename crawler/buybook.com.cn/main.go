package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
	"strings"
)

type PhoneSection struct {
	Section  uint
	Province string
	City     string
	AreaCode string
	Postcode string
	Isp      string
	SimCard  string
}

func main() {
	url := "http://www.buybook.com.cn/lb/1360/all/"
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	dom, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	var records []PhoneSection
	// Find the review items
	dom.Find("table.nlist tr").Each(func(i int, s *goquery.Selection) {
		section := s.Find("td").Eq(0).Text()
		pk, _ := strconv.Atoi(section)
		province := s.Find("td").Eq(1).Text()
		city := s.Find("td").Eq(2).Text()
		areaCode := s.Find("td").Eq(3).Text()
		postcode := s.Find("td").Eq(4).Text()
		isp := s.Find("td").Eq(5).Text()
		simCard := s.Find("td").Eq(6).Text()
		//fmt.Println(i+1, section, province, city, areacode, postcode, isp, simcard)

		if pk > 0 {
			record := PhoneSection{
				Section:  uint(pk),
				Province: strings.TrimSpace(province),
				City:     strings.TrimSpace(city),
				AreaCode: strings.TrimSpace(areaCode),
				Postcode: strings.TrimSpace(postcode),
				Isp:      strings.TrimSpace(isp),
				SimCard:  strings.TrimSpace(simCard),
			}
			records = append(records, record)
		}
	})
	for i, r := range records {
		fmt.Println(i, r)
	}
}

package buybook

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

type GoQueryParser struct {
}

func NewGoQueryParser() *GoQueryParser {
	return &GoQueryParser{}
}

func (r GoQueryParser) Parse(body []byte, keyword string) (p PhoneLoc) {
	var (
		dom *goquery.Document
		err error
	)
	if dom, err = goquery.NewDocumentFromReader(bytes.NewReader(body)); err != nil {
		return
	}
	s := dom.Find("table.nlist tr td:contains('" + keyword + "')").Parent()
	if s.Length() == 0 {
		return
	}

	section := s.Find("td").Eq(0).Text()
	pk, _ := strconv.Atoi(section)
	if pk > 0 {
		p = PhoneLoc{
			Section:  pk,
			Province: strings.TrimSpace(s.Find("td").Eq(1).Text()),
			City:     strings.TrimSpace(s.Find("td").Eq(2).Text()),
			AreaCode: strings.TrimSpace(s.Find("td").Eq(3).Text()),
			Postcode: strings.TrimSpace(s.Find("td").Eq(4).Text()),
			Sp:       strings.TrimSpace(s.Find("td").Eq(5).Text()),
			SimCard:  strings.TrimSpace(s.Find("td").Eq(6).Text()),
		}
	}
	return
}

func (r GoQueryParser) ParseAll(body []byte) (plist []PhoneLoc) {
	var (
		dom *goquery.Document
		err error
	)
	if dom, err = goquery.NewDocumentFromReader(bytes.NewReader(body)); err != nil {
		return
	}

	dom.Find("table.nlist tr").Each(func(i int, s *goquery.Selection) {
		section := s.Find("td").Eq(0).Text()
		pk, _ := strconv.Atoi(section)
		province := s.Find("td").Eq(1).Text()
		city := s.Find("td").Eq(2).Text()
		areaCode := s.Find("td").Eq(3).Text()
		postcode := s.Find("td").Eq(4).Text()
		isp := s.Find("td").Eq(5).Text()
		simCard := s.Find("td").Eq(6).Text()

		if pk > 0 {
			record := PhoneLoc{
				Section:  pk,
				Province: strings.TrimSpace(province),
				City:     strings.TrimSpace(city),
				AreaCode: strings.TrimSpace(areaCode),
				Postcode: strings.TrimSpace(postcode),
				Sp:       strings.TrimSpace(isp),
				SimCard:  strings.TrimSpace(simCard),
			}
			plist = append(plist, record)
		}
	})
	return
}

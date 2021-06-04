package buybook

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	UrlFmt = "http://www.buybook.com.cn/lb/%s/all/"
)

type PhoneLoc struct {
	Section  int    `json:"section"`
	Province string `json:"province"`
	City     string `json:"city"`
	AreaCode string `json:"area_code"`
	Postcode string `json:"postcode"`
	Sp       string `json:"sp"`
	SimCard  string `json:"sim_card"`
}

// Parser 解析器接口，用于解析接口返回的html数据
type Parser interface {
	Parse(body []byte, keyword string) PhoneLoc
	ParseAll(body []byte) []PhoneLoc
}

var Instance *Queryer

type Queryer struct {
	Parser Parser
}

func NewQueryer() *Queryer {
	return &Queryer{
		Parser: NewRegexpParser(),
	}
}

func init() {
	Instance = NewQueryer()
}

func (q *Queryer) Request(url string) (data []byte, err error) {
	var resp *http.Response
	if resp, err = http.Get(url); err != nil {
		return
	}
	defer resp.Body.Close()
	if data, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	return
}

func (q *Queryer) Fetch(phone string) (p PhoneLoc, err error) {
	if len(phone) < 7 {
		return p, errors.New("invalid phone number")
	}

	var html []byte
	if html, err = q.Request(fmt.Sprintf(UrlFmt, phone[0:4])); err != nil {
		return
	}
	p = q.Parser.Parse(html, phone[0:7])
	return
}

func (q *Queryer) FetchSec(phone string) (plist []PhoneLoc, err error) {
	if len(phone) < 7 {
		return plist, errors.New("invalid phone section")
	}
	var html []byte
	if html, err = q.Request(fmt.Sprintf(UrlFmt, phone[0:4])); err != nil {
		return
	}
	plist = q.Parser.ParseAll(html)
	return plist, nil
}

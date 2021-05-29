package taobao

import (
	"errors"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	"net/http"
)

const (
	UrlFmt string = "https://tcc.taobao.com/cc/json/mobile_tel_segment.htm?tel=%s"
)

type Parser interface {
	Parse(text string) PhoneLoc
}

type PhoneLoc struct {
	Mts       string `json:"mts"`
	Province  string `json:"province"`
	CatName   string `json:"catName"`
	TelString string `json:"telString"`
	AreaVid   string `json:"areaVid"`
	IspVid    string `json:"ispVid"`
	Carrier   string `json:"carrier"`
}

type Queryer struct {
	Enc *encoding.Decoder
}

func NewQueryer() *Queryer {
	return &Queryer{
		Enc: simplifiedchinese.GBK.NewDecoder(),
	}
}

func (q *Queryer) Fetch(phone string) (p PhoneLoc, err error) {
	if len(phone) == 7 {
		phone = phone + "0000"
	}

	if len(phone) != 11 {
		err = errors.New("invalid phone number")
		return
	}

	var body []byte
	if body, err = q.Request(phone); err != nil {
		return
	}
	pr := NewStringParser() // 默认使用效率最高的string解析器
	//pr := NewRegexpParser()
	//pr := NewV8Parser()
	p = pr.Parse(string(body))
	return
}

func (q *Queryer) Request(phone string) (body []byte, err error) {
	var resp *http.Response
	if resp, err = http.Get(fmt.Sprintf(UrlFmt, phone)); err != nil {
		return
	}
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	if body, err = q.Enc.Bytes(body); err != nil {
		return
	}
	return
}

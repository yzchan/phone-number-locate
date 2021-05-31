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
	Enc    *encoding.Decoder
	Parser Parser
}

func NewQueryer() *Queryer {
	return &Queryer{
		Enc:    simplifiedchinese.GBK.NewDecoder(),
		Parser: NewStringParser(),
	}
}

var Instance *Queryer

func init() {
	Instance = NewQueryer()
}

func (q *Queryer) Fetch(phone string) (p PhoneLoc, err error) {
	if len(phone) == 7 {
		phone = phone + "0000"
	}

	if len(phone) != 11 {
		err = errors.New("invalid phone number")
		return
	}

	body, err := q.Request(phone)
	if err != nil {
		return
	}
	p = q.Parser.Parse(body)
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

// Parser 解析器接口，用于解析接口返回的jsonp格式的数据
type Parser interface {
	Parse(body []byte) PhoneLoc
}

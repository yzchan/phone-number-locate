package buybook

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	UrlFmt = "http://www.buybook.com.cn/lb/%s/all/"
	ReFmt  = `<tr>
  <td><a href="/sj/{#sec7}/">({#sec7})</a></td>
  <td>([\p{Han}]*)</td>
  <td>([\p{Han}]*)</td>
  <td>(\d*)</td>
  <td>(\d*)</td>
  <td>([\p{Han}]*)</td>
  <td>(.*)
</td>
</tr>`
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

// Parser 解析器接口，用于解析接口返回的jsonp格式的数据
type Parser interface {
	Parse(body []byte) PhoneLoc
}

var Instance *Queryer

type Queryer struct {
	Parser Parser
}

func NewQueryer() *Queryer {
	return &Queryer{}
}

func init() {
	Instance = NewQueryer()
}

func (q *Queryer) Fetch(phone string) (p PhoneLoc, err error) {
	if len(phone) != 7 && len(phone) != 11 {
		err = errors.New("invalid phone number")
		return
	}

	sec := phone[0:7]

	resp, err := http.Get(fmt.Sprintf(UrlFmt, sec[0:4]))
	if err != nil {
		return
	}

	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	re := regexp.MustCompile(strings.Replace(ReFmt, "{#sec7}", sec, -1))
	matched := re.FindAllStringSubmatch(string(html), -1)

	if len(matched) > 0 {
		p.Section, _ = strconv.Atoi(matched[0][1])
		p.Province, p.City, p.AreaCode, p.Postcode, p.Sp, p.SimCard = matched[0][2], matched[0][3], matched[0][4], matched[0][5], matched[0][6], matched[0][7]
	}

	return
}

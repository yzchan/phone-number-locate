package m360

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const UrlFmt = "http://cx.shouji.360.cn/phonearea.php?number=%s"

type PhoneLoc struct {
	Province string `json:"province"`
	City     string `json:"city"`
	Sp       string `json:"sp"`
}

type RespResult struct {
	Code int      `json:"code"`
	Data PhoneLoc `json:"data"`
}

var Instance *Queryer

type Queryer struct {
}

func NewQueryer() *Queryer {
	return &Queryer{}
}

func init() {
	Instance = NewQueryer()
}

func (q *Queryer) Fetch(phone string) (p PhoneLoc, err error) {
	var (
		resp *http.Response
		body []byte
		r    RespResult
	)
	if resp, err = http.Get(fmt.Sprintf(UrlFmt, phone)); err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &r); err != nil {
		return
	}
	p = r.Data
	return
}

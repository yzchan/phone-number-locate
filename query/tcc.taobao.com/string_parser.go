package taobao

import (
	"strings"
)

type StringParser struct {
}

func NewStringParser() *StringParser {
	return &StringParser{}
}

func (s StringParser) Parse(body []byte) (p PhoneLoc) {
	text := string(body[21 : len(body)-2])
	if text == "" {
		return
	}
	text = strings.Replace(text, "'", "", -1)
	text = strings.Replace(text, ",", "", -1)
	text = strings.Replace(text, " ", "", -1)
	text = strings.Replace(text, "\t", "", -1)
	slice := strings.Split(text, "\n")
	p = PhoneLoc{
		Mts:       slice[0][4:],
		Province:  slice[1][9:],
		CatName:   slice[2][8:],
		TelString: slice[3][10:],
		AreaVid:   slice[4][8:],
		IspVid:    slice[5][7:],
		Carrier:   slice[6][8:],
	}
	return
}

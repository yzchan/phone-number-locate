package taobao

import (
	"fmt"
	"regexp"
)

type RegexpParser struct {
	re *regexp.Regexp
}

func NewRegexpParser() *RegexpParser {
	return &RegexpParser{
		re: regexp.MustCompile(`^__GetZoneResult_ = \{\s+mts:'(\d{7})',\s+province:'([\p{Han}]+)',\s+catName:'([\p{Han}]+)',\s+telString:'(\d{11})',\s+areaVid:'(\d+)',\s+ispVid:'(\d+)',\s+carrier:'([\p{Han}]+)'\n\}\n$`),
	}
}

func (r RegexpParser) Parse(body []byte) PhoneLoc {
	fmt.Println("RegexpParser")
	text:=string(body)
	matched := r.re.FindStringSubmatch(text)
	return PhoneLoc{
		Mts:       matched[1],
		Province:  matched[2],
		CatName:   matched[3],
		TelString: matched[4],
		AreaVid:   matched[5],
		IspVid:    matched[6],
		Carrier:   matched[7],
	}
}

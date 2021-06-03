package taobao

import (
	"regexp"
)

type RegexpParser struct {
	re *regexp.Regexp
}

func NewRegexpParser() *RegexpParser {
	return &RegexpParser{
		re: regexp.MustCompile(`^__GetZoneResult_ = \{
[ ]{4}mts:'(\d{7})',
[ ]{4}province:'([\p{Han}]+)',
[ ]{4}catName:'([\p{Han}]+)',
[ ]{4}telString:'(\d{11})',
[\t]areaVid:'(\d+)',
[\t]ispVid:'(\d+)',
[\t]carrier:'([\p{Han}]+)'
\}
$`),
	}
}

func (r RegexpParser) Parse(body []byte) PhoneLoc {
	text := string(body)
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

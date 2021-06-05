package buybook

import (
	"regexp"
	"strconv"
	"strings"
)

const ReFmt = `<tr>
  <td><a href="/sj/\d{7}/">(\d{7})</a></td>
  <td>([\p{Han}]*)</td>
  <td>([\p{Han}]*)</td>
  <td>(\d*)</td>
  <td>(\d*)</td>
  <td>([\p{Han}]*)</td>
  <td>(.*)
</td>
</tr>`

type RegexpParser struct {
	re *regexp.Regexp
}

func NewRegexpParser() *RegexpParser {
	return &RegexpParser{
		re: regexp.MustCompile(ReFmt),
	}
}

func (r RegexpParser) Parse(body []byte, keyword string) (p PhoneLoc) {
	text := string(body)

	r.re = regexp.MustCompile(strings.Replace(ReFmt, "\\d{7}", keyword, -1))
	matched := r.re.FindAllStringSubmatch(text, 1)

	if len(matched) == 0 {
		return
	}
	m := matched[0]
	p.Section, _ = strconv.Atoi(m[1])
	p.Province, p.City, p.AreaCode, p.Postcode, p.Sp, p.SimCard = m[2], m[3], m[4], m[5], m[6], m[7]
	return
}

func (r RegexpParser) ParseAll(body []byte) (plist []PhoneLoc) {
	text := string(body)

	matched := r.re.FindAllStringSubmatch(text, -1)

	if len(matched) == 0 {
		return
	}
	var p PhoneLoc
	for _, m := range matched {
		p.Section, _ = strconv.Atoi(m[1])
		p.Province, p.City, p.AreaCode, p.Postcode, p.Sp, p.SimCard = m[2], m[3], m[4], m[5], m[6], m[7]
		plist = append(plist, p)
	}
	return
}

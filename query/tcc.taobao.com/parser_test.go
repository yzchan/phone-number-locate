package taobao

import (
	"math/rand"
	"testing"
	"time"
)

const JsonpStr string = `__GetZoneResult_ = {
    mts:'1330000',
    province:'广西',
    catName:'中国电信',
    telString:'13300000000',
	areaVid:'30518',
	ispVid:'3399685',
	carrier:'广西电信'
}
`

func BenchmarkStringParser(b *testing.B) {
	b.StopTimer()
	p := NewStringParser()
	rand.Seed(time.Now().UnixNano())
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		p.Parse(JsonpStr)
	}
}
//
//func BenchmarkV8Parser(b *testing.B) {
//	b.StopTimer()
//	p := NewV8Parser()
//	rand.Seed(time.Now().UnixNano())
//	b.StartTimer()
//	for i := 0; i < b.N; i++ {
//		p.Parse(JsonpStr)
//	}
//}
//
//func BenchmarkRegexpParser(b *testing.B) {
//	b.StopTimer()
//	p := NewRegexpParser()
//	rand.Seed(time.Now().UnixNano())
//
//	b.StartTimer()
//	for i := 0; i < b.N; i++ {
//		p.Parse(JsonpStr)
//	}
//}

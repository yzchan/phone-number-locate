package buybook

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

var html []byte

func init() {
	resp, err := http.Get("http://www.buybook.com.cn/lb/1340/all/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	html, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func BenchmarkRegexpParser_Parse(b *testing.B) {
	b.StopTimer()
	parser := NewRegexpParser()
	rand.Seed(time.Now().UnixNano())
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		parser.Parse(html, "1340000")
	}
}

func BenchmarkRegexpParser_ParseAll(b *testing.B) {
	b.StopTimer()
	parser := NewRegexpParser()
	rand.Seed(time.Now().UnixNano())
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		parser.ParseAll(html)
	}
}

func BenchmarkGoQueryParser_Parse(b *testing.B) {
	b.StopTimer()
	parser := NewGoQueryParser()
	rand.Seed(time.Now().UnixNano())
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		parser.Parse(html, "1340000")
	}
}

func BenchmarkGoQueryParser_ParseAll(b *testing.B) {
	b.StopTimer()
	parser := NewGoQueryParser()
	rand.Seed(time.Now().UnixNano())
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		parser.ParseAll(html)
	}
}

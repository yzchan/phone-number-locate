package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"phone-number-locate/work"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	UrlFmt   string = "http://cx.shouji.360.cn/phonearea.php?number=%d"
	FileName string = "data.csv"
)

type RespResult struct {
	Code int `json:"code"`
	Data struct {
		Province string `json:"province"`
		City     string `json:"city"`
		Sp       string `json:"sp"`
	} `json:"data"`
}

type task struct {
	sec int64
	w   *csv.Writer
	rw  sync.RWMutex
}

func (t *task) Task() {
	var resp *http.Response
	var err error
	sec := atomic.LoadInt64(&t.sec)
	// 使用Keep-Alive复用长链接可以提升不少性能，可惜该网站并不支持
	if resp, err = http.Get(fmt.Sprintf(UrlFmt, sec)); err != nil {
		fmt.Println(err)
		return
	}
	result, err := ioutil.ReadAll(resp.Body)
	var r RespResult
	err = json.Unmarshal(result, &r)
	t.rw.Lock()
	_ = t.w.Write([]string{strconv.Itoa(int(sec)), r.Data.Province, r.Data.City, r.Data.Sp})
	t.rw.Unlock()
}

func main() {
	var (
		f   *os.File
		err error
	)
	if f, err = os.OpenFile(FileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666); err != nil {
		panic(err)
	}

	t := &task{
		w: csv.NewWriter(f),
	}
	_ = t.w.Write([]string{"section", "province", "city", "sp"})
	w := work.New(100)
	sections := []int{130}
	bT := time.Now()
	for _, section := range sections {
		for i := 0; i < 10000; i++ {
			atomic.StoreInt64(&t.sec, int64(section*10000+i))
			w.Run(t)
		}
	}
	w.Shutdown()
	t.w.Flush()
	eT := time.Since(bT)
	fmt.Println("Run time: ", eT)
}

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	UrlFmt   string = "http://cx.shouji.360.cn/phonearea.php?number=%d"
	FileName string = "data.csv"
)

var (
	pool chan int
	wg   sync.WaitGroup
	rw   sync.RWMutex
)

type RespResult struct {
	Code int `json:"code"`
	Data struct {
		Province string `json:"province"`
		City     string `json:"city"`
		Sp       string `json:"sp"`
	} `json:"data"`
}

func main() {
	var (
		f   *os.File
		err error
	)

	if f, err = os.OpenFile(FileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666); err != nil {
		panic(err)
	}

	w := csv.NewWriter(f)
	var header = []string{"section", "province", "city", "sp"}
	_ = w.Write(header)
	pool = make(chan int, 500)
	bT := time.Now()

	wg.Add(1)
	go producer([]int{130})

	for i := 0; i < 200; i++ {
		wg.Add(1)
		go consumer(w)
	}

	wg.Wait()
	w.Flush()

	eT := time.Since(bT)
	fmt.Println("Run time: ", eT)
}

func consumer(w *csv.Writer) {
	defer wg.Done()
	var result []byte
	var err error
	var resp *http.Response
	var r RespResult
	for {
		sec, ok := <-pool
		if !ok {
			break
		}
		if resp, err = http.Get(fmt.Sprintf(UrlFmt, sec)); err != nil {
			fmt.Println(err.Error())
			continue
		}
		result, err = ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(result, &r)
		//fmt.Println("consumer:", string(result))
		rw.Lock()
		_ = w.Write([]string{strconv.Itoa(sec), r.Data.Province, r.Data.City, r.Data.Sp})
		rw.Unlock()
	}
}

func producer(sections []int) {
	defer wg.Done()
	for _, section := range sections {
		for i := 0; i < 10000; i++ {
			pool <- section*10000 + i
		}
	}
	close(pool)
}

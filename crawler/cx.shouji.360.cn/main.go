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

var pool chan int

type RespResult struct {
	Code int `json:"code"`
	Data struct {
		Province string `json:"province"`
		City     string `json:"city"`
		Sp       string `json:"sp"`
	} `json:"data"`
}

var wg sync.WaitGroup
var rw sync.RWMutex

func main() {
	f, err := os.OpenFile("data.csv", os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(f)
	var header = []string{"section", "province", "city", "sp"}
	w.Write(header)
	pool = make(chan int, 500)
	wg.Add(1)
	bT := time.Now()
	go producer()
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
	var result []byte
	var err error
	var resp *http.Response
	var r RespResult
	for {
		sec, ok := <-pool
		if !ok {
			break
		}
		if resp, err = http.Get(fmt.Sprintf("http://cx.shouji.360.cn/phonearea.php?number=%d", sec)); err != nil {
			fmt.Println(err.Error())
			continue
		}
		result, err = ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(result, &r)
		fmt.Println("consumer:", string(result))
		rw.Lock()
		w.Write([]string{strconv.Itoa(sec), r.Data.Province, r.Data.City, r.Data.Sp})
		rw.Unlock()
	}
	wg.Done()
}

func producer() {

	// 要采集的号段
	section1 := []int{134}
	section2 := []int{130}
	section3 := []int{133}

	var sections []int
	sections = append(sections, section1...)
	sections = append(sections, section2...)
	sections = append(sections, section3...)
	fmt.Println(len(sections))

	for _, section := range sections {
		for i := 0; i < 10000; i++ {
			//url = fmt.Sprintf("producer: http://cx.shouji.360.cn/phonearea.php?number=%d", section*10000+i)
			//fmt.Println(url)
			pool <- section*10000 + i
		}
	}
	close(pool)
	wg.Done()
}

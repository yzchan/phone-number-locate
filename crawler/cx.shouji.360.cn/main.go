package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var pool chan int

var wg sync.WaitGroup

func main() {
	pool = make(chan int, 2000)
	wg.Add(101)
	bT := time.Now()
	go producer()
	for i := 0; i < 1000; i++ {
		go consumer()
	}
	//time.Sleep(time.Second * 10)
	wg.Wait()
	eT := time.Since(bT)
	fmt.Println("Run time: ", eT)
}

func consumer() {
	var result []byte
	var err error
	var resp *http.Response
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
		fmt.Println("consumer:", string(result))
	}
	wg.Done()
}

func producer() {

	// 要采集的号段
	section1 := []int{134, 135, 136, 137, 138, 139, 150, 151, 152, 157, 158, 159, 188}
	section2 := []int{130, 131, 132, 155, 156, 186}
	section3 := []int{133, 153, 189}

	var sections []int
	sections = append(sections, section1...)
	sections = append(sections, section2...)
	sections = append(sections, section3...)
	fmt.Println(len(sections))

	var url string
	for _, section := range sections {
		for i := 0; i < 1000; i++ {
			url = fmt.Sprintf("producer: http://cx.shouji.360.cn/phonearea.php?number=%d", section*10000+i)
			fmt.Println(url)
			pool <- section*10000 + i
		}
	}
	close(pool)
	wg.Done()
}

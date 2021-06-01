package main

import (
	"encoding/csv"
	"fmt"
	"os"
	m360 "phone-number-locate/query/cx.shouji.360.cn"
	"strconv"
	"sync"
	"time"
)

var (
	FileName = "data.csv"
	pool     chan int
	rw       sync.RWMutex
	wg       sync.WaitGroup
	w        *csv.Writer
)

func main() {
	var (
		f   *os.File
		err error
	)
	if f, err = os.OpenFile(FileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666); err != nil {
		panic(err)
	}
	w = csv.NewWriter(f)
	var header = []string{"section", "province", "city", "sp"}
	w.Write(header)
	pool = make(chan int, 200)
	wg.Add(1)
	bT := time.Now()
	go producer()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go consumer()
	}
	wg.Wait()
	w.Flush()

	eT := time.Since(bT)
	fmt.Println("Run time: ", eT)
}

func consumer() {
	var loc m360.PhoneLoc
	var err error
	for {
		sec, ok := <-pool
		if !ok {
			break
		}
		if loc, err = m360.Instance.Fetch(strconv.Itoa(sec)); err != nil {
			fmt.Println(err)
			continue
		}
		rw.Lock()
		_ = w.Write([]string{strconv.Itoa(sec), loc.Province, loc.City, loc.Sp})
		rw.Unlock()
	}
	wg.Done()
}

func producer() {
	sections := []int{130}
	for _, section := range sections {
		for i := 0; i < 10000; i++ {
			pool <- section*10000 + i
		}
	}
	close(pool)
	wg.Done()
}

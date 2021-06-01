package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"phone-number-locate/kit/work"
	taobao "phone-number-locate/query/tcc.taobao.com"
	"strconv"
	"sync"
	"time"
)

const (
	FileName string = "data.csv"
)

var (
	w  *csv.Writer
	rw sync.RWMutex
)

type task struct {
	sec int
}

func (t *task) Task() {
	var err error
	var loc taobao.PhoneLoc
	if loc, err = taobao.Instance.Fetch(fmt.Sprintf("%d0000", t.sec)); err != nil {
		fmt.Println(err)
		return
	}

	rw.Lock()
	_ = w.Write([]string{
		strconv.Itoa(t.sec),
		loc.Mts,
		loc.Province,
		loc.CatName,
		loc.AreaVid,
		loc.IspVid,
		loc.Carrier,
	})
	rw.Unlock()
}

func main() {
	var (
		f   *os.File
		err error
	)
	if f, err = os.OpenFile(FileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666); err != nil {
		panic(err)
	}
	w = csv.NewWriter(f)

	_ = w.Write([]string{"sec", "mts", "province", "catName", "areaVid", "ispVid", "carrier"})
	wrk := work.New(100)
	sections := []int{130}
	bT := time.Now()
	for _, section := range sections {
		for i := 0; i < 10000; i++ {
			wrk.Run(&task{sec: section*10000 + i})
		}
	}
	wrk.Shutdown()
	w.Flush()
	eT := time.Since(bT)
	fmt.Println("Run time: ", eT)
}

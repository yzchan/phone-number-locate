package main

import (
	"encoding/csv"
	"fmt"
	"os"
	taobao "phone-number-locate/query/tcc.taobao.com"
	"phone-number-locate/work"
	"sync"
	"sync/atomic"
	"time"
)

const (
	FileName string = "data.csv"
)

type task struct {
	sec int64
	m   sync.RWMutex
	w   *csv.Writer
	q   *taobao.Queryer
	p   *taobao.StringParser
}

func (t *task) Task() {
	var err error
	var body []byte
	sec := atomic.LoadInt64(&t.sec)

	if body, err = t.q.Request(fmt.Sprintf("%d0000", sec)); err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(len(string(body)))
	//if len(body) < 24 {
	//	fmt.Println(string(body))
	//	return
	//}
	ploc := t.p.Parse(string(body))
	t.m.Lock()
	_ = t.w.Write([]string{
		ploc.Mts,
		ploc.Province,
		ploc.CatName,
		ploc.AreaVid,
		ploc.IspVid,
		ploc.Carrier,
	})
	t.m.Unlock()
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
		q: taobao.NewQueryer(),
		p: taobao.NewStringParser(),
	}
	_ = t.w.Write([]string{"mts", "province", "catName", "areaVid", "ispVid", "carrier"})
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

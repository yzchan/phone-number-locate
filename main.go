package main

import (
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strconv"
)

type PhoneSection struct {
	Section int    `gorm:"primaryKey;type:int(10)"`
	Prov    string `gorm:"type:varchar(20)"`
	City    string `gorm:"type:varchar(20)"`
	Isp     string `gorm:"type:varchar(20)"`
	Card    string `gorm:"type:varchar(20)"`
	Time    int    `gorm:"type:int(10);index"` // 更新时间
}

type pnl struct {
	conn *gorm.DB
}

var p *pnl

func NewLocation() *pnl {
	var err error
	p = &pnl{}
	p.conn, err = gorm.Open(sqlite.Open("./phone.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return p
}

func (p *pnl) Locate(phoneNum string) (ret PhoneSection, err error) {
	if len(phoneNum) != 11 || phoneNum[0] != 0x31 {
		return ret, errors.New("invalid phone number")
	}
	sectionStr := phoneNum[0:7]
	section, err := strconv.Atoi(sectionStr)
	if err != nil {
		return
	}
	result := p.conn.Take(&ret, section)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ret, errors.New("not found")
	}
	return
}

func main() {
	p := NewLocation()
	if loc, err := p.Locate("1381583****"); err != nil {
		panic(err)
	} else {
		fmt.Println(loc)
	}
}

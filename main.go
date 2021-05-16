package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
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
	rdb  *redis.Client
}

var p *pnl

func NewLocation() *pnl {
	var err error
	p = &pnl{}
	p.conn, err = gorm.Open(sqlite.Open("./phone.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	p.rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return p
}

var ctx = context.Background()

func (p *pnl) Locate(phoneNum string) (ret PhoneSection, err error) {
	if len(phoneNum) != 11 || phoneNum[0] != 0x31 {
		return ret, errors.New("invalid phone number")
	}
	sectionStr := phoneNum[0:7]
	MacStr := phoneNum[0:3]

	decoded, err := p.rdb.HGet(ctx, "PhoneSection:MAC-"+MacStr, sectionStr).Result()

	if err == nil { // 命中缓存
		if err = json.Unmarshal([]byte(decoded), &ret); err != nil {
			return ret, err
		}
		fmt.Println("from cache")
		return ret, nil
	} else { // 未命中缓存
		_ = err
		section, err := strconv.Atoi(sectionStr)
		if err != nil {
			return ret, err
		}
		result := p.conn.Take(&ret, section)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ret, errors.New("not found")
		}
		cacheStr, err := json.MarshalIndent(ret, "", "\t")
		fmt.Println(string(cacheStr))
		if err2 := p.rdb.HSet(ctx, "PhoneSection:MAC-"+MacStr, sectionStr, string(cacheStr)).Err(); err2 != nil {
			return ret, err2
		}
	}
	return ret, nil
}

func main() {
	p := NewLocation()
	if loc, err := p.Locate("1381583****"); err != nil {
		panic(err)
	} else {
		fmt.Println(loc)
	}
}

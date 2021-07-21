package domain

import (
	"github.com/likexian/whois"
	"math"
	"strings"
	"time"
)

func GetDomainExpired() float64 {
	result, err := whois.Whois("baidu.com")
	if err != nil {
		panic(err)
	}
	s := strings.Split(result, "Expiry Date: ")
	s2 := strings.Split(s[1], "\n")
	//去掉多余的回车
	ss := []byte(s2[0])
	ss = ss[:len(ss)-1]
	t2, _ := time.ParseInLocation(time.RFC3339, string(ss), time.Local)
	//计算时间差（天）
	return math.Ceil((t2.Sub(time.Now().Local()).Hours()) / 24)
}

package domain

import (
	"github.com/likexian/whois"
	"log"
	"math"
	"strings"
	"time"
)

// GetDomainExpired 获取域名过期时间
func GetDomainExpired(domainName string) float64 {
	var Dtime float64
	result, err := whois.Whois(domainName)
	if err != nil {
		log.Println(2, "域名过期时间获取err", err)
	}
	s := strings.Split(result, "Expiry Date: ")
	if len(s) == 0 {
		log.Println("获取域名过期超时")
	} else {
		s2 := strings.Split(s[1], "\n")

		//去掉多余的回车
		ss := []byte(s2[0])
		ss = ss[:len(ss)-1]
		t2, err := time.ParseInLocation(time.RFC3339, string(ss), time.Local)
		if err != nil {
			log.Println(2, "时间解析err", err)
		}
		//计算时间差（天）
		Dtime = math.Ceil((t2.Sub(time.Now().Local()).Hours()) / 24)
	}
	return Dtime
}

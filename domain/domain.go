package domain

import (
	"fmt"
	"github.com/likexian/whois"
	"strconv"
	"strings"
	"time"
)

func GetDomainExpired() int {
	result, err := whois.Whois("baidu.com")
	if err != nil {
		panic(err)
	}
	s := strings.Split(result, "Expiry Date: ")
	fmt.Println(s[1])
	todayZero, _ := time.ParseInLocation("2006-01-02 15:04:05","2019-07-16 10:50:04.778", time.Local)
	2026-10-11T11:05:17Z

	s2 := []byte(s[1])
	s3 := string(s2[:4]) + string(s2[5:7]) + string(s2[8:10])
	v1, err := strconv.Atoi(s3)
	//v1 ,err := strconv.ParseFloat(s3,64)
	if err != nil {
		panic(err)
	}
	v2 :=   time.Now().Format("2006-01-02")
	v3 := []byte(v2)
	v4,err := strconv.Atoi( string(v3[:4]) + string(v3[5:7]) +string(v3[8:]))
	if err != nil {
		panic(err)
	}
	fmt.Println(v4,v1)
	now := time.Now().Unix()
	fmt.Println(now)
	return v4-v1
}

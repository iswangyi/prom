package alertinfo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type BackInfo struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}
type Metric struct {
	Name     string `json:"__name__"`
	Instance string `json:"instance"`
	Job      string `json:"job"`
}
type Result struct {
	Metric Metric        `json:"metric"`
	Value  []interface{} `json:"value"`
}
type Data struct {
	ResultType string   `json:"resultType"`
	Result     []Result `json:"result"`
}

func GetAlertLabel() string {
	var backinfo BackInfo
	req, err := http.Get("http://localhost:9090/api/v1/query_range?query=up")
	if err != nil {
		fmt.Println("req err", err)
	}
	data2, err := io.ReadAll(req.Body)
	defer func(data io.ReadCloser) {
		err := data.Close()
		if err != nil {
			fmt.Println("close err", err)
		}
	}(req.Body)
	err = json.Unmarshal(data2, &backinfo)
	if err != nil {
		fmt.Println("json err", err)
	}
	fmt.Println(12345)
	return "123"
}

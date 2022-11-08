package collector

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	procPath = "/proc/loadavg"
)

func parseLoad(data string) (loads []float64, err error) {
	// 访问linux系统下的/proc文件获取信息
	loads = make([]float64, 3)
	parts := strings.Fields(data)
	if len(parts) < 3 {
		return nil, fmt.Errorf("unknwn data from /proc/loadavg")
	}
	for i, v := range parts {
		if loads[i], err = strconv.ParseFloat(v, 64); err != nil {
			return nil, fmt.Errorf("parse Failed")
		}
	}
	return loads, nil
}

func getLoad() (loads []float64, err error) {
	data, err := ioutil.ReadFile(procPath)
	if err != nil {
		return nil, err
	}
	loads, err = parseLoad(string(data))
	if err != nil {
		return nil, err
	}
	return loads, nil
}

package mws

import (
	"fmt"
	"math"
	"time"
)

// 指数退避法
func backOff(cellTime time.Duration, retries float64) {
	var i float64 = 0
	for i <= retries {
		i++
		sleepTime := time.Duration(math.Pow(2, i)) * cellTime
		fmt.Println("sleep ", sleepTime)
		time.Sleep(sleepTime)
	}
}

type BO struct {
	cellTime time.Duration
	retries  float64
	index    float64
}

func (b *BO) Do(f func(), result interface{}) {
	if b.index <= b.retries {
		b.Do(f, result)
	}
}

func test() {
	backOff(time.Second, 20)
}

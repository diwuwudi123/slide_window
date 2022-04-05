package main

import (
	"log"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	for i := 1; i <= 60; i++ {
		nowTime := time.Now().Add(time.Duration(i) * time.Minute)
		nowStr := nowTime.Format("2006-01-02 15:04:05")
		log.Println("nowstr", nowStr, i)
		add(nowTime)
	}
	if getR(time.Now().Add(60*time.Minute)) != 60 {
		t.Error("add num is err")
	}
	if getR(time.Now().Add(80*time.Minute)) != 40 {
		t.Error("add num is err")
	}
}

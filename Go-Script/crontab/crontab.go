package main

import (
	"fmt"
	"github.com/robfig/cron"
	"time"
)

func newWithSeconds() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}


func main() {
	c := newWithSeconds()
	c.AddFunc("0/3 * * * * ? ", func() {
		fmt.Println("f1 start job....", time.Now().Format("2006-01-02 15:04:05"))
	})
	c.AddFunc("3 * * * * ? ", func() {
		fmt.Println("f2 start job....", time.Now().Format("2006-01-02 15:04:05"))
	})
	c.Start()
	select {

	}
}
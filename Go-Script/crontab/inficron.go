package main

import (
	"fmt"
	"strings"
)

// FkZero 格式化数据，去掉多余的0
func FkZero(times string) (fmttime string) {
	if string(times[0]) != "0" {
		return times
	} else if strings.Split(times, "0")[1] == "" {
		fkzero := "0"
		return fkzero
	} else {
		fkzero := strings.Split(times, "0")[1]
		return fkzero
	}
}

// CreateCron 生成crontab语句
func CreateCron(times string) (crontab string) {
	timelists := strings.Split(times, ",")
	if timelists[0] == "d" {
		hours := strings.Split(timelists[1], ":")[0]
		minutes := strings.Split(timelists[1], ":")[1]
		crontab := fmt.Sprintf("* %s %s * * *", FkZero(minutes), FkZero(hours))
		return crontab
	} else if timelists[0] == "w" {
		days := strings.Split(timelists[1], ",")[0]
		hours := strings.Split(strings.Split(timelists[2], ",")[0], ":")[0]
		minutes := strings.Split(strings.Split(timelists[2], ",")[0], ":")[1]
		crontab := fmt.Sprintf("* %s %s * * %s", FkZero(minutes), FkZero(hours), FkZero(days))
		return crontab
	} else if timelists[0] == "m" {
		days := strings.Split(timelists[1], ",")[0]
		hours := strings.Split(strings.Split(timelists[2], ",")[0], ":")[0]
		minutes := strings.Split(strings.Split(timelists[2], ",")[0], ":")[1]
		crontab := fmt.Sprintf("* %s %s %s * *", FkZero(minutes), FkZero(hours), FkZero(days))
		return crontab
	} else {
		if len(times) < 16 {
			mounth := strings.Split(times, "-")[1]
			day := strings.Split(times, "-")[2]
			hours := "0"
			minutes := "0"
			crontab := fmt.Sprintf("0 %s %s %s %s *", FkZero(minutes), FkZero(hours), FkZero(day), FkZero(mounth))
			return crontab
		}
		timelists := strings.Split(times, " ")[0]
		month := strings.Split(timelists, "-")[1]
		day := strings.Split(timelists, "-")[2]
		timework := strings.Split(times, " ")[1]
		hours := strings.Split(timework, ":")[0]
		minutes := strings.Split(timework, ":")[1]
		crontab := fmt.Sprintf("0 %s %s %s %s *", FkZero(minutes), FkZero(hours), FkZero(day), FkZero(month))
		return crontab
	}
}

func main() {
	times := "d,03,12:00"
	crontab := CreateCron(times)
	fmt.Println(crontab)
}

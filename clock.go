package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type data struct {
	index     int
	imaginary float64
}

var times map[int]data

func minSinceMidnight() int {
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local) // % int64(time.Hour*24/time.Second)
	return int(time.Since(midnight)/time.Minute) % int(12*time.Hour/time.Minute)
}

func loadTimeData() error {
	f, err := os.Open("time-data.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		tokens := strings.Split(line, ",")
		key, _ := strconv.ParseInt(tokens[0], 0, 32)
		mi, _ := strconv.ParseInt(tokens[1], 0, 32)
		i, _ := strconv.ParseFloat(tokens[2], 32)

		times[int(key)] = data{index: int(mi), imaginary: i}
	}
	return nil
}

package main

import "time"

func checkDaemon(checker func() error, durationInSeconds int64) {
	ticker := time.NewTicker(time.Second * time.Duration(durationInSeconds))
	for range ticker.C {
		if err := checker(); err != nil {
			panic(err)
		}
	}
}

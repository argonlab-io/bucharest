package utils

import "time"

func RunUntil(callback func() bool, duration time.Duration) {
	timeout := time.Now().Add(duration)
	success := false
	for !success && !time.Now().After(timeout) {
		success = callback()
	}
}

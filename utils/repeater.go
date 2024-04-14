package utils

import (
	"time"
)

func Repeater(function func(), delay int) {
	ticker := time.NewTicker(time.Duration(delay) * time.Second)

	go func() {
		for range ticker.C {
			function()
		}
	}()
}

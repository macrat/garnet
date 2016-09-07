package main

import "fmt"

type Time float64

func (this Time) String() string {
	if float64(this) < 0 {
		return "unknown"
	} else {
		return fmt.Sprintf("%d:%02d", (int)((float64)(this)/60), (int)(this)%60)
	}
}

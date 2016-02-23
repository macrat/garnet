package main

import "fmt"

type Time float64

func (this Time) String() string {
	return fmt.Sprintf("%d:%02d", (int)((float64)(this)/60), (int)(this)%60)
}

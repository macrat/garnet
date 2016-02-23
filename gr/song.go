package main

import (
	"fmt"
	"strconv"
)

type Song struct {
	file    Path
	pos     int
	current bool
	time    Time
}

func NewSong(info map[string]string, current bool) (song *Song, err error) {
	p, err := strconv.Atoi(info["Pos"])
	if err != nil {
		return nil, fmt.Errorf("parse failed song position: \"%s\"", info["Pos"])
	}
	t, err := strconv.Atoi(info["Time"])
	if err != nil {
		return nil, fmt.Errorf("parse failed song duration: \"%s\"", info["Time"])
	}
	return &Song{
		file:    (Path)(info["file"]),
		pos:     p,
		time:    (Time)(t),
		current: current,
	}, nil
}

func (this *Song) String() string {
	if this.current {
		return fmt.Sprintf("% 4s] %s", fmt.Sprintf("[%d", this.pos), this.file)
	} else {
		return fmt.Sprintf("% 4d  %s", this.pos, this.file)
	}
}

func (this *Song) ColoredString() string {
	if this.current {
		return fmt.Sprintf("%3d\033[37;7m %s\033[0m", this.pos, this.file)
	} else {
		return fmt.Sprintf("\033[37;7m%3d\033[0m %s", this.pos, this.file.ColoredString())
	}
}

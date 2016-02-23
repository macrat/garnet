package main

import "fmt"

type Status struct {
	playing                         bool
	repeat, random, single, consume bool
	elapsed                         Time
	current                         *Song
	playlist                        Playlist
}

func NewStatus(status map[string]string, elapsed Time, playlist Playlist) *Status {
	var current *Song
	for _, s := range playlist {
		if s.current {
			current = s
			break
		}
	}

	return &Status{
		playing:  status["state"] == "play",
		repeat:   status["repeat"] == "1",
		random:   status["random"] == "1",
		single:   status["single"] == "1",
		consume:  status["consume"] == "1",
		elapsed:  elapsed,
		current:  current,
		playlist: playlist,
	}
}

func (this *Status) String() (r string) {
	for _, song := range this.playlist {
		r += fmt.Sprintln(song.String())
	}
	if len(this.playlist) > 0 {
		r += "\n"
	}

	flag := func(key string, state bool) string {
		if state {
			return "[" + key + "]"
		} else {
			return " " + key + " "
		}
	}

	var time Time
	if this.current != nil {
		time = this.current.time
	}

	rate := 0
	if this.elapsed > 0 {
		rate = (int)(this.elapsed / time * 100)
	}
	r += fmt.Sprintf("% 6s /% 6s [% 3d%%]%s\n%s\t%s\t%s\t%s",
		this.elapsed,
		time,
		rate,
		map[bool]string{true: " playing", false: ""}[this.playing],
		flag("repeat", this.repeat),
		flag("random", this.random),
		flag("single", this.single),
		flag("consume", this.consume),
	)
	return
}

func (this *Status) ColoredString() (r string) {
	for _, song := range this.playlist {
		r += fmt.Sprintln(song.ColoredString())
	}
	if len(this.playlist) > 0 {
		r += "\n"
	}

	var time Time
	if this.current != nil {
		time = this.current.time
	}

	b2c := map[bool]string{true: "\033[1m", false: "\033[37m"}
	rate := 0
	if this.elapsed > 0 {
		rate = (int)(this.elapsed / time * 100)
	}
	r += fmt.Sprintf("% 6s /% 6s [% 3d%%]%s\n%srepeat\033[0m\t%srandom\033[0m\t%ssingle\033[0m\t%sconsume\033[0m",
		this.elapsed,
		time,
		rate,
		map[bool]string{true: " playing", false: ""}[this.playing],
		b2c[this.repeat],
		b2c[this.random],
		b2c[this.single],
		b2c[this.consume],
	)
	return
}

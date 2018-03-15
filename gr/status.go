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

func (this *Status) PlaybackString() string {
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
	return fmt.Sprintf("% 6s /% 6s [% 3d%%]%s\n%s\t%s\t%s\t%s",
		this.elapsed,
		time,
		rate,
		map[bool]string{true: " playing", false: ""}[this.playing],
		flag("repeat", this.repeat),
		flag("random", this.random),
		flag("single", this.single),
		flag("consume", this.consume),
	)
}

func (this *Status) PlaybackColoredString() string {
	var time Time
	if this.current != nil {
		time = this.current.time
	}

	b2c := map[bool]string{true: "\033[1m", false: "\033[37m"}
	rate := 0
	if this.elapsed > 0 && time > 0 {
		rate = (int)(this.elapsed / time * 100)
	}
	return fmt.Sprintf("% 6s /% 6s [% 3d%%]%s\n%srepeat\033[0m\t%srandom\033[0m\t%ssingle\033[0m\t%sconsume\033[0m",
		this.elapsed,
		time,
		rate,
		map[bool]string{true: " playing", false: ""}[this.playing],
		b2c[this.repeat],
		b2c[this.random],
		b2c[this.single],
		b2c[this.consume],
	)
}

func (this *Status) String() string {
	r := this.playlist.String()

	if len(r) > 1 {
		r += "\n"
	}

	r += this.PlaybackString()

	return r
}

func (this *Status) ColoredString() string {
	r := this.playlist.ColoredString()

	if len(r) > 1 {
		r += "\n"
	}

	r += this.PlaybackColoredString()

	return r
}

type StatusSummary struct {
	Status
}

func (this *Status) Summary() *StatusSummary {
	return &StatusSummary{Status: *this}
}

func (this *StatusSummary) summarize() (left, right int) {
	if this.current == nil {
		left = 0
		right = len(this.playlist)
		if right > 10 {
			right = 10
		}
		return
	}

	left = 0
	if this.current.pos > 10 {
		left = this.current.pos - 10
	}

	right = len(this.playlist)
	if left + 21 < right {
		right = left + 21
	}

	return
}

func (this *StatusSummary) String() string {
	left, right := this.summarize()

	r := this.playlist[left:right].String()

	if right <= len(this.playlist)-1 {
		r += " ...\n"
		r += this.playlist[len(this.playlist)-1].String() + "\n"
	}

	if len(r) > 1 {
		r += "\n"
	}

	r += this.PlaybackString()

	return r
}

func (this *StatusSummary) ColoredString() string {
	left, right := this.summarize()

	r := this.playlist[left:right].ColoredString()

	if right <= len(this.playlist)-1 {
		r += " \033[37m...\033[0m\n"
		r += this.playlist[len(this.playlist)-1].ColoredString() + "\n"
	}

	if len(r) > 1 {
		r += "\n"
	}

	r += this.PlaybackColoredString()

	return r
}

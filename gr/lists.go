package main

import (
	"strings"
	"strconv"

	"golang.org/x/text/width"
)

func preprocess(str string) (r string) {
	for _, c := range strings.ToLower(str) {
		n := width.LookupRune(c).Narrow()
		if n == 0 {
			r += string(c)
		} else {
			r += string(n)
		}
	}
	return
}

func check(target string, queries []string) bool {
	target = preprocess(target)
	for _, q := range queries {
		if !strings.Contains(target, q) {
			return false
		}
	}
	return true
}

type PathList []Path

func (this PathList) Filter(query string) PathList {
	if len(query) == 0 {
		return this
	}

	queries := strings.Split(preprocess(query), " ")
	var result PathList
	for _, path := range this {
		if check(string(path), queries) {
			result = append(result, path)
		}
	}
	return result
}

type Playlist []*Song

func (this Playlist) Filter(query string) Playlist {
	if len(query) == 0 {
		return this
	}

	queries := strings.Split(preprocess(query), " ")
	var result Playlist
	for _, song := range this {
		if check(string(song.file), queries) {
			result = append(result, song)
		}
	}
	return result
}

func (this Playlist) RangeFilter(query string) Playlist {
	if pos, err := strconv.Atoi(query); err == nil && 0 <= pos && pos < len(this) {
		return Playlist{this[pos]}
	}

	if xs := strings.Split(query, "-"); len(xs) == 2 {
		f, fe := strconv.Atoi(xs[0])
		t, te := strconv.Atoi(xs[1])

		if fe != nil && te == nil && xs[0] == "" {
			f = 0
		} else if fe == nil && te != nil && xs[1] == "" {
			t = len(this) - 1
		}
		if fe == nil || te == nil {
			return this[f:t+1]
		}
	}

	return this.Filter(query)
}

func (this Playlist) Has(s Song) bool {
	for _, x := range this {
		if x.file == s.file {
			return true
		}
	}
	return false
}

func (this Playlist) Sub(pl Playlist) Playlist {
	var result Playlist

	for _, s := range this {
		if !pl.Has(*s) {
			result = append(result, s)
		}
	}

	return result
}

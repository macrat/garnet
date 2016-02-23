package main

import (
	"strings"

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

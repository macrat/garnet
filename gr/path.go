package main

import (
	"fmt"
	"strings"
)

type Path string

func (this Path) Split() []string {
	if len(this) == 0 {
		return []string{}
	}
	return strings.Split((string)(this), "/")
}

func (this Path) Directory() string {
	if len(this) == 0 {
		return ""
	}
	pathes := this.Split()
	result := strings.Join(pathes[:len(pathes)-1], "/")
	if result != "" {
		result += "/"
	}
	return result
}

func (this Path) BaseName() string {
	if len(this) == 0 {
		return ""
	}
	pathes := this.Split()
	return pathes[len(pathes)-1]
}

func (this Path) Name() string {
	name := strings.Split(this.BaseName(), ".")
	return strings.Join(name[:len(name)-1], ".")
}

func (this Path) Ext() string {
	name := strings.Split(this.BaseName(), ".")
	return name[len(name)-1]
}

func (this Path) Contains(target Path) bool {
	xs := this.Split()
	ys := target.Split()
	if len(xs) > len(ys) {
		return false
	}

	for i := 0; i < len(xs); i++ {
		if xs[i] != ys[i] {
			return false
		}
	}
	return true
}

func (this Path) String() string {
	return string(this)
}

func (this Path) ColoredString() string {
	return fmt.Sprintf("\033[37m%s\033[0m%s\033[37m.%s\033[0m",
		this.Directory(),
		this.Name(),
		this.Ext(),
	)
}

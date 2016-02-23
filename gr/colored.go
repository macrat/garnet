package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

type ColoredStringer interface {
	fmt.Stringer
	ColoredString() string
}

func coloring(obj []interface{}) (r []interface{}) {
	for _, x := range obj {
		if c, ok := x.(ColoredStringer); ok {
			r = append(r, c.ColoredString())
		} else {
			r = append(r, x)
		}
	}
	return
}

func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	return ok && terminal.IsTerminal(int(f.Fd()))
}

func Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	if isTerminal(w) {
		return fmt.Fprint(w, coloring(a)...)
	} else {
		return fmt.Fprint(w, a...)
	}
}

func Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	if isTerminal(w) {
		return fmt.Fprintln(w, coloring(a)...)
	} else {
		return fmt.Fprintln(w, a...)
	}
}

func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	if isTerminal(w) {
		return fmt.Fprintf(w, format, coloring(a)...)
	} else {
		return fmt.Fprintf(w, format, a...)
	}
}

func Sprint(a ...interface{}) string {
	if isTerminal(os.Stdout) {
		return fmt.Sprint(coloring(a)...)
	} else {
		return fmt.Sprint(a...)
	}
}

func Print(a ...interface{}) (n int, err error) {
	return Fprint(os.Stdout, a...)
}

func Println(a ...interface{}) (n int, err error) {
	return Fprintln(os.Stdout, a...)
}

func Printf(format string, a ...interface{}) (n int, err error) {
	return Fprintf(os.Stdout, format, a...)
}

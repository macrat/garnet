package main

import "testing"

func TestMatchesCheck(t *testing.T) {
	if !check("hello", []string{"h"}) {
		t.Errorf("\"hello\" matches \"h\" but got false")
	}
	if check("hello", []string{"x"}) {
		t.Errorf("\"hello\" not matches \"x\" but got true")
	}
	if !check("hello", []string{"h", "l"}) {
		t.Errorf("\"hello\" matches \"h\" and \"l\" but got false")
	}
	if check("hello", []string{"h", "l", "x"}) {
		t.Errorf("\"hello\" not matches \"h\", \"l\" and \"x\" but got true")
	}
}

func TestPreprocess(t *testing.T) {
	for _, c := range []struct {
		input, expect string
	}{
		{"Hello World", "hello world"},
		{"Ｈｅｌｌｏ　Ｗｏｒｌｄ", "hello world"},
	} {
		r := preprocess(c.input)
		if r != c.expect {
			t.Errorf("\"%s\": expected \"%s\" but got \"%s\"", c.input, c.expect, r)
		}
	}
}

func TestPathListFilter(t *testing.T) {
	xs := PathList{
		"this/is/test",
		"test/file",
		"test/hoge",
		"HOGE/fuga",
	}

	for _, x := range []struct {
		q string
		l int
		e PathList
	}{
		{q: "", e: xs},
		{q: "test", e: PathList{"this/is/test", "test/file", "test/hoge"}},
		{q: "Hoge", e: PathList{"test/hoge", "HOGE/fuga"}},
		{q: "FUGA", e: PathList{"HOGE/fuga"}},
		{q: "test FILE", e: PathList{"test/file"}},
	} {
		f := xs.Filter(x.q)
		if len(f) != len(x.e) {
			t.Errorf("query \"%s\": expected list length %d but got %d", x.q, len(x.e), len(f))
		}
		for i, y := range x.e {
			if f[i] != y {
				t.Errorf("query \"%s\" index %d: expected \"%s\" but got \"%s\"", x.q, i, y, f[i])
			}
		}
	}
}

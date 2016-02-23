package main

import "testing"

func TestSplit(t *testing.T) {
	var path Path = "this/is/test"
	xs := path.Split()

	if len(xs) != 3 {
		t.Errorf("splitted list has incorrect length. (expected: 3, got: %d)", len(xs))
	}

	for i, x := range []string{"this", "is", "test"} {
		if x != xs[i] {
			t.Errorf("splitted list has incorrect item. (expected: \"%s\", got: \"%s\")", x, xs[i])
		}
	}

	path = ""
	if len(path.Split()) != 0 {
		t.Errorf("splitted empty string list has %d item(s).", len(path.Split()))
	}
}

func TestDirectory(t *testing.T) {
	var path Path = "this/is/test.mp3"

	if path.Directory() != "this/is/" {
		t.Errorf("incorrect directory. (expected: \"this/is/\", got: \"%s\")", path.Directory())
	}

	path = "name.aac"
	if path.Directory() != "" {
		t.Errorf("incorrect directory. (expected: \"\", got: \"%s\")", path.Directory())
	}

	path = ""
	if path.Directory() != "" {
		t.Errorf("incorrect directory. (expected: \"\", got: \"%s\")", path.Directory())
	}
}

func TestBaseName(t *testing.T) {
	var path Path = "this/is/test.mp3"

	if path.BaseName() != "test.mp3" {
		t.Errorf("incorrect base name. (expected: \"test.mp3\", got: \"%s\")", path.BaseName())
	}

	path = "name.aac"
	if path.BaseName() != "name.aac" {
		t.Errorf("incorrect base name. (expected: \"name.aac\", got: \"%s\")", path.BaseName())
	}

	path = ""
	if path.BaseName() != "" {
		t.Errorf("incorrect base name. (expected: \"\", got: \"%s\")", path.BaseName())
	}
}

func TestName(t *testing.T) {
	var path Path = "this/is/test.mp3"

	if path.Name() != "test" {
		t.Errorf("incorrect name. (expected: \"test\", got: \"%s\")", path.Name())
	}

	path = "name.hoge..aac"
	if path.Name() != "name.hoge." {
		t.Errorf("incorrect name. (expected: \"name.hoge.\", got: \"%s\")", path.Name())
	}

	path = ""
	if path.Name() != "" {
		t.Errorf("incorrect name. (expected: \"\", got: \"%s\")", path.Name())
	}
}

func TestExt(t *testing.T) {
	var path Path = "this/is/test.mp3"

	if path.Ext() != "mp3" {
		t.Errorf("incorrect extension. (expected: \"mp3\", got: \"%s\")", path.Ext())
	}

	path = "name.hoge..aac"
	if path.Ext() != "aac" {
		t.Errorf("incorrect extension. (expected: \"aac\", got: \"%s\")", path.Ext())
	}

	path = ""
	if path.Ext() != "" {
		t.Errorf("incorrect extension. (expected: \"\", got: \"%s\")", path.Ext())
	}
}

func TestContains(t *testing.T) {
	var a Path = "this"
	var b Path = "this/is"
	var c Path = "this/is/test.mp3"
	var d Path = "this/was"

	for _, x := range []struct {
		x, y Path
		z    bool
	}{
		{x: a, y: a, z: true},
		{x: a, y: b, z: true},
		{x: a, y: c, z: true},
		{x: b, y: a, z: false},
		{x: b, y: b, z: true},
		{x: b, y: c, z: true},
		{x: c, y: a, z: false},
		{x: c, y: b, z: false},
		{x: c, y: c, z: true},
		{x: b, y: d, z: false},
		{x: a, y: b, z: true},
		{x: d, y: c, z: false},
	} {
		if x.x.Contains(x.y) != x.z {
			if x.z {
				t.Errorf("\"%s\" contains \"%s\" but got false", x.x, x.y)
			} else {
				t.Errorf("\"%s\" not contains \"%s\" but got true", x.x, x.y)
			}
		}
	}
}

package main

import "testing"

func TestTimeToString(t *testing.T) {
	for _, x := range []struct {
		t Time
		e string
	}{
		{t: 0, e: "0:00"},
		{t: 0.9, e: "0:00"},
		{t: 1, e: "0:01"},
		{t: 1.1, e: "0:01"},
		{t: 59, e: "0:59"},
		{t: 60, e: "1:00"},
		{t: 61, e: "1:01"},
		{t: 120, e: "2:00"},
		{t: -1, e: "unknown"},
		{t: -0.01, e: "unknown"},
	} {
		if x.t.String() != x.e {
			t.Errorf("Time %f excepted \"%s\" but got \"%s\"", (float64)(x.t), x.e, x.t.String())
		}
	}
}

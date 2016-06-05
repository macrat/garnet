package main

import (
	"strings"
	"testing"
)

func connect(t *testing.T) *Connection {
	conn, err := NewConnection()
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			t.Fatal("Can't connect to server.")
		}
		t.Errorf(err.Error())
	}

	if err == nil && conn == nil {
		t.Error("err is nil but conn is nil too.")
	}

	return conn
}

func statusCheck(t *testing.T, conn *Connection, pos int, playing bool) {
	if st, err := conn.Status(); err != nil {
		t.Error(err.Error())
	} else {
		if st.playing != playing {
			if playing {
				t.Error("failed start playing.")
			} else {
				t.Error("failed stop playing.")
			}
		}

		if playing && st.current == nil {
			t.Error("failed get current song information.")
		} else if st.current.pos != pos {
			t.Errorf("unexpected song playing. expected %d of %d but now playing %d of %d.",
				pos,
				len(st.playlist),
				st.current.pos,
				len(st.playlist),
			)
		}
	}
}

func TestConnection(t *testing.T) {
	handleError := func(err error) {
		if err != nil {
			t.Error(err.Error())
		}
	}

	conn := connect(t)
	defer conn.Close()

	if _, err := conn.Update(); err != nil {
		t.Error(err.Error())
	}

	asis, err := conn.Status()
	if err != nil {
		t.Fatal(err.Error())
	}
	if asis.random {
		handleError(conn.Random())
	}

	handleError(conn.Clear())

	files, err := conn.GetFiles()
	if err != nil {
		t.Error(err.Error())
	} else if files == nil {
		t.Error("GetFiles: err is nil but files is nil too.")
	} else if len(files) < 4 {
		t.Fatal("songs too less. please add more than 4 songs.")
	}

	for i, s := range files {
		if i >= 4 {
			break
		}

		handleError(conn.Add(s))
	}

	if pl, err := conn.Playlist(); err != nil {
		t.Error(err.Error())
	} else if len(pl) != 4 {
		t.Errorf("expected 4 songs in the playlist but got %d songs.", len(pl))
	}

	handleError(conn.Play())
	statusCheck(t, conn, 0, true)

	handleError(conn.Next())
	statusCheck(t, conn, 1, true)

	handleError(conn.Stop())
	statusCheck(t, conn, 1, false)

	handleError(conn.Next())
	statusCheck(t, conn, 2, false)

	handleError(conn.PlayPos(3))
	statusCheck(t, conn, 3, true)

	handleError(conn.Previous())
	statusCheck(t, conn, 2, true)

	handleError(conn.Pause())
	statusCheck(t, conn, 2, false)

	handleError(conn.Previous())
	statusCheck(t, conn, 1, false)

	if err := conn.Move(0, 3); err != nil {
		t.Error(err.Error())
	} else if pl, err := conn.Playlist(); err != nil {
		t.Error(err.Error())
	} else if pl[3].file != files[0] {
		t.Error("failed move song but err is nil")
	}

	if err := conn.Delete(3); err != nil {
		t.Error(err.Error())
	} else if pl, err := conn.Playlist(); err != nil {
		t.Error(err.Error())
	} else if len(pl) != 3 {
		t.Error("failed delete song but err is nil.")
	}

	if err := conn.DeleteRange(1, 3); err != nil {
		t.Error(err.Error())
	} else if pl, err := conn.Playlist(); err != nil {
		t.Error(err.Error())
	} else if len(pl) != 1 {
		t.Error("failed range delete song but err is nil.")
	}

	handleError(conn.Clear())
	if st, err := conn.Status(); err != nil {
		t.Error(err.Error())
	} else if len(st.playlist) != 0 {
		t.Error("failed clear playlist but got err is nil.")
	}

	for _, s := range asis.playlist {
		handleError(conn.Add(s.file))
	}
	if asis.random {
		handleError(conn.Random())
	}

	handleError(conn.Close())
}

func TestModeOption(t *testing.T) {
	conn := connect(t)
	defer conn.Close()

	for i := 0; i < 2; i++ {
		st, err := conn.Status()
		if err != nil {
			t.Fatal(err.Error())
		}

		if err := conn.Repeat(); err != nil {
			t.Error(err.Error())
		} else if newst, err := conn.Status(); err != nil {
			t.Error(err.Error())
		} else if newst.repeat == st.repeat {
			t.Error("failed toggle repeat mode but err is nil.")
		}

		if err := conn.Random(); err != nil {
			t.Error(err.Error())
		} else if newst, err := conn.Status(); err != nil {
			t.Error(err.Error())
		} else if newst.random == st.random {
			t.Error("failed toggle random mode but err is nil.")
		}

		if err := conn.Single(); err != nil {
			t.Error(err.Error())
		} else if newst, err := conn.Status(); err != nil {
			t.Error(err.Error())
		} else if newst.single == st.single {
			t.Error("failed toggle single mode but err is nil.")
		}

		if err := conn.Consume(); err != nil {
			t.Error(err.Error())
		} else if newst, err := conn.Status(); err != nil {
			t.Error(err.Error())
		} else if newst.consume == st.consume {
			t.Error("failed toggle consume mode but err is nil.")
		}
	}
}

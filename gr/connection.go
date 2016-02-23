package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fhs/gompd/mpd"
)

type Connection struct {
	conn *mpd.Client
}

func NewConnection() (conn *Connection, err error) {
	var c *mpd.Client

	host := os.Getenv("MPD_HOST")
	if host == "" {
		c, err = mpd.Dial("tcp", ":6600")
	} else {
		proto := "tcp"
		if strings.HasPrefix(host, "/") {
			proto = "unix"
		}

		c, err = mpd.Dial(proto, host)
	}
	conn = &Connection{conn: c}
	return
}

func (this *Connection) Play() error {
	return this.conn.Play(-1)
}

func (this *Connection) PlayPos(pos int) error {
	return this.conn.Play(pos)
}

func (this *Connection) Pause() error {
	return this.conn.Pause(true)
}

func (this *Connection) Stop() error {
	return this.conn.Stop()
}

func (this *Connection) Next() error {
	if st, err := this.Status(); err != nil {
		return err
	} else if !st.playing {
		if err := this.Play(); err != nil {
			return err
		} else if err := this.Next(); err != nil {
			return err
		}
		return this.Stop()
	}
	return this.conn.Next()
}

func (this *Connection) Previous() error {
	if st, err := this.Status(); err != nil {
		return err
	} else if !st.playing {
		if err := this.Play(); err != nil {
			return err
		} else if err := this.Previous(); err != nil {
			return err
		}
		return this.Stop()
	}
	return this.conn.Previous()
}

func (this *Connection) Clear() error {
	return this.conn.Clear()
}

func (this *Connection) Delete(pos int) error {
	return this.conn.Delete(pos, -1)
}

func (this *Connection) DeleteRange(from int, to int) error {
	return this.conn.Delete(from, to)
}

func (this *Connection) GetFiles() (PathList, error) {
	pathes, err := this.conn.GetFiles()
	if err != nil {
		return nil, err
	}
	var result PathList
	for _, x := range pathes {
		result = append(result, Path(x))
	}
	return result, nil
}

func (this *Connection) Add(path Path) error {
	return this.conn.Add(string(path))
}

func (this *Connection) Move(from int, to int) error {
	return this.conn.Move(from, -1, to)
}

func (this *Connection) Repeat() error {
	st, err := this.Status()
	if err != nil {
		return err
	}
	return this.conn.Repeat(!st.repeat)
}

func (this *Connection) Random() error {
	st, err := this.Status()
	if err != nil {
		return err
	}
	return this.conn.Random(!st.random)
}

func (this *Connection) Single() error {
	st, err := this.Status()
	if err != nil {
		return err
	}
	return this.conn.Single(!st.single)
}

func (this *Connection) Consume() error {
	st, err := this.Status()
	if err != nil {
		return err
	}
	return this.conn.Consume(!st.consume)
}

func (this *Connection) Close() error {
	return this.conn.Close()
}

func (this *Connection) Status() (*Status, error) {
	status, err := this.conn.Status()
	if err != nil {
		return nil, err
	}

	t := 0.0
	if status["elapsed"] != "" {
		t, err = strconv.ParseFloat(status["elapsed"], 64)
		if err != nil {
			return nil, fmt.Errorf("parse failed current elapsed: \"%s\"", status["elapsed"])
		}
	}

	pl, err := this.Playlist()
	if err != nil {
		return nil, err
	}
	return NewStatus(status, (Time)(t), pl), nil
}

func (this *Connection) CurrentSong() (*Song, error) {
	cur, err := this.conn.CurrentSong()
	if err != nil {
		return nil, err
	}
	if cur["file"] != "" {
		return NewSong(cur, true)
	} else {
		return &Song{pos: -1, current: true}, nil
	}
}

func (this *Connection) Playlist() (playlist Playlist, err error) {
	songs, err := this.conn.PlaylistInfo(0, 1048576)
	if err != nil {
		return
	}

	cur, e := this.CurrentSong()
	if e != nil {
		return nil, e
	}

	for i, song := range songs {
		s, err := NewSong(song, i == cur.pos)
		if err != nil {
			return nil, err
		}
		playlist = append(playlist, s)
	}

	return
}

func (this *Connection) Update() (id int, err error) {
	return this.conn.Update("/")
}

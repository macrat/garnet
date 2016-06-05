package main

import (
	"os"
	"strconv"
	"strings"

	"bitbucket.org/macrat/go-lsfmt"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/width"
)

func stringWidth(str string) (w int) {
	for _, c := range norm.NFC.String(str) {
		kind := width.LookupRune(c).Kind()
		if kind == width.EastAsianWide || kind == width.EastAsianFullwidth {
			w += 2
		} else {
			w += 1
		}
	}
	return
}

func handleError(err error) {
	if err != nil {
		Fprintln(os.Stderr, "error:", err.Error())
		os.Exit(1)
	}
}

func main() {
	conn, err := NewConnection()
	handleError(err)
	defer conn.Close()

	showStatus := func() {
		st, err := conn.Status()
		handleError(err)
		Println(st)
	}

	if len(os.Args) <= 1 {
		showStatus()
		return
	}

	type option struct {
		opt, args, help string
		cmd             func(args []string)
	}

	var cmds []option
	cmds = []option{
		{"play", "[POSITION|QUERY]", "Play song.", func(args []string) {
			if len(os.Args) <= 2 {
				handleError(conn.Play())
			} else {
				pl, err := conn.Playlist()
				handleError(err)

				pos, err := strconv.Atoi(args[0])
				if err == nil {
					if pos < 0 || len(pl) <= pos {
						Fprintln(os.Stderr, "out of range")
						os.Exit(1)
					}
					handleError(conn.PlayPos(pos))
				} else {
					t := pl.Filter(strings.Join(args, " "))
					switch len(t) {
					case 0:
						Fprintln(os.Stderr, "no such song")
						os.Exit(1)
					case 1:
						handleError(conn.PlayPos(t[0].pos))
					default:
						Fprintln(os.Stderr, "there is", len(t), "candidates:")
						for _, x := range t {
							Println("", x.file)
						}
						os.Exit(1)
					}
				}
			}
			showStatus()
		}},
		{"pause", "", "Pause playing.", func(args []string) { handleError(conn.Pause()); showStatus() }},
		{"stop", "", "Stop playing.", func(args []string) { handleError(conn.Stop()); showStatus() }},
		{"next", "", "Change to next song of playlist.", func(args []string) { handleError(conn.Next()); showStatus() }},
		{"previous", "", "Change to previous song of playlist.", func(args []string) { handleError(conn.Previous()); showStatus() }},
		{"ls", "[DIRECTORY]", "Show list directories.", func(args []string) {
			var starts Path = ""
			if len(args) > 0 {
				starts = Path(args[0])
			}
			var last string
			var lst []string
			var sizes []int
			files, err := conn.GetFiles()
			handleError(err)
			for _, x := range files {
				if !starts.Contains(x) {
					continue
				}

				xs := x.Split()
				n := xs[len(starts.Split())]
				if n != last {
					if len(xs)-1 != len(starts.Split()) {
						lst = append(lst, n+"/")
						sizes = append(sizes, stringWidth(n)+1)
					} else {
						lst = append(lst, Sprint(Path(n)))
						sizes = append(sizes, stringWidth(n))
					}
					last = n
				}
			}

			if len(lst) == 0 {
				Fprintln(os.Stderr, "no such file or directory.")
				os.Exit(1)
			}

			formatter, err := lsfmt.NewFormatterFile(os.Stdout)
			if err == nil {
				formatter.PrintVerticalWithLength(lst, sizes)
			} else {
				for _, x := range lst {
					Println(x)
				}
			}
		}},
		{"find", "[QUERY]", "Search song.", func(args []string) {
			files, err := conn.GetFiles()
			handleError(err)
			for _, x := range files.Filter(strings.Join(args, " ")) {
				Println(x)
			}
		}},
		{"add", "[QUERY]", "Add songs into playlist.", func(args []string) {
			files, err := conn.GetFiles()
			handleError(err)
			files = files.Filter(strings.Join(args, " "))
			if len(files) == 0 {
				Fprintln(os.Stderr, "no such song")
				os.Exit(1)
			}

			pl, err := conn.Playlist()
			handleError(err)

			lastPos := len(pl)
			Println("added", len(files), "songs:")
			for _, x := range files {
				handleError(conn.Add(x))
				Println("", x)
			}

			st, err := conn.Status()
			handleError(err)
			if !st.playing {
				handleError(conn.PlayPos(lastPos))
			}
		}},
		{"delete", "POSITION|RANGE|QUERY", "Delete songs from playlist.", func(args []string) {
			if len(args) != 1 {
				Fprintln(os.Stderr, "please designation deleting song.")
				os.Exit(1)
			}

			pl, err := conn.Playlist()
			handleError(err)

			t := pl.RangeFilter(strings.Join(args, " "))
			if len(t) == 0 {
				Fprintln(os.Stderr, "no such song")
			} else {
				Println("deleted", len(t), "songs:")
				for _, x := range t {
					Println("", x.file)
				}
				handleError(conn.DeleteAll(t))
			}
		}},
		{"clear", "", "Delete all songs in the playlist.", func(args []string) { handleError(conn.Clear()); showStatus() }},
		{"only", "POSITION|RANGE|QUERY", "Delete all songs except mathed song.", func(args []string) {
			if len(args) != 1 {
				Fprintln(os.Stderr, "please designation deleting song.")
				os.Exit(1)
			}

			pl, err := conn.Playlist()
			handleError(err)

			t := pl.Sub(pl.RangeFilter(strings.Join(args, " ")))
			if len(t) == 0 {
				Fprintln(os.Stderr, "no such song")
			} else {
				Println("deleted", len(t), "songs:")
				for _, x := range t {
					Println("", x.file)
				}
				handleError(conn.DeleteAll(t))
			}
		}},
		{"move", "FROM TO", "Moving song in the playlist.", func(args []string) {
			if len(args) != 2 {
				Fprintln(os.Stderr, "please give `from` position and `to` position.")
				os.Exit(1)
			}
			pl, err := conn.Playlist()
			handleError(err)
			from, err := strconv.Atoi(args[0])
			if err != nil {
				Fprintln(os.Stderr, "parse failed `from` position.")
				os.Exit(1)
			}
			if from < 0 || len(pl) <= from {
				Fprintln(os.Stderr, "`from` position is out of range.")
				os.Exit(1)
			}
			to, err := strconv.Atoi(args[1])
			if err != nil {
				Fprintln(os.Stderr, "parse failed `to` position.")
				os.Exit(1)
			}
			if to < 0 || len(pl) <= to {
				Fprintln(os.Stderr, "`to` position is out of range.")
				os.Exit(1)
			}
			if to == from {
				Fprintln(os.Stderr, "`from` position and `to` position is same.")
				os.Exit(1)
			}
			handleError(conn.Move(from, to))
			showStatus()
		}},
		{"repeat", "", "Toggle repeat mode.", func(args []string) { handleError(conn.Repeat()); showStatus() }},
		{"random", "", "Toggle random mode.", func(args []string) { handleError(conn.Random()); showStatus() }},
		{"single", "", "Toggle single mode.", func(args []string) { handleError(conn.Single()); showStatus() }},
		{"consume", "", "Toggle consume mode.", func(args []string) { handleError(conn.Consume()); showStatus() }},
		{"update", "", "Update database.", func(args []string) {
			id, err := conn.Update()
			handleError(err)
			if id >= 0 {
				Printf("updateing database... #%d\n", id)
			}
		}},
		{"help", "", "Show this message.", func(args []string) {
			putOpts := func(xs []option) {
				var max int
				for _, c := range xs {
					l := len(c.opt) + len(c.args)
					if max < l {
						max = l
					}
				}
				for _, c := range xs {
					Printf(" %s %s%s   %s\n",
						c.opt,
						c.args,
						strings.Repeat(" ", max-len(c.opt)-len(c.args)),
						strings.Replace(c.help, "\n", "\n"+strings.Repeat(" ", max+5), -1),
					)
				}
			}

			if len(args) == 1 {
				var cs []option
				for _, c := range cmds {
					if strings.HasPrefix(c.opt, args[0]) {
						cs = append(cs, c)
					}
				}
				putOpts(cs)
			} else {
				Println("\t\tGarnet")
				Println("\tThe simple client for MPD")
				Println()
				Println("commands:")
				putOpts(cmds)
			}
		}},
	}

	var candidate []option

	for _, cmd := range cmds {
		if os.Args[1] == cmd.opt {
			cmd.cmd(os.Args[2:])
			return
		} else if strings.HasPrefix(cmd.opt, os.Args[1]) {
			candidate = append(candidate, cmd)
		}
	}

	switch len(candidate) {
	case 0:
		Fprintln(os.Stderr, "unknown command:", os.Args[1])
		Fprintf(os.Stderr, "please see help: %s help\n", os.Args[0])
		os.Exit(1)
	case 1:
		candidate[0].cmd(os.Args[2:])
	default:
		Fprint(os.Stderr, "there is some candidates:")
		for _, c := range candidate {
			Fprint(os.Stderr, " ", c.opt)
		}
		Fprintln(os.Stderr)
		os.Exit(1)
	}
}

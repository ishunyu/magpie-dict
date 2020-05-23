package subtitle

import (
	"fmt"
)

type Line struct {
	start uint64
	end   uint64
	text  []string
}

func (l *Line) SetStart(start uint64) {
	l.start = start
}

func (l *Line) SetEnd(end uint64) {
	l.end = end
}

func (l *Line) Append(s string) {
	l.text = append(l.text, s)
}

func (l *Line) String() string {
	return fmt.Sprint(l.start, ",", l.end)
}

type File struct {
	lines []Line
}

func (f *File) Append(l *Line) {
	f.lines = append(f.lines, *l)
}

func (f *File) String() string {
	var s string
	for _, l := range f.lines {
		s += l.String()
		s += "\n"
	}
	return s
}

package parser

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/ishunyu/magpie-dict/src/subtitle"
)

var regex = regexp.MustCompile(`\s*(\d+:\d{2}:\d{2}.\d{3}),(\d+:\d{2}:\d{2}.\d{3})`)

// ParseSBV parses SBV file
func ParseSBV(path string) (*subtitle.File, error) {
	dat, _ := ioutil.ReadFile(path)
	str := string(dat)
	return parseSBV(str), nil
}

func parseSBV(s string) *subtitle.File {

	var file subtitle.File
	var line subtitle.Line

	ss := strings.Split(s, "\n")
	for _, l := range ss {
		_l := strings.TrimRight(l, "\r\n")

		fmt.Println(_l)

		switch {
		case _l == "":
			{
				file.Append(&line)
				line = subtitle.Line{}
			}
		case isTimeline(_l):
			{
				start, end := parseTimeline(_l)
				line.SetStart(start)
				line.SetEnd(end)
			}
		default:
			{
				line.Append(_l)
			}
		}
	}

	return &file
}

func isTimeline(s string) bool {
	res := regex.MatchString(s)
	return res
}

func parseTimeline(s string) (uint64, uint64) {
	_s := strings.Split(s, ",")
	return parseTime(_s[0]), parseTime(_s[1])
}

func parseTime(s string) uint64 {
	hmsS := strings.Split(s, ".")
	S := hmsS[1]
	hms := strings.Split(hmsS[0], ":")
	h := hms[0]
	m := hms[1]
	_s := hms[2]

	v, _ := strconv.ParseUint(h, 10, 64)
	res := v * 60

	v, _ = strconv.ParseUint(m, 10, 64)
	res += v

	res *= 60
	v, _ = strconv.ParseUint(_s, 10, 64)
	res += v

	res *= 1000
	v, _ = strconv.ParseUint(S, 10, 64)
	res += v
	return res
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var errInvalidRange = errors.New("range is invalid")

type interval struct {
	low           float64
	lowInclusive  bool
	high          float64
	highInclusive bool
}

func (i *interval) includes(n float64) bool {
	return ((i.lowInclusive && i.low <= n) || i.low < n) &&
		((i.highInclusive && n <= i.high) || n < i.high)
}

func parseSingleNumber(s string) (*interval, error) {
	i := &interval{}
	var n float64
	var err error
	if n, err = strconv.ParseFloat(s, 64); err != nil {
		return nil, err
	}
	if n >= 0 {
		i.lowInclusive = true
		i.high = n
	} else {
		i.low = n
		i.lowInclusive = false
		i.highInclusive = true
	}
	return i, nil
}

func parseRange(s string) (*interval, error) {
	if !strings.Contains(s, ",") {
		return parseSingleNumber(s)
	}

	switch s[0] {
	case '[', '(':
	default:
		return parseRange(fmt.Sprintf("[%s)", s))
	}

	i := &interval{}
	rngStr := strings.SplitN(s, ",", 2)

	lowRng := rngStr[0]
	highRng := rngStr[1]

	switch lowRng[0] {
	case '[':
		i.lowInclusive = true
	case '(':
	default:
		return nil, errInvalidRange
	}

	var err error
	if i.low, err = strconv.ParseFloat(lowRng[1:], 64); err != nil {
		return nil, fmt.Errorf("not a number: %s", lowRng[1:])
	}

	switch highRng[len(highRng)-1] {
	case ']':
		i.highInclusive = true
	case ')':
	default:
		return nil, errInvalidRange
	}

	highNum := highRng[0 : len(highRng)-1]
	if i.high, err = strconv.ParseFloat(highNum, 64); err != nil {
		return nil, fmt.Errorf("not a number: %s", highNum)
	}

	return i, nil
}

func main() {
	log.SetFlags(0)
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, `usage:
    %s INTERVAL
where INTERVAL is a range of numbers in mathematical notation (e.g. [3,10)).

If a single number n is specified, the range is assumed to be [0,n) or (n,0],
depending on whether n is positive or negative.
`, os.Args[0])
		os.Exit(1)
	}

	intervalStr := os.Args[1]

	var i *interval
	var err error
	if i, err = parseRange(intervalStr); err != nil {
		log.Fatalf("Error parsing range: %s", err)
	}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		var err error
		var n float64
		line := string(s.Bytes())
		if n, err = strconv.ParseFloat(line, 64); err != nil {
			log.Fatalf("Error parsing input: %s", line)
		}
		if i.includes(n) {
			fmt.Println(line)
		}
	}
	if err := s.Err(); err != nil {
		log.Fatalf("Error reading from stdin: %s", err)
	}
}

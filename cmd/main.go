package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/adamliesko/tsp"
)

type lookup struct {
	cityToIndex map[string]tsp.City
	indexToCity []string
}

func getIndex(city string, l *lookup) tsp.City {
	ci, found := l.cityToIndex[city]
	if found {
		return ci
	}
	ci = tsp.City(len(l.cityToIndex))
	l.cityToIndex[city] = ci
	l.indexToCity = append(l.indexToCity, city)
	return ci
}

func readInput() tsp.Problem {
	lookup := &lookup{make(map[string]tsp.City), make([]string, 0, tsp.MAX_CITIES)}
	flights := make([]tsp.Flight, 0, tsp.MAX_FLIGHTS)

	var src string
	stdin := bufio.NewScanner(os.Stdin)
	if stdin.Scan() {
		src = stdin.Text()
		getIndex(src, lookup)
	}
	l := make([]string, 4)
	var i int
	var from, to tsp.City
	var day tsp.Day
	var cost tsp.Money
	for stdin.Scan() {
		customSplit(stdin.Text(), l)
		i, _ = strconv.Atoi(l[2])
		day = tsp.Day(i)
		i, _ = strconv.Atoi(l[3])
		cost = tsp.Money(i)
		from = getIndex(l[0], lookup)
		to = getIndex(l[1], lookup)
		flights = append(flights, tsp.Flight{from, to, day, cost})
	}
	// TODO: remove flights to the start city if not on last day
	p := tsp.NewProblem(flights, lookup.indexToCity, lookup.cityToIndex)
	return p
}

func customSplit(s string, r []string) {
	/* Splits lines of input into 4 parts
	   strictly expects format "{3}[A-Z] {3}[A-Z] \d \d"
	   WARNING: no checks are done at all */
	r[0] = s[:3]
	r[1] = s[4:7]
	pos2 := strings.LastIndexByte(s, ' ')
	r[2] = s[8:pos2]
	r[3] = s[pos2+1:]
}

func main() {
	debug.SetGCPercent(-1)
	timeout := time.After(29 * time.Second)
	problem := readInput()
	solution, err := problem.Solve(timeout)
	if err == nil {
		fmt.Print(solution)
	} else {
		fmt.Println(err)
	}
}

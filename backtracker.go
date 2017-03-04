package tsp

import (
	"sort"
	"fmt"
)

type BacktrackerEngine struct{}

func (d BacktrackerEngine) run(comm comm, buffer *result, task *taskData) {
	trip := make([]Flight, 0)
	possibleFlights := make(map[int][]*Flight, 0) // todo better alloc?

	// create a map of cities to visit - it will be updated in the algo iterations
	toVisit := make(map[City]bool, len(task.problem.cities))
	for _, v := range task.problem.cityToIndex {
		toVisit[v] = true
	}
	currentCity := task.problem.start
	day := 0
	delete(toVisit, currentCity)

	// the meaty loop
	for {
		if len(toVisit) == 0 {
			break
		}

		if _, ok := possibleFlights[day]; !ok {
			possibleFlights[day] = bestFlightsForDay(&task.graph, currentCity, day, toVisit)
		}
		// we have to go back - 1 day
		if _, ok := possibleFlights[day]; !ok || len(possibleFlights[day]) == 0 {
			delete(possibleFlights, day)
			toVisit[currentCity] = true

			fmt.Println(trip)
			// remove last - pop op
			f := trip[len(trip)-1]
			trip = trip[0 : len(trip)-1]

			currentCity = f.From
			day--
		} else {
			// remove first
			f := possibleFlights[day][0]
			possibleFlights[day] = possibleFlights[day][1:]

			trip = append(trip, *f)
			delete(toVisit, f.To)
			currentCity = f.To
			day++
		}

		if len(toVisit) == 0 {
			// todo here we don't have to sort, just take the cheapest
			possibleFinal := bestFlightsForDay(&task.graph, currentCity, day, map[City]bool{task.problem.start: true})

			if len(possibleFlights) > 0 {
				trip = append(trip, *possibleFinal[0])
			} else {
				// backtrack
				// remove last
				f := trip[len(trip)-1]
				trip = trip[:len(trip)-1]

				toVisit[currentCity] = true
				currentCity = f.From
				day--
			}
		}
	}

	// send the result
	comm.isFree()
	for i := 0; i < len(buffer.flights); i++ {
		buffer.flights[i] = trip[i]
	}
	buffer.cost = Cost(buffer.flights)
	comm.resultReady()

	// notify we are done
	comm.done()
}

func bestFlightsForDay(g *Graph, currentCity City, day int, toVisit map[City]bool) []*Flight {
	possibilities := g.data[currentCity][day]
	// if we still have to visit some, then filter only those that we should visit
	if len(toVisit) != 0 {
		var filtered []*Flight
		for _, f := range possibilities {
			if f == nil {
				continue
			}
			if _, ok := toVisit[f.To]; ok {
				filtered = append(filtered, f)
			}
		}
		possibilities = filtered
	}

	// sort the available flights according to the price

	sort.Sort(ByCost(possibilities))
	return possibilities
}

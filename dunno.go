package tsp

type dunno struct{}

func (d dunno) Solve(_ <-chan struct{}, _ Problem) <-chan Solution {
	result := make(chan Solution)
	go func() {
		result <- Solution{}
	}()
	return result
}

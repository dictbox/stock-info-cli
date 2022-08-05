package main

type Grid struct {
	Level int8
	Buy   float64
	Sell  float64
	Hold  int
}

type Stock struct {
	Code  string
	Grids []Grid
}

type Config struct {
	Index  []string
	Stocks []Stock
}

func Map(data []Stock) []string {
	mapped := make([]string, len(data))

	for i, v := range data {
		mapped[i] = v.Code
	}

	return mapped
}

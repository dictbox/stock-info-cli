package main

type Grid struct {
	Level int8
	Buy   float32
	Sell  float32
	Hold  int64
}

type Stock struct {
	Code  string
	Grids []Grid
}

type Config struct {
	Stocks []Stock
}

func Map(data []Stock) []string {
	mapped := make([]string, len(data))

	for i, v := range data {
		mapped[i] = v.Code
	}

	return mapped
}

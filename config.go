package main

type Stock struct {
	Code string
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

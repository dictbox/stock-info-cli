package main

type ExchangeType string

const (
	SH = "SH" //上交所
	SZ = "SZ" //深交所
	HS = "HS" //恒生
	SP = "SP" //新加坡
)

type Grid struct {
	Level int8
	Buy   float64
	Sell  float64
	Hold  int
}

type Stock struct {
	Type  ExchangeType
	Code  string
	Name  string
	Grids []Grid
}

type Config struct {
	Index  []Stock //关注的指数
	Stocks []Stock //关注的标的
}

func Map(data []Stock, typeConvert map[ExchangeType]string) []string {
	mapped := make([]string, len(data))

	for i, v := range data {

		mapped[i] = typeConvert[v.Type] + "." + v.Code
	}

	return mapped
}

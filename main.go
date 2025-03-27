package main

import (
	"fmt"
	"github.com/apoorvam/goterminal"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var mu sync.Mutex
var stocks []string
var indexStocks []string
var result map[string]gjson.Result
var upCount int64
var flatCount int64
var downCount int64
var amountCount float64

func main() {
	config := viper.New()
	config.AddConfigPath(".")
	config.SetConfigName("config")
	config.SetConfigType("json")
	_ = config.ReadInConfig()

	var configJson Config
	_ = config.Unmarshal(&configJson)
	typeConverter := map[ExchangeType]string{
		SH: "1",
		SZ: "0",
		HS: "134",
		SP: "104",
	}
	stocks = Map(configJson.Stocks, typeConverter)
	indexStocks = Map(configJson.Index, typeConverter)
	stocks = append(stocks, indexStocks...)

	count := len(configJson.Stocks)
	lines := count / 2
	if count%2 == 0 {
		lines++
	}

	GetStatsCount()
	GetStockInfo()

	writer := goterminal.New(os.Stdout)
	// start listening for updates and render
	//writer.Start()
	//var writers []io.Writer
	//writers = append(writers, writer)
	//for j := 0; j <= lines; j++ {
	//	writers = append(writers, writer.Newline())
	//}

	chunks := ArrayChunk(configJson.Stocks, 2)

	innerWriter := writer // writers[0]

	for {
		// 打印涨跌家数
		_, _ = fmt.Fprintf(innerWriter, "U:%d|F:%d|D:%d|AMO:%.2f", upCount, flatCount, downCount, amountCount/100000000)
		_, _ = fmt.Fprintln(innerWriter)

		//打印指数
		//innerWriter := writers[0]
		for idx, index := range configJson.Index {
			if idx > 0 && idx%3 == 0 {
				_, _ = fmt.Fprintln(innerWriter)
			}

			info := result[index.Code]
			convert := math.Pow10(int(info.Get("f1").Int()))
			_, _ = fmt.Fprintf(innerWriter, "%6s|%8.3f|%8.3f|%8.3f|%6.2f%%\t",
				index.Code, info.Get("f15").Float()/convert, info.Get("f16").Float()/convert, info.Get("f2").Float()/convert,
				info.Get("f3").Float()/100)
		}

		//goland:noinspection GoUnhandledErrorResult
		fmt.Fprintln(innerWriter)

		for _, chunk := range chunks {
			//innerWriter := writers[i+1]
			for _, v := range chunk {
				info := result[v.Code]
				convert := math.Pow10(int(info.Get("f1").Int()))
				buyPrice := 0.0
				sellPrice := 0.0
				// hold := 0
				if v.Grids != nil {
					price := info.Get("f2").Float() / convert
					for _, g := range v.Grids {
						if g.Buy < price && g.Sell > price {
							buyPrice = g.Buy
							sellPrice = g.Sell
							// hold = g.Hold
							break
						}
					}

					_, _ = fmt.Fprintf(innerWriter, "%6s|%8.3f|%8.3f|%8.3f|%6.2f%%|%8.3f|%8.3f\t",
						v.Code, info.Get("f15").Float()/convert, info.Get("f16").Float()/convert, price,
						info.Get("f3").Float()/100, buyPrice, sellPrice)
				} else {
					_, _ = fmt.Fprintf(innerWriter, "%6s|%8.3f|%8.3f|%8.3f|%6.2f%%|%8.3f|%8.3f\t",
						v.Code, info.Get("f15").Float()/convert, info.Get("f16").Float()/convert, info.Get("f2").Float()/convert,
						info.Get("f3").Float()/100, buyPrice, sellPrice)
				}
			}

			_, _ = fmt.Fprintln(innerWriter)
		}

		_ = writer.Print()

		time.Sleep(time.Second * 5)
		//获取待更新的数据
		GetStatsCount()
		GetStockInfo()
		innerWriter.Clear()
	}

	//GetStockInfo()
}

// GetStockInfo 获取配置中所有标的的信息
func GetStockInfo() {
	//TODO:刷选出对应的CODE

	url := fmt.Sprintf("https://push2.eastmoney.com/api/qt/ulist/get?invt=3&pi=0&pz=30&mpi=2000&secids=%s&fields=f1,f2,f3,f12,f13,f14,f15,f16,f17,f18&po=1", strings.Join(stocks, ","))
	r, e := http.Get(url)
	if e == nil {
		s, _ := io.ReadAll(r.Body)
		body := fmt.Sprintf("%s", s)
		body = strings.TrimRight(strings.TrimLeft(body, "go("), ")")
		//fmt.Println(body)
		bodyResult := gjson.Parse(body)
		total := bodyResult.Get("data.total").Int()
		result = make(map[string]gjson.Result, total)
		bodyResult.Get("data.diff").ForEach(func(key, value gjson.Result) bool {
			name := value.Get("f12").String()
			result[name] = value
			return true
		})
	}
}

// GetStatsCount  获取上涨下跌平居的家数家数
func GetStatsCount() {
	upCount = 0
	flatCount = 0
	downCount = 0
	amountCount = 0
	url := fmt.Sprintf("https://push2.eastmoney.com/api/qt/ulist.np/get?fields=f1,f2,f3,f4,f6,f104,f105,f106,f152&secids=1.000001,0.399001")
	r, e := http.Get(url)
	if e == nil {
		s, _ := io.ReadAll(r.Body)
		body := fmt.Sprintf("%s", s)
		bodyResult := gjson.Parse(body)
		bodyResult.Get("data.diff").ForEach(func(key, value gjson.Result) bool {
			upCount += value.Get("f104").Int()     //上涨家数
			flatCount += value.Get("f106").Int()   //持平家数
			downCount += value.Get("f105").Int()   //下跌家数
			amountCount += value.Get("f6").Float() //成交金额
			return true
		})
	}

}

//func FetchData(i string) {
//	mu.Lock()
//	//s = append(s, i)
//	mu.Unlock()
//	wg.Done()
//}

func ArrayChunk(s []Stock, size int) [][]Stock {
	if size < 1 {
		panic("size: cannot be less than 1")
	}
	length := len(s)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]Stock
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		n = append(n, s[i*size:end])
		i++
	}
	return n
}

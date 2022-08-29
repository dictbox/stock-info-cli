package main

import (
	"fmt"
	"github.com/apoorvam/goterminal"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

//type github struct {
//	Name      string `json:"name"`
//	ShortName string `json:"short_name"`
//	Icons     []struct {
//		Sizes string `json:"sizes"`
//		Src   string `json:"src"`
//		Type  string `json:"type,omitempty"`
//	} `json:"icons"`
//	PreferRelatedApplications bool `json:"prefer_related_applications"`
//	RelatedApplications       []struct {
//		Platform string `json:"platform"`
//		URL      string `json:"url"`
//		ID       string `json:"id"`
//	} `json:"related_applications"`
//}

var wg sync.WaitGroup
var mu sync.Mutex
var stocks []string
var result gjson.Result

func main() {
	config := viper.New()
	config.AddConfigPath(".")
	config.SetConfigName("config")
	config.SetConfigType("json")
	config.ReadInConfig()

	var configJson Config
	config.Unmarshal(&configJson)
	stocks = Map(configJson.Stocks)
	stocks = append(stocks, configJson.Index...)

	count := len(configJson.Stocks)
	lines := count / 2
	if count%2 == 0 {
		lines++
	}

	GetStockInfo()

	//writer := uilive.New()
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
		//打印指数
		//innerWriter := writers[0]
		for idx, indexCode := range configJson.Index {
			if idx > 0 && idx%3 == 0 {
				fmt.Fprintln(innerWriter)
			}

			info := result.Get(indexCode)
			fmt.Fprintf(innerWriter, "%7s|%8.3f|%8.3f|%8.3f|%6.2f%%\t",
				indexCode, info.Get("high").Float(), info.Get("low").Float(), info.Get("price").Float(),
				info.Get("percent").Float()*100)
		}

		fmt.Fprintln(innerWriter)

		for _, chunk := range chunks {
			//innerWriter := writers[i+1]
			for _, v := range chunk {
				info := result.Get(v.Code)
				buyPrice := 0.0
				sellPrice := 0.0
				hold := 0
				if v.Grids != nil {
					price := info.Get("price").Float()
					for _, g := range v.Grids {
						if g.Buy < price && g.Sell > price {
							buyPrice = g.Buy
							sellPrice = g.Sell
							hold = g.Hold
							break
						}
					}

					fmt.Fprintf(innerWriter, "%7s|%8.3f|%8.3f|%8.3f|%6.2f%%|%8.3f|%8.3f|%8d\t",
						v.Code, info.Get("high").Float(), info.Get("low").Float(), price,
						info.Get("percent").Float()*100, buyPrice, sellPrice, hold)
				} else {
					fmt.Fprintf(innerWriter, "%7s|%8.3f|%8.3f|%8.3f|%6.2f%%|%8.3f|%8.3f|%8d\t",
						v.Code, info.Get("high").Float(), info.Get("low").Float(), info.Get("price").Float(),
						info.Get("percent").Float()*100, buyPrice, sellPrice, hold)
				}
			}

			fmt.Fprintln(innerWriter)
		}

		writer.Print()

		time.Sleep(time.Second * 5)
		//获取待更新的数据
		GetStockInfo()
		innerWriter.Clear()
	}

	//
	writer.Reset() // flush and stop rendering

	//r, e := http.Get("https://github.com/manifest.json")
	//if e != nil {
	//	panic(e)
	//}
	//gobj := github{}
	//err := json.NewDecoder(r.Body).Decode(&gobj)
	//if err != nil {
	//	return
	//}
	//fmt.Println(gobj.Icons[0].Sizes)

	//test uilive

	//GetStockInfo()
}

func GetStockInfo() {
	//TODO:刷选出对应的CODE

	url := fmt.Sprintf("https://api.money.126.net/data/feed/%s?callback=go", strings.Join(stocks, ","))
	r, e := http.Get(url)
	if e == nil {
		s, _ := ioutil.ReadAll(r.Body)
		body := fmt.Sprintf("%s", s)
		body = strings.TrimRight(strings.TrimLeft(body, "go("), ")")
		//fmt.Println(body)
		result = gjson.Parse(body)
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

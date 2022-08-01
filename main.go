package main

import (
	"fmt"
	"github.com/gosuri/uilive"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"math"
	"net/http"
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
	stocks = []string{
		"0000001", "1399006",
		"1000651", "0600050",
		"0513050", "0512800",
		"1002701", "0601318",
		"0510050", "0512010",
		"0512980", "0513360",
	}

	count := len(stocks)
	lines := count / 2
	if count%2 == 0 {
		lines++
	}

	GetStockInfo()

	writer := uilive.New()
	// start listening for updates and render
	writer.Start()
	var writers []io.Writer
	writers = append(writers, writer)
	for j := 0; j < lines; j++ {
		writers = append(writers, writer.Newline())
	}

	chunks := ArrayChunk(stocks, 2)

	for {

		for i, chunk := range chunks {
			innerWriter := writers[i]
			for _, v := range chunk {
				info := result.Get(v)
				fmt.Fprintf(innerWriter, "%7s|%8.3f|%8.3f|%8.3f|%6.2f%%\t",
					v, info.Get("high").Float(), info.Get("low").Float(), info.Get("price").Float(),
					info.Get("percent").Float()*100)
			}

			fmt.Fprintln(innerWriter)
		}

		time.Sleep(time.Second * 10)
		//获取待更新的数据
		GetStockInfo()
	}

	//
	writer.Stop() // flush and stop rendering

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

func ArrayChunk(s []string, size int) [][]string {
	if size < 1 {
		panic("size: cannot be less than 1")
	}
	length := len(s)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]string
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

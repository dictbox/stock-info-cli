package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"hash"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"sync"
)

type github struct {
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	Icons     []struct {
		Sizes string `json:"sizes"`
		Src   string `json:"src"`
		Type  string `json:"type,omitempty"`
	} `json:"icons"`
	PreferRelatedApplications bool `json:"prefer_related_applications"`
	RelatedApplications       []struct {
		Platform string `json:"platform"`
		URL      string `json:"url"`
		ID       string `json:"id"`
	} `json:"related_applications"`
}

var wg sync.WaitGroup
var mu sync.Mutex
var stockInfo hash.Hash
var stocks []string

func main() {
	stocks = []string{"sz002701"}
	url := "https://api.money.126.net/data/feed/%s?callback=go"

	for _, v := range stocks {
		if strings.HasPrefix(v, "sz") {
			url = fmt.Sprintf(url, strings.Replace(v, "sz", "1", 1))
		}

		r, e := http.Get(url)
		if e == nil {
			s, _ := ioutil.ReadAll(r.Body)
			body := fmt.Sprintf("%s", s)
			body = strings.TrimRight(strings.TrimLeft(body, "go("), ")")
			fmt.Println(body)
			result := gjson.Parse(body)

			fmt.Println(result.Get("1002701.name"))

		}
	}

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
	//writer := uilive.New()
	//// start listening for updates and render
	//writer.Start()

	//GetStockInfo()
	//count := len(s)
	//
	//lines := count / 3
	//if count%3 == 0 {
	//	lines++
	//}
	//
	//var writers []io.Writer
	//writers = append(writers, writer)
	//for j := 0; j < lines; j++ {
	//	writers = append(writers, writer.Newline())
	//}
	//
	//k := 0
	//
	//for {
	//	chunks := ArrayChunk(s, 3)
	//
	//	for i, chunk := range chunks {
	//		innerWriter := writers[i]
	//		for _, v := range chunk {
	//			fmt.Fprintf(innerWriter, "%d_%d/%3.3f/%3.3f/%3.3f%%  ", i, k, 3.005, 2.85, float64(v))
	//		}
	//		fmt.Fprintln(innerWriter)
	//	}
	//
	//	time.Sleep(time.Second * 2)
	//	k++
	//	//获取待更新的数据
	//	//GetStockInfo()
	//}
	//
	//writer.Stop() // flush and stop rendering
}

func GetStockInfo() {
	for _, v := range stocks {
		wg.Add(1)
		go FetchData(v)
	}
	wg.Wait()
}

func FetchData(i string) {
	mu.Lock()
	//s = append(s, i)
	mu.Unlock()
	wg.Done()
}

func ArrayChunk(s []int, size int) [][]int {
	if size < 1 {
		panic("size: cannot be less than 1")
	}
	length := len(s)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]int
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

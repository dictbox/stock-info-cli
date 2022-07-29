package main

import (
	"encoding/json"
	"fmt"
	"github.com/gosuri/uilive"
	"net/http"
	"time"
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

func main() {
	r, e := http.Get("https://github.com/manifest.json")
	if e != nil {
		panic(e)
	}
	gobj := github{}
	err := json.NewDecoder(r.Body).Decode(&gobj)
	if err != nil {
		return
	}
	fmt.Println(gobj.Icons[0].Sizes)

	//test uilive
	writer := uilive.New()
	// start listening for updates and render
	writer.Start()

	for i := 0; i <= 100; i++ {
		fmt.Fprintf(writer, "Downloading.. (%d/%d) GB\n", i, 100)
		time.Sleep(time.Millisecond * 5)
	}

	fmt.Fprintln(writer, "Finished: Downloaded 100GB")
	writer.Stop() // flush and stop rendering
}

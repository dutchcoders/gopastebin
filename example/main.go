package main

import (
	"fmt"
	"github.com/dutchcoders/gopastebin"
	"io/ioutil"
	"net/url"
	"time"
)

func main() {
	baseURL, _ := url.Parse("https://scrape.pastebin.com/")

	pc := pastebin.New(baseURL)

	for {
		pastes, err := pc.Recent(10)
		if err != nil {
			fmt.Println(err.Error())
		}

		for _, paste := range pastes {
			fmt.Printf("%#v\n", paste)

			raw, err := pc.GetRaw(paste.Key)
			if err != nil {
				fmt.Println("Error", err.Error())
				continue
			}

			defer raw.Close()

			b, err := ioutil.ReadAll(raw)

			fmt.Println("----------------")
			fmt.Printf("%#v\n", string(b))
			fmt.Println("----------------")
		}
		time.Sleep(time.Second * 60)
	}
}

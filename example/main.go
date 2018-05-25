package main

import (
	"bytes"
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
			fmt.Println("----------------")
			fmt.Println("Scrape URL: ", paste.ScrapeURL)
			fmt.Println("Full URL: ", paste.FullURL)
			fmt.Println("Key: ", paste.Key)
			fmt.Println("Title: ", paste.Title)
			fmt.Println("User: ", paste.User)
			fmt.Println("Syntax: ", paste.Syntax)

			raw, err := pc.GetRaw(paste.Key)
			if err != nil {
				fmt.Println("Error", err.Error())
				continue
			}

			defer raw.Close()

			b, err := ioutil.ReadAll(raw)
			b = bytes.Replace([]byte(b), []byte(`\r\n`), []byte("\n"), -1)
			b = bytes.Replace([]byte(b), []byte(`\n`), []byte("\n"), -1)

			fmt.Println("----------------")
			fmt.Printf("%s\n", b)
		}
		time.Sleep(time.Second * 60)
	}
}

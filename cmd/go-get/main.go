package main

import (
	"errors"
	"flag"
	"fmt"
	"go-get/internal/parser"
	"go-get/internal/saver"
	"log"
	"net/url"
	"os"
	"sync"
)

var urls []string

func init() {
	fmt.Println("Get urls ...")

	if err := getUrls(); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	preparedUrls := make(chan string)
	httpResponses := make(chan *parser.HttpResponse)

	var wg sync.WaitGroup

	go func() {
		for url := range preparedUrls {
			resp, err := parser.DoHttp(url)

			if err != nil {
				log.Println(err.Error())
				wg.Done()
				continue
			}

			log.Printf("Page %s has been loaded | latency %d\n", url, resp.Latency)

			httpResponses <- resp
		}
	}()

	go func() {
		for resp := range httpResponses {
			fileName, err := saver.SaveResponseToFile(resp)

			if err != nil {
				log.Println(err.Error())
				wg.Done()
				continue
			}

			log.Printf("Page %s has been saved to -> %s\n", resp.Url, fileName)

			wg.Done()
		}
	}()

	for _, rawUrl := range urls {
		wg.Add(1)

		url, _ := url.Parse(rawUrl)

		if url.Scheme == "" {
			preparedUrls <- "https://" + rawUrl
			continue
		}

		preparedUrls <- rawUrl
	}
	close(preparedUrls)

	wg.Wait()
}

func getUrls() error {
	var link string
	flag.StringVar(&link, "urls", "", "urls for parse")
	flag.Parse()
	urls = flag.Args()

	if len(urls) == 0 {
		return errors.New("You don`t added a urls")
	}

	return nil
}

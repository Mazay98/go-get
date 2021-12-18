package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

var links []string
var chanBuf int

func main() {
	fmt.Println("Get links ...")

	var link string
	flag.StringVar(&link, "links", "", "links for parse")
	flag.IntVar(&chanBuf, "ch-buff", 10, "size buffer channel")
	flag.Parse()

	links = flag.Args()
	if len(links) == 0 {
		log.Fatal("You don`t added a links")
	}

	preparedUrls := make(chan string, chanBuf)

	var wg sync.WaitGroup

	go func() {
		defer wg.Done()
		for preparedUrl := range preparedUrls {
			log.Print(preparedUrl)

			err := parseAndSaveToFile(preparedUrl)
			if err != nil {
				log.Println(err.Error())
				wg.Done()
				continue
			}
		}
	}()

	for _, rawLink := range links {
		preparedUrl, err := url.Parse(rawLink)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if preparedUrl.Scheme == "" {
			rawLink = "https://" + rawLink
		}

		preparedUrls <- rawLink
		wg.Add(1)
	}
	close(preparedUrls)

	wg.Wait()
}

func parseAndSaveToFile(url string) error {
	t := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Failed to get <%s>\nError: %s\n", url, err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	log.Printf("Page %s has been loaded | latency %d\n", url, time.Since(t).Milliseconds())

	fileTypes, err := mime.ExtensionsByType(resp.Header.Get("Content-Type"))
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("go-get_file%s", fileTypes[0])
	i := 0

	for {
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			break
		}

		i++
		fileName = fmt.Sprintf("go-get_file (%d)%s", i, fileTypes[0])
	}

	err = ioutil.WriteFile(fileName, body, 0666)
	if err != nil {
		return fmt.Errorf("Failed to get body document <%s>\nError: %s\n", url, err.Error())
	}

	log.Printf("Page %s has been saved to -> %s\n", url, fileName)

	return nil
}

package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"sync"
)

var baseUrl string = "https://geoftp.ibge.gov.br/organizacao_do_territorio/malhas_territoriais/malhas_de_setores_censitarios__divisoes_intramunicipais/censo_2010/setores_censitarios_kmz/"

func extractLinks() []string {

	response, err := http.Get(baseUrl)

	if err != nil {
		panic(err)
	}

	var links []string

	tokenizer := html.NewTokenizer(response.Body)

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			break
		}

		token := tokenizer.Token()

		if tokenType == html.StartTagToken && token.Data == "a" {
			link := extractHref(token)
			links = append(links, link)

			fmt.Println(link)
		}

	}

	return links

}

func extractHref(token html.Token) string {
	var link string

	for _, attr := range token.Attr {

		if attr.Key == "href" {
			if len(attr.Val) > 4 && attr.Val[len(attr.Val)-4:] == ".kmz" {
				link = attr.Val
				break
			}
		}

	}

	return link
}

func downloadKmz(link string, wg *sync.WaitGroup) {

	defer wg.Done()
	fmt.Println("Downloading " + link)

	if len(link) == 0 {
		return
	}

	download, err := http.Get(baseUrl + link)

	if err != nil {
		panic(err)
	}

	arquivo, err := os.Create("kmz/" + link)

	if err != nil {
		panic(err)
	}

	defer arquivo.Close()

	size, err := io.Copy(arquivo, download.Body)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Arquivo %v salvo com %v bytes\n", link, size)

}

func downloadKmzs() {

	links := extractLinks()
	var wg sync.WaitGroup

	for _, link := range links {
		wg.Add(1)
		go downloadKmz(link, &wg)
	}

	wg.Wait()
}

func main() {
	downloadKmzs()
}

package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"strconv"
)

type mapNodesOperation func(nodes []*cdp.Node) []string

func main() {
	var nodes = getNodesBySelector("section.supported-document", "tbody.list>tr>td")
	var mappedNodes = mapNodes(nodes, mapToCountries)
	var divided = divideIntoChunksOf25(mappedNodes)
	for i, d := range divided {
		writeToFile(format(d, i), "countries"+strconv.Itoa(i)+".json")
	}
}

func getNodesBySelector(waitVisibleSelector string, nodesSelector string) []*cdp.Node {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var nodes []*cdp.Node
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://onfido.com/supported-documents/"),
		chromedp.WaitVisible(waitVisibleSelector),
		chromedp.Nodes(nodesSelector, &nodes, chromedp.ByQueryAll),
	); err != nil {
		log.Println(err)
	}
	return nodes
}

func mapNodes(nodes []*cdp.Node, fn mapNodesOperation) []string {
	return fn(nodes)
}
func format(countries []string, iteration int) string {
	formatted := "{\n  \"dev-onfido-supported-countries\": [\n"
	for i, c := range countries {
		var id = 25*iteration + i
		formatted = formatted + "{\n\"PutRequest\":{\n\"Item\":{\n\"id\":{\n\"N\":\"" + strconv.Itoa(id) + "\"\n},\n\"country\":{\n\"S\":\"" + c + "\"\n}\n}\n}\n},"
	}
	formatted = formatted + "\n]\n}"
	return formatted
}
func writeToFile(content string, fileName string) {
	fmt.Println("Writing file " + fileName)
	f, errCreate := os.Create(fileName)
	if errCreate != nil {
		fmt.Println(errCreate)
		return
	}

	_, errWrite := f.WriteString(content)
	if errWrite != nil {
		fmt.Println(errWrite)
		f.Close()
		return
	}

	errClose := f.Close()
	if errClose != nil {
		fmt.Println(errClose)
		return
	}
}

func divideIntoChunksOf25(countries []string) [][]string {
	chunkSize := 25
	var divided [][]string

	for i := 0; i < len(countries); i += chunkSize {
		end := i + chunkSize
		if end > len(countries) {
			end = len(countries)
		}
		divided = append(divided, countries[i:end])
	}
	return divided
}

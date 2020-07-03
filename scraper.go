package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"log"
	"os"
)

type mapNodesOperation func(nodes []*cdp.Node) []string

func main() {
	var nodes = getNodesBySelector("section.supported-document", "tbody.list>tr>td")
	var mappedNodes = mapNodes(nodes, mapToCountries)
	writeToFile(mappedNodes, "countries.txt")
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

func writeToFile(countries []string, fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, c := range countries {
		_, err := f.WriteString(c + "\n")
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
	}

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

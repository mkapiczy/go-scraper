package main

import "github.com/chromedp/cdproto/cdp"

type row struct {
	country      string
	documentType string
}

func mapToCountries(nodes []*cdp.Node) []string {
	var rows = getSupportedDocumentRows(nodes)
	var countries = getCountriesForDrivingLicenseDocumentType(rows)
	return countries
}

func getSupportedDocumentRows(nodes []*cdp.Node) []row {
	var rows []row
	var nodesInChunksOfFour = divideNodesIntoChunksOfFour(nodes)
	for _, n := range nodesInChunksOfFour {
		rows = append(rows, row{n[0].Children[0].NodeValue, n[2].Children[0].NodeValue})
	}
	return rows
}

func divideNodesIntoChunksOfFour(nodes []*cdp.Node) [][]*cdp.Node {
	chunkSize := 4
	var divided [][]*cdp.Node

	for i := 0; i < len(nodes); i += chunkSize {
		end := i + chunkSize
		if end > len(nodes) {
			end = len(nodes)
		}
		divided = append(divided, nodes[i:end])
	}
	return divided
}

func getCountriesForDrivingLicenseDocumentType(rows []row) []string {
	var countries []string
	for _, r := range rows {
		if r.documentType == "Driving Licence" {
			countries = append(countries, r.country)
		}
	}
	return countries
}

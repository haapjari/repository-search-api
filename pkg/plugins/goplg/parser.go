package goplg

import (
	"encoding/json"
	"errors"
	"strings"
)

type Parser struct {
}

func NewParser() *Parser {
	return new(Parser)
}

func (p *Parser) ParseRepository(url string) (string, string) {
	if strings.Contains(url, "github.com") {
		urlArray := strings.Split(url, "/")

		return urlArray[1], urlArray[2]
	}

	return "", ""
}

func (p *Parser) ParseSourceGraphResponse(data string) (map[string]interface{}, error) {
	var responseAsJsonMap map[string]interface{}
	var err error

	err = json.Unmarshal([]byte(string(data)), &responseAsJsonMap)
	if err != nil {

		return nil, err
	}

	dataArray := responseAsJsonMap["data"]
	if dataArray == nil {
		err = errors.New("unable to find 'data' element from response")

		return nil, err
	}

	dataMap := dataArray.(map[string]interface{})
	if dataMap == nil {
		err = errors.New("unable to convert array to map")

		return nil, err
	}

	searchArray := dataMap["search"]
	if searchArray == nil {
		err = errors.New("unable to find 'search' element from response")

		return nil, err
	}

	searchMap := searchArray.(map[string]interface{})
	if searchMap == nil {
		err = errors.New("unable to convert array to map")

		return nil, err
	}

	resultsArray := searchMap["results"]
	if resultsArray == nil {
		err = errors.New("unable to find 'results' element from response")

		return nil, err
	}

	resultsMap := resultsArray.(map[string]interface{})
	if resultsMap == nil {
		err = errors.New("unable to convert array to map")

		return nil, err
	}

	return resultsMap, nil
}

// parseLibrariesFromModFile parses the library names from the given go.mod file and returns them as a slice of strings.
func parseLibrariesFromModFile(modFile string) []string {
	// Split the mod file into lines
	lines := strings.Split(modFile, "\n")

	// Initialize a slice to hold the require statements
	var requires []string

	// Flag to track whether we are inside a require block
	inRequireBlock := false

	// Iterate over the lines of the mod file
	for _, line := range lines {
		// If the line contains "require (", set the flag to true
		if strings.Contains(line, "require (") {
			inRequireBlock = true
			continue
		}

		// If the line contains ")", set the flag to false
		if strings.Contains(line, ")") {
			inRequireBlock = false
			continue
		}

		// If we are inside a require block, add the line to the slice of require statements
		if inRequireBlock {
			requires = append(requires, strings.TrimSpace(line))
		}
	}

	// Remove empty strings from the slice of require statements
	filteredRequires := filterEmpty(requires)

	// Initialize a buffer slice to hold the filtered libraries
	buf := make([]string, len(filteredRequires))
	bufIdx := 0 // Current index of the buffer slice

	// Iterate over the filtered require statements
	for _, lib := range filteredRequires {
		// If the statement does not contain "=>", it is a library and should be added to the buffer slice
		if !strings.Contains(lib, "=>") {
			parts := strings.Split(lib, " // ")
			buf[bufIdx] = parts[0]
			bufIdx++
		}
	}

	// Trim the buffer slice to the length of the filtered libraries
	return buf[:bufIdx]
}

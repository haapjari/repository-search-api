package goplg

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Parser struct {
}

func NewParser() *Parser {
	return new(Parser)
}

func (p *Parser) ParseRepositoryName(url string) string {
	fmt.Println(url)

	return ""
}

func (p *Parser) ParseRepositoryOwner(url string) string {
	fmt.Println(url)

	return ""
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

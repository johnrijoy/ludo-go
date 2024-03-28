package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

const Version = "0.1.0"

var (
	apiUrl          = "https://pipedapi.kavin.rocks"
	instanceListApi = "https://piped-instances.kavin.rocks"
	oldApiUrl       string
	pipedList       []PipedInstance
)

type PipedInstance struct {
	Name     string
	ApiUrl   string
	ProxyUrl string
}

func (instance PipedInstance) String() string {
	return fmt.Sprintf("%-30s| %-50s| %-50s", instance.Name, instance.ApiUrl, instance.ProxyUrl)
}

func SetPipedApi(val string) error {
	oldApiUrl = apiUrl
	apiUrl = val
	return nil
}

func GetPipedApi() string {
	return apiUrl
}

func GetOldPipedApi() string {
	return oldApiUrl
}

func GetPipedInstanceList() ([]PipedInstance, error) {
	if len(pipedList) > 0 {
		log.Println("Instance already loaded")
		return pipedList, nil
	}

	log.Println("Fetching instance list")
	resp, err := http.Get(instanceListApi)
	if err != nil {
		return nil, err
	}

	log.Println("Resp status: ", resp.Status)

	if resp.StatusCode != 200 {
		return nil, errors.New("[GetPipedInstanceList] bad response from api")
	}

	var response interface{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	instList, ok := response.([]interface{})
	if !ok {
		return nil, errors.New("response is not of expected format")
	}

	apiList := make([]PipedInstance, 0)
	for _, inst := range instList {
		var pipedInstance PipedInstance
		if name, ok := getValue(inst, path{"name"}).(string); ok {
			pipedInstance.Name = name
		}
		if url, ok := getValue(inst, path{"api_url"}).(string); ok {
			pipedInstance.ApiUrl = url
		}
		if proxy, ok := getValue(inst, path{"image_proxy_url"}).(string); ok {
			pipedInstance.ProxyUrl = proxy
		}
		apiList = append(apiList, pipedInstance)
	}

	log.Println("Instance loaded")
	return apiList, nil
}

// Helpers //

func trimList[T any](inputList []T, offset int, limit int) []T {
	outputList := inputList

	if offset > 0 && offset < len(inputList) {
		outputList = outputList[offset:]
	}

	if limit > 0 && offset >= 0 && limit < len(inputList)-offset {
		outputList = outputList[:limit]
	}

	return outputList
}

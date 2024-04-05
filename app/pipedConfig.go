package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/magiconair/properties"
)

var Piped PipedConfig

type PipedConfig struct {
	apiUrl          string
	instanceListApi string
	oldApiUrl       string
	pipedList       []PipedInstance
}

type PipedInstance struct {
	Name     string
	ApiUrl   string
	ProxyUrl string
}

func (instance PipedInstance) String() string {
	return fmt.Sprintf("%-30s| %-50s| %-50s", instance.Name, instance.ApiUrl, instance.ProxyUrl)
}

func (p *PipedConfig) SetPipedApi(val string) error {
	p.oldApiUrl = p.apiUrl
	p.apiUrl = val
	return nil
}

func (p *PipedConfig) GetPipedApi() string {
	return p.apiUrl
}

func (p *PipedConfig) GetOldPipedApi() string {
	return p.oldApiUrl
}

func (p *PipedConfig) GetPipedInstanceList() ([]PipedInstance, error) {
	if len(p.pipedList) > 0 {
		log.Println("Instance already loaded")
		return p.pipedList, nil
	}

	log.Println("Fetching instance list")
	resp, err := http.Get(p.instanceListApi)
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
	if len(apiList) > 0 {
		p.pipedList = apiList
	}

	return apiList, nil
}

func setPipedConfig(props *properties.Properties) {
	Piped.apiUrl = props.MustGetString(pipedUrlKey)
	Piped.instanceListApi = props.MustGetString(instanceListUrlKey)
}

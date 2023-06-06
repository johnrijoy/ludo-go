package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/raitonoberu/ytmusic"
)

func musicSearch(search string) string {
	searchClient := ytmusic.Search(search)
	result, err := searchClient.Next()
	checkErr(err)

	//jsonstr, _ := json.MarshalIndent(result, "", "    ")
	// fmt.Println(string(jsonstr))

	for _, val := range result.Tracks {
		arts := ""
		for _, art := range val.Artists {
			arts += art.Name
		}
		fmt.Printf("%v-%v %v\n", val.Title, arts, val.VideoID)
	}

	target := "https://pipedapi.kavin.rocks/streams/" + result.Tracks[0].VideoID

	fmt.Println(target)

	resp, err := http.Get(target)
	checkErr(err)

	// bB, err := io.ReadAll(resp.Body)
	// checkErr(err)

	var response interface{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	checkErr(err)

	trackUrl := getValue(response, path{"audioStreams", 0, "url"})
	val := trackUrl.(string)

	fmt.Println(val)
	return val
}

type path []interface{}

func getValue(source interface{}, path path) interface{} {
	value := source
	for _, element := range path {
		mustBreak := false
		switch element.(type) {
		case string:
			if val, ok := value.(map[string]interface{})[element.(string)]; ok {
				value = val
			} else {
				value = nil
				mustBreak = true
			}
		case int:
			if len(value.([]interface{})) > element.(int) {
				value = value.([]interface{})[element.(int)]
			} else {
				value = nil
				mustBreak = true
			}
		}
		if mustBreak {
			break
		}
	}
	return value
}

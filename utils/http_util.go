package utils

import (
	"bytes"
	"encoding/json"
	"game_assistantor/model"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"
)

func DoRequest(address string, method string, data map[string]interface{}) (res model.HttpResponse, err error) {
	client := http.Client{}

	var bs []byte
	var req *http.Request
	bs, err = json.Marshal(data)
	if err != nil {
		log.Info().Msgf("fail to marshal data, error is: %v", err)
		return
	}
	req, err = http.NewRequest(address, method, bytes.NewReader(bs))
	if err != nil {
		log.Info().Msgf("fail to make request, error is: %v", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Info().Msgf("fail to request data, error is: %v", err)
		return
	}

	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Info().Msgf("fail to read response data, error is: %v", err)
		return
	}

	err = json.Unmarshal(respBytes, &res)
	if err != nil {
		log.Info().Msgf("fail to decode response data, error is: %v", err)
		return
	}
	return
}

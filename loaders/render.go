package loaders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Oscar170/reverse-proxy/models"
)

var RenderApiHost = "http://localhost:3000"

func marshalComponents(components []models.Replace) string {
	requestBody := "["
	lastIndex := len(components) - 1
	for i, c := range components {
		requestBody = fmt.Sprintf(`%s{"component":"%s","props":%s}`, requestBody, c.Component.Name, c.Component.Props)
		if lastIndex != i {
			requestBody = requestBody + ","
		}
	}

	requestBody = requestBody + "]"

	return requestBody
}

// Load connect to the render api
func Load(components []models.Replace) ([]models.CompoentRendered, error) {
	resp, err := http.Post(
		RenderApiHost+"/es_ES/multirender",
		"application/json",
		bytes.NewBuffer([]byte(marshalComponents(components))),
	)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	rendered := make([]models.CompoentRendered, 0)
	err = json.Unmarshal(body, &rendered)

	return rendered, nil
}

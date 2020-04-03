package effect

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Oscar170/reverse-proxy/models"
)

var RenderApiHost = "http://localhost:3000"

// Load connect to the render api to render the components
func Load(components []models.Replace) ([]models.CompoentRendered, error) {
	requestBody, _ := json.Marshal(components)

	resp, err := http.Post(
		RenderApiHost+"/es_ES/multirender",
		"application/json",
		bytes.NewBuffer(requestBody),
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

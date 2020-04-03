package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/Oscar170/reverse-proxy/models"
	"github.com/Oscar170/reverse-proxy/parser"
)

var ServerToHydrate = "http://127.0.0.1:8080"
var RenderApiHost = "http://localhost:3000"

func findComponentsToRender(html string) []models.Replace {
	components := make([]models.Replace, 0)
	findReg := regexp.MustCompile(`@rerender\(.*\)`)
	valuesReg := regexp.MustCompile(`\"([a-zA-z]*)\"\, (\{.*\})`)

	matches := findReg.FindAllString(html, -1)

	for _, tag := range matches {
		values := valuesReg.FindStringSubmatch(tag)

		components = append(components, models.Replace{
			Tag: tag,
			Component: models.Component{
				Name:  values[1],
				Props: values[2],
			}})
	}

	return components
}

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

func renderComponents(components []models.Replace) ([]models.CompoentRendered, error) {
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

func hydrateDocument(html string, toReplace []models.Replace, renderedComponents []models.CompoentRendered) string {
	cssInline := ""
	intiState := make(map[string]string)
	for i, replace := range toReplace {
		component := renderedComponents[i]
		html = strings.Replace(html, replace.Tag, component.Html, 1)
		cssInline = cssInline + component.Css
		intiState[replace.Name] = string(component.InitState)
	}

	bInitState, _ := json.Marshal(intiState)

	html = strings.Replace(html, "<!––css_inline_hook––>", cssInline, 1)
	html = strings.Replace(
		html,
		"<!––init_state_hook––>",
		fmt.Sprintf(`<script type="application/json" id="INIT_STATE">%s</script>`, string(bInitState)),
		1,
	)

	return html
}

func rerender(html string) string {
	components := parser.ComponentsParser(html)
	renderedComponents, err := renderComponents(components)

	if err != nil {
		panic(err)
	}

	return hydrateDocument(html, components, renderedComponents)
}

func handleRequestAndRedirect(proxy *httputil.ReverseProxy) func(res http.ResponseWriter, req *http.Request) {
	proxy.ModifyResponse = func(resp *http.Response) error {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		newHtml := rerender(string(b))

		fmt.Println(newHtml)
		body := ioutil.NopCloser(strings.NewReader(newHtml))
		resp.Body = body
		resp.ContentLength = int64(len(newHtml))
		resp.Header.Set("Content-Length", strconv.Itoa(len(newHtml)))

		return nil
	}

	return proxy.ServeHTTP
}

func main() {
	url, _ := url.Parse(ServerToHydrate)
	http.HandleFunc("/", handleRequestAndRedirect(httputil.NewSingleHostReverseProxy(url)))
	http.ListenAndServe(":3030", nil)
}

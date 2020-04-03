package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/Oscar170/reverse-proxy/effect"
	"github.com/Oscar170/reverse-proxy/morphism"
)

var ServerToHydrate = "http://127.0.0.1:8080"

func rerender(html string) string {
	components := morphism.ComponentsParser(html)
	renderedComponents, err := effect.Load(components)

	if err != nil {
		panic(err)
	}

	return morphism.Hydrate(html, components, renderedComponents)
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

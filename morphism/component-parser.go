package morphism

import (
	"encoding/json"
	"regexp"

	"github.com/Oscar170/reverse-proxy/models"
)

// ComponentsParser extracts the render tags
func ComponentsParser(html string) []models.Replace {
	components := make([]models.Replace, 0)
	findReg := regexp.MustCompile(`@rerender\(.*\)`)
	valuesReg := regexp.MustCompile(`\{.*\}`)

	matches := findReg.FindAllString(html, -1)

	for _, tag := range matches {
		values := valuesReg.FindAllString(tag, 1)

		component := models.Component{}
		json.Unmarshal([]byte(values[0]), &component)

		components = append(components, models.Replace{
			Tag:       tag,
			Component: component,
		})
	}

	return components
}

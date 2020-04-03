package morphism

import (
	"regexp"

	"github.com/Oscar170/reverse-proxy/models"
)

// ComponentsParser extracts the render tags
func ComponentsParser(html string) []models.Replace {
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

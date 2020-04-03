package morphism

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Oscar170/reverse-proxy/models"
)

// Hydrate populate the html using the rendered components
func Hydrate(html string, toReplace []models.Replace, renderedComponents []models.CompoentRendered) string {
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

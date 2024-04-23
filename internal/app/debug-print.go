package app

import (
	"strings"

	"golang.org/x/net/html"
)

func NodeToString(node *html.Node) string {
	var buf strings.Builder
	html.Render(&buf, node)
	return buf.String()
}

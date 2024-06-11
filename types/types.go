package types

import (
	"regexp"
	"strings"
)

func normalise(v string) string {
	re := regexp.MustCompile(`[^a-z0-9]`)

	return re.ReplaceAllString(strings.ToLower(v), "")
}

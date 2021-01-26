package types

import (
	"strings"
)

func normalise(v string) string {
	return strings.ToLower(strings.ReplaceAll(v, " ", ""))
}

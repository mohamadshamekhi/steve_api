package slug

import (
	"regexp"
	"strings"
)

func GenerateSlug(name string) string {
	slug := strings.ReplaceAll(name, " ", "-")
	slug = strings.ReplaceAll(slug, "‌", "-")

	re := regexp.MustCompile(`[^\p{L}\p{N}-]+`)
	slug = re.ReplaceAllString(slug, "")

	re = regexp.MustCompile("-{2,}")
	slug = re.ReplaceAllString(slug, "-")

	return strings.Trim(slug, "-")
}

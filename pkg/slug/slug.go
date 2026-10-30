package slug

import (
	"regexp"
	"strings"
)

func GenerateSlug(name string) string {
	slug := digitsToASCII.Replace(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "‌", "-")

	re := regexp.MustCompile(`[^\p{L}\p{N}-]+`)
	slug = re.ReplaceAllString(slug, "")

	re = regexp.MustCompile("-{2,}")
	slug = re.ReplaceAllString(slug, "-")

	return strings.Trim(slug, "-")
}

// digitsToASCII turns Persian and Arabic-Indic digits into 0-9, so a title
// typed as "iPhone 1۷" never yields a slug that mixes digit scripts.
var digitsToASCII = strings.NewReplacer(
	"۰", "0", "۱", "1", "۲", "2", "۳", "3", "۴", "4",
	"۵", "5", "۶", "6", "۷", "7", "۸", "8", "۹", "9",
	"٠", "0", "١", "1", "٢", "2", "٣", "3", "٤", "4",
	"٥", "5", "٦", "6", "٧", "7", "٨", "8", "٩", "9",
)

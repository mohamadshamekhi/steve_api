// Package phonenumber canonicalises Iranian mobile numbers so that lookups and
// dedup checks compare apples to apples regardless of how the client typed them.
package phonenumber

import "strings"

// digitFolds maps Persian and Arabic-Indic digits to ASCII.
var digitFolds = map[rune]rune{
	'۰': '0', '۱': '1', '۲': '2', '۳': '3', '۴': '4',
	'۵': '5', '۶': '6', '۷': '7', '۸': '8', '۹': '9',
	'٠': '0', '١': '1', '٢': '2', '٣': '3', '٤': '4',
	'٥': '5', '٦': '6', '٧': '7', '٨': '8', '٩': '9',
}

// Normalize returns a canonical form of an Iranian mobile number:
//   - Persian/Arabic digits folded to ASCII
//   - spaces, dashes, dots and parentheses stripped
//   - +98 / 0098 / 98 country prefix rewritten to a leading 0
//
// A well-formed Iranian mobile becomes "09XXXXXXXXX". Input it can't confidently
// interpret is returned trimmed but otherwise untouched, so callers never end up
// comparing against an empty string.
func Normalize(raw string) string {
	var b strings.Builder
	for _, r := range strings.TrimSpace(raw) {
		if folded, ok := digitFolds[r]; ok {
			b.WriteRune(folded)
			continue
		}
		switch r {
		case ' ', '-', '.', '(', ')', '\t', '‌':
			continue
		case '+':
			continue
		default:
			b.WriteRune(r)
		}
	}
	s := b.String()

	switch {
	case strings.HasPrefix(s, "0098"):
		s = "0" + strings.TrimPrefix(s, "0098")
	case strings.HasPrefix(s, "98") && len(s) == 12:
		s = "0" + strings.TrimPrefix(s, "98")
	case strings.HasPrefix(s, "9") && len(s) == 10:
		s = "0" + s
	}
	return s
}

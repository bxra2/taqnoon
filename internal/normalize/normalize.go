// Package normalize provides string normalization helpers shared by the
// ingest pipeline and the search handlers. The semantics mirror the frontend's
// removeDiacritics() utility and the original ingest normalization functions.
package normalize

import "strings"

// RemoveDiacritics strips Arabic diacritical marks (harakat) from s.
func RemoveDiacritics(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r >= 0x064B && r <= 0x065F: // Arabic tashkeel
		case r == 0x0670: // superscript alef
		case r >= 0x06D6 && r <= 0x06ED: // other Arabic combining marks
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Latin normalizes a latin-script string: lowercase + trim.
func Latin(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// Arabic normalizes an Arabic-script string: trim + strip diacritics.
func Arabic(s string) string {
	return RemoveDiacritics(strings.TrimSpace(s))
}

// Text normalizes free-form text: trim + lowercase + strip diacritics.
func Text(s string) string {
	return RemoveDiacritics(strings.ToLower(strings.TrimSpace(s)))
}

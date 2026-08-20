package ingest

import "strings"

// removeDiacritics strips Arabic diacritical marks (harakat), mirroring the
// frontend's removeDiacritics() utility.
func removeDiacritics(s string) string {
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

func normLatin(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func normArabic(s string) string {
	return removeDiacritics(strings.TrimSpace(s))
}

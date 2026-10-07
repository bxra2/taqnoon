package ingest

import "github.com/bxra2/taqnun/internal/normalize"

// removeDiacritics strips Arabic diacritical marks (harakat), mirroring the
// frontend's removeDiacritics() utility.
func removeDiacritics(s string) string {
	return normalize.RemoveDiacritics(s)
}

func normLatin(s string) string {
	return normalize.Latin(s)
}

func normArabic(s string) string {
	return normalize.Arabic(s)
}

func normText(s string) string {
	return normalize.Text(s)
}

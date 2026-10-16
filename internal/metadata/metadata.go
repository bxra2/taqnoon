package metadata

import "regexp"

// Meta holds optional enrichment for a glossary.
type Meta struct {
	BookURL     string
	ImageURL    string
	WikidataURL string
}

var archiveDetails = regexp.MustCompile(`(?i)archive\.org/details/([^/?#]+)`)

// archiveID returns the archive.org item identifier if glossaryURL points at an
// archive.org details page, otherwise ok=false.
func archiveID(glossaryURL string) (string, bool) {
	m := archiveDetails.FindStringSubmatch(glossaryURL)
	if m == nil {
		return "", false
	}
	return m[1], true
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// Resolve returns the enrichment for a glossary URL. Precedence: explicit
// curated entry (by full URL), then curated entry by archive identifier, then
// archive.org derivation for the book/item page and cover thumbnail.
func Resolve(glossaryURL string) Meta {
	var m Meta
	if cm, ok := byURL[glossaryURL]; ok {
		m = cm
	}
	if id, ok := archiveID(glossaryURL); ok {
		if cm, ok := byID[id]; ok {
			m.BookURL = firstNonEmpty(m.BookURL, cm.BookURL)
			m.ImageURL = firstNonEmpty(m.ImageURL, cm.ImageURL)
			m.WikidataURL = firstNonEmpty(m.WikidataURL, cm.WikidataURL)
		}
		if m.BookURL == "" {
			m.BookURL = "https://archive.org/details/" + id
		}
		if m.ImageURL == "" {
			m.ImageURL = "https://archive.org/services/img/" + id
		}
	}
	return m
}

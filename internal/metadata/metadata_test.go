package metadata

import "testing"

func TestResolveArchive(t *testing.T) {
	cases := []struct {
		url  string
		book string
		img  string
	}{
		{"https://archive.org/details/ARB2010AREN/", "https://archive.org/details/ARB2010AREN", "https://archive.org/services/img/ARB2010AREN"},
		{"https://archive.org/details/SOI2000ENAR/mode/2up", "https://archive.org/details/SOI2000ENAR", "https://archive.org/services/img/SOI2000ENAR"},
		{"https://archive.org/details/DAM2018ENAR-B/mode/2up", "https://archive.org/details/DAM2018ENAR-B", "https://archive.org/services/img/DAM2018ENAR-B"},
		{"https://archive.org/details/20200515_20200515_2258/mode/2up", "https://archive.org/details/20200515_20200515_2258", "https://archive.org/services/img/20200515_20200515_2258"},
	}
	for _, c := range cases {
		m := Resolve(c.url)
		if m.BookURL != c.book {
			t.Errorf("%s: book = %q, want %q", c.url, m.BookURL, c.book)
		}
		if m.ImageURL != c.img {
			t.Errorf("%s: image = %q, want %q", c.url, m.ImageURL, c.img)
		}
	}
}

func TestCuratedWikidataResolves(t *testing.T) {
	expect := map[string]string{
		"https://archive.org/details/ARB2010AREN/":               "https://www.wikidata.org/wiki/Q2859359",
		"https://archive.org/details/ALECSO2011AREN":             "https://www.wikidata.org/wiki/Q4164494",
		"https://archive.org/details/NAH2012ENFRAR/":             "https://www.wikidata.org/wiki/Q12193946",
		"https://archive.org/details/DAM2018ENAR-B/mode/2up":     "https://www.wikidata.org/wiki/Q108593221",
		"https://sdaia.gov.sa/ar/MediaCenter/KnowledgeCenter/ResearchLibrary/SDAIAPublications15.pdf": "https://www.wikidata.org/wiki/Q66982800",
	}
	for url, want := range expect {
		got := Resolve(url).WikidataURL
		if got != want {
			t.Errorf("%s: wikidata = %q, want %q", url, got, want)
		}
	}
}

func TestResolveNonArchive(t *testing.T) {
	m := Resolve("https://arabterm.org/index.php?id=41&L=1&tx_3m5techdict_pi1[filterCategory]=1")
	if m.BookURL != "" || m.ImageURL != "" || m.WikidataURL != "" {
		t.Errorf("non-archive URL should have no enrichment, got %+v", m)
	}
}

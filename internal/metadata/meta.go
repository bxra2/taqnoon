package metadata

// Curated enrichment keyed by the exact glossary URL (for non-archive sources)
// and by archive.org item identifier (for archive sources). Keep keys exact.
var byURL = map[string]Meta{
	"https://sdaia.gov.sa/ar/MediaCenter/KnowledgeCenter/ResearchLibrary/SDAIAPublications15.pdf": {
		WikidataURL: "https://www.wikidata.org/wiki/Q66982800",
	},
}

var byID = map[string]Meta{
	"ARB2010AREN": { // ArabEyes القاموس التقني المعلوماتي
		WikidataURL: "https://www.wikidata.org/wiki/Q2859359",
	},
	"ALECSO2011AREN": { // معجم تقانة المعلومات الموحد (ALECSO)
		WikidataURL: "https://www.wikidata.org/wiki/Q4164494",
	},
	"NAH2012ENFRAR": { // المنظمة العربية للترجمة
		WikidataURL: "https://www.wikidata.org/wiki/Q12193946",
	},
	"DAM2018ENAR-B": { // معجم مصطلحات الرياضيات (مجمع دمشق 2018)
		WikidataURL: "https://www.wikidata.org/wiki/Q108593221",
	},
}

package localizations

type LocalizationSource struct {
	URL         string `json:"url"`
	FileType    string `json:"fileType"` // "po" (default) or "txt"
	ProjectEn   string `json:"projectEn"`
	ProjectAr   string `json:"projectAr"`
	PublisherEn string `json:"publisherEn"`
	PublisherAr string `json:"publisherAr"`
}

type PublisherSources struct {
	PublisherEn string           `json:"publisherEn"`
	PublisherAr string           `json:"publisherAr"`
	Projects    []ProjectSource `json:"projects"`
}

type ProjectSource struct {
	URL       string `json:"url"`
	ProjectEn string `json:"projectEn"`
	ProjectAr string `json:"projectAr"`
}

// SourcesByPublisher returns Sources grouped by publisher, preserving the
// original order of first appearance for both publishers and their projects.
func SourcesByPublisher() []PublisherSources {
	idx := make(map[string]int, len(Sources))
	out := make([]PublisherSources, 0)
	for _, s := range Sources {
		i, ok := idx[s.PublisherEn]
		if !ok {
			out = append(out, PublisherSources{
				PublisherEn: s.PublisherEn,
				PublisherAr: s.PublisherAr,
			})
			i = len(out) - 1
			idx[s.PublisherEn] = i
		}
		out[i].Projects = append(out[i].Projects, ProjectSource{
			URL:       s.URL,
			ProjectEn: s.ProjectEn,
			ProjectAr: s.ProjectAr,
		})
	}
	return out
}

var Sources = []LocalizationSource{
	// GNOME
	{URL: "https://gitlab.gnome.org/GNOME/gnome-shell/-/raw/main/po/ar.po", FileType: "po", ProjectAr: "صدفة جنوم", ProjectEn: "gnome-shell", PublisherAr: "جنوم", PublisherEn: "GNOME"},
	{URL: "https://gitlab.gnome.org/GNOME/gtk/-/raw/main/po/ar.po", FileType: "po", ProjectAr: "جي تي ك", ProjectEn: "gtk", PublisherAr: "جنوم", PublisherEn: "GNOME"},
	{URL: "https://gitlab.gnome.org/GNOME/nautilus/-/raw/main/po/ar.po", FileType: "po", ProjectAr: "نوتيلاس", ProjectEn: "nautilus", PublisherAr: "جنوم", PublisherEn: "GNOME"},
	{URL: "https://gitlab.gnome.org/GNOME/gnome-control-center/-/raw/main/po/ar.po", FileType: "po", ProjectAr: "وحدة تحكم جنوم", ProjectEn: "gnome-control-center", PublisherAr: "جنوم", PublisherEn: "GNOME"},
	{URL: "https://gitlab.gnome.org/GNOME/gnome-session/-/raw/main/po/ar.po", FileType: "po", ProjectAr: "جنوم سيشن", ProjectEn: "gnome-session", PublisherAr: "جنوم", PublisherEn: "GNOME"},
	{URL: "https://gitlab.gnome.org/GNOME/gimp/-/raw/master/po/ar.po", FileType: "po", ProjectAr: "جيمب", ProjectEn: "gimp", PublisherAr: "جنوم", PublisherEn: "GNOME"},
	{URL: "https://gitlab.gnome.org/GNOME/epiphany/-/raw/main/po/ar.po", FileType: "po", ProjectAr: "إبيفاني", ProjectEn: "epiphany", PublisherAr: "جنوم", PublisherEn: "GNOME"},

	// KDE
	{URL: "https://raw.githubusercontent.com/KDE/krusader/refs/heads/master/po/ar/krusader.po", FileType: "po", ProjectAr: "كروسيدر", ProjectEn: "krusader", PublisherAr: "كيدي", PublisherEn: "KDE"},
	{URL: "https://raw.githubusercontent.com/KDE/kcalc/refs/heads/master/po/ar/kcalc.po", FileType: "po", ProjectAr: "كالك", ProjectEn: "kcalc", PublisherAr: "كيدي", PublisherEn: "KDE"},
	{URL: "https://raw.githubusercontent.com/KDE/konsole/refs/heads/master/po/ar/konsole.po", FileType: "po", ProjectAr: "كونسل", ProjectEn: "konsole", PublisherAr: "كيدي", PublisherEn: "KDE"},
	{URL: "https://raw.githubusercontent.com/KDE/krita/refs/heads/master/po/ar/krita.po", FileType: "po", ProjectAr: "كريتا", ProjectEn: "krita", PublisherAr: "كيدي", PublisherEn: "KDE"},
	{URL: "https://raw.githubusercontent.com/KDE/dolphin/refs/heads/master/po/ar/dolphin.po", FileType: "po", ProjectAr: "دولفين", ProjectEn: "dolphin", PublisherAr: "كيدي", PublisherEn: "KDE"},
	{URL: "https://raw.githubusercontent.com/KDE/okular/refs/heads/master/po/ar/okular.po", FileType: "po", ProjectAr: "أوكيلار", ProjectEn: "okular", PublisherAr: "كيدي", PublisherEn: "KDE"},
	{URL: "https://raw.githubusercontent.com/KDE/kdenlive/refs/heads/master/po/ar/kdenlive.po", FileType: "po", ProjectAr: "ك دن لايف", ProjectEn: "kdenlive", PublisherAr: "كيدي", PublisherEn: "KDE"},
	{URL: "https://raw.githubusercontent.com/KDE/kate/refs/heads/master/po/ar/kate.po", FileType: "po", ProjectAr: "كات", ProjectEn: "kate", PublisherAr: "كيدي", PublisherEn: "KDE"},

	// Cinnamon
	{URL: "https://raw.githubusercontent.com/linuxmint/cinnamon-translations/refs/heads/master/po-export/cinnamon/cinnamon-ar.po", ProjectAr: "سينامن", ProjectEn: "cinnamon", PublisherAr: "لينكس مينت", PublisherEn: "Linux Mint"},
	{URL: "https://raw.githubusercontent.com/linuxmint/cinnamon-translations/refs/heads/master/po-export/nemo/nemo-ar.po", ProjectAr: "نيمو", ProjectEn: "nemo", PublisherAr: "لينكس مينت", PublisherEn: "Linux Mint"},
	{URL: "https://raw.githubusercontent.com/linuxmint/cinnamon-translations/refs/heads/master/po-export/cinnamon-control-center/cinnamon-control-center-ar.po", ProjectAr: "وحدة تحكم سينامن", ProjectEn: "cinnamon control center", PublisherAr: "لينكس مينت", PublisherEn: "Linux Mint"},

	// Obsidian
	{URL: "https://raw.githubusercontent.com/obsidianmd/obsidian-translations/master/translations/ar.txt", FileType: "txt", ProjectAr: "أوبسيديان", ProjectEn: "Obsidian", PublisherAr: "أوبسيديان", PublisherEn: "Obsidian"},

	// Audacity
	{URL: "https://raw.githubusercontent.com/audacity/audacity/refs/heads/master/au3/locale/ar.po", FileType: "po", ProjectAr: "أوداسيتي", ProjectEn: "Audacity", PublisherAr: "أوداسيتي", PublisherEn: "Audacity"},
}

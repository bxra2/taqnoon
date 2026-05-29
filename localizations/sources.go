package localizations

type LocalizationSource struct {
	URL         string
	Type        string // "po" (default) or "txt"
	ProjectEn   string
	ProjectAr   string
	PublisherEn string
	PublisherAr string
}

var Sources = []LocalizationSource{
	// GNOME
	{URL: "https://gitlab.gnome.org/GNOME/gnome-shell/-/raw/main/po/ar.po", ProjectAr: "صدفة جنوم", ProjectEn: "gnome-shell", PublisherAr: "جنوم", PublisherEn: "GNOME"},
	{URL: "https://gitlab.gnome.org/GNOME/gtk/-/raw/main/po/ar.po", ProjectAr: "جي تي ك", ProjectEn: "gtk", PublisherAr: "جنوم", PublisherEn: "GNOME"},
	{URL: "https://gitlab.gnome.org/GNOME/nautilus/-/raw/main/po/ar.po", ProjectAr: "نوتيلاس", ProjectEn: "nautilus", PublisherAr: "جنوم", PublisherEn: "GNOME"},
	{URL: "https://gitlab.gnome.org/GNOME/gnome-control-center/-/raw/main/po/ar.po", ProjectAr: "وحدة تحكم جنوم", ProjectEn: "gnome-control-center", PublisherAr: "جنوم", PublisherEn: "GNOME"},
	{URL: "https://gitlab.gnome.org/GNOME/gnome-session/-/raw/main/po/ar.po", ProjectAr: "جنوم سيشن", ProjectEn: "gnome-session", PublisherAr: "جنوم", PublisherEn: "GNOME"},
	{URL: "https://gitlab.gnome.org/GNOME/gimp/-/raw/master/po/ar.po", ProjectAr: "جيمب", ProjectEn: "gimp", PublisherAr: "جنوم", PublisherEn: "GNOME"},
	{URL: "https://gitlab.gnome.org/GNOME/epiphany/-/raw/main/po/ar.po", Type: "po", ProjectAr: "إبيفاني", ProjectEn: "epiphany", PublisherAr: "جنوم", PublisherEn: "GNOME"},

	// KDE
	{URL: "https://raw.githubusercontent.com/KDE/krusader/refs/heads/master/po/ar/krusader.po", ProjectAr: "كروسيدر", ProjectEn: "krusader", PublisherAr: "كيدي", PublisherEn: "KDE"},
	{URL: "https://raw.githubusercontent.com/KDE/kcalc/refs/heads/master/po/ar/kcalc.po", ProjectAr: "كالك", ProjectEn: "kcalc", PublisherAr: "كيدي", PublisherEn: "KDE"},
	{URL: "https://raw.githubusercontent.com/KDE/konsole/refs/heads/master/po/ar/konsole.po", ProjectAr: "كونسل", ProjectEn: "konsole", PublisherAr: "كيدي", PublisherEn: "KDE"},
	{URL: "https://raw.githubusercontent.com/KDE/krita/refs/heads/master/po/ar/krita.po", ProjectAr: "كريتا", ProjectEn: "krita", PublisherAr: "كيدي", PublisherEn: "KDE"},
	{URL: "https://raw.githubusercontent.com/KDE/dolphin/refs/heads/master/po/ar/dolphin.po", ProjectAr: "دولفين", ProjectEn: "dolphin", PublisherAr: "كيدي", PublisherEn: "KDE"},
	{URL: "https://raw.githubusercontent.com/KDE/okular/refs/heads/master/po/ar/okular.po", ProjectAr: "أوكيلار", ProjectEn: "okular", PublisherAr: "كيدي", PublisherEn: "KDE"},
	{URL: "https://raw.githubusercontent.com/KDE/kdenlive/refs/heads/master/po/ar/kdenlive.po", ProjectAr: "ك دن لايف", ProjectEn: "kdenlive", PublisherAr: "كيدي", PublisherEn: "KDE"},
	{URL: "https://raw.githubusercontent.com/KDE/kate/refs/heads/master/po/ar/kate.po", ProjectAr: "كات", ProjectEn: "kate", PublisherAr: "كيدي", PublisherEn: "KDE"},

	// Cinnamon
	{URL: "https://raw.githubusercontent.com/linuxmint/cinnamon-translations/refs/heads/master/po-export/cinnamon/cinnamon-ar.po", ProjectAr: "سينامن", ProjectEn: "cinnamon", PublisherAr: "لينكس مينت", PublisherEn: "Linux Mint"},
	{URL: "https://raw.githubusercontent.com/linuxmint/cinnamon-translations/refs/heads/master/po-export/nemo/nemo-ar.po", ProjectAr: "نيمو", ProjectEn: "nemo", PublisherAr: "لينكس مينت", PublisherEn: "Linux Mint"},
	{URL: "https://raw.githubusercontent.com/linuxmint/cinnamon-translations/refs/heads/master/po-export/cinnamon-control-center/cinnamon-control-center-ar.po", ProjectAr: "وحدة تحكم سينامن", ProjectEn: "cinnamon control center", PublisherAr: "لينكس مينت", PublisherEn: "Linux Mint"},

	// Obsidian
	{URL: "https://raw.githubusercontent.com/obsidianmd/obsidian-translations/master/translations/ar.txt", Type: "txt", ProjectAr: "أوبسيديان", ProjectEn: "Obsidian", PublisherAr: "أوبسيديان", PublisherEn: "Obsidian"},

	// Audacity
	{URL: "https://raw.githubusercontent.com/audacity/audacity/refs/heads/master/au3/locale/ar.po", Type: "po", ProjectAr: "أوداسيتي", ProjectEn: "Audacity", PublisherAr: "أوداسيتي", PublisherEn: "Audacity"},
}

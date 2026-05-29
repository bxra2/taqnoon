package localizations

type LocalizationSource struct {
	url         string
	fileType    string // "po" (default) or "txt"
	projectEn   string
	projectAr   string
	publisherEn string
	publisherAr string
}

var Sources = []LocalizationSource{
	// GNOME
	{url: "https://gitlab.gnome.org/GNOME/gnome-shell/-/raw/main/po/ar.po", fileType: "po", projectAr: "صدفة جنوم", projectEn: "gnome-shell", publisherAr: "جنوم", publisherEn: "GNOME"},
	{url: "https://gitlab.gnome.org/GNOME/gtk/-/raw/main/po/ar.po", fileType: "po", projectAr: "جي تي ك", projectEn: "gtk", publisherAr: "جنوم", publisherEn: "GNOME"},
	{url: "https://gitlab.gnome.org/GNOME/nautilus/-/raw/main/po/ar.po", fileType: "po", projectAr: "نوتيلاس", projectEn: "nautilus", publisherAr: "جنوم", publisherEn: "GNOME"},
	{url: "https://gitlab.gnome.org/GNOME/gnome-control-center/-/raw/main/po/ar.po", fileType: "po", projectAr: "وحدة تحكم جنوم", projectEn: "gnome-control-center", publisherAr: "جنوم", publisherEn: "GNOME"},
	{url: "https://gitlab.gnome.org/GNOME/gnome-session/-/raw/main/po/ar.po", fileType: "po", projectAr: "جنوم سيشن", projectEn: "gnome-session", publisherAr: "جنوم", publisherEn: "GNOME"},
	{url: "https://gitlab.gnome.org/GNOME/gimp/-/raw/master/po/ar.po", fileType: "po", projectAr: "جيمب", projectEn: "gimp", publisherAr: "جنوم", publisherEn: "GNOME"},
	{url: "https://gitlab.gnome.org/GNOME/epiphany/-/raw/main/po/ar.po", fileType: "po", projectAr: "إبيفاني", projectEn: "epiphany", publisherAr: "جنوم", publisherEn: "GNOME"},

	// KDE
	{url: "https://raw.githubusercontent.com/KDE/krusader/refs/heads/master/po/ar/krusader.po", fileType: "po", projectAr: "كروسيدر", projectEn: "krusader", publisherAr: "كيدي", publisherEn: "KDE"},
	{url: "https://raw.githubusercontent.com/KDE/kcalc/refs/heads/master/po/ar/kcalc.po", fileType: "po", projectAr: "كالك", projectEn: "kcalc", publisherAr: "كيدي", publisherEn: "KDE"},
	{url: "https://raw.githubusercontent.com/KDE/konsole/refs/heads/master/po/ar/konsole.po", fileType: "po", projectAr: "كونسل", projectEn: "konsole", publisherAr: "كيدي", publisherEn: "KDE"},
	{url: "https://raw.githubusercontent.com/KDE/krita/refs/heads/master/po/ar/krita.po", fileType: "po", projectAr: "كريتا", projectEn: "krita", publisherAr: "كيدي", publisherEn: "KDE"},
	{url: "https://raw.githubusercontent.com/KDE/dolphin/refs/heads/master/po/ar/dolphin.po", fileType: "po", projectAr: "دولفين", projectEn: "dolphin", publisherAr: "كيدي", publisherEn: "KDE"},
	{url: "https://raw.githubusercontent.com/KDE/okular/refs/heads/master/po/ar/okular.po", fileType: "po", projectAr: "أوكيلار", projectEn: "okular", publisherAr: "كيدي", publisherEn: "KDE"},
	{url: "https://raw.githubusercontent.com/KDE/kdenlive/refs/heads/master/po/ar/kdenlive.po", fileType: "po", projectAr: "ك دن لايف", projectEn: "kdenlive", publisherAr: "كيدي", publisherEn: "KDE"},
	{url: "https://raw.githubusercontent.com/KDE/kate/refs/heads/master/po/ar/kate.po", fileType: "po", projectAr: "كات", projectEn: "kate", publisherAr: "كيدي", publisherEn: "KDE"},

	// Cinnamon
	{url: "https://raw.githubusercontent.com/linuxmint/cinnamon-translations/refs/heads/master/po-export/cinnamon/cinnamon-ar.po", projectAr: "سينامن", projectEn: "cinnamon", publisherAr: "لينكس مينت", publisherEn: "Linux Mint"},
	{url: "https://raw.githubusercontent.com/linuxmint/cinnamon-translations/refs/heads/master/po-export/nemo/nemo-ar.po", projectAr: "نيمو", projectEn: "nemo", publisherAr: "لينكس مينت", publisherEn: "Linux Mint"},
	{url: "https://raw.githubusercontent.com/linuxmint/cinnamon-translations/refs/heads/master/po-export/cinnamon-control-center/cinnamon-control-center-ar.po", projectAr: "وحدة تحكم سينامن", projectEn: "cinnamon control center", publisherAr: "لينكس مينت", publisherEn: "Linux Mint"},

	// Obsidian
	{url: "https://raw.githubusercontent.com/obsidianmd/obsidian-translations/master/translations/ar.txt", fileType: "txt", projectAr: "أوبسيديان", projectEn: "Obsidian", publisherAr: "أوبسيديان", publisherEn: "Obsidian"},

	// Audacity
	{url: "https://raw.githubusercontent.com/audacity/audacity/refs/heads/master/au3/locale/ar.po", fileType: "po", projectAr: "أوداسيتي", projectEn: "Audacity", publisherAr: "أوداسيتي", publisherEn: "Audacity"},
}

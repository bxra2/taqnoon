package normalize

import (
	"strings"
	"testing"
)

func TestRemoveDiacritics(t *testing.T) {
	cases := map[string]string{
		"حاسوب":       "حاسوب",
		"حَاسُوب":     "حاسوب",
		"قَامُوسٌ":    "قاموس",
		"بِسْمِ اللَّهِ": "بسم الله",
	}
	for in, want := range cases {
		if got := RemoveDiacritics(in); got != want {
			t.Errorf("RemoveDiacritics(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestLatin(t *testing.T) {
	cases := map[string]string{
		"Computer":      "computer",
		"  COMPUTER  ":  "computer",
		"\tDatabase\n":  "database",
	}
	for in, want := range cases {
		if got := Latin(in); got != want {
			t.Errorf("Latin(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestArabic(t *testing.T) {
	cases := map[string]string{
		"حاسوب":      "حاسوب",
		" حَاسُوب ":  "حاسوب",
		"  قَامُوسٌ": "قاموس",
	}
	for in, want := range cases {
		if got := Arabic(in); got != want {
			t.Errorf("Arabic(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestText(t *testing.T) {
	cases := map[string]string{
		"Computer":          "computer",
		"  Computer Sci  ":  "computer sci",
		"حَاسُوب  آلِي":     "حاسوب  آلي",
	}
	for in, want := range cases {
		if got := Text(in); got != want {
			t.Errorf("Text(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestParityWithOldIngest ensures the exported helpers still match the
// original ingest semantics (lowercase+trim for latin, strip+trim for arabic).
func TestParityWithOldIngest(t *testing.T) {
	oldLatin := func(s string) string { return strings.ToLower(strings.TrimSpace(s)) }
	oldArabic := func(s string) string { return RemoveDiacritics(strings.TrimSpace(s)) }

	inputs := []string{
		"  Computer  ",
		"\tDATA\n",
		" حَاسُوب ",
		"  قَامُوسٌ  ",
		"مُصْطَلَح",
	}

	for _, in := range inputs {
		if got, want := Latin(in), oldLatin(in); got != want {
			t.Errorf("Latin parity %q: got %q, want %q", in, got, want)
		}
		if got, want := Arabic(in), oldArabic(in); got != want {
			t.Errorf("Arabic parity %q: got %q, want %q", in, got, want)
		}
	}
}

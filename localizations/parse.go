package localizations

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Entry struct {
	ID          string `json:"id"`
	English     string `json:"english"`
	Arabic      string `json:"arabic"`
	GlossaryEn  string `json:"glossaryEn"`
	GlossaryAr  string `json:"glossaryAr"`
	GlossaryURL string `json:"glossaryUrl"`
	PublisherAr string `json:"publisherAr"`
	PublisherEn string `json:"publisherEn"`
	TURL        string `json:"tURL"`
}

var (
	msgidRe  = regexp.MustCompile(`msgid\s+"(.*)"`)
	msgstrRe = regexp.MustCompile(`msgstr\s+"(.*)"`)
)

func parsePo(text string) []struct{ english, arabic string } {
	var out []struct{ english, arabic string }
	var msgid, msgstr string

	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "msgid") {
			if m := msgidRe.FindStringSubmatch(line); m != nil {
				msgid = m[1]
			} else {
				msgid = ""
			}
		}
		if strings.HasPrefix(line, "msgstr") {
			if m := msgstrRe.FindStringSubmatch(line); m != nil {
				msgstr = m[1]
			} else {
				msgstr = ""
			}
			if msgid != "" {
				out = append(out, struct{ english, arabic string }{msgid, msgstr})
				msgid, msgstr = "", ""
			}
		}
	}
	return out
}

func parseTxt(text string) []struct{ english, arabic string } {
	var out []struct{ english, arabic string }
	var original, translation string

	flush := func() {
		if original != "" {
			out = append(out, struct{ english, arabic string }{original, translation})
		}
		original, translation = "", ""
	}

	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "[") {
			flush()
			continue
		}
		switch {
		case strings.HasPrefix(line, "original="):
			original = strings.TrimPrefix(line, "original=")
		case strings.HasPrefix(line, "translation="):
			translation = strings.TrimPrefix(line, "translation=")
		}
	}
	flush()
	return out
}

func fetchAndParse(ctx context.Context, src LocalizationSource, counter *uint64) ([]Entry, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, src.URL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d for %s", resp.StatusCode, src.URL)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var raw []struct{ english, arabic string }
	if src.FileType == "txt" {
		raw = parseTxt(string(body))
	} else {
		raw = parsePo(string(body))
	}

	entries := make([]Entry, 0, len(raw))
	for _, p := range raw {
		id := atomic.AddUint64(counter, 1)
		entries = append(entries, Entry{
			ID:          fmt.Sprintf("%s%s%d", src.PublisherEn, p.english, id),
			English:     p.english,
			Arabic:      p.arabic,
			GlossaryEn:  src.ProjectEn,
			GlossaryAr:  src.ProjectAr,
			GlossaryURL: src.URL,
			PublisherAr: src.PublisherAr,
			PublisherEn: src.PublisherEn,
			TURL:        src.URL,
		})
	}
	return entries, nil
}

const (
	concurrency = 4
	cacheTTL    = 6 * time.Hour
)

type cache struct {
	mu       sync.Mutex
	value    []Entry
	at       time.Time
	inflight chan struct{}
	err      error
}

var localizationsCache = &cache{}

func FetchAll(ctx context.Context) ([]Entry, error) {
	localizationsCache.mu.Lock()
	if localizationsCache.value != nil && time.Since(localizationsCache.at) < cacheTTL {
		v := localizationsCache.value
		localizationsCache.mu.Unlock()
		return v, nil
	}
	if localizationsCache.inflight != nil {
		ch := localizationsCache.inflight
		localizationsCache.mu.Unlock()
		<-ch
		localizationsCache.mu.Lock()
		v, err := localizationsCache.value, localizationsCache.err
		localizationsCache.mu.Unlock()
		return v, err
	}
	done := make(chan struct{})
	localizationsCache.inflight = done
	localizationsCache.mu.Unlock()

	value, err := fetchAllNow(ctx)

	localizationsCache.mu.Lock()
	if err == nil {
		localizationsCache.value = value
		localizationsCache.at = time.Now()
		localizationsCache.err = nil
	} else {
		localizationsCache.err = err
	}
	localizationsCache.inflight = nil
	localizationsCache.mu.Unlock()
	close(done)

	return value, err
}

func fetchAllNow(ctx context.Context) ([]Entry, error) {
	var counter uint64
	results := make([][]Entry, len(Sources))

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	for i, src := range Sources {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, src LocalizationSource) {
			defer wg.Done()
			defer func() { <-sem }()
			entries, err := fetchAndParse(ctx, src, &counter)
			if err != nil {
				log.Printf("localizations: fetch failed for %s: %v", src.URL, err)
				return
			}
			results[i] = entries
		}(i, src)
	}
	wg.Wait()

	total := 0
	for _, r := range results {
		total += len(r)
	}
	merged := make([]Entry, 0, total)
	for _, r := range results {
		merged = append(merged, r...)
	}
	return merged, nil
}

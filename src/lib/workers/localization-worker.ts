import { removeDiacritics } from '$lib/utils/removeDiacritics';
import { fetchAndParseAllPoFiles } from '../utils/fetchAllPoFiles';
let entries: { english: string; arabic: string }[] = [];
let isLoaded = false

self.onmessage = async (e) => {
    if (e.data.type === 'load') {
        if (isLoaded) {
            self.postMessage({ type: 'ready', count: entries.length });
            return;
        }
        try {
            const allEntries = await fetchAndParseAllPoFiles();
            entries = allEntries;
            isLoaded = true
            self.postMessage({ type: 'ready', count: entries.length })
        } catch (error) {
            self.postMessage({ type: 'error', message: 'فشل تحميل الترجمات' })
        }
        return
    }

    if (e.data.type === 'search') {
        if (!isLoaded) {
            self.postMessage({ type: 'error', message: 'الترجمات لم تحمل بعد' })
            return
        }
        const lower = e.data.query.toLowerCase();
        const normalizedQuery = removeDiacritics(e.data.query);

        let results;

        if (e.data.exact) {
            results = entries.filter(
                (entry) =>
                    entry.english.toLowerCase() === lower ||
                    removeDiacritics(entry.arabic) === normalizedQuery
            )
        } else {
            const exactMatches = entries.filter(
                (entry) =>
                    entry.english.toLowerCase() === lower ||
                    removeDiacritics(entry.arabic) === normalizedQuery
            )

            const partialMatches = entries.filter(
                (entry) =>
                    !exactMatches.includes(entry) && (
                        entry.english.toLowerCase().includes(lower) ||
                        removeDiacritics(entry.arabic).includes(normalizedQuery)
                    )
            );

            results = [...exactMatches, ...partialMatches]
        }

        self.postMessage({ type: 'results', results })
    }
}

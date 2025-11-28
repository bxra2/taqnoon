import { removeDiacritics } from '$lib/utils/removeDiacritics';
import { fetchAndParseAllPoFiles } from '../utils/fetchAllPoFiles';
let entries: { english: string; arabic: string }[] = [];

self.onmessage = async (e) => {
    if (e.data.type === 'load') {
        const allEntries = await fetchAndParseAllPoFiles();
        entries = allEntries;

        self.postMessage({ type: 'loaded', count: entries.length });
    }

    if (e.data.type === 'search') {
        const lower = e.data.query.toLowerCase();
        const normalizedQuery = removeDiacritics(e.data.query);

        let results;

        if (e.data.exact) {
            results = entries.filter(
                (entry) =>
                    entry.english.toLowerCase() === lower ||
                    entry.arabic.toLowerCase() === normalizedQuery
            )
        } else {
            const exactMatches = entries.filter(
                (entry) =>
                    entry.english.toLowerCase() === lower ||
                    entry.arabic.toLowerCase() === normalizedQuery
            )

            const partialMatches = entries.filter(
                (entry) =>
                    !exactMatches.includes(entry) && (
                        entry.english.toLowerCase().includes(lower) ||
                        entry.arabic.toLowerCase().includes(normalizedQuery)
                    )
            );

            results = [...exactMatches, ...partialMatches]
        }

        self.postMessage({ type: 'results', results })
    }
};

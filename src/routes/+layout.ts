// src/routes/(app)/loadData.js
export const prerender = false;
export const trailingSlash = 'always';

export async function load() {
    try {
        // 1️⃣ Fetch the index file listing all publisher JSON files
        const indexRes = await fetch('/data/publishers/index.json');
        const files = await indexRes.json(); // ["scs/scs-1.json", "scs/scs-2.json", ...]

        const allData = [];
        const publisherMap = new Map();

        // 2️⃣ Fetch each JSON file dynamically
        for (const file of files) {
            const res = await fetch(`/data/publishers/${file}`);
            const json = await res.json();

            if (!json.fileData || !Array.isArray(json.entries)) continue;

            const { glossaryEn, glossaryAr, glossaryUrl, publisherEn, publisherAr, publisherUrl } = json.fileData;

            for (const entry of json.entries) {
                const fullEntry = {
                    ...entry,
                    glossaryEn,
                    glossaryAr,
                    glossaryUrl,
                    publisherEn,
                    publisherAr,
                    publisherUrl,
                };
                allData.push(fullEntry);

                const pubKey = publisherEn || publisherAr;
                if (!pubKey) continue;

                if (!publisherMap.has(pubKey)) {
                    publisherMap.set(pubKey, {
                        publisherEn,
                        publisherAr,
                        publisherUrl,
                        glossaries: new Map(),
                    });
                }

                const publisher = publisherMap.get(pubKey);
                const glossaryKey = glossaryEn || glossaryAr;

                if (glossaryKey && !publisher.glossaries.has(glossaryKey)) {
                    publisher.glossaries.set(glossaryKey, { glossaryEn, glossaryAr, glossaryUrl });
                }
            }
        }

        // 3️⃣ Convert maps to arrays for easier usage in Svelte components
        const termData = allData;
        const publishersData = Array.from(publisherMap.values()).map((pub) => ({
            ...pub,
            glossaries: Array.from(pub.glossaries.values()),
        }));

        return { termData, publishersData };
    } catch (error) {
        console.error('Error loading data:', error);
        return { termData: [], publishersData: [] };
    }
}
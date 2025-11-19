export async function load() {
    try {
        const modules = import.meta.glob('/src/lib/data/**/*.json', { eager: true });
        const allData = [];
        const publisherMap = new Map();

        for (const path in modules) {
            const json = modules[path].default;

            if (!json.fileData || !Array.isArray(json.entries)) continue;

            const fileMeta = json.fileData;
            const {
                glossaryEn,
                glossaryAr,
                glossaryUrl,
                publisherEn,
                publisherAr,
                publisherUrl,
            } = fileMeta;

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
                    publisher.glossaries.set(glossaryKey, {
                        glossaryEn,
                        glossaryAr,
                        glossaryUrl,
                    });
                }
            }
        }

        const termData = allData;
        const publishersData = Array.from(publisherMap.values()).map((pub) => ({
            ...pub,
            glossaries: Array.from(pub.glossaries.values()),
        }));

        return { termData, publishersData };
    } catch (error) {
        console.error('Error loading data:', error);
        // Return empty data instead of undefined to prevent errors in components
        return {
            termData: [],
            publishersData: []
        };
    }
}
export const ssr = false;
export const prerender = false;

export async function load({ fetch }) {
    try {
        const [publishersRes, localizationSourcesRes] = await Promise.all([
            fetch('/api/publishers'),
            fetch('/api/localizations/sources')
        ]);
        if (!publishersRes.ok) throw new Error(`publishers ${publishersRes.status}`);
        if (!localizationSourcesRes.ok) throw new Error(`localizations ${localizationSourcesRes.status}`);
        const [publishersData, localizationSourcesData] = await Promise.all([
            publishersRes.json(),
            localizationSourcesRes.json()
        ]);
        return { publishersData, localizationSourcesData };
    } catch (error) {
        console.error('Error loading data from API:', error);
        return { publishersData: [], localizationSourcesData: [] };
    }
}

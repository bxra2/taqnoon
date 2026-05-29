export const ssr = false;
export const prerender = false;

export async function load({ fetch }) {
    try {
        const [termsRes, publishersRes, localizationSourcesRes] = await Promise.all([
            fetch('/api/terms'),
            fetch('/api/publishers'),
            fetch('/api/localizations/sources')
        ]);
        if (!termsRes.ok) throw new Error(`terms ${termsRes.status}`);
        if (!publishersRes.ok) throw new Error(`publishers ${publishersRes.status}`);
        if (!localizationSourcesRes.ok) throw new Error(`localizations ${localizationSourcesRes.status}`);
        const [termData, publishersData, localizationSourcesData] = await Promise.all([
            termsRes.json(),
            publishersRes.json(),
            localizationSourcesRes.json()
        ]);
        return { termData, publishersData, localizationSourcesData };
    } catch (error) {
        console.error('Error loading data from API:', error);
        return { termData: [], publishersData: [], localizationSourcesData: [] };
    }
}

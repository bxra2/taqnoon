export const ssr = false;
export const prerender = false;

export async function load({ fetch }) {
    try {
        const [termsRes, publishersRes] = await Promise.all([
            fetch('/api/terms'),
            fetch('/api/publishers')
        ]);
        if (!termsRes.ok) throw new Error(`terms ${termsRes.status}`);
        if (!publishersRes.ok) throw new Error(`publishers ${publishersRes.status}`);
        const [termData, publishersData] = await Promise.all([
            termsRes.json(),
            publishersRes.json()
        ]);
        return { termData, publishersData };
    } catch (error) {
        console.error('Error loading data from API:', error);
        return { termData: [], publishersData: [] };
    }
}

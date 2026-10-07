export const ssr = false;

export async function load({ fetch }) {
    try {
        const res = await fetch('/api/terms');
        if (!res.ok) throw new Error(`terms ${res.status}`);
        const termData = await res.json();
        return { termData };
    } catch (error) {
        console.error('Error loading terms from API:', error);
        return { termData: [] };
    }
}

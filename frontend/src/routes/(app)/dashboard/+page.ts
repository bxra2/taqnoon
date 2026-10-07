export const ssr = false;

export async function load({ fetch }) {
    try {
        const res = await fetch('/api/stats');
        if (!res.ok) throw new Error(`stats ${res.status}`);
        const stats = await res.json();
        return { stats };
    } catch (error) {
        console.error('Error loading stats from API:', error);
        return { stats: null };
    }
}

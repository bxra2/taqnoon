export const ssr = false;
export const prerender = false;

export async function load({ fetch }) {
    try {
        const res = await fetch('/api/suggest');
        if (!res.ok) throw new Error(`suggest ${res.status}`);
        const suggestions = await res.json();
        return { suggestions };
    } catch (error) {
        console.error('Error loading suggestions from API:', error);
        return { suggestions: [] };
    }
}

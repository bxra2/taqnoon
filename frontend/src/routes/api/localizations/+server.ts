import { json } from '@sveltejs/kit';
import { getCachedLocalizations } from '$lib/server/localizations';

export async function GET() {
    try {
        const entries = await getCachedLocalizations();
        return json(entries, {
            headers: {
                'cache-control': 'public, max-age=21600, stale-while-revalidate=86400'
            }
        });
    } catch (err) {
        console.error('Failed to load localizations:', err);
        return json({ error: 'failed to load localizations' }, { status: 500 });
    }
}

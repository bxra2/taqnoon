import { fetchAndParseAllPoFiles } from '$lib/utils/fetchAllPoFiles';

const TTL_MS = 6 * 60 * 60 * 1000;

const globalKey = '__taqnoonLocalizations';
const store = globalThis[globalKey] ?? (globalThis[globalKey] = { value: null, at: 0, inflight: null });

export function getCachedLocalizations() {
    const now = Date.now();
    if (store.value && now - store.at < TTL_MS) return Promise.resolve(store.value);
    if (store.inflight) return store.inflight;
    store.inflight = (async () => {
        try {
            const value = await fetchAndParseAllPoFiles();
            store.value = value;
            store.at = Date.now();
            return value;
        } finally {
            store.inflight = null;
        }
    })();
    return store.inflight;
}

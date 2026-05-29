import { MongoClient } from 'mongodb';
import { env } from '$env/dynamic/private';

const TTL_MS = 5 * 60 * 1000;

const globalKey = '__taqnoonMongo';
const store = /** @type {{ client?: MongoClient, db?: import('mongodb').Db, terms?: CacheEntry, publishers?: CacheEntry }} */ (
    globalThis[globalKey] ?? (globalThis[globalKey] = {})
);

/** @typedef {{ value: any[], at: number, inflight: Promise<any[]> | null }} CacheEntry */

async function connect() {
    if (store.db) return store.db;
    const uri = env.MONGODB_URI || 'mongodb://localhost:27017/taqnoon';
    store.client = new MongoClient(uri);
    await store.client.connect();
    store.db = store.client.db('taqnoon');
    return store.db;
}

export async function getDb() {
    return connect();
}

export async function getTermsCollection() {
    const database = await connect();
    return database.collection('terms');
}

export async function getPublishersCollection() {
    const database = await connect();
    return database.collection('publishers');
}

/**
 * @param {'terms' | 'publishers'} key
 * @param {() => Promise<any[]>} loader
 * @returns {Promise<any[]>}
 */
async function cached(key, loader) {
    const entry = store[key] ?? (store[key] = { value: null, at: 0, inflight: null });
    const now = Date.now();
    if (entry.value && now - entry.at < TTL_MS) return entry.value;
    if (entry.inflight) return entry.inflight;
    entry.inflight = (async () => {
        try {
            const value = await loader();
            entry.value = value;
            entry.at = Date.now();
            return value;
        } finally {
            entry.inflight = null;
        }
    })();
    return entry.inflight;
}

export function getCachedTerms() {
    return cached('terms', async () => {
        const c = await getTermsCollection();
        return c.find({}, { projection: { _id: 0 } }).toArray();
    });
}

export function getCachedPublishers() {
    return cached('publishers', async () => {
        const c = await getPublishersCollection();
        return c.find({}, { projection: { _id: 0 } }).toArray();
    });
}

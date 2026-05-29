import { getCachedTerms, getCachedPublishers } from '$lib/server/db';

export async function load() {
    try {
        const [termData, publishersData] = await Promise.all([
            getCachedTerms(),
            getCachedPublishers()
        ]);
        return { termData, publishersData };
    } catch (error) {
        console.error('Error loading data from MongoDB:', error);
        return { termData: [], publishersData: [] };
    }
}
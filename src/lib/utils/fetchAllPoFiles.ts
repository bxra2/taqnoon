import { parsePo } from '$lib/utils/parsePo';
import { CinnamonUrls, GnomeUrls, KDEUrls } from '../data/localizationUrls';


const poSources = [
  ...GnomeUrls,...CinnamonUrls,...KDEUrls
]

async function runWithConcurrency(tasks, concurrency = 4) {
  const results = [];
  let i = 0;

  async function worker() {
    while (i < tasks.length) {
      const idx = i++;
      try {
        const r = await tasks[idx]();
        results[idx] = { status: 'fulfilled', value: r };
      } catch (err) {
        results[idx] = { status: 'rejected', reason: err };
      }
    }
  }

  const workers = Array.from({ length: Math.max(1, concurrency) }, () => worker());
  await Promise.all(workers);
  return results;
}

export async function fetchAndParseAllPoFiles({ concurrency = 4 } = {}) {
  const tasks = poSources.map((src) => async () => {
    const res = await fetch(src.url);
    if (!res.ok) throw new Error(`Failed to fetch ${src.url}: ${res.status}`);
    const text = await res.text();
    return { src, text };
  });

  const settled = await runWithConcurrency(tasks, concurrency);

  const allEntries = [];
  let idCounter = 1;

  for (const item of settled) {
    if (item?.status === 'fulfilled') {
      const { src, text } = item.value;
      try {
        const parsed = parsePo(text);
        for (const p of parsed) {
          allEntries.push({
            id: src.publisherEn+p.english+idCounter++,
            english: p.english,
            arabic: p.arabic,
            glossaryEn: src.projectEn,
            glossaryAr: src.projectAr,
            glossaryUrl: src.url,
            publisherAr: src.publisherAr,
            publisherEn: src.publisherEn,
            termURL: src.url
          });
        }
      } catch (err) {
        console.error(`Failed parsing ${src.url}`, err);
      }
    } else {
      console.warn('Fetch failed:', item?.reason);
    }
  }

  return allEntries;
}

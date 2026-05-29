import { writable } from 'svelte/store';

export const workerReady = writable(false);
export const workerError = writable('');

let workerInstance: Worker | null = null;

export function getLocalizationWorker(): Worker | null {
    return workerInstance;
}

export function initializeWorker() {
    if (workerInstance) {
        return workerInstance;
    }

    const worker = new Worker(
        new URL ('$lib/workers/localization-worker.ts', import.meta.url),
        { type: 'module' }
    );

    worker.onmessage = (event) => {
        // console.log('Worker message received:', event.data.type);
        
        if (event.data.type === 'ready') {
            // console.log(`Worker loaded ${event.data.count} entries`);
            workerReady.set(true);
            workerError.set('');
        } else if (event.data.type === 'error') {
            console.error('Worker error:', event.data.message);
            workerError.set(event.data.message);
            workerReady.set(false);
        }
        // Note: 'results' messages are handled by individual components
    };

    worker.postMessage({ type: 'load' });
    workerInstance = worker;

    return worker;
}

export function terminateWorker() {
    if (workerInstance) {
        workerInstance.terminate()
        workerInstance = null
        workerReady.set(false)
        workerError.set('')
    }
}
<script lang="ts">
    // Loading.svelte
    export let message = 'جاري التحميل...'
</script>

<div class="flex flex-col items-center justify-center py-12">
    <div class="loading-spinner"></div>
    <p class="mt-4 text-gray-600 dark:text-gray-400">{message}</p>
</div>

<style>
    .loading-spinner {
        width: 48px;
        height: 48px;
        border: 4px solid var(--accent);
        border-top: 4px solid var(--text-color);
        border-radius: 50%;
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        0% {
            transform: rotate(0deg);
        }
        100% {
            transform: rotate(360deg);
        }
    }

    :global(.dark) .loading-spinner {
        border-color: #374151;
        border-top-color: #60a5fa;
    }
</style>
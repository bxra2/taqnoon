<script lang="ts">
    import { browser } from '$app/environment';
    import { initializeWorker } from '$lib/stores/localizationWorker';
    import '$src/app.css';

    let { children } = $props();

    if (browser) {
        initializeWorker();
    }
</script>

{@render children()}

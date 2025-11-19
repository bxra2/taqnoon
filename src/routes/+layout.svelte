<script lang="ts">
    import favicon from '$lib/assets/favicon.svg'
    import DarkMode from '$lib/components/DarkMode.svelte'
    import Glossaries from '$lib/components/Glossaries.svelte'
    import Loading from '$lib/components/Loading.svelte'
    import '../app.css'
    import { goto } from '$app/navigation'
    import { navigating } from '$app/stores'

    let { children } = $props()
    let query = $state('')
    let searchAlign = $state(true)

    function changeDirection(event: Event) {
        const target = event.target as HTMLInputElement | null
        if (target && target.value) {
            const arabicRegex = /[\u0600-\u06FF]/
            searchAlign = arabicRegex.test(target.value)
        } else {
            searchAlign = false
        }
    }

    function search() {
        const params = new URLSearchParams({ q: query })
        goto(`/search?${params.toString()}`)
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === 'Enter') {
            event.preventDefault()
            search()
        }
    }
</script>

<svelte:head>
    <link rel="icon" href={favicon} />
</svelte:head>

<div class="max-w-3xl mx-auto p-6">
    <div dir="rtl" class="flex items-center justify-between mb-6">
        <h1 class="text-3xl font-bold text-center">
            <a href="/">تقنون</a>
        </h1>
        <div class="flex gap-2 items-center">
            <Glossaries />
            <DarkMode />
        </div>
    </div>

    <!-- Search Bar -->
    <div class="flex gap-2 mb-6">
        <input
            type="text"
            bind:value={query}
            onkeydown={handleKeydown}
            placeholder="ابحث عن مصطلحات تقنية..."
            class="flex-1 p-3 border rounded"
            dir={searchAlign ? 'rtl' : 'ltr'}
            oninput={changeDirection}
        />
        <button
            onclick={search}
            class="px-6 py-3 search-button rounded"
            disabled={!query}
        >
            ابحث
        </button>
    </div>

    <!-- Loading State -->
    {#if $navigating}
        <Loading message="جاري تحميل المعاجم..." />
    {:else}
        {@render children()}
    {/if}
</div>

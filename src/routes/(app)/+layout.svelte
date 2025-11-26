<script lang="ts">
    import favicon from '$lib/assets/favicon.svg'
    import DarkMode from '$lib/components/DarkMode.svelte'
    import Glossaries from '$lib/components/Glossaries.svelte'
    import Loading from '$lib/components/Loading.svelte'
    import '$src/app.css'
    import { goto } from '$app/navigation'
    import { navigating, page } from '$app/stores'
    import { isRTL } from '$lib/utils/changeDirection'

    let { children } = $props()
    let query = $state($page.url.searchParams.get('q'))
    let isHome = $derived($page.url.pathname === '/')
    let exactMatch = $state($page.url.searchParams.get('exact') === 'true')
    let includeDescription = $state($page.url.searchParams.get('desc') === 'true')
    let lookupLocalization = $state($page.url.searchParams.get('lookup') === 'true')
    
    let searchAlign = $derived(query ? isRTL(query) : true)

    function search() {
        const params = new URLSearchParams({ q: query })
        if (exactMatch) {
            params.set('exact', 'true')
        }
        if (includeDescription) {
            params.set('desc', 'true')
        }
        if (lookupLocalization) {
            params.set('lookup', 'true')
        }
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

<div
    class="max-w-3xl mx-auto p-6 {isHome
        ? 'min-h-screen flex flex-col justify-center'
        : ''}"
>
    <div dir="rtl" class="flex items-center justify-between mb-6">
        <h1
            class="text-3xl font-bold text-center {isHome
                ? 'text-6xl'
                : 'text-3xl'}"
        >
            <a href="/">تقنون</a>
        </h1>
        <div class="flex gap-2 items-center">
            <Glossaries />
            <DarkMode />
        </div>
    </div>

    <!-- Search Bar -->
    <div class="flex flex-col gap-2 mb-6 border p-4 rounded-lg search-bar">
        <div class="flex gap-2">
            <input
                type="text"
                bind:value={query}
                onkeydown={handleKeydown}
                placeholder="ابحث عن مصطلحات تقنية..."
                class="flex-1 p-3 border rounded"
                dir={searchAlign ? 'rtl' : 'ltr'}
            />
            <button
                onclick={search}
                class="px-6 py-3 search-button rounded"
                disabled={!query}
            >
                ابحث
            </button>
        </div>
         <div class="flex items-center justify-center my-2 gap-6">
            <input type="checkbox" bind:checked={exactMatch} id="exact" />
            <label for="exact" class="flex items-center gap-2 text-sm">
                بحث مطابق تماماً
            </label>
            <input type="checkbox" bind:checked={includeDescription} id="desc" />

            <label for="desc" class="flex items-center gap-2 text-sm">
                يشمل الوصف
            </label>
            <input type="checkbox" bind:checked={lookupLocalization} id="lookup" />
            <label for="lookup" class="flex items-center gap-2 text-sm cursor-pointer">
                بحث في ترجمات البرمجيات مفتوحة المصدر
            </label>
        </div>
    </div>

    <!-- Loading State -->
    {#if $navigating}
        <Loading message="جاري تحميل المعاجم..." />
    {:else}
        {@render children()}
    {/if}
</div>

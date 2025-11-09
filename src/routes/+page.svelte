<script lang="ts">
    import DarkMode from '$lib/components/DarkMode.svelte'
    import Paginator from '$lib/components/Paginator.svelte'
    import TermCard from '$lib/components/TermCard.svelte'
    import { onMount } from 'svelte'

    let query = $state('')
    let searchQuery = $state('') // The actual query used for filtering
    let data = $state([])
    let currentPage = $state(1)
    let limit = $state(10)

    onMount(async () => {
        const modules = import.meta.glob('/src/lib/data/**/*.json')
        const allData = []

        for (const path in modules) {
            const mod = await modules[path]()
            const json = mod.default

            if (json.fileData && Array.isArray(json.entries)) {
                for (const entry of json.entries) {
                    allData.push({
                        ...entry,
                        glossaryEn: json.fileData.glossaryEn,
                        glossaryAr: json.fileData.glossaryAr,
                        url: json.fileData.url,
                        publisherAr: json.fileData.publisherAr,
                        publisherEn: json.fileData.publisherEn,
                    })
                }
            }
        }

        data = allData
    })

    // Filtered results based on search query
    let results = $derived.by(() => {
        if (!searchQuery.trim()) return []

        const lower = searchQuery.toLowerCase()
        return data.filter((item) =>
            Object.values(item).some(
                (value) =>
                    typeof value === 'string' &&
                    value.toLowerCase().includes(lower)
            )
        )
    })

    // Paginated results
    let paginatedResults = $derived.by(() => {
        const startIndex = (currentPage - 1) * limit
        const endIndex = startIndex + limit
        return results.slice(startIndex, endIndex)
    })

    // Reset to page 1 when search query changes
    $effect(() => {
        searchQuery // track search query changes
        currentPage = 1
    })

    // Handle search submit
    function handleSearch() {
        searchQuery = query
    }

    // Handle Enter key
    function handleKeydown(event: KeyboardEvent) {
        if (event.key === 'Enter') {
            handleSearch()
        }
    }
</script>

<div class="max-w-3xl mx-auto p-6">
    <div dir="rtl" class="flex align-center justify-between">
        <h1 class="text-3xl font-bold mb-6 text-center">تقنون</h1>
        <DarkMode />
    </div>

    <div class="flex gap-2 mb-6">
        <input
            type="text"
            bind:value={query}
            onkeydown={handleKeydown}
            placeholder="ابحث عن مصطلحات تقنية..."
            class="flex-1 p-3 border rounded"
            dir="rtl"
        />
        <button
            onclick={handleSearch}
            class="px-6 py-3 search-button text-white rounded"
        >
            بحث
        </button>
    </div>

    {#if results.length > 0}
        <!-- <div class="flex align-center justify-center mb-6">
            <Paginator
                bind:currentPage
                bind:limit
                count={results.length}
                onPageChange={(page, lim) => {
                    window.scrollTo({ top: 0, behavior: 'smooth' })
                }}
                onLimitChange={(page, lim) => {
                    window.scrollTo({ top: 0, behavior: 'smooth' })
                }}
            />
        </div> -->

        <ul class="space-y-4">
            {#each paginatedResults as term (term.English || term.Arabic)}
                <TermCard
                    french={term.French}
                    german={term.German}
                    arabic={term.Arabic}
                    english={term.English}
                    desc={term.Description}
                    glossaryEn={term.glossaryEn}
                    glossaryAr={term.glossaryAr}
                    collectionName={term.publisherEn}
                    collectionNameAr={term.publisherAr}
                    collectionURL={term.publisherUrl}
                    sourceURL={term.url}
                    {term}
                />
            {/each}
        </ul>

        <div class="flex align-center justify-center mt-6">
            <Paginator
                bind:currentPage
                bind:limit
                count={results.length}
                onPageChange={(page, lim) => {
                    window.scrollTo({ top: 0, behavior: 'smooth' })
                }}
            />
        </div>
    {:else if searchQuery}
        <p class="text-center text-gray-500">لا توجد نتائج</p>
    {:else}
        <p class="text-center text-gray-500">أدخل كلمة للبحث</p>
    {/if}
</div>
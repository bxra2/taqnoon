<script lang="ts">
    import DarkMode from '$lib/components/DarkMode.svelte'
    import Paginator from '$lib/components/Paginator.svelte'
    import TermCard from '$lib/components/TermCard.svelte'
    import { onMount } from 'svelte'

    let query = $state('')
    let data = $state([])
    let results = $state([])
    let currentPage = $state(1)
    let limit = $state(10)
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
                        glossaryUrl: json.fileData.glossaryUrl,
                        publisherAr: json.fileData.publisherAr,
                        publisherEn: json.fileData.publisherEn,
                        publisherUrl: json.fileData.publisherUrl,
                    })
                }
            }
        }

        data = allData
    })

    let paginatedResults = $derived.by(() => {
        const startIndex = (currentPage - 1) * limit
        const endIndex = startIndex + limit
        return results.slice(startIndex, endIndex)
    })

    function removeDiacritics(text: string): string {
    return text.normalize('NFKD').replace(/[\u064B-\u065F\u0670\u06D6-\u06ED]/g, '')
    }


function handleSearch() {
	if (!query.trim()) {
		results = []
		return
	}

	const lower = query.toLowerCase()
	const normalizedQuery = removeDiacritics(query)
	results = data
		.filter((item) => {
			const english = item.english?.toLowerCase() || ''
			const arabic = removeDiacritics(item.arabic || '')
			return english.includes(lower) || arabic.includes(normalizedQuery)
		})
		.sort((a, b) => {
			const aEng = a.english?.toLowerCase() || ''
			const bEng = b.english?.toLowerCase() || ''
			const aAr = removeDiacritics(a.arabic || '')
			const bAr = removeDiacritics(b.arabic || '')

			// exact matches first
			const aExact = aEng === lower || aAr === normalizedQuery
			const bExact = bEng === lower || bAr === normalizedQuery
			if (aExact && !bExact) return -1
			if (bExact && !aExact) return 1

			// starts-with next
			const aStarts = aEng.startsWith(lower) || aAr.startsWith(normalizedQuery)
			const bStarts = bEng.startsWith(lower) || bAr.startsWith(normalizedQuery)
			if (aStarts && !bStarts) return -1
			if (bStarts && !aStarts) return 1

			// otherwise keep original order
			return 0
		})

	currentPage = 1
}
    function handleKeydown(event: KeyboardEvent) {
        if (event.key === 'Enter') {
            event.preventDefault()
            handleSearch()
        }
    }
</script>

<div class="max-w-3xl mx-auto p-6">
    <div dir="rtl" class="flex align-center justify-between mb-6">
        <h1 class="text-3xl font-bold text-center">تقنون</h1>
        <DarkMode />
    </div>

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
            onclick={handleSearch}
            class="px-6 py-3 search-button text-white rounded"
        >
            بحث
        </button>
    </div>

    {#if results.length > 0}
        <div class="flex align-center justify-between mb-6">
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
        </div>

        <ul class="space-y-4">
            {#each paginatedResults as term, i (i)}
                <TermCard {term} />
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
    {:else if query}
        <p class="text-center text-gray-500">لا توجد نتائج</p>
    {:else}
        <p class="text-center text-gray-500">أدخل كلمة للبحث</p>
    {/if}
</div>

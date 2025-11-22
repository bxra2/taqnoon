<script lang="ts">
    import { page } from '$app/stores'
    import Loading from '$lib/components/Loading.svelte'
    import Paginator from '$lib/components/Paginator.svelte'
    import TermCard from '$lib/components/TermCard.svelte'

    let { data } = $props()

    let currentPage = $state(1)
    let limit = $state(10)
    let results = $state([])
    let isSearching = $state(false)
    let allTerms = [...data.termData || []]
    function scrollToTop() {
        window.scrollTo({ top: 0, behavior: 'smooth' })
    }

    $effect(() => {
        const q = $page.url.searchParams.get('q')
        if (q) {
            handleSearch(q, allTerms)
        } else {
            results = []
        }
    })

    let paginatedResults = $derived.by(() => {
        const startIndex = (currentPage - 1) * limit
        const endIndex = startIndex + limit
        return results.slice(startIndex, endIndex)
    })

    // Show empty state only when not searching and query exists but no results
    let showEmptyState = $derived(
        !isSearching &&
            results.length === 0 &&
            $page.url.searchParams.get('q')?.trim()
    )

    function removeDiacritics(text: string): string {
        return text
            .normalize('NFKD')
            .replace(/[\u064B-\u065F\u0670\u06D6-\u06ED]/g, '')
    }

    async function handleSearch(query: string, terms) {
        if (!query.trim()) {
            results = []
            console.log('trimimmedd')
            return
        }

        isSearching = true
        await new Promise((resolve) => setTimeout(resolve, 100))

        const lower = query.toLowerCase()
        const normalizedQuery = removeDiacritics(query)

        results = terms
            .filter((item) => {
                const english = item.english?.toLowerCase() || ''
                const arabic = removeDiacritics(item.arabic || '')
                return (
                    english.includes(lower) || arabic.includes(normalizedQuery)
                )
            })
            .sort((a, b) => {
                const aEng = a.english?.toLowerCase() || ''
                const bEng = b.english?.toLowerCase() || ''
                const aAr = removeDiacritics(a.arabic || '')
                const bAr = removeDiacritics(b.arabic || '')

                // Exact matches first
                const aExact = aEng === lower || aAr === normalizedQuery
                const bExact = bEng === lower || bAr === normalizedQuery
                if (aExact && !bExact) return -1
                if (bExact && !aExact) return 1

                // Starts-with next
                const aStarts =
                    aEng.startsWith(lower) || aAr.startsWith(normalizedQuery)
                const bStarts =
                    bEng.startsWith(lower) || bAr.startsWith(normalizedQuery)
                if (aStarts && !bStarts) return -1
                if (bStarts && !aStarts) return 1

                // Otherwise keep original order
                return 0
            })

        // Reset to first page when new search results come in
        currentPage = 1
        isSearching = false
    }
</script>

<div class="mb-4">
    {#if results.length > 0}
        <div class="mb-4">
            <Paginator
                bind:currentPage
                bind:limit
                count={results.length}
                onPageChange={scrollToTop}
                onLimitChange={scrollToTop}
            />
        </div>

        <ul class="space-y-4">
            {#each paginatedResults as term, i (term.id ?? i)}
                <TermCard {term} />
            {/each}
        </ul>

        <div class="mt-4">
            <Paginator
                bind:currentPage
                bind:limit
                count={results.length}
                onPageChange={scrollToTop}
            />
        </div>
    {:else if isSearching}
        <Loading message="جاري البحث..." />
    {:else if showEmptyState}
        <p class="text-center text-gray-500">لا توجد نتائج</p>
    {/if}
</div>

<script lang="ts">
    import { page } from '$app/stores'
    import Loading from '$lib/components/Loading.svelte'
    import Paginator from '$lib/components/Paginator.svelte'
    import TermCard from '$lib/components/TermCard.svelte'
    import Slider from '$src/lib/components/Slider.svelte'

    let { data } = $props()

    let currentPage = $state(1)
    let limit = $state(10)
    let showSlider = $state(false)
    let results = $state([])
    let isSearching = $state(false)
    let allTerms = [...(data.termData || [])]
    function scrollToTop() {
        window.scrollTo({ top: 0, behavior: 'smooth' })
    }

    $effect(() => {
        const q = $page.url.searchParams.get('q')
        const exact = $page.url.searchParams.get('exact') === 'true'
        if (q) {
            handleSearch(q, allTerms, exact)
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

    async function handleSearch(query: string, terms, exact: boolean = false) {
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
                
                if (exact) {
                    return english === lower || arabic === normalizedQuery
                }
                
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
    let currentPublishers = $derived(
        (() => {
            const grouped: Record<string, Set<string>> = {}

            for (const term of results) {
                const publisher = term.publisherAr
                const glossary = term.glossaryAr

                if (!grouped[publisher]) {
                    grouped[publisher] = new Set()
                }

                grouped[publisher].add(glossary)
            }

            // Convert Sets back to arrays
            return Object.fromEntries(
                Object.entries(grouped).map(([k, v]) => [k, Array.from(v)])
            )
        })()
    )

    $effect(() => {
        console.log(results)
        console.log('currentPublishers', currentPublishers)
    })
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
            <button onclick={() => (showSlider = true)}>
                <div class="mt-4 dropdown-icon"></div>
            </button>
        </div>
        <Slider bind:open={showSlider}>
            <h2 class="text-3xl mt-8 mb-6">عدد النتائج : {results.length}</h2>
            <h2 class="text-3xl mt-8 mb-6">المعاجم</h2>

            {#each Object.entries(currentPublishers) as [publisher, glossaries]}
                <div class="publisher-group">
                    <h3 class="mt-4 text-xl">{publisher}</h3>
                    <ul>
                        {#each glossaries as glossary}
                            <li class="m-4"> - {glossary}</li>
                        {/each}
                    </ul>
                </div>
            {/each}
            <!-- {#each currentPublishers as publisher}
                <h3 class="mt-4 text-xl">{publisher}</h3>
            {/each} -->
        </Slider>
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

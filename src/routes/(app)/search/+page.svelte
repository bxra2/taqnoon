<script lang="ts">
    import { page } from '$app/state'
    import Loading from '$lib/components/Loading.svelte'
    import Paginator from '$lib/components/Paginator.svelte'
    import TermCard from '$lib/components/TermCard.svelte'
    import Slider from '$src/lib/components/Slider.svelte'
    import { onMount } from 'svelte'
    import { removeDiacritics } from '$src/lib/utils/removeDiacritics.js'
    import { getLocalizationWorker, workerReady } from '$lib/stores/localizationWorker'

    let { data } = $props()

    let currentPage = $state(1)
    let limit = $state(10)
    let showSlider = $state(false)
    let results = $state([])
    let isSearching = $state(false)
    let allTerms = [...(data.termData || [])]
    let localizationResults = $state([])
    let isWorkerSearching = $state(false)
    let localizationWorker: Worker | null = null
    let messageHandler: ((e: MessageEvent) => void) | undefined

    function scrollToTop() {
        window.scrollTo({ top: 0, behavior: 'smooth' })
    }

    onMount(() => {
        // Get the worker from the shared store
        localizationWorker = getLocalizationWorker()
        
        if (!localizationWorker) {
            console.warn('Localization worker not initialized')
            return
        }

        messageHandler = (e: MessageEvent) => {
            if (e.data.type === 'results') {
                console.log(`Received ${e.data.results.length} localization results`)
                localizationResults = e.data.results
                isWorkerSearching = false
            }
            if (e.data.type === 'error') {
                console.error('Worker error:', e.data.message)
                isWorkerSearching = false
            }
        }
        
        localizationWorker.addEventListener('message', messageHandler)
        
        return () => {
            if (localizationWorker && messageHandler) {
                localizationWorker.removeEventListener('message', messageHandler)
            }
        }
    })

    $effect(() => {
        const q = page.url.searchParams.get('q')
        const exact = page.url.searchParams.get('exact') === 'true'
        const desc = page.url.searchParams.get('desc') === 'true'
        const lookup = page.url.searchParams.get('lookup') === 'true'
        if (q) {
            handleSearch(q, allTerms, exact, desc, lookup)
        } else {
            results = []
            localizationResults = []
        }
    })

    let mergedResults = $derived.by(() => {
        if (localizationResults.length === 0) {
            return results
        }
        
        const uniqueTerms = new Map()
        
        for (const term of results) {
            uniqueTerms.set(term.id, term)
        }
        
        for (const term of localizationResults) {
            if (uniqueTerms.has(term.id)) {
                uniqueTerms.set(term.id, {
                    ...uniqueTerms.get(term.id),
                    ...term,
                    fromLocalization: true
                })
            } else {
                uniqueTerms.set(term.id, { ...term, fromLocalization: true })
            }
        }
        
        const merged = Array.from(uniqueTerms.values())
        return merged
    })

    let paginatedResults = $derived.by(() => {
        const startIndex = (currentPage - 1) * limit
        const endIndex = startIndex + limit
        return mergedResults.slice(startIndex, endIndex)
    })
    // Show empty state only when not searching and query exists but no results
    let showEmptyState = $derived(
        !isSearching &&
            !isWorkerSearching &&
            mergedResults.length === 0 &&
            page.url.searchParams.get('q')?.trim() &&
            (!page.url.searchParams.get('lookup') || $workerReady)
    )

    async function handleSearch(
        query: string,
        terms,
        exact: boolean = false,
        includeDescription: boolean = false,
        lookupLocalization: boolean = false
    ) {
        if (!query.trim()) {
            results = []
            localizationResults = []
            return
        }

        isSearching = true
        
        if (!lookupLocalization) {
            localizationResults = []
        }
        
        await new Promise((resolve) => setTimeout(resolve, 100))

        if (lookupLocalization && localizationWorker && $workerReady) {
            isWorkerSearching = true
            localizationWorker.postMessage({
                type: 'search',
                query,
                exact,
            })
        } else if (lookupLocalization && !$workerReady) {
            console.warn('Worker not ready yet for localization search')
        }

        const lower = query.toLowerCase()
        const normalizedQuery = removeDiacritics(query)

        results = terms
            .filter((item) => {
                const english = item.english?.toLowerCase() || ''
                const arabic = removeDiacritics(item.arabic || '')

                const hasDescription = !!item.description?.trim()

                if (includeDescription && !hasDescription) {
                    return false
                }

                const description = hasDescription
                    ? removeDiacritics(item.description).toLowerCase()
                    : ''
                if (exact) {
                    return english === lower || arabic === normalizedQuery
                }

                return (
                    english.includes(lower) ||
                    arabic.includes(normalizedQuery) ||
                    (hasDescription && description.includes(normalizedQuery))
                )
            })
            .sort((a, b) => {
                const aEng = a.english?.toLowerCase() || ''
                const bEng = b.english?.toLowerCase() || ''
                const aAr = removeDiacritics(a.arabic || '')
                const bAr = removeDiacritics(b.arabic || '')

                const aExact = aEng === lower || aAr === normalizedQuery
                const bExact = bEng === lower || bAr === normalizedQuery
                if (aExact && !bExact) return -1
                if (bExact && !aExact) return 1

                const aStarts =
                    aEng.startsWith(lower) || aAr.startsWith(normalizedQuery)
                const bStarts =
                    bEng.startsWith(lower) || bAr.startsWith(normalizedQuery)
                if (aStarts && !bStarts) return -1
                if (bStarts && !aStarts) return 1

                return 0
            }).map((item, index) => ({...item,
                id: item.id || `${item.publisherEn}${item.english}${index}`
            }))
            // add id to results
        currentPage = 1
        isSearching = false
    }
    let currentPublishers = $derived(
        (() => {
            const grouped: Record<string, Set<string>> = {}

            for (const term of mergedResults) {
                const publisher = term.publisherAr
                const glossary = term.glossaryAr

                if (!grouped[publisher]) {
                    grouped[publisher] = new Set()
                }

                grouped[publisher].add(glossary)
            }

            return Object.fromEntries(
                Object.entries(grouped).map(([k, v]) => [k, Array.from(v)])
            )
        })()
    )
</script>

<div class="mb-4">
    {#if mergedResults.length > 0}
        <div class="mb-4">
            <Paginator
                bind:currentPage
                bind:limit
                count={mergedResults.length}
                onPageChange={scrollToTop}
                onLimitChange={scrollToTop}
            />
            <button onclick={() => (showSlider = true)}>
                <div class="mt-4 dropdown-icon"></div>
            </button>
        </div>
        <Slider bind:open={showSlider}>
            <h2 class="text-3xl mt-8 mb-6">
                عدد النتائج : {mergedResults.length}
            </h2>
            {#if localizationResults.length > 0}
                <h2 class="text-3xl mt-8 mb-6">
                    النتائج المعجمية: {results.length}
                </h2>
                <h2 class="text-3xl mt-8 mb-6">
                    النتائج المترجمة: {localizationResults.length}
                </h2>
            {/if}
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
        </Slider>
        {#if isWorkerSearching}
            <div class="mb-4 p-3 bg-blue-50 dark:bg-blue-900/20 rounded border border-blue-200 dark:border-blue-800 text-center">
                <span class="text-sm text-blue-700 dark:text-blue-300">
                    جاري البحث في الترجمات...
                </span>
            </div>
        {/if}
        <ul class="space-y-4">
            {#each paginatedResults as term, i (term.id ?? i)}
                <TermCard {term} />
            {/each}
        </ul>
        <div class="mt-4">
            <Paginator
                bind:currentPage
                bind:limit
                count={mergedResults.length}
                onPageChange={scrollToTop}
            />
        </div>
    {:else if isSearching || isWorkerSearching}
        <Loading message="جاري البحث..." />
        {#if isWorkerSearching && !isSearching}
            <p class="text-center text-sm text-gray-600 dark:text-gray-400 mt-2">
                البحث في الترجمات قد يستغرق بعض الوقت...
            </p>
        {/if}
    {:else if showEmptyState}
        <p class="text-center text-gray-500">لا توجد نتائج</p>
    {/if}
</div>

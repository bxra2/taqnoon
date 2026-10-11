<script lang="ts">
    import { page } from '$app/state'
    import { goto } from '$app/navigation'
    import Loading from '$lib/components/Loading.svelte'
    import Paginator from '$lib/components/Paginator.svelte'
    import TermCard from '$lib/components/TermCard.svelte'
    import Slider from '$src/lib/components/Slider.svelte'
    import { onMount } from 'svelte'
    import { getLocalizationWorker, workerReady } from '$lib/stores/localizationWorker'

    let currentPage = $state(1)
    let limit = $state(10)
    let showSlider = $state(false)
    let results = $state<any[]>([])
    let total = $state(0)
    let isSearching = $state(false)
    let localizationResults = $state<any[]>([])
    let isWorkerSearching = $state(false)
    let localizationPage = $state(1)
    let localizationLimit = $state(10)
    let localizationWorker: Worker | null = null
    let messageHandler: ((e: MessageEvent) => void) | undefined
    let abortController: AbortController | null = null

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

    async function loadResults(q: string, exact: boolean, desc: boolean) {
        if (abortController) {
            abortController.abort()
        }
        const controller = new AbortController()
        abortController = controller

        isSearching = true

        const offset = (currentPage - 1) * limit
        const params = new URLSearchParams({
            q,
            exact: String(exact),
            desc: String(desc),
            limit: String(limit),
            offset: String(offset),
        })

        try {
            const res = await fetch(`/api/search?${params.toString()}`, {
                signal: controller.signal,
            })
            if (!res.ok) throw new Error(`search ${res.status}`)
            const data = await res.json()
            if (controller.signal.aborted) return
            results = data.items ?? []
            total = data.total ?? 0
            isSearching = false
        } catch (error: any) {
            if (error?.name === 'AbortError') return
            console.error('Search error:', error)
            results = []
            total = 0
            isSearching = false
        }
    }

    $effect(() => {
        const params = page.url.searchParams
        const q = params.get('q')?.trim()
        const exact = params.get('exact') === 'true'
        const desc = params.get('desc') === 'true'
        const lookup = params.get('lookup') === 'true'

        const urlPage = parseInt(params.get('page') || '1', 10)
        const urlLimit = parseInt(params.get('limit') || '10', 10)
        currentPage = isNaN(urlPage) || urlPage < 1 ? 1 : urlPage
        limit = isNaN(urlLimit) || urlLimit < 1 ? 10 : urlLimit

        if (q) {
            loadResults(q, exact, desc)
            if (lookup && $workerReady) {
                isWorkerSearching = true
                getLocalizationWorker()?.postMessage({ type: 'search', query: q, exact })
            } else if (!lookup) {
                localizationResults = []
                isWorkerSearching = false
            }
        } else {
            abortController?.abort()
            results = []
            total = 0
            localizationResults = []
            isSearching = false
            isWorkerSearching = false
        }
    })

    function updatePagination(newPage: number, newLimit: number) {
        const params = new URLSearchParams(page.url.searchParams)
        if (newPage > 1) {
            params.set('page', String(newPage))
        } else {
            params.delete('page')
        }
        if (newLimit !== 10) {
            params.set('limit', String(newLimit))
        } else {
            params.delete('limit')
        }
        goto(`/search?${params.toString()}`)
    }

    function onPageChange(newPage: number, newLimit: number) {
        updatePagination(newPage, newLimit)
    }

    function onLimitChange(newPage: number, newLimit: number) {
        updatePagination(newPage, newLimit)
    }

    let paginatedLocalization = $derived(
        localizationResults.slice(
            (localizationPage - 1) * localizationLimit,
            localizationPage * localizationLimit
        )
    )

    // Show empty state only when not searching and query exists but no results
    let showEmptyState = $derived(
        !isSearching &&
            !isWorkerSearching &&
            results.length === 0 &&
            localizationResults.length === 0 &&
            !!page.url.searchParams.get('q')?.trim() &&
            (!page.url.searchParams.get('lookup') || $workerReady)
    )

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

            return Object.fromEntries(
                Object.entries(grouped).map(([k, v]) => [k, Array.from(v)])
            )
        })()
    )
</script>

<div class="mb-4">
    {#if results.length > 0}
        <div class="mb-4">
            <Paginator
                bind:currentPage
                bind:limit
                count={total}
                onPageChange={onPageChange}
                onLimitChange={onLimitChange}
            />
            <button onclick={() => (showSlider = true)} aria-label="Open slider">
                <div class="mt-4 dropdown-icon"></div>
            </button>
        </div>
        <Slider bind:open={showSlider}>
            <h2 class="text-3xl mt-8 mb-6">
                عدد النتائج : {total}
            </h2>
            {#if localizationResults.length > 0}
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
            {#each results as term, i (term.id ?? i)}
                <TermCard {term} />
            {/each}
        </ul>
        <div class="mt-4">
            <Paginator
                bind:currentPage
                bind:limit
                count={total}
                onPageChange={onPageChange}
                onLimitChange={onLimitChange}
            />
        </div>
    {/if}

    {#if localizationResults.length > 0}
        <div class="mt-8">
            <h2 class="text-2xl font-bold mb-4">
                نتائج ترجمات البرمجيات مفتوحة المصدر
            </h2>
            <Paginator
                bind:currentPage={localizationPage}
                bind:limit={localizationLimit}
                count={localizationResults.length}
            />
            <ul class="space-y-4">
                {#each paginatedLocalization as term, i (term.id ?? i)}
                    <TermCard {term} />
                {/each}
            </ul>
            <div class="mt-4">
                <Paginator
                    bind:currentPage={localizationPage}
                    bind:limit={localizationLimit}
                    count={localizationResults.length}
                />
            </div>
        </div>
    {/if}

    {#if results.length === 0 && localizationResults.length === 0}
        {#if isSearching || isWorkerSearching}
            <Loading message="جاري البحث..." />
            {#if isWorkerSearching && !isSearching}
                <p class="text-center text-sm text-gray-600 dark:text-gray-400 mt-2">
                    البحث في الترجمات قد يستغرق بعض الوقت...
                </p>
            {/if}
        {:else if showEmptyState}
            <p class="text-center text-gray-500">لا توجد نتائج</p>
        {/if}
    {/if}
</div>

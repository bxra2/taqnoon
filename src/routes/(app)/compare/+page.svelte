<script lang="ts">
    import { getContext } from 'svelte'
    import { searchTerms } from '$src/lib/utils/search'
    import TermCard from '$lib/components/TermCard.svelte'
    import { isRTL } from '$lib/utils/changeDirection'
    import { Plus, X, Search } from 'lucide-svelte'

    import type { Term } from '$src/lib/utils/search'
    import Link_2 from 'lucide-svelte/icons/link-2'

    let { data } = $props()
    const compareState = getContext('compareState')

    // Local state for options (could be moved to context if needed globally)
    let partialMatch = $state(false)
    let includeDescription = $state(false)
    let wordResults = $state([])


    type GroupedPublisher = {
        name: string
        url?: string
        glossaries: {
            name: string
            terms: Term[]
        }[]
    }

    function groupResults(terms: Term[]): GroupedPublisher[] {
        const groups: Record<string, {
            url?: string,
            glossaries: Record<string, Term[]>
        }> = {}

        for (const term of terms) {
            // Prioritize Arabic names for display if available
            const pubName = term.publisherAr || term.publisherEn || 'غير معروف'
            const glossaryName = term.glossaryAr || term.glossaryEn || 'عام'

        if (!groups[pubName]) {
            groups[pubName] = {
                url: term.publisherUrl,
                glossaries: {}
            }
        }
        if (!groups[pubName].glossaries[glossaryName]) {
            groups[pubName].glossaries[glossaryName] = []
        }
        groups[pubName].glossaries[glossaryName].push(term)
        }

        // Convert to array and sort
        return Object.entries(groups)
            .sort(([a], [b]) => a.localeCompare(b, 'ar'))
            .map(([pubName, data]) => ({
                name: pubName,
                url: data.url,
                glossaries: Object.entries(data.glossaries)
                    .sort(([a], [b]) => a.localeCompare(b, 'ar'))
                    .map(([glossaryName, terms]) => ({
                        name: glossaryName,
                        terms
                    })),
            }))
    }

    async function handleCompare() {
        if (!data.termData) return
        
        compareState.isSearching = true
        wordResults = JSON.parse(JSON.stringify(compareState.queries))
        // Normalize queries to ensure we have results for each input index
        const activeQueries = compareState.queries
        
        // Perform searches concurrently
        // We map each query to a search result
        const searchPromises = activeQueries.map(async (q) => {
            if (!q.trim()) return []
            // Add a small artificial delay to yield to UI if needed, or just run sync
            // Since searchTerms is sync, we can just run it. 
            // If it's heavy, we might want to wrap in timeout or use worker. 
            // For now, sync is fine as per original implementation preference.
            return searchTerms(q, data.termData, { exact: !partialMatch, includeDescription })
        })
        const results = await Promise.all(searchPromises)
        compareState.results = results
        compareState.isSearching = false
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === 'Enter') {
            handleCompare()
        }
    }

    $effect(() => {
        console.log(compareState.results)
        console.log(wordResults)
    })
</script>

<div class="flex flex-col gap-8">
    <!-- Controls Section -->
    <div class="p-6 rounded-2xl shadow-sm border border-gray-100 dark:border-stone-700">
        <div class="flex flex-col gap-4">
            <div class="flex flex-wrap gap-4 items-end">
                {#each compareState.queries as query, i}
                    <div class="flex-1 min-w-[200px] flex flex-col gap-2">
                        <label for="q-{i}" class="text-sm font-medium ">مصطلح {i + 1}</label>
                        <div class="relative flex items-center justify-between">
                            <div class="flex-1 min-w-[200px] flex flex-col gap-2 relative">
                            <input
                                    id="q-{i}"
                                    type="text"
                                    bind:value={compareState.queries[i]}
                                    onkeydown={handleKeydown}
                                    dir={query.trim() === '' || isRTL(query) ? 'rtl' : 'ltr'}
                                    placeholder="ابحث..."
                                    class="w-full pl-10 focus:ring-2 focus:ring-primary/20 transition-all"
                                />
                            </div>

                            {#if compareState.queries.length > 2}
                            <button
                                onclick={() => compareState.removeQuery(i)}
                                class="absolute left-2 hover:text-red-500 p-1 m-1 cursor-pointer"
                                title="إزالة"
                            >
                                <X size={16} />
                            </button>
                            {/if}
                        </div>
                    </div>
                {/each}

                {#if compareState.queries.length < 5}
                    <button
                        onclick={() => compareState.addQuery()}
                        class="p-2 mb-[1px] border-2 border-dashed border-gray-300 rounded-xl hover:border-primary hover:text-primary transition-colors"
                        title="إضافة مصطلح للمقارنة"
                    >
                        <Plus size={24} />
                    </button>
                {/if}

                <!-- search when all are full -->
                <button
                    onclick={handleCompare}
                    disabled={compareState.queries.some(q => !q.trim())}
                    class="search-button px-8 py-3.5 mb-[1px] bg-primary text-white rounded-xl hover:bg-primary/90 transition-colors flex items-center gap-2 font-bold shadow-lg shadow-primary/20"
                >
                    <Search size={20} />
                    <span>مقارنة</span>
                </button>
            </div>

            <!-- Options -->
            <div class="flex items-center gap-6 mt-2 pt-4 border-t border-gray-100 dark:border-stone-700">
                <label class="flex items-center gap-2 cursor-pointer">
                    <input type="checkbox" bind:checked={partialMatch} class="w-4 h-4 rounded border-gray-300 text-primary focus:ring-primary" />
                    <span class="text-sm">تطابق جزئي</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                    <input type="checkbox" bind:checked={includeDescription} class="w-4 h-4 rounded  focus:ring-primary" />
                    <span class="text-sm">بحث في الوصف</span>
                </label>
            </div>
        </div>
    </div>

    <!-- Results Grid -->
    {#if compareState.results.some(r => r.length > 0)}
        <div class="grid gap-6 items-start" style="grid-template-columns: repeat({compareState.queries.length}, minmax(0, 1fr))">
            {#each compareState.results as columnResults, i}
                <div class="flex flex-col gap-4">
                    <h3 class="font-bold text-center accent p-2 rounded-lg text-lg truncate mb-2" title={wordResults[i]}>
                        {wordResults[i] || `مصطلح ${i + 1}`}
                        <span class="text-sm font-normal opacity-60">({columnResults.length})</span>
                    </h3>
                    
                    {#if columnResults.length > 0}
                        <div class="flex flex-col gap-6">
                            {#each groupResults(columnResults) as publisherGroup}
                                <div class="publisher-section">
                                    <div class="flex items-center gap-4 mb-4">
                                        <div class="h-px flex-1 bg-gray-200 dark:bg-stone-700"></div>
                                        <h2 class="font-bold flex items-center gap-2 text-gray-700 dark:text-gray-300 text-lg whitespace-nowrap px-2">
                                            <a href={publisherGroup.url} target="_blank" rel="noopener noreferrer"><Link_2 size={20} /></a>
                                            {publisherGroup.name}
                                        </h2>
                                        <div class="h-px flex-1 bg-gray-200 dark:bg-stone-700"></div>
                                    </div>
                                
                                    {#each publisherGroup.glossaries as glossary}
                                        <div class="glossary-section mb-6 last:mb-0">
                                            {#if glossary.name !== 'عام'}
                                                <h5 class="text-sm text-gray-400 mb-3 mr-2 font-medium">{glossary.name}</h5>
                                            {/if}
                                            <div class="space-y-4">
                                                {#each glossary.terms as term}
                                                    <TermCard {term} />
                                                {/each}
                                            </div>
                                        </div>
                                    {/each}
                                </div>
                            {/each}
                        </div>
                    {:else if compareState.queries[i]}
                        <div class="text-center p-8 text-gray-400 rounded-xl border border-dashed">
                            لا توجد نتائج
                        </div>
                    {/if}
                </div>
            {/each}
        </div>
    {:else if compareState.isSearching}
        <div class="text-center py-20">
            <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 dark:border-white"></div>
            <p class="mt-4 text-gray-500">جاري البحث...</p>
        </div>
    {:else if compareState.results.length > 0}
         <div class="text-center py-20 text-gray-400">
            لا توجد نتائج مطابقة لأي من المصطلحات
        </div>
    {/if}
</div>

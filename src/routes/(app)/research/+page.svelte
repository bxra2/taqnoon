<script lang="ts">
  import { onMount, untrack } from "svelte"
  import {
    getLocalizationWorker,
    workerReady,
  } from "$lib/stores/localizationWorker"
  import Loading from "$lib/components/Loading.svelte"
  import { removeDiacritics } from "$src/lib/utils/removeDiacritics.js"
  import { isRTL } from "$src/lib/utils/changeDirection.js"
  import { goto } from "$app/navigation"
  import { page } from "$app/state"

  let { data } = $props()

  // Worker State
  let localizationWorker: Worker | null = null
  let messageHandler: ((e: MessageEvent) => void) | undefined
  let workerInterval: any = null

  // UI State
  // Initialize query from URL if present
  let query = $state(page.url.searchParams.get("q") || "")

  let isLoadingStats = $state(false)
  let isWorkerSearching = $state(false)
  let error = $state("")
  let lastSearchedTerm = $state("")

  // Data State
  let finalStats = $state<any>(null)
  let mergedResults = $state<any[]>([])

  // Internal Data
  let localResults = $state<any[]>([])
  let workerResults = $state<any[]>([])
  let workerStats = $state<any>(null)

  // Pagination
  let currentPage = $state(1)
  let limit = $state(10)
  let searchAlign = $derived(query ? isRTL(query) : true)

  // React to URL changes
  $effect(() => {
    const term = page.url.searchParams.get("q") || ""

    untrack(() => {
      if (term !== lastSearchedTerm) {
        lastSearchedTerm = term
        query = term

        if (term) {
          performSearch(term)
        } else {
          // Reset if no term
          finalStats = null
          mergedResults = []
          isLoadingStats = false
          isWorkerSearching = false
        }
      }
    })
  })

  // --- Utils ---
  function normalize(str: string) {
    return removeDiacritics(str || "")
      .toLowerCase()
      .trim()
  }

  function calculateStats(term: string, items: any[]) {
    const normalizedTerm = normalize(term)

    // Filter exact matches (English or Arabic)
    const exactMatches = items.filter((item) => {
      return (
        normalize(item.english) === normalizedTerm ||
        normalize(item.arabic) === normalizedTerm
      )
    })

    const translationCounts: Record<string, number> = {}
    const translationSources: Record<string, string[]> = {}
    const publishers: Record<string, number> = {}

    exactMatches.forEach((match) => {
      const rawAr = removeDiacritics(match.arabic || "").trim()
      if (!rawAr) return

      // Split by delimiters: ; ؛ , ، /
      const cleaned = rawAr.replace(/[0-9]/g, "")

      const synonyms = cleaned
        .split(/::| او |[;؛,،/]+/)
        .map((s) => s.trim())
        .filter(Boolean)

      synonyms.forEach((ar) => {
        translationCounts[ar] = (translationCounts[ar] || 0) + 1

        const source = match.glossaryAr || match.publisherAr || "Dictionary"
        if (!translationSources[ar]) translationSources[ar] = []
        if (!translationSources[ar].includes(source))
          translationSources[ar].push(source)
      })

      const pub = match.publisherAr || "Dictionary"
      // Count publisher once per entry, not per synonym
      publishers[pub] = (publishers[pub] || 0) + 1
    })

    const translations = Object.entries(translationCounts).map(
      ([arabic, count]) => ({
        arabic,
        count,
        sources: translationSources[arabic],
      }),
    )

    const publisherList = Object.entries(publishers).map(
      ([publisher, count]) => ({
        publisher,
        count,
      }),
    )

    return {
      totalOccurrences: exactMatches.length,
      translations,
      publishers: publisherList,
    }
  }

  function mergeStatsData(local: any, worker: any) {
    if (!local && !worker) return null

    const combined = {
      totalOccurrences:
        (local?.totalOccurrences || 0) + (worker?.totalOccurrences || 0),
      translations: [] as any[],
      publishers: [] as any[],
    }

    // Merge Translations
    const transMap = new Map<string, any>()

    const addTrans = (list: any[]) => {
      if (!list) return
      list.forEach((t) => {
        const existing = transMap.get(t.arabic)
        if (existing) {
          existing.count += t.count
          // Merge sources uniquely
          const set = new Set([...existing.sources, ...t.sources])
          existing.sources = Array.from(set)
        } else {
          transMap.set(t.arabic, { ...t })
        }
      })
    }

    addTrans(local?.translations)
    addTrans(worker?.translations)

    combined.translations = Array.from(transMap.values())
      .sort((a, b) => b.count - a.count)
      .map((t) => ({
        ...t,
        percentage:
          combined.totalOccurrences > 0
            ? Math.round((t.count / combined.totalOccurrences) * 100)
            : 0,
      }))

    // Merge Publishers
    const pubMap = new Map<string, number>()
    const addPubs = (list: any[]) => {
      if (!list) return
      list.forEach((p) => {
        pubMap.set(p.publisher, (pubMap.get(p.publisher) || 0) + p.count)
      })
    }

    addPubs(local?.publishers)
    addPubs(worker?.publishers)

    combined.publishers = Array.from(pubMap.entries())
      .map(([publisher, count]) => ({ publisher, count }))
      .sort((a, b) => b.count - a.count)

    return combined
  }

  // --- Search Actions ---
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === "Enter") handleSearchNavigation()
  }

  async function handleSearchNavigation() {
    if (!query.trim()) return
    const currentParam = page.url.searchParams.get("q")
    // If already on this search, just return (or force refresh if needed, but usually not)
    if (currentParam === query.trim()) return

    goto(`/research?q=${encodeURIComponent(query.trim())}`)
  }

  async function performSearch(term?: string) {
    const q = term || query
    if (!q.trim()) return

    // Clear any pending worker wait
    if (workerInterval) {
      clearInterval(workerInterval)
      workerInterval = null
    }

    isLoadingStats = true
    isWorkerSearching = true
    error = ""
    currentPage = 1

    finalStats = null
    workerStats = null
    workerResults = []

    // 1. Local Search & Stats
    const allTerms = data?.termData || []
    const localStats = calculateStats(q, allTerms)

    // Simple fuzzy search for local results list
    const lower = q.toLowerCase()
    localResults = allTerms
      .filter(
        (t) =>
          t.english?.toLowerCase().includes(lower) ||
          (t.arabic &&
            removeDiacritics(t.arabic).includes(removeDiacritics(q))),
      )
      .map((item, idx) => ({ ...item, id: item.id || `local-${idx}` }))

    // Initial render with local data while waiting for worker
    finalStats = mergeStatsData(localStats, null)
    refreshMergedResults()

    // 2. Worker Request
    if (localizationWorker && $workerReady) {
      // NOTE: We only need search results now, we calculate stats on client
      localizationWorker.postMessage({
        type: "search",
        query: q,
        exact: false, // We fetch broad results to calculate stats from EXACT matches in client
      })
    } else {
      waitForWorker(q)
    }
  }

  function waitForWorker(term: string) {
    if (workerInterval) clearInterval(workerInterval)

    workerInterval = setInterval(() => {
      if (localizationWorker && $workerReady) {
        clearInterval(workerInterval)
        workerInterval = null
        localizationWorker.postMessage({
          type: "search",
          query: term,
          exact: false,
        })
      }
    }, 500)
  }

  function refreshMergedResults() {
    // De-duplicate by ID or content
    const unique = new Map()
    localResults.forEach((r) => unique.set(r.id, r))

    workerResults.forEach((r) => {
      // If worker result has same id as local (unlikely given different sources, but possible if id generation overlaps)
      // or just append
      if (!unique.has(r.id)) {
        unique.set(r.id, { ...r, fromLocalization: true })
      }
    })

    mergedResults = Array.from(unique.values())
  }

  // --- Pagination ---
  let paginatedResults = $derived.by(() => {
    const start = (currentPage - 1) * limit
    return mergedResults.slice(start, start + limit)
  })

  function scrollToTop() {
    const el = document.getElementById("results-header")
    if (el) el.scrollIntoView({ behavior: "smooth" })
  }

  onMount(() => {
    localizationWorker = getLocalizationWorker()
    if (!localizationWorker) console.warn("Worker not init")

    messageHandler = (e: MessageEvent) => {
      if (e.data.type === "stats_results") {
        // Legacy handler if worker still sends it
        workerStats = e.data.stats
        const allTerms = data?.termData || []
        const localStats = calculateStats(query, allTerms)
        finalStats = mergeStatsData(localStats, workerStats)
        isLoadingStats = false
      }

      if (e.data.type === "results") {
        // Worker Search Results
        workerResults = e.data.results || []
        isWorkerSearching = false
        refreshMergedResults()

        // Calculate Worker Stats on Client
        const calculatedWorkerStats = calculateStats(query, workerResults)
        const allTerms = data?.termData || []
        const localStats = calculateStats(query, allTerms)

        finalStats = mergeStatsData(localStats, calculatedWorkerStats)
        isLoadingStats = false
      }

      if (e.data.type === "error") {
        // Even if worker errors, we preserve local stats
        isLoadingStats = false
        isWorkerSearching = false
      }
    }

    if (localizationWorker)
      localizationWorker.addEventListener("message", messageHandler)
    return () =>
      localizationWorker?.removeEventListener("message", messageHandler!)
  })
</script>

<div
  class="flex flex-col justify-center align-center max-w-4xl mx-auto p-6 space-y-8"
>
  <!-- Search Bar -->
  <div class="flex gap-2 w-full max-w-2xl border rounded-lg border p-6">
    <input
      type="text"
      bind:value={query}
      onkeydown={handleKeydown}
      placeholder="قارن ترجمات مصطلح واحد..."
      class="flex-1 p-3 border rounded"
      dir={searchAlign ? "rtl" : "ltr"}
    />
    <button
      onclick={handleSearchNavigation}
      class="search-button px-8 py-3.5 mb-[1px] bg-primary text-white rounded-xl hover:bg-primary/90 transition-colors flex items-center gap-2 font-bold shadow-lg shadow-primary/20"
      disabled={!query}
    >
      ابحث
    </button>
  </div>

  <!-- Stats -->
  {#if finalStats && finalStats.translations.length > 0}
    <div class="space-y-6 animate-in fade-in slide-in-from-bottom-2">
      <h2 class="text-2xl font-bold">الاحصائيات المجمعة</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="p-6 rounded-xl shadow-sm border text-center">
          <h3 class="text-sm mb-2">إجمالي التكرارات</h3>
          <p class="text-4xl font-bold">
            {finalStats.totalOccurrences}
          </p>
          {#if finalStats.totalOccurrences === 0}
            <p class="text-sm text-gray-400 mt-2">
              لم يتم العثور على هذا المصطلح
            </p>
          {/if}
        </div>
        <div class="p-6 rounded-xl shadow-sm border text-center">
          <h3 class="text-sm text-gray-500 mb-2">عدد الترجمات المقترحة</h3>
          <p class="text-4xl font-bold text-secondary">
            {finalStats.translations.length}
          </p>
        </div>
        <div class="p-6 rounded-xl shadow-sm border text-center">
          <h3 class="text-sm text-gray-500 mb-2">الترجمة الأكثر شيوعاً</h3>
          <p class="text-2xl font-bold text-green-600 truncate">
            {finalStats.translations[0]?.arabic || "-"}
          </p>
          <p class="text-xs text-gray-400 mt-1">
            نسبة الاستخدام: {finalStats.translations[0]?.percentage || 0}%
          </p>
        </div>
      </div>

      <!-- Charts -->
      <div class="p-6 rounded-xl shadow-sm border">
        <h3 class="font-bold mb-4 text-xl">توزيع الترجمات</h3>
        <div class="space-y-6">
          {#each finalStats.translations as tr}
            <div>
              <div class="flex justify-between items-end mb-1">
                <span class="font-bold text-lg">{tr.arabic}</span>
                <span class="font-mono text-sm opacity-60 text-left" dir="ltr"
                  >{tr.count} occurrences ({tr.percentage}%)</span
                >
              </div>
              <div class="h-4 muted rounded-full overflow-hidden">
                <div
                  class="h-full bg-primary transition-all duration-1000"
                  style="width: {tr.percentage}%"
                ></div>
              </div>
              <div class="text-xs text-gray-400 mt-1 flex flex-wrap gap-1">
                {#each tr.sources as source}
                  <span
                    class="accent px-1.5 py-0.5 rounded border border-gray-100 dark:border-stone-700"
                    >{source}</span
                  >
                {/each}
              </div>
            </div>
          {/each}
        </div>
      </div>

      <!-- Publisher Stats -->
      <div class="p-6 rounded-xl shadow-sm border">
        <h3 class="font-bold mb-4 text-xl">توزيع المصادر</h3>
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
          {#each finalStats.publishers as pub}
            <div
              class="flex justify-between items-center p-3 accent rounded-lg"
            >
              <span class="font-medium">{pub.publisher}</span>
              <span
                class="muted px-2 py-0.5 rounded text-xs font-bold shadow-sm"
                >{pub.count}</span
              >
            </div>
          {/each}
        </div>
      </div>
    </div>
  {:else if isLoadingStats}
    <Loading message="جاري تحليل البيانات..." />
  {:else if query && finalStats.translations.length === 0}
    <p class="text-center text-gray-500 py-12">لا توجد بيانات لهذا المصطلح.</p>
  {/if}
</div>

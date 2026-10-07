<script lang="ts">
  import '$src/app.css'
  import { goto } from '$app/navigation'
  import { isRTL } from '$lib/utils/changeDirection'
  import DarkMode from '$src/lib/components/DarkMode.svelte'
  import Loading from '$src/lib/components/Loading.svelte'
  import { workerReady, workerError } from '$lib/stores/localizationWorker'

  let { data } = $props()
  let query = $state('')
  let exactMatch = $state(false)
  let includeDescription = $state(false)
  let lookupLocalization = $state(false)
  let searchAlign = $derived(query ? isRTL(query) : true)

  let randomTerms = $state<string[]>([])

  $effect(() => {
    if (data.suggestions && data.suggestions.length > 0 && randomTerms.length === 0) {
      const pool = [
        ...new Set(
          data.suggestions
            .map((s: any) => {
              const ar = s.arabic?.trim()
              if (ar && !ar.includes(' ')) return ar
              return s.english?.trim()
            })
            .filter(Boolean),
        ),
      ]

      // Fisher-Yates shuffle
      for (let i = pool.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1))
        ;[pool[i], pool[j]] = [pool[j], pool[i]]
      }

      randomTerms = pool.slice(0, 10)
    }
  })

  function search() {
    if (!query) return
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

<div class="absolute top-4 left-4">
  <DarkMode />
</div>

<div class="min-h-screen flex flex-col justify-start sm:justify-center items-center p-4 sm:p-6 pt-20 sm:pt-0">
  <div class="flex flex-col items-center gap-3 sm:gap-4 mb-6 sm:mb-8">
    <img class="home-logo" src="/icons/taqnun.svg" alt="تقنون" />
    <h2 class="text-base sm:text-lg md:text-2xl text-center font-bold px-2 sm:px-4">
      محرك بحث جامع للمعاجم التقنية لسهولة البحث و المقارنة فيها
    </h2>
  </div>
  <div class="flex flex-col gap-2 w-full max-w-3xl sm:max-w-2xl mb-4 sm:mb-6 border p-3 sm:p-4 rounded-lg border" dir="rtl">
    <div class="flex gap-2 w-full max-w-2xl mb-6" dir="rtl">
      <input
        type="text"
        bind:value={query}
        onkeydown={handleKeydown}
        placeholder="ابحث عن مصطلحات تقنية..."
        class="flex-1 p-3 border rounded"
        dir={searchAlign ? 'rtl' : 'ltr'}
      />
    </div>
    <div class="flex flex-col sm:flex-row w-full max-w-2xl mb-6 px-1 gap-3 sm:gap-6" dir="rtl">
      <label class="flex items-center gap-2 text-sm cursor-pointer">
        <input type="checkbox" bind:checked={exactMatch} />
        بحث مطابق
      </label>
      <label class="flex items-center gap-2 text-sm cursor-pointer">
        <input type="checkbox" bind:checked={includeDescription} />
        يشمل الوصف
      </label>
      <label
        class="flex items-center gap-2 text-sm cursor-pointer"
        class:opacity-50={!$workerReady}
      >
        <input
          type="checkbox"
          bind:checked={lookupLocalization}
          disabled={!$workerReady}
        />
        <span class="break-words whitespace-normal">
          بحث في ترجمات البرمجيات مفتوحة المصدر
        </span>
      </label>
    </div>
    <div class="flex gap-2 w-full max-w-2xl justify-center" dir="rtl">
      <button
        onclick={search}
        class="search-button px-5 py-2.5 sm:px-8 sm:py-3.5 mb-[1px] bg-primary text-white rounded-xl hover:bg-primary/90 transition-colors flex items-center gap-2 font-bold shadow-lg shadow-primary/20"
        disabled={!query}
      >
        ابحث
      </button>
    </div>

    {#if randomTerms.length > 0}
      <div class="w-full pt-4 mt-2 border-t border-gray-100 dark:border-stone-700">
        <p class="text-xs text-center text-gray-400 mb-3">مقترحات للبحث</p>
        <div class="flex flex-wrap justify-center gap-2">
          {#each randomTerms as term}
            <button
              class="px-3 py-1.5 text-sm accent cursor-pointer hover:bg-primary/10 hover:text-primary rounded-full transition-colors border border-gray-200 dark:border-stone-700"
              onclick={() => {
                query = term
                search()
              }}
            >
              {term}
            </button>
          {/each}
        </div>
      </div>
    {/if}
  </div>
  {#if !$workerReady && !$workerError}
    <span class="text-xs opacity-60">
      <Loading message="جاري التحميل..." />
    </span>
  {:else if $workerError}
    <span class="text-xs text-red-500">({$workerError})</span>
  {/if}
</div>
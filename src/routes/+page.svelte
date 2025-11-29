<script lang="ts">
    import '$src/app.css'
    import { goto } from '$app/navigation'
    import { isRTL } from '$lib/utils/changeDirection'
    import Glossaries from '$src/lib/components/Glossaries.svelte'
    import DarkMode from '$src/lib/components/DarkMode.svelte'
    import Loading from '$src/lib/components/Loading.svelte'
    import { workerReady, workerError } from '$lib/stores/localizationWorker'

    let query = $state('')
    let exactMatch = $state(false)
    let includeDescription = $state(false)
    let lookupLocalization = $state(false)
    let searchAlign = $derived(query ? isRTL(query) : true)

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

<div class="min-h-screen flex flex-col justify-center items-center p-6">
    <div class="flex flex-col items-center gap-4 mb-8">
        <img class="home-logo" src="/icons/taqnun.svg" alt="تقنون" />
        <h2 class="text-lg sm:text-xl md:text-2xl text-center font-bold px-4">
            محرك بحث جامع للمعاجم التقنية لسهولة البحث و المقارنة فيها
        </h2>
    </div>
    <div class="flex flex-col gap-2 w-full max-w-2xl mb-6 border p-4 rounded-lg search-bar" dir="rtl">
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
        <div class="flex w-full max-w-2xl mb-6 px-1 gap-6" dir="rtl">
            <label class="flex items-center gap-2 text-sm cursor-pointer">
                <input type="checkbox" bind:checked={exactMatch} />
                بحث مطابق
            </label>
            <label class="flex items-center gap-2 text-sm cursor-pointer">
                <input type="checkbox" bind:checked={includeDescription} />
                يشمل الوصف
            </label>
            <label class="flex items-center gap-2 text-sm cursor-pointer">
                <input
                    type="checkbox"
                    bind:checked={lookupLocalization}
                    disabled={!$workerReady}
                />
                <span class:opacity-50={!$workerReady}>
                    بحث في ترجمات البرمجيات مفتوحة المصدر
                </span>
            </label>
        </div>
        <div class="flex gap-2 w-full max-w-2xl justify-center" dir="rtl">
            <Glossaries />
            <button
                onclick={search}
                class="px-6 py-3 search-button rounded"
                disabled={!query}
            >
                ابحث
            </button>
        </div>
    </div>
    {#if !$workerReady && !$workerError}
        <span class="text-xs opacity-60">
            <Loading message="جاري التحميل..." />
        </span>
    {:else if $workerError}
        <span class="text-xs text-red-500">({$workerError})</span>
    {/if}
</div>
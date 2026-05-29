<script lang="ts">
    import { setContext } from 'svelte'
    import { writable } from 'svelte/store'
    import { page } from '$app/state'
    import DarkMode from '$lib/components/DarkMode.svelte'

    let { children } = $props()

    // Shared state for the compare route
    // Using a simple store or state object if Svelte 5 runes allow deep reactivity easily
    // But since we want to share it, context is good.
    
    // We can use a rune-based state object and pass it via context
    class CompareState {
        queries = $state(['', ''])
        results = $state([])
        isSearching = $state(false)
        
        addQuery() {
            if (this.queries.length < 5) {
                this.queries.push('')
            }
        }
        
        removeQuery(index: number) {
            if (this.queries.length > 2) {
                this.queries.splice(index, 1)
            }
        }

        updateQuery(index: number, value: string) {
            this.queries[index] = value
        }
    }

    const compareState = new CompareState()
    
    setContext('compareState', compareState)

</script>

<div class="min-h-screen text-gray-900 dark:text-gray-100 flex flex-col">
    <!-- Custom Header for Compare Route -->
    <header class="border-b border-gray-200 dark:border-stone-800 p-4">
        <div class="max-w-7xl mx-auto flex items-center justify-between">
            <h1 class="text-2xl font-bold almahdi">
                <a href="/">تقنون</a>
                <span class="text-gray-400 mx-2">|</span>
                <span class="text-lg">مقارنة المصطلحات</span>
            </h1>
            
            <div class="flex gap-4 items-center">
                 <DarkMode />
            </div>
        </div>
    </header>

    <main class="flex-1 max-w-7xl mx-auto w-full p-6">
        {@render children()}
    </main>
</div>

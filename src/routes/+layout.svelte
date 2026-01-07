<script lang="ts">
  import { browser } from '$app/environment';
  import { initializeWorker } from '$lib/stores/localizationWorker';
  import { Search, LayoutDashboard, GitCompare, BookCopy } from 'lucide-svelte';
  import '$src/app.css';

  let { children } = $props();
  let open = $state(true);

  if (browser) {
    initializeWorker();
  }
</script>

<div class="flex h-screen">
  <!-- Sidebar -->
  <aside
    class="overflow-hidden border rounded-3xl
           transition-all duration-200 ease-in-out border-r"
    class:w-64={open}
    class:w-16={!open}
  >
    <!-- Toggle -->
    <button
      class="p-4 text-2xl opacity-70 hover:opacity-100"
      onclick={() => (open = !open)}
    >
      ☰
    </button>

    <!-- Nav -->
    <nav class="mt-4 space-y-2">
      <a href="/" class="nav-item truncate flex items-center gap-3 px-4 py-2 text-sm opacity-80 hover:opacity-100
           hover:bg-stone-800 rounded">
        <Search size={20} />
        {#if open}<span>بحث</span>{/if}
      </a>

      <a href="/dashboard" class="nav-item truncate flex items-center gap-3 px-4 py-2 text-sm opacity-80 hover:opacity-100
           hover:bg-stone-800 rounded">
        <LayoutDashboard size={20} />
        {#if open}<span>لوحة معلومات</span>{/if}
      </a>

      <a href="/compare" class="nav-item truncate flex items-center gap-3 px-4 py-2 text-sm opacity-80 hover:opacity-100
           hover:bg-stone-800 rounded">
        <GitCompare size={20} />
        {#if open}<span>مقارنة كلمات</span>{/if}
      </a>

       <a href="/glossaries" class="nav-item truncate flex items-center gap-3 px-4 py-2 text-sm opacity-80 hover:opacity-100
           hover:bg-stone-800 rounded">
        <BookCopy size={20} />
        {#if open}<span>المعاجم</span>{/if}
      </a>
    </nav>
  </aside>

  <!-- Page -->
  <main class="flex-1 overflow-auto p-6">
    {@render children()}
  </main>
</div>

<style>
  /* .nav-item {
    @apply flex items-center gap-3 px-4 py-2
           text-sm opacity-80 hover:opacity-100
           hover:bg-slate-800 rounded;
  } */
</style>

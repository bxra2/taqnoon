<script lang="ts">
  import { browser } from "$app/environment";
  import { initializeWorker } from "$lib/stores/localizationWorker";
  import {
    Search,
    LayoutDashboard,
    GitCompare,
    BookCopy,
    ChartBar,
  } from "lucide-svelte";
  import "$src/app.css";

  let { children } = $props();
  let open = $state(true);

  if (browser) {
    initializeWorker();
  }
</script>

<div class="flex h-screen">
  <!-- Sidebar -->
  <aside
    class="overflow-hidden rounded-3xl
           transition-all duration-200 ease-in-out border-r slider"
    class:w-64={open}
    class:w-16={!open}
  >
    <!-- Toggle -->
    <button
      class="p-4 text-4xl opacity-70 hover:opacity-100 margin-auto"
      onclick={() => (open = !open)}
    >
      ☰
    </button>

    <!-- Nav -->
    <nav
      class="mt-4 space-y-2 flex flex-col"
      class:items-start={open}
      class:items-center={!open}
    >
      <a
        href="/dashboard"
        class="nav-item truncate flex items-center gap-3 px-4 py-2 text-sm hover-accent rounded"
      >
        <LayoutDashboard size={20} />
        {#if open}<span>لوحة معلومات</span>{/if}
      </a>
      <a
        href="/"
        class="nav-item truncate flex items-center gap-3 px-4 py-2 text-sm hover-accent rounded"
      >
        <Search size={20} />
        {#if open}<span>بحث</span>{/if}
      </a>
      <a
        href="/glossaries"
        class="nav-item truncate flex items-center gap-3 px-4 py-2 text-sm hover-accent rounded"
      >
        <BookCopy size={20} />
        {#if open}<span>المعاجم</span>{/if}
      </a>

      <a
        href="/compare"
        class="nav-item truncate flex items-center gap-3 px-4 py-2 text-sm hover-accent rounded"
      >
        <GitCompare size={20} />
        {#if open}<span>مقارنة كلمات</span>{/if}
      </a>

      <a
        href="/research"
        class="nav-item truncate flex items-center gap-3 px-4 py-2 text-sm hover-accent rounded"
      >
        <ChartBar size={20} />
        {#if open}<span>بحث متقدم</span>{/if}
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
  .slider {
    border-radius: 20px 0 0 20px;
    border-left: 1px solid var(--border-color);
    box-shadow: -2px 0 8px rgba(0, 0, 0, 0.2);
  }
</style>

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

  let isMobile = $state(false);
  let open = $state(true);
  let wasDesktop = true;

  $effect(() => {
    if (!browser) return;
    initializeWorker();

    // Set initial mobile state (SSR-safe: starts as desktop, corrects post-mount)
    const nowMobile = window.innerWidth < 768;
    isMobile = nowMobile;
    open = !nowMobile;
    wasDesktop = !nowMobile;

    const onResize = () => {
      const nowMobile = window.innerWidth < 768;
      if (nowMobile && wasDesktop) open = false;
      if (!nowMobile && !wasDesktop) open = true;
      wasDesktop = !nowMobile;
      isMobile = nowMobile;
    };
    window.addEventListener("resize", onResize);
    return () => window.removeEventListener("resize", onResize);
  });
</script>

<svelte:window
  onkeydown={(e) => e.key === "Escape" && isMobile && (open = false)}
/>

{#if isMobile && !open}
  <button
    class="fixed top-4 right-4 z-60 p-4 text-4xl opacity-70 hover:opacity-100"
    onclick={() => (open = true)}
    aria-label="Open sidebar"
  >
    ☰
  </button>
{/if}

{#if isMobile && open}
  <button
    class="fixed inset-0 z-50 bg-black/15"
    onclick={() => (open = false)}
    aria-label="Close sidebar"
  ></button>
{/if}

<div class="flex h-[100dvh]">
  <!-- Sidebar -->
  <aside
    class="overflow-hidden slider
           fixed top-0 right-0 h-[100dvh] z-60 w-80 bg-[var(--bg-color)]
           md:relative md:h-auto md:z-auto
           transition-transform duration-300"
    class:md:w-64={open}
    class:md:w-16={!open}
    class:translate-x-full={isMobile && !open}
    class:translate-x-0={isMobile && open}
  >
    <!-- Toggle (desktop only) -->
    <button
      class="p-4 text-4xl opacity-70 hover:opacity-100 mx-auto hidden md:block"
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
  <main class="flex-1 overflow-auto p-6" inert={isMobile && open}>
    {@render children()}
  </main>
</div>

<style>
  .slider {
    border-radius: 20px 0 0 20px;
    border-left: 1px solid var(--border-color);
    box-shadow: -2px 0 8px rgba(0, 0, 0, 0.2);
  }
</style>

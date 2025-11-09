<script lang="ts">
    import { toArabicIndic } from '$lib/utils/numerals'

    interface Props {
        currentPage?: number
        limit?: number
        count: number
        onPageChange?: (page: number, limit: number) => void
        onLimitChange?: (page: number, limit: number) => void
    }

    let {
        currentPage = $bindable(1),
        limit = $bindable(10),
        count,
        onPageChange,
        onLimitChange
    }: Props = $props()

    const pagesToShow = 2 // Show 2 pages before and after the current page
    const limitsList = [5, 10, 20, 50]

    // Calculate total pages
    let totalPages = $derived(Math.ceil(count / limit))
    
    // Calculate the page range to display
    let pageRange = $derived(calculatePageRange(currentPage, totalPages, pagesToShow))

    function calculatePageRange(current: number, total: number, range: number) {
        let start = Math.max(1, current - range)
        let end = Math.min(total, current + range)

        // Expand range if we have room
        if (end - start < range * 2) {
            if (start === 1) {
                end = Math.min(total, start + range * 2)
            } else if (end === total) {
                start = Math.max(1, end - range * 2)
            }
        }

        const pages = []
        for (let i = start; i <= end; i++) {
            pages.push(i)
        }

        return pages
    }

    function goToPage(page: number) {
        if (page < 1 || page > totalPages || page === currentPage) return
        currentPage = page
        onPageChange?.(page, limit)
    }

    function changeLimit(newLimit: number) {
        if (newLimit === limit) return
        limit = newLimit
        currentPage = 1 // Reset to page 1 when changing limit
        onLimitChange?.(1, newLimit)
    }
</script>

<nav class="pagination" aria-label="نتائج البحث التصفح">
    <div class="pages" role="navigation">
        <!--  الصفحة السابقة  -->
        <button
            class="page-btn"
            class:disabled={currentPage === 1}
            disabled={currentPage === 1}
            onclick={() => goToPage(currentPage - 1)}
            aria-label="الصفحة السابقة"
        >
            السابق
        </button>

        <!--  الصفحة الاولي  -->
        {#if currentPage > pagesToShow + 2}
            <button
                class="page-btn"
                onclick={() => goToPage(1)}
                aria-label="الصفحة 1"
            >
                {toArabicIndic(1)}
            </button>
            <span class="ellipsis" aria-hidden="true">...</span>
        {/if}

        <!--  الصفحات  -->
        {#each pageRange as page}
            <button
                class="page-btn"
                class:active={page === currentPage}
                disabled={page === currentPage}
                onclick={() => goToPage(page)}
                aria-label={`الصفحة ${page}`}
                aria-current={page === currentPage ? 'page' : undefined}
            >
                {toArabicIndic(page)}
            </button>
        {/each}

        <!--  الصفحة الاخيرة  -->
        {#if currentPage + pagesToShow + 1 < totalPages}
            <span class="ellipsis" aria-hidden="true">...</span>
            <button
                class="page-btn"
                onclick={() => goToPage(totalPages)}
                aria-label={`الصفحة ${totalPages}`}
            >
                {toArabicIndic(totalPages)}
            </button>
        {/if}

        <!--  الصفحة التالية  -->
        <button
            class="page-btn"
            class:disabled={currentPage === totalPages}
            disabled={currentPage === totalPages}
            onclick={() => goToPage(currentPage + 1)}
            aria-label="الصفحة التالية"
        >
            التالي
        </button>
    </div>

    <div class="limit">
        {#if limitsList[0] <= count}
            <span>عدد النتائج</span>

            {#each limitsList as limitItem}
                {#if limitItem <= count}
                    <button
                        class="limit-btn"
                        class:active={limit === limitItem}
                        disabled={limit === limitItem}
                        onclick={() => changeLimit(limitItem)}
                        aria-label={`عرض ${limitItem} نتيجة`}
                    >
                        {toArabicIndic(limitItem)}
                    </button>
                {/if}
            {/each}
        {/if}
    </div>
</nav>

<style>
    .pagination {
        display: flex;
        justify-content: space-around;
        align-items: center;
        margin-top: 20px;
        flex-wrap: wrap;
        gap: 20px;
    }

    .pages {
        display: flex;
        justify-content: center;
        gap: 10px;
        align-items: center;
        flex-wrap: wrap;
    }

    .limit {
        gap: 10px;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-wrap: wrap;
    }

    .page-btn,
    .limit-btn {
        padding: 6px 8px;
        border: 1px solid #ccc;
        border-radius: 4px;
        cursor: pointer;
        color: black;
        background-color: white;
        box-shadow: 3px 3px 0 #4a4a4a;
        font-family: inherit;
        font-size: inherit;
        transition: all 0.2s ease;
    }

    .page-btn:disabled,
    .limit-btn:disabled {
        cursor: not-allowed;
        opacity: 0.5;
    }

    .page-btn.active,
    .limit-btn.active {
        background-color: #4a4a4a;
        color: white;
        cursor: not-allowed;
    }

    .page-btn:not(:disabled):hover,
    .limit-btn:not(:disabled):hover {
        background-color: #eeeeee;
        transform: translateY(-1px);
    }

    .page-btn.active:hover,
    .limit-btn.active:hover {
        background-color: #6e6e6e;
        transform: none;
    }

    .page-btn:focus-visible,
    .limit-btn:focus-visible {
        outline: 2px solid #4a4a4a;
        outline-offset: 2px;
    }

    .ellipsis {
        padding: 0 4px;
        color: #666;
        user-select: none;
    }

    @media (max-width: 640px) {
        .pagination {
            flex-direction: column;
            gap: 15px;
        }

        .page-btn,
        .limit-btn {
            padding: 8px 12px;
        }
    }
</style>
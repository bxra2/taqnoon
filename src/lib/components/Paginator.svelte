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
        onLimitChange,
    }: Props = $props()

    const pagesToShow = 1
    const limitsList = [5, 10, 20, 50]

    let totalPages = $derived(Math.ceil(count / limit))

    let pageRange = $derived(
        calculatePageRange(currentPage, totalPages, pagesToShow)
    )

    function calculatePageRange(
        currentPage: number,
        totalPages: number,
        pagesToShow: number
    ) {
        let start = Math.max(1, currentPage - pagesToShow)
        let end = Math.min(totalPages, currentPage + pagesToShow)

        // Try to expand range to full width if possible
        const desiredCount = pagesToShow * 2 + 1
        let actualCount = end - start + 1

        if (actualCount < desiredCount) {
            const remaining = desiredCount - actualCount

            // Shift start left if possible
            start = Math.max(1, start - remaining)
            actualCount = end - start + 1

            // Shift end right if still not enough
            if (actualCount < desiredCount) {
                end = Math.min(totalPages, end + (desiredCount - actualCount))
            }
        }

        // Build pages array without duplicates
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
        {#if limitsList[0] < count}
            <label for="limit-select">عدد النتائج</label>

            <select
                id="limit-select"
                class="limit-select"
                bind:value={limit}
                onchange={(e) => changeLimit(Number(e.target.value))}
                aria-label="تغيير عدد النتائج"
            >
                {#each limitsList as limitItem}
                    {#if limitItem <= count}
                        <option value={limitItem}>
                            {toArabicIndic(limitItem)}
                        </option>
                    {/if}
                {/each}
            </select>
        {/if}
    </div>
</nav>

<style>
    .pagination {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-top: 0px;
        flex-wrap: wrap;
        gap: 10px;
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

    .page-btn {
        padding: 6px 8px;
        border: 1px solid #ccc;
        border-radius: 4px;
        cursor: pointer;
        /* color: black;
        background-color: white; */
        box-shadow: 3px 3px 0 #4a4a4a;
        font-family: inherit;
        font-size: 1rem;
        transition: all 0.2s ease;
    }

    .page-btn:disabled {
        cursor: not-allowed;
        opacity: 0.5;
    }

    .page-btn.active {
        background-color: #4a4a4a;
        color: white;
        cursor: not-allowed;
    }

    .page-btn:not(:disabled):hover {
        opacity: 0.9;
        /* background-color: #eeeeee; */
        transform: translateY(-1px);
    }

    .page-btn.active:hover {
        background-color: #6e6e6e;
        transform: none;
    }

    .page-btn:focus-visible {
        outline: 2px solid #4a4a4a;
        outline-offset: 2px;
    }

    .ellipsis {
        padding: 0 4px;
        color: #666;
        user-select: none;
    }

    .limit-select {
        padding: 6px 8px;
        border: 1px solid #ccc;
        border-radius: 4px;
        /* background-color: white; */
        box-shadow: 3px 3px 0 #4a4a4a;
        font-family: inherit;
        font-size: 1rem;
        cursor: pointer;
        transition: all 0.2s ease;
    }

    .limit-select:focus-visible {
        outline: 2px solid #4a4a4a;
        outline-offset: 2px;
    }

    @media (max-width: 640px) {
        .pagination {
            flex-direction: column;
            gap: 15px;
        }

        .page-btn {
            padding: 4px 6px;
        }

        .limit,
        .pages {
            gap: 8px;
        }
    }
</style>

<script>
    let {
        showPublishers = $bindable(false),
        results = $bindable([]),
        query = $bindable(''),
    } = $props()

    function togglePublishers() {
        showPublishers = !showPublishers
        if (showPublishers) {
            results = []
            query = ''
        }
    }
</script>

<button
    onclick={togglePublishers}
    class="px-3 py-2 search-button rounded"
    title="المعاجم والناشرين"
>
    المعاجم
</button>

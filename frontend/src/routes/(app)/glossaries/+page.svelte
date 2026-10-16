<script>
    let { data } = $props()

    let { publishersData, localizationSourcesData } = data

    const gridCols =
        'grid grid-cols-[2.5rem_minmax(0,1fr)_5rem_4rem] gap-3 items-center'

    const COLLAPSE_THRESHOLD = 10

    let expanded = $state({})

    function toggle(id) {
        expanded[id] = !expanded[id]
    }

    function visibleGlossaries(publisher) {
        const list = publisher.glossaries || []
        if (list.length <= COLLAPSE_THRESHOLD || expanded[publisher.id]) {
            return list
        }
        return list.slice(0, COLLAPSE_THRESHOLD)
    }
</script>

<div class="space-y-6" dir="rtl">
    <h2 class="text-2xl text-center font-bold mb-4">الناشرون والمعاجم</h2>
    {#each publishersData as publisher}
        <div class="my-6 border rounded-3xl overflow-hidden">
            <div class="flex justify-between items-center p-4">
                <p class="text-2xl">
                    <a href={publisher.publisherUrl} class="no-underline">
                        {publisher.publisherAr}
                    </a>
                </p>
                <p class="text-2xl" dir="ltr">
                    <a href={publisher.publisherUrl} class="no-underline">
                        {publisher.publisherEn}
                    </a>
                </p>
            </div>
            <div
                class="{gridCols} px-4 py-2 border-b text-sm opacity-60"
            >
                <span aria-hidden="true" class="w-10"></span>
                <span>المعجم</span>
                <span class="text-center">ويكيداتا</span>
                <span class="text-center">الكتاب</span>
            </div>
            <ul class="divide-y">
                {#each visibleGlossaries(publisher) as glossary}
                    <li
                        class="{gridCols} px-4 py-3"
                    >
                        <div
                            class="h-14 w-10 shrink-0 rounded overflow-hidden border"
                        >
                            {#if glossary.imageUrl}
                                <img
                                    src={glossary.imageUrl}
                                    alt={glossary.glossaryAr}
                                    class="h-full w-full object-cover"
                                    loading="lazy"
                                />
                            {:else}
                                <div
                                    class="h-full w-full flex items-center justify-center opacity-40"
                                >
                                    📖
                                </div>
                            {/if}
                        </div>
                        <div class="min-w-0">
                            <a
                                href={glossary.glossaryUrl}
                                class="no-underline font-semibold block truncate"
                            >
                                {glossary.glossaryAr}
                            </a>
                            <a
                                href={glossary.glossaryUrl}
                                class="no-underline text-sm opacity-60 block truncate"
                                dir="ltr"
                            >
                                {glossary.glossaryEn}
                            </a>
                        </div>
                        <div class="text-center text-sm">
                            {#if glossary.wikidataUrl}
                                <a
                                    href={glossary.wikidataUrl}
                                    target="_blank"
                                    rel="noreferrer"
                                    class="no-underline opacity-70 hover:opacity-100"
                                >
                                    ويكيداتا
                                </a>
                            {:else}
                                <span aria-hidden="true"></span>
                            {/if}
                        </div>
                        <div class="text-center text-sm">
                            <a
                                href={glossary.bookUrl ||
                                    glossary.glossaryUrl}
                                target="_blank"
                                rel="noreferrer"
                                class="no-underline opacity-70 hover:opacity-100"
                            >
                                الكتاب
                            </a>
                        </div>
                    </li>
                {/each}
            </ul>
            {#if publisher.glossaries.length > COLLAPSE_THRESHOLD}
                <button
                    class="w-full py-2 text-sm opacity-70 hover:opacity-100 border-t"
                    onclick={() => toggle(publisher.id)}
                >
                    {expanded[publisher.id]
                        ? 'عرض أقل'
                        : `عرض الكل (${publisher.glossaries.length})`}
                </button>
            {/if}
        </div>
    {/each}
    <h2 class="text-2xl text-center font-bold mb-4">مشاريع التوطين</h2>
    {#each localizationSourcesData as publisher}
        <div class="my-6 border rounded-3xl overflow-hidden">
            <div class="flex justify-between items-center p-4">
                <p class="text-2xl">{publisher.publisherAr}</p>
                <p class="text-2xl" dir="ltr">{publisher.publisherEn}</p>
            </div>
            <div class="flex justify-between items-center px-4 py-2 border-b">
                <p class="text-xl">- المشاريع :</p>
                <p class="text-xl" dir="ltr">- projects :</p>
            </div>
            <ul class="divide-y">
                {#each publisher.projects as project}
                    <li
                        class="flex justify-between items-center px-4 py-3"
                    >
                        <span>
                            <a href={project.url} class="no-underline">
                                {project.projectAr}
                            </a>
                        </span>
                        <span class="text-left" dir="ltr">
                            <a
                                href={project.url}
                                class="no-underline text-left"
                            >
                                {project.projectEn}
                            </a>
                        </span>
                    </li>
                {/each}
            </ul>
        </div>
    {/each}
</div>

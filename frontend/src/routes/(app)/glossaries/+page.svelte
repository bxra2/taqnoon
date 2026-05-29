<script>
    let { data } = $props()

    let { publishersData } = data
</script>

<div class="space-y-6" dir="rtl">
    <h2 class="text-2xl text-center font-bold mb-4">الناشرون والمعاجم</h2>
    {#each publishersData as publisher}
        <div class="my-6 border p-4 rounded-3xl">
            <div class="flex justify-between items-center mb-6 mt-4">
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
            <div class="flex justify-between items-center mb-6 mt-4">
                <p class="text-xl">- المعاجم :</p>
                <p class="text-xl" dir="ltr">- glossaries :</p>
            </div>
            <ul>
                {#each publisher.glossaries as glossary}
                    <li class="flex justify-between items-center mx-2">
                        <span>
                            <a href={glossary.glossaryUrl} class="no-underline">
                                {glossary.glossaryAr}
                            </a>
                        </span>
                        <span class="text-left" dir="ltr">
                            <a
                                href={glossary.glossaryUrl}
                                class="no-underline text-left"
                            >
                                {glossary.glossaryEn}
                            </a>
                        </span>
                    </li>
                {/each}
            </ul>
        </div>
    {/each}
</div>

<!-- partial -->
<script lang="ts">
    import { marked } from 'marked'
    export interface Term {
        english: string
        arabic: string
        french: string
        german: string
        description?: string
        desc?: string // ← add this
        tURL: string
        publisherAr: string
        publisherEn: string
        publisherUrl: string
        glossaryAr: string
        glossaryEn: string
        glossaryUrl: string
    }
    const { term }: { term: Term } = $props()

    const desc = $derived(term.description ?? term.desc)
</script>

<article class="term-card">
    <header>
        <h2 class="en">
            {#if term.tURL}
                <a target="_blank" href={term.tURL}>{term.english}</a>
            {:else}
                {term.english}
            {/if}
        </h2>
        <h3 class="ar">
            {#if term.tURL}
                <a target="_blank" href={term.tURL}>{term.arabic}</a>
            {:else}
                {term.arabic}
            {/if}
        </h3>
    </header>
    {#if desc}
        <b>الوصف: </b>
        <p class="desc">{@html marked.parse(desc)}</p>
    {/if}

    {#if term.german || term.french}
        <section class="translations">
            {#if term.french}<p><b>فرنسي:</b> {term.french}</p>{/if}
            {#if term.german}<p><b>ألماني:</b> {term.german}</p>{/if}
        </section>
    {/if}

    <footer>
        <small class="arb" dir="rtl">
            {#if term.glossaryAr}
                <span>
                    <b>المعجم:&ensp;</b>
                    <a href={term.glossaryUrl} target="_blank">
                        {term.glossaryAr}
                    </a>
                </span>
                <br />
            {/if}
            {#if term.publisherAr}
                <span>
                    <b>الناشر:&ensp;</b>
                    <a href={term.publisherUrl} target="_blank">
                        {term.publisherAr}
                    </a>
                </span>
            {/if}
        </small>
        <small dir="ltr">
            {#if term.glossaryEn}
                <span>
                    <b>Glossary:&ensp;</b>
                    <a href={term.glossaryUrl} target="_blank">
                        {term.glossaryEn}
                    </a>
                </span>
                <br />
            {/if}
            {#if term.publisherEn}
                <span>
                    <b>Publisher:&ensp;</b>
                    <a href={term.publisherUrl} target="_blank">
                        {term.publisherEn}
                    </a>
                </span>
            {/if}
        </small>
    </footer>
</article>

<style>
</style>

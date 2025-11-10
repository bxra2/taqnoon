<!-- partial -->
<script lang="ts">
    export interface Term {
        english: string
        arabic: string
        french: string
        german: string
        description: string
        termURL: string
        publisherAr: string
        publisherEn: string
        publisherUrl: string
        glossaryAr: string
        glossaryEn: string
        glossaryUrl: string
    }
    const { term }: { term: Term } = $props()

    console.log(term)
</script>

<article class="term-card">
    <header>
        <h2 class="en"><a href={term.termURL}>{term.english}</a></h2>
        <h3 class="ar"><a href={term.termURL}>{term.arabic}</a></h3>
    </header>

    {#if term.german || term.french}
        <section class="translations">
            {#if term.french}<p><b>فرنسي:</b> {term.french}</p>{/if}
            {#if term.german}<p><b>ألماني:</b> {term.german}</p>{/if}
        </section>
    {/if}

    {#if term.description}
        <p class="desc">{term.description}</p>
    {/if}

    <footer>
        {#if term.glossaryAr}
            <small class="arb" dir="rtl">
                <span>
                    <b>المعجم: </b>
                    <a href={term.glossaryUrl} target="_blank">
                        {term.glossaryAr}
                    </a>
                </span>
                <br />
                {#if term.publisherAr}
                    <span>
                        <b>الناشر: </b>
                        <a href={term.publisherUrl} target="_blank">
                            {term.publisherAr}
                        </a>
                    </span>
                {/if}
            </small>
        {/if}
        {#if term.glossaryEn}
            <small dir="ltr">
                <span>
                    <b>Glossary: </b>
                    <a href={term.glossaryUrl} target="_blank">
                        {term.glossaryEn}
                    </a>
                </span>
                <br />
                {#if term.publisherEn}
                    <span>
                        <b>Publisher: </b>
                        <a href={term.publisherUrl} target="_blank">
                            {term.publisherEn}
                        </a>
                    </span>
                {/if}
            </small>
        {/if}
    </footer>
</article>

<style>
</style>

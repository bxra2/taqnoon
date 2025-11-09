<!-- partial -->
<script>
    export let french,
        arabic,
        german,
        english,
        desc,
        glossaryEn,
        glossaryAr,
        glossaryUrl,
        publisherNameAr,
        publisherName,
        publisherURL,
        term
</script>

<article class="term-card">
    <header>
        <h2 class="en">{english}</h2>
        <h3 class="ar">{arabic}</h3>
    </header>

    {#if german || french}
        <section class="translations">
            {#if french}<p><b>French:</b> {french}</p>{/if}
            {#if german}<p><b>German:</b> {german}</p>{/if}
        </section>
    {/if}

    {#if desc}
        <p class="desc">{desc}</p>
    {/if}

    <footer>
        {#if glossaryEn}
            <small dir="ltr">
                <span>
                    <b>Glossary: </b>
                    {glossaryEn}
                </span>
                <br />
                {#if term.publisherEn}
                    <span>
                        <b>Publisher: </b>
                        {term.publisherEn}
                    </span>
                {/if}
                {#if term.glossaryUrl}
                    <span>
                        <a href={term.glossaryUrl} target="_blank">Source ↗</a>
                    </span>
                {/if}
            </small>
        {/if}
        {#if glossaryAr}
            <small dir="rtl">
                <span>
                    <b>المعجم: </b>
                    {glossaryAr}
                </span>
                <br />
                {#if term.publisherAr}
                    <span>
                        <b>الناشر: </b>
                        {term.publisherAr}
                    </span>
                {/if}
                {#if term.glossaryUrl}
                    <span>
                        <a href={term.glossaryUrl} target="_blank">المصدر ↗</a>
                    </span>
                {/if}
            </small>
        {/if}
    </footer>
</article>

<style>
    .term-card {
        background-color: var(--paper);
        border: 1px solid var(--border);
        padding: 1rem 1.2rem;
        border-radius: 6px;
        margin-bottom: 1.2rem;
        max-width: 1200px;
    }

    header {
        border-bottom: 1px solid var(--border);
        padding-bottom: 0.3rem;
        margin-bottom: 0.6rem;
    }

    .en {
        font-size: 1.3rem;
        font-weight: 600;
        color: var(--ink);
    }

    .ar {
        font-size: 1.2rem;
        direction: rtl;
        color: var(--muted);
        font-family: 'Noto Serif Arabic', serif;
    }

    .translations p {
        margin: 0.3rem 0;
        font-size: 0.9rem;
        color: var(--muted);
    }

    .desc {
        font-size: 0.95rem;
        color: var(--ink);
        margin-top: 0.6rem;
        direction: rtl;
    }

    footer {
        border-top: 1px dotted var(--border);
        margin-top: 0.8rem;
        padding-top: 0.4rem;
        font-size: 0.8rem;
        color: var(--muted);
        display: flex;
        align-items: center;
        justify-content: space-between;
    }
</style>

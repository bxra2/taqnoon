<script lang="ts">
    let { open = $bindable(), children } = $props()
    function close() {
        open = !open
    }
</script>

{#if open}
    <button 
        class="overlay" 
        onclick={close} 
        aria-label="Close slider"
        type="button"
    ></button>
{/if}

<div class="slider {open ? 'open' : ''}">
    <button class="close-btn" onclick={close} aria-label="Close">×</button>
    <div class="slider-content">
        {@render children?.()}
    </div>
</div>

<style>
    .overlay {
        position: fixed;
        inset: 0;
        background: rgba(0, 0, 0, 0.1);
        z-index: 50;
        border: none;
        width: 100%;
        height: 100%;
        cursor: pointer;
    }

    .slider {
        position: fixed;
        top: 0;
        right: -350px;
        width: 350px;
        height: 100vh;
        background: var(--bg-color);
        padding: 20px;
        border-radius: 20px 0 0 20px;
        border-left: 1px solid var(--border-color);
        box-shadow: -2px 0 8px rgba(0, 0, 0, 0.2);
        transition: right 0.3s ease;
        z-index: 60;
        display: flex;
        flex-direction: column;
    }

    .slider.open {
        right: 0;
    }
    .slider-content {
        flex: 1;
        overflow-y: auto;
        margin-top: 40px;
        scrollbar-width: thin;
    }

    .close-btn {
        position: absolute;
        left: 20px;
        top: 10px;
        background: none;
        border: none;
        font-size: 1.8rem;
        cursor: pointer;
        color: var(--text-color);
        z-index: 61;
    }
</style>

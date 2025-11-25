<script lang="ts">
    let { open = $bindable(), children } = $props()
    function close() {
        open = !open
    }
</script>

<!-- Overlay -->
{#if open}
    <div class="overlay" onclick={close}></div>
{/if}

<!-- Slider Window -->
<div class="slider {open ? 'open' : ''}">
    <button class="close-btn" onclick={close}>×</button>
    {@render children?.()}
</div>

<style>
    .overlay {
        position: fixed;
        inset: 0;
        background: rgba(0, 0, 0, 0.1);
        z-index: 50;
    }

    .slider {
        position: fixed;
        top: 0;
        right: -350px;
        width: 350px;
        height: 100vh;
        background: var(--bg-color);
        padding: 20px;
        border-radius: 20px;
        border: 1px solid var(--border-color);
        box-shadow: -2px 0 8px rgba(0, 0, 0, 0.2);
        transition: right 0.3s ease;
        z-index: 60;
    }

    .slider.open {
        right: 0;
    }

    .close-btn {
        position: absolute;
        left: 20px;
        top: 10px;
        background: none;
        border: none;
        font-size: 1.8rem;
        cursor: pointer;
    }
</style>

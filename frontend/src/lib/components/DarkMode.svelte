<script lang="ts">
    import { onMount } from 'svelte'

    let isDark = $state(false)
    let mounted = $state(false)

    onMount(() => {
        const savedTheme = localStorage.getItem('theme')

        if (savedTheme) {
            isDark = savedTheme === 'dark'
            document.documentElement.setAttribute('data-theme', savedTheme)
        } else {
            const prefersDark = window.matchMedia(
                '(prefers-color-scheme: dark)'
            ).matches
            isDark = prefersDark
            const theme = isDark ? 'dark' : 'light'
            document.documentElement.setAttribute('data-theme', theme)
            localStorage.setItem('theme', theme)
        }

        mounted = true
    })

    // Update theme when toggle changes (only after mount)
    $effect(() => {
        if (mounted) {
            const theme = isDark ? 'dark' : 'light'
            document.documentElement.setAttribute('data-theme', theme)
            localStorage.setItem('theme', theme)
        }
    })
</script>

<div class="switch">
    <label class="theme-switch" for="checkbox">
        <input
            type="checkbox"
            id="checkbox"
            bind:checked={isDark}
            aria-label="Toggle dark mode"
        />
        <div class="slider round">
            <div class="icon-container">
                <img
                    src="/icons/sun.svg"
                    alt="Light mode"
                    class="sun-icon"
                    class:visible={!isDark}
                />
                <img
                    src="/icons/moon.svg"
                    alt="Dark mode"
                    class="moon-icon"
                    class:visible={isDark}
                />
            </div>
        </div>
    </label>
</div>

<style>
    .switch {
        height: 40px;
    }

    .theme-switch {
        display: inline-block;
        height: 34px;
        position: relative;
        width: 60px;
    }

    .theme-switch input {
        display: none;
    }

    .slider {
        background-color: #ccc;
        bottom: 0;
        cursor: pointer;
        left: 0;
        position: absolute;
        right: 0;
        top: 0;
        transition: background-color 0.4s;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .slider:before {
        background-color: #fff;
        bottom: 4px;
        content: '';
        height: 26px;
        left: 4px;
        position: absolute;
        transition: transform 0.4s;
        width: 26px;
        z-index: 2;
    }

    .icon-container {
        position: relative;
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .slider img {
        width: 18px;
        height: 18px;
        object-fit: contain;
        position: absolute;
        opacity: 0;
        transition: opacity 0.3s ease;
        pointer-events: none;
    }

    .slider img.visible {
        opacity: 1;
        transition: none;
    }

    .sun-icon {
        left: 8px;
        fill: #24201a;
        filter: #24201a;
        z-index: 3;
    }

    .moon-icon {
        right: 8px;
        fill: #24201a;
        filter: #24201a;
        z-index: 3;
    }

    input:checked + .slider {
        background-color: #2b2b2b;
    }

    input:checked + .slider:before {
        transform: translateX(26px);
    }

    .slider.round {
        border-radius: 34px;
    }

    .slider.round:before {
        border-radius: 50%;
    }

    .theme-switch:hover .slider {
        opacity: 0.9;
    }

    .theme-switch input:focus + .slider {
        box-shadow: 0 0 0 2px rgba(66, 153, 225, 0.5);
    }
</style>

<script lang="ts">
    import { onMount } from 'svelte'

    let isDark = $state(false)
    let mounted = $state(false)

    // Initialize theme on mount
    onMount(() => {
        // Get saved theme or default to light
        const savedTheme = localStorage.getItem('theme')
        
        if (savedTheme) {
            isDark = savedTheme === 'dark'
            document.documentElement.setAttribute('data-theme', savedTheme)
        } else {
            // Optional: Check system preference
            const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
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
        />
        <div class="slider round"></div>
    </label>
</div>

<style>
    .switch {
        height: 70px;
        /* position: fixed; */
        right: 50px;
        top: 50px;
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
        transition: 0.4s;
    }

    .slider:before {
        background-color: #fff;
        bottom: 4px;
        content: '';
        height: 26px;
        left: 4px;
        position: absolute;
        transition: 0.4s;
        width: 26px;
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
</style> 

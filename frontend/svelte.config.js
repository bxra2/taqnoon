import adapter from '@sveltejs/adapter-static'

export default {
    kit: {
        alias: {
            $src: './src',
            $routes: './src/routes',
        },
        adapter: adapter({
            pages: 'dist',
            assets: 'dist',
            fallback: 'index.html', // needed for SPA routing
        }),
    },
}

import adapter from '@sveltejs/adapter-static';

export default {
    kit: {
        adapter: adapter({
            pages: 'build',
            assets: 'build',
            fallback: 'index.html',
        }),

        alias: {
            $src: './src',
            $routes: './src/routes',
        },

        prerender: {
            handleHttpError: 'ignore' // prevents build crashes on dynamic routes
        }
    }
};
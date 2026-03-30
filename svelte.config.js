import adapter from '@sveltejs/adapter-cloudflare';

export default {
    kit: {
        adapter: adapter(),
        alias: { $src: './src', $routes: './src/routes' },
        prerender: { handleHttpError: 'ignore' }
    }
};

export function convertMarkdownLinksToButtons(input: string): string {
    return input.replace(/\[([^\]]+)]\(([^)]+)\)/g, (_match, text, url) => {
        return `${text} <button onclick="window.open('${url}', '_blank')">🔗</button>`;
    });
}

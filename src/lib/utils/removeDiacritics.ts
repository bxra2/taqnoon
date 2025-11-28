export function removeDiacritics(text: string): string {
    return text
        .normalize('NFKD')
        .replace(/[\u064B-\u065F\u0670\u06D6-\u06ED]/g, '')
}
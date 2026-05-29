export function isRTL(text: string): boolean {
    const arabicRegex = /[\u0600-\u06FF]/
    return arabicRegex.test(text)
}
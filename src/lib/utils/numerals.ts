export function toArabicIndic(str: Number) {
    const arabicIndicDigits = ['٠', '١', '٢', '٣', '٤', '٥', '٦', '٧', '٨', '٩']
    return str
        .toString()
        .replace(/\d/g, (digit: string) => arabicIndicDigits[+digit])
}

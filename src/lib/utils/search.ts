import { removeDiacritics } from '$src/lib/utils/removeDiacritics.js'

export interface SearchOptions {
    exact?: boolean
    includeDescription?: boolean
}

export interface Term {
    id?: string
    english?: string
    arabic?: string
    description?: string
    publisherEn?: string
    link?: string
    // Allow other properties
    [key: string]: any
}


export function searchTerms(query: string, terms: Term[], options: SearchOptions = {}): Term[] {
    const { exact = false, includeDescription = false } = options
    
    if (!query || !query.trim()) {
        return []
    }

    const lower = query.toLowerCase()
    const normalizedQuery = removeDiacritics(query)

    return terms
        .filter((item) => {
            const english = item.english?.toLowerCase() || ''
            const arabic = removeDiacritics(item.arabic || '')

            const hasDescription = !!item.description?.trim()

            if (includeDescription && !hasDescription) {
                return false
            }

            const description = hasDescription
                ? removeDiacritics(item.description).toLowerCase()
                : ''
            
            if (exact) {
                return english === lower || arabic === normalizedQuery
            }

            return (
                english.includes(lower) ||
                arabic.includes(normalizedQuery) ||
                (hasDescription && description.includes(normalizedQuery))
            )
        })
        .sort((a, b) => {
            const aEng = a.english?.toLowerCase() || ''
            const bEng = b.english?.toLowerCase() || ''
            const aAr = removeDiacritics(a.arabic || '')
            const bAr = removeDiacritics(b.arabic || '')

            const aExact = aEng === lower || aAr === normalizedQuery
            const bExact = bEng === lower || bAr === normalizedQuery
            if (aExact && !bExact) return -1
            if (bExact && !aExact) return 1

            const aStarts =
                aEng.startsWith(lower) || aAr.startsWith(normalizedQuery)
            const bStarts =
                bEng.startsWith(lower) || bAr.startsWith(normalizedQuery)
            if (aStarts && !bStarts) return -1
            if (bStarts && !aStarts) return 1

            return 0
        })
        .map((item, index) => ({
            ...item,
            id: item.id || `${item.publisherEn}${item.english}${index}`
        }))
}

export function parseTxTFiles(text: string) {
    const entries: { id: string; english: string; arabic: string }[] = [];
    const lines = text.split('\n');

    let original = '';
    let translation = '';
    let idCounter = 1;

    const flush = () => {
        if (original) {
            entries.push({
                id: original + idCounter++,
                english: original,
                arabic: translation,
            });
        }
        original = '';
        translation = '';
    };

    for (const line of lines) {
        if (line.startsWith('[')) {
            flush();
            continue;
        }
        if (line.startsWith('original=')) {
            original = line.slice('original='.length);
        } else if (line.startsWith('translation=')) {
            translation = line.slice('translation='.length);
        }
    }

    flush();

    return entries;
}

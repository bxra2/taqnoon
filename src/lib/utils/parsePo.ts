export function parsePo(poText: string) {
    const entries = [];
    let msgid = '';
    let msgstr = '';

    const lines = poText.split('\n');

    let idCounter = 1;
    const publisherAr = 'جنوم'; 
    const publisherEn = 'GNOME';

    for (const line of lines) {
    // Skip comments
    if (line.startsWith('#')) continue;

    if (line.startsWith('msgid')) {
        msgid = line.match(/msgid\s+"(.*)"/)?.[1] ?? '';
    }

    if (line.startsWith('msgstr')) {
        msgstr = line.match(/msgstr\s+"(.*)"/)?.[1] ?? '';

        if (msgid) {
            entries.push({
                id: publisherEn+idCounter++,   
                english: msgid,
                arabic: msgstr,
                publisherAr,   
                publisherEn 
            });

            msgid = '';
            msgstr = '';
        }
    }
}

    return entries;
}

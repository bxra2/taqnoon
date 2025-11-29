export function parsePo(poText: string) {
    const entries = [];
    let msgid = '';
    let msgstr = '';

    const lines = poText.split('\n');

    let idCounter = 1;

    for (const line of lines) {
    // Skip comments
    if (line.startsWith('#')) continue;
    // get msgid (english)
    if (line.startsWith('msgid')) {
        msgid = line.match(/msgid\s+"(.*)"/)?.[1] ?? '';
    }
    // get msgstr (arabic)
    if (line.startsWith('msgstr')) {
        msgstr = line.match(/msgstr\s+"(.*)"/)?.[1] ?? '';

        if (msgid) {
            entries.push({
                id: msgid+idCounter++,
                english: msgid,
                arabic: msgstr,
            });

            msgid = '';
            msgstr = '';
        }
    }
}

    return entries;
}

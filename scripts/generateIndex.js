// scripts/generateIndex.js
import fs from 'fs';
import path from 'path';

const baseDir = path.resolve('../static/data/publishers');
const indexFile = path.join(baseDir, 'index.json');

function walk(dir) {
    let results = [];
    const list = fs.readdirSync(dir, { withFileTypes: true });
    for (const file of list) {
        const fullPath = path.join(dir, file.name);
        if (file.isDirectory()) {
            results = results.concat(walk(fullPath));
        } else if (file.name.endsWith('.json')) {
            // store relative path from baseDir
            results.push(path.relative(baseDir, fullPath).replace(/\\/g, '/'));
        }
    }
    return results;
}

const allFiles = walk(baseDir);
fs.writeFileSync(indexFile, JSON.stringify(allFiles, null, 2), 'utf-8');
console.log(`Generated index.json with ${allFiles.length} files.`);
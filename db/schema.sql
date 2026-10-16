CREATE TABLE IF NOT EXISTS publisher (
    id INTEGER PRIMARY KEY,
    publisherEn TEXT,
    publisherAr TEXT,
    publisherUrl TEXT
);

CREATE TABLE IF NOT EXISTS domain (
    id INTEGER PRIMARY KEY,
    domainAr TEXT NOT NULL,
    domainEn TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS glossary (
    id INTEGER PRIMARY KEY,
    glossaryEn TEXT,
    glossaryAr TEXT,
    glossaryUrl TEXT,
    publisher_id INTEGER REFERENCES publisher(id),
    domain_id INTEGER REFERENCES domain(id),
    license TEXT,
    image_url TEXT,
    book_url TEXT,
    wikidata_url TEXT,
    imported_at TEXT
);

CREATE TABLE IF NOT EXISTS term (
    id INTEGER PRIMARY KEY,
    glossary_id INTEGER REFERENCES glossary(id),
    english TEXT,
    arabic TEXT,
    french TEXT,
    german TEXT,
    term_url TEXT,
    description TEXT,
    description_norm TEXT,
    english_norm TEXT,
    arabic_norm TEXT,
    french_norm TEXT,
    german_norm TEXT
);

CREATE INDEX IF NOT EXISTS idx_term_english ON term(english);
CREATE INDEX IF NOT EXISTS idx_term_arabic ON term(arabic);
CREATE INDEX IF NOT EXISTS idx_term_french ON term(french);
CREATE INDEX IF NOT EXISTS idx_term_german ON term(german);
CREATE INDEX IF NOT EXISTS idx_term_english_norm ON term(english_norm);
CREATE INDEX IF NOT EXISTS idx_term_arabic_norm ON term(arabic_norm);
CREATE INDEX IF NOT EXISTS idx_term_glossary_id ON term(glossary_id);
CREATE INDEX IF NOT EXISTS idx_glossary_publisher_id ON glossary(publisher_id);
CREATE INDEX IF NOT EXISTS idx_glossary_domain_id ON glossary(domain_id);

CREATE TABLE IF NOT EXISTS source_file (
    path TEXT PRIMARY KEY,
    file_hash TEXT NOT NULL,
    glossary_id INTEGER NOT NULL REFERENCES glossary(id) ON DELETE CASCADE,
    ingested_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_source_file_glossary_id ON source_file(glossary_id);


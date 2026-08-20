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
    english_norm TEXT,
    arabic_norm TEXT,
    french_norm TEXT,
    german_norm TEXT
);

CREATE INDEX IF NOT EXISTS idx_term_english ON term(english);
CREATE INDEX IF NOT EXISTS idx_term_arabic ON term(arabic);
CREATE INDEX IF NOT EXISTS idx_term_french ON term(french);
CREATE INDEX IF NOT EXISTS idx_term_german ON term(german);
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

INSERT INTO publisher (publisherEn, publisherAr, publisherUrl)
SELECT 'United Nations', 'منظمة الأمم المتحدة', 'https://www.un.org'
WHERE NOT EXISTS (SELECT 1 FROM publisher WHERE publisherAr = 'منظمة الأمم المتحدة');
UPDATE publisher SET publisherEn = COALESCE(publisherEn, 'United Nations'), publisherUrl = COALESCE(publisherUrl, 'https://www.un.org') WHERE publisherAr = 'منظمة الأمم المتحدة';

INSERT INTO domain (id, domainAr, domainEn) VALUES
    (1, 'تقنية المعلومات', 'Information Technology'),
    (2, 'القانون الدولي', 'International Law'),
    (3, 'الطب', 'Medicine')
ON CONFLICT(id) DO UPDATE SET
    domainAr = excluded.domainAr,
    domainEn = excluded.domainEn;

INSERT INTO glossary (id, glossaryEn, glossaryAr, glossaryUrl, publisher_id, domain_id, license, imported_at) VALUES
    (1, 'Computing Terms Glossary', 'معجم مصطلحات الحوسبة', NULL, (SELECT id FROM publisher WHERE publisherAr = 'منظمة الأمم المتحدة'), 1, 'CC BY 4.0', datetime('now')),
    (2, 'Legal Terms Glossary', 'معجم المصطلحات القانونية', NULL, (SELECT id FROM publisher WHERE publisherAr = 'منظمة الأمم المتحدة'), 2, 'CC BY 4.0', datetime('now')),
    (3, 'Medical Glossary', 'المعجم الطبي', NULL, (SELECT id FROM publisher WHERE publisherAr = 'منظمة الأمم المتحدة'), 3, 'CC BY 4.0', datetime('now'))
ON CONFLICT(id) DO UPDATE SET
    glossaryEn = excluded.glossaryEn,
    glossaryAr = excluded.glossaryAr,
    glossaryUrl = excluded.glossaryUrl,
    publisher_id = excluded.publisher_id,
    domain_id = excluded.domain_id,
    license = excluded.license;

INSERT INTO term (id, glossary_id, english, arabic, french, german, term_url, description) VALUES
    (1, 1, 'Algorithm', 'خوارزمية', 'Algorithme', 'Algorithmus', 'https://www.un.org/example/1', 'خطوات محدّدة لحل مسألة أو إنجاز مهمة'),
    (2, 1, 'Database', 'قاعدة بيانات', 'Base de données', 'Datenbank', NULL, 'مجموعة منظمة من البيانات يمكن الوصول إليها وإدارتها'),
    (3, 1, 'Network', 'شبكة', 'Réseau', 'Netzwerk', NULL, 'ربط مجموعة من الأجهزة لتبادل البيانات'),
    (4, 1, 'Compiler', 'مُصرِّف', 'Compilateur', 'Compiler', NULL, 'برنامج يحوّل الشيفرة المصدرية إلى لغة الآلة'),
    (5, 2, 'Treaty', 'معاهدة', 'Traité', 'Vertrag', NULL, 'اتفاقية دولية ملزمة بين الدول'),
    (6, 2, 'Jurisdiction', 'الاختصاص القضائي', 'Compétence', 'Gerichtsbarkeit', NULL, 'سلطة المحكمة في نظر الدعوى والفصل فيها'),
    (7, 2, 'Sovereignty', 'سيادة', 'Souveraineté', 'Souveränität', NULL, NULL),
    (8, 3, 'Diagnosis', 'تشخيص', 'Diagnostic', 'Diagnose', NULL, NULL),
    (9, 1, 'Operating System', 'نظام التشغيل', 'Système d''exploitation', 'Betriebssystem', NULL, NULL),
    (10, 3, 'Anatomy', 'علم التشريح', 'Anatomie', 'Anatomie', NULL, NULL)
ON CONFLICT(id) DO UPDATE SET
    glossary_id = excluded.glossary_id,
    english = excluded.english,
    arabic = excluded.arabic,
    french = excluded.french,
    german = excluded.german,
    term_url = excluded.term_url,
    description = excluded.description;


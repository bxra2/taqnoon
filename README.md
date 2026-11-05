# تقنون

> ذاكرة بحثية للمصطلح التقني العربي: نجمع المعاجم، نحفظ سياقها، ونقارن ترجماتها.

**TaQnun** — a searchable, comparative collection of Arabic technical and scientific terminology.

[العربية](#العربية) · [English](#english)

---

## العربية

### ما هو تقنون؟

تقنون مشروع توثيقي وبحثي يجمع المصطلحات التقنية والعلمية العربية من معاجم ومسارد مستقلة، ويعيد عرضها في مكان واحد للبحث والمقارنة. الفكرة ليست إصدار «ترجمة معيارية» واحدة، بل حفظ تعدد الآراء: لكل مصطلح مصدره وناشره ومعجمه، وقد تختلف الترجمات العربية للمصطلح الإنجليزي نفسه من مصدر إلى آخر. لذلك يعرض تقنون هذه الصيغ إلى جانب بعضها، ويحافظ على الإسناد الذي يسمح للقارئ بمعرفة من قال ماذا.

يناسب هذا الباحثين والمترجمين والطلاب والمحررين التقنيين الذين يحتاجون إلى رؤية سريعة لاستخدامات المصطلح عبر أكثر من مرجع، مع إمكانية العودة إلى المصدر الأصلي.

### ما الذي تقدمه الواجهة الحالية؟

- **بحث عربي/إنجليزي** مع خيار التطابق التام أو البحث كجزء من الكلمة، وخيار تضمين الوصف، وترشيح حسب الناشر أو المعجم، وتصفح النتائج على صفحات.
- **الصفحة الرئيسية** تعرض اقتراحات عشوائية للبحث من `/api/suggest`.
- **فهرس الناشرين والمعاجم (`/glossaries`)** يعرض الناشرين ومعاجمهم من `/api/publishers`، ومشاريع التوطين من `/api/localizations/sources`، مع تفاصيل الغلاف وروابط الكتاب وويكيداتا عند توفرها.
- **لوحة المعلومات (`/dashboard`)** تعرض إحصاءات إجمالية من `/api/stats`: إجمالي المصطلحات، والمصطلحات العربية والفرنسية والألمانية، وعدد المعاجم.
- **بحث الخادم عبر `/api/search`** هو المسار الذي ينفذ البحث الدقيق أو الجزئي في المصطلحات، والبحث الاختياري في الأوصاف. يرشح النتائج بحسب الناشرين (`publisher`) وأسماء المعاجم (`glossary`)، ويدعم `limit` و`offset`. يطابق الحقلين الإنجليزي والعربي فقط، ولا يطابق الفرنسي أو الألماني حتى لو كانت هذه الحقول معروضة في بطاقة المصطلح.
- **بحث اختياري في ترجمات البرمجيات مفتوحة المصدر** يتم عبر عامل ويب (Web Worker) داخل المتصفح. يطلب العامل `/api/localizations` من نفس أصل الموقع ثم يطابق النتائج محليًا؛ لذلك لا يمر هذا البحث عبر `/api/search`. أما جلب ملفات المصادر البعيدة (`.po`/`.txt`) فيتم على الخادم، ولذلك يقع اعتماد الشبكة وإمكانية الوصول إلى تلك المصادر على جانب الخادم لا على المتصفح.
- **المقارنة (Compare)** تتيح مقارنة 2 إلى 5 مصطلحات في وقت واحد. تحمّل الصفحة كامل مجموعة المصطلحات من `/api/terms` (من دون تصفح صفحات) وتجري المطابقة في المتصفح، مع خيار التطابق الجزئي والبحث في الوصف. لا تستخدم المقارنة `/api/search` في المطابقة.
- **البحث الإحصائي (Research)** يحمّل أيضًا كامل المصطلحات من `/api/terms`، ثم يحسب تكرار صيغ الترجمة العربية للمصطلح وتوزيع المصادر التي وردت فيها. وعند توفر بحث الترجمات مفتوحة المصدر، تُدمج نتائجه مع الإحصاءات المحلية. المطابقة هنا تتم في المتصفح أيضًا، وليست عبر `/api/search`.

| المسار | مصدر البيانات | مكان المطابقة |
| --- | --- | --- |
| البحث العام | `/api/search` | الخادم |
| ترجمات البرمجيات | `/api/localizations` | عامل ويب في المتصفح |
| المقارنة | كامل `/api/terms` | المتصفح |
| البحث الإحصائي | كامل `/api/terms` مع ترجمات اختيارية | المتصفح |

### البنية والتقنيات

- واجهة أمامية بـ **SvelteKit** و**Svelte 5**، تُبنى بـ Vite.
- خادم HTTP بـ **Go** وإطار **Fiber**.
- قاعدة بيانات **SQLite** محلية.
- خادم التطوير في Vite يمرر طلبات `/api` إلى `http://localhost:8080`.

### التشغيل من جذر المستودع

المتطلبات: إصدار Go المعلن في `go.mod` هو **1.25.9**، مع Node.js وnpm لبناء الواجهة. لا تتطلب عملية التشغيل المحلية أي بيانات اعتماد؛ المنفذ الافتراضي هو `8080` ويمكن تغييره بمتغير البيئة `PORT`. ولا تُضِف ملف `.env` إلى المستودع.

للبناء والتشغيل من جذر المشروع:

```sh
cd frontend
npm ci
npm run build
cd ..
go run .
```

ثم افتح `http://localhost:8080`.

عند بدء `go run .` من جذر المستودع يقوم الخادم بالخطوات التالية:

1. يقرأ ملف `.env` الاختياري إن وجد، مع أفضلية للمتغيرات الموجودة مسبقًا في البيئة.
2. يفتح قاعدة البيانات `./data.db` أو ينشئها.
3. يطبق مخطط قاعدة البيانات والترحيلات.
4. يستورد ملفات JSON المدرجة في `db/data/glossaries/index.json`.
5. يخدم واجهات `/api` والواجهة المبنية من `frontend/dist`، مع إعادة التوجيه إلى `index.html` لمسارات الواجهة (SPA fallback).

يلاحظ أن المسارات النسبية أعلاه تُحسب من دليل التشغيل الحالي؛ شغّل الأمر من جذر المستودع أو وفر هذه المسارات في بيئة التشغيل.

### تطوير الواجهة

شغّل الخادم الخلفي والواجهة في نافذتين منفصلتين.

النافذة الأولى، من جذر المشروع:

```sh
go run .
```

النافذة الثانية:

```sh
cd frontend
npm run dev
```

افتح العنوان الذي يطبعه Vite. ستمر طلبات `/api` تلقائيًا إلى خادم Go على المنفذ `8080` بفضل إعداد الوكيل في `frontend/vite.config.js`. أثناء تطوير الواجهة لا تحتاج إلى إعادة بناء `frontend/dist`، لكن النسخة التي يخدمها خادم Go على المنفذ `8080` تبقى معتمدة على آخر بناء ثابت.

### استيعاب البيانات

- ملف `db/data/glossaries/index.json` هو **المرجع الحصري** لقائمة الملفات المستوردة. لا يفترض هذا المستند عددًا ثابتًا من الملفات أو المعاجم أو الناشرين أو المصطلحات؛ المخزون الفعلي هو ما يسرده هذا الملف وملفات `fileData` الخاصة به.
- قد يمثل أكثر من ملف في القائمة أجزاءً موضوعية أو مجلدات من عمل واحد؛ ولذلك تُفصل هوية المعجم عن ملف المصدر.
- كل ملف JSON يحتوي كائن `fileData` لبيانات المصدر (المعجم والناشر والروابط) ومصفوفة `entries` للمصطلحات.
- يحسب الاستيراد بصمة **SHA-256** لكل ملف. إذا لم تتغير البصمة، يُتخطى الملف. وإذا تغير، يُستورد ضمن معاملة (transaction) تستبدل مجموعة مصطلحات ذلك المصدر، فلا تبقى نسخة قديمة مبعثرة عند الفشل.
- تسجل قاعدة البيانات مسار كل ملف مصدر وبصمته والمعجم المرتبط به ووقت الاستيراد في جدول `source_file`.

تسوية البحث تتم عند الاستيراد: تُحوّل الحروف اللاتينية إلى أحرف صغيرة، وتُزال التشكيلات العربية من الحقول العربية، وتُخزن نسخ معيارية من الحقول للبحث. لا يعني ذلك توحيد النصوص المعروضة أو حذف اختلافات المصادر.

### نموذج البيانات ومصطلحاته

| المفهوم | المعنى | ملاحظات |
| --- | --- | --- |
| **الناشر (Publisher)** | الجهة أو المؤسسة التي أصدرت المعجم | له اسم عربي وإنجليزي ورابط، وقد يتبع له أكثر من معجم. |
| **المعجم (Glossary)** | العمل المرجعي أو المسرد نفسه | له اسمان عربي وإنجليزي ورابط مصدر، وبيانات إثراء اختيارية مثل صورة الغلاف ورابط الكتاب وويكيداتا. |
| **المجال (Domain)** | التصنيف الموضوعي للمعجم | الجدول موجود في المخطط، لكن الاستيراد الحالي لا يملأ هذا الحقل. |
| **المصطلح (Term)** | المدخل المعجمي الواحد | يتضمن الإنجليزي والعربي، وقد يتضمن الفرنسي والألماني والوصف ورابط المصطلح. |
| **ملف المصدر (Source file)** | ملف JSON الذي جاءت منه البيانات | يحفظ المسار والبصمة والمعجم ووقت الاستيراد. |

يمكن أن توجد **صيغ لغوية إضافية** (فرنسية وألمانية) ووصف ورابط خاص بالمصطلح، لكنها اختيارية وتختلف من مصدر إلى آخر. كما أن ترخيص المعجم موجود كحقل في المخطط، إلا أن الاستيراد الحالي لا يضبطه. الإسناد الأساسي هو الناشر والمعجم والمصدر المعروض مع كل مصطلح.

المصادر متنوعة وقد تتعارض في الترجمة أو التهجئة أو مستوى التفصيل. لذلك فالمقارنة في تقنون **استكشافية**؛ هي تعرض ما ورد في المصادر ولا تمنح أي صيغة صفة معيارية أو نهائية.

للاطلاع على المخزون الحالي وبيانات كل مصدر، ارجع إلى `db/data/glossaries/index.json` وإلى كائن `fileData` داخل ملفات JSON. ويقدم `frontend/README.md` قائمة مختارة أقصر بالمصادر.

### واجهة API

جميع المسارات التالية من نوع `GET` وتعيد JSON. الغرض منها موضح أدناه، أما شكل الحقول التفصيلي فيتبع التنفيذ الحالي ولا يُفترض هنا مخططًا إضافيًا غير موثق.

| المسار | الوظيفة |
| --- | --- |
| `/api/publishers` | الناشرون ومعاجمهم. |
| `/api/terms` | جميع المصطلحات، دون تصفح صفحات. |
| `/api/search` | بحث في المصطلحات والأوصاف مع ترشيح وتصفح. |
| `/api/stats` | إحصاءات إجمالية عن المصطلحات والمعاجم واللغات. |
| `/api/suggest` | اقتراحات عشوائية مخزنة مؤقتًا من المدخلات التي تحتوي على قيمة من كلمة واحدة في العربية أو الإنجليزية على الأقل، وليس بالضرورة في اللغتين معًا. |
| `/api/localizations` | مدخلات يجلبها الخادم من ملفات ترجمة بعيدة بصيغتي `.po` و`.txt` ويحللها ويخزنها مؤقتًا قبل تقديمها للمتصفح. |
| `/api/localizations/sources` | مشاريع الترجمة مصنفة بحسب الناشر. |

#### معاملات `/api/search`

| المعامل | القيمة | السلوك |
| --- | --- | --- |
| `q` | نص البحث | إذا كان فارغًا تُعاد نتيجة فارغة. |
| `exact` | `true` أو `1` | مطابقة تامة للاسم الإنجليزي أو العربي بعد التسوية. |
| `desc` | `true` أو `1` | تضمين البحث في الوصف في وضع البحث الجزئي. |
| `publisher` | معرفات ناشرين رقمية مفصولة بفواصل | تهدف إلى قصر النتائج على هؤلاء الناشرين؛ راجع التنبيه أدناه. |
| `glossary` | أسماء معاجم عربية مفصولة بفواصل | مطابقة تامة لاسم `glossaryAr`. |
| `limit` | عدد صحيح | القيمة الافتراضية 10، وتُحصر بين 1 و100. |
| `offset` | عدد صحيح | القيمة الافتراضية 0، والقيم السالبة تصبح 0. |

**تنبيه حول مرشحات البحث:** يتوفر `publisher` و`glossary`، لكنهما لا يضمنان تقييد كل النتائج في التنفيذ الحالي لأن مرشحات SQL تُضاف من دون تجميع شرط OR. لذلك تنطبق المرشحات على آخر فرع من فروع OR فقط: فرع الحقل العربي في وضعي التطابق التام والبحث الجزئي من دون `desc`؛ وعند تفعيل `desc` في وضع البحث الجزئي تنطبق على فرع الوصف فقط. أما مطابقات الحقل الإنجليزي فتتجاوز المرشحات في جميع الأوضاع، كما تتجاوزها مطابقات الحقل العربي في وضع البحث الجزئي مع `desc`. هذه محدودية معروفة في التنفيذ الحالي، لذا تعامل مع النتائج المفلترة بحذر.

`/api/search` يبحث في الحقلين الإنجليزي والعربي (والوصف عند تفعيل `desc`)، وليس في الفرنسي أو الألماني. أما `/api/localizations` فيجلب الخادم ملفات الترجمة من مصادر بعيدة ويحللها ويخزنها مؤقتًا؛ لذلك يتأثر توفر هذه النتائج بتوفر الشبكة وصلاحية تلك المصادر لدى الخادم.

### النشر والقيود الحالية

- لا تتضمن الواجهة الحالية تسجيل دخول أو إدارة حسابات مستخدمين.
- يوجد `frontend/Dockerfile`، لكنه مخصص للواجهة الأمامية فقط، ويعتمد على `node build` بينما مخرجات البناء الحالية عبر المحول الثابت (`adapter-static`) هي `frontend/dist`. لذلك لا ينبغي اعتباره مسارًا متكاملًا وجاهزًا لنشر الخادم والواجهة معًا؛ نشر الإنتاج يحتاج إلى مطابقة مخرجات البناء مع طريقة التشغيل أو إلى إعداد حاويات مصحح.
- تُعرض بيانات الإسناد كما وردت في المصادر. راجع صفحة كل مصدر لمعرفة شروط استخدام محتواه؛ هذا المستند يشير إلى بيانات الإسناد ولا يحدد حقوق إعادة الاستخدام.
- اختلاف المصادر ليس خللًا في البيانات، بل جزء من طبيعة المعاجم المتعددة التي يقوم عليها تقنون.

### أين تجد ماذا

| المسار | المحتوى |
| --- | --- |
| `main.go` | نقطة البدء، تسجيل مسارات API، وخدمة الواجهة الثابتة. |
| `handlers/` | معالجات واجهة API. |
| `ingest/` | قراءة بيان المصادر، والتحقق من البصمات، والاستيراد. |
| `internal/` | حزم التطبيع (`normalize`) وإثراء البيانات الوصفية (`metadata`). |
| `localizations/` | جلب ملفات الترجمة البعيدة بصيغتي `.po`/`.txt` وتحليلها وتخزينها مؤقتًا خلف مسارات التوطين. |
| `db/` | مخطط SQLite والترحيلات والاتصال. |
| `models/` | أشكال البيانات التي تعيدها الـ API. |
| `frontend/` | تطبيق SvelteKit ومخرجاته المبنية. |
| `db/data/glossaries/` | ملفات المصادر وبيانها `index.json`. |

### الاختبارات

من جذر المستودع:

```sh
go test ./...
```

توجد ملفات اختبار حاليًا في `internal/normalize` و`internal/metadata` فقط. لا توجد ملفات اختبار في `ingest` أو `handlers` أو `localizations`، لذلك لا تغطي الاختبارات الحالية الاستيراد أو معالجات API أو جلب التوطين. ينبغي أن تبقى هذه الاختبارات ناجحة قبل أي تغيير في التطبيع أو إثراء البيانات الوصفية.

---

## English

### What is TaQnun?

TaQnun is a documentation and research project that brings Arabic technical and scientific terminology together from separate dictionaries and glossaries and presents it in one place for searching and comparison. Its aim is not to produce a single "normative" translation, but to preserve the diversity of usage: every term carries its source, publisher, and glossary, and different sources may translate the same English term differently. TaQnun therefore surfaces those variants side by side and keeps the attribution needed to see who said what.

It is useful for researchers, translators, students, and technical editors who need a quick view of how a term is used across multiple references, with a path back to the original source.

### What the current interface provides

- **Arabic/English search** with exact or substring matching, an optional description lookup, filtering by publisher or glossary, and paginated results.
- **Homepage** displays random search suggestions from `/api/suggest`.
- **Publisher and glossary catalog (`/glossaries`)** lists publishers and their glossaries from `/api/publishers`, plus localization projects from `/api/localizations/sources`, with cover details and book and Wikidata links when available.
- **Dashboard (`/dashboard`)** shows aggregate statistics from `/api/stats`: total terms, Arabic, French, and German terms, and the number of glossaries.
- **Server-side search through `/api/search`** is the path that performs exact or substring term lookup and optional description lookup. It filters by publisher IDs (`publisher`) and Arabic glossary names (`glossary`) and supports `limit` and `offset`. It matches the English and Arabic fields only, not French or German, even though those fields may appear on a term card.
- **Optional open-source localization lookup** runs through a browser Web Worker. The worker requests same-origin `/api/localizations` and matches the results locally, so this lookup does not go through `/api/search`. The server performs the remote fetches of the `.po`/`.txt` source files, so the network dependency and remote-source availability apply to server-side fetching, not the browser.
- **Compare** lets you compare 2 to 5 terms at once. The page loads the full term set from `/api/terms` (unpaginated) and matches in the browser, with options for partial matching and description search. Compare does not use `/api/search` for its matching.
- **Research** also loads the full term set from `/api/terms`, then calculates how often each Arabic translation variant occurs and which sources use it. When the optional localization lookup is available, its results are merged with the local statistics. Matching happens in the browser here as well, not through `/api/search`.

| Path | Data source | Matching location |
| --- | --- | --- |
| General search | `/api/search` | Server |
| Open-source localizations | `/api/localizations` | Browser Web Worker |
| Compare | Full `/api/terms` | Browser |
| Research | Full `/api/terms` plus optional localizations | Browser |

### Architecture and stack

- **SvelteKit** and **Svelte 5** frontend, built with Vite.
- **Go** HTTP API using **Fiber**.
- Local **SQLite** database.
- The Vite development server proxies `/api` requests to `http://localhost:8080`.

### Running from the repository root

Requirements: the Go version declared in `go.mod` is **1.25.9**, plus Node.js and npm for the frontend. No credentials are required for a local run; the default port is `8080` and can be changed with the `PORT` environment variable. Do not commit `.env` to the repository.

To build and run from the project root:

```sh
cd frontend
npm ci
npm run build
cd ..
go run .
```

Then open `http://localhost:8080`.

When `go run .` starts from the repository root, the server does the following:

1. Reads an optional `.env` file if present; values already present in the environment take precedence.
2. Opens `./data.db`, creating it if needed.
3. Applies the database schema and migrations.
4. Ingests the JSON files listed in `db/data/glossaries/index.json`.
5. Serves the `/api` routes and the built frontend from `frontend/dist`, with an SPA fallback to `index.html` for client-side routes.

The relative paths above are resolved from the current working directory; run the command from the repository root or provide those paths in the runtime environment.

### Frontend development

Run the backend and frontend in two separate terminals.

Terminal one, from the project root:

```sh
go run .
```

Terminal two:

```sh
cd frontend
npm run dev
```

Open the address printed by Vite. Requests to `/api` are proxied automatically to the Go server on port `8080` by the configuration in `frontend/vite.config.js`. While developing the UI you do not need to rebuild `frontend/dist`, but the version served by the Go server on port `8080` remains tied to the last static build.

### Data ingestion

- `db/data/glossaries/index.json` is the **authoritative** list of ingested files. This document does not assume a fixed number of files, dictionaries, glossaries, publishers, or terms; the actual inventory is whatever that manifest and each file's `fileData` describe.
- Multiple manifest files may be topic slices or volumes of one logical work; glossary identity is therefore kept separate from source-file identity.
- Each JSON file contains a `fileData` object with source metadata (glossary and publisher names and links) and an `entries` array of terms.
- Ingestion computes a **SHA-256** hash for each file. If the hash is unchanged, the file is skipped. If it has changed, ingestion runs in a transaction and replaces that source's term set, so a partial failure does not leave a stale mixture behind.
- The database records each source file's path, hash, linked glossary, and ingestion time in the `source_file` table.

Search normalization happens during ingestion: Latin text is lowercased, Arabic diacritics are removed from Arabic fields, and normalized versions of the fields are stored for lookup. This does not flatten the displayed text or erase differences between sources.

### Data model and terminology

| Concept | Meaning | Notes |
| --- | --- | --- |
| **Publisher** | The organization or institution that issued the glossary | Has Arabic and English names and a URL, and may own more than one glossary. |
| **Glossary** | The reference work or word list itself | Has Arabic and English names and a source URL, plus optional enrichment such as a cover image, book link, and Wikidata link. |
| **Domain** | The subject classification of the glossary | The table exists in the schema, but current ingestion leaves this field unset. |
| **Term** | A single dictionary entry | Contains English and Arabic, and may contain French, German, a description, and a term URL. |
| **Source file** | The JSON file the data came from | Stores the path, hash, linked glossary, and ingestion time. |

**Additional language variants** (French and German), a description, and a term-specific link are optional and vary by source. The glossary license also exists as a schema field, but current ingestion does not set it. The primary attribution is the publisher, glossary, and source shown with each term.

Sources are varied and may disagree on translation, spelling, or level of detail. Comparison in TaQnun is therefore **exploratory**: it shows what the sources say and does not treat any variant as canonical or final.

For the current inventory and per-source metadata, see `db/data/glossaries/index.json` and the `fileData` object in each JSON file. `frontend/README.md` provides a shorter selected-source list.

### API

All routes below are `GET` and return JSON. Their purposes are described here; exact field shapes follow the current implementation, and no additional undocumented schema is implied.

| Route | Purpose |
| --- | --- |
| `/api/publishers` | Publishers and their glossaries. |
| `/api/terms` | All terms, unpaginated. |
| `/api/search` | Term and description search with filters and pagination. |
| `/api/stats` | Aggregate counts for terms, glossaries, and languages. |
| `/api/suggest` | Random cached suggestions from entries with a single-word value in at least one of Arabic or English, not necessarily a single word in both. |
| `/api/localizations` | Entries fetched, parsed, and cached by the server from remote `.po` and `.txt` localization files before being served to the browser. |
| `/api/localizations/sources` | Localization projects grouped by publisher. |

#### `/api/search` parameters

| Parameter | Value | Behavior |
| --- | --- | --- |
| `q` | Search text | An empty query returns an empty result. |
| `exact` | `true` or `1` | Exact match on the normalized English or Arabic name. |
| `desc` | `true` or `1` | Includes description search in non-exact mode. |
| `publisher` | Comma-separated numeric publisher IDs | Intended to restrict results to those publishers; see caveat below. |
| `glossary` | Comma-separated Arabic glossary names | Exact match on `glossaryAr`. |
| `limit` | Integer | Defaults to 10 and is clamped to 1–100. |
| `offset` | Integer | Defaults to 0; negative values become 0. |

**Search filter caveat:** `publisher` and `glossary` are available, but in the current implementation they do not reliably constrain every result because SQL filters are appended without grouping the OR predicate. Filters therefore apply only to the last OR branch: the Arabic-field branch in exact and substring modes without `desc`; when `desc` is enabled in substring mode, they apply only to the description branch. English matches bypass the filters in all modes; Arabic matches also bypass them in substring-with-desc mode. This is a known implementation limitation, so treat filtered results with care.

`/api/search` searches the English and Arabic fields (and the description when `desc` is enabled), not French or German. `/api/localizations` is backed by the server fetching, parsing, and caching localization files from remote sources, so the availability of those results depends on network availability and the health of those sources on the server side.

### Deployment and current limitations

- The current interface does not include login or user account management.
- `frontend/Dockerfile` exists, but it is frontend-only and runs `node build`, while the current documented frontend build uses the static adapter and outputs to `frontend/dist`. It should not be treated as a complete, ready-to-use full-stack deployment path; production deployment needs the build output and runtime to match, or a corrected container setup.
- Attribution is shown as provided by the sources. Check each source's page for the terms that apply to its content; this document points to attribution data and does not determine reuse rights.
- Differences between sources are not a data defect; they are part of the multi-dictionary nature that TaQnun is designed around.

### Where things live

| Path | Contents |
| --- | --- |
| `main.go` | Entry point, API route registration, and static frontend serving. |
| `handlers/` | API handlers. |
| `ingest/` | Manifest reading, hash checks, and ingestion. |
| `internal/` | Normalization (`normalize`) and metadata enrichment (`metadata`) packages. |
| `localizations/` | Remote `.po`/`.txt` localization fetching, parsing, and caching behind the localization routes. |
| `db/` | SQLite schema, migrations, and connection setup. |
| `models/` | Data shapes returned by the API. |
| `frontend/` | SvelteKit application and its build output. |
| `db/data/glossaries/` | Source files and the `index.json` manifest. |

### Tests

From the repository root:

```sh
go test ./...
```

Test files currently exist only in `internal/normalize` and `internal/metadata`. There are no test files in `ingest`, `handlers`, or `localizations`, so the current tests do not cover ingestion, the API handlers, or localization fetching. These tests should remain green before changes to normalization or metadata.

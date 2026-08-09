-- +goose Up
-- +goose StatementBegin
UPDATE content SET data = data || '{
  "bioRu": "Fullstack-разработчик: PHP 8.4/Laravel и React/TypeScript. До разработки несколько лет занималась серверной инфраструктурой — поэтому мои проекты доезжают до продакшена вместе с Docker, CI/CD и мониторингом, а не остаются в репозитории. Веду open-source тикет-систему TicketHub. Открыта к позициям fullstack-разработчика и системного аналитика.",
  "bioEn": "Full-stack developer working in PHP 8.4/Laravel and React/TypeScript. I spent several years on server infrastructure before development — which is why my projects reach production with Docker, CI/CD, and monitoring instead of staying in the repo. I maintain TicketHub, an open-source ticketing system. Open to full-stack and systems-analyst roles.",
  "bioFullRu": "Пишу на PHP/Laravel и React/TypeScript. До разработки несколько лет занималась серверной инфраструктурой, поэтому мои проекты доезжают до продакшена вместе с Docker, CI/CD и мониторингом, а не остаются в репозитории.\n\nНачинала младшим сисадмином в ноябре 2021 в Российской государственной детской библиотеке. 100+ рабочих мест на Windows, Linux и macOS, легаси-сайты на PHP, которых годами никто не трогал. Первые скрипты на Bash и Python написала в тот же месяц — если одно и то же приходилось делать руками дважды, писала скрипт. Эта привычка осталась. Веб-разработка пришла оттуда же: сайты организации ломались, а чинить их было некому — я исправила около 30 багов и ускорила главный сайт на 20%.\n\nС марта 2024 — веб-разработчик и системный администратор там же, руковожу командой из двух админов. Внедрила CI/CD на GitHub Actions: деплой внутренних сервисов упал с 2 часов до 15 минут, обновления перестали требовать простоя. Держу SLA 99,9% на Prometheus + Grafana. Собрала библиотеку из 20+ скриптов автоматизации на Python, Bash и PowerShell — рутина быстрее на 40%, ручных ошибок меньше на 60%, команде это экономит около 15 часов в неделю.\n\nС ноября 2024 параллельно веду коммерческие проекты под ключ — от постановки задачи до работающего продакшена с мониторингом. Флагман — TicketHub, open-source тикет-система на PHP 8.4 под GPL: REST API на ~45 эндпоинтов, security-слой по OWASP, N+1 в канбане с 201 запроса до 2. Переносила легаси с PHP 5.x и jQuery на PHP 8.x и React, подключала ЮKassa, Т-Банк, amoCRM и Bitrix24. Всё это живёт на моём VPS вместе с ещё тремя приложениями за общим Caddy с health-gate и авто-откатом.\n\nВ 2026 закончила бакалавриат по прикладной информатике (ИИ и анализ данных) в МУ им. Витте. Опубликовала статью в журнале «Моделирование, оптимизация и информационные технологии» (перечень ВАК, категория К2) — ансамблевые методы в диагностике сердечно-сосудистых заболеваний, CatBoost с ROC-AUC 0,948 на выборке из 1 904 пациентов.\n\nМой пунктик — автоматизация: если что-то делается руками дважды, на третий раз это делает скрипт. Честное предупреждение — перфекционистка. Иногда сложнее не написать код, а дать проекту выйти до того, как я отполирую каждую мелочь.",
  "bioFullEn": "I write PHP/Laravel and React/TypeScript. Before development I spent several years on server infrastructure, which is why my projects reach production together with Docker, CI/CD, and monitoring instead of staying in the repo.\n\nI started as a junior sysadmin in November 2021 at the Russian State Children''s Library. 100+ workstations on Windows, Linux, and macOS, and legacy PHP sites nobody had touched in years. I wrote my first Bash and Python scripts that same month — every time I did the same thing by hand twice, I scripted it. That habit stuck. Web development came from the same place: the organization''s sites kept breaking and there was no one to fix them, so I fixed about 30 bugs and sped up the main site by 20%.\n\nSince March 2024 I''ve been a web developer and system administrator there, leading a team of two admins. I rolled out CI/CD on GitHub Actions: deploys of internal services dropped from 2 hours to 15 minutes and updates stopped requiring downtime. I hold a 99.9% SLA on Prometheus + Grafana. I built a library of 20+ automation scripts in Python, Bash, and PowerShell — routine work is 40% faster, manual errors are down 60%, and the team saves around 15 hours a week.\n\nSince November 2024 I have been running commercial projects end to end in parallel — from the brief to a working production with monitoring. The flagship is TicketHub, an open-source PHP 8.4 ticketing system under GPL: a ~45-endpoint REST API, an OWASP-aligned security layer, and Kanban N+1 cut from 201 queries to 2. I migrated legacy from PHP 5.x and jQuery to PHP 8.x and React, and integrated YooKassa, T-Bank, amoCRM, and Bitrix24. All of it runs on my VPS alongside three other apps behind a shared Caddy with a health gate and auto-rollback.\n\nIn 2026 I finished my bachelor''s in Applied Informatics (AI & Data Analysis) at Moscow Witte University. I published a paper in Modeling, Optimization and Information Technology (Russian HAC list, category K2) on ensemble methods for cardiovascular disease diagnosis — CatBoost at 0.948 ROC-AUC on a sample of 1,904 patients.\n\nMy thing is automation: if something gets done by hand twice, the third time a script does it. Honest caveat — I''m a perfectionist. Sometimes the hard part isn''t writing the code; it''s letting the project ship before I''ve polished every edge.",
  "stats": {
    "projects": 11,
    "githubStars": 0,
    "certificates": 1,
    "habrArticles": 0,
    "yearsExperience": 4
  }
}'::jsonb
WHERE section = 'about';
-- +goose StatementEnd

-- +goose StatementBegin
UPDATE content SET data = data || '{
  "subtitleRu": "Fullstack-разработчик: PHP/Laravel и React/TypeScript. Пришла из системного администрирования — довожу проекты до продакшена с Docker, CI/CD и мониторингом.",
  "subtitleEn": "Full-stack developer: PHP/Laravel and React/TypeScript. I came from system administration — I take projects all the way to production with Docker, CI/CD, and monitoring.",
  "titleAnimated": {
    "ru": ["Fullstack-разработчик", "Системный аналитик", "ML-исследователь", "Системный администратор"],
    "en": ["Full-Stack Developer", "Systems Analyst", "ML Researcher", "System Administrator"]
  }
}'::jsonb
WHERE section = 'hero';
-- +goose StatementEnd

-- +goose StatementBegin
UPDATE content SET data = data || '{
  "subtitleRu": "Рассматриваю позиции fullstack-разработчика и системного аналитика — фултайм, удалённо или гибрид. Беру и фриланс-проекты. Отвечаю в течение суток.",
  "subtitleEn": "Open to full-stack and systems-analyst roles — full-time, remote or hybrid — and to freelance projects. I reply within a day."
}'::jsonb
WHERE section = 'contactCTA';
-- +goose StatementEnd

UPDATE content SET data = '[
  "PHP 8.4", "Laravel 12", "React 19", "Next.js 16", "TypeScript", "Go 1.25",
  "Python", "FastAPI", "PostgreSQL", "MySQL", "Redis", "Docker",
  "GitHub Actions", "CI/CD", "Caddy", "Prometheus", "REST API",
  "Meilisearch", "CatBoost", "Telegram Bots"
]'::jsonb
WHERE section = 'marqueeItems';

-- +goose StatementBegin
UPDATE work_history SET
  position = '{"en": "Full-Stack Developer", "ru": "Fullstack-разработчик"}'::jsonb,
  start_date = '2024-11-01T00:00:00Z',
  is_current = TRUE,
  technologies = ARRAY['PHP 8.4','Laravel 12','React 19','TypeScript','Next.js','Go','Python','FastAPI','MySQL','PostgreSQL','Redis','Docker','GitHub Actions','Caddy'],
  description = '{
    "en": "Commercial projects end to end — from the brief to a running production with monitoring. Flagship: TicketHub, an open-source PHP 8.4 ticketing system — a ~45-endpoint REST API, an OWASP-aligned security layer, and Kanban N+1 cut from 201 queries to 2 per page.",
    "ru": "Веду коммерческие проекты под ключ: от постановки задачи до работающего продакшена с мониторингом. Флагман — TicketHub, open-source тикет-система на PHP 8.4: REST API на ~45 эндпоинтов, security-слой по OWASP, N+1 в канбане с 201 запроса до 2 на страницу."
  }'::jsonb,
  full_description = '{
    "en": "The freelance track runs in parallel with the main job at RGDB — evenings and weekends. Projects range from small Telegram bots to multi-service web apps. Weekly iterations, GitHub from day one, a demo at the end of each iteration. The portfolio on this site lists the open, shareable pieces — commercial client work stays under NDA by default. Payment model: 50/50 (upfront + on delivery), 5 free revisions, then paid. The full process is described on the homepage under \"How I Work\".",
    "ru": "Фриланс идёт параллельно основной работе в РГДБ — вечера и выходные. Проекты от маленьких Telegram-ботов до многосервисных веб-приложений. Недельные итерации, GitHub с первого дня, демо в конце каждой итерации. В портфолио на этом сайте — открытые вещи, которыми можно делиться; коммерческая работа для клиентов по умолчанию под NDA. Модель оплаты: 50/50 (предоплата + при сдаче), 5 бесплатных правок, дальше платно. Полное описание процесса — на главной в блоке «Как я работаю»."
  }'::jsonb,
  achievements = '[
    {
      "en": "Designed and shipped TicketHub — an open-source PHP 8.4 ticketing system under GPL: REST API v1 with ~45 endpoints and 947 lines of documentation, live in production at tickethub.lavier.tech.",
      "ru": "Спроектировала и выпустила TicketHub — open-source тикет-систему на PHP 8.4 под GPL: REST API v1 на ~45 эндпоинтов и 947 строк документации, в продакшене на tickethub.lavier.tech."
    },
    {
      "en": "Built an OWASP-aligned security layer: SHA-256 token hashing, granular permissions, CIDR IP whitelist, SSRF protection, and brute-force lockout.",
      "ru": "Собрала security-слой по OWASP: SHA-256-хеширование токенов, granular permissions, CIDR IP whitelist, защита от SSRF и brute-force lockout."
    },
    {
      "en": "Optimized N+1 queries in the Kanban module: 201 → 2 SQL queries per page.",
      "ru": "Оптимизировала N+1-запросы в канбан-модуле: 201 → 2 SQL-запроса на страницу."
    },
    {
      "en": "Built 5+ web apps from scratch on React + TypeScript and PHP/Laravel (Node.js, PostgreSQL, MySQL): online stores, customer portals, admin panels, Telegram bots.",
      "ru": "Сделала с нуля 5+ веб-приложений на React + TypeScript и PHP/Laravel (Node.js, PostgreSQL, MySQL): интернет-магазины, личные кабинеты, админ-панели, Telegram-боты."
    },
    {
      "en": "Migrated a legacy site off a dying PHP host three weeks before shutdown: Docker copy of prod, DB migrations, critical modules rewritten in PHP 8.4 strict types. The cutover fit into a 4-minute window; the single bug that surfaced was rolled back and fixed within the hour.",
      "ru": "Перевела legacy-сайт со старого PHP-хостинга за 3 недели до его закрытия: Docker-копия прода, миграции БД, критичные модули переписаны на PHP 8.4 strict types. Само переключение уложилось в 4-минутное окно; единственный всплывший баг откатили и починили в течение часа."
    },
    {
      "en": "Migrated three more sites from PHP 5.x and jQuery to PHP 8.x and React + Tailwind — load times 40–60% faster.",
      "ru": "Перенесла ещё 3 сайта с PHP 5.x и jQuery на PHP 8.x и React + Tailwind — загрузка ускорилась на 40–60%."
    },
    {
      "en": "Integrated payments and CRMs: YooKassa, T-Bank (Tinkoff), amoCRM, Bitrix24, OAuth 2.0, plus email, SMS, and Telegram notifications.",
      "ru": "Подключала платежи и CRM: ЮKassa, Т-Банк (Тинькофф), amoCRM, Bitrix24, OAuth 2.0; уведомления в email, SMS и Telegram."
    },
    {
      "en": "Set up production infrastructure end to end: VPS, Docker, GitHub Actions, Caddy/Nginx with automatic SSL, database backups, and healthchecks alerting the owner in Telegram.",
      "ru": "Разворачивала production-инфраструктуру под ключ: VPS, Docker, GitHub Actions, Caddy/Nginx с автоматическим SSL, бэкапы БД, healthcheck с алертами владельцу в Telegram."
    }
  ]'::jsonb
WHERE id = '169b6f34-1813-4378-b7aa-ac1292b446a4';
-- +goose StatementEnd

-- +goose StatementBegin
UPDATE work_history SET
  position = '{"en": "Web Developer / System Administrator", "ru": "Веб-разработчик / системный администратор"}'::jsonb,
  technologies = ARRAY['PHP','Python','Bash','PowerShell','Docker','GitHub Actions','CI/CD','Prometheus','Grafana','Linux','Windows Server','Nginx','Git'],
  description = '{
    "en": "I maintain and extend the library''s PHP web services and rolled out CI/CD on GitHub Actions — deploys went from 2 hours to 15 minutes. I own the 99.9% SLA for the server infrastructure (Prometheus + Grafana) and lead a team of two system administrators.",
    "ru": "Сопровождаю и дорабатываю веб-сервисы библиотеки на PHP, внедрила CI/CD на GitHub Actions — деплой сократился с 2 часов до 15 минут. Отвечаю за SLA 99,9% серверной инфраструктуры (Prometheus + Grafana) и руковожу командой из двух системных администраторов."
  }'::jsonb,
  full_description = '{
    "en": "Second role at the Russian State Children''s Library. Promoted from junior sysadmin in March 2024. I combine development and infrastructure: the library''s PHP web services, CI/CD, monitoring, and the internal network, plus a team of two administrators. The focus shifted from ticket-by-ticket work to building systems that don''t generate tickets in the first place.",
    "ru": "Вторая роль в РГДБ — повысили из младших сисадминов в марте 2024. Совмещаю разработку и инфраструктуру: веб-сервисы библиотеки на PHP, CI/CD, мониторинг, внутренняя сеть, плюс команда из двух администраторов. Фокус сместился с работы по тикетам на построение систем, которые тикеты не генерируют в принципе."
  }'::jsonb,
  achievements = '[
    {
      "en": "Maintain and extend the library''s web services in PHP — fixes, new modules, incident analysis.",
      "ru": "Сопровождаю и дорабатываю веб-сервисы библиотеки на PHP: правки, новые модули, разбор инцидентов."
    },
    {
      "en": "Rolled out CI/CD on GitHub Actions for internal web services: deploy time dropped from 2 hours to 15 minutes, and updates stopped requiring downtime.",
      "ru": "Внедрила CI/CD на GitHub Actions для внутренних веб-сервисов: деплой сократился с 2 часов до 15 минут, обновления перестали требовать простоя."
    },
    {
      "en": "Hold a 99.9% SLA for the server infrastructure: Prometheus + Grafana, dashboards, alerts on CPU, memory, disk, and network.",
      "ru": "Отвечаю за SLA 99,9% серверной инфраструктуры: Prometheus + Grafana, дашборды, алерты по CPU, памяти, дискам и сети."
    },
    {
      "en": "Lead an IT team of 2 system administrators: project allocation, 1-on-1s, mentoring a junior specialist.",
      "ru": "Руковожу IT-командой из 2 системных администраторов: распределяю проекты, провожу 1-on-1, менторю junior-специалиста."
    },
    {
      "en": "Built a library of 20+ automation scripts in Python, Bash, and PowerShell: routine operations ~40% faster, manual errors down ~60%, saving the team around 15 hours a week.",
      "ru": "Собрала библиотеку из 20+ скриптов автоматизации на Python, Bash и PowerShell: рутинные операции ускорились на ~40%, ручных ошибок меньше на ~60%, команде это экономит около 15 часов в неделю."
    },
    {
      "en": "Rebuilt the network infrastructure from scratch — network map, tidy patch cabinets, cable labeling. Network diagnostics 25% faster, MTTR down 30%.",
      "ru": "Реорганизовала сетевую инфраструктуру с нуля: карта сети, порядок в коммутационных шкафах, маркировка кабелей. Диагностика сетевых проблем быстрее на 25%, MTTR ниже на 30%."
    },
    {
      "en": "Set up Docker for dev environments — a new developer spins up a test bench in 10 minutes instead of half a day of manual setup.",
      "ru": "Настроила Docker для dev-окружений: новый разработчик поднимает тестовый стенд за 10 минут вместо полдня ручной настройки."
    },
    {
      "en": "Ran a technical audit of museum-space management systems and wrote the architecture documentation from scratch — the organization stopped depending on the vendor.",
      "ru": "Провела технический аудит систем управления музейными пространствами и написала документацию архитектуры с нуля: организация перестала зависеть от вендора."
    },
    {
      "en": "Built an IT knowledge base of 50+ articles — new hires reach solo work in one week instead of three.",
      "ru": "Выстроила базу знаний на 50+ статей — новые сотрудники выходят на самостоятельную работу за неделю вместо трёх."
    }
  ]'::jsonb
WHERE id = '927b49e1-cfe9-4891-9a57-5e0e4a4542d6';
-- +goose StatementEnd

-- +goose StatementBegin
UPDATE work_history SET
  achievements = jsonb_set(
    achievements,
    '{5}',
    '{
      "en": "Maintained and patched the org''s web services: legacy PHP, HTML, CSS. Fixed ~30 bugs and sped up the main site by 20%. This was my first commercial web-development experience.",
      "ru": "Доработала и чинила веб-сервисы организации: legacy-код на PHP, HTML, CSS. Исправила ~30 багов, ускорила загрузку основного сайта на 20%. Это был мой первый коммерческий опыт в веб-разработке."
    }'::jsonb
  )
WHERE id = '796a6651-a01f-46c9-93c4-3da245d9468b'
  AND jsonb_array_length(achievements) = 6;
-- +goose StatementEnd

-- +goose StatementBegin
UPDATE education SET
  description = '{
    "en": "Studied alongside work. Focus on machine learning, data engineering, and applied AI. Coursework: predictive maintenance for a metro-train air compressor on the MetroPT-3 dataset (ROC-AUC 0.993). The research on ensemble methods for cardiovascular disease diagnosis was published in a Russian HAC-listed journal (MOIT, category K2, 2026).",
    "ru": "Училась параллельно с работой. Фокус — машинное обучение, дата-инженерия и прикладной ИИ. Курсовая — predictive maintenance: раннее детектирование утечек воздуха в компрессоре поезда метро на датасете MetroPT-3 (ROC-AUC 0,993). Исследование по ансамблевым методам в диагностике сердечно-сосудистых заболеваний вышло статьёй в журнале из перечня ВАК (МОИТ, категория К2, 2026)."
  }'::jsonb,
  related_project_slugs = ARRAY['heart-disease-ml-benchmark','metropt-predictive-maintenance','genreneuro','children-book-dataset']
WHERE institution->>'ru' LIKE '%Витте%';
-- +goose StatementEnd

UPDATE certificates SET issue_date = '2026-01'
WHERE id = '893315c2-6271-4441-bc3e-6dbb5bd39888';

-- +goose StatementBegin
INSERT INTO publications (id, title, journal, year, doi, url, abstract, sort_order) VALUES (
  'd3e8b1a7-2c94-4f60-a5d8-7b1e6c40f923',
  '{
    "en": "Ensemble Machine Learning Methods for Predictive Diagnostics of Cardiovascular Diseases: Comparative Analysis on a Multi-Center Dataset",
    "ru": "Ансамблевые методы машинного обучения для прогностической диагностики сердечно-сосудистых заболеваний: сравнительный анализ на многоцентровой выборке"
  }'::jsonb,
  '{
    "en": "Modeling, Optimization and Information Technology (MOIT), 2026, Vol. 14, No. 6(57) · Russian HAC list, category K2",
    "ru": "Моделирование, оптимизация и информационные технологии (МОИТ), 2026, т. 14, № 6(57) · перечень ВАК, категория К2"
  }'::jsonb,
  '2026',
  '10.26102/2310-6018/2026.57.6.017',
  'https://doi.org/10.26102/2310-6018/2026.57.6.017',
  '{
    "en": "Eight machine-learning algorithms compared on a combined multi-center sample from six databases (n = 1,904). CatBoost reached ROC-AUC 0.948 [0.922–0.966] and F1 0.884 with the best calibration; an ablation study showed that seven features retain 97.5% of full-model performance. Pairwise comparison used DeLong and McNemar tests with Bonferroni correction. Co-authors: D.I. Veselov, N.A. Andriyanov.",
    "ru": "Сравнение восьми алгоритмов машинного обучения на объединённой многоцентровой выборке из шести баз данных (n = 1 904). CatBoost показал ROC-AUC 0,948 [0,922–0,966] и F1 0,884 при лучшей калибровке; аблация показала, что 7 признаков сохраняют 97,5% качества полной модели. Попарное сравнение — тесты ДеЛонга и МакНемара с поправкой Бонферрони. Соавторы: Веселов Д.И., Андриянов Н.А."
  }'::jsonb,
  1
)
ON CONFLICT (id) DO UPDATE SET
  title = EXCLUDED.title,
  journal = EXCLUDED.journal,
  year = EXCLUDED.year,
  doi = EXCLUDED.doi,
  url = EXCLUDED.url,
  abstract = EXCLUDED.abstract,
  sort_order = EXCLUDED.sort_order;
-- +goose StatementEnd

-- +goose StatementBegin
UPDATE projects SET
  description = jsonb_set(
    jsonb_set(
      description,
      '{ru}',
      to_jsonb(
        replace(
          replace(description->>'ru',
            'Сравнительный анализ алгоритмов машинного обучения для диагностики сердечно-сосудистых заболеваний на многоцентровой выборке',
            'Ансамблевые методы машинного обучения для прогностической диагностики сердечно-сосудистых заболеваний: сравнительный анализ на многоцентровой выборке'),
          'на рецензии',
          'опубликована в журнале «Моделирование, оптимизация и информационные технологии» (перечень ВАК, категория К2), 2026, т. 14, № 6(57), DOI 10.26102/2310-6018/2026.57.6.017'
        )
      )
    ),
    '{en}',
    to_jsonb(
      replace(
        replace(description->>'en',
          'Сравнительный анализ алгоритмов машинного обучения для диагностики сердечно-сосудистых заболеваний на многоцентровой выборке',
          'Ансамблевые методы машинного обучения для прогностической диагностики сердечно-сосудистых заболеваний: сравнительный анализ на многоцентровой выборке'),
        'is under peer review',
        'was published in Modeling, Optimization and Information Technology (Russian HAC list, category K2), 2026, Vol. 14, No. 6(57), DOI 10.26102/2310-6018/2026.57.6.017'
      )
    )
  )
WHERE slug = 'heart-disease-ml-benchmark';
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO projects (
  id, slug, title, short_desc, description, category, status, tags, tech_stack,
  goal_desc, github_url, demo_url, site_url, thumbnail, featured, sort_order,
  problem, approach, outcome, tech_choices, highlights,
  timeline_started, timeline_shipped
) VALUES (
  'b1f4c8a2-5d3e-4a7b-9c61-0e2f7a4d8b13',
  'lavier-tech',
  '{"en": "lavier.tech", "ru": "lavier.tech"}'::jsonb,
  '{
    "en": "The site you are reading: Next.js 16 + Go 1.25, all content bilingual in Postgres JSONB, an admin panel behind it, push-based ISR, and a deploy with a health gate and auto-rollback",
    "ru": "Сайт, который вы сейчас читаете: Next.js 16 + Go 1.25, весь контент двуязычный в Postgres JSONB, за ним админка, push-ревалидация ISR и деплой с health-gate и авто-откатом"
  }'::jsonb,
  '{
    "en": "A personal-brand site that is also a working product. The frontend is Next.js 16 (App Router) organized by Feature-Sliced Design, the backend is a small Go 1.25 service, and Postgres 16 stores every piece of content as JSONB with an {en, ru} pair inside — so both languages come from one row and never drift apart. Everything visible on the site is editable from the admin panel: projects with their case-study blocks, work history, education, certificates, publications, services, the homepage sections, and the contact form schema itself. Saving from the admin panel calls a revalidation endpoint, so a static page is rebuilt within seconds instead of waiting for the next deploy or an ISR timer. The Go side also proxies Umami analytics with a singleflight cache, syncs GitHub stars into project cards, rate-limits public endpoints, and signs the admin session. Deploys run from GitHub Actions to GHCR, run goose migrations, then hold a health gate for up to 120 seconds — if the new images do not answer, the workflow reads .last-good-sha and rolls production back on its own.",
    "ru": "Персональный сайт, который одновременно рабочий продукт. Фронт — Next.js 16 (App Router) по Feature-Sliced Design, бэк — небольшой сервис на Go 1.25, Postgres 16 хранит весь контент как JSONB с парой {en, ru} внутри: оба языка берутся из одной строки и не расходятся между собой. Всё, что видно на сайте, редактируется из админки: проекты вместе с блоками кейса, опыт работы, образование, сертификаты, публикации, услуги, секции главной и даже схема формы обратной связи. Сохранение из админки дёргает ручку ревалидации, поэтому статическая страница пересобирается за секунды, а не ждёт следующего деплоя или таймера ISR. Go-часть ещё проксирует аналитику Umami с singleflight-кешем, подтягивает звёзды GitHub в карточки проектов, ограничивает частоту запросов к публичным ручкам и подписывает админскую сессию. Деплой идёт из GitHub Actions в GHCR, прогоняет goose-миграции и держит health-gate до 120 секунд — если новые образы не отвечают, workflow читает .last-good-sha и откатывает прод сам."
  }'::jsonb,
  'web',
  'open_source',
  ARRAY['portfolio','nextjs','golang','postgresql','i18n','seo','docker'],
  ARRAY['Next.js 16','React 19','TypeScript','Tailwind CSS','Go 1.25','PostgreSQL 16','Docker','GitHub Actions','Caddy'],
  '{
    "en": "Build a bilingual personal site I can edit without touching code or waiting for a deploy — and deploy it the same way I would deploy a client project.",
    "ru": "Сделать двуязычный личный сайт, который можно править без кода и без ожидания деплоя — и катить его так же, как я катила бы клиентский проект."
  }'::jsonb,
  'https://github.com/laviercasey/portfolio-site',
  'https://lavier.tech',
  'https://lavier.tech',
  '/images/projects/lavier-tech.svg',
  TRUE,
  5,
  '{
    "en": "Most developer portfolios hardcode their content in JSX. Fixing a typo means an edit, a commit, and a deploy; adding a second language doubles every file and guarantees the translations drift apart within a month. And a personal site is exactly the place where content changes constantly — a new project, a new number in a case study, a new job title.",
    "ru": "Большинство портфолио разработчиков держат контент прямо в JSX. Поправить опечатку — это правка, коммит и деплой; добавить второй язык — удвоить каждый файл и через месяц гарантированно получить разъехавшиеся переводы. А личный сайт — как раз то место, где контент меняется постоянно: новый проект, новая цифра в кейсе, новая должность."
  }'::jsonb,
  '{
    "en": "Content lives in Postgres, not in the codebase. Every translatable field is a JSONB object with en and ru keys, so a project or a job entry is one row that carries both languages. The Go backend serves it, the Next.js frontend renders it statically, and the admin panel writes it back and triggers on-demand revalidation for exactly the paths that changed. The frontend follows Feature-Sliced Design so entities, features, and widgets stay separable. Deployment mirrors the client projects: GitHub Actions builds both images, pushes to GHCR, dumps the database before every deploy, runs goose migrations, and only marks the release good after a health gate passes.",
    "ru": "Контент живёт в Postgres, а не в кодовой базе. Каждое переводимое поле — JSONB-объект с ключами en и ru, так что проект или запись об опыте это одна строка, несущая оба языка. Go-бэкенд их отдаёт, Next.js рендерит статикой, админка пишет обратно и запускает ревалидацию ровно тех путей, которые изменились. Фронт разложен по Feature-Sliced Design, чтобы entities, features и widgets оставались разделимыми. Деплой повторяет клиентские проекты: GitHub Actions собирает оба образа, пушит в GHCR, снимает дамп БД перед каждым выкатом, прогоняет goose-миграции и помечает релиз хорошим только после прохождения health-gate."
  }'::jsonb,
  '{
    "en": "Content edits ship in seconds without a deploy, and both languages stay in sync because they are literally stored side by side. Eight admin sections cover everything the public site renders. The deploy pipeline has a 120-second health gate and rolls itself back to .last-good-sha on failure, so a bad release does not become downtime. The site also runs the SEO plumbing that clients ask for: sitemap, robots, JSON-LD, Open Graph images generated per project, hreflang pairs, and IndexNow pings on publish.",
    "ru": "Правки контента выкатываются за секунды и без деплоя, а оба языка не разъезжаются, потому что физически лежат рядом. Восемь разделов админки покрывают всё, что рендерит публичная часть. В пайплайне деплоя стоит health-gate на 120 секунд и авто-откат к .last-good-sha, поэтому плохой релиз не превращается в простой. Заодно на сайте собрана вся SEO-обвязка, которую просят клиенты: sitemap, robots, JSON-LD, Open Graph-картинки под каждый проект, hreflang-пары и пинг IndexNow при публикации."
  }'::jsonb,
  '[
    {
      "tech": "Postgres JSONB for i18n",
      "reason": {
        "en": "One row per entity carrying {en, ru} beats parallel tables or parallel JSON files. Adding a language is a key, not a migration, and a half-translated field is visible immediately instead of silently falling back.",
        "ru": "Одна строка на сущность с парой {en, ru} лучше параллельных таблиц или параллельных JSON-файлов. Добавить язык — это ключ, а не миграция, а недопереведённое поле видно сразу, а не проваливается молча в фолбэк."
      }
    },
    {
      "tech": "Next.js 16 + push-based ISR",
      "reason": {
        "en": "Static pages with on-demand revalidation give static speed with CMS behavior. The admin panel pushes the exact paths to rebuild, so I never had to pick a revalidate interval and hope it is short enough.",
        "ru": "Статические страницы с ревалидацией по требованию дают скорость статики и поведение CMS. Админка сама сообщает, какие пути пересобрать, так что не приходится подбирать интервал revalidate и надеяться, что он достаточно короткий."
      }
    },
    {
      "tech": "Go 1.25 backend",
      "reason": {
        "en": "A single small binary next to Postgres — fast start, tiny image, and no second Node runtime to keep alive. It also makes the analytics proxy trivial: singleflight collapses concurrent Umami requests into one upstream call.",
        "ru": "Один маленький бинарь рядом с Postgres — быстрый старт, крошечный образ и никакого второго Node-рантайма. Заодно упрощается прокси аналитики: singleflight схлопывает параллельные запросы к Umami в один поход наверх."
      }
    },
    {
      "tech": "Caddy + health-gated deploy",
      "reason": {
        "en": "Caddy handles TLS automatically and fronts several apps on the same VPS. The health gate is the part that matters: a release is only recorded as good after the containers answer, so rollback is a file read instead of a night of debugging.",
        "ru": "Caddy сам держит TLS и фронтит несколько приложений на одном VPS. Важна именно часть с health-gate: релиз считается хорошим только после того, как контейнеры ответили, поэтому откат — это чтение файла, а не ночь отладки."
      }
    }
  ]'::jsonb,
  '[
    { "label": { "en": "Content languages", "ru": "Языка контента" }, "value": "2" },
    { "label": { "en": "Admin sections", "ru": "Разделов в админке" }, "value": "8" },
    { "label": { "en": "Deploy health gate", "ru": "Окно health-gate" }, "value": "120 s" },
    { "label": { "en": "Rollback on failure", "ru": "Откат при падении" }, "value": "auto" }
  ]'::jsonb,
  '2026-04-27',
  NULL
)
ON CONFLICT (slug) DO NOTHING;
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO projects (
  id, slug, title, short_desc, description, category, status, tags, tech_stack,
  goal_desc, github_url, thumbnail, featured, sort_order,
  problem, approach, outcome, tech_choices, highlights,
  timeline_started, timeline_shipped
) VALUES (
  'c7a2e9d4-6b18-4f35-8e70-1d93b5c2a640',
  'metropt-predictive-maintenance',
  '{"en": "MetroPT Predictive Maintenance", "ru": "MetroPT — предиктивное обслуживание"}'::jsonb,
  '{
    "en": "Early detection of air leaks in a metro-train compressor 30 minutes before failure — 1.5M sensor readings, 16 models compared, tuned GradientBoosting at 0.993 ROC-AUC",
    "ru": "Раннее детектирование утечек воздуха в компрессоре поезда метро за 30 минут до отказа — 1,5 млн показаний сенсоров, 16 моделей в сравнении, оптимизированный GradientBoosting с ROC-AUC 0,993"
  }'::jsonb,
  '{
    "en": "Coursework for the Machine Learning and Data Analysis module, built on the open MetroPT-3 dataset from the UCI repository: 1,516,948 time-series points from 15 sensors of a Porto metro-train air compressor, sampled every 10 seconds over seven months, with four recorded air-leak failures. The task is binary classification with a 30-minute lead time — given a 10-minute sliding window of sensor readings, decide whether the compressor is failing or will start failing within the next half hour. Raw readings were aggregated into 50,563 windows with 58 features (six statistics over seven analog sensors plus two over eight digital ones), split strictly by time with no shuffling, and validated with TimeSeriesSplit. Seven classical algorithms and nine ensembles were compared under one protocol, then the winner went through RandomizedSearch, threshold tuning, and feature selection. Unsupervised methods were run alongside as a sanity check: PCA needed 11 components for 95% of variance, and K-Means concentrated 24.1% of positives in one cluster against a 1.34% base rate.",
    "ru": "Курсовая по дисциплине «Машинное обучение и анализ данных» на открытом датасете MetroPT-3 из репозитория UCI: 1 516 948 точек временного ряда с 15 сенсоров воздушного компрессора поезда метро Порту, шаг 10 секунд, семь месяцев наблюдений и четыре зарегистрированных отказа типа Air Leak. Задача — бинарная классификация с горизонтом упреждения 30 минут: по скользящему 10-минутному окну показаний понять, находится ли компрессор в отказе или отказ начнётся в ближайшие полчаса. Сырые показания агрегированы в 50 563 окна с 58 признаками (шесть статистик по семи аналоговым сенсорам плюс две по восьми цифровым), разбиение строго временное без перемешивания, валидация через TimeSeriesSplit. Семь классических алгоритмов и девять ансамблей прогнаны одинаковым протоколом, победитель прошёл RandomizedSearch, подбор порога и отбор признаков. Параллельно как проверка на вменяемость запущены методы без учителя: PCA потребовалось 11 компонент на 95% дисперсии, а K-Means собрал в одном кластере 24,1% позитивов против базовых 1,34%."
  }'::jsonb,
  'research',
  'completed',
  ARRAY['machine-learning','predictive-maintenance','time-series','python','shap'],
  ARRAY['Python','scikit-learn','Gradient Boosting','XGBoost','LightGBM','CatBoost','pandas','SHAP'],
  '{
    "en": "Catch a compressor air leak half an hour before it becomes a failure, on a badly imbalanced time series — and be honest about what the metrics really say.",
    "ru": "Поймать утечку воздуха в компрессоре за полчаса до отказа на сильно несбалансированном временном ряде — и честно разобраться, что на самом деле говорят метрики."
  }'::jsonb,
  'https://github.com/laviercasey/metropt-predictive-maintenance',
  '/images/projects/metropt.svg',
  FALSE,
  8,
  '{
    "en": "Only about 2% of the windows are positive, and neighbouring windows are autocorrelated — so a random split leaks the future into training and every accuracy number becomes meaningless. On top of that, at the default 0.5 threshold the strongest models look useless: RandomForest scores F1 = 0.000 and GradientBoosting F1 = 0.231, which would normally get them thrown out.",
    "ru": "Позитивных окон около 2%, а соседние окна автокоррелированы — при случайном разбиении будущее протекает в обучение, и любая цифра accuracy теряет смысл. Сверху на пороге 0,5 самые сильные модели выглядят бесполезными: у RandomForest F1 = 0,000, у GradientBoosting F1 = 0,231 — обычно после такого модель выбрасывают."
  }'::jsonb,
  '{
    "en": "Time-based split (train up to 1 June 2020, test from that date), TimeSeriesSplit with three folds inside training, and PR-AUC as the primary metric because ROC-AUC flatters imbalanced data. Sixteen models under an identical protocol, then RandomizedSearch plus explicit threshold tuning on the winner — the threshold is a hyperparameter here, not a default. SHAP on the final model to check the reasoning rather than trust the score.",
    "ru": "Временное разбиение (train до 1 июня 2020, test с этой даты), TimeSeriesSplit на три фолда внутри обучения и PR-AUC как основная метрика, потому что ROC-AUC льстит несбалансированным данным. Шестнадцать моделей по одинаковому протоколу, затем RandomizedSearch и явный подбор порога у победителя — порог здесь гиперпараметр, а не значение по умолчанию. SHAP на финальной модели, чтобы проверить логику, а не поверить в цифру."
  }'::jsonb,
  '{
    "en": "The tuned GradientBoosting with a fitted threshold reaches ROC-AUC 0.993, PR-AUC 0.889, F1 0.878 (precision 0.873 / recall 0.882) on the held-out test period — up from F1 0.231 for the same algorithm at the default threshold. The most interesting result is not the score: SHAP put Oil_temperature_min first with a mean |SHAP| of 0.361, while linear correlation had pointed at DV_pressure_mean. Pearson correlation ignores feature interactions and SHAP does not — which is exactly why the feature the correlation matrix liked was not the one the model relied on.",
    "ru": "Оптимизированный GradientBoosting с подобранным порогом даёт на отложенном тестовом периоде ROC-AUC 0,993, PR-AUC 0,889 и F1 0,878 (precision 0,873 / recall 0,882) — против F1 0,231 у того же алгоритма на пороге по умолчанию. Самое интересное не в цифре: SHAP поставил на первое место Oil_temperature_min со средним |SHAP| 0,361, тогда как линейная корреляция указывала на DV_pressure_mean. Корреляция Пирсона не учитывает взаимодействия признаков, а SHAP учитывает — поэтому признак, который нравился корреляционной матрице, оказался не тем, на который опиралась модель."
  }'::jsonb,
  '[
    {
      "tech": "Time-based split + TimeSeriesSplit",
      "reason": {
        "en": "The windows overlap and neighbouring rows are nearly identical. A shuffled split would put the same failure on both sides of the line and return a beautiful, fake score.",
        "ru": "Окна перекрываются, соседние строки почти одинаковые. Перемешанный сплит положил бы один и тот же отказ по обе стороны границы и вернул красивую фальшивую цифру."
      }
    },
    {
      "tech": "Threshold tuning as a step, not a default",
      "reason": {
        "en": "At 0.5 the best model scored F1 0.231 and looked like a failure. Fitting the threshold on validation turned the same model into F1 0.878 — on imbalanced data the cutoff carries as much weight as the algorithm.",
        "ru": "На 0,5 лучшая модель давала F1 0,231 и выглядела провалом. Подбор порога на валидации превратил ту же модель в F1 0,878 — на несбалансированных данных отсечка весит не меньше алгоритма."
      }
    },
    {
      "tech": "PR-AUC as the primary metric",
      "reason": {
        "en": "With ~2% positives, ROC-AUC stays high even for a model that misses most failures. PR-AUC is the one that moves when recall on the rare class is actually bad.",
        "ru": "При ~2% позитивов ROC-AUC остаётся высоким даже у модели, которая пропускает большинство отказов. Двигается именно PR-AUC, когда recall по редкому классу действительно плохой."
      }
    },
    {
      "tech": "SHAP over correlation",
      "reason": {
        "en": "The correlation matrix and the model disagreed about which sensor mattered. SHAP explains the model that will actually run, so it wins that argument.",
        "ru": "Корреляционная матрица и модель разошлись в том, какой сенсор важен. SHAP объясняет ту модель, которая реально поедет в работу, поэтому спор выигрывает он."
      }
    }
  ]'::jsonb,
  '[
    { "label": { "en": "ROC-AUC (test)", "ru": "ROC-AUC (тест)" }, "value": "0.993" },
    { "label": { "en": "PR-AUC (test)", "ru": "PR-AUC (тест)" }, "value": "0.889" },
    { "label": { "en": "Sensor readings", "ru": "Показаний сенсоров" }, "value": "1.5M" },
    { "label": { "en": "Lead time", "ru": "Горизонт упреждения" }, "value": "30 min" }
  ]'::jsonb,
  '2026-05-01',
  '2026-05-03'
)
ON CONFLICT (slug) DO NOTHING;
-- +goose StatementEnd

UPDATE projects SET sort_order = 1,  featured = TRUE  WHERE slug = 'tickethub';
UPDATE projects SET sort_order = 2,  featured = TRUE  WHERE slug = 'cosmos-explorer';
UPDATE projects SET sort_order = 3,  featured = TRUE  WHERE slug = 'ecoshop';
UPDATE projects SET sort_order = 4,  featured = TRUE  WHERE slug = 'med-reminder-bot';
UPDATE projects SET sort_order = 5,  featured = TRUE  WHERE slug = 'lavier-tech';
UPDATE projects SET sort_order = 6,  featured = TRUE  WHERE slug = 'fluxrouter';
UPDATE projects SET sort_order = 7,  featured = FALSE WHERE slug = 'heart-disease-ml-benchmark';
UPDATE projects SET sort_order = 8,  featured = FALSE WHERE slug = 'metropt-predictive-maintenance';
UPDATE projects SET sort_order = 9,  featured = FALSE WHERE slug = 'genreneuro';
UPDATE projects SET sort_order = 10, featured = FALSE WHERE slug = 'children-book-dataset';
UPDATE projects SET sort_order = 11, featured = FALSE WHERE slug = 'diafilm-parts-finder';

-- +goose Down
DELETE FROM projects WHERE slug IN ('lavier-tech', 'metropt-predictive-maintenance');
DELETE FROM publications WHERE id = 'd3e8b1a7-2c94-4f60-a5d8-7b1e6c40f923';
UPDATE certificates SET issue_date = '2024-09' WHERE id = '893315c2-6271-4441-bc3e-6dbb5bd39888';

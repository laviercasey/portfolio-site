# UTM Cheatsheet / Шпаргалка по UTM

## What is UTM? / Что такое UTM?

**EN:** UTM tags are query-string parameters (`utm_source`, `utm_medium`,
`utm_campaign`) appended to a link so analytics can attribute each visit to
the channel that sent it. Umami records them automatically with every pageview.

**RU:** UTM-метки — это query-параметры (`utm_source`, `utm_medium`,
`utm_campaign`), которые добавляются к ссылке, чтобы аналитика понимала,
откуда пришёл посетитель. Umami фиксирует их автоматически при каждом
просмотре страницы.

## Standard values

Use only these canonical `utm_source` / `utm_medium` pairs — keeps the admin
widget readable and prevents fragmentation (`hh` vs `HH.ru` vs `headhunter`).

| Channel | `utm_source` | `utm_medium` |
| --- | --- | --- |
| HH.ru resume profile | `hh` | `resume` |
| LinkedIn profile | `linkedin` | `profile` |
| LinkedIn post | `linkedin` | `social` |
| Telegram channel | `telegram` | `social` |
| Telegram direct message | `telegram` | `direct` |
| Instagram bio | `instagram` | `bio` |
| Instagram stories | `instagram` | `social` |
| Habr article | `habr` | `article` |
| GitHub README | `github` | `readme` |
| Email signature | `email` | `signature` |
| Friend referral | `friend` | `referral` |

## Example link

```
https://lavier.tech/ru/projects/ml-case?utm_source=hh&utm_medium=resume&utm_campaign=2026_q2_search
```

## `utm_campaign` naming convention

Format: `{year}_{quarter}_{theme}` — lowercase, underscores, no spaces.

Examples: `2026_q2_search`, `2026_launch_newsite`, `portfolio_refresh`.

## Cleanup behaviour

UTM params stay visible in the address bar for ~1.5 seconds after page load,
then the `UtmCleanup` client provider strips them via `history.replaceState`
(no navigation event, no extra pageview). Umami has already recorded the
pageview with the UTMs by then — no extra work on your side.

UTM-параметры видны в адресной строке ~1.5 секунды после загрузки, затем
клиентский провайдер `UtmCleanup` удаляет их через `history.replaceState`
(без перехода, без повторного просмотра). К этому моменту Umami уже
зафиксировал визит с UTM-метками — никаких дополнительных действий не нужно.

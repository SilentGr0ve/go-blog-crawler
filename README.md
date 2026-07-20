# go-blog-crawler

> Topic Specific CLI Web Crawler для Go-блогов, ищущий релевантные статьи

Учебный проект, посвященный изучению возможностей работы с конкурентной средой на языке Go.
Особое внимание уделено работе с паттернами конкурентности.

## Возможности

**Реализованные:**
- Настраиваемое структурированное логирование в stderr и файл (text/json)
- CLI с persistent-флагами, наследуемыми подкомандами

**Запланированные:**
- Настройка глубины обхода и лимита страниц
- robots.txt и per-host rate limit для вежливого обхода
- Экспорт результатов в JSON
- Состояние в Redis для продолжения обхода после перезапуска

## Технологии

- **Go** 1.26 — основной язык
- **cobra** — CLI с подкомандами и persistent-флагами
- **log/slog** — структурированное логирование
- **Redis** — внешнее хранилище состояния
- **Docker / docker-compose** — воспроизводимое окружение

## Быстрый старт

```bash
git clone https://github.com/SilentGr0ve/go-blog-crawler
cd go-blog-crawler
make build 
./bin/crawler --help # показать список команд
```

## Команды

CLI построен на cobra. Корневая команда `crawler` содержит
persistent-флаги (`--log-level`, `--log-dir`), которые наследуются
всеми подкомандами.

| Команда | Описание | Статус        |
|---|---|---------------|
| `crawl` | Обход от seed-списка, поиск релевантных страниц | в разработке  |
| `resume` | Продолжить обход из сохранённого в Redis состояния | запланировано |
| `stats` | Метрики текущего или последнего обхода | запланировано |
| `export` | Экспорт результатов в JSON или CSV | запланировано |

Пример запуска (после реализации `crawl`):

```bash
./crawler crawl \
--seeds seeds.txt \
--keywords keywords.txt \
--depth 2 \
--max-pages 100 \
--concurrency 5
```

## Архитектура

**Будет готово позже**

## Что демонстрирует проект

- **Конкуренция**: goroutines, channels, worker pool, WaitGroup, Mutex
- **Логирование**: slog с настраиваемыми handler'ами
- **CLI**: cobra с persistent/local флагами
- **Средства обхода**: robots.txt, rate limit
- **Дисциплина Git**: Conventional Commits, feature-ветки
- **Внешнее состояние**: Redis

## Roadmap

- [x] Шаг 0: скелет (cobra + slog + структура)
- [ ] Шаг 1: последовательный обход
- [ ] Шаг 2: worker pool 
- [ ] Шаг 3: Redis-хранилище
- [ ] Шаг 4: robots.txt + rate limit
- [ ] Шаг 5: релевантность по ключевым словам
- [ ] Шаг 6: подкоманды resume/stats/export
- [ ] Шаг 7: graceful shutdown
- [ ] Шаг 8: финализация (Dockerfile, CI, тесты)

## Git-дисциплина

Проект использует [Conventional Commits](https://www.conventionalcommits.org/):
каждый коммит имеет тип, scope и subject. Крупные фичи — в отдельных
ветках `feat/<название>` с PR в `main`.

**Примеры:**

\`\`\`
feat(crawler): add worker pool with configurable concurrency
fix(fetcher): handle context cancellation during HTTP request
docs(readme): document CLI subcommands
\`\`\`
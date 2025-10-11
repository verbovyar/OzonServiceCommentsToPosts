# OzonProject — GraphQL сервис постов и комментариев

Учебный сервис на **Go**, реализующий систему постов и комментариев в стиле Reddit / Habr.  
Использует **GraphQL (gqlgen)** и **PostgreSQL** с возможностью работы в **in-memory** режиме.  
Готов к запуску через **Docker Compose**.

---

## Возможности

- Просмотр списка постов с пагинацией (`limit`, `offset`)
- Просмотр поста с комментариями
- Возможность отключить комментарии к посту
- Иерархические комментарии (вложенность без ограничений)
- Пагинация комментариев
- Поддержка GraphQL Subscriptions (асинхронная доставка новых комментариев)

---

## Технологии

- **Go 1.22+**
- **gqlgen** — GraphQL-сервер
- **PostgreSQL 15**
- **pgx / pgxpool** — работа с базой
- **pgxmock** — unit-тесты
- **Docker + docker-compose** — сборка и запуск

---

## Запуск проекта

### Через Docker Compose

```bash
docker-compose up --build
```

## Примеры запросов

### Создать пост

```gql
mutation {
  createPost(title: "Hello", content: "My first post", author: "Alice", commentsEnabled: true) {
    id
    title
    author
  }
}
```

### Добавить комментарий

```gql
mutation {
  createComment(postID: "1", author: "Bob", content: "Nice post!") {
    id
    content
    author
  }
}
```

### Подписка на новые комментарии

```gql
subscription {
  commentAdded(postID: "1") {
    id
    author
    content
  }
}
```

### Replay к комментарию

```gql
mutation {
  createComment(
    postID: "1"
    parentID: "123"
    author: "Bob"
    content: "Отвечаю на комментарий 123"
  ) {
    id
    postId
    parentId
    author
    content
    createdAt
  }
}
```

### Пост со всеми комментами 

```gql
query {
  post(id: "1") {
    id
    title
    author
    createdAt
    commentsEnabled

    # Верхний уровень
    comments(limit: 50, offset: 0) {
      id
      author
      content
      createdAt
      parentId

      # 1-й уровень вложенности
      children(limit: 50, offset: 0) {
        id
        author
        content
        createdAt
        parentId

        # 2-й уровень вложенности
        children(limit: 50, offset: 0) {
          id
          author
          content
          createdAt
          parentId

          # 3-й уровень вложенности
          children(limit: 50, offset: 0) {
            id
            author
            content
            createdAt
            parentId
          }
        }
      }
    }
  }
}
```

## Структура проекта

```pgsql
.
├── cmd/
│   └── service/              # Точка входа
├── graph/                    # GraphQL схема и резолверы (gqlgen)
├── config/                   # Конфиг файл
├── internal/
│   ├── models/               # Модели данных
│   ├── storage/              # Хранилище на PostgreSQL и in memory
│   ├── service/              # Бизнес-логика
|   ├── validation/           # Валидация
|   ├── utils/                # Утилиты
|   ├── pubsub/               # (Subscribe/Unsubscribe/Publish)
├── migrations/               # SQL миграции
├── pkg/
├── docker-compose.yml
├── Dockerfile
└── README.md
```
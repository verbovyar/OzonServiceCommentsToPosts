**REST API для блога с системой комментариев** — высокопроизводительное Go-приложение с in-memory хранилищем и поддержкой древовидных комментариев.

Создать пост
mutation {
  createPost(title: "Hello", content: "First Post", author: "Yaroslav", commentsEnabled: true) {
    id title commentsEnabled
  }
}

Список постов
query {
  posts(limit: 10, offset: 0) { id title author createdAt }
}

Получить пост с комментами верхнего уровня
query {
  post(id: "1") {
    id title
    comments(limit: 20, offset: 0) {
      id author content createdAt
    }
  }
}

Добавить коммент
mutation {
  createComment(postId: "1", author: "Sergey", content: "Nice post") {
    id postId content
  }
}

Дочерние комменты конерктного узла
query {
  post(id: "1") {
    comments(limit: 10, offset: 0) {
      id
      children(limit: 10, offset: 0) { id content }
    }
  }
}
------------------------------------------------------------------------------
Добавить коммент верхнего уровня
mutation {
  createComment(postId: "1", author: "Elena", content: "My son") {
    id
  }
}
Ответ на этот коммент
mutation {
  createComment(postId: "1", parentId: "1", author: "Yaroslav", content: "Reply to Elena") {
    id
  }
}
Еще уровень вложенности 
mutation {
  createComment(postId: "1", parentId: "2", author: "Alina", content: "Nested reply") {
    id
  }
}

Читатем теперь
Получить верхний уровень под постом
query {
  post(id: "1") {
    comments(limit: 20, offset: 0) {
      id author content createdAt
      children(limit: 10, offset: 0) {
        id author content
      }
    }
  }
}
Подгружать детей по клику
query {
  post(id: "1") {
    comments(parentId: "1", limit: 10, offset: 0) {
      id content
    }
  }
}
---------------------------
Подписка на пост
subscription {
  commentAdded(postId: "1") {
    id
    postId
    parentId
    author
    content
    createdAt
  }
}

В другой вкладке 
mutation {
  createComment(postId: "1", author: "Roma", content: "Hello via subscriptions!") {
    id
  }
}

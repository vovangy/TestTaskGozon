# TestTaskGozon

## Для запуска приложения понадобится

- Docker и Docker Compose: используется версия "3.9" для docker-compose.yml.
- Находясь в директории с docker-compose.yml введите команду ниже. Она запустит приложение
```sh
make run
``` 
- Команда ниже выключает приложение
```sh
make stop
``` 

## Не реализовано
1. Хранилище in memmory, только postgresql
2. Пагинация комментариев
3. Тесты написаны только на уровне репозитория


## API Сервиса
Приложение запускается на localhost:8080
# Регистрация пользователя
1. username - Логин - String
2. password - Пароль - String
```graphql
mutation Registry {
  signUp(input: {
    username: "Vova",
    password: "123456"
  }) {
    authToken {
      accessToken
      expiredAt
    }
  }
}
```
- Ответом будет являться
```graphql
{
  "data": {
    "signUp": {
      "authToken": {
        "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTY4Mjg0MjgsImlkIjozfQ.vC0cKH-qTHL3-ouaUHb9d2WqWHflhCionTDArfO19-4",
        "expiredAt": "2024-05-27 16:47:08"
      }
    }
  }
}
```
- Откуда можно взять accessToken и проставить у себя в http headers для мутаций требующих авторизацию, а именно
```graphql
{
  "Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTY4Mjg0MjgsImlkIjozfQ.vC0cKH-qTHL3-ouaUHb9d2WqWHflhCionTDArfO19-4"
}
```

# Авторизация пользователя
1. username - Логин - String
2. password - Пароль - String
```graphql
mutation Login {
  signIn(input: {
    username: "Vova",
    password: "123456"
  }) {
    authToken {
      accessToken
      expiredAt
    }
  }
}
```
- Ответом будет являться
```graphql
{
  "data": {
    "signUp": {
      "authToken": {
        "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTY4Mjg0MjgsImlkIjozfQ.vC0cKH-qTHL3-ouaUHb9d2WqWHflhCionTDArfO19-4",
        "expiredAt": "2024-05-27 16:47:08"
      }
    }
  }
}
```
- Откуда можно взять accessToken и проставить у себя в http headers для мутаций требующих авторизацию, а именно
```graphql
{
  "Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTY4Mjg0MjgsImlkIjozfQ.vC0cKH-qTHL3-ouaUHb9d2WqWHflhCionTDArfO19-4"
}
```
# Создание поста !Требуется авторизация
1. title - Название поста - String
2. content - Содержимое поста - String
3. IsCommented - Можно ли комментировать пост - Boolean

```graphql
mutation PostCreate {
  createPost(input: {
    title: "alloo",
    content: "Baaaaaa",
    IsCommented: false
  }
  ) {
     title
     id
     content
    author {
      id
      username
    }
    }
  }
```
- Ответом будет являться
```graphql
{
  "data": {
    "createPost": {
      "title": "alloo",
      "id": "14",
      "content": "Baaaaaa",
      "author": {
        "id": "1",
        "username": "AVova"
      }
    }
  }
}
```
# Создание комментария !Требуется авторизация
1. postId - Id поста - String
2. content - Содержимое комментария - String
3. parentCommentId - Id К какому комментарию принадлежит, если комментарий принадлежит к посту, parentCommentId указать в 0 - String
```graphql
mutation CommentCreate {
  createComment(input: {
    postId: 4,
    content: "Baaaaaa",
    parentCommentId: 0
  }
  ) {
     id
     content
    author {
      id
      username
    }
    }
  }
```
- Ответом будет являться
```graphql
{
  "data": {
    "createComment": {
      "id": "38",
      "content": "Baaaaaa",
      "author": {
        "id": "1",
        "username": "AVova"
      }
    }
  }
}
```
# Блокировка комментариев к посту !Требуется авторизация
1. postId - Id поста - String
```graphql
mutation {
  blockComments(postId: "5")
}
```
- Ответом будет являться
```graphql
{
  "data": {
    "blockComments": "Comments blocked"
  }
}
```
# Получить список постов
```graphql
query {
  posts {
    title
    content
    comments {
      content
      }
    }
}
```
# Получить информацию о посте
1. postId - Id поста - String
```graphql
query {
  post(id :4) {
    id
    title
    content
    author {
      username
    }
    comments {
      id
      replies {
        id
      }
    }
  }
}
```
# Ограничения на вложенность у комментариев нет, то есть вложенность replies может быть любой
```graphql
query {
  post(id :4) {
    id
    title
    content
    author {
      username
    }
    comments {
      id
      replies {
        id
        replies {
            id
            replies {
                id
                replies {
                    id
                }
            }
        }
      }
    }
  }
}
```

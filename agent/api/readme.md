# API агента мониторинга

## Обзор

API агента мониторинга предоставляет доступ к метрикам системы и позволяет настраивать его работу.

## Создание документации

Для создания документации по комментариям:

- установить библиотеку swag: `go install github.com/swaggo/swag/cmd/swag@latest`
- выполнить команду генерации из папки agent: `swag init -g ./internal/transport/handler.go -o ./api`

## Документация

Полная спецификация API доступна в формате OpenAPI (Swagger) в файле `swagger.yaml`.

Для просмотра документации в удобном виде:

1. Используйте [Swagger Editor](https://editor.swagger.io/)
2. Скопируйте содержимое `swagger.yaml` в редактор
3. Или запустите Swagger UI локально:

```bash
docker run -p 8090:8080 -e SWAGGER_JSON=/api/swagger.yaml -v $(pwd)/api:/api swaggerapi/swagger-ui
```

После успешного запуска Swagger UI будет доступен по адресу http://localhost:8090.
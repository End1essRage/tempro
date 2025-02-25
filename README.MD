# Tempro - Project Templating Tool

Утилита для быстрого создания проектов из шаблонов с поддержкой множества языков и компонентов.

## комманды
tempro <lang> <command> [params]
| Команда | Описание |
|:--------|:---------|
|list|список темлпейтов и шардов|
|init|создает проект|
|apply|применяет шард|
|revert|отменяет последнюю команду|
|env gen|создает файл .env с необходимыми переменными|
|env add|Добавляет переменную|

## 📦 Установка
```bash
go install github.com/end1essrage/tempro@latest

# Шаблоны / ассеты
tempro list 
tempro go list

# Создание проекта
tempro go init -n myModule	Создать новый проект

### Репозиторий темлпейтов
github.com/end1essrage/tempro-templates

## Файл конфигурации
~/.tempro/config.yaml

sources:
  - name: "official"
    url: "https://github.com/tempro-templates"

  - name: "custom"
    url: "git@github.com:my-company/templates.git"
    auth: ssh-key

#Добавить источник
tempro source add my-templates git@github.com:user/my-templates.git

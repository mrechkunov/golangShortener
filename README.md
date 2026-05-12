# go-musthave-shortener-tpl

Шаблон репозитория для трека «Сервис сокращения URL».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m v2 template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/v2 .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

## Структура проекта

Приведённая в этом репозитории структура проекта является рекомендуемой, но не обязательной.

Это лишь пример организации кода, который поможет вам в реализации сервиса.

При необходимости можно вносить изменения в структуру проекта, использовать любые библиотеки и предпочитаемые структурные паттерны организации кода приложения, например:
- **DDD** (Domain-Driven Design)
- **Clean Architecture**
- **Hexagonal Architecture**
- **Layered Architecture**
File: shortener
Type: inuse_space
Time: 2026-05-12 15:16:44 MSK
Showing nodes accounting for 519.85kB, 50.71% of 1025.05kB total
      flat  flat%   sum%        cum   cum%
  519.03kB 50.63% 50.63%   519.03kB 50.63%  runtime.mallocgc
  512.88kB 50.03% 100.67%   512.88kB 50.03%  regexp/syntax.map.init.1
 -512.05kB 49.95% 50.71%  -512.05kB 49.95%  github.com/mrechkunov/golangShortener.git/internal/service.SetIsDeleted
         0     0% 50.71%   512.88kB 50.03%  regexp/syntax.init
         0     0% 50.71%   512.88kB 50.03%  runtime.doInit (inline)
         0     0% 50.71%   512.88kB 50.03%  runtime.doInit1
         0     0% 50.71%   512.88kB 50.03%  runtime.main
         0     0% 50.71%   519.03kB 50.63%  runtime.newobject
         0     0% 50.71%   519.03kB 50.63%  runtime.procresize
         0     0% 50.71%   519.03kB 50.63%  runtime.rt0_go
         0     0% 50.71%   519.03kB 50.63%  runtime.schedinit

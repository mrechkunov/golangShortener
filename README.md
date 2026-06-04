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

Анализ следующий:
Запускаем сервер, даем нагрузку клиентом в это время снимаем профиль памяти ./profiles/base1.proff
после анализа виим, что encoding/hex.EncodeToString находится в топе. идем в handler и результаты выполнения энкодинга записываем в переменную и далее используем ее по коду (в оригинали энкодинг выполнялся дважды при формировании 
ShortURL := baseResultAdress + "/" + hex.EncodeToString(hash[:4])
и при записи в хранилище
err = repository.GetStorage().SetData(hex.EncodeToString(hash[:4]), req.URL, cookie.Value))
после повторяем тестирование и снимаем профиль памяти ./profiles/result1.proff
далее делаем go tool pproff -diff_base, результаты ниже 

File: shortener
Type: inuse_space
Time: 2026-05-15 10:36:42 MSK
Entering interactive mode (type "help" for commands,  "o" for options)
(pprof) top
Showing nodes accounting for -7562.82kB, 25.74% of 29376.28kB total
Dropped 2 nodes (cum <= 146.88kB)
Showing top 10 nodes out of 45
      flat  flat%   sum%        cum   cum%
-3584.27kB 12.20% 12.20% -3584.27kB 12.20%  encoding/hex.EncodeToString
-3343.56kB 11.38% 23.58% -3343.56kB 11.38%  github.com/mrechkunov/golangShortener.git/internal/repository.(*SafeMap).SetData
-1025.01kB  3.49% 27.07% -1025.01kB  3.49%  runtime.mallocgc
  902.59kB  3.07% 24.00%   902.59kB  3.07%  compress/flate.NewWriter (inline)
 -512.56kB  1.74% 25.74%  -512.56kB  1.74%  sync.(*Pool).pinSlow
 -512.02kB  1.74% 27.49%  -512.02kB  1.74%  encoding/json.(*decodeState).literalStore
  512.01kB  1.74% 25.74%   512.01kB  1.74%  encoding/json.(*scanner).pushParseState
         0     0% 25.74%   902.59kB  3.07%  compress/gzip.(*Writer).Write
         0     0% 25.74%   902.59kB  3.07%  encoding/json.(*Encoder).Encode
         0     0% 25.74%  -512.02kB  1.74%  encoding/json.(*decodeState).object
(pprof) 

после чего запустили сервер в режиме работы с БД и сняли профиль памяти ./profiles/result2.pprof
diff_base ниже

File: shortener
Type: inuse_space
Time: 2026-05-15 10:37:50 MSK
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for -26423.14kB, 89.95% of 29376.28kB total
Showing top 10 nodes out of 55
      flat  flat%   sum%        cum   cum%
-12259.71kB 41.73% 41.73% -12259.71kB 41.73%  github.com/mrechkunov/golangShortener.git/internal/repository.(*SafeMap).SetData
-8192.59kB 27.89% 69.62% -8192.59kB 27.89%  encoding/hex.EncodeToString
-2560.08kB  8.71% 78.34% -2560.08kB  8.71%  encoding/json.(*decodeState).literalStore
-1805.17kB  6.14% 84.48% -2898.68kB  9.87%  compress/flate.NewWriter
-1093.51kB  3.72% 88.20% -1093.51kB  3.72%  compress/flate.(*compressor).initDeflate (inline)
 -512.56kB  1.74% 89.95%  -512.56kB  1.74%  sync.(*Pool).pinSlow
  512.50kB  1.74% 88.20%   512.50kB  1.74%  go.uber.org/zap/internal/bufferpool.init.NewPool.func1
  512.05kB  1.74% 86.46%   512.05kB  1.74%  github.com/golang-migrate/migrate/v4.(*Migrate).lock.func2
 -512.05kB  1.74% 88.20%  -512.05kB  1.74%  github.com/mrechkunov/golangShortener.git/internal/service.SetIsDeleted
 -512.01kB  1.74% 89.95%  -512.01kB  1.74%  runtime.mallocgc
(pprof) 

выводы: тяжелые функции лучше всего делать один раз и сохранять результаты в переменную
работа с БД занимает меньше памяти на стеке.


Для установки значений переменным 
var buildVersion string = "N/A"
var buildDate string = "N/A"
var buildCommit string = "N/A"
при компиляции, необходимо добавить флаг
 -ldflags "-X main.<var>=<value>"
 например:
 go build -ldflags "-X main.buildVersion=v1.0.1" ./cmd/shortener/.
 
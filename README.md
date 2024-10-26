# SSO

Сервис на Go, реализующий авторизацию и аутентификацию пользователей.

Прото-файлы можно найти в репозитории: [sso_proto](https://github.com/sariya23/sso_proto).

## Usage 📖

Сервис предоставляет три RPC-метода:

- `Register`
- `Login`
- `IsAdmin`

Сигнатуры методов описаны в прото-файлах: [sso_proto](https://github.com/sariya23/sso_proto).

Реализация клиента может отличаться в зависимости от используемого языка.

### Python 🐍

1. Склонируйте репозиторий с прото-файлами в ваш проект:
    
    ```shell
    git clone git@github.com:sariya23/sso_proto.git
    ```

2. Установите зависимости (внутри или вне виртуального окружения):

    ```shell
    pip install grpcio grpcio-tools
    ```

3. Сгенерируйте Python-классы на основе прото-файлов:
    ```shell
    python3 -m grpc_tools.protoc -I ./sso_proto/proto/sso --python_out=. --grpc_python_out=. ./sso_proto/proto/sso/sso.proto
    ```
    - `-I` — путь к директории с прото-файлами
    - `--python_out` — директория для сохранения сгенерированных классов
    - `--grpc_python_out` — директория для сохранения кода gRPC
    - `./sso_proto/proto/sso/sso.proto` — путь к прото-файлу

4. Создайте клиент:
    ```py
    import grpc

    from sso_pb2_grpc import AuthStub
    import sso_pb2 as sso_pb

    def run():
        with grpc.insecure_channel("localhost:44044") as ch:
            client = AuthStub(ch)
            req = sso_pb.RegisterRequest(
                email="email@yandex.ru",
                password="strongpassword",
            )
            response = client.Register(req)
            print(response)

    if __name__ == "__main__":
        run()
    ```

### Go 🚀

```go
import (
    "context"
    "fmt"
    ssov1 "github.com/sariya23/sso_proto/gen/sso"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    ctx := context.Background()
    conn, err := grpc.Dial("localhost:44044", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        panic(err)
    }

    client := ssov1.NewAuthClient(conn)
    response, err := client.Register(ctx, &ssov1.RegisterRequest{Email: "email@yandex.ru", Password: "password"})
    if err != nil {
        panic(err)
    }
    fmt.Printf("%+v\n", response) // user_id: <int64>
}
```

## Локальный запуск 🖥️
Настройте локальный конфиг. Конфигурация базы данных находится в файле (нужно создать) `./config/db.yaml`, пример - `./config/db.example.yaml`.

Создайте файл `./.env` для хранения конфигурации БД для Docker. Пример — `./.env.example`.

> **Важно! Параметры в файлах .env и .yaml должны совпадать.**

### Docker 📦
Выполните команду в корне проекта:
```shell
docker-compose --env-file=.env build 
docker-compose --env-file=.env up -d
```
Команда создаст БД с параметрами из `.env`, применит миграции из `./migrations` и запустит приложение на `localhost:44044`.

### Запуск без Docker 🏗️
1. Создайте вручную БД PostgreSQL и укажите её имя в db_name внутри `.yaml`.

2. Примените миграции, выполнив команду:
    ```shell
    make migrate
    ```

3. Запустите приложение:
    ```shell
    make run
    ```
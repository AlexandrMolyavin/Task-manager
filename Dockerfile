# Используйте официальный образ golang для первой стадии сборки
#FROM golang:alpine

# Установите рабочую директорию внутри контейнера
#WORKDIR /app

# Копируйте исходный код приложения в контейнер
#COPY . .

# Скачайте модули Go
#RUN go mod download

# Соберите приложение
#RUN go build -o ./myapi


# Откройте порт, на котором будет слушать приложение
#EXPOSE 8080

# Команда для запуска приложения
#CMD ["/app/myapi"]
FROM golang:alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

WORKDIR /app

RUN go build -o myapp .

WORKDIR /app

EXPOSE 8080

CMD ["./myapp"]
FROM golang:latest
LABEL authors="ilia"

WORKDIR /DuckieDatasetServer

COPY . /DuckieDatasetServer

CMD ["go", "run", "main.go", "-o", "/Data"]

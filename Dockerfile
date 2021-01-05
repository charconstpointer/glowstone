
FROM golang:1.15-alpine AS build

WORKDIR /src/
COPY . /src/
ENV GO111MODULE=on
RUN CGO_ENABLED=0 go build -o /bin/tunnel /src/cmd/tunnel/main.go

FROM scratch
COPY --from=build /bin/tunnel /bin/tunnel
ENTRYPOINT ["/bin/tunnel"]
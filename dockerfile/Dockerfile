FROM golang:1.17.1 as build
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN mkdir bin && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-w -s" -gcflags "-N -l" -o bin/http_stadys

FROM alpine
COPY --from=build /app/bin/* .
CMD ["./http_stadys"]
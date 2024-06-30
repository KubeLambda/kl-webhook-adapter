FROM golang:alpine AS build

RUN apk --no-cache add \
    gcc \
    g++ \
    make \
    git


WORKDIR /go/src/app

COPY ./cmd ./cmd
COPY ./configs ./configs
COPY ./internal ./internal


COPY go.mod .
COPY go.sum .
COPY *.go ./
RUN go mod tidy
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/webhook-adapter ./main.go


FROM alpine:3.17
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /go/src/app/bin /app
COPY --from=build /go/src/app/configs /app/configs


EXPOSE 8080

CMD /app/webhook-adapter run --deployment=production

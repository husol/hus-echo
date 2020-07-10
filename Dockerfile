FROM golang:1.14.4-alpine AS go_builder
ADD . /project
WORKDIR /project
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags "-s -w" -o ./bin/server ./cmd/main.go

FROM alpine:3.10
RUN apk --no-cache add ca-certificates tzdata bash && \
    cp /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime
COPY --from=go_builder /project/bin/server ./server
COPY --from=go_builder /project/scripts/wait-for-it.sh ./scripts/wait-for-it.sh
COPY --from=go_builder /project/migrations ./scripts/wait-for-it.sh

EXPOSE 3000
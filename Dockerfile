FROM golang:alpine AS build

ARG app=yarb-ig-producer

RUN apk add --update --no-cache tzdata

WORKDIR /app
COPY . .
RUN go mod init $app && go mod tidy && CGO_ENABLED=0 go build -ldflags "-s -w"

FROM scratch

ENV TZ Europe/Moscow

COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /app/$app /app/$app

USER 1000:1000

ENTRYPOINT ["/app/$app"]

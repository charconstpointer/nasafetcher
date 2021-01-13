FROM golang:1.15-alpine AS build
RUN apk --no-cache add ca-certificates
WORKDIR /src/
COPY . /src/
ENV GO111MODULE=on
RUN CGO_ENABLED=0 go build -o /bin/url-collector /src/cmd/url-collector/main.go

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/url-collector /bin/url-collector
ENTRYPOINT ["/bin/url-collector"]
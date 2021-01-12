FROM golang:1.15-alpine AS build

WORKDIR /src/
COPY . /src/
ENV GO111MODULE=on
RUN CGO_ENABLED=0 go build -o /bin/url-collector /src/cmd/url-collector/main.go

FROM scratch
COPY --from=build /bin/url-collector /bin/url-collector
ENTRYPOINT ["/bin/url-collector"]
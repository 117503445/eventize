FROM golang:1.22 as build
WORKDIR /workspace
RUN go env -w GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o app ./cmd/eventize

FROM gcr.io/distroless/static-debian12 as prod
WORKDIR /workspace
COPY --from=build /workspace/app app
ENTRYPOINT [ "./app"]
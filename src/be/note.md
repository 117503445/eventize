go install github.com/spf13/cobra-cli@latest
go run cmd/eventize/main.go
go build ./cmd/eventize

go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2


zxh404.vscode-proto3
pacman -Sy protobuf

pacman -Sy --noconfirm websocat
websocat ws://127.0.0.1:9090/ws

go run -mod=mod entgo.io/ent/cmd/ent new Event
go generate ./ent
FROM 117503445/dev-golang

RUN go install github.com/twitchtv/twirp/protoc-gen-twirp@latest && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && pacman -Syu --noconfirm protobuf npm yarn rsync && yarn config set registry https://registry.npmmirror.com/ && yarn global add twirpscript
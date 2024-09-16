package main

import (
	"github.com/117503445/goutils"
	"github.com/rs/zerolog/log"
)

func main() {
	goutils.InitZeroLog()

	log.Debug().Msg("debug")
	goutils.Exec("docker build -t 117503445/eventize-dev .")

	goutils.Exec("docker run --rm -v /root/workspace/eventize:/workspace --entrypoint /workspace/src/dev/build_proto.sh 117503445/eventize-dev")

}

package main

import (
	"fmt"

	"github.com/117503445/goutils"
	"github.com/rs/zerolog/log"
)

func main() {
	goutils.InitZeroLog()
	goutils.ExecOpt.DumpOutput = true

	log.Debug().Msg("debug")
	goutils.Exec("docker build -t 117503445/eventize-dev .")

	rootDir, err := goutils.FindGitRepoRoot()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get git root dir")
	}

	// docker run -it --rm -v $PWD:/workspace --entrypoint fish 117503445/eventize-dev
	cmd := fmt.Sprintf("docker run --rm -v %v:/workspace --entrypoint /workspace/src/dev/build_proto.sh 117503445/eventize-dev", rootDir)
	goutils.Exec(cmd)

}

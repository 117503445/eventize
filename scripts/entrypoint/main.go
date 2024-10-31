package main

import (
	"github.com/117503445/goutils"
	"github.com/rs/zerolog/log"
)

func main() {
	goutils.InitZeroLog()

	goutils.ExecOpt.DumpOutput = true

	log.Debug().Msg("Entrypoint")

	goutils.Exec("dockerd", goutils.WithDumpOutput{})
}

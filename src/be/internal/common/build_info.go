package common

import (
	_ "embed"
	"encoding/json"

	"github.com/rs/zerolog/log"
)

//go:embed build_info.json
var buildInfoText string

type BuildInfo struct {
	Version string `json:"version"`
	Tag     string `json:"tag"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
	Dirty   bool   `json:"dirty"`
}

var buildInfo *BuildInfo

func init() {
	buildInfo = &BuildInfo{}
	if err := json.Unmarshal([]byte(buildInfoText), buildInfo); err != nil {
		log.Fatal().Err(err).Msg("failed to unmarshal build info")
	}
}

func GetBuildInfo() *BuildInfo {
	return buildInfo
}

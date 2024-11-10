package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/117503445/goutils"
	"github.com/rs/zerolog/log"
)

func main() {
	goutils.InitZeroLog()
	goutils.ExecOpt.DumpOutput = true

	rootDir, err := goutils.FindGitRepoRoot()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get git root dir")
	}
	// dDir := rootDir + "/src/d"   // dev cli dir
	sDir := rootDir + "/scripts" // scripts dir

	goutils.ExecOpt.Cwd = rootDir
	goutils.Exec("docker compose up -d", goutils.WithCwd(sDir))

	// goutils.Exec("docker build -t 117503445/eventize-dev .", goutils.WithCwd(dDir))

	// docker run -it --rm -v $PWD:/workspace --entrypoint fish 117503445/eventize-dev
	// goutils.Exec(fmt.Sprintf("docker run --rm -v %v:/workspace --entrypoint /workspace/src/d/build_proto.sh 117503445/eventize-dev", rootDir))
	goutils.Exec("docker compose exec --no-TTY --workdir /workspace dev-builder /workspace/src/d/build_proto.sh", goutils.WithCwd(sDir))

	goutils.Exec("cp src/common/service.pb.ts src/fe/src/rpc/service.pb.ts", goutils.WithCwd(rootDir))

	goutils.Exec("docker compose exec --no-TTY fe-dev pnpm build", goutils.WithCwd(sDir))

	// - `-a`：归档模式，表示递归复制文件，并保留文件的权限、时间戳等属性。
	// - `-v`：详细模式，显示同步过程中的信息。
	// - `--delete`：删除目标目录（B）中不存在的源目录（A）中的文件。
	// - `--progress`：在传输过程中显示进度信息。
	// - `--stats`：在传输结束后显示传输统计信息。
	goutils.Exec(fmt.Sprintf("rsync -av --delete --progress --stats %v/ %v/", "./src/fe/dist", "./src/be/cmd/server/dist"), goutils.WithCwd(rootDir))

	// 最新 commit
	r, _ := goutils.Exec("git rev-parse HEAD")
	commit := strings.TrimSuffix(r.Stdout, "\n")
	log.Debug().Str("commit", commit).Msg("commit")

	// 此 commit 对应的 tag
	r, _ = goutils.Exec("git tag --points-at HEAD")
	tagOutput := strings.TrimSuffix(r.Stdout, "\n")
	tags := strings.Split(tagOutput, "\n")
	if len(tags) > 1 {
		log.Warn().Str("tags", tagOutput).Msg("more than one tag")
	}
	tag := tags[0]
	// log.Debug().Str("tag", tag).Msg("tag")
	// 是否有未提交的修改
	dirty := false
	r, _ = goutils.Exec("git status --porcelain")
	if r.Stdout != "" {
		dirty = true
	}
	// log.Debug().Bool("dirty", dirty).Msg("dirty")

	version := ""
	if tag != "" {
		version = tag
	} else {
		version = commit
	}

	if dirty {
		version = version + "-dirty"
	}

	// log.Debug().Str("version", version).Msg("version")

	buildInfo := map[string]interface{}{
		"commit":  commit,
		"tag":     tag,
		"dirty":   dirty,
		"version": version,
		"date":    time.Now().Format("2006-01-02 15:04:05"),
	}

	log.Debug().Interface("buildInfo", buildInfo).Msg("buildInfo")

	if err = goutils.WriteJSON(rootDir+"/src/be/internal/common/build_info.json", buildInfo); err != nil {
		log.Fatal().Err(err).Msg("failed to write build_info.json")
	}

	goutils.Exec("docker compose exec --no-TTY eventize-dev go build ./cmd/server", goutils.WithCwd(sDir))
	goutils.Exec("docker compose exec --no-TTY eventize-dev go build ./cmd/agent", goutils.WithCwd(sDir))
}

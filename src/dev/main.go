package main

import (
	"fmt"
	// "os"

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
	goutils.ExecOpt.Cwd = rootDir

	// docker run -it --rm -v $PWD:/workspace --entrypoint fish 117503445/eventize-dev
	cmd := fmt.Sprintf("docker run --rm -v %v:/workspace --entrypoint /workspace/src/dev/build_proto.sh 117503445/eventize-dev", rootDir)
	goutils.Exec(cmd)

	// srcDir := fmt.Sprintf("%s/src", rootDir)
	// feDir := fmt.Sprintf("%s/%s", srcDir, "fe")

	goutils.Exec("docker compose exec --no-TTY fe-dev pnpm build")


	// feDistDir := fmt.Sprintf("%s/src/fe/dist", rootDir)
	// beDistDir := fmt.Sprintf("%s/src/be/cmd/server/dist", rootDir)
	// // Remove the existing dist directory
	// if err = os.RemoveAll(beDistDir); err != nil {
	// 	log.Fatal().Err(err).Msg("failed to remove existing dist directory")
	// }

	// // Copy the new dist directory from fe to be
	// goutils.Exec("cp -r ./src/fe/dist ./src/be/cmd/server/")

	// - `-a`：归档模式，表示递归复制文件，并保留文件的权限、时间戳等属性。
	// - `-v`：详细模式，显示同步过程中的信息。
	// - `--delete`：删除目标目录（B）中不存在的源目录（A）中的文件。
	// - `--progress`：在传输过程中显示进度信息。
	// - `--stats`：在传输结束后显示传输统计信息。
	goutils.Exec(fmt.Sprintf("rsync -av --delete --progress --stats %v/ %v/", "./src/fe/dist", "./src/be/cmd/server/dist"))
}

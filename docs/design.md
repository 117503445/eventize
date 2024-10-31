https://github.com/117503445/eventize 基于事件驱动的运维平台

## 动机

我平时会有一些运维任务，比如更新一堆节点的 vscode server, 修改 clash 规则后同步到所有运行 clash 的节点，更新系统软件包等等。现有的 Ansible 能基于 SSH 实现批量远程命令执行，但是有 2 个缺点。
- 某些节点可能处于关机状态，控制端需要轮询直到节点开机，很不优雅，轮询时间设长了有延迟，设短了开销大。
- 控制端不一定能通过 SSH 直接连接到被控端。比如控制端在公网，被控端在内网。
所以想设计一个事件驱动的运维平台。控制端跑一个 Web 服务，被控端开机后启动一个 agent，和控制端建立长连接，就可以实现节点上线后触发任务、公网控制端触发内网节点任务了。

## 设计

分为 event action machine task 四大模块。

- event 就是事件，比如 OnConnected 新的 machine 连上、cron 定时任务、webhook 调用控制端 API、用户在前端手动触发任务等。
- action 就是实际要执行的任务，形式可以就是脚本文件。
- machine 指受控机器节点。
- task 就是 action 的某次执行，有任务状态（成功/进行中/失败）、结果等。

流程就是某个 event 会触发多个 machine 上的多个 action，并创建对应的 task。

事件分为 cron，webhook，OnConnected

将事件参数、受控机器 传入 filter，判断是否触发和事件相关联的任务

任务就是执行某个受控端某个路径的脚本

event

action

machine

task 某个 event 触发某个 machine 的 action

## 目标功能特性

- agent 支持不同架构 / 操作系统
- agent 支持通过 Docker 安装，并使用 SSH 获取本机 Shell
- agent 自动升级
- 只需要 server 在公网，agent 可以在内网
- 使用在线 IDE + GitOps，管理 Go 脚本
- 失去连接时，缓存任务
- 自动删除重复的 pending 任务
- 便捷易用的 bucket，基于文件系统


## 可参考技术

https://github.com/gofiber/fiber
https://github.com/graph-gophers/graphql-go

https://github.com/air-verse/air 实时编译

https://github.com/goreleaser/goreleaser 打包为各种平台的二进制文件

## 任务示例

更新 VSCode
1. 下载最新 server
2. 替换到服务器的目标路径
3. 执行命令：更新插件
	
更新fq配置
1. 拉取节点信息
2. 结合域名列表，生成配置
3. 分发配置文件
4. 重启

更新系统
1. 运行系统命令

更新 fish 配置
1. 更新文件

HTTPS 证书管理
1. 申请
2. 推送

家庭网络测速

## 竞品分析

https://github.com/distribworks/dkron

非常靠谱的分布式任务管理系统

https://github.com/i4de/go-ops

任务类型只有 shell, powershell, python, perl。我希望支持 go，并有在线 IDE、GitOps。

https://github.com/ssbeatty/oms


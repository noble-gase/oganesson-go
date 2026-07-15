# og - 项目脚手架

支持创建 `HTTP` 和 `gRPC` 以及 `MCP` 和 `Agent` 项目，支持「单应用」和「多应用」模式

> - 项目结构采用[标准布局](https://github.com/golang-standards/project-layout)
> - 配合 `protoc-gen-og`，支持使用 `proto` 定义API
> - MCP 服务基于 [mcp-go](https://github.com/mark3labs/mcp-go) 构建
> - Agent 服务基于 [adk-go](https://github.com/google/adk-go) 构建

## 安装

```shell
go install github.com/noble-gase/oganesson/cmd/og@latest
```

## 创建项目

<details>
<summary>点击展开</summary>

```shell
og init # 当前目录初始化
og new demo # 创建demo项目
.
├── cmd
│   ├── config.toml
│   └── main.go
├── internal
│   └── app
│       ├── cmd
│       ├── config
│       ├── handler
│       ├── router
│       └── service
├── pkg
│   └── ...
├── Dockerfile
├── dockerun.sh
├── go.mod
├── go.sum
└── README.md
```

</details>

## 创建应用

<details>
<summary>点击展开</summary>

> 多应用项目适用，需在项目根目录执行（即：`go.mod` 所在目录）

```shell
og app foo bar # 创建两个HTTP应用 -- foo 和 bar
.
├── api
│   ├── bar
│   └── foo
├── cmd
│   ├── bar
│   │   ├── config.toml
│   │   └── main.go
│   └── foo
│       ├── config.toml
│       └── main.go
├── internal
│   └── app
│       ├── bar
│       └── foo
├── pkg
│   └── ...
├── foo.dockerfile
├── foo.dockerun.sh
├── bar.dockerfile
├── bar.dockerun.sh
├── go.mod
├── go.sum
└── README.md
```

</details>

## 创建Ent实例

<details>
<summary>点击展开</summary>

#### 单实例

```shell
og ent
.
├── api
├── cmd
├── internal
│   ├── app
│   └── ent
│       ├── schema
│       └── generate.go
├── pkg
│   └── ...
├── go.mod
├── go.sum
└── README.md
```

#### 多实例

```shell
og ent foo bar # 创建Ent实例 -- foo 和 bar
.
├── api
├── cmd
├── internal
│   ├── app
│   └── ent
│       ├── foo
│       │   ├── schema
│       │   └── generate.go
│       └── bar
│           ├── schema
│           └── generate.go
├── pkg
│   └── ...
├── go.mod
├── go.sum
└── README.md
```

</details>

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"charm.land/huh/v2"
	"github.com/noble-gase/oganesson/internal"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
)

func main() {
	cmd := &cobra.Command{
		Use:     "og",
		Short:   "project scaffold",
		Long:    "project scaffold, quickly create a Go project",
		Version: "v0.2.1",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Use == "new" && len(args) != 0 {
				if err := os.MkdirAll(args[0], 0o775); err != nil {
					log.Fatalln("🐛 Mkdir failed:", internal.FmtErr(err))
				}
			}
		},
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println("🐹 Welcome to use noble-gase[Go] scaffold")
		},
	}
	// 注册命令
	cmd.AddCommand(Init(), New(), App(), Ent())
	// 执行
	if err := cmd.Execute(); err != nil {
		log.Fatalln("🐛 Cmd execute failed:", internal.FmtErr(err))
	}
}

type projectKind string

const (
	kindHTTP      projectKind = "http"
	kindHTTPProto projectKind = "http-proto"
	kindGRPC      projectKind = "grpc"
	kindMCP       projectKind = "mcp"
	kindAgent     projectKind = "agent"
)

func Init() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "initialize a project",
		Run: func(_ *cobra.Command, args []string) {
			workDir := "."

			// 判断是否存在go.mod
			_, err := os.Stat("go.mod")
			if err == nil || !os.IsNotExist(err) {
				log.Fatalln("👿 The current directory already exists go.mod, please confirm!")
			}

			dirname, err := internal.GetCurDirName()
			if err != nil {
				log.Fatalln("🐛 Failed to get the current directory:", internal.FmtErr(err))
			}

			kind, mod, apps, err := projectFormSurvey(dirname)
			if err != nil {
				log.Fatalln("🐛 Failed to get project information:", internal.FmtErr(err))
			}

			// 创建项目文件
			fmt.Println("🐹 Create project files")

			switch kind {
			case kindHTTP:
				internal.InitHttpProject(workDir, mod, false, apps...)
			case kindHTTPProto:
				internal.InitHttpProject(workDir, mod, true, apps...)
			case kindGRPC:
				internal.InitGrpcProject(workDir, mod, apps...)
			case kindMCP:
				internal.InitMcpProject(workDir, mod, apps...)
			case kindAgent:
				internal.InitAgentProject(workDir, mod, apps...)
			default:
				log.Fatalln("🐛 Invalid project type:", kind)
			}

			// go mod init
			fmt.Println("⌛️ go mod init")

			modInit := exec.Command("go", "mod", "init", mod)
			modInit.Dir = workDir
			if err := modInit.Run(); err != nil {
				log.Fatalln("🐛 go mod init failed:", internal.FmtErr(err))
			}

			// go mod tidy
			fmt.Println("⌛️ go mod tidy")

			modTidy := exec.Command("go", "mod", "tidy")
			modTidy.Dir = workDir
			modTidy.Stderr = os.Stderr
			if err := modTidy.Run(); err != nil {
				log.Fatalln("🐛 go mod tidy failed:", internal.FmtErr(err))
			}
			fmt.Println("🐹 Project creation completed! please read README")
		},
	}
	return cmd
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "create a project",
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("must specify a project name")
			}
			return nil
		},
		Example: internal.CmdExamples("og new demo"),
		Run: func(_ *cobra.Command, args []string) {
			workDir := args[0]

			// 判断目录是否为空
			if path, ok := internal.IsDirEmpty(workDir); !ok {
				log.Fatalf("👿 The directory(%s) is not empty, please confirm!", path)
			}

			kind, mod, apps, err := projectFormSurvey(workDir)
			if err != nil {
				log.Fatalln("🐛 Failed to get project information:", internal.FmtErr(err))
			}

			// 创建项目文件
			fmt.Println("🐹 Create project files")

			switch kind {
			case kindHTTP:
				internal.InitHttpProject(workDir, mod, false, apps...)
			case kindHTTPProto:
				internal.InitHttpProject(workDir, mod, true, apps...)
			case kindGRPC:
				internal.InitGrpcProject(workDir, mod, apps...)
			case kindMCP:
				internal.InitMcpProject(workDir, mod, apps...)
			case kindAgent:
				internal.InitAgentProject(workDir, mod, apps...)
			default:
				log.Fatalln("🐛 Invalid project type:", kind)
			}

			// go mod init
			fmt.Println("⌛️ go mod init")

			modInit := exec.Command("go", "mod", "init", mod)
			modInit.Dir = workDir
			if err := modInit.Run(); err != nil {
				log.Fatalln("🐛 go mod init failed:", internal.FmtErr(err))
			}

			// go mod tidy
			fmt.Println("⌛️ go mod tidy")

			modTidy := exec.Command("go", "mod", "tidy")
			modTidy.Dir = workDir
			modTidy.Stderr = os.Stderr
			if err := modTidy.Run(); err != nil {
				log.Fatalln("🐛 go mod tidy failed:", internal.FmtErr(err))
			}
			fmt.Println("🐹 Project creation completed! please read README")
		},
	}
	return cmd
}

func App() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app",
		Short: "create an app",
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("must specify an app name")
			}
			return nil
		},
		Example: internal.CmdExamples(
			"og app foo",
			"og app foo bar",
		),
		Run: func(_ *cobra.Command, args []string) {
			fmt.Println("⌛️ Parse go.mod")

			// 读取 go.mod 文件
			data, err := os.ReadFile("go.mod")
			if err != nil {
				log.Fatalln("🐛 Failed to read go.mod file:", internal.FmtErr(err))
			}

			// 解析 go.mod 文件
			f, err := modfile.Parse("go.mod", data, nil)
			if err != nil {
				log.Fatalln("🐛 Failed to parse go.mod file:", internal.FmtErr(err))
			}

			kind, err := appFormSurvey()
			if err != nil {
				log.Fatalln("🐛 Failed to get app information:", internal.FmtErr(err))
			}

			// 创建应用文件
			fmt.Println("🐹 Create app files")

			switch kind {
			case kindHTTP:
				for _, name := range args {
					if path, ok := internal.IsDirEmpty("internal/app/" + name); !ok {
						log.Fatalf("👿 The directory(%s) is not empty, please confirm!", path)
					}
					internal.InitHttpApp(f.Module.Mod.Path, name, false)
				}
			case kindHTTPProto:
				for _, name := range args {
					if path, ok := internal.IsDirEmpty("internal/app/" + name); !ok {
						log.Fatalf("👿 The directory(%s) is not empty, please confirm!", path)
					}
					internal.InitHttpApp(f.Module.Mod.Path, name, true)
				}
			case kindGRPC:
				for _, name := range args {
					if path, ok := internal.IsDirEmpty("internal/app/" + name); !ok {
						log.Fatalf("👿 The directory(%s) is not empty, please confirm!", path)
					}
					internal.InitGrpcApp(f.Module.Mod.Path, name)
				}
			case kindMCP:
				for _, name := range args {
					if path, ok := internal.IsDirEmpty("internal/app/" + name); !ok {
						log.Fatalf("👿 The directory(%s) is not empty, please confirm!", path)
					}
					internal.InitMcpApp(f.Module.Mod.Path, name)
				}
			case kindAgent:
				for _, name := range args {
					if path, ok := internal.IsDirEmpty("internal/app/" + name); !ok {
						log.Fatalf("👿 The directory(%s) is not empty, please confirm!", path)
					}
					internal.InitAgentApp(f.Module.Mod.Path, name)
				}
			default:
				log.Fatalln("🐛 Invalid app type:", kind)
			}

			// go mod tidy
			fmt.Println("⌛️ go mod tidy")

			modTidy := exec.Command("go", "mod", "tidy")
			modTidy.Stderr = os.Stderr
			if err := modTidy.Run(); err != nil {
				log.Fatalln("🐛 go mod tidy failed:", internal.FmtErr(err))
			}
			fmt.Println("🐹 App creation completed! please read README")
		},
	}
	return cmd
}

func Ent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ent",
		Short: "create an ent module",
		Example: internal.CmdExamples(
			"👉 -- default instance --",
			"og ent",
			"",
			"👉 -- specify name --",
			"og ent foo",
			"og ent foo bar",
		),
		Run: func(_ *cobra.Command, args []string) {
			fmt.Println("⌛️ Parse go.mod")

			// 读取 go.mod 文件
			data, err := os.ReadFile("go.mod")
			if err != nil {
				log.Fatalln("🐛 Failed to read go.mod file:", internal.FmtErr(err))
			}

			// 解析 go.mod 文件
			f, err := modfile.Parse("go.mod", data, nil)
			if err != nil {
				log.Fatalln("🐛 Failed to parse go.mod file:", internal.FmtErr(err))
			}

			// 创建Ent文件
			fmt.Println("🐹 Create ent file")

			if len(args) != 0 {
				for _, name := range args {
					if path, ok := internal.IsDirEmpty("internal/ent/" + name); !ok {
						log.Fatalf("👿 The directory(%s) is not empty, please confirm!", path)
					}
					internal.InitEnt(f.Module.Mod.Path, name)
				}
			} else {
				if path, ok := internal.IsDirEmpty("internal/ent"); !ok {
					log.Fatalf("👿 The directory(%s) is not empty, please confirm!", path)
				}
				internal.InitEnt(f.Module.Mod.Path)
			}

			// go mod tidy
			fmt.Println("⌛️ go mod tidy")

			modTidy := exec.Command("go", "mod", "tidy")
			modTidy.Stderr = os.Stderr
			if err := modTidy.Run(); err != nil {
				log.Fatalln("🐛 go mod tidy failed:", internal.FmtErr(err))
			}

			// ent generate
			fmt.Println("⌛️ Ent generate")

			if len(args) != 0 {
				for _, name := range args {
					entGen := exec.Command("go", "generate", "./internal/ent/"+name)
					if err := entGen.Run(); err != nil {
						log.Fatalln("🐛 Ent generate failed:", internal.FmtErr(err))
					}
				}
			} else {
				entGen := exec.Command("go", "generate", "./internal/ent")
				if err := entGen.Run(); err != nil {
					log.Fatalln("🐛 Ent generate failed:", internal.FmtErr(err))
				}
			}

			// go mod tidy
			fmt.Println("⌛️ go mod tidy")

			modClean := exec.Command("go", "mod", "tidy")
			modClean.Stderr = os.Stderr
			if err := modClean.Run(); err != nil {
				log.Fatalln("🐛 go mod tidy failed:", internal.FmtErr(err))
			}
			fmt.Println("🐹 Ent module creation completed! please read README")
		},
	}
	return cmd
}

func projectFormSurvey(defaultMod string) (kind projectKind, mod string, apps []string, err error) {
	var appsRaw string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[projectKind]().
				Title("Project type").
				Options(
					huh.NewOption("HTTP", kindHTTP),
					huh.NewOption("HTTP (proto)", kindHTTPProto),
					huh.NewOption("gRPC", kindGRPC),
					huh.NewOption("MCP", kindMCP),
					huh.NewOption("Agent", kindAgent),
				).
				Value(&kind).
				Validate(func(k projectKind) error {
					if len(k) == 0 {
						return errors.New("please choose a project type")
					}
					return nil
				}),

			huh.NewInput().
				Title("Module name").
				Placeholder(defaultMod).
				Value(&mod),

			huh.NewInput().
				Title("Apps (comma separated, optional)").
				Placeholder("foo,bar").
				Value(&appsRaw),
		),
	)

	if err = form.Run(); err != nil {
		return
	}

	mod = strings.TrimSpace(mod)
	if len(mod) == 0 {
		mod = defaultMod
	}

	if len(strings.TrimSpace(appsRaw)) != 0 {
		for items := range strings.SplitSeq(appsRaw, ",") {
			if name := strings.TrimSpace(items); len(name) != 0 {
				apps = append(apps, name)
			}
		}
	}
	return
}

func appFormSurvey() (kind projectKind, err error) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[projectKind]().
				Title("App type").
				Options(
					huh.NewOption("HTTP", kindHTTP),
					huh.NewOption("HTTP (proto)", kindHTTPProto),
					huh.NewOption("gRPC", kindGRPC),
					huh.NewOption("MCP", kindMCP),
					huh.NewOption("Agent", kindAgent),
				).
				Value(&kind).
				Validate(func(k projectKind) error {
					if len(k) == 0 {
						return errors.New("please choose a project type")
					}
					return nil
				}),
		),
	)
	err = form.Run()
	return
}

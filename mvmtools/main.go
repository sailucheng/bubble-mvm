package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/spf13/cobra"
)

const (
	searchGoModMaxDepth = 10
)

var (
	normal    = lipgloss.NewStyle().Foreground(lipgloss.Color("152")).Render
	highlight = lipgloss.NewStyle().Foreground(lipgloss.Color("154")).MarginTop(1).Render
)

var (
	verbose bool
	outPath string
	name    string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "mvm",
		Short: "Generate cli tool",
		Long:  "Mvm is cli generate tool for mvm project.",
	}
	genCmd := buildGenerateCmd()
	rootCmd.AddCommand(genCmd)
	rootCmd.SilenceErrors = true
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Ops!", err)
		os.Exit(1)
	}
}

func buildGenerateCmd() *cobra.Command {
	genCmd := &cobra.Command{
		Use:   "gen",
		Long:  "Generate mvm components",
		Short: "Generate mvm components",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if len(name) == 0 {
				return fmt.Errorf("--name is required, e.g., 'gen --name myProject'")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := buildContext()
			if err != nil {
				return err
			}
			return projectGenerate(ctx)
		},
	}
	genCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "Component name")
	genCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	genCmd.PersistentFlags().StringVarP(&outPath, "out", "o", "", "Component out path")

	controllerCmd := &cobra.Command{
		Use:   "controller",
		Long:  "Generate controller components",
		Short: "Generate controller",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := buildContext()
			if err != nil {
				return err
			}
			return controllerGenerate(ctx)
		},
	}

	viewCmd := &cobra.Command{
		Use:   "view",
		Long:  "Generate view components",
		Short: "Generate view components",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := buildContext()

			if err != nil {
				return err
			}
			return viewGenerate(ctx)
		},
	}

	modelCmd := &cobra.Command{
		Use:   "model",
		Long:  "Generate model components",
		Short: "Generate model components",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := buildContext()

			if err != nil {
				return err
			}
			return modelGenerate(ctx)
		},
	}

	genCmd.AddCommand(controllerCmd, viewCmd, modelCmd)
	return genCmd
}

func buildContext() (*context, error) {
	// prepare logging
	logger := loggerFactory(verbose)
	ctx := context{
		rawName: name,
		Name:    toPascalCase(name),
		verbose: verbose,
		logger:  logger,
	}
	var (
		module string
		err    error
	)
	// check path exists and it's must be a directory.
	// if not exists then create it.
	// if user didn't specify out path , use current path.
	ctx.path, err = checkPath(outPath)
	if err != nil {
		return nil, fmt.Errorf("cannot get outpath %w", err)
	}

	// find go.mod file in parents.
	ctx.Root, module, err = findGoMod(ctx.path, searchGoModMaxDepth)
	if err != nil {
		return nil, err
	}

	ctx.Module, err = calculateModule(module, ctx.rawName, ctx.Root, ctx.path)
	if err != nil {
		return nil, err
	}
	return &ctx, nil
}

type context struct {
	// 用户指定项目的根目录,向上查找go.mod所在的父级文件夹
	Root string
	// 符合golang命名规范的项目名称
	Name    string
	verbose bool
	logger  *logger
	// 指定输出目录
	path string
	// 用户指定的项目名称，用于生成目录
	rawName string
	Module  string
}
type task struct {
	source       string
	generateMode mode
}

func createTask(source string, m mode) task {
	return task{
		source:       source,
		generateMode: m,
	}
}

func viewGenerate(ctx *context) error {
	name := strings.ToLower(string(ctx.Name[0])) + ctx.Name[1:]
	return generate(ctx, createTask(name+"View.go", ModeView))
}
func modelGenerate(ctx *context) error {
	//file first chat is lowercase
	name := strings.ToLower(string(ctx.Name[0])) + ctx.Name[1:]
	return generate(ctx, createTask(name+"Model.go", ModeModel))
}

func controllerGenerate(ctx *context) error {
	//file first chat is lowercase
	name := strings.ToLower(string(ctx.Name[0])) + ctx.Name[1:]
	return generate(ctx, createTask(name+"Controller.go", ModeController))
}
func projectGenerate(ctx *context) error {
	//file first chat is lowercase
	name := strings.ToLower(string(ctx.Name[0])) + ctx.Name[1:]
	return generate(ctx,
		createTask("main.go", ModeMain),
		createTask(name+"Model.go", ModeModel),
		createTask(name+"Controller.go", ModeController),
		createTask(name+"View.go", ModeView),
	)
}

func generate(ctx *context, tasks ...task) error {
	// use parallel
	var wg sync.WaitGroup

	var errs []error
	for _, task := range tasks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := generateSourceFile(task.source, task.generateMode, ctx); err != nil {
				errs = append(errs, err)
			}
		}()
	}
	wg.Wait()
	// aggregate errors
	if len(errs) > 0 {
		var b strings.Builder
		b.WriteString("generate error\n")
		for i, err := range errs {
			fmt.Fprintf(&b, "%d %s\n", i, err.Error())
		}
		return errors.New(b.String())
	}
	return nil
}

func generateSourceFile(sourcePath string, mod mode, ctx *context) error {
	var (
		folder string
		b      strings.Builder
	)
	if mod != ModeMain {
		folder = strings.ToLower(mod.String()) + "s"
	}
	sourcePath = filepath.Join(ctx.path, folder, sourcePath)
	ctx.logger.log("creating folder %s", filepath.Dir(sourcePath))
	if err := mkdir(filepath.Dir(sourcePath)); err != nil {
		return err
	}
	ctx.logger.log("creating file %s", sourcePath)
	fs, err := os.Create(sourcePath)
	if err != nil {
		return fmt.Errorf("cannot create %s, %v", sourcePath, err)
	}
	defer fs.Close()

	if err := parseTemplate(ctx, templateMap[mod], &b, fs); err != nil {
		return fmt.Errorf("generate %s failed, %s", mod, err)
	}

	if ctx.verbose {
		content := b.String()
		ctx.logger.log("%s\n%s", highlight(fmt.Sprintf("generate %s successfully.\n", sourcePath)), normal(content))
	}
	return nil
}

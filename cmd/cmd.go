package cmd

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Execute() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	cobra.EnableCommandSorting = false
	root.SetUsageTemplate(helpTemplate)
	return fixHelpCommand(root).ExecuteContext(ctx)
}

var root = &cobra.Command{
	Use:   "metaman",
	Short: "媒体管理工具",
}

func AddToRoot(cmd *cobra.Command) {
	root.AddCommand(fixHelpCommand(cmd))
}

func fixHelpCommand(cmd *cobra.Command) *cobra.Command {
	// cmd.CompletionOptions.HiddenDefaultCmd = true
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultCompletionCmd()
	cmd.InitDefaultHelpFlag()

	for _, c := range cmd.Commands() {
		switch c.Name() {
		case "help":
			c.Short = "显示帮助"
		case "completion":
			c.Short = "为指定的shell生成自动完成脚本"
		}
	}

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Name == "help" {
			f.Usage = "显示帮助"
			// f.Hidden = true
		}
	})

	return cmd
}

const helpTemplate = `Usage:{{if .Runnable}}
{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
{{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
{{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Available Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
{{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

使用 "{{.CommandPath}} [command] --help" 获取更多关于此命令的信息.{{end}}

`

// `命令格式:{{if .Runnable}}
// {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
// {{.CommandPath}} [命令]{{end}}{{if gt (len .Aliases) 0}}

// 别名:
// {{.NameAndAliases}}{{end}}{{if .HasExample}}

// 例子:
// {{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

// 命令:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
// {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

// {{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
// {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

// 附加命令:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
// {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

// 命令参数:
// {{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

// 全局参数:
// {{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

// 更多帮助:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
// {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

// 使用 "{{.CommandPath}} [命令] --help" 获取更多关于此命令的信息.{{end}}
// `

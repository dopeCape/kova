package main

import (
	install "github.com/dopeCape/kova/cmd/cli/install"
	"github.com/dopeCape/kova/cmd/cli/interactive"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "kova-cli [command] [flags]"}
	rootCmd.AddCommand(install.InstallCmd, interactive.InteractiveCmd)
	rootCmd.Execute()
}

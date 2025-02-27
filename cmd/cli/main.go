package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "flow",
	Short: "TemplatesFlow CLI",
	Long: `A sample CLI application that showcases:
- Command/subcommand structure
- Flag parsing
- Auto-generated help and usage information
- Shell completion support for commands and flags`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Welcome to your CLI app!")
	},
}

// completionCmd generates shell completion scripts.
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `To load completions:

Bash:
  $ source <(app completion bash)

Zsh:
  $ source <(app completion zsh)

Fish:
  $ app completion fish | source

PowerShell:
  PS> app completion powershell | Out-String | Invoke-Expression
`,
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		switch args[0] {
		case "bash":
			err = rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			err = rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			err = rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			err = rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}

		log.Fatalf("%s\n", err)
	},
}

func main() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(completionCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%s\n", err)
	}
}

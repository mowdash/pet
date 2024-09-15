package cmd

import (
	"fmt"
	"os"

	"github.com/knqyf263/pet/config"
	"github.com/spf13/cobra"
	"gopkg.in/alessio/shellescape.v1"
)

// var delimiter string

// searchCmd represents the search command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete snippet",
	Long:  `Delete snippet interactively (default filtering tool: peco)`,
	RunE:  delete,
}

func delete(cmd *cobra.Command, args []string) (err error) {
	flag := config.Flag

	var options []string
	if flag.Query != "" {
		options = append(options, fmt.Sprintf("--query %s", shellescape.Quote(flag.Query)))
	}

	snippetFile, err := selectFile(options, flag.FilterTag)
	if err != nil {
		return err
	}

	if snippetFile == "" {
		fmt.Println("No snippet file selected")
		return nil
	} else {
		os.Remove(snippetFile)
		fmt.Println("deleted", snippetFile)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVarP(&config.Flag.Color, "color", "", false,
		`Enable colorized output (only fzf)`)
	deleteCmd.Flags().StringVarP(&config.Flag.Query, "query", "q", "",
		`Initial value for query`)
	deleteCmd.Flags().StringVarP(&config.Flag.FilterTag, "tag", "t", "",
		`Filter tag`)
	deleteCmd.Flags().StringVarP(&config.Flag.Delimiter, "delimiter", "d", "; ",
		`Use delim as the command delimiter character`)
}

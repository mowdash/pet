package cmd

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/knqyf263/pet/config"
	"github.com/knqyf263/pet/snippet"
	petSync "github.com/knqyf263/pet/sync"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/alessio/shellescape.v1"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit snippet file",
	Long:  `Edit snippet file (default: opened by vim)`,
	RunE:  edit,
}

func edit(cmd *cobra.Command, args []string) (err error) {
	flag := config.Flag
	editor := config.Conf.General.Editor
	snippetFile := config.Conf.General.SnippetFile

	var options []string
	if flag.Query != "" {
		options = append(options, fmt.Sprintf("--query %s", shellescape.Quote(flag.Query)))
	}

	if len(config.Conf.General.SnippetDirs) > 0 {
		snippetFile, err = selectFile(options, flag.FilterTag)
		if err != nil {
			return err
		}
	}

	if snippetFile == "" {
		return errors.New("No sippet file seleted")
	}

	// file content before editing
	before := fileContent(snippetFile)

	err = editFile(editor, snippetFile, 0)
	if err != nil {
		return
	}

	// file content after editing
	after := fileContent(snippetFile)

	var afterSnippets snippet.Snippets
	err = toml.Unmarshal([]byte(after), &afterSnippets)
	if err != nil {
		return err
	}

	// something changed so just delete and write a new file
	if before != after {
		if err = os.Remove(snippetFile); err != nil {
			return err
		}
		if err = afterSnippets.Save(); err != nil {
			return err
		}
	}

	if config.Conf.Gist.AutoSync {
		return petSync.AutoSync(snippetFile)
	}

	return nil
}

func fileContent(fname string) string {
	data, _ := os.ReadFile(fname)
	return string(data)
}

func init() {
	RootCmd.AddCommand(editCmd)
	editCmd.Flags().StringVarP(&config.Flag.Query, "query", "q", "",
		`Initial value for query`)
	editCmd.Flags().StringVarP(&config.Flag.FilterTag, "tag", "t", "",
		`Filter tag`)
}

package cmd

import (
	"fmt"
	"strings"

	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/cli/cli/v2/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List available commands.",
		Long:    "Displays list of parsed command blocks, their name, number of commands in a block, and description from a given markdown file, such as README.md.",
		RunE: func(cmd *cobra.Command, args []string) error {
			blocks, err := getCodeBlocks()
			if err != nil {
				return err
			}

			// there might be a way to get this from cmd
			io := iostreams.System()
			table := utils.NewTablePrinter(io)

			// table header
			table.AddField(strings.ToUpper("Name"), nil, nil)
			table.AddField(strings.ToUpper("First Command"), nil, nil)
			table.AddField(strings.ToUpper("# of Commands"), nil, nil)
			table.AddField(strings.ToUpper("Description"), nil, nil)
			table.EndRow()

			for _, block := range blocks {
				table.AddField(block.Name(), nil, nil)
				table.AddField(block.Line(0), nil, nil)
				table.AddField(fmt.Sprintf("%d", block.LineCount()), nil, nil)
				table.AddField(block.Intro(), nil, nil)
				table.EndRow()
			}

			return errors.Wrap(table.Render(), "failed to render")
		},
	}

	return &cmd
}

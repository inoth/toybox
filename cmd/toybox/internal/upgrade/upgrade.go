package upgrade

import (
	"fmt"

	"github.com/inoth/toybox/cmd/toybox/internal/base"
	"github.com/spf13/cobra"
)

// CmdUpgrade represents the upgrade command.
var CmdUpgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the toybox tools",
	Long:  "Upgrade the toybox tools. Example: toybox upgrade",
	Run:   Run,
}

func Run(_ *cobra.Command, _ []string) {
	err := base.GoInstall(
		"github.com/inoth/toybox/cmd/toybox@latest",
	)
	if err != nil {
		fmt.Println(err)
	}
}

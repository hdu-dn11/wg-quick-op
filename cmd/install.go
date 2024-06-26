package cmd

import (
	"github.com/hdu-dn11/wg-quick-op/daemon"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install wg-quick-op to /usr/sbin/wg-quick-op",
	Run: func(cmd *cobra.Command, args []string) {
		daemon.Install()
		daemon.AddService()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

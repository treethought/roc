package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/treethought/roc"
)

var config string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the roc kernel",
	Run: func(cmd *cobra.Command, args []string) {
		k := roc.NewKernel()

		fmt.Println("loading space definitions")
		spaces := roc.LoadSpaces(config)
		k.Register(spaces...)

		ctx := roc.NewRequestContext("res://hello", roc.Source)
        resp, err := k.Dispatch(ctx)
        if err != nil {
            fmt.Println(err)
        }
        fmt.Print(resp)

	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&config, "config", "c", "config.yaml", "config file")

}

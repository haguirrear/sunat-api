package comprobante

import (
	"github.com/haguirrear/sunatapi/cmd"
	"github.com/spf13/cobra"
)

var ComprobanteCmd = &cobra.Command{
	Use:   "comprobante",
	Short: "Realizar operaciones con un comprobante de SUNAT",
	Long:  "Realizar operaciones con un comprobante de SUNAT",
}

func init() {
	cmd.RootCmd.AddCommand(ComprobanteCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reciboCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reciboCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

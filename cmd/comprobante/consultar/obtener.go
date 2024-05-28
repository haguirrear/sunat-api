package consultar

import (
	"fmt"
	"os"

	root "github.com/haguirrear/sunatapi/cmd"
	"github.com/haguirrear/sunatapi/cmd/comprobante"
	"github.com/haguirrear/sunatapi/pkg/sunat"
	"github.com/spf13/cobra"
)

var outputFolder string

var ConsultarCmd = &cobra.Command{
	Use:   "obtener [Número de Ticket]",
	Short: "Consulta un comprobante enviado por medio de su número de Ticket ",
	Long: `Consulta un comprobante enviado por medio de su número de Ticket.
Descarga la guia enviada y en caso de error genera un archivo {numGuia_error.txt}`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ticket := args[0]

		token, err := sunat.GetToken(root.ConfigData.AuthBaseURL, sunat.AuthParams{
			ClientID:     root.ConfigData.ClientID,
			ClientSecret: root.ConfigData.ClientSecret,
			Password:     root.ConfigData.Password,
			Username:     root.ConfigData.User,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		receipt, err := sunat.GetReceipt(root.ConfigData.BaseURL, token, ticket)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		if receipt.ReceiptCertificate == "" {
			fmt.Fprintln(os.Stderr, "error: Obtained receipt has empty ReceiptCertificate")
			os.Exit(1)
		}

		if receipt.IsError() {
			errorLine := fmt.Sprintf("Error Code: %s | Detail: %s", receipt.ResponseCode, receipt.Error.Detail)
			errorFileName := fmt.Sprintf("%s_error.txt", ticket)
			if err := os.WriteFile(errorFileName, []byte(errorLine), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "error: Could not write error file %s with content: %s\nBecause of error: %v\n", errorFileName, errorLine, err)
			}

		}

		if err := sunat.SaveReceipt(receipt.ReceiptCertificate, outputFolder); err != nil {
			fmt.Fprintf(os.Stderr, "error saving receipt: %v\n", err)
		}
	},
}

func init() {
	comprobante.ComprobanteCmd.AddCommand(ConsultarCmd)

	ConsultarCmd.Flags().StringVarP(&outputFolder, "output-folder", "o", ".", "Output folder where to save the receipt. Defaults to current folder.")
}

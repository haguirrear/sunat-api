package enviar

import (
	"fmt"
	"os"

	root "github.com/haguirrear/sunatapi/cmd"
	"github.com/haguirrear/sunatapi/cmd/comprobante"
	"github.com/haguirrear/sunatapi/pkg/sunat"
	"github.com/spf13/cobra"
)

var EnviarCmd = &cobra.Command{
	Use:   "enviar [flags] <ruta recibo>",
	Short: "Envía un comprobante (XML) a SUNAT usando la API REST",
	Long: `Envía un comprobante (XML) a SUNAT usando la API REST. 
El archivo XML debe tener el nombre de acuerdo al formato establecido por SUNAT
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		receipPath := args[0]
		rFile, err := os.Open(receipPath)
		defer rFile.Close()

		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

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

		ticket, err := sunat.ZipAndSendReceipt(root.ConfigData.BaseURL, token, receipPath, rFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Recibo enviado correctamente!")
		fmt.Printf("Se generó el ticket: %s\n", ticket)
	},
}

func init() {
	comprobante.ComprobanteCmd.AddCommand(EnviarCmd)
}

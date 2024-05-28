package procesar

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	root "github.com/haguirrear/sunatapi/cmd"
	"github.com/haguirrear/sunatapi/cmd/comprobante"
	"github.com/haguirrear/sunatapi/pkg/sunat"
	"github.com/spf13/cobra"
)

var errorFolder string
var outputFolder string
var retries int

var ProcesarCmd = &cobra.Command{
	Use:   "procesar [recibo xml para enviar a SUNAT]",
	Short: "Envia un comprobante y luego consulta el mismo usando el API REST de SUNAT",
	Long: `Envia un comprobante y luego consulta el mismo

Espera un momento a que SUNAT haya procesado el comprobante y luego obtiene la respuesta.
En caso de éxito guarda el comprobante procesado, en caso de error guarda un archivo {codComprobante_error.txt} con el error`,
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

		var receipt sunat.GetReceiptResponse

		for i := 0; i < retries; i++ {
			receipt, err = sunat.GetReceipt(root.ConfigData.BaseURL, token, ticket)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}

			switch {
			case receipt.IsProcessing():
				log.Printf("El comprobante esta siendo procesado por sunat...\n")
				time.Sleep(200 * time.Millisecond)
				continue
			case receipt.IsError() || receipt.IsSuccess():
				break
			}
		}

		if receipt.ReceiptCertificate == "" {
			fmt.Fprintln(os.Stderr, "error: Obtained receipt has empty ReceiptCertificate")
			os.Exit(1)
		}

		if receipt.IsError() {
			errorLine := fmt.Sprintf("Error Code: %s | Detail: %s", receipt.ResponseCode, receipt.Error.Detail)
			errorFileName := fmt.Sprintf("%s_error.txt", strings.Split(filepath.Base(receipPath), ".")[0])
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
	comprobante.ComprobanteCmd.AddCommand(ProcesarCmd)
	ProcesarCmd.Flags().StringVarP(&outputFolder, "output-folder", "o", ".", "Carpeta donde guardar el ticket de SUNAT. Si no es proporcionada se guardará en la carpeta actual")
	ProcesarCmd.Flags().StringVarP(&errorFolder, "error-folder", "e", ".", "Carpeta donde guardar el mensaje de error si es que sucede un error")
	ProcesarCmd.Flags().IntVar(&retries, "retries", 3, "Numero de reintentos si el ticket aun sigue procesandose")
}

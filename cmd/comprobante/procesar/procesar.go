package procesar

import (
	"context"
	"encoding/json"
	"fmt"
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

const (
	pollTimeout = 10 * time.Second
)

var ProcesarCmd = &cobra.Command{
	Use:   "procesar [recibo xml para enviar a SUNAT]",
	Short: "Envia un comprobante y luego consulta el mismo usando el API REST de SUNAT",
	Long: `Envia un comprobante y luego consulta el mismo

Espera un momento a que SUNAT haya procesado el comprobante y luego obtiene la respuesta.
En caso de éxito guarda el comprobante procesado, en caso de error guarda un archivo {codComprobante_error.txt} con el error`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s := sunat.Sunat{Logger: root.GetLogger()}
		receipPath := args[0]
		rFile, err := os.Open(receipPath)
		defer rFile.Close()

		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		token, err := s.GetToken(root.ConfigData.AuthBaseURL, sunat.AuthParams{
			ClientID:     root.ConfigData.ClientID,
			ClientSecret: root.ConfigData.ClientSecret,
			Password:     root.ConfigData.Password,
			Username:     root.ConfigData.User,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		ticket, err := s.ZipAndSendReceipt(root.ConfigData.BaseURL, token, receipPath, rFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Recibo enviado correctamente!")
		fmt.Printf("Se generó el ticket: %s\n", ticket)

		fmt.Printf("El comprobante esta siendo procesado por sunat...\n")

		ctx, cancel := context.WithTimeout(context.Background(), pollTimeout)
		defer cancel()

		receipt, err := s.PollReceipt(ctx, root.ConfigData.BaseURL, token, ticket)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if receipt.IsError() {
			errorLine := fmt.Sprintf("Error Code: %s | Detail: %s", receipt.ResponseCode, receipt.Error.Detail)
			fmt.Fprintln(os.Stderr, errorLine)

			receiptFileName := strings.Split(filepath.Base(receipPath), ".")[0]
			errorFileName := fmt.Sprintf("%s_error.txt", receiptFileName)
			errorFileName = filepath.Join(errorFolder, errorFileName)

			if err := os.MkdirAll(errorFolder, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "error ensuring error folder exists: %v\n", err)
			}

			fmt.Fprintf(os.Stderr, "Guardando error en: %s\n", errorFileName)
			if err := os.WriteFile(errorFileName, []byte(errorLine), 0664); err != nil {
				fmt.Fprintf(os.Stderr, "error: Could not write error file %s: %v\n", errorFileName, err)
			}
		}

		if receipt.ReceiptCertificate == "" {
			fmt.Fprintln(os.Stderr, "Se recibió un comprobante vacío")

			r, err := json.MarshalIndent(receipt, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot show response: %v\n", err)
				os.Exit(1)
			}

			fmt.Fprintf(os.Stderr, "Respuesta de Sunat:\n%s\n", string(r))

			os.Exit(1)
		}

		fmt.Println("Comprobante obtenido exitosamente!")
		if err := sunat.SaveReceipt(receipt.ReceiptCertificate, outputFolder); err != nil {
			fmt.Fprintf(os.Stderr, "error saving receipt: %v\n", err)
		}
	},
}

func init() {
	comprobante.ComprobanteCmd.AddCommand(ProcesarCmd)
	ProcesarCmd.Flags().StringVarP(&outputFolder, "output-folder", "o", ".", "Carpeta donde guardar el ticket de SUNAT. Si no es proporcionada se guardará en la carpeta actual")
	ProcesarCmd.Flags().StringVarP(&errorFolder, "error-folder", "e", ".", "Carpeta donde guardar el mensaje de error si es que sucede un error")
}

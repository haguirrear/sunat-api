package procesar

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	root "github.com/haguirrear/sunatapi/cmd"
	"github.com/haguirrear/sunatapi/cmd/comprobante"
	"github.com/haguirrear/sunatapi/pkg/sunat"
	"github.com/haguirrear/sunatapi/pkg/ui/spinner"
	"github.com/spf13/cobra"
)

var errorFolder string
var outputFolder string
var ticketStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#d2ad5f"))
var errorDetailStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true).BorderForeground(lipgloss.Color("63")).Padding(1, 3)

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
			s.Logger.Error(err.Error())
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
			s.Logger.Error(err.Error())
			os.Exit(1)
		}

		fmt.Println("Recibo enviado correctamente!")
		fmt.Printf("Se generó el ticket: %s\n", ticketStyle.Render(ticket))

		var spinnerProgram *tea.Program
		if root.VerboseCount == 0 {
			spinnerProgram = tea.NewProgram(spinner.NewSpinner("El comprobante está siendo procesado por SUNAT"))

			go func() {
				if _, err := spinnerProgram.Run(); err != nil {
					cobra.CheckErr(err)
				}
			}()
		} else {
			s.Logger.Print("El comprobante está siendo procesado por SUNAT...")
		}

		ctx, cancel := context.WithTimeout(context.Background(), pollTimeout)
		defer cancel()

		time.Sleep(2 * time.Second)
		receipt, err := s.PollReceipt(ctx, root.ConfigData.BaseURL, token, ticket)

		if root.VerboseCount == 0 {
			if err := spinnerProgram.ReleaseTerminal(); err != nil {
				s.Logger.Errorf("There was a problem releasing the terminal: %v\n", err)
			}
		}

		if err != nil {
			s.Logger.Error(err.Error())
			os.Exit(1)
		}

		if receipt.IsError() {
			s.Logger.Error("Ocurrió un error recepcionando el recibo procesado")
			s.Logger.SetIndentation(1)
			errorLine := errorDetailStyle.Render(fmt.Sprintf("Error Code: %s | Detail: %s", receipt.ResponseCode, receipt.Error.Detail))
			s.Logger.Error(errorLine)
			s.Logger.ClearIndentation()

			receiptFileName := strings.Split(filepath.Base(receipPath), ".")[0]
			errorFileName := fmt.Sprintf("%s_error.txt", receiptFileName)
			errorFileName = filepath.Join(errorFolder, errorFileName)
			errorAbs, err := filepath.Abs(errorFileName)
			if err == nil {
				errorFileName = errorAbs
			}

			if err := os.MkdirAll(errorFolder, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "error ensuring error folder exists: %v\n", err)
			}

			s.Logger.Printf("Guardando error en: %s", ticketStyle.Render(errorFileName))
			if err := os.WriteFile(errorFileName, []byte(errorLine), 0664); err != nil {
				fmt.Fprintf(os.Stderr, "error: Could not write error file %s: %v\n", errorFileName, err)
			}
		}

		if receipt.ReceiptCertificate == "" {
			s.Logger.Warn("Se recibió un comprobante vacío")
			s.Logger.SetIndentation(1)

			r, err := json.MarshalIndent(receipt, "", "  ")
			if err != nil {
				s.Logger.Errorf("Cannot show response: %v", err)
				os.Exit(1)
			}

			s.Logger.Warn("Respuesta de Sunat:")
			s.Logger.Warn(errorDetailStyle.Render(string(r)))
			s.Logger.ClearIndentation()

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

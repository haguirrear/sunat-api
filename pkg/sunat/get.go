package sunat

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const (
	TicketErrorResponseCode      = "99"
	TicketSuccessResponseCode    = "0"
	TIcketProcessingResponseCode = "98"
)

type TicketError struct {
	NumError string `json:"numError"`
	Detail   string `json:"desError"`
}

type GetReceiptResponse struct {
	ResponseCode       string      `json:"codRespuesta"`
	Error              TicketError `json:"error"`
	ReceiptCertificate string      `json:"arcCdr"`
	CdrGenerated       string      `json:"indCdrGenerado"`
}

func (r GetReceiptResponse) IsError() bool {
	return r.ResponseCode == TicketErrorResponseCode
}

func (r GetReceiptResponse) IsSuccess() bool {
	return r.ResponseCode == TicketSuccessResponseCode
}

func (r GetReceiptResponse) IsProcessing() bool {
	return r.ResponseCode == TIcketProcessingResponseCode
}
func (r GetReceiptResponse) IsCdrGenerated() bool {
	b, err := strconv.ParseBool(r.CdrGenerated)
	if err != nil {
		log.Printf("error parsing indCdrGenerado from response ('%s') into a boolean: %w", r.CdrGenerated, err)
		return false
	}

	return b
}

func GetReceipt(baseURL string, token string, ticket string) (GetReceiptResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	reqURL := fmt.Sprintf("%s/v1/contribuyente/gem/comprobantes/envios/%s", baseURL, ticket)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return GetReceiptResponse{}, fmt.Errorf("error building request for getting receipt: %s: %w", ticket, err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := client.Do(req)
	if err != nil {
		return GetReceiptResponse{}, fmt.Errorf("error getting receipt %s: %w", ticket, err)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return GetReceiptResponse{}, fmt.Errorf("error parsing response body while getting receipt %s: %w", ticket, err)
	}

	if res.StatusCode >= 400 {
		return GetReceiptResponse{}, fmt.Errorf("error getting receipt %s: %s | %s", ticket, res.Status, string(body))
	}

	var resBody GetReceiptResponse
	if err := json.Unmarshal(body, &resBody); err != nil {
		return GetReceiptResponse{}, fmt.Errorf("error parsing body while getting receipt %s, status %s: %w\nBody: %s", ticket, res.Status, err, string(body))
	}

	return resBody, nil
}

func SaveReceipt(receiptB64 string, outputFolder string) error {
	rByte, err := base64.StdEncoding.DecodeString(receiptB64)
	if err != nil {
		return fmt.Errorf("error decoding receipt from base64: %w", err)
	}

	bReader := bytes.NewReader(rByte)

	zipReader, err := zip.NewReader(bReader, int64(bReader.Len()))
	if err != nil {
		return fmt.Errorf("error reading receipt zip: %w", err)
	}

	if len(zipReader.File) == 0 {
		return fmt.Errorf("error with receipt zip: No files found in the zip archive")
	}

	contentReader, err := zipReader.File[0].Open()
	defer contentReader.Close()

	if err != nil {
		return fmt.Errorf("error reading receipt in zip file: %w", err)
	}

	content, err := io.ReadAll(contentReader)
	if err != nil {
		return fmt.Errorf("error reading receipt content: %w", err)
	}

	destFilePath := filepath.Join(outputFolder, zipReader.File[0].Name)

	outFile, err := os.Create(destFilePath)
	defer outFile.Close()

	if err != nil {
		return fmt.Errorf("error creating receipt output file: %w", err)
	}

	if _, err := outFile.Write(content); err != nil {
		return fmt.Errorf("error writing receipt file: %w", err)
	}

	return nil
}

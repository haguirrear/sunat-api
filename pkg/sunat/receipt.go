package sunat

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var ErrorFileNotFound = errors.New("File not found")

func ZipAndSendReceipt(baseURL, authToken, receiptPath string, receiptFile io.Reader) (numTicket string, err error) {
	zipFile, err := createSingleFileZip(receiptPath, receiptFile)

	if err != nil {
		return "", fmt.Errorf("error sending receipt %s: %w", receiptPath, err)
	}

	zipHash, err := HashFileContent(zipFile)
	if err != nil {
		return "", fmt.Errorf("error sendig receipt %s: %w", receiptPath, err)
	}

	zipBase64, err := EncodeFileBase64(zipFile)
	if err != nil {
		return "", fmt.Errorf("error sending receipt %s: %w", receiptPath, err)
	}

	params := SendReceiptParams{
		ReceiptFilePath:    receiptPath,
		ZipFileHash:        zipHash,
		ZipFileBase64:      zipBase64,
		AuthorizationToken: authToken,
	}

	res, err := SendReceipt(baseURL, params)
	if err != nil {
		return "", err
	}

	return res.NumTicket, err
}

type SendReceiptParams struct {
	ReceiptFilePath    string
	ZipFileBase64      string
	ZipFileHash        string
	AuthorizationToken string
}

type SendReceiptResponse struct {
	NumTicket string `json:"numTicket"`
}

func SendReceipt(baseURL string, params SendReceiptParams) (SendReceiptResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	filename := filepath.Base(params.ReceiptFilePath)
	fileWithoutExt := strings.Split(filename, ".")[0]
	payloadMap := map[string]any{
		"archivo": map[string]any{
			"nomArchivo": fileWithoutExt + ".zip",
			"arcGreZip":  params.ZipFileBase64,
			"hashZip":    params.ZipFileHash,
		},
	}

	payload, err := json.Marshal(payloadMap)
	if err != nil {
		return SendReceiptResponse{}, fmt.Errorf("error building send receipt payload: %w", err)
	}

	reqURL := fmt.Sprintf("%s/v1/contribuyente/gem/comprobantes/%s", baseURL, fileWithoutExt)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(payload))
	if err != nil {
		return SendReceiptResponse{}, fmt.Errorf("error building request for send receipt: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", params.AuthorizationToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return SendReceiptResponse{}, fmt.Errorf("error sending receipt %s: %w", params.ReceiptFilePath, err)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return SendReceiptResponse{}, fmt.Errorf("error sending receipt %s while parsing response body: %w", params.ReceiptFilePath, err)
	}

	if res.StatusCode >= 400 {
		return SendReceiptResponse{}, fmt.Errorf("error sending receipt %s: %s | %s", params.ReceiptFilePath, res.Status, string(body))
	}

	var bodyParsed SendReceiptResponse
	if err := json.Unmarshal(body, &bodyParsed); err != nil {
		return SendReceiptResponse{}, fmt.Errorf("error parsing send receipt response body '%s': %w", string(body), err)
	}

	return bodyParsed, nil
}

// Encodes the content of a file in base64
func EncodeFileBase64(file io.Reader) (string, error) {
	bContent, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error encoding file in base64 while reading file: %w", err)
	}

	return base64.StdEncoding.EncodeToString(bContent), nil
}

// Hashes the content of a file with SHA-256
func HashFileContent(file io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", fmt.Errorf("error generating hash for zip file: %w", err)
	}

	return string(h.Sum(nil)), nil
}

// Creates an in memory zipFile with one content (the file in the argument)
// The caller is responsible to close the file after using it
func createSingleFileZip(fileToCompressPath string, file io.Reader) (zipFile *bytes.Buffer, err error) {
	if file == nil {
		return nil, fmt.Errorf("nil file passed to createSingleFileZip")
	}

	buf := new(bytes.Buffer)

	log.Println("creating zip")
	zipWriter := zip.NewWriter(buf)
	defer func() {
		errzip := zipWriter.Close()
		if errzip != nil && err == nil {
			err = fmt.Errorf("error closing zip file: %w", err)
		}
	}()

	log.Println("creating first file inside zip")
	zw, err := zipWriter.Create(filepath.Base(fileToCompressPath))
	if err != nil {
		return nil, fmt.Errorf("error adding %s to zip file: %w", fileToCompressPath, err)
	}

	log.Println("adding file to zip archive")
	if _, err := io.Copy(zw, file); err != nil {
		return nil, fmt.Errorf("error adding %s to zip file: %w", fileToCompressPath, err)
	}

	return buf, nil
}

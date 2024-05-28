package sunat

import (
	"net/http"

	loghttp "github.com/motemen/go-loghttp"
)

var client = &http.Client{
	Transport: &loghttp.Transport{},
}

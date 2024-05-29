package sunat

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/haguirrear/sunatapi/pkg/logger"
	// loghttp "github.com/motemen/go-loghttp"
)

var client = &http.Client{
	// Transport: &loghttp.Transport{},
}

func (s Sunat) doRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	s.Logger.Debugf("-> Request %s", req.URL.String())
	for k, v := range req.Header {
		for _, vv := range v {
			s.Logger.Tracef("Header '%s': '%s'", k, vv)
		}
	}

	if req.Body != nil {
		logReqBody(s.Logger, req)
	}

	res, err := client.Do(req)

	if err != nil {
		return res, err
	}

	s.Logger.Debugf("<- Response %s", res.Status)

	// for k, v := range res.Header {
	// 	for _, vv := range v {
	// 		s.Logger.Debugf("Header '%s': '%s'", k, vv)
	//
	// 	}
	// }

	if res.Body != nil && res.Body != http.NoBody {
		logResBody(s.Logger, res)
	}

	return res, err
}

func logReqBody(logger *logger.Logger, req *http.Request) {

	b, err := req.GetBody()
	if err != nil {
		return
	}
	bodyBytes, err := io.ReadAll(b)
	if err != nil {
		return
	}
	logger.Tracef("Body: %s", string(bodyBytes))
}

func logResBody(logger *logger.Logger, res *http.Response) {
	var buf bytes.Buffer
	b := res.Body
	if _, err := buf.ReadFrom(b); err != nil {
		return
	}
	if err := b.Close(); err != nil {
		return
	}
	res.Body = io.NopCloser(&buf)

	bCopy := io.NopCloser(bytes.NewReader(buf.Bytes()))
	bodyBytes, err := io.ReadAll(bCopy)
	if err != nil {
		return
	}

	isJson := strings.Contains(res.Header.Get("Content-Type"), "application/json")
	if isJson {
		bodyJson := prettyPrintJson(string(bodyBytes))
		logger.Tracef("Body: %s", bodyJson)

	} else {
		logger.Tracef("Body: %s", string(bodyBytes))
	}

}

func prettyPrintJson(jsonString string) string {
	var temp any

	if err := json.Unmarshal([]byte(jsonString), &temp); err != nil {
		return jsonString
	}

	r, err := json.MarshalIndent(temp, "", "  ")
	if err != nil {
		return jsonString
	}

	return string(r)
}

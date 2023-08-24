package azgohacks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type DumpFullPolicy struct {
	Prefix string
}

func (p DumpFullPolicy) Do(req *policy.Request) (*http.Response, error) {
	fmt.Printf("\n\n===> BEGIN: REQUEST (%s) <===\n\n", p.Prefix)

	requestBytes, err := httputil.DumpRequestOut(req.Raw(), false)

	if err != nil {
		return nil, err
	}

	fmt.Println(string(requestBytes))
	fmt.Printf("Body: %s\n", string(FormatRequestBytes(req.Raw())))
	fmt.Printf("\n\n===> END: REQUEST (%s)<===\n\n", p.Prefix)

	resp, err := req.Next()

	if err != nil {
		return nil, err
	}

	fmt.Printf("\n\n===> BEGIN: RESPONSE (%s) <===\n\n", p.Prefix)

	responseBytes, err := httputil.DumpResponse(resp, false)

	if err != nil {
		return nil, err
	}

	fmt.Println(string(responseBytes))
	fmt.Printf("Body: %s\n", string(FormatResponseBytes(resp)))

	fmt.Printf("\n\n===> END: RESPONSE (%s) <===\n\n", p.Prefix)
	return resp, err
}

func FormatRequestBytes(req *http.Request) []byte {
	requestBytes, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err)
	}

	req.Body = io.NopCloser(bytes.NewBuffer(requestBytes))
	return FormatBytes(requestBytes)
}

func FormatResponseBytes(resp *http.Response) []byte {
	requestBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(requestBytes))
	return FormatBytes(requestBytes)
}

func FormatBytes(body []byte) []byte {
	var m *map[string]any
	var l *[]any

	candidates := []any{&m, &l}

	for _, v := range candidates {
		err := json.Unmarshal(body, v)

		if err != nil {
			continue
		}

		if err == nil {
			formattedBytes, err := json.MarshalIndent(v, "  ", "  ")

			if err != nil {
				continue
			}

			return formattedBytes
		}
	}

	// if we can't format it we'll just give it back.
	return body
}

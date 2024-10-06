package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Proxy(c echo.Context, resp *http.Response, err error) error {
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			c.Response().Header().Add(key, value)
		}
	}

	c.Response().WriteHeader(resp.StatusCode)

	_, err = io.Copy(c.Response(), resp.Body)
	return err
}

func Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	if _, err := writer.Write(data); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Decompress(data []byte) (map[string]interface{}, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	decompressedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to read decompressed data: %w", err)
	}

	var result map[string]interface{}
	if err = json.Unmarshal(decompressedData, &result); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to unmarshal JSON data: %w", err)
	}

	return result, nil
}

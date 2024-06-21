package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ReadJSON(c *gin.Context, data interface{}) error {
	maxBytes := 1048576
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, int64(maxBytes))

	doc := json.NewDecoder(c.Request.Body)
	err := doc.Decode(data)
	if err != nil {
		return err
	}

	err = doc.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single json value")
	}

	return nil
}

func WriteJSON(c *gin.Context, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			c.Writer.Header()[key] = value
		}
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)
	_, err = c.Writer.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func ErrorJSON(c *gin.Context, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return WriteJSON(c, statusCode, payload)
}

func ExtractResponseBody(c *gin.Context, response *http.Response) (JsonResponse, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		ErrorJSON(c, err)
		return JsonResponse{}, err
	}

	var respBody JsonResponse
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		ErrorJSON(c, err)
		return JsonResponse{}, err
	}
	return respBody, nil
}

func MakeHTTPRequest(c *gin.Context, method string, url string, payload ...interface{}) (*JsonResponse, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return nil, errors.New("received non-accepted status code")
	}

	respBody, err := ExtractResponseBody(c, response)
	if err != nil {
		return nil, err
	}

	return &respBody, nil
}

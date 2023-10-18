package utils

import (
	"encoding/json"
	"github.com/spf13/cast"
	"io"
	"jtrans/tbox/models"
	"net/http"
	"strings"
)

func DoRequest(method string, url string, headers map[string]string, query map[string]string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if len(query) > 0 {
		q := req.URL.Query()
		for key, value := range query {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	//fmt.Println(req.URL.String())
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return client.Do(req)
}

func UnmarshalJson[T any](resp *http.Response, target *T) error {
	d := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	return d.Decode(target)
}

func FromCookiesString(cookies string) map[string]string {
	tokens := strings.Split(cookies, ";")
	cookiesMap := make(map[string]string)
	for _, token := range tokens {
		kv := strings.Split(token, "=")
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		cookiesMap[key] = value
	}
	return cookiesMap
}

type RequestStream struct {
	data       []byte
	bufferSize int
	curIndex   int
	total      int
	onProgress models.UploadProgressHandler
}

func NewRequestStream(data []byte, bufferSize int, onProgress models.UploadProgressHandler) *RequestStream {
	return &RequestStream{
		data:       data,
		bufferSize: bufferSize,
		curIndex:   0,
		onProgress: onProgress,
		total:      len(data),
	}
}

func (s *RequestStream) Read(p []byte) (n int, err error) {
	if s.curIndex >= s.total {
		return 0, io.EOF
	}
	size := len(p)
	if size > s.bufferSize {
		size = s.bufferSize
	}

	if size == 0 {
		return 0, nil
	}

	if size+s.curIndex > s.total {
		size = s.total - s.curIndex
	}

	copy(p, s.data[s.curIndex:s.curIndex+size])
	s.curIndex += size
	s.onProgress(cast.ToInt64(s.curIndex), cast.ToInt64(s.total))
	return size, nil
}

func IsSuccessStatusCode(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}

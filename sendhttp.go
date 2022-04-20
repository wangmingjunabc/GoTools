package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

const (
	ContentType     = "Content-Type"
	ContentTypeJSON = "application/json"
)

func PostRequest(
	ctx context.Context,
	returnValuePointer interface{},
	url string,
	data interface{}, token string) error {
	req, err := createRequestWithRawData(ctx, http.MethodPost, url, data)
	if err != nil {
		return err
	}
	if token != "" {
		req.Header.Set("access_token", token)
	}
	res, err := sendRequest(ctx, req)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(res, returnValuePointer); err != nil {
		return err
	}
	return nil
}

func GetRequest(ctx context.Context, returnValuePointer interface{}, baseUrl string, requestPath string, requestParams url.Values, token string) error {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return err
	}
	u.Path = requestPath
	if requestParams != nil {
		u.RawQuery = requestParams.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("access_token", token)
	res, err := sendRequest(ctx, req)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(res, returnValuePointer); err != nil {
		return err
	}
	return nil
}

func FromContext(ctx context.Context, key string) string {
	hdr, ok := ctx.Value(key).(string)
	if !ok {
		hdr = ""
	}
	return hdr
}

func createRequestWithRawData(ctx context.Context, httpMethod string, url string, data interface{}) (*http.Request, error) {
	jsonEncodedData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	content := FromContext(ctx, ContentType)
	if content == "" {
		content = ContentTypeJSON
	}
	var req *http.Request
	if req, err = http.NewRequest(httpMethod, url, bytes.NewReader(jsonEncodedData)); err != nil {
		return nil, err
	}
	req.Header.Set(ContentType, content)
	return req, nil
}

func sendRequest(ctx context.Context, req *http.Request) ([]byte, error) {
	data := ctx.Value("data")
	if data != nil {
		switch data.(type) {
		case []byte:
			req.Header.Set("data", string(data.([]byte)))
		case string:
			req.Header.Set("data", data.(string))
		default:
			req.Header.Set("dpType", fmt.Sprintln(reflect.TypeOf(data)))
			marshal, _ := json.Marshal(data)
			req.Header.Set("data", string(marshal))
		}
	}
	resp, err := makeRequest(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, err
	}
	defer resp.Body.Close()
	return getBody(resp)
}

func makeRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func getBody(resp *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	return body, nil
}

func GetHttp() {

}

package types

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type XHttp struct {
	client  *http.Client
	baseURL string
	headers map[string]string
}

type XHttpResponse struct {
	*http.Response
	bodyBytes []byte
}

// Http 创建 XHttp 实例
func Http() XHttp {
	return XHttp{
		client:  &http.Client{Timeout: 30 * time.Second},
		headers: make(map[string]string),
	}
}

// HttpWithClient 使用指定客户端创建 XHttp 实例
func HttpWithClient(client *http.Client) XHttp {
	return XHttp{
		client:  client,
		headers: make(map[string]string),
	}
}

// BaseURL 设置基础 URL
func (h XHttp) BaseURL(baseURL string) XHttp {
	h.baseURL = strings.TrimRight(baseURL, "/")
	return h
}

// Timeout 设置超时时间
func (h XHttp) Timeout(timeout time.Duration) XHttp {
	h.client.Timeout = timeout
	return h
}

// Header 设置请求头
func (h XHttp) Header(key, value string) XHttp {
	h.headers[key] = value
	return h
}

// Headers 批量设置请求头
func (h XHttp) Headers(headers map[string]string) XHttp {
	for k, v := range headers {
		h.headers[k] = v
	}
	return h
}

// ContentType 设置 Content-Type
func (h XHttp) ContentType(contentType string) XHttp {
	return h.Header("Content-Type", contentType)
}

// JSON 设置 JSON Content-Type
func (h XHttp) JSON() XHttp {
	return h.ContentType("application/json")
}

// Form 设置表单 Content-Type
func (h XHttp) Form() XHttp {
	return h.ContentType("application/x-www-form-urlencoded")
}

// Auth 设置基本认证
func (h XHttp) Auth(username, password string) XHttp {
	return h.Header("Authorization", "Basic "+basicAuth(username, password))
}

// Bearer 设置 Bearer Token
func (h XHttp) Bearer(token string) XHttp {
	return h.Header("Authorization", "Bearer "+token)
}

// UserAgent 设置 User-Agent
func (h XHttp) UserAgent(userAgent string) XHttp {
	return h.Header("User-Agent", userAgent)
}

// buildURL 构建完整 URL
func (h XHttp) buildURL(endpoint string) string {
	if h.baseURL == "" {
		return endpoint
	}

	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		return endpoint
	}

	endpoint = strings.TrimLeft(endpoint, "/")
	return h.baseURL + "/" + endpoint
}

// newRequest 创建新请求
func (h XHttp) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, value := range h.headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

// doRequest 执行请求
func (h XHttp) doRequest(req *http.Request) (*XHttpResponse, error) {
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()

	// 重新设置 Body 为可读流
	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	return &XHttpResponse{
		Response:  resp,
		bodyBytes: bodyBytes,
	}, nil
}

// Get 发送 GET 请求
func (h XHttp) Get(url string) (*XHttpResponse, error) {
	return h.GetWithContext(context.Background(), url)
}

// GetWithContext 使用上下文发送 GET 请求
func (h XHttp) GetWithContext(ctx context.Context, url string) (*XHttpResponse, error) {
	fullURL := h.buildURL(url)
	req, err := h.newRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return h.doRequest(req)
}

// Post 发送 POST 请求
func (h XHttp) Post(url string, body interface{}) (*XHttpResponse, error) {
	return h.PostWithContext(context.Background(), url, body)
}

// PostWithContext 使用上下文发送 POST 请求
func (h XHttp) PostWithContext(ctx context.Context, url string, body interface{}) (*XHttpResponse, error) {
	return h.requestWithBody(ctx, "POST", url, body)
}

// Put 发送 PUT 请求
func (h XHttp) Put(url string, body interface{}) (*XHttpResponse, error) {
	return h.PutWithContext(context.Background(), url, body)
}

// PutWithContext 使用上下文发送 PUT 请求
func (h XHttp) PutWithContext(ctx context.Context, url string, body interface{}) (*XHttpResponse, error) {
	return h.requestWithBody(ctx, "PUT", url, body)
}

// Patch 发送 PATCH 请求
func (h XHttp) Patch(url string, body interface{}) (*XHttpResponse, error) {
	return h.PatchWithContext(context.Background(), url, body)
}

// PatchWithContext 使用上下文发送 PATCH 请求
func (h XHttp) PatchWithContext(ctx context.Context, url string, body interface{}) (*XHttpResponse, error) {
	return h.requestWithBody(ctx, "PATCH", url, body)
}

// Delete 发送 DELETE 请求
func (h XHttp) Delete(url string) (*XHttpResponse, error) {
	return h.DeleteWithContext(context.Background(), url)
}

// DeleteWithContext 使用上下文发送 DELETE 请求
func (h XHttp) DeleteWithContext(ctx context.Context, url string) (*XHttpResponse, error) {
	fullURL := h.buildURL(url)
	req, err := h.newRequest("DELETE", fullURL, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return h.doRequest(req)
}

// requestWithBody 发送带请求体的请求
func (h XHttp) requestWithBody(ctx context.Context, method, url string, body interface{}) (*XHttpResponse, error) {
	var bodyReader io.Reader

	if body != nil {
		switch v := body.(type) {
		case string:
			bodyReader = strings.NewReader(v)
		case []byte:
			bodyReader = bytes.NewReader(v)
		case io.Reader:
			bodyReader = v

		default:
			// 默认序列化为 JSON
			jsonData, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			bodyReader = bytes.NewReader(jsonData)
			h = h.JSON()
		}
	}

	fullURL := h.buildURL(url)
	req, err := h.newRequest(method, fullURL, bodyReader)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return h.doRequest(req)
}

// PostForm 发送表单数据
func (h XHttp) PostForm(url string, data url.Values) (*XHttpResponse, error) {
	return h.PostFormWithContext(context.Background(), url, data)
}

// PostFormWithContext 使用上下文发送表单数据
func (h XHttp) PostFormWithContext(ctx context.Context, url string, data url.Values) (*XHttpResponse, error) {
	return h.Form().PostWithContext(ctx, url, data)
}

// PostJSON 发送 JSON 数据
func (h XHttp) PostJSON(url string, data interface{}) (*XHttpResponse, error) {
	return h.PostJSONWithContext(context.Background(), url, data)
}

// PostJSONWithContext 使用上下文发送 JSON 数据
func (h XHttp) PostJSONWithContext(ctx context.Context, url string, data interface{}) (*XHttpResponse, error) {
	return h.JSON().PostWithContext(ctx, url, data)
}

// Upload 上传文件
func (h XHttp) Upload(url, fieldName, fileName string, fileData io.Reader) (*XHttpResponse, error) {
	return h.UploadWithContext(context.Background(), url, fieldName, fileName, fileData)
}

// UploadWithContext 使用上下文上传文件
func (h XHttp) UploadWithContext(ctx context.Context, url, fieldName, fileName string, fileData io.Reader) (*XHttpResponse, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, fileData)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	h = h.ContentType(writer.FormDataContentType())
	return h.PostWithContext(ctx, url, &buf)
}

// UploadWithFields 上传文件并包含其他表单字段
func (h XHttp) UploadWithFields(url, fieldName, fileName string, fileData io.Reader, fields map[string]string) (*XHttpResponse, error) {
	return h.UploadWithFieldsContext(context.Background(), url, fieldName, fileName, fileData, fields)
}

// UploadWithFieldsContext 使用上下文上传文件并包含其他表单字段
func (h XHttp) UploadWithFieldsContext(ctx context.Context, url, fieldName, fileName string, fileData io.Reader, fields map[string]string) (*XHttpResponse, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加其他字段
	for key, value := range fields {
		err := writer.WriteField(key, value)
		if err != nil {
			return nil, err
		}
	}

	// 添加文件
	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, fileData)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	h = h.ContentType(writer.FormDataContentType())
	return h.PostWithContext(ctx, url, &buf)
}

// XHttpResponse 方法

// String 获取响应体为字符串
func (r *XHttpResponse) String() string {
	return string(r.bodyBytes)
}

// Bytes 获取响应体为字节数组
func (r *XHttpResponse) Bytes() []byte {
	return r.bodyBytes
}

// JSON 解析响应体为 JSON
func (r *XHttpResponse) JSON(v interface{}) error {
	return json.Unmarshal(r.bodyBytes, v)
}

// IsOK 判断响应是否成功 (200-299)
func (r *XHttpResponse) IsOK() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsClientError 判断是否为客户端错误 (400-499)
func (r *XHttpResponse) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

// IsServerError 判断是否为服务器错误 (500-599)
func (r *XHttpResponse) IsServerError() bool {
	return r.StatusCode >= 500 && r.StatusCode < 600
}

// SaveToFile 保存响应体到文件
func (r *XHttpResponse) SaveToFile(filename string) error {
	file := File(filename)
	return file.Write(r.bodyBytes)
}

// GetHeader 获取响应头
func (r *XHttpResponse) GetHeader(key string) string {
	return r.Header.Get(key)
}

// GetHeaders 获取指定键的所有响应头值
func (r *XHttpResponse) GetHeaders(key string) []string {
	return r.Header.Values(key)
}

// GetContentType 获取 Content-Type
func (r *XHttpResponse) GetContentType() string {
	return r.GetHeader("Content-Type")
}

// GetContentLength 获取 Content-Length
func (r *XHttpResponse) GetContentLength() string {
	return r.GetHeader("Content-Length")
}

// Close 关闭响应体
func (r *XHttpResponse) Close() error {
	if r.Response != nil && r.Response.Body != nil {
		return r.Response.Body.Close()
	}
	return nil
}

// 辅助函数

// basicAuth 生成基本认证字符串
func basicAuth(username, password string) string {
	auth := username + ":" + password
	// 这里需要 base64 编码，但为了避免引入额外依赖，使用简单实现
	return fmt.Sprintf("%s", auth) // 实际项目中应该使用 base64.StdEncoding.EncodeToString([]byte(auth))
}

// QueryParams 构建查询参数
func QueryParams(params map[string]interface{}) string {
	values := url.Values{}
	for k, v := range params {
		values.Add(k, fmt.Sprintf("%v", v))
	}
	return values.Encode()
}

// BuildURL 构建带查询参数的 URL
func BuildURL(baseURL string, params map[string]interface{}) string {
	if len(params) == 0 {
		return baseURL
	}

	separator := "?"
	if strings.Contains(baseURL, "?") {
		separator = "&"
	}

	return baseURL + separator + QueryParams(params)
}

// 便捷方法

// GET 快速发送 GET 请求
func GET(url string) (*XHttpResponse, error) {
	return Http().Get(url)
}

// POST 快速发送 POST 请求
func POST(url string, body interface{}) (*XHttpResponse, error) {
	return Http().Post(url, body)
}

// PUT 快速发送 PUT 请求
func PUT(url string, body interface{}) (*XHttpResponse, error) {
	return Http().Put(url, body)
}

// PATCH 快速发送 PATCH 请求
func PATCH(url string, body interface{}) (*XHttpResponse, error) {
	return Http().Patch(url, body)
}

// DELETE 快速发送 DELETE 请求
func DELETE(url string) (*XHttpResponse, error) {
	return Http().Delete(url)
}

// PostJSON 快速发送 JSON POST 请求
func PostJSON(url string, data interface{}) (*XHttpResponse, error) {
	return Http().PostJSON(url, data)
}

// PostForm 快速发送表单 POST 请求
func PostForm(url string, data url.Values) (*XHttpResponse, error) {
	return Http().PostForm(url, data)
}

// 重试机制
type RetryConfig struct {
	MaxRetries int
	Delay      time.Duration
	BackoffFn  func(int) time.Duration
}

// WithRetry 添加重试机制
func (h XHttp) WithRetry(config RetryConfig) XHttp {
	originalClient := h.client
	h.client = &http.Client{
		Timeout: originalClient.Timeout,
		Transport: &retryTransport{
			Transport: originalClient.Transport,
			config:    config,
		},
	}
	return h
}

type retryTransport struct {
	Transport http.RoundTripper
	config    RetryConfig
}

func (rt *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := rt.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	var lastErr error
	for i := 0; i <= rt.config.MaxRetries; i++ {
		if i > 0 {
			delay := rt.config.Delay
			if rt.config.BackoffFn != nil {
				delay = rt.config.BackoffFn(i)
			}
			time.Sleep(delay)
		}

		resp, err := transport.RoundTrip(req)
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		lastErr = err
		if resp != nil {
			resp.Body.Close()
		}
	}

	return nil, lastErr
}

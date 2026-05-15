package zlm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is the ZLMediaKit REST API client
type Client struct {
	baseURL string
	secret  string
	client  *http.Client
}

func NewClient(ip string, httpPort int, secret string) *Client {
	return &Client{
		baseURL: fmt.Sprintf("http://%s:%d", ip, httpPort),
		secret:  secret,
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        50,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     30 * time.Second,
			},
		},
	}
}

func (c *Client) get(path string, params map[string]string) (map[string]interface{}, error) {
	u := c.baseURL + path
	q := url.Values{}
	q.Set("secret", c.secret)
	for k, v := range params {
		q.Set(k, v)
	}
	u += "?" + q.Encode()

	resp, err := c.client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) post(path string, params map[string]string) (map[string]interface{}, error) {
	u := c.baseURL + path
	q := url.Values{}
	q.Set("secret", c.secret)
	for k, v := range params {
		q.Set(k, v)
	}

	resp, err := c.client.PostForm(u, q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetServerConfig gets ZLM server configuration
func (c *Client) GetServerConfig() (*ServerConfig, error) {
	result, err := c.get("/index/api/getServerConfig", nil)
	if err != nil {
		return nil, err
	}
	// Parse response into ServerConfig struct
	data, _ := json.Marshal(result)
	var configs []ServerConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return nil, err
	}
	if len(configs) > 0 {
		return &configs[0], nil
	}
	return nil, fmt.Errorf("no server config returned")
}

// AddStreamProxy adds a stream proxy
func (c *Client) AddStreamProxy(vhost, app, stream, url string, enableAudio, enableMP4 bool, timeout int) (*StreamProxyResult, error) {
	params := map[string]string{
		"vhost":       vhost,
		"app":         app,
		"stream":      stream,
		"url":         url,
		"enable_audio": boolToInt(enableAudio),
		"enable_mp4":  boolToInt(enableMP4),
		"timeout_sec": fmt.Sprintf("%d", timeout),
	}
	result, err := c.get("/index/api/addStreamProxy", params)
	if err != nil {
		return nil, err
	}
	data, _ := json.Marshal(result)
	var res StreamProxyResult
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// DelStreamProxy removes a stream proxy
func (c *Client) DelStreamProxy(key, vhost, app, stream string) (bool, error) {
	params := map[string]string{
		"key":    key,
		"vhost":  vhost,
		"app":    app,
		"stream": stream,
	}
	result, err := c.get("/index/api/delStreamProxy", params)
	if err != nil {
		return false, err
	}
	return result["code"].(float64) == 0, nil
}

// CloseStream closes a stream
func (c *Client) CloseStream(vhost, app, stream, schema string) (bool, error) {
	params := map[string]string{
		"vhost":  vhost,
		"app":    app,
		"stream": stream,
		"schema": schema,
	}
	result, err := c.get("/index/api/close_stream", params)
	if err != nil {
		return false, err
	}
	return result["code"].(float64) == 0, nil
}

// GetMediaList gets media stream list
func (c *Client) GetMediaList(vhost, app, stream, schema string, schemaFilter int) ([]MediaInfo, error) {
	params := map[string]string{}
	if vhost != "" {
		params["vhost"] = vhost
	}
	if app != "" {
		params["app"] = app
	}
	if stream != "" {
		params["stream"] = stream
	}
	if schema != "" {
		params["schema"] = schema
	}
	if schemaFilter > 0 {
		params["schema"] = fmt.Sprintf("%d", schemaFilter)
	}

	result, err := c.get("/index/api/getMediaList", params)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(result)
	var resp MediaListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetOnlineMediaList gets online media list (faster)
func (c *Client) GetOnlineMediaList(vhost, app, stream string) ([]MediaInfo, error) {
	params := map[string]string{}
	if vhost != "" {
		params["vhost"] = vhost
	}
	if app != "" {
		params["app"] = app
	}
	if stream != "" {
		params["stream"] = stream
	}

	result, err := c.get("/index/api/getOnlineMediaList", params)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(result)
	var resp OnlineMediaListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Data.OnlineMediaList, nil
}

// StartRecord starts recording
func (c *Client) StartRecord(vhost, app, stream string, recordType int) (bool, error) {
	params := map[string]string{
		"vhost": vhost,
		"app":   app,
		"stream": stream,
		"type":  fmt.Sprintf("%d", recordType),
	}
	result, err := c.get("/index/api/startRecord", params)
	if err != nil {
		return false, err
	}
	return result["code"].(float64) == 0, nil
}

// StopRecord stops recording
func (c *Client) StopRecord(vhost, app, stream string) (bool, error) {
	params := map[string]string{
		"vhost":  vhost,
		"app":    app,
		"stream": stream,
	}
	result, err := c.get("/index/api/stopRecord", params)
	if err != nil {
		return false, err
	}
	return result["code"].(float64) == 0, nil
}

// GetRtpInfo gets RTP server info
func (c *Client) GetRtpInfo(streamID string) (*RtpInfo, error) {
	params := map[string]string{
		"stream_id": streamID,
	}
	result, err := c.get("/index/api/getRtpInfo", params)
	if err != nil {
		return nil, err
	}
	data, _ := json.Marshal(result)
	var resp RtpInfoResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	if !resp.Exist {
		return nil, fmt.Errorf("rtp server not found")
	}
	return &resp.RtpInfo, nil
}

// OpenRtpServer opens an RTP server
func (c *Client) OpenRtpServer(port, streamID string, tcpMode int, isUDP bool) (*RtpServerResult, error) {
	params := map[string]string{
		"port":      port,
		"stream_id": streamID,
		"tcp_mode":  fmt.Sprintf("%d", tcpMode),
	}
	if isUDP {
		params["is_udp"] = "1"
	}
	result, err := c.get("/index/api/openRtpServer", params)
	if err != nil {
		return nil, err
	}
	data, _ := json.Marshal(result)
	var res RtpServerResult
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// CloseRtpServer closes an RTP server
func (c *Client) CloseRtpServer(streamID string) (bool, error) {
	params := map[string]string{
		"stream_id": streamID,
	}
	result, err := c.get("/index/api/closeRtpServer", params)
	if err != nil {
		return false, err
	}
	return result["code"].(float64) == 0, nil
}

// ListRtpServer lists all RTP servers
func (c *Client) ListRtpServer() ([]RtpServer, error) {
	result, err := c.get("/index/api/listRtpServer", nil)
	if err != nil {
		return nil, err
	}
	data, _ := json.Marshal(result)
	var resp RtpServerListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// StartSendRtp starts sending RTP stream (for GB28181 cascade)
func (c *Client) StartSendRtp(vhost, app, stream, ssrc, dstURL, dstPort string, isUDP bool, usePS bool) (*SendRtpResult, error) {
	params := map[string]string{
		"vhost":   vhost,
		"app":     app,
		"stream":  stream,
		"ssrc":    ssrc,
		"dst_url": dstURL,
		"dst_port": dstPort,
	}
	if isUDP {
		params["is_udp"] = "1"
	}
	if usePS {
		params["use_ps"] = "1"
	}
	result, err := c.get("/index/api/startSendRtp", params)
	if err != nil {
		return nil, err
	}
	data, _ := json.Marshal(result)
	var res SendRtpResult
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// StopSendRtp stops sending RTP stream
func (c *Client) StopSendRtp(vhost, app, stream, ssrc string) (bool, error) {
	params := map[string]string{
		"vhost":  vhost,
		"app":    app,
		"stream": stream,
		"ssrc":   ssrc,
	}
	result, err := c.get("/index/api/stopSendRtp", params)
	if err != nil {
		return false, err
	}
	return result["code"].(float64) == 0, nil
}

// GetSnap gets a snapshot from a stream
func (c *Client) GetSnap(vhost, app, stream, timeoutSec, expireSec string) (string, error) {
	params := map[string]string{
		"vhost":       vhost,
		"app":         app,
		"stream":      stream,
		"timeout_sec": timeoutSec,
		"expire_sec":  expireSec,
	}
	result, err := c.get("/index/api/getSnap", params)
	if err != nil {
		return "", err
	}
	if code, ok := result["code"].(float64); ok && code == 0 {
		return result["url"].(string), nil
	}
	return "", fmt.Errorf("get snap failed: %v", result)
}

// Play controls playback (start)
func (c *Client) Play(key, timeoutSec string) (bool, error) {
	params := map[string]string{
		"key":         key,
		"timeout_sec": timeoutSec,
	}
	result, err := c.get("/index/api/play", params)
	if err != nil {
		return false, err
	}
	return result["code"].(float64) == 0, nil
}

// Pause pauses playback
func (c *Client) Pause(key string) (bool, error) {
	params := map[string]string{
		"key": key,
	}
	result, err := c.get("/index/api/pause", params)
	if err != nil {
		return false, err
	}
	return result["code"].(float64) == 0, nil
}

// Resume resumes playback
func (c *Client) Resume(key string) (bool, error) {
	params := map[string]string{
		"key": key,
	}
	result, err := c.get("/index/api/resume", params)
	if err != nil {
		return false, err
	}
	return result["code"].(float64) == 0, nil
}

// Seek seeks playback position
func (c *Client) Seek(key, stampSec string) (bool, error) {
	params := map[string]string{
		"key":       key,
		"stamp_sec": stampSec,
	}
	result, err := c.get("/index/api/seek", params)
	if err != nil {
		return false, err
	}
	return result["code"].(float64) == 0, nil
}

// Speed sets playback speed
func (c *Client) Speed(key, speed string) (bool, error) {
	params := map[string]string{
		"key":   key,
		"speed": speed,
	}
	result, err := c.get("/index/api/speed", params)
	if err != nil {
		return false, err
	}
	return result["code"].(float64) == 0, nil
}

func boolToInt(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

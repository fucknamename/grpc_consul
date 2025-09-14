package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 参考：https://www.jianshu.com/p/30d5138ceee6

// func GenHttpHeader() map[string][]string {
// 	// var Header = make(map[string][]string, 0)
// 	Header := make(map[string][]string)
// 	Header["User-Agent"] = []string{RandomUserAgent()}
// 	Header["Connection"] = []string{"keep-alive"}
// 	Header["Accept"] = []string{"*/*"}
// 	Header["Accept-Encoding"] = []string{"gzip, deflate"}
// 	//....
// 	return Header
// }

func HttpPost(
	apiurl,
	data string,
	connTimeoutMs,
	serveTimeoutMs int,
	proxyurl,
	contentType string,
	header map[string]string) ([]byte, error) {

	fName := "HttpPost"

	var client *http.Client
	if proxyurl != "" {
		client = genProxyClient(connTimeoutMs, serveTimeoutMs, proxyurl)
	} else {
		client = genNormalClient(connTimeoutMs, serveTimeoutMs)
	}

	// body := strings.NewReader(data)
	body := bytes.NewBuffer([]byte(data))

	reqest, _ := http.NewRequest("POST", apiurl, body)
	if len(header) > 0 {
		for k, v := range header {
			reqest.Header.Set(k, v)
		}
	}
	if len(contentType) > 0 {
		reqest.Header.Set("Content-Type", contentType)
	} else {
		reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	response, err := client.Do(reqest)
	if err != nil {
		return nil, fmt.Errorf("%s http failed, POST url:%s, reason:%s", fName, apiurl, err.Error())
	}
	defer response.Body.Close()

	var errMsg string
	if response.StatusCode != 200 {
		errMsg += fmt.Sprintf("http status %d ", response.StatusCode)
	}

	res_body, err := io.ReadAll(response.Body)
	if err != nil {
		errMsg += fmt.Sprintf("can't read response %s ", err.Error())
	}

	if len(errMsg) > 0 {
		errMsg += fmt.Sprintf("body :%s", string(res_body))
		err = errors.New(errMsg)
	}

	if err != nil {
		return res_body, fmt.Errorf("%s %s", fName, err.Error())
	}

	return res_body, nil
}

func genNormalClient(connTimeoutMs, serveTimeoutMs int) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Duration(connTimeoutMs)*time.Millisecond)
				if err != nil {
					return nil, err
				}
				err = c.SetDeadline(time.Now().Add(time.Duration(serveTimeoutMs) * time.Millisecond))
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		},
	}

	return client
}

func genProxyClient(connTimeoutMs, serveTimeoutMs int, proxyurl string) *http.Client {
	u, _ := url.Parse(proxyurl)

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Duration(connTimeoutMs)*time.Millisecond)
				if err != nil {
					return nil, err
				}
				err = c.SetDeadline(time.Now().Add(time.Duration(serveTimeoutMs) * time.Millisecond))
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			Proxy: http.ProxyURL(u),
		},
	}

	return client
}

func HttpGet(url string, connTimeoutMs, serveTimeoutMs int,
	headerPara map[string]string) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Duration(connTimeoutMs)*time.Millisecond)
				if err != nil {
					return nil, err
				}
				err = c.SetDeadline(time.Now().Add(time.Duration(serveTimeoutMs) * time.Millisecond))
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		},
	}

	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if len(headerPara) > 0 {
		for ind, val := range headerPara {
			reqest.Header.Add(ind, val)
		}
	}
	response, err := client.Do(reqest)
	if err != nil {
		return nil, fmt.Errorf("http failed, GET url:%s, reason:%s", url, err.Error())
	}
	defer response.Body.Close()

	var errMsg string
	if response.StatusCode != 200 {
		errMsg += fmt.Sprintf("http status %d ", response.StatusCode)
	}

	res_body, err := io.ReadAll(response.Body)
	if err != nil {
		errMsg += fmt.Sprintf("can't read response %s ", err.Error())
	}

	if len(errMsg) > 0 {
		errMsg += fmt.Sprintf("url :%s ", url)
		err = errors.New(errMsg)
	}

	return res_body, err
}

func HttpHParaPost(url string, headPara map[string]string, data string, connTimeoutMs int, serveTimeoutMs int) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Duration(connTimeoutMs)*time.Millisecond)
				if err != nil {
					return nil, err
				}
				err = c.SetDeadline(time.Now().Add(time.Duration(serveTimeoutMs) * time.Millisecond))
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		},
	}

	body := strings.NewReader(data)
	reqest, _ := http.NewRequest("POST", url, body)
	// 请求头检查，大于0就设置
	if len(headPara) > 0 {
		for k, v := range headPara {
			reqest.Header.Add(k, v)
		}
	}
	if _, ok := headPara["Content-Type"]; !ok {
		reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	response, err := client.Do(reqest)
	if err != nil {
		return nil, fmt.Errorf("http failed, POST url:%s, reason:%s", url, err.Error())
	}
	defer response.Body.Close()

	var errMsg string
	if response.StatusCode != 200 {
		errMsg += fmt.Sprintf("http status %d ", response.StatusCode)
	}

	res_body, err := io.ReadAll(response.Body)
	if err != nil {
		errMsg += fmt.Sprintf("can't read response %s ", err.Error())
	}

	if len(errMsg) > 0 {
		errMsg += fmt.Sprintf("url :%s ", url)
		err = errors.New(errMsg)
	}

	return res_body, err
}

func HttpFormPost(payurl string, postdata map[string]interface{}, contentType string) ([]byte, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	for k, v := range postdata {
		_ = writer.WriteField(k, fmt.Sprintf("%v", v))
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", payurl, payload)
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	} else {
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

// https://stackoverflow.com/questions/27234861/correct-way-of-getting-clients-ip-addresses-from-http-request?noredirect=1&lq=1
// Get the IP address of the server's connected user.
func GetUserIP(httpWriter http.ResponseWriter, httpServer *http.Request) {
	var userIP string
	if len(httpServer.Header.Get("CF-Connecting-IP")) > 1 {
		userIP = httpServer.Header.Get("CF-Connecting-IP")
		fmt.Println(net.ParseIP(userIP))
	} else if len(httpServer.Header.Get("X-Forwarded-For")) > 1 {
		userIP = httpServer.Header.Get("X-Forwarded-For")
		fmt.Println(net.ParseIP(userIP))
	} else if len(httpServer.Header.Get("X-Real-IP")) > 1 {
		userIP = httpServer.Header.Get("X-Real-IP")
		fmt.Println(net.ParseIP(userIP))
	} else {
		userIP = httpServer.RemoteAddr
		if strings.Contains(userIP, ":") {
			fmt.Println(net.ParseIP(strings.Split(userIP, ":")[0]))
		} else {
			fmt.Println(net.ParseIP(userIP))
		}
	}
}

// func GetClientIp(r *http.Request) string {
// 	clientIp := ""
// 	realIps := r.Header.Get("X-Forwarded-For")
// 	if realIps != "" && len(realIps) != 0 && !strings.EqualFold("unknown", realIps) {
// 		ipArray := strings.Split(realIps, ",")
// 		clientIp = ipArray[0]
// 	}
// 	if clientIp == "" {
// 		clientIp = r.Header.Get("Proxy-Client-IP")
// 	}
// 	if clientIp == "" {
// 		clientIp = r.Header.Get("WL-Proxy-Client-IP")
// 	}
// 	if clientIp == "" {
// 		clientIp = r.Header.Get("HTTP_CLIENT_IP")
// 	}
// 	if clientIp == "" {
// 		clientIp = r.Header.Get("HTTP_X_FORWARDED_FOR")
// 	}
// 	if clientIp == "" {
// 		clientIp = r.Header.Get("X-Real-IP")
// 	}
// 	if clientIp == "" {
// 		clientIp = getRemoteIp(r)
// 	}
// 	return clientIp
// }

// func getRemoteIp(r *http.Request) string {
// 	array, _ := MatchString(`(.+):(\d+)`, r.RemoteAddr)
// 	if len(array) > 1 {
// 		return strings.Trim(array[1], "[]")
// 	}
// 	return r.RemoteAddr
// }

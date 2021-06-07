// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package k8s

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/clusterdialer"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/modules/scheduler/conf"
	"github.com/erda-project/erda/pkg/customhttp"
)

const (
	Error   = '1'
	Ping    = '2'
	Pong    = '3'
	Input   = '4'
	Output  = '5'
	GetSize = '6'
	Size    = '7'
	SetSize = '8'
)

type Winsize struct {
	Rows uint16 // ws_row: Number of rows (in cells)
	Cols uint16 // ws_col: Number of columns (in cells)
	X    uint16 // ws_xpixel: Width in pixels
	Y    uint16 // ws_ypixel: Height in pixels
}

var passRegexp = regexp.MustCompile(`(?i)(?:pass|secret)[^\s\0]*=([^\s\0]+)`)

func hidePassEnv(b []byte) []byte {
	return passRegexp.ReplaceAllFunc(b, func(a []byte) []byte {
		if i := bytes.IndexByte(a, '='); i != -1 {
			for i += 2; i < len(a); i++ { // keep first
				a[i] = '*'
			}
		}
		return a
	})
}

func (k *Kubernetes) Terminal(namespace, podName, containerName string, upperConn *websocket.Conn) {
	f := func(cols, rows uint16) (*websocket.Conn, error) {
		path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/exec", namespace, podName)
		s := `stty cols %d rows %d; s=/bin/sh; if [ -f /bin/bash ]; then s=/bin/bash; fi; `
		if conf.TerminalSecurity() {
			s += "if [ `id -un` != dice ]; then su -l dice -s $s; exit $?; fi; "
		}
		s += "$s"
		cmd := url.QueryEscape(fmt.Sprintf(s, cols, rows))
		query := "command=sh&command=-c&command=" + cmd + "&container=" + containerName + "&stdin=1&stdout=1&tty=true"

		header := http.Header{}
		dialer := websocket.DefaultDialer
		dialer.TLSClientConfig = &tls.Config{}

		execURL := url.URL{
			Scheme:   "wss",
			Host:     k.manageConfig.Address,
			Path:     path,
			RawQuery: query,
		}

		if !(k.manageConfig.Type == apistructs.ManageInet || k.manageConfig.Type == "") {
			dialer.TLSClientConfig.InsecureSkipVerify = k.manageConfig.Insecure
			if !k.manageConfig.Insecure {
				caBytes, err := base64.StdEncoding.DecodeString(k.manageConfig.CaData)
				if err != nil {
					logrus.Errorf("ca bytes load error: %v", err)
					return nil, err
				}
				pool := x509.NewCertPool()
				pool.AppendCertsFromPEM(caBytes)
				dialer.TLSClientConfig.RootCAs = pool
			}
		}

		switch k.manageConfig.Type {
		case apistructs.ManageProxy:
			dialer.NetDialContext = clusterdialer.DialContext(k.clusterName)
			header.Add("Authorization", "Bearer "+k.manageConfig.Token)
		case apistructs.ManageToken:
			header.Add("Authorization", "Bearer "+k.manageConfig.Token)
		case apistructs.ManageCert:
			certData, err := base64.StdEncoding.DecodeString(k.manageConfig.CertData)
			if err != nil {
				logrus.Errorf("decode cert data error: %v", err)
				return nil, err
			}
			keyData, err := base64.StdEncoding.DecodeString(k.manageConfig.KeyData)
			if err != nil {
				logrus.Errorf("decode key data error: %v", err)
				return nil, err
			}
			pair, err := tls.X509KeyPair(certData, keyData)
			if err != nil {
				logrus.Errorf("load X509Key pair error: %v", err)
				return nil, err
			}
			dialer.TLSClientConfig.Certificates = []tls.Certificate{pair}
		case apistructs.ManageInet:
			fallthrough
		default:
			req, err := customhttp.NewRequest("GET", k.manageConfig.Address, nil)
			if err != nil {
				logrus.Errorf("failed to customhttp.NewRequest: %v", err)
				return nil, err
			}
			req.Header.Add("X-Portal-Websocket", "on")
			if req.Header.Get("X-Portal-Host") != "" {
				req.Header.Add("Host", req.Header.Get("X-Portal-Host"))
			}

			execURL.Scheme = "ws"
			execURL.Host = req.URL.Host
			header = req.Header
		}

		conn, resp, err := dialer.Dial(execURL.String(), header)
		if err != nil {
			logrus.Errorf("failed to connect to %s: %v", execURL.String(), err)
			if resp == nil {
				return nil, err
			}
			logrus.Debugf("connect to %s request info: %+v", execURL.String(), resp.Request)
			respBody, _ := ioutil.ReadAll(resp.Body)
			logrus.Debugf("connect to %s response body: %s", execURL.String(), string(respBody))
			return nil, err
		}
		return conn, nil
	}

	var conn *websocket.Conn
	var waitConn sync.WaitGroup
	waitConn.Add(1)
	var setsizeOnce sync.Once
	var wait sync.WaitGroup
	wait.Add(2)
	go func() {
		defer func() {
			wait.Done()
			if conn != nil {
				conn.Close()
			}
			upperConn.Close()
		}()
		waitConn.Wait()
		for {
			if conn == nil {
				return
			}
			tp, m, err := conn.ReadMessage()
			if err != nil {
				return
			}
			m[0] = Output
			if conf.TerminalSecurity() {
				m = hidePassEnv(m)
			}
			if err := upperConn.WriteMessage(tp, m); err != nil {
				return
			}
		}
	}()
	go func() {
		defer func() {
			wait.Done()
			if conn != nil {
				conn.Close()
			}
			upperConn.Close()
		}()
		for {
			tp, m, err := upperConn.ReadMessage()
			if err != nil {
				return
			}
			switch m[0] {
			case Input:
				if conn == nil {
					continue
				}
				m[0] = 0 // k8s 协议, stdin = 0
				if err := conn.WriteMessage(tp, m); err != nil {
					return
				}
			case SetSize:
				var err error
				setsizeOnce.Do(func() {
					var v Winsize
					err = json.Unmarshal(m[1:], &v)
					if err != nil {
						return
					}
					conn, err = f(v.Cols, v.Rows)
					waitConn.Done()
					if err != nil {
						logrus.Errorf("failed to connect k8s exec ws: %v", err)
						return
					}
				})
				if err != nil {
					return
				}
			default:
				continue
			}
		}
	}()
	wait.Wait()
}

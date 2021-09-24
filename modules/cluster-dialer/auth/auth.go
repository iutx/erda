// Copyright (c) 2021 Terminus, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	credentialpb "github.com/erda-project/erda-proto-go/core/services/authentication/credentials/accesskey/pb"
	"github.com/erda-project/erda/pkg/secret"
	"github.com/erda-project/erda/pkg/secret/validator"
)

type Option func(authorizer *Authorizer)

type Authorizer struct {
	Credential credentialpb.AccessKeyServiceServer
}

func New(opts ...Option) *Authorizer {
	a := Authorizer{}

	for _, opt := range opts {
		opt(&a)
	}

	return &a
}

func WithCredentialClient(credential credentialpb.AccessKeyServiceServer) Option {
	return func(a *Authorizer) {
		a.Credential = credential
	}
}

func (a *Authorizer) Authorizer(req *http.Request) (string, bool, error) {
	// inner proxy not need auth
	if req.URL.Path == "/clusterdialer" {
		return "proxy", true, nil
	}

	clusterKey := req.Header.Get("X-Erda-Cluster-Key")
	auth := req.Header.Get("Authorization")

	logrus.Debugf("get auth info: %s", auth)

	ak := getAccessKeyIDFromAuthorization(auth)
	if ak == "" {
		return "", false, fmt.Errorf("get accesskey error, access key id is empty")
	}

	logrus.Debugf("get access key id: %s", ak)

	akSkResp, err := a.Credential.QueryAccessKeys(context.Background(), &credentialpb.QueryAccessKeysRequest{
		AccessKey: ak,
	})
	if err != nil {
		return "", false, fmt.Errorf("get key pair error: %v", err)
	}

	if akSkResp.Total == 0 {
		return "", false, fmt.Errorf("aksk pair not found, access key: %s", ak)
	} else if akSkResp.Total > 1 {
		return "", false, fmt.Errorf("aksk get error, duplidate result, access key: %s", ak)
	}

	hv := validator.NewHMACValidator(secret.AkSkPair{
		AccessKeyID: ak,
		SecretKey:   akSkResp.Data[0].SecretKey,
	}, validator.WithMaxExpireInterval(math.MaxInt64))

	res := hv.Verify(req)

	if !res.Ok {
		err = fmt.Errorf("validator error, message: %v", res.Message)
		logrus.Error(err)
		return "", false, err
	}

	return clusterKey, res.Ok, nil
}

func getAccessKeyIDFromAuthorization(authorization string) string {
	keys := strings.Split(authorization, "&")
	for _, key := range keys {
		kv := strings.Split(key, "=")
		if len(kv) != 2 {
			continue
		}
		if kv[0] == "X-Erda-Ak" {
			return kv[1]
		}
	}
	return ""
}

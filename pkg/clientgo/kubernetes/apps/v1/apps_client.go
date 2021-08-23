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

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes/scheme"
	appsv1client "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/rest"

	"github.com/erda-project/erda/pkg/clientgo/restclient"
)

// NewAppClient creates a new AppsV1Client for the given addr.
func NewAppClient(addr string) (*appsv1client.AppsV1Client, error) {

	var (
		client rest.Interface
		err    error
		config *rest.Config
	)
	if addr != "" {
		config = restclient.GetDefaultConfig("")
		config.GroupVersion = &appsv1.SchemeGroupVersion
		client, err = restclient.NewInetRESTClient(addr, config)
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		config.APIPath = "/apis"
		config.GroupVersion = &appsv1.SchemeGroupVersion
		config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
		rest.SetKubernetesDefaults(config)
		client, err = rest.RESTClientFor(config)
	}
	if err != nil {
		return nil, err
	}
	return appsv1client.New(client), nil
}

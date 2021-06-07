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

package manager

import (
	"github.com/erda-project/erda/pkg/jsonstore"
	"reflect"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/bundle"
	"github.com/erda-project/erda/pkg/clientgo"
)

// Interface .
type Interface interface {
	GetClusterClientSet(clusterName string) (*clientgo.ClientSet, error)
	GetCluster(clusterName string) (*apistructs.ClusterInfo, error)
	CreateCluster(cr *apistructs.ClusterCreateRequest) error
	UpdateCluster(cr *apistructs.ClusterUpdateRequest) error
	DeleteCluster(clusterName string) error
}

type define struct{}

func (d *define) Services() []string { return []string{"cluster-manager"} }
func (d *define) Types() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Interface)(nil)).Elem(),
	}
}
func (d *define) Description() string { return "cluster-manager" }
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

type provider struct {
	js     jsonstore.JsonStore
	bundle *bundle.Bundle
}

func (p *provider) Init(ctx servicehub.Context) error {
	js, err := jsonstore.New()
	if err != nil {
		return err
	}

	p.bundle = bundle.New(bundle.WithCMDB())
	p.js = js

	return nil
}

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	return p
}

func init() {
	servicehub.RegisterProvider("cluster-manager", &define{})
}

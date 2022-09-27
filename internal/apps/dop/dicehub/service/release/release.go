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

package release

import (
	"github.com/erda-project/erda/bundle"
	"github.com/erda-project/erda/internal/apps/dop/dicehub/dbclient"
	imagedb "github.com/erda-project/erda/internal/apps/dop/dicehub/image/db"
	"github.com/pkg/errors"
)

const (
	// AliYunRegistry 阿里云registry前缀
	AliYunRegistry = "registry.cn-hangzhou.aliyuncs.com"
)

// Release Release操作封装
type Release struct {
	db      *dbclient.DBClient
	bdl     *bundle.Bundle
	imageDB *imagedb.ImageConfigDB
}

// Option 定义 Release 对象的配置选项
type Option func(*Release)

// New 新建 Release 实例，操作 Release 资源
func New(options ...Option) *Release {
	app := &Release{}
	for _, op := range options {
		op(app)
	}
	return app
}

// WithDBClient 配置 db client
func WithDBClient(db *dbclient.DBClient) Option {
	return func(a *Release) {
		a.db = db
	}
}

// WithImageDBClient 配置 db client
func WithImageDBClient(db *imagedb.ImageConfigDB) Option {
	return func(a *Release) {
		a.imageDB = db
	}
}

// WithBundle 配置 bundle
func WithBundle(bdl *bundle.Bundle) Option {
	return func(a *Release) {
		a.bdl = bdl
	}
}

// GetDiceYAML 获取dice.yml内容
func (r *Release) GetDiceYAML(orgID int64, releaseID string) (string, error) {
	release, err := r.db.GetRelease(releaseID)
	if err != nil {
		return "", err
	}
	if orgID != 0 && release.OrgID != orgID { // when calling internally，orgID is 0
		return "", errors.Errorf("release not found")
	}

	return release.Dice, nil
}

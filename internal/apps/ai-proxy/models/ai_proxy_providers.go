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

// Code generated by erda-cli. DO NOT EDIT.
// Source: bin/erda-cli gorm gen -f .erda/ai-proxy/migrations/ai-proxy/20230808-providers.sql -o internal/apps/ai-proxy/models

package models

import (
	"time"

	"github.com/erda-project/erda-infra/providers/mysql/v2/plugins/fields"
)

// AIProxyProviders is the table ai_proxy_providers
type AIProxyProviders struct {
	ID          fields.UUID      `gorm:"column:id;type:char(36)" json:"iD" yaml:"iD"`
	CreatedAt   time.Time        `gorm:"column:created_at;type:datetime" json:"createdAt" yaml:"createdAt"`
	UpdatedAt   time.Time        `gorm:"column:updated_at;type:datetime" json:"updatedAt" yaml:"updatedAt"`
	DeletedAt   fields.DeletedAt `gorm:"column:deleted_at;type:datetime" json:"deletedAt" yaml:"deletedAt"`
	Name        string           `gorm:"column:name;type:varchar(128)" json:"name" yaml:"name"`
	InstanceID  string           `gorm:"column:instance_id;type:varchar(512)" json:"instanceID" yaml:"instanceID"`
	Host        string           `gorm:"column:host;type:varchar(128)" json:"host" yaml:"host"`
	Scheme      string           `gorm:"column:scheme;type:varchar(16)" json:"scheme" yaml:"scheme"`
	Description string           `gorm:"column:description;type:varchar(1024)" json:"description" yaml:"description"`
	DocSite     string           `gorm:"column:doc_site;type:varchar(512)" json:"docSite" yaml:"docSite"`
	AesKey      string           `gorm:"column:aes_key;type:char(16)" json:"aesKey" yaml:"aesKey"`
	APIKey      string           `gorm:"column:api_key;type:varchar(512)" json:"aPIKey" yaml:"aPIKey"`
	Metadata    string           `gorm:"column:metadata;type:longtext" json:"metadata" yaml:"metadata"`
}

// TableName returns the table name ai_proxy_providers
func (*AIProxyProviders) TableName() string { return "ai_proxy_providers" }

type AIProxyProvidersList []*AIProxyProviders

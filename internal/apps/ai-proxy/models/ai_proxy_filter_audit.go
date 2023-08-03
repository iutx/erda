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
// Source: .erda/ai-proxy/migrations/ai-proxy/20230425-filter-audit.sql

package models

import (
	"time"

	"github.com/erda-project/erda-infra/providers/mysql/v2/plugins/fields"
)

// AIProxyFilterAudit is the table ai_proxy_filter_audit
type AIProxyFilterAudit struct {
	ID                  fields.UUID      `gorm:"column:id;type:char(36)" json:"iD" yaml:"iD"`
	CreatedAt           time.Time        `gorm:"column:created_at;type:datetime" json:"createdAt" yaml:"createdAt"`
	UpdatedAt           time.Time        `gorm:"column:updated_at;type:datetime" json:"updatedAt" yaml:"updatedAt"`
	DeletedAt           fields.DeletedAt `gorm:"column:deleted_at;type:datetime" json:"deletedAt" yaml:"deletedAt"`
	APIKeySHA256        string           `gorm:"column:api_key_sha256;type:char(64)" json:"aPIKeySHA256" yaml:"aPIKeySHA256"`
	Username            string           `gorm:"column:username;type:varchar(128)" json:"username" yaml:"username"`
	PhoneNumber         string           `gorm:"column:phone_number;type:varchar(32)" json:"phoneNumber" yaml:"phoneNumber"`
	JobNumber           string           `gorm:"column:job_number;type:varchar(32)" json:"jobNumber" yaml:"jobNumber"`
	Email               string           `gorm:"column:email;type:varchar(64)" json:"email" yaml:"email"`
	DingtalkStaffID     string           `gorm:"column:dingtalk_staff_id;type:varchar(64)" json:"dingtalkStaffID" yaml:"dingtalkStaffID"`
	SessionID           string           `gorm:"column:session_id;type:varchar(64)" json:"sessionID" yaml:"sessionID"`
	ChatType            string           `gorm:"column:chat_type;type:varchar(32)" json:"chatType" yaml:"chatType"`
	ChatTitle           string           `gorm:"column:chat_title;type:varchar(128)" json:"chatTitle" yaml:"chatTitle"`
	ChatID              string           `gorm:"column:chat_id;type:varchar(64)" json:"chatID" yaml:"chatID"`
	Source              string           `gorm:"column:source;type:varchar(128)" json:"source" yaml:"source"`
	ProviderName        string           `gorm:"column:provider_name;type:varchar(128)" json:"providerName" yaml:"providerName"`
	ProviderInstanceID  string           `gorm:"column:provider_instance_id;type:varchar(512)" json:"providerInstanceID" yaml:"providerInstanceID"`
	Model               string           `gorm:"column:model;type:varchar(128)" json:"model" yaml:"model"`
	OperationID         string           `gorm:"column:operation_id;type:varchar(128)" json:"operationID" yaml:"operationID"`
	Prompt              string           `gorm:"column:prompt;type:mediumtext" json:"prompt" yaml:"prompt"`
	Completion          string           `gorm:"column:completion;type:longtext" json:"completion" yaml:"completion"`
	Metadata            string           `gorm:"column:metadata;type:longtext" json:"metadata" yaml:"metadata"`
	XRequestID          string           `gorm:"column:x_request_id;type:varchar(64)" json:"xRequestID" yaml:"xRequestID"`
	RequestAt           time.Time        `gorm:"column:request_at;type:datetime" json:"requestAt" yaml:"requestAt"`
	ResponseAt          time.Time        `gorm:"column:response_at;type:datetime" json:"responseAt" yaml:"responseAt"`
	RequestContentType  string           `gorm:"column:request_content_type;type:varchar(32)" json:"requestContentType" yaml:"requestContentType"`
	RequestBody         string           `gorm:"column:request_body;type:longtext" json:"requestBody" yaml:"requestBody"`
	ResponseContentType string           `gorm:"column:response_content_type;type:varchar(32)" json:"responseContentType" yaml:"responseContentType"`
	ResponseBody        string           `gorm:"column:response_body;type:longtext" json:"responseBody" yaml:"responseBody"`
	UserAgent           string           `gorm:"column:user_agent;type:varchar(128)" json:"userAgent" yaml:"userAgent"`
	Server              string           `gorm:"column:server;type:varchar(32)" json:"server" yaml:"server"`
	Status              string           `gorm:"column:status;type:varchar(32)" json:"status" yaml:"status"`
	StatusCode          int64            `gorm:"column:status_code;type:int(11)" json:"statusCode" yaml:"statusCode"`
}

// TableName returns the table name ai_proxy_filter_audit
func (*AIProxyFilterAudit) TableName() string { return "ai_proxy_filter_audit" }

type AIProxyFilterAuditList []*AIProxyFilterAudit

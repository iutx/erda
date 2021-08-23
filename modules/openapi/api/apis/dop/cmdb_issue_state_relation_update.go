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

package dop

import (
	"net/http"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/modules/openapi/api/apis"
)

var CMDB_ISSUE_STATE_RELATION_UPDATE = apis.ApiSpec{
	Path:         "/api/issues/actions/update-state-relation",
	BackendPath:  "/api/issues/actions/update-state-relation",
	Host:         "dop.marathon.l4lb.thisdcos.directory:9527",
	Scheme:       "http",
	Method:       http.MethodPut,
	CheckLogin:   true,
	CheckToken:   true,
	RequestType:  apistructs.IssueStateUpdateRequest{},
	ResponseType: apistructs.IssueStateRelationGetResponse{},
	IsOpenAPI:    true,
	Doc:          "summary: 更新 工作流 ",
}

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

package core_services

import (
	"net/http"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/modules/openapi/api/apis"
)

var CMDB_ORG_RESOURCE_GET = apis.ApiSpec{
	Path:         "/api/orgs/actions/fetch-resources",
	BackendPath:  "/api/orgs/actions/fetch-resources",
	Host:         "core-services.marathon.l4lb.thisdcos.directory:9526",
	Scheme:       "http",
	Method:       http.MethodGet,
	CheckLogin:   true,
	CheckToken:   true,
	ResponseType: apistructs.OrgResourceInfo{},
	IsOpenAPI:    true,
	Doc:          "summary: 获取企业资源使用",
}

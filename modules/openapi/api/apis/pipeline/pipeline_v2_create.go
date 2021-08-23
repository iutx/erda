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

package pipeline

import (
	"net/http"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/modules/openapi/api/apis"
)

var PIPELINE_V2_CREATE = apis.ApiSpec{
	Path:         "/api/v2/pipelines",
	BackendPath:  "/api/v2/pipelines",
	Host:         "pipeline.marathon.l4lb.thisdcos.directory:3081",
	Scheme:       "http",
	Method:       http.MethodPost,
	IsOpenAPI:    true,
	CheckLogin:   true,
	CheckToken:   true,
	RequestType:  apistructs.PipelineCreateRequestV2{},
	ResponseType: apistructs.PipelineCreateResponse{},
	Doc:          "summary: v2 版本创建 pipeline",
}

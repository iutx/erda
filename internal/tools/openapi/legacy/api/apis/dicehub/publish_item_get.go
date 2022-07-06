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

package dicehub

import "github.com/erda-project/erda/internal/tools/openapi/legacy/api/apis"

var PUBLISH_ITEM_GET = apis.ApiSpec{
	Path:          "/api/publish-items/<id>",
	BackendPath:   "/api/publish-items/<id>",
	Host:          "erda-server.marathon.l4lb.thisdcos.directory:9095",
	Scheme:        "http",
	Method:        "GET",
	TryCheckLogin: true,
	Doc:           "summary: 查询发布内容",
}

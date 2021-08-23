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

package config_manage

import "github.com/erda-project/erda/modules/openapi/api/apis"

var CONFIG_MANAGE_DEPLOY_CONFIG_GET = apis.ApiSpec{
	Path:        "/api/configmanage/configs/deployment",
	BackendPath: "/api/config/deployment",
	Host:        "dop.marathon.l4lb.thisdcos.directory:9527",
	Scheme:      "http",
	Method:      "GET",
	CheckLogin:  true,
	Doc:         "summary: 获取指定命名空间的所有部署配置",
}

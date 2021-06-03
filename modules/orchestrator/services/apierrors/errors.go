// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

// Package apierrors 定义了错误列表
package apierrors

import (
	"github.com/erda-project/erda/pkg/http/httpserver/errorresp"
)

// runtime errors
var (
	ErrCreateRuntime   = err("ErrCreateRuntime", "创建应用实例失败")
	ErrDeleteRuntime   = err("ErrDeleteRuntime", "删除应用实例失败")
	ErrDeployRuntime   = err("ErrDeployRuntime", "部署失败")
	ErrRollbackRuntime = err("ErrRollbackRuntime", "回滚失败")
	ErrListRuntime     = err("ErrListRuntime", "查询应用实例列表失败")
	ErrGetRuntime      = err("ErrGetRuntime", "查询应用实例失败")
	ErrUpdateRuntime   = err("ErrUpdateRuntime", "更新应用实例失败")
	ErrReferRuntime    = err("ErrReferRuntime", "查询应用实例引用集群失败")
	ErrKillPod         = err("ErrKillPod", "kill pod 失败")
)

var (
	ErrListInstance = err("ErrListInstance", "获取实例列表失败")
)

// deployment errors
var (
	ErrCancelDeployment     = err("ErrCancelDeployment", "取消部署失败")
	ErrListDeployment       = err("ErrListDeployment", "查询部署列表失败")
	ErrGetDeployment        = err("ErrGetDeployment", "查询部署失败")
	ErrApproveDeployment    = err("ErrApproveDeployment", "审批部署失败")
	ErrDeployStagesAddons   = err("ErrDeployStagesAddons", "部署addon失败")
	ErrDeployStagesServices = err("ErrDeployStagesServices", "部署service失败")
	ErrDeployStagesDomains  = err("ErrDeployStagesDomains", "部署domain失败")
)

// domain errors
var (
	ErrListDomain   = err("ErrListDomain", "查询域名列表失败")
	ErrUpdateDomain = err("ErrUpdateDomain", "更新域名失败")
)

var (
	ErrCreateAddon     = err("ErrCreateAddon", "创建 addon 失败")
	ErrUpdateAddon     = err("ErrUpdateAddon", "更新 addon 失败")
	ErrFetchAddon      = err("ErrFetchAddon", "获取 addon 详情失败")
	ErrDeleteAddon     = err("ErrDeleteAddon", "删除 addon 失败")
	ErrListAddon       = err("ErrListAddon", "获取 addon 列表失败")
	ErrListAddonMetris = err("ErrListAddonMetris", "获取 addon 监控失败")
)

var (
	ErrMigrationLog = err("ErrMigrationLog", "查询migration日志失败")
)

var (
	ErrOrgLog = err("ErrOrgLog", "查询容器日志失败")
)

var (
	ErrProjectResource = err("ErrProjectResource", "查询项目资源失败")
)
var (
	ErrClusterResource = err("ErrClusterResource", "查询集群资源失败")
)

var (
	ErrGetAppWorkspaceReleases = err("ErrGetAppWorkspaceReleases", "查询环境可部署制品失败")
)

var (
	ErrAddonYmlExport = err("ErrAddonYmlExport", "addonyml 导出")
	ErrAddonYmlImport = err("ErrAddonYmlImport", "addonyml 导入")
)

func err(template, defaultValue string) *errorresp.APIError {
	return errorresp.New(errorresp.WithTemplateMessage(template, defaultValue))
}

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

package configcenter

import (
	"github.com/erda-project/erda/modules/msp/instance/db"
	"github.com/erda-project/erda/modules/msp/resource/deploy/handlers"
	"github.com/erda-project/erda/modules/msp/resource/utils"
)

func (p *provider) IsMatch(tmc *db.Tmc) bool {
	return tmc.Engine == handlers.ResourceConfigCenter
}

func (p *provider) BuildTmcInstanceConfig(tmcInstance *db.Instance, serviceGroupDeployResult interface{},
	clusterConfig map[string]string, additionalConfig map[string]string) map[string]string {
	config := map[string]string{}
	options := map[string]string{}
	utils.JsonConvertObjToType(tmcInstance.Options, &options)
	if v, ok := options["NACOS_ADDRESS"]; ok {
		config["CONFIGCENTER_ADDRESS"] = v
	}
	return config
}

func (p *provider) DoApplyTmcInstanceTenant(req *handlers.ResourceDeployRequest, resourceInfo *handlers.ResourceInfo,
	tmcInstance *db.Instance, tenant *db.InstanceTenant, clusterConfig map[string]string) (map[string]string, error) {

	instanceOptions := map[string]string{}
	_ = utils.JsonConvertObjToType(tmcInstance.Options, &instanceOptions)

	tenantOptions := map[string]string{}
	_ = utils.JsonConvertObjToType(tenant.Options, &tenantOptions)

	addr := instanceOptions["NACOS_ADDRESS"]
	user := instanceOptions["NACOS_USER"]
	pwd := instanceOptions["NACOS_PASSWORD"]
	namespaceId := tenantOptions["NACOS_TENANT_ID"]
	groupName := tenantOptions["applicationName"]

	p.saveDefaultConfigToNacos(clusterConfig["DICE_CLUSTER_NAME"], addr, user, pwd, namespaceId, groupName)

	key, _ := p.TmcIniDb.GetMicroServiceEngineJumpKey(tmcInstance.Engine)
	console := map[string]string{
		"tenantGroup":             tenant.TenantGroup,
		"tenantId":                tenant.ID,
		"key":                     key,
		"CONFIGCENTER_TENANT_ID":  tenant.ID,
		"CONFIGCENTER_GROUP_NAME": groupName,
	}

	str, _ := utils.JsonConvertObjToString(console)
	config := map[string]string{
		"CONFIGCENTER_TENANT_ID":   tenant.ID,
		"CONFIGCENTER_TENANT_NAME": namespaceId,
		"CONFIGCENTER_GROUP_NAME":  groupName,
		"PUBLIC_HOST":              str,
	}

	return config, nil
}

func (p *provider) DeleteTenant(tenant *db.InstanceTenant, tmcInstance *db.Instance, clusterConfig map[string]string) error {

	instanceOptions := map[string]string{}
	utils.JsonConvertObjToType(tmcInstance.Options, &instanceOptions)

	tenantConfig := map[string]string{}
	utils.JsonConvertObjToType(tenant.Config, &tenantConfig)

	addr := instanceOptions["NACOS_ADDRESS"]
	user := instanceOptions["NACOS_USER"]
	password := instanceOptions["NACOS_PASSWORD"]
	namespaceId := tenantConfig["CONFIGCENTER_TENANT_NAME"]
	groupName := tenantConfig["CONFIGCENTER_GROUP_NAME"]

	p.deleteDefaultConfigFromNacos(clusterConfig["DICE_CLUSTER_NAME"], addr, user, password, namespaceId, groupName)

	return p.DefaultDeployHandler.DeleteTenant(tenant, tmcInstance, clusterConfig)
}

func (p *provider) saveDefaultConfigToNacos(clusterName string, addr string, user string, pwd string, namespaceId string, groupName string) {
	nacosClient := utils.NewNacosClient(clusterName, addr, user, pwd)
	nacosClient.SaveConfig(namespaceId, groupName, "application.yml", "##")
}

func (p *provider) deleteDefaultConfigFromNacos(clusterName, addr, user, pwd, namespaceId, groupName string) {
	nacosClient := utils.NewNacosClient(clusterName, addr, user, pwd)
	nacosClient.DeleteConfig(namespaceId, groupName)
}

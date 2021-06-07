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

package cluster

import (
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/modules/scheduler/impl/cluster/clusterutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	unitTestInet      = "inet://unit-test.terminus.io?ssl=on/kubernetes.default.svc.cluster.local"
	unitTestApiServer = "https://kuberentes.default.svc.cluster.local"
	unitTestToken     = "[unit-test] token value"
	unitTestCertData  = "[unit-test] cert data value"
	unitTestKeyData   = "[unit-test] key data value"
	unitTestCaData    = "[unit-test] ca data value"
)

func TestSetDefaultClusterConfig(t *testing.T) {
	c1 := ClusterInfo{
		ClusterName: "x1",
		Kind:        "MARATHON",
		Options: map[string]string{
			"ADDR":      "http://dcos.x1.cn/service/marathon",
			"ENABLETAG": "true",
		},
	}

	setDefaultClusterConfig(&c1)
	assert.Equal(t, "true", c1.Options["ENABLETAG"])
	assert.Equal(t, "true", c1.Options["ENABLE_ORG"])
	assert.Equal(t, "true", c1.Options["ENABLE_WORKSPACE"])

	c2 := ClusterInfo{
		ClusterName: "x2",
		Kind:        "MARATHON",
		Options: map[string]string{
			"ADDR": "http://dcos.x2.cn/service/marathon",
		},
	}

	setDefaultClusterConfig(&c2)
	assert.Equal(t, "true", c2.Options["ENABLETAG"])
	assert.Equal(t, "true", c2.Options["ENABLE_ORG"])
	assert.Equal(t, "true", c2.Options["ENABLE_WORKSPACE"])

	c3 := ClusterInfo{
		ClusterName: "x3",
		Kind:        "METRONOME",
		Options: map[string]string{
			"ADDR": "http://dcos.x3.cn/service/metronome",
		},
	}

	setDefaultClusterConfig(&c3)
	assert.Equal(t, "true", c3.Options["ENABLETAG"])
	assert.Equal(t, "true", c3.Options["ENABLE_ORG"])
	assert.Equal(t, "true", c3.Options["ENABLE_WORKSPACE"])
}

func TestPatchK8SConfig(t *testing.T) {
	cur1 := &ClusterInfo{
		Kind: clusterutil.JobKindK8S,
		Options: map[string]string{
			"ADDR":      "",
			"TYPE":      apistructs.ManageToken,
			"CERT_DATA": "illegal data",
		},
	}
	new1 := &apistructs.ClusterInfo{
		SchedConfig: &apistructs.ClusterSchedConfig{
			MasterURL:         unitTestInet,
			CPUSubscribeRatio: "10",
		},
		ManageConfig: &apistructs.ManageConfig{
			Type:     apistructs.ManageToken,
			Address:  unitTestApiServer,
			CaData:   unitTestCaData,
			CertData: unitTestCertData,
			KeyData:  unitTestKeyData,
			Token:    unitTestToken,
		},
	}

	err := patchK8SConfig(cur1, new1)
	assert.Nil(t, err)
	assert.Equal(t, unitTestApiServer, cur1.Options["ADDR"])
	assert.Equal(t, "", cur1.Options["CERT_DATA"])
	assert.Equal(t, "", cur1.Options["DEV_CPU_SUBSCRIBE_RATIO"])
}

func TestCreateK8sOptions(t *testing.T) {

	c := apistructs.ClusterInfo{
		SchedConfig: &apistructs.ClusterSchedConfig{
			MasterURL: unitTestInet,
		},
		ManageConfig: &apistructs.ManageConfig{
			Address:  unitTestApiServer,
			CaData:   unitTestCaData,
			CertData: unitTestCertData,
			KeyData:  unitTestKeyData,
			Token:    unitTestToken,
		},
	}

	options := createK8sOptions(c)
	assert.Equal(t, unitTestInet, options["ADDR"])
	assert.Equal(t, "", options["TOKEN"])

	c.ManageConfig.Type = apistructs.ManageToken
	options = createK8sOptions(c)
	assert.Equal(t, unitTestApiServer, options["ADDR"])
	assert.Equal(t, unitTestToken, options["TOKEN"])
	assert.Equal(t, unitTestCaData, options["CA_DATA"])
	assert.Equal(t, "", options["CERT_DATA"])
	assert.Equal(t, "", options["KEY_DATA"])

	c.ManageConfig.Type = apistructs.ManageProxy
	options = createK8sOptions(c)
	assert.Equal(t, unitTestApiServer, options["ADDR"])
	assert.Equal(t, unitTestToken, options["TOKEN"])
	assert.Equal(t, unitTestCaData, options["CA_DATA"])
	assert.Equal(t, "", options["CERT_DATA"])
	assert.Equal(t, "", options["KEY_DATA"])

	c.ManageConfig.Type = apistructs.ManageCert
	options = createK8sOptions(c)
	assert.Equal(t, unitTestApiServer, options["ADDR"])
	assert.Equal(t, "", options["TOKEN"])
	assert.Equal(t, unitTestCaData, options["CA_DATA"])
	assert.Equal(t, unitTestCertData, options["CERT_DATA"])
	assert.Equal(t, unitTestKeyData, options["KEY_DATA"])

	c.ManageConfig.Insecure = true
	options = createK8sOptions(c)
	assert.Equal(t, unitTestApiServer, options["ADDR"])
	assert.Equal(t, "", options["TOKEN"])
	assert.Equal(t, "", options["CA_DATA"])
	assert.Equal(t, unitTestCertData, options["CERT_DATA"])
	assert.Equal(t, unitTestKeyData, options["KEY_DATA"])
}

package manager

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/modules/scheduler/conf"
	"github.com/erda-project/erda/pkg/clientgo"
	"github.com/erda-project/erda/pkg/clientgo/restclient"
	"github.com/erda-project/erda/pkg/httputil"
)

const (
	clusterK8SSuffix = "-service-k8s"
)

// GetClusterClientSet Get cluster client-go clientSet
func (p *provider) GetClusterClientSet(clusterName string) (*clientgo.ClientSet, error) {
	ci, err := p.bundle.GetCluster(clusterName)
	if err != nil {
		return nil, err
	}

	if ci.ManageConfig == nil {
		return clientgo.New(ci.SchedConfig.MasterURL)
	}

	switch ci.ManageConfig.Type {
	case apistructs.ManageProxy:
		rc, err := restclient.GetDialerRestConfig(clusterName, ci.ManageConfig)
		if err != nil {
			return nil, err
		}
		return clientgo.NewClientSet(rc)
	case apistructs.ManageToken, apistructs.ManageCert:
		rc, err := restclient.GetRestConfig(ci.ManageConfig)
		if err != nil {
			return nil, err
		}
		return clientgo.NewClientSet(rc)
	case apistructs.ManageInet:
		fallthrough
	default:
		return clientgo.New(ci.SchedConfig.MasterURL)
	}
}

// GetCluster
func (p *provider) GetCluster(clusterName string) (*apistructs.ClusterInfo, error) {
	ci, err := p.bundle.GetCluster(clusterName)
	if err != nil {
		return nil, err
	}

	// If cpu subscribe ratio in every workspace is diff from etcd, update to db.
	if ci.SchedConfig.DevCPUSubscribeRatio == "" || ci.SchedConfig.TestCPUSubscribeRatio == "" ||
		ci.SchedConfig.StagingCPUSubscribeRatio == "" {

		var exeConfig conf.ExecutorConfig

		// Get executor config from etcd.
		key := conf.CLUSTERS_CONFIG_PATH + clusterName + clusterK8SSuffix
		if err = p.js.Get(context.Background(), key, &exeConfig); err != nil {
			return nil, fmt.Errorf("key: %s get error: %v", key, err)
		}

		// Cover cpu subscribe ratio.
		ci.SchedConfig.DevCPUSubscribeRatio = exeConfig.Options["DEV_CPU_SUBSCRIBE_RATIO"]
		ci.SchedConfig.TestCPUSubscribeRatio = exeConfig.Options["TEST_CPU_SUBSCRIBE_RATIO"]
		ci.SchedConfig.StagingCPUSubscribeRatio = exeConfig.Options["STAGING_CPU_SUBSCRIBE_RATIO"]

		if err = p.UpdateCluster(&apistructs.ClusterUpdateRequest{
			Name:            clusterName,
			OrgID:           ci.OrgID,
			SchedulerConfig: ci.SchedConfig,
		}); err != nil {
			return nil, err
		}
	}

	return ci, nil
}

// CreateCluster Create cluster with bundle.
func (p *provider) CreateCluster(cr *apistructs.ClusterCreateRequest) error {
	// Create cluster with bundle.
	// cluster-manager -> bundle -> cmdb -> eventBox -> scheduler (hook)
	//                               |                     |
	//                           *database*             *etcd*
	return p.bundle.CreateCluster(cr, http.Header{
		httputil.InternalHeader: []string{"cluster-manager"},
	})
}

// UpdateCluster Update cluster with bundle.
func (p *provider) UpdateCluster(cr *apistructs.ClusterUpdateRequest) error {
	if cr == nil {
		return errors.New("cluster update request is nil")
	}

	// Default header, must provider internal and org header.
	header := http.Header{
		httputil.InternalHeader: []string{"cluster-manager"},
		httputil.OrgHeader:      []string{strconv.Itoa(cr.OrgID)},
	}

	// If doesn't provider orgID, get first.
	if cr.OrgID == 0 {
		ci, err := p.bundle.GetCluster(cr.Name)
		if err != nil {
			return err
		}
		cr.OrgID = ci.OrgID
		header[httputil.OrgHeader] = []string{strconv.Itoa(ci.OrgID)}
	}

	// The process of update is same to create.
	return p.bundle.UpdateCluster(*cr, header)
}

// DeleteCluster Delete cluster with bundle.
func (p *provider) DeleteCluster(clusterName string) error {
	return p.bundle.DeleteCluster(clusterName, http.Header{
		httputil.InternalHeader: []string{"cluster-manager"},
	})
}

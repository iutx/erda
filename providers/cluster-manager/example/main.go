package main

import (
	"context"
	"fmt"
	"github.com/erda-project/erda/modules/scheduler/conf"
	"github.com/erda-project/erda/pkg/jsonstore"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda/apistructs"
	manager "github.com/erda-project/erda/providers/cluster-manager"
)

// define Represents the definition of provider and provides some information
type define struct{}

// Declare what services the provider provides
func (d *define) Services() []string { return []string{"example"} }

// Declare which services the provider depends on
func (d *define) Dependencies() []string { return []string{"cluster-manager"} }

// Describe information about this provider
func (d *define) Description() string { return "example" }

// Return a provider creator
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

type provider struct {
	Manager manager.Interface
}

const (
	clusterK8SSuffix    = "-service-k8s"
	clusterK8SJobSuffix = "-job-k8s"
)

func (p *provider) Init(ctx servicehub.Context) error {
	js, err := jsonstore.New()
	if err != nil {
		return err
	}
	var exeConfig conf.ExecutorConfig

	// Get executor config from etcd.
	key := conf.CLUSTERS_CONFIG_PATH + "terminus-dev" + clusterK8SSuffix
	if err = js.Get(context.Background(), key, &exeConfig); err != nil {
		return fmt.Errorf("key: %s get error: %v", key, err)
	}

	exeConfig.Options["ENABLETAG"] = "true"

	if err = js.Put(context.Background(), key, &exeConfig); err != nil {
		return err
	}

	return nil
}

func (p *provider) Usage() error {
	clusterName := "fake-cluster"
	manageConfig := &apistructs.ManageConfig{
		Type:    "token",
		Address: "https://127.0.0.1:6443",
		CaData:  "ca data",
		Token:   "token",
	}
	newManageConfig := &apistructs.ManageConfig{
		Type:     "cert",
		Address:  "https://127.0.0.1:6443",
		CaData:   "ca data",
		CertData: "cert data",
		KeyData:  "key data",
	}

	if err := p.Manager.CreateCluster(&apistructs.ClusterCreateRequest{
		OrgID:          2,
		Name:           clusterName,
		CloudVendor:    "alicloud-ecs",
		DisplayName:    clusterName,
		Description:    "Demo",
		Type:           "k8s",
		Logo:           " ",
		WildcardDomain: fmt.Sprintf("%s.demo.io", clusterName),
		SchedulerConfig: &apistructs.ClusterSchedConfig{
			MasterURL: "inet://ingress-nginx.kube-system.svc.cluster.local?direct=on&ssl=on/kubernetes.default.svc.cluster.local",
		},
		ManageConfig: manageConfig,
	}); err != nil {
		return err
	}

	clientSet, err := p.Manager.GetClusterClientSet(clusterName)
	if err != nil {
		return err
	}

	ns, err := clientSet.K8sClient.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	fmt.Println(len(ns.Items))

	if err = p.Manager.UpdateCluster(&apistructs.ClusterUpdateRequest{
		Name:         clusterName,
		ManageConfig: newManageConfig,
	}); err != nil {
		return err
	}

	ci, err := p.Manager.GetCluster(clusterName)
	if err != nil {
		return err
	}

	fmt.Println(ci.SchedConfig)

	if err = p.Manager.DeleteCluster(clusterName); err != nil {
		return err
	}

	return nil
}

func init() {
	servicehub.RegisterProvider("example", &define{})
}

func main() {
	hub := servicehub.New()
	hub.Run("example", "/Users/tianxiang.liu/work/go/erda/providers/cluster-manager/example/example.yaml", os.Args...)
}

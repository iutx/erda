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

package clusterinfo

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/erda-project/erda/modules/scheduler/events"
	"github.com/erda-project/erda/pkg/clientgo/kubernetes"
	"github.com/erda-project/erda/pkg/dlock"
	"github.com/erda-project/erda/pkg/jsonstore"
	"github.com/erda-project/erda/pkg/strutil"
)

const (
	// diceCMNamespace dice configmap namespace
	diceCMNamespace = "default"
	// clusterInfoConfigMapName cluster info configmap name
	clusterInfoConfigMapName = "dice-cluster-info"
	// addonsConfigMapName addon configmap name
	addonsConfigMapName = "dice-addons-info"
	// clusterInfoPrefix 是集群配置信息在 etcd 中的路径前缀
	clusterInfoPrefix = "/dice/scheduler/clusterinfo/"
	// dlockKeyPrefix 分布式锁前缀，每个集群一把锁
	dlockKeyPrefix = "/dice/scheduler/dlock/clusterinfo/"
	// loopSyncTimeout 一次同步的超时时间
	loopSyncTimeout = 10 * time.Minute
)

const (
	netportalURLPrefix  = "inet://"
	netportalURLKeyName = "NETPORTAL_URL"
	// DiceClusterName dice 集群名
	DiceClusterName                = "DICE_CLUSTER_NAME"
	ENABLE_SPECIFIED_K8S_NAMESPACE = "ENABLE_SPECIFIED_K8S_NAMESPACE"
)

// diceCIODiscardKeys 需要丢弃的
var diceCIDiscardKeys = []string{
	"DICE_SSH_PASSWORD",
	"DICE_SSH_USER",
	"ETCD_ENDPOINTS",
}

// diceAddonsInfoKeys
var diceAddonsInfoKeys = []string{
	"REGISTRY_ADDR",
	"NEXUS_ADDR",
	"NEXUS_USERNAME",
	"NEXUS_PASSWORD",
	"SOLDIER_ADDR",
	"MS_NACOS_HOST",
	"MS_NACOS_PORT",
	"MS_MYSQL_HOST",
	"MS_MYSQL_PORT",
	"MS_MYSQL_USER",
	"MS_MYSQL_PASSWORD",
	"MS_MYSQL_DATABASE",
	"MS_POSTGRESQL_HOST",
	"MS_POSTGRESQL_PORT",
	"MS_POSTGRESQL_USER",
	"MS_POSTGRESQL_PASSWORD",
	"MS_POSTGRESQL_DATABASE",
	"ISTIO_ALIYUN",
	"ISTIO_INSTALLED",
	"ISTIO_VERSION",
}

// ClusterInfo is the object to encapsulate cluster info
type ClusterInfo struct {
	load_mutex  sync.Mutex
	store       jsonstore.JsonStore   // store cluster info
	lock        *dlock.DLock          // distributed lock
	clusterName string                // cluster name
	data        map[string]string     // the data of cluster info
	addr        string                // k8s master address
	k8sClient   *kubernetes.Clientset // client-go client
}

// Option configures an ClusterInfo
type Option func(*ClusterInfo)

// New news an ClusterInfo
func New(clusterName string, options ...Option) (*ClusterInfo, error) {
	cm := &ClusterInfo{}

	for _, op := range options {
		op(cm)
	}

	cm.clusterName = clusterName

	// json store
	store, err := jsonstore.New()
	if err != nil {
		return nil, errors.Errorf("failed to new json store for clusterInfo: %v", err)
	}
	cm.store = store

	// distributed lock
	lockKey := strutil.Concat(dlockKeyPrefix, clusterName)
	lock, err := dlock.New(lockKey, func() {})
	if err != nil {
		return nil, errors.Errorf("failed to new dlock: %s, error: %v", lockKey, err)
	}
	cm.lock = lock

	return cm, nil
}

func WithKubernetesClient(client *kubernetes.Clientset) Option {
	return func(ci *ClusterInfo) {
		ci.k8sClient = client
	}
}

// WithCompleteParams provides an Option
func WithCompleteParams(addr string) Option {
	return func(ci *ClusterInfo) {
		ci.addr = addr
	}
}

// Load load clusterInfo for specified cluster
func (ci *ClusterInfo) Load() error {
	ci.load_mutex.Lock()
	defer ci.load_mutex.Unlock()

	namespace := metav1.NamespaceDefault

	if os.Getenv(ENABLE_SPECIFIED_K8S_NAMESPACE) != "" {
		namespace = os.Getenv(ENABLE_SPECIFIED_K8S_NAMESPACE)
	}
	var (
		cm      *corev1.ConfigMap
		addonCM *corev1.ConfigMap
		err     error
	)

	cm, err = ci.k8sClient.CoreV1().ConfigMaps(namespace).Get(context.Background(), clusterInfoConfigMapName, metav1.GetOptions{})
	if err != nil {
		return errors.Errorf("failed to get %s configMap, clusterName: %s, (%v)",
			clusterInfoConfigMapName, ci.clusterName, err)
	}

	// ignore specified fields
	for _, key := range diceCIDiscardKeys {
		delete(cm.Data, key)
	}

	ci.data = cm.Data

	addonCM, err = ci.k8sClient.CoreV1().ConfigMaps(namespace).Get(context.Background(), addonsConfigMapName, metav1.GetOptions{})
	if err != nil {
		return errors.Errorf("failed to get %s configMap, clusterName: %s, (%v)",
			addonsConfigMapName, ci.clusterName, err)
	}

	// add registry addr
	for _, key := range diceAddonsInfoKeys {
		if _, ok := addonCM.Data[key]; ok {
			ci.data[key] = addonCM.Data[key]
		}
	}

	// If addr is inet format, sync address to etcd, other component use it to visit cluster.
	// TODO: Compatible with inet protocol, will remove in the future version.
	if ci.data != nil && strings.HasPrefix(ci.addr, "inet://") {
		netPortalURL, err := parseNetPortalURL(ci.addr)
		if err != nil {
			logrus.Errorf("failed to parse netportal address, (%v)", err)
		}
		ci.data[netportalURLKeyName] = netPortalURL
	}

	return nil
}

// Get get cluster info
func (ci *ClusterInfo) Get() (map[string]string, error) {
	if len(ci.data) == 0 {
		if err := ci.Load(); err != nil {
			return nil, errors.Errorf("failed to load cluster info, clusterName: %s, (%v)",
				ci.clusterName, err)
		}
	}

	return ci.data, nil
}

// SyncStore Sync clusterInfo from configMap to other store(e.g. ETCD).
func (ci *ClusterInfo) SyncStore() error {
	if ci.clusterName == "" {
		return errors.New("cluster name is null")
	}

	if len(ci.data) == 0 {
		return errors.New("clusterInfo data is null")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// DLock
	cleanup, err := events.OnlyOne(ctx, ci.lock)
	defer cleanup()

	if err != nil {
		return errors.Errorf("failed to lock for cluster info, clusterName: %s, error: %v",
			ci.clusterName, err)
	}

	key := strutil.Concat(clusterInfoPrefix, ci.clusterName)
	if err := ci.store.Put(ctx, key, ci.data); err != nil {
		return errors.Errorf("failed to put cluster info to json store, key: %s (%v)", key, err)
	}

	return nil
}

// LoopLoadAndSync load cluster info sync
func (ci *ClusterInfo) LoopLoadAndSync(ctx context.Context, sync bool) {
	var loadErr error

	for {
		if loadErr := ci.Load(); loadErr != nil {
			logrus.Errorf("failed to loop load cluster info, (%v)", loadErr)
		}

		if sync && (loadErr == nil) {
			if err := ci.SyncStore(); err != nil {
				logrus.Errorf("failed to loop sync cluster info, cluster: %s, (%v)", ci.clusterName, err)
			}
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(loopSyncTimeout):
			continue
		}
	}
}

// parseNetPortalURL：parse netPortal url, e.g. inet://abc?ssl=on&direct=on/123/qq?a=b
func parseNetPortalURL(url string) (string, error) {
	if !strings.HasPrefix(url, "inet://") {
		return "", errors.New("no prefix: inet://")
	}

	url = strings.TrimPrefix(url, netportalURLPrefix)
	url = strings.Replace(url, "//", "/", -1)

	parts := strings.SplitN(url, "/", 3)
	if len(parts) < 2 {
		return "", errors.Errorf("invalid addr: %s", url)
	}

	return strutil.Concat(netportalURLPrefix, parts[0]), nil
}

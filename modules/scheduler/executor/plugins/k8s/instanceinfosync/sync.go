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

package instanceinfosync

import (
	"context"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/erda-project/erda/bundle"
	"github.com/erda-project/erda/modules/scheduler/instanceinfo"
	"github.com/erda-project/erda/pkg/clientgo/kubernetes"
	"github.com/erda-project/erda/pkg/strutil"
)

//Synchronization strategy:
// 0. Periodically list all deployment, statefulset, and pod states
// 1. watch deployment, statefulset, pod
// TODO: 2. watch event for more detail messages

//type eventUtils interface {
//	watch pod events in all namespaces, use ctx to cancel
//WatchPodEventsAllNamespaces(ctx context.Context, callback func(*corev1.Event)) error
//}

type Syncer struct {
	clusterName string
	dbUpdater   *instanceinfo.Client
	bdl         *bundle.Bundle
	client      *kubernetes.Clientset
}

func NewSyncer(clusterName string, db *instanceinfo.Client, bdl *bundle.Bundle, client *kubernetes.Clientset) *Syncer {
	return &Syncer{
		clusterName: clusterName,
		dbUpdater:   db,
		bdl:         bdl,
		client:      client,
	}
}

func (s *Syncer) Sync(ctx context.Context) {
	s.listSync(ctx)
	s.watchSync(ctx)
	s.gc(ctx)
	<-ctx.Done()
}

func (s *Syncer) listSync(ctx context.Context) {
	go s.listSyncPod(ctx)
}

func (s *Syncer) watchSync(ctx context.Context) {
	go s.watchSyncPod(ctx)
	go s.watchSyncEvent(ctx)
}

func (s *Syncer) gc(ctx context.Context) {
	go s.gcDeadInstances(ctx)
	go s.gcServices(ctx)
}

func (s *Syncer) listSyncDeployment(ctx context.Context) {
	var cont *string

	for {
		wait := waitSeconds(cont)
		select {
		case <-ctx.Done():
			return
		case <-time.After(wait):
		}

		options := metav1.ListOptions{
			FieldSelector: "metadata.namespace!=default,metadata.namespace!=kube-system",
			Limit:         100,
		}

		if cont != nil {
			options.Continue = *cont
		}

		deployList, err := s.client.AppsV1().Deployments(metav1.NamespaceAll).List(context.Background(), options)
		if err != nil {
			logrus.Errorf("failed to list deployments: %v", err)
			cont = nil
			continue
		}

		// if returned continue = nil, means that this is the last part of the list
		if deployList.ListMeta.Continue != "" {
			cont = &deployList.ListMeta.Continue
		}

		if err := updateStatelessServiceDeployment(s.dbUpdater, deployList, false); err != nil {
			logrus.Errorf("failed to update statless-service serviceinfo: %v", err)
			continue
		}
		if err := updateAddonDeployment(s.dbUpdater, deployList, false); err != nil {
			logrus.Errorf("failed to update addon serviceinfo: %v", err)
			continue
		}
	}
}

func (s *Syncer) listSyncStatefulSet(ctx context.Context) {
	var cont *string

	for {
		wait := waitSeconds(cont)
		select {
		case <-ctx.Done():
			return
		case <-time.After(wait):
		}
		waitSeconds(cont)

		options := metav1.ListOptions{
			FieldSelector: "metadata.namespace!=default,metadata.namespace!=kube-system",
			Limit:         100,
		}

		if cont != nil {
			options.Continue = *cont
		}

		stsList, err := s.client.AppsV1().StatefulSets(metav1.NamespaceAll).List(context.TODO(), options)
		if err != nil {
			logrus.Errorf("failed to list statefulset: %v", err)
			cont = nil
			continue
		}

		if stsList.ListMeta.Continue != "" {
			cont = &stsList.ListMeta.Continue
		}

		if err := updateAddonStatefulSet(s.dbUpdater, stsList, false); err != nil {
			logrus.Errorf("failed to update addon serviceinfo: %v", err)
			continue
		}
	}
}

func (s *Syncer) listSyncPod(ctx context.Context) {
	var initUpdateTime time.Time

	for {
		wait := waitSeconds(nil)
		select {
		case <-ctx.Done():
			return
		case <-time.After(wait):
		}
		initUpdateTime = time.Now()
		logrus.Infof("start listpods for cluster: %s", s.clusterName)

		podList, err := s.client.CoreV1().Pods(metav1.NamespaceAll).List(context.Background(), metav1.ListOptions{
			FieldSelector: strutil.Join([]string{
				"metadata.namespace!=default",
				"metadata.namespace!=kube-system"}, ","),
		})

		if err != nil {
			logrus.Errorf("failed to list pod: %v", err)
			continue
		}
		logrus.Infof("listpods(%d) for: %s", len(podList.Items), s.clusterName)
		if err := updatePodAndInstance(s.dbUpdater, podList, false, nil); err != nil {
			logrus.Errorf("failed to update instanceinfo: %v", err)
			continue
		}
		logrus.Infof("export podlist info start, cluster: %s", s.clusterName)
		exportPodErrInfo(s.bdl, podList)
		logrus.Infof("export podlist info end, cluster: %s", s.clusterName)
		logrus.Infof("updatepods for: %s, cluster", s.clusterName)
		// it is last part of pod list, so execute gcAliveInstancesInDB
		// GcAliveInstancesInDB is triggered after every 2 complete traversals
		cost := int(time.Now().Sub(initUpdateTime).Seconds())
		if err := gcAliveInstancesInDB(s.dbUpdater, cost, s.clusterName); err != nil {
			logrus.Errorf("failed to gcAliveInstancesInDB: %v", err)
		}
		cost2 := int(time.Now().Sub(initUpdateTime).Seconds())
		if err := gcPodsInDB(s.dbUpdater, cost2, s.clusterName); err != nil {
			logrus.Errorf("failed to gcPodsInDB: %v", err)
		}
		logrus.Infof("gcAliveInstancesInDB for cluster: %s", s.clusterName)
	}
}

func (s *Syncer) watchSyncDeployment(ctx context.Context) {
	addOrUpdate, del := updateDeploymentOnWatch(s.dbUpdater)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(10) * time.Second):
		}

		if err := WatchDeployInAllNamespace(ctx, s.client, addOrUpdate, addOrUpdate, del); err != nil {
			logrus.Errorf("failed to watch update deployment: %v", err)
		}
	}
}

func (s *Syncer) watchSyncStatefulSet(ctx context.Context) {
	addOrUpdate, del := updateAddonStatefulSetOnWatch(s.dbUpdater)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(10) * time.Second):
		}
		if err := WatchStsInAllNamespace(ctx, s.client, addOrUpdate, addOrUpdate, del); err != nil {
			logrus.Errorf("failed to watch update statefulset: %v", err)
		}
	}
}

func (s *Syncer) watchSyncPod(ctx context.Context) {
	addOrUpdate, del := updatePodOnWatch(s.bdl, s.dbUpdater)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(10) * time.Second):
		}
		if err := WatchPodInAllNamespace(ctx, s.client, addOrUpdate, addOrUpdate, del); err != nil {
			logrus.Errorf("failed to watch update pod: %v, cluster: %s", err, s.clusterName)
		}
	}
}

func (s *Syncer) watchSyncEvent(ctx context.Context) {
	callback := func(e *corev1.Event) {
		if e.Type == "Normal" {
			return
		}
		ns := e.InvolvedObject.Namespace
		name := e.InvolvedObject.Name

		pod, err := s.client.CoreV1().Pods(ns).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			logrus.Errorf("failed to get pod: %s/%s", ns, name)
			return
		}

		if err := updatePodAndInstance(s.dbUpdater, &corev1.PodList{Items: []corev1.Pod{*pod}}, false,
			map[string]*corev1.Event{pod.Namespace + "/" + pod.Name: e}); err != nil {
			logrus.Errorf("failed to updatepod: %v", err)
			return
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(10) * time.Second):
		}
		if err := WatchPodEventsAllNamespaces(ctx, s.client, callback); err != nil {
			logrus.Errorf("failed to watch event: %v, addr: %s", err, s.clusterName)
		}
	}
}

func (s *Syncer) gcDeadInstances(ctx context.Context) {
	if err := gcDeadInstancesInDB(s.dbUpdater); err != nil {
		logrus.Errorf("failed to gcInstancesInDB: %v", err)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(24) * time.Hour):
		}
		if err := gcDeadInstancesInDB(s.dbUpdater); err != nil {
			logrus.Errorf("failed to gcInstancesInDB: %v", err)
		}
	}
}

func (s *Syncer) gcServices(ctx context.Context) {
	if err := gcServicesInDB(s.dbUpdater); err != nil {
		logrus.Errorf("failed to gcServicesInDB: %v", err)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(24) * time.Hour):
		}
		if err := gcServicesInDB(s.dbUpdater); err != nil {
			logrus.Errorf("failed to gcServicesInDB: %v", err)
		}
	}
}

func waitSeconds(cont *string) time.Duration {
	randSec := rand.Intn(5)
	wait := time.Duration(180+randSec) * time.Second
	if cont == nil {
		wait = time.Duration(60+randSec) * time.Second
	}
	return wait
}

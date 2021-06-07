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
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	wathchtool "k8s.io/client-go/tools/watch"

	"github.com/erda-project/erda/modules/scheduler/instanceinfo"
	"github.com/erda-project/erda/pkg/clientgo/kubernetes"
	"github.com/erda-project/erda/pkg/strutil"
)

func updateAddonStatefulSet(dbClient *instanceinfo.Client, stsList *appsv1.StatefulSetList, delete bool) error {
	r := dbClient.ServiceReader()
	w := dbClient.ServiceWriter()

	for _, sts := range stsList.Items {
		var (
			cluster     string
			orgName     string
			orgID       string
			projectName string
			projectID   string
			workspace   string
			addonName   string

			phase     instanceinfo.ServicePhase
			message   string
			startedAt time.Time

			// 参考 terminus.io/dice/dice/docs/dice-env-vars.md
			envsuffixmap = map[string]*string{
				"DICE_CLUSTER_NAME": &cluster,
				"DICE_ORG_NAME":     &orgName,
				"DICE_ORG_ID":       &orgID,
				"DICE_PROJECT_NAME": &projectName,
				"DICE_PROJECT_ID":   &projectID,
				"DICE_WORKSPACE":    &workspace,
				"DICE_ADDON_NAME":   &addonName,
			}
		)

		// -------------------------------
		// 1. Get all needed information from statefulset
		// -------------------------------
		for _, env := range sts.Spec.Template.Spec.Containers[0].Env {
			for k, v := range envsuffixmap {
				if strutil.HasSuffixes(env.Name, k) {
					*v = env.Value
				}
			}
		}

		// If the content in envsuffixmap is empty, don't update this sts to DB
		// Because the services initiated by the dice deployment process should have these environment variables
		skipThisSts := false
		for _, v := range envsuffixmap {
			if *v == "" {
				skipThisSts = true
				break
			}
		}
		if skipThisSts {
			continue
		}
		phase = instanceinfo.ServicePhaseUnHealthy
		if sts.Spec.Replicas == nil || *sts.Spec.Replicas == sts.Status.ReadyReplicas {
			phase = instanceinfo.ServicePhaseHealthy
		}
		startedAt = sts.ObjectMeta.CreationTimestamp.Time

		// -------------------------------
		// 2. Update or create ServiceInfo record
		// -------------------------------
		svcs, err := r.ByOrgName(orgName).
			ByProjectName(projectName).
			ByWorkspace(workspace).
			ByServiceType("addon").
			Do()
		if err != nil {
			return err
		}
		serviceInfo := instanceinfo.ServiceInfo{
			Cluster:     cluster,
			OrgName:     orgName,
			OrgID:       orgID,
			ProjectName: projectName,
			ProjectID:   projectID,
			Workspace:   workspace,
			ServiceType: "addon",
			Phase:       phase,
			Message:     message,
			StartedAt:   startedAt,
		}
		switch len(svcs) {
		case 0:
			if delete {
				break
			}
			if err := w.Create(&serviceInfo); err != nil {
				return err
			}
		default:
			for _, svc := range svcs {
				serviceInfo.ID = svc.ID
				if delete {
					if err := w.Delete(svc.ID); err != nil {
						return err
					}
				} else {
					if err := w.Update(serviceInfo); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func updateAddonStatefulSetOnWatch(db *instanceinfo.Client) (func(*appsv1.StatefulSet), func(*appsv1.StatefulSet)) {
	addOrUpdateFunc := func(sts *appsv1.StatefulSet) {
		if err := updateAddonStatefulSet(db,
			&appsv1.StatefulSetList{Items: []appsv1.StatefulSet{*sts}}, false); err != nil {
			logrus.Errorf("failed to update statefulset: %v", err)
		}
	}
	deleteFunc := func(sts *appsv1.StatefulSet) {
		if err := updateAddonStatefulSet(db,
			&appsv1.StatefulSetList{Items: []appsv1.StatefulSet{*sts}}, true); err != nil {
			logrus.Errorf("failed to update(delete) statefulset: %v", err)
		}
	}
	return addOrUpdateFunc, deleteFunc
}

func WatchStsInAllNamespace(ctx context.Context, client *kubernetes.Clientset, addFunc, updateFunc, deleteFunc func(*appsv1.StatefulSet)) error {
	curStsList, err := client.AppsV1().StatefulSets(metav1.NamespaceAll).List(ctx, metav1.ListOptions{Limit: 10})
	if err != nil {
		logrus.Errorf("list statefulset (limit: 10) error: %v", err)
		return err
	}

	retryWatcher, err := wathchtool.NewRetryWatcher(curStsList.GetResourceVersion(), &cache.ListWatch{
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = "metadata.namespace!=kube-system"
			return client.AppsV1().Deployments(metav1.NamespaceAll).Watch(ctx, options)
		},
	})

	if err != nil {
		logrus.Errorf("watch statefulset expand kube-system namespace error: %v", err)
		return err
	}

	defer func() {
		retryWatcher.Stop()
		logrus.Info("watch appV1.statefulset resource done")
		retryWatcher.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case res, ok := <-retryWatcher.ResultChan():
			if !ok {
				logrus.Error("watch appV1.statefulset resource closed unexpectedly")
				break
			}

			if res.Object == nil {
				continue
			}

			sts := appsv1.StatefulSet{}
			byteSts, err := json.Marshal(res.Object)
			if err != nil {
				logrus.Errorf("failed to marshal event obj, err: %v", err)
				continue
			}
			if err = json.Unmarshal(byteSts, &sts); err != nil {
				logrus.Errorf("failed to unmarshal event obj to deploy obj, err: %v", err)
				continue
			}

			switch res.Type {
			case watch.Added:
				addFunc(&sts)
			case watch.Modified:
				updateFunc(&sts)
			case watch.Deleted:
				deleteFunc(&sts)
				logrus.Infof("ignore event: %v, %v", watch.Error, watch.Bookmark)
			}
		}
	}
}

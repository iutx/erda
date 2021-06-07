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

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	wathchtool "k8s.io/client-go/tools/watch"

	"github.com/erda-project/erda/pkg/clientgo/kubernetes"
	"github.com/erda-project/erda/pkg/strutil"
)

func WatchPodEventsAllNamespaces(ctx context.Context, client *kubernetes.Clientset, callback func(*corev1.Event)) error {
	curEventList, err := client.CoreV1().Events(metav1.NamespaceAll).List(ctx, metav1.ListOptions{Limit: 10})
	if err != nil {
		logrus.Errorf("list event (limit: 10) error: %v", err)
		return err
	}

	retryWatcher, err := wathchtool.NewRetryWatcher(curEventList.GetResourceVersion(), &cache.ListWatch{
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = strutil.Join([]string{
				"metadata.namespace!=kube-system",
				"metadata.namespace!=kube-public",
				"involvedObject.kind=Pod",
			}, ",")
			return client.CoreV1().Events(metav1.NamespaceAll).Watch(ctx, options)
		},
	})

	if err != nil {
		logrus.Errorf("watch pod event expand kube-system,kube-public namespace error: %v", err)
		return err
	}

	defer func() {
		retryWatcher.Stop()
		logrus.Info("watch pod event done")
		retryWatcher.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case res, ok := <-retryWatcher.ResultChan():
			if !ok {
				logrus.Error("watch pod event closed unexpectedly")
				break
			}

			if res.Object == nil {
				continue
			}

			event := corev1.Event{}
			byteEvent, err := json.Marshal(res.Object)
			if err != nil {
				logrus.Errorf("failed to marshal event obj, err: %v", err)
				continue
			}
			if err = json.Unmarshal(byteEvent, &event); err != nil {
				logrus.Errorf("failed to unmarshal event obj to event obj, err: %v", err)
				continue
			}

			callback(&event)
		}
	}
}

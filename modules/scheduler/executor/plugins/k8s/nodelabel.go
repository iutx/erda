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

package k8s

import (
	"context"
	"encoding/json"
	"github.com/erda-project/erda/modules/scheduler/executor/executortypes"
	"github.com/erda-project/erda/modules/scheduler/schedulepolicy/labelconfig"
	"github.com/erda-project/erda/pkg/strutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

func (k *Kubernetes) IPToHostname(ip string) string {
	nodeList, err := k.k8sClient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return ""
	}
	for _, node := range nodeList.Items {
		for _, addr := range node.Status.Addresses {
			if addr.Type == v1.NodeInternalIP && addr.Address == ip {
				return node.Name
			}
		}
	}
	return ""
}

// SetNodeLabels set the labels of k8s node
func (k *Kubernetes) SetNodeLabels(_ executortypes.NodeLabelSetting, hosts []string, labels map[string]string) error {
	// contents in 'hosts' maybe hostname or internalIP, it should be unified into hostname
	nodeList, err := k.k8sClient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		logrus.Errorf("failed to list nodes: %v", err)
		return err
	}
	updatedHosts := make([]string, 0)
	for _, host := range hosts {
		for _, node := range nodeList.Items {
			add := false
			for _, addr := range node.Status.Addresses {

				if addr.Address == host {
					add = true
					break
				}
			}
			if add {
				updatedHosts = append(updatedHosts, node.Name)
			}
		}
	}

	for _, host := range updatedHosts {
		prefixedLabels := map[string]*string{}
		node, err := k.k8sClient.CoreV1().Nodes().Get(context.Background(), host, metav1.GetOptions{})
		if err != nil {
			return err
		}

		// 1. unset all 'dice/' labels
		for k := range node.ObjectMeta.Labels {
			if !strutil.HasPrefixes(k, labelconfig.K8SLabelPrefix) {
				continue
			}
			prefixedLabels[k] = nil
		}

		// 2. set labels in param 'labels'
		for k := range labels {
			v := labels[k]
			prefixedKey := k
			if !strutil.HasPrefixes(prefixedKey, labelconfig.K8SLabelPrefix) {
				prefixedKey = strutil.Concat(labelconfig.K8SLabelPrefix, k)
			}
			prefixedLabels[prefixedKey] = &v
		}

		// 3. set them
		var patch struct {
			Metadata struct {
				Labels map[string]*string `json:"labels"` // Use '*string' to cover 'null' case
			} `json:"metadata"`
		}

		patch.Metadata.Labels = prefixedLabels
		patchData, err := json.Marshal(patch)
		if err != nil {
			return err
		}

		_, err = k.k8sClient.CoreV1().Nodes().Patch(context.Background(), host, types.MergePatchType, patchData, metav1.PatchOptions{})
		if err != nil {
			return err
		}

		return nil
	}
	return nil
}

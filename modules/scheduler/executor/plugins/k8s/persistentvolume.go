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
	"strings"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/strutil"
)

func (k *Kubernetes) DeletePV(sg *apistructs.ServiceGroup) error {
	if !IsGroupStateful(sg) {
		return nil
	}
	// todo:
	for _, service := range sg.Services {
		for _, bind := range service.Binds {
			hostPath := bind.HostPath
			// Find local disk
			if strings.HasPrefix(hostPath, "/") || len(hostPath) == 0 {
				continue
			}
			// todo: The pv name rule is uniformly produced by a certain function
			pvName := strutil.Concat("lp-", sg.ID, "-")
			if len(hostPath) > 8 {
				pvName = strutil.Concat(pvName, hostPath[:8])
			} else {
				pvName = strutil.Concat(pvName, hostPath)
			}

			// todo: Confirm that the PV is bound to the corresponding PVC of the service under the runtime
			list, err := k.k8sClient.CoreV1().PersistentVolumes().List(context.Background(), metav1.ListOptions{})
			if err != nil {
				logrus.Errorf("failed to list pv, runtime: %s, pv: %s, (%v)", sg.ID, pvName, err)
				continue
			}
			for i := range list.Items {
				if !strings.HasPrefix(list.Items[i].Name, pvName) {
					continue
				}
				logrus.Infof("succeed to got pvName: %s, phase: %v", list.Items[i].Name, list.Items[i].Status.Phase)
				if err := k.k8sClient.CoreV1().PersistentVolumes().Delete(context.Background(), list.Items[i].Name, metav1.DeleteOptions{}); err != nil {
					logrus.Errorf("failed to delete pv name: %s, (%v)", list.Items[i].Name, err)
				}
			}
		}
	}
	return nil
}

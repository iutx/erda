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

package resourceinfo

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/clientgo/kubernetes"
	"github.com/erda-project/erda/pkg/strutil"
)

type ResourceInfo struct {
	client *kubernetes.Clientset
}

func New(client *kubernetes.Clientset) *ResourceInfo {
	return &ResourceInfo{client: client}
}

// PARAM brief: Does not provide cpuusage, memusage data, reducing the overhead of calling k8sapi
func (ri *ResourceInfo) Get(brief bool) (apistructs.ClusterResourceInfoData, error) {
	podList := &corev1.PodList{}

	if !brief {
		var err error
		podList, err = ri.client.CoreV1().Pods(metav1.NamespaceAll).List(context.Background(), metav1.ListOptions{
			FieldSelector: "status.phase!=Succeeded,status.phase!=Failed",
		})
		if err != nil {
			return apistructs.ClusterResourceInfoData{}, nil
		}
	}

	nodeList, err := ri.client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		logrus.Errorf("failed to list nodes: %v", err)
		return apistructs.ClusterResourceInfoData{}, nil
	}

	podMap := splitPodsByNodeName(podList)
	nodeResourceInfoMap := map[string]*apistructs.NodeResourceInfo{}

	for _, no := range nodeList.Items {
		var ip net.IP
		for _, addr := range no.Status.Addresses {
			if addr.Type == corev1.NodeInternalIP {
				ip = net.ParseIP(addr.Address)
			}
		}

		// ignore, if internalIP not found
		if ip == nil {
			continue
		}

		nodeResourceInfoMap[ip.String()] = &apistructs.NodeResourceInfo{}
		info := nodeResourceInfoMap[ip.String()]
		info.Labels = nodeLabels(&no)
		info.Ready = nodeReady(&no)
		cpuAllocatable, err := strconv.ParseFloat(fmt.Sprintf("%f", no.Status.Allocatable.Cpu().AsDec()), 64)
		if err != nil {
			return apistructs.ClusterResourceInfoData{}, err
		}

		memAllocatable, _ := no.Status.Allocatable.Memory().AsInt64()
		info.CPUAllocatable = cpuAllocatable
		info.MemAllocatable = memAllocatable
		pods := podMap[no.Name]

		podList := &corev1.PodList{Items: pods}
		reqs, limits := getPodsTotalRequestsAndLimits(podList)
		cpuReqs, cpuLimit, memReqs, memLimit := reqs[corev1.ResourceCPU], limits[corev1.ResourceCPU], reqs[corev1.ResourceMemory], limits[corev1.ResourceMemory]
		cpuReqsNum, err := strconv.ParseFloat(fmt.Sprintf("%f", cpuReqs.AsDec()), 64)
		if err != nil {
			return apistructs.ClusterResourceInfoData{}, err
		}

		memReqsNum, _ := memReqs.AsInt64()
		cpuLimitNum, err := strconv.ParseFloat(fmt.Sprintf("%f", cpuLimit.AsDec()), 64)
		if err != nil {
			return apistructs.ClusterResourceInfoData{}, err
		}

		memLimitNum, _ := memLimit.AsInt64()
		info.CPUReqsUsage = cpuReqsNum
		info.CPULimitUsage = cpuLimitNum
		info.MemReqsUsage = memReqsNum
		info.MemLimitUsage = memLimitNum
	}

	return apistructs.ClusterResourceInfoData{Nodes: nodeResourceInfoMap}, nil
}

func splitPodsByNodeName(podlist *corev1.PodList) map[string][]corev1.Pod {
	podmap := map[string][]corev1.Pod{}
	for i := range podlist.Items {
		if _, ok := podmap[podlist.Items[i].Spec.NodeName]; ok {
			podmap[podlist.Items[i].Spec.NodeName] = append(podmap[podlist.Items[i].Spec.NodeName], podlist.Items[i])
		} else {
			podmap[podlist.Items[i].Spec.NodeName] = []corev1.Pod{podlist.Items[i]}
		}
	}
	return podmap
}

func nodeLabels(n *corev1.Node) []string {
	r := []string{}
	for k := range n.ObjectMeta.Labels {
		if strutil.HasPrefixes(k, "dice/") {
			r = append(r, k)
		}
	}
	return r
}

func nodeReady(n *corev1.Node) bool {
	for _, cond := range n.Status.Conditions {
		if cond.Type == corev1.NodeReady {
			if cond.Status == "True" {
				return true
			}
		}
	}
	return false
}

// copy from kubectl
func getPodsTotalRequestsAndLimits(podList *corev1.PodList) (reqs map[corev1.ResourceName]resource.Quantity, limits map[corev1.ResourceName]resource.Quantity) {
	reqs, limits = map[corev1.ResourceName]resource.Quantity{}, map[corev1.ResourceName]resource.Quantity{}
	for _, pod := range podList.Items {
		podReqs, podLimits := PodRequestsAndLimits(&pod)
		for podReqName, podReqValue := range podReqs {
			if value, ok := reqs[podReqName]; !ok {
				reqs[podReqName] = podReqValue.DeepCopy()
			} else {
				value.Add(podReqValue)
				reqs[podReqName] = value
			}
		}
		for podLimitName, podLimitValue := range podLimits {
			if value, ok := limits[podLimitName]; !ok {
				limits[podLimitName] = podLimitValue.DeepCopy()
			} else {
				value.Add(podLimitValue)
				limits[podLimitName] = value
			}
		}
	}
	return
}

// copy from kubectl
// PodRequestsAndLimits returns a dictionary of all defined resources summed up for all
// containers of the pod.
func PodRequestsAndLimits(pod *corev1.Pod) (reqs, limits corev1.ResourceList) {
	reqs, limits = corev1.ResourceList{}, corev1.ResourceList{}
	for _, container := range pod.Spec.Containers {
		addResourceList(reqs, container.Resources.Requests)
		addResourceList(limits, container.Resources.Limits)
	}
	// init containers define the minimum of any resource
	for _, container := range pod.Spec.InitContainers {
		maxResourceList(reqs, container.Resources.Requests)
		maxResourceList(limits, container.Resources.Limits)
	}
	return
}

// copy from kubectl
// addResourceList adds the resources in newList to list
func addResourceList(list, new corev1.ResourceList) {
	for name, quantity := range new {
		if value, ok := list[name]; !ok {
			list[name] = quantity.DeepCopy()
		} else {
			value.Add(quantity)
			list[name] = value
		}
	}
}

// copy from kubectl
// maxResourceList sets list to the greater of list/newList for every resource
// either list
func maxResourceList(list, new corev1.ResourceList) {
	for name, quantity := range new {
		if value, ok := list[name]; !ok {
			list[name] = quantity.DeepCopy()
			continue
		} else {
			if quantity.Cmp(value) > 0 {
				list[name] = quantity.DeepCopy()
			}
		}
	}
}

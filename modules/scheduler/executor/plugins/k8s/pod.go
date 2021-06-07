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
	"github.com/erda-project/erda/pkg/strutil"
	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	None                  UnreadyReason = "None"
	ImagePullFailed       UnreadyReason = "ImagePullFailed"
	InsufficientResources UnreadyReason = "InsufficientResources"
	Unschedulable         UnreadyReason = "Unschedulable"
	ProbeFailed           UnreadyReason = "ProbeFailed"
	ContainerCannotRun    UnreadyReason = "ContainerCannotRun"
)

type UnreadyReason string

type PodStatus struct {
	Reason  UnreadyReason
	Message string
}

func getUnreadyPodReason(pod *corev1.Pod) (UnreadyReason, string) {
	// image pull failed
	for _, container := range pod.Status.ContainerStatuses {
		if container.State.Waiting != nil && container.State.Waiting.Reason == "ImagePullBackOff" {
			return ImagePullFailed, container.State.Waiting.Message
		}
	}

	// insufficient resources
	for _, cond := range pod.Status.Conditions {
		if cond.Type == "PodScheduled" &&
			cond.Status == corev1.ConditionFalse &&
			cond.Reason == "Unschedulable" &&
			strutil.Contains(cond.Message, "nodes are available") &&
			strutil.Contains(cond.Message, "Insufficient memory", "Insufficient cpu") {
			return InsufficientResources, cond.Message
		}
	}

	// other unschedule reasons
	for _, cond := range pod.Status.Conditions {
		if cond.Type == "PodScheduled" && cond.Status == corev1.ConditionFalse && cond.Reason == "Unschedulable" {
			return Unschedulable, cond.Message
		}
	}

	// container cannot run
	if pod.Status.Phase == "Running" {
		for _, cond := range pod.Status.Conditions {
			if cond.Type == "ContainersReady" && cond.Status == corev1.ConditionFalse && cond.Reason == "ContainersNotReady" {
				for _, container := range pod.Status.ContainerStatuses {
					if container.Ready == false &&
						container.Started != nil &&
						*container.Started == false &&
						container.LastTerminationState.Terminated != nil &&
						container.LastTerminationState.Terminated.Reason == "ContainerCannotRun" {
						return ContainerCannotRun, container.LastTerminationState.Terminated.Message
					}
				}
			}
		}
	}
	// probe failed
	if pod.Status.Phase == "Running" {
		for _, cond := range pod.Status.Conditions {
			if cond.Type == "ContainersReady" && cond.Status == corev1.ConditionFalse && cond.Reason == "ContainersNotReady" {
				return ProbeFailed, cond.Message
			}
		}
	}

	return None, ""
}

func (k *Kubernetes) GetNamespacedPodsStatus(namespace string) ([]PodStatus, error) {
	pods, err := k.k8sClient.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	r := make([]PodStatus, 0)
	for _, pod := range pods.Items {
		reason, message := getUnreadyPodReason(&pod)
		switch reason {
		case None:
		case ImagePullFailed:
			r = append(r, PodStatus{Reason: ImagePullFailed, Message: message})
		case InsufficientResources:
			r = append(r, PodStatus{Reason: InsufficientResources, Message: message})
		case Unschedulable:
			r = append(r, PodStatus{Reason: Unschedulable, Message: message})
		case ProbeFailed:
			r = append(r, PodStatus{Reason: ProbeFailed, Message: message})
		case ContainerCannotRun:
			r = append(r, PodStatus{Reason: ContainerCannotRun, Message: message})
		}
	}
	return r, nil
}

func (k *Kubernetes) killPod(namespace, podName string) error {
	return k.k8sClient.CoreV1().Pods(namespace).Delete(context.Background(), podName, metav1.DeleteOptions{})
}

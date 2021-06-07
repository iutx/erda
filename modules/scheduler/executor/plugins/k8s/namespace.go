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

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/strutil"
)

// MakeNamespace Generate a Namespace name
// Each runtime corresponds to a k8s namespace on k8s,
// format is ${runtimeNamespace}--${runtimeName}
func MakeNamespace(sg *apistructs.ServiceGroup) string {
	if IsGroupStateful(sg) {
		// Create a new namespace for the servicegroup that needs to be split into multiple statefulsets, that is, add the group- prefix
		if v, ok := sg.Labels[groupNum]; ok && v != "" && v != "1" {
			return strutil.Concat("group-", sg.Type, "--", sg.ID)
		}
	}
	return strutil.Concat(sg.Type, "--", sg.ID)
}

// CreateNamespace create namespace
func (k *Kubernetes) CreateNamespace(nsName string, sg *apistructs.ServiceGroup) error {
	notfound, err := k.NotfoundNamespace(nsName)
	if err != nil {
		return err
	}

	if !notfound {
		if sg.ProjectNamespace != "" {
			return nil
		}
		return errors.Errorf("failed to create namespace, ns: %s, (namespace already exists)", nsName)
	}

	labels := map[string]string{}

	if sg.Labels["service-mesh"] == "on" {
		labels["istio-injection"] = "enabled"
	}

	ns := GenerateNamespaceWithLabels(nsName, labels)

	_, err = k.k8sClient.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	// Create imagePullSecret under this namespace
	if err = k.NewRuntimeImageSecret(nsName, sg); err != nil {
		logrus.Errorf("failed to create imagePullSecret, namespace: %s, (%v)", ns, err)
	}

	return nil
}

// UpdateNamespace
func (k *Kubernetes) UpdateNamespace(nsName string, sg *apistructs.ServiceGroup) error {
	notfound, err := k.NotfoundNamespace(nsName)
	if err != nil {
		return err
	}
	if notfound {
		return errors.Errorf("not found ns: %v", nsName)
	}

	labels := map[string]string{}

	if sg.Labels["service-mesh"] == "on" {
		labels["istio-injection"] = "enabled"
	}

	ns := GenerateNamespaceWithLabels(nsName, labels)

	_, err = k.k8sClient.CoreV1().Namespaces().Update(context.Background(), ns, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// NotfoundNamespace not found namespace
func (k *Kubernetes) NotfoundNamespace(ns string) (bool, error) {
	_, err := k.k8sClient.CoreV1().Namespaces().Get(context.Background(), ns, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return true, nil
		}
		return false, err
	}

	return false, nil
}

// GenerateNamespaceWithLabels generate namespace with labels
func GenerateNamespaceWithLabels(ns string, labels map[string]string) *apiv1.Namespace {
	return &apiv1.Namespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   ns,
			Labels: labels,
		},
	}
}

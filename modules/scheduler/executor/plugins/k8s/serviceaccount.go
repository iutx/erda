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
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (k *Kubernetes) updateDefaultServiceAccountForImageSecret(namespace, secretName string) error {
	// Try to create first, then update after failure
	// k8s will automatically create the default serviceaccount, but there will be a delay, resulting in failure to update the probability.
	newSa := newServiceAccount(defaultServiceAccountName, namespace, []string{secretName})

	_, err := k.k8sClient.CoreV1().ServiceAccounts(namespace).Create(context.Background(), newSa, metav1.CreateOptions{})
	if err != nil {
		for {
			sa, err := k.k8sClient.CoreV1().ServiceAccounts(namespace).Get(context.Background(), defaultServiceAccountName, metav1.GetOptions{})
			if err != nil {
				return err
			}

			sa.ImagePullSecrets = append(sa.ImagePullSecrets, corev1.LocalObjectReference{
				Name: secretName,
			})

			patchData, err := json.Marshal(sa)
			if err != nil {
				return errors.Errorf("failed to Patch serviceaccount, name: %s, (%v)", sa.Name, err)
			}

			_, err = k.k8sClient.CoreV1().ServiceAccounts(sa.Namespace).Patch(context.Background(), sa.Name,
				types.StrategicMergePatchType, patchData, metav1.PatchOptions{})

			if !k8serrors.IsConflict(err) {
				return err
			}
		}
	}

	return nil
}

func newServiceAccount(name, namespace string, imageSecrets []string) *corev1.ServiceAccount {
	sa := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	for _, is := range imageSecrets {
		sa.ImagePullSecrets = append(sa.ImagePullSecrets, corev1.LocalObjectReference{
			Name: is,
		})
	}

	return sa
}

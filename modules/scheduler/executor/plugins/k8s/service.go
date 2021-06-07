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
	"fmt"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/istioctl"
	"github.com/erda-project/erda/pkg/strutil"
)

// CreateService create k8s service
func (k *Kubernetes) CreateService(service *apistructs.Service) error {
	svc := newService(service)
	_, err := k.k8sClient.CoreV1().Services(svc.Namespace).Create(context.Background(), svc, metav1.CreateOptions{})
	return err
}

// PutService update k8s service
func (k *Kubernetes) PutService(svc *apiv1.Service) error {
	_, err := k.k8sClient.CoreV1().Services(svc.Namespace).Update(context.TODO(), svc, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// GetService get k8s service
func (k *Kubernetes) GetService(namespace, name string) (*apiv1.Service, error) {
	return k.k8sClient.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

// DeleteService delete k8s service
func (k *Kubernetes) DeleteService(namespace, name string) error {
	return k.k8sClient.CoreV1().Services(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
}

// Port changes in the service description will cause changes in services and ingress
func (k *Kubernetes) updateService(service *apistructs.Service) error {
	// Service.Ports is empty, indicating that no service is expected
	if len(service.Ports) == 0 {
		// There is a service before the update, if there is no service, delete the servicece
		if err := k.DeleteService(service.Namespace, service.Name); err != nil {
			return err
		}

		return nil
	}

	svc, getErr := k.GetService(service.Namespace, service.Name)
	if getErr != nil && !k8serrors.IsNotFound(getErr) {
		return errors.Errorf("failed to get service, name: %s, (%v)", service.Name, getErr)
	}

	// If not found, create a new k8s service
	if k8serrors.IsNotFound(getErr) {
		if err := k.CreateService(service); err != nil {
			return err
		}
	} else {
		// If there is service before the update, if there is a service, then update the service
		desiredService := newService(service)

		if diffServiceMetadata(desiredService, svc) {
			desiredService.ResourceVersion = svc.ResourceVersion
			desiredService.Spec.ClusterIP = svc.Spec.ClusterIP

			if err := k.PutService(desiredService); err != nil {
				return err
			}
		}
	}

	return nil
}

func newService(service *apistructs.Service) *apiv1.Service {
	if len(service.Ports) == 0 {
		return nil
	}

	appValue := service.Env[KeyOriginServiceName]
	if appValue == "" {
		appValue = service.Name
	}

	k8sService := &apiv1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: service.Namespace,
			Labels:    make(map[string]string),
		},
		Spec: apiv1.ServiceSpec{
			// TODO: type?
			//Type: ServiceTypeLoadBalancer,
			Selector: map[string]string{
				"app": appValue,
			},
		},
	}

	setServiceLabelSelector(service, k8sService)

	for i, port := range service.Ports {
		k8sService.Spec.Ports = append(k8sService.Spec.Ports, apiv1.ServicePort{
			// TODO: name?
			Name: strutil.Concat(strings.ToLower(port.Protocol), "-", strconv.Itoa(i)),
			Port: int32(port.Port),
			// The user on Dice only fills in Port, that is, Port (port exposed by service) and targetPort (port exposed by container) are the same
			TargetPort: intstr.FromInt(port.Port),
			// Append protocol feature, Protocol Type contains TCP, UDP, SCTP
			Protocol: port.L4Protocol,
		})
	}
	return k8sService
}

// TODO: Complete me.
func diffServiceMetadata(left, right *apiv1.Service) bool {
	// compare the fields of Metadata and Spec
	if !reflect.DeepEqual(left.Labels, right.Labels) {
		return true
	}

	if !reflect.DeepEqual(left.Spec.Ports, right.Spec.Ports) {
		return true
	}

	if !reflect.DeepEqual(left.Spec.Selector, right.Spec.Selector) {
		return true
	}

	return false
}

// actually get deployment's names list, as k8s service would not be created
// if no ports exposed
func (k *Kubernetes) listServiceName(namespace string, labelSelector map[string]string) ([]string, error) {
	srvNames := make([]string, 0)

	deployList, err := k.k8sClient.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector).String(),
	})

	if err != nil {
		return srvNames, err
	}

	for _, item := range deployList.Items {
		srvNames = append(srvNames, item.Name)
	}

	dss, err := k.k8sClient.AppsV1().DaemonSets(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector).String(),
	})
	if err != nil {
		return srvNames, err
	}

	for _, item := range dss.Items {
		srvNames = append(srvNames, item.Name)
	}

	return srvNames, nil
}

func getServiceName(service *apistructs.Service) string {
	if service.Env[ProjectNamespace] == "true" {
		return service.Env[ProjectNamespaceServiceNameNameKey]
	}
	return service.Name
}

func setServiceLabelSelector(service *apistructs.Service, k8sService *apiv1.Service) {
	if v, ok := service.Env[ProjectNamespace]; ok && v == "true" && service.Name == service.Env[ProjectNamespaceServiceNameNameKey] {
		k8sService.Spec.Selector[LabelServiceGroupID] = service.Env[KeyServiceGroupID]
	}
}

func (k *Kubernetes) deleteRuntimeServiceWithProjectNamespace(service apistructs.Service) error {
	if service.Env[ProjectNamespace] == "true" {
		err := k.DeleteService(service.Namespace, service.Env[ProjectNamespaceServiceNameNameKey])
		if err != nil {
			return fmt.Errorf("delete service %s error: %v", service.Name, err)
		}
	}
	if k.istioEngine != istioctl.EmptyEngine {
		err := k.istioEngine.OnServiceOperator(istioctl.ServiceDelete, &service)
		if err != nil {
			return fmt.Errorf("delete istio resource error: %v", err)
		}
	}
	return nil
}

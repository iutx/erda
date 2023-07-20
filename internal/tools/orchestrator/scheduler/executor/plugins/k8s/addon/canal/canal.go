// Copyright (c) 2021 Terminus, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package canal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/internal/tools/orchestrator/scheduler/executor/plugins/k8s/addon"
	canalv1 "github.com/erda-project/erda/internal/tools/orchestrator/scheduler/executor/plugins/k8s/addon/canal/v1"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/schedule/schedulepolicy/constraintbuilders"
	"github.com/erda-project/erda/pkg/strutil"
)

type CanalOperator struct {
	cs     kubernetes.Interface
	k8s    addon.K8SUtil
	ns     addon.NamespaceUtil
	secret addon.SecretUtil
	pvc    addon.PVCUtil
	client *httpclient.HTTPClient
}

func (c *CanalOperator) Name(sg *apistructs.ServiceGroup) string {
	s := sg.Services[0].Env["CANAL_SERVER"]
	if s == "" {
		s = "canal-" + sg.ID[:10]
	}
	return s
}
func (c *CanalOperator) Namespace(sg *apistructs.ServiceGroup) string {
	return sg.ProjectNamespace
}
func (c *CanalOperator) NamespacedName(sg *apistructs.ServiceGroup) string {
	return c.Namespace(sg) + "/" + c.Name(sg)
}

func New(cs kubernetes.Interface, k8s addon.K8SUtil, ns addon.NamespaceUtil, secret addon.SecretUtil, pvc addon.PVCUtil, client *httpclient.HTTPClient) *CanalOperator {
	return &CanalOperator{
		cs:     cs,
		k8s:    k8s,
		ns:     ns,
		secret: secret,
		pvc:    pvc,
		client: client,
	}
}

func (c *CanalOperator) IsSupported() bool {
	res, err := c.client.Get(c.k8s.GetK8SAddr()).
		Path("/apis/database.erda.cloud/v1").Do().RAW()
	if err == nil {
		defer res.Body.Close()
		var b []byte
		b, err = io.ReadAll(res.Body)
		if err == nil {
			return bytes.Contains(b, []byte("canals"))
		}
	}
	logrus.Errorf("failed to query /apis/database.erda.cloud/v1, host: %s, err: %v",
		c.k8s.GetK8SAddr(), err)
	return false
}

func (c *CanalOperator) Validate(sg *apistructs.ServiceGroup) error {
	operator, ok := sg.Labels["USE_OPERATOR"]
	if !ok {
		return fmt.Errorf("[BUG] sg need USE_OPERATOR label")
	}
	if strings.ToLower(operator) != "canal" {
		return fmt.Errorf("[BUG] value of label USE_OPERATOR should be 'canal'")
	}
	if len(sg.Services) != 1 {
		return fmt.Errorf("illegal services num: %d", len(sg.Services))
	}
	if sg.Services[0].Name != "canal" {
		return fmt.Errorf("illegal service: %s, should be 'canal'", sg.Services[0].Name)
	}

	if sg.Services[0].Env["CANAL_DESTINATION"] == "" {
		return fmt.Errorf("illegal service: %s, need env 'CANAL_DESTINATION'", sg.Services[0].Name)
	}
	if sg.Services[0].Env["canal.instance.master.address"] == "" {
		return fmt.Errorf("illegal service: %s, need env 'canal.instance.master.address'", sg.Services[0].Name)
	}
	if sg.Services[0].Env["canal.instance.dbUsername"] == "" {
		return fmt.Errorf("illegal service: %s, need env 'canal.instance.dbUsername'", sg.Services[0].Name)
	}
	if sg.Services[0].Env["canal.instance.dbPassword"] == "" {
		return fmt.Errorf("illegal service: %s, need env 'canal.instance.dbPassword'", sg.Services[0].Name)
	}
	return nil
}

func (c *CanalOperator) Convert(sg *apistructs.ServiceGroup) interface{} {
	canal := sg.Services[0]

	scheinfo := sg.ScheduleInfo2
	scheinfo.Stateful = true
	affinity := constraintbuilders.K8S(&scheinfo, nil, nil, nil).Affinity

	resources := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{},
		Limits:   corev1.ResourceList{},
	}
	if canal.Resources.Cpu != 0 {
		cpu := resource.MustParse(strutil.Concat(strconv.Itoa(int(canal.Resources.Cpu*1000)), "m"))
		resources.Requests[corev1.ResourceCPU] = cpu
	}
	if canal.Resources.MaxCPU != 0 {
		maxCpu := resource.MustParse(strutil.Concat(strconv.Itoa(int(canal.Resources.MaxCPU*1000)), "m"))
		resources.Limits[corev1.ResourceCPU] = maxCpu
	}
	if canal.Resources.Mem != 0 {
		mem := resource.MustParse(strutil.Concat(strconv.Itoa(int(canal.Resources.Mem)), "Mi"))
		resources.Requests[corev1.ResourceMemory] = mem
	}
	if canal.Resources.MaxMem != 0 {
		maxMem := resource.MustParse(strutil.Concat(strconv.Itoa(int(canal.Resources.MaxMem)), "Mi"))
		resources.Limits[corev1.ResourceMemory] = maxMem
	}

	v := "v1.1.5"
	if canal.Env["CANAL_VERSION"] != "" {
		v = canal.Env["CANAL_VERSION"]
		if !strings.HasPrefix(v, "v") {
			v = "v" + v
		}
	}

	props := ""
	canalOptions := make(map[string]string)
	for k, v := range canal.Env {
		switch k {
		case "canal.zkServers", "canal.zookeeper.flush.period":
			canalOptions[k] = v
		default:
			if strings.HasPrefix(k, "canal.") {
				props += k + "=" + v + "\n"
			}
		}
	}

	obj := &canalv1.Canal{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "database.erda.cloud/v1",
			Kind:       "Canal",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.Name(sg),
			Namespace: c.Namespace(sg),
		},
		Spec: canalv1.CanalSpec{
			Version:      v,
			Image:        canal.Image,
			Replicas:     canal.Scale,
			Affinity:     &affinity,
			Resources:    resources,
			Labels:       make(map[string]string),
			CanalOptions: canalOptions,
			Annotations: map[string]string{
				"_destination": canal.Env["CANAL_DESTINATION"],
				"_props":       props,
			},
		},
	}

	addon.SetAddonLabelsAndAnnotations(canal, obj.Spec.Labels, obj.Spec.Annotations)

	return obj
}

func (c *CanalOperator) Create(k8syml interface{}) error {
	obj, ok := k8syml.(*canalv1.Canal)
	if !ok {
		return fmt.Errorf("[BUG] this k8syml should be *canalv1.Canal")
	}
	if err := c.ns.Exists(obj.Namespace); err != nil {
		if err := c.ns.Create(obj.Namespace, nil); err != nil {
			return err
		}
	}

	destination := obj.Spec.Annotations["_destination"]
	props := obj.Spec.Annotations["_props"]
	delete(obj.Spec.Annotations, "_destination")
	delete(obj.Spec.Annotations, "_props")

	{ // canal cm
		var b bytes.Buffer
		res, err := c.client.Get(c.k8s.GetK8SAddr()).
			Path(fmt.Sprintf("/api/v1/namespaces/%s/configmaps/%s", obj.Namespace, obj.Name)).
			Do().
			Body(&b)
		if err != nil {
			return fmt.Errorf("failed to get canal cm, %s/%s, err: %s, body: %s",
				obj.Namespace, obj.Name, err.Error(), b.String())
		}
		if res.IsNotfound() {
			cm := &corev1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "ConfigMap",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      obj.Name,
					Namespace: obj.Namespace,
				},
				Data: map[string]string{
					destination: props,
				},
			}
			{ //create cm
				var b bytes.Buffer
				res, err := c.client.Post(c.k8s.GetK8SAddr()).
					Path(fmt.Sprintf("/api/v1/namespaces/%s/configmaps", obj.Namespace)).
					JSONBody(cm).
					Do().
					Body(&b)
				if err != nil {
					return fmt.Errorf("failed to create canal cm, %s/%s, err: %s, body: %s",
						obj.Namespace, obj.Name, err.Error(), b.String())
				}
				if !res.IsOK() {
					return fmt.Errorf("failed to create canal cm, %s/%s, statuscode: %d, body: %s",
						obj.Namespace, obj.Name, res.StatusCode(), b.String())
				}
			}
		} else if !res.IsOK() {
			return fmt.Errorf("failed to get canal cm, %s/%s, statuscode: %d, body: %s",
				obj.Namespace, obj.Name, res.StatusCode(), b.String())
		} else {
			cm := &corev1.ConfigMap{}
			err := json.Unmarshal(b.Bytes(), cm)
			if err != nil {
				return fmt.Errorf("failed to unmarshal canal cm, %s/%s, err: %s, body: %s",
					obj.Namespace, obj.Name, err.Error(), b.String())
			}
			if cm.Data == nil {
				cm.Data = map[string]string{}
			}
			cm.Data[destination] = props

			{ //update cm
				var b bytes.Buffer
				res, err := c.client.Put(c.k8s.GetK8SAddr()).
					Path(fmt.Sprintf("/api/v1/namespaces/%s/configmaps/%s", obj.Namespace, obj.Name)).
					JSONBody(cm).
					Do().
					Body(&b)
				if err != nil {
					return fmt.Errorf("failed to update canal cm, %s/%s, err: %s, body: %s",
						obj.Namespace, obj.Name, err.Error(), b.String())
				}
				if !res.IsOK() {
					return fmt.Errorf("failed to update canal cm, %s/%s, statuscode: %d, body: %s",
						obj.Namespace, obj.Name, res.StatusCode(), b.String())
				}
			}
		}
	}

	{ // canal crd
		var b bytes.Buffer
		res, err := c.client.Get(c.k8s.GetK8SAddr()).
			Path(fmt.Sprintf("/apis/database.erda.cloud/v1/namespaces/%s/canals/%s", obj.Namespace, obj.Name)).
			Do().
			Body(&b)
		if err != nil {
			return fmt.Errorf("failed to get canal, %s/%s, err: %s, body: %s",
				obj.Namespace, obj.Name, err.Error(), b.String())
		}
		if res.IsNotfound() {
			var b bytes.Buffer
			res, err := c.client.Post(c.k8s.GetK8SAddr()).
				Path(fmt.Sprintf("/apis/database.erda.cloud/v1/namespaces/%s/canals", obj.Namespace)).
				JSONBody(obj).
				Do().
				Body(&b)
			if err != nil {
				return fmt.Errorf("failed to create canal, %s/%s, err: %s, body: %s",
					obj.Namespace, obj.Name, err.Error(), b.String())
			}
			if !res.IsOK() {
				return fmt.Errorf("failed to create canal, %s/%s, statuscode: %d, body: %s",
					obj.Namespace, obj.Name, res.StatusCode(), b.String())
			}
		} else if !res.IsOK() {
			return fmt.Errorf("failed to get canal, %s/%s, statuscode: %d, body: %s",
				obj.Namespace, obj.Name, res.StatusCode(), b.String())
		} else {
			//TODO: wait canal cm destination
			time.Sleep(5 * time.Second)
		}
	}

	labels := map[string]string{
		"addon": "canal",
		"group": obj.Name,
	}

	// Create Service NodePort
	npService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-x-np", obj.Name),
			Namespace: obj.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:        "canal",
					Protocol:    corev1.ProtocolTCP,
					AppProtocol: nil,
					Port:        11111,
					TargetPort: intstr.IntOrString{
						IntVal: 11111,
					},
				},
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	_, err := c.cs.CoreV1().Services(obj.Namespace).Create(context.Background(), npService, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create canal %s/%s node port service, %v", obj.Namespace, obj.Name, err)
	}

	return nil
}

func (c *CanalOperator) Inspect(sg *apistructs.ServiceGroup) (*apistructs.ServiceGroup, error) {
	var b bytes.Buffer
	res, err := c.client.Get(c.k8s.GetK8SAddr()).
		Path(fmt.Sprintf("/apis/database.erda.cloud/v1/namespaces/%s/canals/%s", c.Namespace(sg), c.Name(sg))).
		Do().
		Body(&b)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect canal, %s, err: %s",
			c.NamespacedName(sg), err.Error())
	}
	if !res.IsOK() {
		return nil, fmt.Errorf("failed to inspect canal, %s, statuscode: %d, body: %s",
			c.NamespacedName(sg), res.StatusCode(), b.String())
	}
	obj := new(canalv1.Canal)
	if err := json.NewDecoder(&b).Decode(obj); err != nil {
		return nil, err
	}

	canalsvc := &(sg.Services[0])

	if obj.Status.Color == canalv1.Green {
		canalsvc.Status = apistructs.StatusHealthy
		sg.Status = apistructs.StatusHealthy
	} else {
		canalsvc.Status = apistructs.StatusUnHealthy
		sg.Status = apistructs.StatusUnHealthy
	}

	canalsvc.Vip = strings.Join([]string{
		obj.BuildName("x"),
		obj.Namespace,
		"svc.cluster.local",
	}, ".")

	externalEp := &apistructs.ExternalEndpoint{
		Hosts: make([]string, 0),
		Ports: make([]int32, 0),
	}

	npSvc, err := c.cs.CoreV1().Services(obj.Namespace).Get(context.Background(),
		fmt.Sprintf("%s-x-np", obj.Name), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if len(npSvc.Spec.Ports) != 0 && npSvc.Spec.Ports[0].NodePort != 0 {
		externalEp.Ports = append(externalEp.Ports, npSvc.Spec.Ports[0].NodePort)
	}

	pods, err := c.cs.CoreV1().Pods(obj.Namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: labels.FormatLabels(map[string]string{
			"group": obj.Name,
			"addon": "canal",
		}),
	})
	if err != nil {
		return nil, err
	}

	logrus.Infof("canal-nodeport pod get %d", len(pods.Items))

	for _, pod := range pods.Items {
		if pod.Status.HostIP == "" {
			return nil, fmt.Errorf("canal-nodeport pod %s/%s is not ready", pod.Namespace, pod.Name)
		}
		externalEp.Hosts = append(externalEp.Hosts, pod.Status.HostIP)
		canalsvc.ExternalEndpoint = externalEp
	}

	if obj.Spec.Replicas > 1 {
		options := obj.Spec.CanalOptions
		_, zkServersOk := options["canal.zkServers"]
		_, registerIpOk := options["canal.register.ip"]
		port, _ := options["canal.port"]

		if zkServersOk {
			needUpdate := false
			if !registerIpOk {
				needUpdate = true
				options["canal.register.ip"] = canalsvc.ExternalEndpoint.Hosts[0]
			}
			if canalsvc.ExternalEndpoint.Ports[0] != 0 && port != strconv.Itoa(int(canalsvc.ExternalEndpoint.Ports[0])) {
				needUpdate = true
				nodePort := canalsvc.ExternalEndpoint.Ports[0]
				options["canal.port"] = strconv.Itoa(int(nodePort))
				if npSvc.Spec.Ports[0].Port != nodePort {
					npSvc.Spec.Ports[0].Port = nodePort
					npSvc.Spec.Ports[0].TargetPort = intstr.IntOrString{
						IntVal: nodePort,
					}
					npSvc.ResourceVersion = ""
					_, err := c.cs.CoreV1().Services(obj.Namespace).Update(context.Background(), npSvc, metav1.UpdateOptions{})
					if err != nil {
						return nil, fmt.Errorf("canal-nodeport failed to update custom node port, %d, err: %v", nodePort, err)
					}
				}

			}
			if needUpdate {
				obj.Spec.CanalOptions = options
				var b bytes.Buffer
				res, err := c.client.Put(c.k8s.GetK8SAddr()).
					Path(fmt.Sprintf("/apis/database.erda.cloud/v1/namespaces/%s/canals/%s", obj.Namespace, obj.Name)).
					JSONBody(obj).
					Do().
					Body(&b)
				if err != nil {
					return nil, fmt.Errorf("failed to update canal, %s/%s, err: %s, body: %s",
						obj.Namespace, obj.Name, err.Error(), b.String())
				}
				if !res.IsOK() {
					return nil, fmt.Errorf("failed to update canal, %s/%s, statuscode: %d, body: %s",
						obj.Namespace, obj.Name, res.StatusCode(), b.String())
				}
			}
		}
	}

	//TODO: check canal cm destination
	time.Sleep(5 * time.Second)

	return sg, nil
}

func (c *CanalOperator) Remove(sg *apistructs.ServiceGroup) error {
	canal := sg.Services[0]
	destination := canal.Env["CANAL_DESTINATION"]

	{ // canal cm
		var b bytes.Buffer
		res, err := c.client.Get(c.k8s.GetK8SAddr()).
			Path(fmt.Sprintf("/api/v1/namespaces/%s/configmaps/%s", c.Namespace(sg), c.Name(sg))).
			Do().
			Body(&b)
		if err != nil {
			return fmt.Errorf("failed to get canal cm, %s/%s, err: %s, body: %s",
				c.Namespace(sg), c.Name(sg), err.Error(), b.String())
		}
		if res.IsNotfound() {
			return nil
		} else if !res.IsOK() {
			return fmt.Errorf("failed to get canal cm, %s/%s, statuscode: %d, body: %s",
				c.Namespace(sg), c.Name(sg), res.StatusCode(), b.String())
		} else {
			cm := &corev1.ConfigMap{}
			err := json.Unmarshal(b.Bytes(), cm)
			if err != nil {
				return fmt.Errorf("failed to unmarshal canal cm, %s/%s, err: %s, body: %s",
					c.Namespace(sg), c.Name(sg), err.Error(), b.String())
			}

			if _, ok := cm.Data[destination]; !ok {
				return nil
			}
			delete(cm.Data, destination)

			var b bytes.Buffer
			res, err := c.client.Put(c.k8s.GetK8SAddr()).
				Path(fmt.Sprintf("/api/v1/namespaces/%s/configmaps/%s", c.Namespace(sg), c.Name(sg))).
				JSONBody(cm).
				Do().
				Body(&b)
			if err != nil {
				return fmt.Errorf("failed to update canal cm, %s/%s, err: %s, body: %s",
					c.Namespace(sg), c.Name(sg), err.Error(), b.String())
			}
			if !res.IsOK() {
				return fmt.Errorf("failed to update canal cm, %s/%s, statuscode: %d, body: %s",
					c.Namespace(sg), c.Name(sg), res.StatusCode(), b.String())
			}

			if len(cm.Data) > 0 {
				return nil
			}
		}
	}

	var b bytes.Buffer
	res, err := c.client.Delete(c.k8s.GetK8SAddr()).
		Path(fmt.Sprintf("/apis/database.erda.cloud/v1/namespaces/%s/canals/%s", c.Namespace(sg), c.Name(sg))).
		Do().
		Body(&b)
	if err != nil {
		return fmt.Errorf("failed to remove canal, %s, err: %s",
			c.NamespacedName(sg), err.Error())
	}
	if !res.IsOK() {
		return fmt.Errorf("failed to remove canal, %s, statuscode: %d, body: %s",
			c.NamespacedName(sg), res.StatusCode(), b.String())
	}

	namespace := c.Namespace(sg)
	name := c.Name(sg)
	npName := fmt.Sprintf("%s-x-np", name)

	_, err = c.cs.CoreV1().Services(namespace).Get(context.Background(), npName, metav1.GetOptions{})
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to get service %s, err: %v", c.NamespacedName(sg), err)
		}
		return nil
	}

	err = c.cs.CoreV1().Services(namespace).Delete(context.Background(), fmt.Sprintf("%s-x-np", name),
		metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete mysql %s node port service, %v", c.Name(sg), err)
	}

	return nil
}

func (c *CanalOperator) Update(k8syml interface{}) error {
	//TODO
	return fmt.Errorf("canaloperator not impl Update yet")
}

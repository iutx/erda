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

package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/erda-project/erda/pkg/clientgo"
	"strings"

	"github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/modules/scheduler/executor/plugins/k8s/addon"
	"github.com/erda-project/erda/modules/scheduler/schedulepolicy/constraintbuilders"
	commonv1 "github.com/erda-project/erda/pkg/clientgo/apis/elasticsearch/common/v1"
	v1 "github.com/erda-project/erda/pkg/clientgo/apis/elasticsearch/v1"
	"github.com/erda-project/erda/pkg/clientgo/customclient"
	"github.com/erda-project/erda/pkg/clientgo/kubernetes"
	"github.com/erda-project/erda/pkg/strutil"
)

type ElasticsearchAndSecret struct {
	v1.Elasticsearch
	corev1.Secret
}

type ElasticsearchOperator struct {
	k8sClient    *kubernetes.Clientset
	customClient *customclient.Clientset
	overcommit   addon.OvercommitUtil
	imageSecret  addon.ImageSecretUtil
}

func New(clientSet *clientgo.ClientSet, overcommit addon.OvercommitUtil, imageSecret addon.ImageSecretUtil) *ElasticsearchOperator {
	return &ElasticsearchOperator{
		k8sClient:    clientSet.K8sClient,
		customClient: clientSet.CustomClient,
		overcommit:   overcommit,
		imageSecret:  imageSecret,
	}
}

// IsSupported Determine whether to support  elasticseatch operator
func (eo *ElasticsearchOperator) IsSupported() bool {
	_, err := eo.customClient.ElasticsearchV1().Elasticsearches(metav1.NamespaceDefault).List(context.Background(),
		metav1.ListOptions{Limit: 1})
	if err != nil {
		logrus.Errorf("failed to query resource elasticsearch, err: %v", err)
		return false
	}

	return true
}

// Validate Verify the legality of the ServiceGroup transformed from diceyml
func (eo *ElasticsearchOperator) Validate(sg *apistructs.ServiceGroup) error {
	operator, ok := sg.Labels["USE_OPERATOR"]
	if !ok {
		return fmt.Errorf("[BUG] sg need USE_OPERATOR label")
	}
	if strutil.ToLower(operator) != "elasticsearch" {
		return fmt.Errorf("[BUG] value of label USE_OPERATOR should be 'elasticsearch'")
	}
	if _, ok := sg.Labels["VERSION"]; !ok {
		return fmt.Errorf("[BUG] sg need VERSION label")
	}
	if _, err := convertJsontToMap(sg.Services[0].Env["config"]); err != nil {
		return fmt.Errorf("[BUG] config Abnormal format of ENV")
	}
	return nil
}

// Convert Convert sg to cr, which is kubernetes yaml
func (eo *ElasticsearchOperator) Convert(sg *apistructs.ServiceGroup) interface{} {
	svc0 := sg.Services[0]
	scname := "dice-local-volume"
	// 官方建议，将堆内和堆外各设置一半, Xmx 和 Xms 设置的是堆内
	esMem := int(convertMiToMB(svc0.Resources.Mem) / 2)
	svc0.Env["ES_JAVA_OPTS"] = fmt.Sprintf("-Xms%dm -Xmx%dm", esMem, esMem)
	svc0.Env["ES_USER"] = "elastic"
	svc0.Env["ES_PASSWORD"] = svc0.Env["requirepass"]
	svc0.Env["SELF_PORT"] = "9200"

	scheinfo := sg.ScheduleInfo2
	scheinfo.Stateful = true
	affinity := constraintbuilders.K8S(&scheinfo, nil, nil, nil).Affinity.NodeAffinity

	nodeSets := eo.NodeSetsConvert(svc0, scname, affinity)
	es := v1.Elasticsearch{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Elasticsearch",
			APIVersion: "elasticsearch.k8s.elastic.co/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      sg.ID,
			Namespace: genK8SNamespace(sg.Type, sg.ID),
		},
		Spec: v1.ElasticsearchSpec{
			HTTP: commonv1.HTTPConfig{
				TLS: commonv1.TLSOptions{
					SelfSignedCertificate: &commonv1.SelfSignedCertificate{
						Disabled: true,
					},
				},
			},
			Version: sg.Labels["VERSION"],
			Image:   svc0.Image,
			NodeSets: []v1.NodeSet{
				nodeSets,
			},
		},
	}

	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-es-elastic-user", sg.ID),
			Namespace: genK8SNamespace(sg.Type, sg.ID),
		},
		Data: map[string][]byte{
			"elastic": []byte(svc0.Env["requirepass"]),
		},
	}

	return ElasticsearchAndSecret{Elasticsearch: es, Secret: secret}
}

func (eo *ElasticsearchOperator) Create(k8syml interface{}) error {
	elasticsearchAndSecret, ok := k8syml.(ElasticsearchAndSecret)
	if !ok {
		return fmt.Errorf("[BUG] this k8syml should be elasticsearchAndSecret")
	}

	elasticsearch := elasticsearchAndSecret.Elasticsearch
	secret := elasticsearchAndSecret.Secret

	_, err := eo.k8sClient.CoreV1().Namespaces().Create(context.Background(), &apiv1.Namespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: elasticsearch.Namespace,
		},
	}, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}

	_, err = eo.k8sClient.CoreV1().Secrets(secret.Namespace).Create(context.Background(), &secret, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}

	logrus.Debugf("es operator, start to create image secret, sepc:%+v", elasticsearch)
	if err := eo.imageSecret.NewImageSecret(elasticsearch.Namespace); err != nil {
		return err
	}

	_, err = eo.customClient.ElasticsearchV1().Elasticsearches(elasticsearch.Namespace).Create(context.Background(),
		&elasticsearch, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create elasticsearch, %s/%s, err: %v", elasticsearch.Namespace, elasticsearch.Name, err)
	}

	return nil
}

func (eo *ElasticsearchOperator) Inspect(sg *apistructs.ServiceGroup) (*apistructs.ServiceGroup, error) {
	nsName := genK8SNamespace(sg.Type, sg.ID)

	stsList, err := eo.k8sClient.AppsV1().StatefulSets(nsName).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	svcList, err := eo.k8sClient.CoreV1().Services(nsName).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var elasticsearch *apistructs.Service
	if sg.Services[0].Name == "elasticsearch" {
		elasticsearch = &(sg.Services[0])
	}
	for _, sts := range stsList.Items {
		if sts.Spec.Replicas == nil {
			elasticsearch.Status = apistructs.StatusUnknown
		} else if *sts.Spec.Replicas == sts.Status.ReadyReplicas {
			elasticsearch.Status = apistructs.StatusHealthy
		} else {
			elasticsearch.Status = apistructs.StatusUnHealthy
		}
	}

	for _, svc := range svcList.Items {
		if strings.Contains(svc.Name, "es-http") {
			elasticsearch.Vip = strutil.Join([]string{svc.Name, svc.Namespace, "svc.cluster.local"}, ".")
		}
	}

	if elasticsearch.Status == apistructs.StatusHealthy {
		sg.Status = apistructs.StatusHealthy
	} else {
		sg.Status = apistructs.StatusUnHealthy
	}
	return sg, nil
}

func (eo *ElasticsearchOperator) Remove(sg *apistructs.ServiceGroup) error {
	nsName := genK8SNamespace(sg.Type, sg.ID)

	if err := eo.customClient.ElasticsearchV1().Elasticsearches(nsName).Delete(context.Background(), sg.ID,
		metav1.DeleteOptions{}); err != nil {
		logrus.Errorf("failed to delete elasticsearch: %s: %v", nsName, err)
		return err
	}

	if err := eo.k8sClient.CoreV1().Namespaces().Delete(context.Background(), nsName, metav1.DeleteOptions{}); err != nil {
		logrus.Errorf("failed to delete namespace: %s: %v", nsName, err)
		return fmt.Errorf("failed to delete namespace: %s: %v", nsName, err)
	}
	return nil
}

// Update
// secret The update will not be performed, and a restart is required due to the update of the static password. (You can improve the multi-user authentication through the user management machine with perfect service)
func (eo *ElasticsearchOperator) Update(k8syml interface{}) error {
	elasticsearchAndSecret, ok := k8syml.(ElasticsearchAndSecret)
	if !ok {
		return fmt.Errorf("[BUG] this k8syml should be elasticsearchAndSecret")
	}

	elasticsearch := elasticsearchAndSecret.Elasticsearch
	secret := elasticsearchAndSecret.Secret

	_, err := eo.k8sClient.CoreV1().Namespaces().Get(context.Background(), elasticsearch.Namespace, metav1.GetOptions{})
	if err != nil {
		return err
	}

	_, err = eo.k8sClient.CoreV1().Secrets(secret.Namespace).Create(context.Background(), &secret, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}

	es, err := eo.Get(elasticsearch.Namespace, elasticsearch.Name)
	if err != nil {
		return err
	}
	// set last resource version.
	// fix error: "metadata.resourceVersion: Invalid value: 0x0: must be specified for an update"
	elasticsearch.ObjectMeta.ResourceVersion = es.ObjectMeta.ResourceVersion
	// set annotations controller-version
	// fix error: Resource was created with older version of operator, will not take action
	elasticsearch.ObjectMeta.Annotations = es.ObjectMeta.Annotations

	if _, err = eo.customClient.ElasticsearchV1().Elasticsearches(elasticsearch.Namespace).Update(context.Background(),
		&elasticsearch, metav1.UpdateOptions{}); err != nil {
		logrus.Errorf("failed to update elasticsearchs, %s/%s, err: %v", elasticsearch.Namespace, elasticsearch.Name, err)
		return err
	}

	return nil
}

func genK8SNamespace(namespace, name string) string {
	return strutil.Concat(namespace, "--", name)
}

func (eo *ElasticsearchOperator) NodeSetsConvert(svc apistructs.Service, scname string, affinity *corev1.NodeAffinity) v1.NodeSet {
	config, _ := convertJsontToMap(svc.Env["config"])
	nodeSets := v1.NodeSet{
		Name:  "addon",
		Count: int32(svc.Scale),
		Config: &commonv1.Config{
			Data: config,
		},
		PodTemplate: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				Affinity: &corev1.Affinity{NodeAffinity: affinity},
				Containers: []corev1.Container{
					{
						Name: "elasticsearch",
						Env:  envs(svc.Env),
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								"cpu": resource.MustParse(
									fmt.Sprintf("%dm", int(1000*eo.overcommit.CPUOvercommit(svc.Resources.Cpu)))),
								"memory": resource.MustParse(
									fmt.Sprintf("%dMi", int(svc.Resources.Mem))),
							},
							Limits: corev1.ResourceList{
								"cpu": resource.MustParse(
									fmt.Sprintf("%dm", int(1000*svc.Resources.Cpu))),
								"memory": resource.MustParse(
									fmt.Sprintf("%dMi", int(svc.Resources.Mem))),
							},
						},
					},
				},
			},
		},
		VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "elasticsearch-data",
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							"storage": resource.MustParse("10Gi"),
						},
					},
					StorageClassName: &scname,
				},
			},
		},
	}
	return nodeSets
}

func envs(envs map[string]string) []corev1.EnvVar {
	r := []corev1.EnvVar{}
	for k, v := range envs {
		r = append(r, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	selfHost := apiv1.EnvVar{
		Name: "SELF_HOST",
		ValueFrom: &apiv1.EnvVarSource{
			FieldRef: &apiv1.ObjectFieldSelector{
				APIVersion: "v1",
				FieldPath:  "status.podIP",
			},
		},
	}
	addonNodeID := apiv1.EnvVar{
		Name: "ADDON_NODE_ID",
		ValueFrom: &apiv1.EnvVarSource{
			FieldRef: &apiv1.ObjectFieldSelector{
				APIVersion: "v1",
				FieldPath:  "metadata.name",
			},
		},
	}

	r = append(r, selfHost, addonNodeID)

	return r
}

// convertMiToMB Convert MiB to MB
//1 MiB = 1.048576 MB
func convertMiToMB(mem float64) float64 {
	return mem * 1.048576
}

// convertJsontToMap Convert json to map
func convertJsontToMap(str string) (map[string]interface{}, error) {
	var tempMap map[string]interface{}
	if str == "" {
		return tempMap, nil
	}

	if err := json.Unmarshal([]byte(str), &tempMap); err != nil {
		return tempMap, err
	}
	return tempMap, nil
}

// Get get elasticsearchs resource information
func (eo *ElasticsearchOperator) Get(namespace, name string) (*v1.Elasticsearch, error) {
	es, err := eo.customClient.ElasticsearchV1().Elasticsearches(namespace).Get(context.Background(),
		name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get elasticsearchs: %s/%s, err: %v", namespace, name, err)
	}

	return es, nil
}

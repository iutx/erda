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

package redis

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/internal/tools/orchestrator/scheduler/executor/plugins/k8s/addon"
	"github.com/erda-project/erda/internal/tools/orchestrator/scheduler/executor/plugins/k8s/k8sapi"
	"github.com/erda-project/erda/internal/tools/orchestrator/scheduler/executor/util"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/schedule/schedulepolicy/constraintbuilders"
	"github.com/erda-project/erda/pkg/schedule/schedulepolicy/constraintbuilders/constraints"
	"github.com/erda-project/erda/pkg/strutil"
)

var (
	redisExporterImage = util.GetAddonPublicRegistry() + "/retag/redis-exporter:v1.45.0"
)

type RedisOperator struct {
	cs          kubernetes.Interface
	k8s         addon.K8SUtil
	deployment  addon.DeploymentUtil
	statefulset addon.StatefulsetUtil
	ns          addon.NamespaceUtil
	service     addon.ServiceUtil
	overcommit  addon.OvercommitUtil
	secret      addon.SecretUtil
	client      *httpclient.HTTPClient
}

func NewRedisOperator(
	cs kubernetes.Interface,
	k8sutil addon.K8SUtil,
	deploy addon.DeploymentUtil,
	sts addon.StatefulsetUtil,
	service addon.ServiceUtil,
	ns addon.NamespaceUtil,
	overcommit addon.OvercommitUtil,
	secret addon.SecretUtil,
	client *httpclient.HTTPClient) *RedisOperator {
	return &RedisOperator{
		k8s:         k8sutil,
		deployment:  deploy,
		statefulset: sts,
		service:     service,
		ns:          ns,
		overcommit:  overcommit,
		secret:      secret,
		client:      client,
		cs:          cs,
	}
}

func (ro *RedisOperator) IsSupported() bool {
	resp, err := ro.client.Get(ro.k8s.GetK8SAddr()).
		Path("/apis/databases.spotahome.com/v1").
		Do().
		DiscardBody()
	if err != nil {
		logrus.Errorf("failed to query /apis/databases.spotahome.com/v1, host: %v, err: %v",
			ro.k8s.GetK8SAddr(), err)
		return false
	}
	if !resp.IsOK() {
		return false
	}
	return true
}

// Validate 检查
func (ro *RedisOperator) Validate(sg *apistructs.ServiceGroup) error {
	operator, ok := sg.Labels["USE_OPERATOR"]
	if !ok {
		return fmt.Errorf("[BUG] sg need USE_OPERATOR label")
	}
	if strutil.ToLower(operator) != svcNameRedis {
		return fmt.Errorf("[BUG] value of label USE_OPERATOR should be 'redis'")
	}
	if len(sg.Services) != 2 {
		return fmt.Errorf("illegal services num: %d", len(sg.Services))
	}
	if sg.Services[0].Name != svcNameRedis && sg.Services[0].Name != svcNameSentinel {
		return fmt.Errorf("illegal service: %+v, should be one of [redis, sentinel]", sg.Services[0])
	}
	if sg.Services[1].Name != svcNameRedis && sg.Services[1].Name != svcNameSentinel {
		return fmt.Errorf("illegal service: %+v, should be one of [redis, sentinel]", sg.Services[1])
	}
	var redis apistructs.Service
	if sg.Services[0].Name == svcNameRedis {
		redis = sg.Services[0]
	}
	// if sg.Services[0].Name == svcNameSentinel {
	// 	sentinel = sg.Services[0]
	// }
	if sg.Services[1].Name == svcNameRedis {
		redis = sg.Services[1]
	}
	// if sg.Services[1].Name == svcNameSentinel {
	// 	sentinel = sg.Services[1]
	// }
	if _, ok := redis.Env["requirepass"]; !ok {
		return fmt.Errorf("redis service not provide 'requirepass' env")
	}
	return nil
}

type redisFailoverAndSecret struct {
	RedisFailover
	corev1.Secret
}

func (ro *RedisOperator) Convert(sg *apistructs.ServiceGroup) interface{} {
	svc0 := sg.Services[0]
	svc1 := sg.Services[1]
	var redis RedisSettings
	var sentinel SentinelSettings
	var redisService apistructs.Service

	scheinfo := sg.ScheduleInfo2
	scheinfo.Stateful = true
	podsLabels := []constraints.PodLabelsForAffinity{
		{PodLabels: map[string]string{"redisfailovers.databases.spotahome.com/name": sg.Dice.ID}}}
	affinity := constraintbuilders.K8S(&scheinfo, nil, podsLabels, nil).Affinity

	switch svc0.Name {
	case svcNameRedis:
		redis = ro.convertRedis(svc0, &affinity)
		redisService = svc0
	case svcNameSentinel:
		sentinel = convertSentinel(svc0, &affinity)
	}
	switch svc1.Name {
	case svcNameRedis:
		redis = ro.convertRedis(svc1, &affinity)
		redisService = svc1
	case svcNameSentinel:
		sentinel = convertSentinel(svc1, &affinity)
	}

	labels := make(map[string]string)
	annotations := make(map[string]string)
	addon.SetAddonLabelsAndAnnotations(svc0, labels, annotations)
	addon.SetAddonLabelsAndAnnotations(svc1, labels, annotations)

	password := redisService.Env["requirepass"]
	if password == apistructs.AddonRedisEmptyPassword {
		password = ""
	}

	if addonId, ok := redisService.Env["ADDON_ID"]; ok {
		labels["ADDON_ID"] = addonId
	}
	if clusterName, ok := redisService.Env["DICE_CLUSTER_NAME"]; ok {
		labels["DICE_CLUSTER_NAME"] = clusterName
	}

	rf := RedisFailover{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "databases.spotahome.com/v1",
			Kind:       "RedisFailover",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        sg.ID,
			Namespace:   genK8SNamespace(sg.Type, sg.ID),
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: RedisFailoverSpec{
			Redis:    redis,
			Sentinel: sentinel,
			Auth:     AuthSettings{SecretPath: "redis-password"},
		},
	}
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "redis-password",
			Namespace: genK8SNamespace(sg.Type, sg.ID),
		},
		Data: map[string][]byte{
			"password": []byte(password),
		},
	}
	return redisFailoverAndSecret{RedisFailover: rf, Secret: secret}

}

func (ro *RedisOperator) Create(k8syml interface{}) error {
	redisAndSecret, ok := k8syml.(redisFailoverAndSecret)
	if !ok {
		return fmt.Errorf("[BUG] this k8syml should be redisFailoverAndSecret")
	}

	redis := redisAndSecret.RedisFailover
	secret := redisAndSecret.Secret

	labels := map[string]string{
		"app.kubernetes.io/name": redis.Name,
	}
	portName := "fake-port"

	// Create Service NodePort
	rfrFakeNodePort := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("rfr-%s-np", redis.Name),
			Namespace: redis.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     portName,
					Protocol: corev1.ProtocolTCP,
					Port:     6379,
					TargetPort: intstr.IntOrString{
						IntVal: 6379,
					},
				},
			},
			Selector: labels,
			Type:     corev1.ServiceTypeNodePort,
		},
	}

	// Create Service NodePort
	rfsFakeNodePort := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("rfs-%s-np", redis.Name),
			Namespace: redis.Namespace,
			Labels: map[string]string{
				apistructs.ServiceTypeLabel: apistructs.ServiceTypeNodePort,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     portName,
					Protocol: corev1.ProtocolTCP,
					Port:     26379,
					TargetPort: intstr.IntOrString{
						IntVal: 26379,
					},
				},
			},
			Selector: labels,
			Type:     corev1.ServiceTypeNodePort,
		},
	}

	if err := ro.ns.Exists(redis.Namespace); err != nil {
		if err := ro.ns.Create(redis.Namespace, nil); err != nil {
			return err
		}
	}

	rfrNewService, err := ro.cs.CoreV1().Services(redis.Namespace).Create(context.Background(), rfrFakeNodePort, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create redis fake %s node port service, %v", redis.Name, err)
	}
	_, err = ro.cs.CoreV1().Services(redis.Namespace).Create(context.Background(), rfsFakeNodePort, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create redis fake sentinal %s node port service, %v", redis.Name, err)
	}

	defer func() {
		if err := ro.cs.CoreV1().Services(redis.Namespace).Delete(context.Background(), rfrNewService.Name,
			metav1.DeleteOptions{}); err != nil {
			logrus.Errorf("failed to delete redis fake port %s/%s node port service", redis.Namespace, rfrNewService.Name)
		}
	}()

	redis.Spec.Redis.HostNetwork = true
	redis.Spec.Redis.Port = rfrNewService.Spec.Ports[0].NodePort

	if err := ro.secret.CreateIfNotExist(&secret); err != nil {
		return err
	}
	var b bytes.Buffer
	resp, err := ro.client.Post(ro.k8s.GetK8SAddr()).
		Path(fmt.Sprintf("/apis/databases.spotahome.com/v1/namespaces/%s/redisfailovers", redis.Namespace)).
		JSONBody(redis).
		Do().
		Body(&b)
	if err != nil {
		return fmt.Errorf("failed to create redisfailover, %s/%s, err: %v", redis.Namespace, redis.Name, err)
	}
	if !resp.IsOK() {
		return fmt.Errorf("failed to create redisfailover, %s/%s, statuscode: %v, body: %v",
			redis.Namespace, redis.Name, resp.StatusCode(), b.String())
	}
	return nil
}

func (ro *RedisOperator) Inspect(sg *apistructs.ServiceGroup) (*apistructs.ServiceGroup, error) {
	namespace := genK8SNamespace(sg.Type, sg.ID)
	deploylist, err := ro.deployment.List(namespace, nil)
	if err != nil {
		return nil, err
	}
	stslist, err := ro.statefulset.List(namespace)
	if err != nil {
		return nil, err
	}
	svclist, err := ro.service.List(namespace, nil)
	if err != nil {
		return nil, err
	}

	// TODO: get redis info dynamically
	// current 2 redis + 3 sentinel
	externalSentinelEps := &apistructs.ExternalEndpoint{
		Hosts: make([]string, 0, 3),
		Ports: make([]int32, 1),
	}

	externalRedisEps := &apistructs.ExternalEndpoint{
		Hosts: make([]string, 0, 2),
		Ports: make([]int32, 1),
	}

	for _, svc := range svclist.Items {
		if strings.Contains(svc.Name, "np") && svc.Spec.Ports[0].NodePort != 0 {
			externalSentinelEps.Ports[0] = svc.Spec.Ports[0].NodePort
		}
	}

	pods, err := ro.cs.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		logrus.Errorf("failed to list namespace %s, pods: %v", namespace, err)
		return nil, err
	}

	for _, item := range pods.Items {
		if item.Spec.NodeName == "" {
			continue
		}

		node, err := ro.cs.CoreV1().Nodes().Get(context.Background(), item.Spec.NodeName, metav1.GetOptions{})
		if err != nil {
			logrus.Errorf("redis-nodeport, failed to get node %s, err: %v", item.Spec.NodeName, err)
			return nil, err
		}
		nodeIP := ""
		for _, addr := range node.Status.Addresses {
			if addr.Type == corev1.NodeInternalIP {
				nodeIP = addr.Address
			}
		}
		if nodeIP == "" {
			return nil, fmt.Errorf("node %s can't get internal ip", node.Name)
		}

		val, ok := item.Labels["app.kubernetes.io/component"]
		if !ok {
			continue
		}
		switch val {
		case "redis":
			externalRedisEps.Hosts = append(externalRedisEps.Hosts, nodeIP)
		case "sentinel":
			externalSentinelEps.Hosts = append(externalSentinelEps.Hosts, nodeIP)
		}
	}

	var redis, sentinel *apistructs.Service
	if sg.Services[0].Name == svcNameRedis {
		redis = &(sg.Services[0])
	}
	if sg.Services[1].Name == svcNameRedis {
		redis = &(sg.Services[1])
	}
	if sg.Services[0].Name == svcNameSentinel {
		sentinel = &(sg.Services[0])
	}
	if sg.Services[1].Name == svcNameSentinel {
		sentinel = &(sg.Services[1])
	}
	for _, deploy := range deploylist.Items {
		for _, cond := range deploy.Status.Conditions {
			if cond.Type == appsv1.DeploymentAvailable {
				if cond.Status == corev1.ConditionTrue {
					sentinel.Status = apistructs.StatusHealthy
				} else {
					sentinel.Status = apistructs.StatusUnHealthy
				}
			}
		}
	}
	for _, sts := range stslist.Items {
		if sts.Spec.Replicas == nil {
			redis.Status = apistructs.StatusUnknown
		} else if *sts.Spec.Replicas == sts.Status.ReadyReplicas {
			redis.Status = apistructs.StatusHealthy
		} else {
			redis.Status = apistructs.StatusUnHealthy
		}

		for _, envVar := range sts.Spec.Template.Spec.Containers[0].Env {
			if envVar.Name == "REDIS_PORT" {
				redisPort, err := strconv.Atoi(envVar.Value)
				if err == nil {
					externalRedisEps.Ports = append(externalRedisEps.Ports, int32(redisPort))
				}
			}
		}
	}

	for _, svc := range svclist.Items {
		sentinel.Vip = strutil.Join([]string{svc.Name, svc.Namespace, "svc.cluster.local"}, ".")
	}
	if redis.Status == apistructs.StatusHealthy && sentinel.Status == apistructs.StatusHealthy {
		sg.Status = apistructs.StatusHealthy
	} else {
		sg.Status = apistructs.StatusUnHealthy
	}

	redis.ExternalEndpoint = externalRedisEps
	sentinel.ExternalEndpoint = externalSentinelEps
	return sg, nil
}

func (ro *RedisOperator) Remove(sg *apistructs.ServiceGroup) error {
	k8snamespace := genK8SNamespace(sg.Type, sg.ID)
	var b bytes.Buffer
	resp, err := ro.client.Delete(ro.k8s.GetK8SAddr()).
		Path(fmt.Sprintf("/apis/databases.spotahome.com/v1/namespaces/%s/redisfailovers/%s", k8snamespace, sg.ID)).
		JSONBody(k8sapi.DeleteOptions).
		Do().
		Body(&b)
	if err != nil {
		return fmt.Errorf("failed to delele redisfailover: %s/%s, err: %v", sg.Type, sg.ID, err)
	}
	if !resp.IsOK() {
		if resp.IsNotfound() {
			return nil
		}
		return fmt.Errorf("failed to delete redisfailover: %s/%s, statuscode: %v, body: %v",
			sg.Type, sg.ID, resp.StatusCode(), b.String())
	}

	if err := ro.ns.Delete(k8snamespace); err != nil {
		logrus.Errorf("failed to delete namespace: %s: %v", k8snamespace, err)
		return nil
	}

	npName := fmt.Sprintf("rfs-%s-np", sg.ID)

	_, err = ro.cs.CoreV1().Services(k8snamespace).Get(context.Background(), npName, metav1.GetOptions{})
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to get service %s/%s, err: %v", k8snamespace, npName, err)
		}
		return nil
	}

	err = ro.cs.CoreV1().Services(k8snamespace).Delete(context.Background(), npName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete redis sentinel %s node port service, %v", npName, err)
	}

	return nil
}

// Update 支持镜像、环境变量、资源、副本数(不能为0) 等更新
// 副本数限制参考: https://github.com/spotahome/redis-operator/blob/master/api/redisfailover/v1/validate.go
func (ro *RedisOperator) Update(k8syml interface{}) error {
	// TODO: scale replicas to 0 not work，will change to default value(3)：  https://github.com/spotahome/redis-operator/blob/master/api/redisfailover/v1/validate.go
	redisAndSecret, ok := k8syml.(redisFailoverAndSecret)
	if !ok {
		return fmt.Errorf("[BUG] this k8syml should be redisFailoverAndSecret")
	}

	redis := redisAndSecret.RedisFailover
	if err := ro.ns.Exists(redis.Namespace); err != nil {
		return fmt.Errorf("namespace %s for redisoperator is not existed", redis.Namespace)
	}

	// 更新  RedisFailover (副本数不能为 0，否则将)
	var oldRedis RedisFailover
	resp, err := ro.client.Get(ro.k8s.GetK8SAddr()).
		Path(fmt.Sprintf("/apis/databases.spotahome.com/v1/namespaces/%s/redisfailovers/%s", redis.Namespace, redis.Name)).
		Do().JSON(&oldRedis)
	if err != nil {
		return fmt.Errorf("failed to update redisfailover, %s/%s, get redisfailover failed, err: %v, body: %v ", redis.Namespace, redis.Name, err, string(resp.Body()))
	}
	if !resp.IsOK() {
		return fmt.Errorf("failed to update redisfailover, %s/%s, get redisfailover is not OK, err: %v, body: %v ", redis.Namespace, redis.Name, err, string(resp.Body()))
	}

	// update redis
	if redis.Spec.Redis.Replicas >= 0 {
		oldRedis.Spec.Redis.Replicas = redis.Spec.Redis.Replicas
	}
	oldRedis.Spec.Redis.Resources = redis.Spec.Redis.Resources
	oldRedis.Spec.Redis.Image = redis.Spec.Redis.Image
	oldRedis.Spec.Redis.CustomConfig = redis.Spec.Redis.CustomConfig
	oldRedis.Spec.Redis.Command = redis.Spec.Redis.Command
	oldRedis.Spec.Redis.ShutdownConfigMap = redis.Spec.Redis.ShutdownConfigMap
	oldRedis.Spec.Redis.Exporter = redis.Spec.Redis.Exporter
	oldRedis.Spec.Redis.ImagePullSecrets = redis.Spec.Redis.ImagePullSecrets
	oldRedis.Spec.Redis.Envs = redis.Spec.Redis.Envs

	// update sentinels
	if redis.Spec.Sentinel.Replicas >= 0 {
		oldRedis.Spec.Sentinel.Replicas = redis.Spec.Sentinel.Replicas
	}
	oldRedis.Spec.Sentinel.Resources = redis.Spec.Sentinel.Resources
	oldRedis.Spec.Sentinel.Resources = redis.Spec.Sentinel.Resources
	oldRedis.Spec.Sentinel.Image = redis.Spec.Sentinel.Image
	oldRedis.Spec.Sentinel.CustomConfig = redis.Spec.Sentinel.CustomConfig
	oldRedis.Spec.Sentinel.Command = redis.Spec.Sentinel.Command
	oldRedis.Spec.Sentinel.Exporter = redis.Spec.Sentinel.Exporter
	oldRedis.Spec.Sentinel.ImagePullSecrets = redis.Spec.Sentinel.ImagePullSecrets
	oldRedis.Spec.Sentinel.Envs = redis.Spec.Sentinel.Envs

	var b bytes.Buffer
	resp, err = ro.client.Put(ro.k8s.GetK8SAddr()).
		Path(fmt.Sprintf("/apis/databases.spotahome.com/v1/namespaces/%s/redisfailovers/%s", oldRedis.Namespace, oldRedis.Name)).
		JSONBody(oldRedis).
		Do().
		Body(&b)
	if err != nil {
		return fmt.Errorf("failed to update redisfailover, %s/%s, err: %v, body: %v ", oldRedis.Namespace, oldRedis.Name, err, b.String())
	}
	if !resp.IsOK() {
		return fmt.Errorf("failed to update redisfailover, %s/%s, statuscode: %v, body: %v", oldRedis.Namespace, oldRedis.Name, resp.StatusCode(), b.String())
	}

	return nil
}

func (ro *RedisOperator) convertRedis(svc apistructs.Service, affinity *corev1.Affinity) RedisSettings {
	settings := RedisSettings{}
	settings.Affinity = affinity
	settings.Envs = svc.Env
	settings.Replicas = int32(svc.Scale)
	settings.Resources = corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			"cpu": resource.MustParse(
				fmt.Sprintf("%dm", int(1000*ro.overcommit.CPUOvercommit(svc.Resources.Cpu)))),
			"memory": resource.MustParse(
				fmt.Sprintf("%dMi", ro.overcommit.MemoryOvercommit(int(svc.Resources.Mem)))),
		},
		Limits: corev1.ResourceList{
			"cpu": resource.MustParse(
				fmt.Sprintf("%dm", int(1000*svc.Resources.Cpu))),
			"memory": resource.MustParse(
				fmt.Sprintf("%dMi", int(svc.Resources.Mem))),
		},
	}
	settings.Exporter = RedisExporter{
		Enabled: true,
		Image:   redisExporterImage,
	}
	settings.Image = svc.Image
	settings.CustomConfig = []string{}
	return settings
}

func convertSentinel(svc apistructs.Service, affinity *corev1.Affinity) SentinelSettings {
	settings := SentinelSettings{}
	settings.Affinity = affinity
	settings.Envs = svc.Env
	settings.Replicas = int32(svc.Scale)
	settings.Resources = corev1.ResourceRequirements{
		Requests: corev1.ResourceList{ // sentinel Not over-provisioned, because it should already occupy very little resources
			"cpu": resource.MustParse(
				fmt.Sprintf("%dm", int(1000*svc.Resources.Cpu))),
			"memory": resource.MustParse(
				fmt.Sprintf("%dMi", int(svc.Resources.Mem))),
		},
		Limits: corev1.ResourceList{
			"cpu": resource.MustParse(
				fmt.Sprintf("%dm", int(1000*svc.Resources.Cpu))),
			"memory": resource.MustParse(
				fmt.Sprintf("%dMi", int(svc.Resources.Mem))),
		},
	}

	configs := []string{
		"down-after-milliseconds 12000",
		"failover-timeout 12000",
	}

	password := svc.Env["requirepass"]
	if password == apistructs.AddonRedisEmptyPassword {
		password = ""
	} else {
		configs = append(configs, fmt.Sprintf("auth-pass %s", password))
	}

	settings.CustomConfig = configs
	settings.Image = svc.Image
	return settings
}

func genK8SNamespace(namespace, name string) string {
	return strutil.Concat(namespace, "--", name)
}

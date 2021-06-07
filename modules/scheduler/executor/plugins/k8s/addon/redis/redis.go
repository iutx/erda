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

package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/modules/scheduler/executor/plugins/k8s/addon"
	"github.com/erda-project/erda/modules/scheduler/schedulepolicy/constraintbuilders"
	redisfailoverv1 "github.com/erda-project/erda/pkg/clientgo/apis/redisfailover/v1"
	"github.com/erda-project/erda/pkg/clientgo/customclient"
	"github.com/erda-project/erda/pkg/clientgo/kubernetes"
	"github.com/erda-project/erda/pkg/strutil"
)

type RedisOperator struct {
	k8sClient    *kubernetes.Clientset
	customClient *customclient.Clientset
	overcommit   addon.OvercommitUtil
}

func NewRedisOperator(overcommit addon.OvercommitUtil) *RedisOperator {
	return &RedisOperator{
		overcommit: overcommit,
	}
}

func (ro *RedisOperator) IsSupported() bool {
	_, err := ro.customClient.RedisfailoverV1().RedisFailovers(metav1.NamespaceDefault).List(context.Background(), metav1.ListOptions{
		Limit: 1,
	})

	if err != nil {
		logrus.Errorf("failed to query resource redisfailover, err: %v", err)
		return false
	}

	return true
}

// Validate
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
	redisfailoverv1.RedisFailover
	corev1.Secret
}

func (ro *RedisOperator) Convert(sg *apistructs.ServiceGroup) interface{} {
	var (
		redis        redisfailoverv1.RedisSettings
		sentinel     redisfailoverv1.SentinelSettings
		redisService apistructs.Service
	)

	svc0 := sg.Services[0]
	svc1 := sg.Services[1]

	scheinfo := sg.ScheduleInfo2
	scheinfo.Stateful = true
	affinity := constraintbuilders.K8S(&scheinfo, nil, nil, nil).Affinity.NodeAffinity

	switch svc0.Name {
	case svcNameRedis:
		redis = ro.convertRedis(svc0, affinity)
		redisService = svc0
	case svcNameSentinel:
		sentinel = convertSentinel(svc0, affinity)
	}
	switch svc1.Name {
	case svcNameRedis:
		redis = ro.convertRedis(svc1, affinity)
		redisService = svc1
	case svcNameSentinel:
		sentinel = convertSentinel(svc1, affinity)
	}

	rf := redisfailoverv1.RedisFailover{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "databases.spotahome.com/v1",
			Kind:       "RedisFailover",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      sg.ID,
			Namespace: genK8SNamespace(sg.Type, sg.ID),
		},
		Spec: redisfailoverv1.RedisFailoverSpec{
			Redis:    redis,
			Sentinel: sentinel,
			Auth:     redisfailoverv1.AuthSettings{SecretPath: "redis-password"},
		},
	}
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "redis-password",
			Namespace: genK8SNamespace(sg.Type, sg.ID),
		},
		Data: map[string][]byte{
			"password": []byte(redisService.Env["requirepass"]),
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

	_, err := ro.k8sClient.CoreV1().Namespaces().Create(context.Background(), &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: redis.Namespace,
		},
	}, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}

	_, err = ro.k8sClient.CoreV1().Secrets(secret.Namespace).Create(context.Background(), &secret, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}

	_, err = ro.customClient.RedisfailoverV1().RedisFailovers(redis.Namespace).Create(context.Background(), &redis)
	if err != nil {
		logrus.Errorf("failed to create redisfailover, %s/%s, err: %v", redis.Namespace, redis.Name, err)
		return err
	}

	return nil
}

func (ro *RedisOperator) Inspect(sg *apistructs.ServiceGroup) (*apistructs.ServiceGroup, error) {
	nsName := genK8SNamespace(sg.Type, sg.ID)

	deployList, err := ro.k8sClient.AppsV1().Deployments(nsName).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	stsList, err := ro.k8sClient.AppsV1().StatefulSets(nsName).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	svcList, err := ro.k8sClient.CoreV1().Services(nsName).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
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

	if sentinel == nil || redis == nil {
		return nil, errors.New("sentinel or redis is nil")
	}

	for _, deploy := range deployList.Items {
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
	for _, sts := range stsList.Items {
		if sts.Spec.Replicas == nil {
			redis.Status = apistructs.StatusUnknown
		} else if *sts.Spec.Replicas == sts.Status.ReadyReplicas {
			redis.Status = apistructs.StatusHealthy
		} else {
			redis.Status = apistructs.StatusUnHealthy
		}
	}

	for _, svc := range svcList.Items {
		sentinel.Vip = strutil.Join([]string{svc.Name, svc.Namespace, "svc.cluster.local"}, ".")
	}
	if redis.Status == apistructs.StatusHealthy && sentinel.Status == apistructs.StatusHealthy {
		sg.Status = apistructs.StatusHealthy
	} else {
		sg.Status = apistructs.StatusUnHealthy
	}
	return sg, nil
}

func (ro *RedisOperator) Remove(sg *apistructs.ServiceGroup) error {
	nsName := genK8SNamespace(sg.Type, sg.ID)

	err := ro.customClient.RedisfailoverV1().RedisFailovers(nsName).Delete(context.Background(), sg.ID, &metav1.DeleteOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("failed to delele redisfailover: %s/%s, err: %v", sg.Type, sg.ID, err)
	}

	if err := ro.k8sClient.CoreV1().Namespaces().Delete(context.Background(), nsName, metav1.DeleteOptions{}); err != nil {
		logrus.Errorf("failed to delete namespace: %s: %v", nsName, err)
		return nil
	}

	return nil
}

func (ro *RedisOperator) Update(k8syml interface{}) error {
	// TODO:
	return fmt.Errorf("redisoperator not impl Update yet")
}

func (ro *RedisOperator) convertRedis(svc apistructs.Service, affinity *corev1.NodeAffinity) redisfailoverv1.RedisSettings {
	settings := redisfailoverv1.RedisSettings{}
	settings.Affinity = &corev1.Affinity{NodeAffinity: affinity}
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
	settings.Image = svc.Image
	return settings
}

func convertSentinel(svc apistructs.Service, affinity *corev1.NodeAffinity) redisfailoverv1.SentinelSettings {
	settings := redisfailoverv1.SentinelSettings{}
	settings.Affinity = &corev1.Affinity{NodeAffinity: affinity}
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
	settings.CustomConfig = []string{
		fmt.Sprintf("auth-pass %s", svc.Env["requirepass"]),
		"down-after-milliseconds 12000",
		"failover-timeout 12000",
	}
	settings.Image = svc.Image
	return settings
}

func genK8SNamespace(namespace, name string) string {
	return strutil.Concat(namespace, "--", name)
}

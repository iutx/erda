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

package scheme

import (
	rocketmqv1alpha1 "erda.cloud/rocketmq/api/v1alpha1"
	elasticsearchv1 "github.com/elastic/cloud-on-k8s/v2/pkg/apis/elasticsearch/v1"
	kedav1alpha1 "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	sparkoperatorv1beta2 "github.com/kubeflow/spark-operator/pkg/apis/sparkoperator.k8s.io/v1beta2"
	flinkoperatoryv1beta1 "github.com/spotify/flink-on-k8s-operator/apis/flinkcluster/v1beta1"
	istioconfigv1alpha2 "istio.io/client-go/pkg/apis/config/v1alpha2"
	istionetworkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	istionetworkingv1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	istiorbacv1alpha1 "istio.io/client-go/pkg/apis/rbac/v1alpha1"
	istiosecv1beta1 "istio.io/client-go/pkg/apis/security/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/apimachinery/pkg/runtime"
	k8sschema "k8s.io/client-go/kubernetes/scheme"

	openyurtv1alpha1 "github.com/erda-project/erda/pkg/k8sclient/apis/openyurt/v1alpha1"
)

// LocalSchemeBuilder register crd scheme
var LocalSchemeBuilder = runtime.SchemeBuilder{
	openyurtv1alpha1.AddToScheme,
	k8sschema.AddToScheme,
	elasticsearchv1.AddToScheme,
	sparkoperatorv1beta2.AddToScheme,
	istioconfigv1alpha2.AddToScheme,
	istionetworkingv1beta1.AddToScheme,
	istionetworkingv1alpha3.AddToScheme,
	istiorbacv1alpha1.AddToScheme,
	istiosecv1beta1.AddToScheme,
	flinkoperatoryv1beta1.AddToScheme,
	apiextensions.AddToScheme,
	rocketmqv1alpha1.AddToScheme,
	kedav1alpha1.AddToScheme,
}

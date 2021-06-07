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

package addon

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/erda-project/erda/apistructs"
)

type AddonOperator interface {
	// Determine whether this cluster supports addon operator
	IsSupported() bool
	// Verify the legality of sg converted from diceyml
	Validate(*apistructs.ServiceGroup) error
	// Convert sg to cr. Because of the different definitions of cr, use interface() here
	Convert(*apistructs.ServiceGroup) interface{}
	// the cr converted by Convert in k8s deploying
	Create(interface{}) error
	// Check running status
	Inspect(*apistructs.ServiceGroup) (*apistructs.ServiceGroup, error)

	Remove(*apistructs.ServiceGroup) error

	Update(interface{}) error
}

type K8SUtil interface {
	GetK8SAddr() string
}

type ImageSecretUtil interface {
	//The secret used to pull the mirror under the default namespace is copied to the current ns,
	// Then add this secret to the default sa of this ns
	NewImageSecret(ns string) error
}

type HealthcheckUtil interface {
	NewHealthcheckProbe(*apistructs.Service) *corev1.Probe
}

type OvercommitUtil interface {
	CPUOvercommit(limit float64) float64
	MemoryOvercommit(limit int) int
}

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

package processors_test

import (
	"encoding/json"
	"fmt"

	"github.com/erda-project/erda-proto-go/core/monitor/metric/pb"
	"github.com/erda-project/erda/modules/extensions/loghub/metrics/analysis/processors"
	_ "github.com/erda-project/erda/modules/extensions/loghub/metrics/analysis/processors/regex" //
)

func ExampleProcessors() {
	var (
		scopeID string = "terminus"
		tags           = map[string]string{
			"dice_org_id":         "1",
			"dice_application_id": "2",
			"dice_service_name":   "abc",
		}
		metricName = "test_metric"
	)
	cfg, _ := json.Marshal(map[string]interface{}{
		"pattern": "(d+)",
		"keys": []*pb.FieldDefine{
			{
				Key:  "ip",
				Type: "string",
			},
		},
	})
	ps := processors.New()
	err := ps.Add(scopeID, tags, metricName, "regexp", cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	list := ps.Find("", scopeID, map[string]string{
		"dice_org_id":         "1",
		"dice_application_id": "2",
	})
	if len(list) != 0 {
		fmt.Println("Find error")
		return
	}
	list = ps.Find("", scopeID, map[string]string{
		"dice_org_id":         "1",
		"dice_application_id": "3",
		"dice_service_name":   "abc",
	})
	if len(list) != 0 {
		fmt.Println("Find error")
		return
	}
	list = ps.Find("", scopeID, map[string]string{
		"dice_org_id":         "1",
		"dice_application_id": "2",
		"dice_service_name":   "abc",
	})
	if len(list) <= 0 {
		fmt.Println("Find error")
		return
	}
	fmt.Printf("Find %d\n", len(list))

	// Output:
	// Find 1
}

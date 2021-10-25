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

package resource

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/erda-project/erda-infra/providers/i18n"
	"github.com/erda-project/erda/apistructs"
	_interface "github.com/erda-project/erda/modules/cmp/cmp_interface"
)

type nopTranslator struct{}

func (t nopTranslator) Get(lang i18n.LanguageCodes, key, def string) string { return key }

func (t nopTranslator) Text(lang i18n.LanguageCodes, key string) string { return key }

func (t nopTranslator) Sprintf(lang i18n.LanguageCodes, key string, args ...interface{}) string {
	return fmt.Sprintf(key, args...)
}

func TestResource_getGauge(t *testing.T) {
	type fields struct {
		Ctx    context.Context
		Server _interface.Provider
		I18N   i18n.Translator
		Lang   i18n.LanguageCodes
	}
	type args struct {
		request *apistructs.GaugeRequest
		resp    *apistructs.ResourceResp
	}
	request := &apistructs.GaugeRequest{}
	resp := &apistructs.ResourceResp{}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantData map[string]*GaugeData
	}{
		// TODO: Add test cases.
		{
			name: "test",
			fields: fields{
				I18N: nopTranslator{},
			},
			args: args{
				request: request,
				resp:    resp,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				I18N: tt.fields.I18N,
			}
			if gotData := r.getGauge(tt.args.request, tt.args.resp); !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("getGauge() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

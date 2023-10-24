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

package sheet_issue

import (
	"fmt"
	"sort"

	"github.com/erda-project/erda-proto-go/dop/issue/core/pb"
	"github.com/erda-project/erda/internal/apps/dop/providers/issue/core/query/issueexcel/vars"
	"github.com/erda-project/erda/pkg/excel"
)

// convertOldIssueSheet
// old 也有两个版本：
// - 最老的版本，只有 18 个基础字段
// - 其他版本，有 21 个基础字段 + 自定义字段
func convertOldIssueSheet(data *vars.DataForFulfill, sheet [][]string) ([]vars.IssueSheetModel, error) {
	// convert by column fixed index
	info := NewIssueSheetModelCellInfoByColumns()
	m := info.M
	addM := func(m map[IssueSheetColumnUUID]excel.Column, uuid IssueSheetColumnUUID, s string) {
		uuid.AutoComplete()
		info.Add(uuid, s)
	}
	// handle custom fields
	if len(sheet) == 0 {
		return nil, nil
	}
	// remove empty rows
	removeEmptySheetRows(&sheet)
	// columnLen, 计算到标题行的第一个非空 cell，因为有些单元格是手动删除过数据的
	var columnLen int
	for _, cell := range sheet[0] {
		if cell == "" {
			break
		}
		columnLen++
	}
	switch true {
	case columnLen >= oldExcelFormatCustomFieldRowColumnIndexFrom:
	case columnLen == oldOldExcelFormatColumnLen:
	default:
		return nil, fmt.Errorf("invalid column len: %d, please check excel", columnLen)
	}
	// auto fill empty cells
	autoFillEmptyRowCells(&sheet, columnLen)
	// try to match custom field name to issue type, because the order of custom field is not fixed
	var customFieldNames []string
	var columnIndexAndPropertyTypeMap map[int]pb.PropertyIssueTypeEnum_PropertyIssueType
	if columnLen > oldExcelFormatCustomFieldRowColumnIndexFrom {
		customFieldNames = sheet[0][oldExcelFormatCustomFieldRowColumnIndexFrom:]
		_columnIndexAndPropertyTypeMap, err := tryToMatchCustomFieldNameToIssueType(customFieldNames, data.CustomFieldMapByTypeName)
		if err != nil {
			return nil, fmt.Errorf("failed to match custom field name to issue type, err: %v", err)
		}
		columnIndexAndPropertyTypeMap = _columnIndexAndPropertyTypeMap
	}
	for rowIdx, row := range sheet {
		if rowIdx == 0 {
			continue
		}
		// get issue type first
		issueType, err := parseStringIssueType(data, sheet[rowIdx][oldExcelFormatIndexOfIssueType])
		if err != nil {
			return nil, fmt.Errorf("failed to parse issue type, err: %v", err)
		}
		// get columnLen
		columnLen := len(row)
		for columnIdx := 0; columnIdx < oldExcelFormatCustomFieldRowColumnIndexFrom; columnIdx++ {
			if columnIdx >= columnLen {
				continue
			}
			s := row[columnIdx]
			switch columnIdx {
			case 0: // ID
				addM(m, NewIssueSheetColumnUUID(fieldCommon, fieldID), s)
			case 1: // Title
				addM(m, NewIssueSheetColumnUUID(fieldCommon, fieldIssueTitle), s)
			case 2: // Content
				addM(m, NewIssueSheetColumnUUID(fieldCommon, fieldContent), s)
			case 3: // State
				addM(m, NewIssueSheetColumnUUID(fieldCommon, FieldState), s)
			case 4: // Creator
				addM(m, NewIssueSheetColumnUUID(fieldCommon, fieldCreatorName), s)
			case 5: // Assignee
				addM(m, NewIssueSheetColumnUUID(fieldCommon, fieldAssigneeName), s)
			case 6: // Owner
				addM(m, NewIssueSheetColumnUUID(fieldBugOnly, fieldOwnerName), s)
			case 7: // TaskType or BugSource
				switch *issueType {
				case pb.IssueTypeEnum_TASK:
					addM(m, NewIssueSheetColumnUUID(fieldTaskOnly, fieldTaskType), s)
				case pb.IssueTypeEnum_BUG:
					addM(m, NewIssueSheetColumnUUID(fieldBugOnly, fieldSource), s)
				}
			case 8: // Priority
				addM(m, NewIssueSheetColumnUUID(fieldCommon, fieldPriority), s)
			case 9: // IterationName
				addM(m, NewIssueSheetColumnUUID(fieldCommon, fieldIterationName), s)
			case 10: // Complexity
				addM(m, NewIssueSheetColumnUUID(fieldCommon, fieldComplexity), s)
			case 11: // Severity
				addM(m, NewIssueSheetColumnUUID(fieldCommon, fieldSeverity), s)
			case 12: // Labels
				addM(m, NewIssueSheetColumnUUID(fieldCommon, fieldLabels), s)
			case 13: // IssueType
				addM(m, NewIssueSheetColumnUUID(fieldCommon, fieldIssueType), s)
			case 14: // PlanFinishedAt
				addM(m, NewIssueSheetColumnUUID(fieldCommon, fieldPlanFinishedAt), s)
			case 15: // CreatedAt
				addM(m, NewIssueSheetColumnUUID(fieldCommon, fieldCreatedAt), s)
			case 16: // ConnectionIssueIDs
				addM(m, NewIssueSheetColumnUUID(fieldCommon, fieldConnectionIssueIDs), s)
			case 17: // EstimateTime
				addM(m, NewIssueSheetColumnUUID(fieldCommon, fieldEstimateTime), s)
			case 18: // FinishedAt
				addM(m, NewIssueSheetColumnUUID(fieldCommon, fieldFinishAt), s)
			case 19: // StartAt
				addM(m, NewIssueSheetColumnUUID(fieldCommon, fieldPlanStartedAt), s)
			case 20: // ReopenCount
				addM(m, NewIssueSheetColumnUUID(fieldBugOnly, fieldReopenCount), s)
			default:
			}
		}
		// handle custom fields
		for i, propertyType := range columnIndexAndPropertyTypeMap {
			s := row[i+oldExcelFormatCustomFieldRowColumnIndexFrom]
			switch propertyType {
			case pb.PropertyIssueTypeEnum_REQUIREMENT:
				addM(m, NewIssueSheetColumnUUID(fieldRequirementOnly, fieldCustomFields, customFieldNames[i]), s)
			case pb.PropertyIssueTypeEnum_TASK:
				addM(m, NewIssueSheetColumnUUID(fieldTaskOnly, fieldCustomFields, customFieldNames[i]), s)
			case pb.PropertyIssueTypeEnum_BUG:
				addM(m, NewIssueSheetColumnUUID(fieldBugOnly, fieldCustomFields, customFieldNames[i]), s)
			}
		}
	}
	models, err := decodeMapToIssueSheetModel(data, info)
	if err != nil {
		return nil, fmt.Errorf("failed to decode old excel format map to issue sheet model, err: %v", err)
	}
	return models, nil
}

const (
	oldExcelFormatCustomFieldRowColumnIndexFrom = 21
	oldExcelFormatIndexOfIssueType              = 13
	oldOldExcelFormatColumnLen                  = 18
)

// tryToMatchCustomFieldNameToIssueType
// 由于之前全量导出时:
// - 没有区分类型
// - 而且是根据类型 WaitGroup 并发执行，没有顺序
// - 同一类型内保证了顺序
// 如果存在一个自定义字段被多个 issue type 使用，只能尽可能匹配
// 如果字段按顺序都能匹配上（字段名、顺序），则匹配成功
// 特殊情况，如果只有一个自定义字段，且这个字段被 3 个类型都引用了，则无法保证正确性。解决方案：用户可以手动调整模板字段顺序，原则就是 需求 > 任务 > 缺陷
func tryToMatchCustomFieldNameToIssueType(cfNames []string, customFieldMap map[pb.PropertyIssueTypeEnum_PropertyIssueType]map[string]*pb.IssuePropertyIndex) (
	map[int]pb.PropertyIssueTypeEnum_PropertyIssueType, error) {

	genOrders := func(typeOrders ...pb.PropertyIssueTypeEnum_PropertyIssueType) []*pb.IssuePropertyIndex {
		var result []*pb.IssuePropertyIndex
		for _, t := range typeOrders {
			var cfs []*pb.IssuePropertyIndex
			for _, cf := range customFieldMap[t] {
				cfs = append(cfs, cf)
			}
			sort.SliceStable(cfs, func(i, j int) bool {
				return cfs[i].Index < cfs[j].Index
			})
			result = append(result, cfs...)
		}
		return result
	}

	// 所有可能的顺序
	possibleCfNameOrders := [][]*pb.IssuePropertyIndex{
		// all types
		genOrders(pb.PropertyIssueTypeEnum_REQUIREMENT, pb.PropertyIssueTypeEnum_TASK, pb.PropertyIssueTypeEnum_BUG),
		genOrders(pb.PropertyIssueTypeEnum_REQUIREMENT, pb.PropertyIssueTypeEnum_BUG, pb.PropertyIssueTypeEnum_TASK),
		genOrders(pb.PropertyIssueTypeEnum_TASK, pb.PropertyIssueTypeEnum_REQUIREMENT, pb.PropertyIssueTypeEnum_BUG),
		genOrders(pb.PropertyIssueTypeEnum_TASK, pb.PropertyIssueTypeEnum_BUG, pb.PropertyIssueTypeEnum_REQUIREMENT),
		genOrders(pb.PropertyIssueTypeEnum_BUG, pb.PropertyIssueTypeEnum_REQUIREMENT, pb.PropertyIssueTypeEnum_TASK),
		genOrders(pb.PropertyIssueTypeEnum_BUG, pb.PropertyIssueTypeEnum_TASK, pb.PropertyIssueTypeEnum_REQUIREMENT),
		// two types
		genOrders(pb.PropertyIssueTypeEnum_REQUIREMENT, pb.PropertyIssueTypeEnum_TASK),
		genOrders(pb.PropertyIssueTypeEnum_REQUIREMENT, pb.PropertyIssueTypeEnum_BUG),
		genOrders(pb.PropertyIssueTypeEnum_TASK, pb.PropertyIssueTypeEnum_REQUIREMENT),
		genOrders(pb.PropertyIssueTypeEnum_TASK, pb.PropertyIssueTypeEnum_BUG),
		genOrders(pb.PropertyIssueTypeEnum_BUG, pb.PropertyIssueTypeEnum_REQUIREMENT),
		genOrders(pb.PropertyIssueTypeEnum_BUG, pb.PropertyIssueTypeEnum_TASK),
		// one type
		genOrders(pb.PropertyIssueTypeEnum_REQUIREMENT),
		genOrders(pb.PropertyIssueTypeEnum_TASK),
		genOrders(pb.PropertyIssueTypeEnum_BUG),
	}

	// 顺序匹配
	var found bool
	var matchedOrders []*pb.IssuePropertyIndex
	for _, expectOrders := range possibleCfNameOrders {
		if len(expectOrders) != len(cfNames) {
			continue
		}
		allMatched := true
		for i := range expectOrders {
			expect := expectOrders[i]
			actual := cfNames[i]
			if expect.PropertyName != actual {
				allMatched = false
				break
			}
		}
		if allMatched {
			found = true
			matchedOrders = expectOrders
			break
		}
		continue
	}
	if !found {
		return nil, fmt.Errorf("custom field name order not matched")
	}

	// 匹配成功，返回字段索引和 issue type 的映射关系
	result := make(map[int]pb.PropertyIssueTypeEnum_PropertyIssueType)
	for i, cf := range matchedOrders {
		result[i] = cf.PropertyIssueType
	}
	return result, nil
}

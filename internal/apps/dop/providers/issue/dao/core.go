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

package dao

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/erda-project/erda-proto-go/dop/issue/core/pb"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/internal/apps/dop/providers/issue/core/common"
	"github.com/erda-project/erda/pkg/arrays"
	"github.com/erda-project/erda/pkg/database/dbengine"
	"github.com/erda-project/erda/pkg/strutil"
)

const (
	// day = 8 * 60(minute)
	day = 480
)

// Issue .
type Issue struct {
	dbengine.BaseModel

	PlanStartedAt  *time.Time // 计划开始时间
	PlanFinishedAt *time.Time // 计划结束时间
	ProjectID      uint64     // 所属项目 ID
	IterationID    int64      // 所属迭代 ID
	AppID          *uint64    // deprecated, 所属应用 ID
	RequirementID  *int64     // deprecated, 所属需求 ID
	Type           string     // issue 类型
	Title          string     // 标题
	Content        string     // 内容
	State          int64      // 状态
	Priority       string     // 优先级
	Complexity     string     // 复杂度
	Severity       string     // 严重程度
	Creator        string     // issue 创建者 ID
	Assignee       string     // 分配到 issue 的人，即当前处理人
	Source         string     // issue创建的来源，目前只有工单使用
	ManHour        string     // 工时信息
	External       bool       // Deprecated. 用来区分是通过ui还是bundle创建的 // 2023-05-30 Deprecated. Always true now
	Deleted        bool       // 是否已删除
	Stage          string     // bug阶段 or 任务类型 的值
	Owner          string     // 负责人

	FinishTime   *time.Time // 实际结束时间
	ExpiryStatus ExpireType
	ReopenCount  int
	StartTime    *time.Time
}

type ExpireType string

func (e ExpireType) String() string {
	return string(e)
}

const (
	ExpireTypeUndefined      ExpireType = "Undefined"
	ExpireTypeExpired        ExpireType = "Expired"
	ExpireTypeExpireIn1Day   ExpireType = "ExpireIn1Day"
	ExpireTypeExpireIn2Days  ExpireType = "ExpireIn2Days"
	ExpireTypeExpireIn7Days  ExpireType = "ExpireIn7Days"
	ExpireTypeExpireIn30Days ExpireType = "ExpireIn30Days"
	ExpireTypeExpireInFuture ExpireType = "ExpireInFuture"
	ExpireTypeUnfinished     ExpireType = "Unfinished"
)

var ExpireTypes = []ExpireType{ExpireTypeUndefined, ExpireTypeExpired, ExpireTypeExpireIn1Day, ExpireTypeExpireIn2Days, ExpireTypeExpireIn7Days, ExpireTypeExpireIn30Days, ExpireTypeExpireInFuture}

func GetExpiryStatus(planFinishedAt *time.Time, timeBase time.Time) ExpireType {
	if planFinishedAt == nil {
		return ExpireTypeUndefined
	}
	if planFinishedAt.Before(timeBase) {
		return ExpireTypeExpired
	} else if planFinishedAt.Before(timeBase.Add(1 * 24 * time.Hour)) {
		return ExpireTypeExpireIn1Day
	} else if planFinishedAt.Before(timeBase.Add(2 * 24 * time.Hour)) {
		return ExpireTypeExpireIn2Days
	} else if planFinishedAt.Before(timeBase.Add(7 * 24 * time.Hour)) {
		return ExpireTypeExpireIn7Days
	} else if planFinishedAt.Before(timeBase.Add(30 * 24 * time.Hour)) {
		return ExpireTypeExpireIn30Days
	} else {
		return ExpireTypeExpireInFuture
	}
}

func (Issue) TableName() string {
	return "dice_issues"
}

type IssueSummary struct {
	Total       int                  `json:"total,omitempty"`
	IssueType   apistructs.IssueType `json:"issue_type,omitempty"`
	State       int64                `json:"state,omitempty"`
	IterationID int64                `json:"iteration_id,omitempty"`
}

// GetCanUpdateFields 获取所有可以被主动更新的字段
func (i *Issue) GetCanUpdateFields() map[string]interface{} {
	return map[string]interface{}{
		"plan_started_at":  i.PlanStartedAt,
		"plan_finished_at": i.PlanFinishedAt,
		"app_id":           i.AppID,
		"title":            i.Title,
		"content":          i.Content,
		"state":            i.State,
		"priority":         i.Priority,
		"complexity":       i.Complexity,
		"severity":         i.Severity,
		"assignee":         i.Assignee,
		"iteration_id":     i.IterationID,
		"man_hour":         i.ManHour,
		"source":           i.Source,
		"stage":            i.Stage,
		"owner":            i.Owner,
		"finish_time":      i.FinishTime,
		"expiry_status":    i.ExpiryStatus,
	}
}

// CreateIssue 创建
func (client *DBClient) CreateIssue(issue *Issue) error {
	return client.Create(issue).Error
}

// UpdateIssue 更新
func (client *DBClient) UpdateIssue(id uint64, fields map[string]interface{}) error {
	issue := Issue{}
	issue.ID = id

	return client.Model(&issue).Updates(fields).Error
}

// UpdateIssues 批量更新 issue
func (client *DBClient) UpdateIssues(requirementID uint64, fields map[string]interface{}) error {
	return client.Model(Issue{}).
		Where("requirement_id = ?", requirementID).
		Where("type = ?", apistructs.IssueTypeTask).
		Updates(fields).Error
}

func (client *DBClient) BatchUpdateIssueIterationIDByIterationID(iterationID uint64, ID int64) error {
	return client.Model(Issue{}).Where("deleted = 0").Where("iteration_id = ?", iterationID).Updates(map[string]interface{}{"iteration_id": ID}).Error
}

// UpdateIssueType 转换issueType
func (client *DBClient) UpdateIssueType(issue *Issue) error {
	return client.Model(Issue{}).Save(issue).Error
}

// GetBatchUpdateIssues 获取待批量更新记录，生成活动记录
func (client *DBClient) GetBatchUpdateIssues(req *pb.BatchUpdateIssueRequest) ([]Issue, error) {
	var issues []Issue
	sql := client.Model(Issue{}).Where("type = ?", req.Type.String()).Where("deleted = 0")
	if len(req.CurrentIterationIDs) > 0 {
		sql = sql.Where("iteration_id in (?)", req.CurrentIterationIDs)
	}
	if !req.All && len(req.Ids) > 0 {
		sql = sql.Where("id in (?)", req.Ids)
	}
	if req.All && req.Mine {
		sql = sql.Where("assignee = ?", req.IdentityInfo.UserID)
	}
	if err := sql.Find(&issues).Error; err != nil {
		return nil, err
	}

	return issues, nil
}

// BatchUpdateIssues 批量更新 issue
func (client *DBClient) BatchUpdateIssues(req *pb.BatchUpdateIssueRequest) error {
	sql := client.Model(Issue{}).Where("type = ?", req.Type.String())
	if len(req.CurrentIterationIDs) > 0 {
		sql = sql.Where("iteration_id in (?)", req.CurrentIterationIDs)
	}
	if req.All && req.Mine {
		sql = sql.Where("assignee = ?", req.IdentityInfo.UserID)
	}
	if !req.All && len(req.Ids) > 0 {
		sql = sql.Where("id in (?)", req.Ids)
	}

	var issue Issue
	if req.State != 0 {
		issue.State = req.State
	}
	if req.Owner != "" {
		issue.Owner = req.Owner
	}
	if req.Assignee != "" {
		issue.Assignee = req.Assignee
	}
	if req.NewIterationID != 0 {
		issue.IterationID = req.NewIterationID
	}

	return sql.Updates(issue).Error
}

// DeleteIssue 删除
func (client *DBClient) DeleteIssue(id uint64) error {
	return client.Model(Issue{}).Where("id = ?", id).Update("deleted", 1).Error
}

// BatchDeleteIssues 批量删除
func (client *DBClient) BatchDeleteIssues(ids []uint64) error {
	return client.Model(Issue{}).Where("id IN (?)", ids).Update("deleted", 1).Error
}

func (client *DBClient) BatchDeleteIssueByIterationID(iterationID uint64) error {
	sql := fmt.Sprintf(`
		UPDATE dice_issues AS dis
		LEFT JOIN dice_issue_testcase_relations AS ditr ON ditr.issue_id = dis.id
		LEFT JOIN dice_issue_relation AS dir ON dir.issue_id = dis.id
		LEFT JOIN dice_issue_property_relation AS dipr ON dipr.issue_id = dis.id
		LEFT JOIN erda_issue_state_transition AS eist ON eist.issue_id = dis.id
		SET dis.deleted = 1
		WHERE dis.iteration_id = ?;
	`)
	return client.Exec(sql, iterationID).Error
}

// DeleteIssuesByRequirement .
func (client *DBClient) DeleteIssuesByRequirement(requirementID uint64) error {
	return client.Model(Issue{}).Where("requirement_id = ?", requirementID).Update("deleted", 1).Error
}

// IssueStateCount 需求下任务统计
func (client *DBClient) IssueStateCount(requirementIDs []int64) ([]apistructs.RequirementGroupResult, error) {
	var result []apistructs.RequirementGroupResult
	if err := client.Raw("select requirement_id as id, state, COUNT(*) as count from dice_issues where type = ? and requirement_id in (?) and deleted = ? group by requirement_id, state",
		apistructs.IssueTypeTask, requirementIDs, 0).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// PagingIssues 分页查询issue
// queryIDs = true 表示不管req.IDs是否为空都要根据ID查询
func (client *DBClient) PagingIssues(req pb.PagingIssueRequest, queryIDs bool) ([]Issue, uint64, error) {
	var (
		total  uint64
		issues []Issue
	)
	// 有了待办事件后，迭代id一定不能等于0
	cond := Issue{}
	if req.ProjectID > 0 {
		cond.ProjectID = req.ProjectID
	}
	if req.AppID != nil && *req.AppID > 0 {
		cond.AppID = req.AppID
	}
	if req.RequirementID != nil && *req.RequirementID > 0 {
		cond.RequirementID = req.RequirementID
	}
	sql := client.DB
	if req.CustomPanelID != 0 {
		joinSQL := "LEFT OUTER JOIN dice_issue_panel on dice_issues.id=dice_issue_panel.issue_id"
		sql = sql.Joins(joinSQL).Where("relation = ?", req.CustomPanelID)
	}
	if req.NotIncluded {
		sql = sql.Joins("LEFT JOIN dice_issue_relation b ON dice_issues.id = b.related_issue and b.type = ?", apistructs.IssueRelationInclusion).
			Where("b.id IS NULL")
	}
	sql = sql.Where(cond).Where("deleted = ?", 0)
	if len(req.IDs) > 0 {
		sql = sql.Where("dice_issues.id IN (?)", req.IDs)
	} else if queryIDs {
		return nil, 0, nil
	}
	if len(req.IterationIDs) > 0 {
		sql = sql.Where("iteration_id in (?)", req.IterationIDs)
	}
	if len(req.ProjectIDs) > 0 {
		sql = sql.Where("project_id in (?)", req.ProjectIDs)
	}
	if len(req.Type) > 0 {
		sql = sql.Where("dice_issues.type IN (?)", req.Type)
	}
	if len(req.Creator) > 0 {
		sql = sql.Where("creator IN (?)", req.Creator)
	}
	if len(req.Assignee) > 0 {
		sql = sql.Where("assignee IN (?)", req.Assignee)
	}
	if len(req.Priority) > 0 {
		sql = sql.Where("priority IN (?)", req.Priority)
	}
	if len(req.Complexity) > 0 {
		sql = sql.Where("complexity IN (?)", req.Complexity)
	}
	if len(req.Severity) > 0 {
		sql = sql.Where("severity IN (?)", req.Severity)
	}
	if len(req.State) > 0 {
		sql = sql.Where("state IN (?)", req.State)
	}
	if len(req.StateBelongs) > 0 {
		sql = sql.Joins(joinState).Where("dice_issue_state.belong IN (?)", req.StateBelongs)
	}
	if len(req.Owner) > 0 {
		sql = sql.Where("owner IN (?)", req.Owner)
	}
	if len(req.BugStage) > 0 {
		sql = sql.Where("stage IN (?)", req.BugStage)
	} else if len(req.TaskType) > 0 {
		sql = sql.Where("stage IN (?)", req.TaskType)
	}
	if req.ReopenedCountGte > 0 {
		sql = sql.Where("dice_issues.reopen_count >= ?", req.ReopenedCountGte)
	}
	if len(req.ExceptIDs) > 0 {
		sql = sql.Where("dice_issues.id NOT IN (?)", req.ExceptIDs)
	}
	if req.StartCreatedAt > 0 {
		startCreatedAt := time.Unix(req.StartCreatedAt/1000, 0)
		sql = sql.Where("dice_issues.created_at >= ?", startCreatedAt)
	}
	if req.EndCreatedAt > 0 {
		endCreatedAt := time.Unix(req.EndCreatedAt/1000, 0)
		sql = sql.Where("dice_issues.created_at <= ?", endCreatedAt)
	}
	if req.IsEmptyPlanFinishedAt {
		sql = sql.Where("plan_finished_at IS NULL")
	}
	if req.StartFinishedAt > 0 && !req.IsEmptyPlanFinishedAt {
		startFinishedAt := time.Unix(req.StartFinishedAt/1000, 0)
		sql = sql.Where("plan_finished_at >= ?", startFinishedAt)
	}
	if req.EndFinishedAt > 0 && !req.IsEmptyPlanFinishedAt {
		endFinishedAt := time.Unix(req.EndFinishedAt/1000, 0)
		sql = sql.Where("plan_finished_at <= ?", endFinishedAt)
	}

	if req.StartClosedAt > 0 {
		startClosedAt := time.Unix(req.StartClosedAt/1000, 0)
		sql = sql.Where("finish_time >= ?", startClosedAt)
	}
	if req.EndClosedAt > 0 {
		endClosedAt := time.Unix(req.EndClosedAt/1000, 0)
		sql = sql.Where("finish_time <= ?", endClosedAt)
	}

	if req.Title != "" {
		title := strings.ReplaceAll(req.Title, "%", "\\%")
		if _, err := strutil.Atoi64(title); err == nil {
			sql = sql.Where("title LIKE ? OR dice_issues.id LIKE ?", "%"+title+"%", "%"+title+"%")
		} else {
			sql = sql.Where("title LIKE ?", "%"+title+"%")
		}
	}
	if req.Source != "" {
		sql = sql.Where("source LIKE ?", "%"+req.Source+"%")
	}
	if req.External {
		sql = sql.Where("external = 1")
	} else {
		sql = sql.Where("external = 0")
	}
	if req.OrderBy != "" {
		if req.Asc {
			sql = sql.Order(fmt.Sprintf("%s", req.OrderBy))
		} else {
			sql = sql.Order(fmt.Sprintf("%s DESC", req.OrderBy))
		}
	} else {
		sql = sql.Order("dice_issues.id DESC")
	}

	if len(req.Participant) > 0 {
		sql = sql.Joins(joinParticipant).Where("erda_issue_subscriber.user_id in (?)", req.Participant).Select("distinct dice_issues.*")
	}
	if req.OnlyIdResult {
		sql = sql.Select("dice_issues.id")
	}

	offset := (req.PageNo - 1) * req.PageSize
	if err := sql.Offset(offset).Limit(req.PageSize).Find(&issues).
		// reset offset & limit before count
		Offset(0).Limit(-1).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return issues, total, nil
}

// GetIssue 查询事件
func (client *DBClient) GetIssue(id int64) (Issue, error) {
	var issue Issue
	err := client.Where("deleted = 0").Where("id = ?", id).First(&issue).Error
	return issue, err
}

// GetIssue 查询事件
func (client *DBClient) GetIssueWithOutDelete(id int64) (Issue, error) {
	var issue Issue
	err := client.First(&issue, id).Error
	return issue, err
}

func (client *DBClient) GetIssuesIDByIterationID(iterationID uint64) (ids []uint64, err error) {
	var issues []Issue
	err = client.Model(Issue{}).Where("deleted = 0").Where("iteration_id = ?", iterationID).Select("id").Find(&issues).Error
	for _, v := range issues {
		ids = append(ids, v.ID)
	}
	return ids, err
}

func (client *DBClient) ListIssueByIDs(issueIDs []int64) ([]Issue, error) {
	var issues []Issue
	if err := client.Where("`id` IN (?)", issueIDs).Find(&issues).Error; err != nil {
		return nil, err
	}
	return issues, nil
}

// ListIssue 查询事件列表
func (client *DBClient) ListIssue(req pb.IssueListRequest) ([]Issue, error) {
	var issues []Issue
	cond := Issue{}
	if req.ProjectID > 0 {
		cond.ProjectID = req.ProjectID
	}
	if req.IterationID > 0 {
		cond.IterationID = req.IterationID
	}
	// if req.AppID != nil && *req.AppID > 0 {
	// 	cond.AppID = req.AppID
	// }
	// if req.RequirementID != nil && *req.RequirementID > 0 {
	// 	cond.RequirementID = req.RequirementID
	// }

	sql := client.Where(cond)
	// var ts []string
	// for _, i := range req.Type {
	// 	ts = append(ts, i.String())
	// }
	if len(req.Type) > 0 {
		sql = sql.Where("type IN (?)", req.Type)
	}
	if len(req.Assignee) > 0 {
		sql = sql.Where("assignee IN (?)", req.Assignee)
	}
	if len(req.State) > 0 {
		sql = sql.Where("state IN (?)", req.State)
	}
	sql = sql.Where("deleted = ?", 0).Order("id DESC")

	if len(req.IDs) > 0 {
		sql = sql.Where("id IN (?)", req.IDs)
	}
	if req.OnlyIdResult {
		sql = sql.Select("id")
	}

	sql = sql.Find(&issues)
	if err := sql.Error; err != nil {
		return nil, err
	}

	return issues, nil
}

// GetIssueSummary 获取迭代相关的issue统计信息
func (client *DBClient) GetIssueSummary(iterationID int64, task, bug, requirement []int64) apistructs.ISummary {
	var reqDoneCountIDsLine, reqUnDoneCountIDsLine, taskDoneCountIDsLine, taskUnDoneCountIDsLine, bugDoneCountIDsLine, bugUnDoneCountIDsLine []Issue
	// 已完成需求数
	client.Model(Issue{}).Select("id").Where("type = ?", apistructs.IssueTypeRequirement).Where("deleted = ?", 0).
		Where("state in (?)", requirement).Where("iteration_id = ?", iterationID).Find(&reqDoneCountIDsLine)
	// 总需求数
	client.Model(Issue{}).Select("id").Where("type = ?", apistructs.IssueTypeRequirement).Where("deleted = ?", 0).
		Where("iteration_id = ?", iterationID).Find(&reqUnDoneCountIDsLine)
	// 已完成任务数
	client.Model(Issue{}).Select("id").Where("type = ?", apistructs.IssueTypeTask).Where("deleted = ?", 0).
		Where("state in (?)", task).Where("iteration_id = ?", iterationID).Find(&taskDoneCountIDsLine)
	// 总任务数
	client.Model(Issue{}).Select("id").Where("type = ?", apistructs.IssueTypeTask).Where("deleted = ?", 0).
		Where("iteration_id = ?", iterationID).Find(&taskUnDoneCountIDsLine)
	// 已完成bug数
	client.Model(Issue{}).Select("id").Where("type = ?", apistructs.IssueTypeBug).Where("deleted = ?", 0).
		Where("state in (?)", bug).Where("iteration_id = ?", iterationID).Find(&bugDoneCountIDsLine)
	// 总bug数
	client.Model(Issue{}).Select("id").Where("type = ?", apistructs.IssueTypeBug).Where("deleted = ?", 0).
		Where("iteration_id = ?", iterationID).Find(&bugUnDoneCountIDsLine)

	reqDoneCountIDs := issuesToUint64(reqDoneCountIDsLine)
	reqUnDoneCountIDs := issuesToUint64(reqUnDoneCountIDsLine)
	taskDoneCountIDs := issuesToUint64(taskDoneCountIDsLine)
	taskUnDoneCountIDs := issuesToUint64(taskUnDoneCountIDsLine)
	bugDoneCountIDs := issuesToUint64(bugDoneCountIDsLine)
	bugUnDoneCountIDs := issuesToUint64(bugUnDoneCountIDsLine)
	reqUnDoneCountIDs = arrays.DifferenceSet(reqUnDoneCountIDs, reqDoneCountIDs)
	taskUnDoneCountIDs = arrays.DifferenceSet(taskUnDoneCountIDs, taskDoneCountIDs)
	bugUnDoneCountIDs = arrays.DifferenceSet(bugUnDoneCountIDs, bugDoneCountIDs)

	return apistructs.ISummary{
		Requirement: apistructs.ISummaryState{
			Done:   len(reqDoneCountIDsLine),
			UnDone: len(reqUnDoneCountIDsLine) - len(reqDoneCountIDsLine),
		},
		Task: apistructs.ISummaryState{
			Done:   len(taskDoneCountIDsLine),
			UnDone: len(taskUnDoneCountIDsLine) - len(taskDoneCountIDsLine),
		},
		Bug: apistructs.ISummaryState{
			Done:   len(bugDoneCountIDsLine),
			UnDone: len(bugUnDoneCountIDsLine) - len(bugDoneCountIDsLine),
		},
		ReqDoneCountIDs:    reqDoneCountIDs,
		ReqUnDoneCountIDs:  reqUnDoneCountIDs,
		TaskDoneCountIDs:   taskDoneCountIDs,
		TaskUnDoneCountIDs: taskUnDoneCountIDs,
		BugDoneCountIDs:    bugDoneCountIDs,
		BugUnDoneCountIDs:  bugUnDoneCountIDs,
	}
}

func lineToUint64(lines []Line) (issueIDs []uint64) {
	for _, v := range lines {
		issueIDs = append(issueIDs, v.ID)
	}
	return
}

func issuesToUint64(issues []Issue) (issueIDs []uint64) {
	for _, v := range issues {
		issueIDs = append(issueIDs, v.ID)
	}
	return
}

func (client *DBClient) ListIssueSummaryStates(projectID uint64, iterationIDS []int64) ([]IssueSummary, error) {
	var summaries []IssueSummary
	if err := client.Model(Issue{}).Where("project_id = ?", projectID).
		Where("iteration_id in (?)", iterationIDS).
		Where("deleted = ?", 0).
		Group("iteration_id,state,type").
		Select("count(*) as total,state,type as issue_type,iteration_id").
		Scan(&summaries).Error; err != nil {
		return nil, err
	}
	return summaries, nil
}

// GetIssueByIssueIDs 通过issueid获取issue
func (client *DBClient) GetIssueByIssueIDs(issueIDs []uint64) ([]Issue, error) {
	var issues []Issue
	if err := client.Where("deleted = 0").Where("id in ( ? )", issueIDs).Find(&issues).Error; err != nil {
		return nil, err
	}

	return issues, nil
}

// GetIssueByState 通过状态获取issue
func (client *DBClient) GetIssueByState(state int64) (*Issue, error) {
	var issues Issue
	if err := client.Where("deleted = 0").Where("state = ?", state).Find(&issues).Error; err != nil {
		return nil, err
	}
	return &issues, nil
}

// IssueStateCount2 需求关联的任务状态统计
func (client *DBClient) IssueStateCount2(issueIDs []uint64) ([]apistructs.RequirementGroupResult, error) {
	var result []apistructs.RequirementGroupResult
	if err := client.Raw("select id, state from dice_issues where id in (?) and deleted = ?",
		issueIDs, 0).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// GetIssuesManHour get issues mam hour
func (client *DBClient) GetIssuesManHour(req apistructs.IssuesStageRequest) (apistructs.IssueManHourResponse, error) {
	var (
		issues             []Issue
		ans                      = make(map[string]int64)
		sumElapsedTime     int64 = 0
		sumEstimateTime    int64 = 0
		estimateDayGtOne   int64 = 0
		esitmateDayGtTwo   int64 = 0
		esitmateDayGtThree int64 = 0
	)
	sql := client.Table("dice_issues")
	if len(req.StatisticRange) > 0 {
		if req.StatisticRange == "project" {
			sql = sql.Where("project_id in (?)", req.RangeID)
		}
		if req.StatisticRange == "iteration" {
			sql = sql.Where("iteration_id in (?)", req.RangeID)
		}
	}
	if req.Assignee != 0 {
		sql = sql.Where("assignee = ?", req.Assignee)
	}
	if req.Owner != 0 {
		sql = sql.Where("owner = ?", req.Owner)
	}
	if len(req.IssueType) > 0 {
		sql = sql.Where("type = ?", req.IssueType)
	}
	if len(req.StateIDs) > 0 {
		sql = sql.Where("state in (?)", req.StateIDs)
	}
	if err := sql.Where("deleted = ?", 0).Where("type = ?", req.IssueType).Find(&issues).Error; err != nil {
		return apistructs.IssueManHourResponse{}, err
	}
	// 工时，单位与数据库一致 （人分）

	var total, estimateMinute, elapsedMinute uint64
	var avgEstimateMinute, avgElapsedMinute float64
	for _, each := range issues {
		ret := apistructs.IssueManHour{}
		if each.ManHour == "" {
			continue
		}
		total++
		err := json.Unmarshal([]byte(each.ManHour), &ret)
		if err != nil {
			return apistructs.IssueManHourResponse{}, err
		}
		ans[each.Stage] += ret.ElapsedTime
		sumElapsedTime += ret.ElapsedTime
		sumEstimateTime += ret.EstimateTime
		estimateDay := float64(ret.EstimateTime) / day
		if estimateDay > 1 {
			estimateDayGtOne++
		}
		if estimateDay > 2 {
			esitmateDayGtTwo++
		}
		if estimateDay > 3 {
			esitmateDayGtThree++
		}
		estimateMinute += uint64(ret.EstimateTime)
		elapsedMinute += uint64(ret.ElapsedTime)
	}
	if total > 0 {
		avgEstimateMinute = float64(estimateMinute) / float64(total)
		avgElapsedMinute = float64(elapsedMinute) / float64(total)
	}
	return apistructs.IssueManHourResponse{
		DesignManHour:               ans["design"],
		DevManHour:                  ans["dev"],
		TestManHour:                 ans["test"],
		ImplementManHour:            ans["implement"],
		DeployManHour:               ans["deploy"],
		OperatorManHour:             ans["operator"],
		SumElapsedTime:              sumElapsedTime,
		SumEstimateTime:             sumEstimateTime,
		EstimateManDayGtOneDayNum:   estimateDayGtOne,
		EstimateManDayGtTwoDayNum:   esitmateDayGtTwo,
		EstimateManDayGtThreeDayNum: esitmateDayGtThree,

		Total:               total,
		AvgEstimateMinute:   avgEstimateMinute,
		AvgElapsedMinute:    avgElapsedMinute,
		TotalEstimateMinute: estimateMinute,
		TotalElapsedMinute:  elapsedMinute,
	}, nil
}

// GetIssueByRange 通过迭代或项目获取issue Bug
func (client *DBClient) GetIssueBugByRange(req apistructs.IssuesStageRequest) ([]Issue, float32, error) {
	var (
		issue []Issue
		total float32
	)
	sql := client.Table("dice_issues").Where("deleted = ?", 0).Where("type = ?", apistructs.IssueTypeBug)
	if len(req.StatisticRange) > 0 {
		if req.StatisticRange == "project" {
			sql = sql.Where("project_id in (?)", req.RangeID)
		}
		if req.StatisticRange == "iteration" {
			sql = sql.Where("iteration_id in (?)", req.RangeID)
		}
	}
	if err := sql.Scan(&issue).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return issue, total, nil
}

// GetReceiversByIssueID get receivers of issue event
func (client *DBClient) GetReceiversByIssueID(issueID int64) ([]string, error) {
	var receivers []string
	issue, err := client.GetIssue(issueID)
	if err != nil {
		return nil, err
	}
	subscribers, err := client.GetIssueSubscribersSliceByIssueID(issueID)
	if err != nil {
		return nil, err
	}
	receivers = append(receivers, issue.Assignee)
	receivers = append(receivers, subscribers...)

	return strutil.DedupSlice(receivers), nil
}

// GetIssueNumByPros query by IssuePagingRequest and group by project_id to get special issue num
func (client *DBClient) GetIssueNumByPros(projectIDS []uint64, req apistructs.IssuePagingRequest) ([]apistructs.IssueNum, error) {
	var (
		res []apistructs.IssueNum
	)

	sql := client.Table("dice_issues").Select("count(id) as issue_num, project_id").Where("deleted = ?", 0)
	if len(req.IDs) > 0 {
		sql = sql.Where("dice_issues.id IN (?)", req.IDs)
	}
	if len(req.IterationIDs) > 0 {
		sql = sql.Where("iteration_id in (?)", req.IterationIDs)
	}
	if len(req.Type) > 0 {
		sql = sql.Where("type IN (?)", req.Type)
	}
	if len(req.Creators) > 0 {
		sql = sql.Where("creator IN (?)", req.Creators)
	}
	if len(req.Assignees) > 0 {
		sql = sql.Where("assignee IN (?)", req.Assignees)
	}
	if len(req.Priority) > 0 {
		sql = sql.Where("priority IN (?)", req.Priority)
	}
	if len(req.Complexity) > 0 {
		sql = sql.Where("complexity IN (?)", req.Complexity)
	}
	if len(req.Severity) > 0 {
		sql = sql.Where("severity IN (?)", req.Severity)
	}
	if len(req.State) > 0 {
		sql = sql.Where("state IN (?)", req.State)
	}
	if len(req.Owner) > 0 {
		sql = sql.Where("owner IN (?)", req.Owner)
	}
	if len(req.BugStage) > 0 {
		sql = sql.Where("stage IN (?)", req.BugStage)
	} else if len(req.TaskType) > 0 {
		sql = sql.Where("stage IN (?)", req.TaskType)
	}
	if len(req.ExceptIDs) > 0 {
		sql = sql.Not("id", req.ExceptIDs)
	}
	if req.StartCreatedAt > 0 {
		startCreatedAt := time.Unix(req.StartCreatedAt/1000, 0)
		sql = sql.Where("dice_issues.created_at >= ?", startCreatedAt)
	}
	if req.EndCreatedAt > 0 {
		endCreatedAt := time.Unix(req.EndCreatedAt/1000, 0)
		sql = sql.Where("dice_issues.created_at <= ?", endCreatedAt)
	}
	if req.IsEmptyPlanFinishedAt {
		sql = sql.Where("plan_finished_at IS NULL")
	}
	if req.StartFinishedAt > 0 && !req.IsEmptyPlanFinishedAt {
		startFinishedAt := time.Unix(req.StartFinishedAt/1000, 0)
		sql = sql.Where("plan_finished_at >= ?", startFinishedAt)
	}
	if req.EndFinishedAt > 0 && !req.IsEmptyPlanFinishedAt {
		endFinishedAt := time.Unix(req.EndFinishedAt/1000, 0)
		sql = sql.Where("plan_finished_at <= ?", endFinishedAt)
	}

	if req.StartClosedAt > 0 {
		startClosedAt := time.Unix(req.StartClosedAt/1000, 0)
		sql = sql.Where("finish_time >= ?", startClosedAt)
	}
	if req.EndClosedAt > 0 {
		endClosedAt := time.Unix(req.EndClosedAt/1000, 0)
		sql = sql.Where("finish_time <= ?", endClosedAt)
	}
	if len(projectIDS) > 0 {
		sql = sql.Where("project_id IN (?)", projectIDS)
	}

	if req.Title != "" {
		title := strings.ReplaceAll(req.Title, "%", "\\%")
		if _, err := strutil.Atoi64(title); err == nil {
			sql = sql.Where("title LIKE ? OR dice_issues.id LIKE ?", "%"+title+"%", "%"+title+"%")
		} else {
			sql = sql.Where("title LIKE ?", "%"+title+"%")
		}
	}
	if req.Source != "" {
		sql = sql.Where("source LIKE ?", "%"+req.Source+"%")
	}
	if req.External {
		sql = sql.Where("external = 1")
	} else {
		sql = sql.Where("external = 0")
	}

	if err := sql.Group("project_id").Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

var joinState = "LEFT JOIN dice_issue_state ON dice_issues.state = dice_issue_state.id"
var joinParticipant = "LEFT JOIN erda_issue_subscriber on dice_issues.id = erda_issue_subscriber.issue_id"

type IssueExpiryStatus struct {
	IssueNum     uint64
	ProjectID    uint64
	ExpiryStatus ExpireType
}

func (client *DBClient) GetIssueExpiryStatusByProjects(req apistructs.WorkbenchRequest) ([]IssueExpiryStatus, error) {
	sql := client.issueExpiryStatusQuery(req)
	var res []IssueExpiryStatus
	// query with matched projects
	if err := sql.Select("count(dice_issues.id) as issue_num, dice_issues.project_id, dice_issues.expiry_status").
		Group("dice_issues.project_id, dice_issues.expiry_status").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (client *DBClient) issueExpiryStatusQuery(req apistructs.WorkbenchRequest) *gorm.DB {
	sql := client.Table("dice_issues").Joins(joinState)
	sql = sql.Where("deleted = 0").Where("assignee IN (?) AND dice_issue_state.belong IN (?)", req.Assignees, req.StateBelongs)
	if len(req.Type) > 0 {
		sql = sql.Where("type IN (?)", req.Type)
	}
	if len(req.ProjectIDs) > 0 {
		sql = sql.Where("dice_issues.project_id IN (?)", req.ProjectIDs)
	}
	return sql
}

func (client *DBClient) GetIssuesByProject(req apistructs.IssuePagingRequest) ([]Issue, uint64, error) {
	var res []Issue
	sql := client.Table("dice_issues").Joins(joinState)
	sql = sql.Where("deleted = 0").Where("dice_issues.project_id = ?", req.ProjectID)
	if len(req.Type) > 0 {
		sql = sql.Where("type IN (?)", req.Type)
	}
	if len(req.Assignees) > 0 {
		sql = sql.Where("assignee IN (?)", req.Assignees)
	}
	if len(req.StateBelongs) > 0 {
		sql = sql.Where("dice_issue_state.belong IN (?)", req.StateBelongs)
	}
	if req.OrderBy != "" {
		if req.Asc {
			sql = sql.Order(fmt.Sprintf("%s", req.OrderBy))
		} else {
			sql = sql.Order(fmt.Sprintf("%s DESC", req.OrderBy))
		}
	}
	var total uint64
	offset := (req.PageNo - 1) * req.PageSize
	if err := sql.Offset(offset).Limit(req.PageSize).Find(&res).Offset(0).Limit(-1).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return res, total, nil
}

var expireTypes = []ExpireType{
	ExpireTypeUndefined, ExpireTypeExpired, ExpireTypeExpireIn1Day,
	ExpireTypeExpireIn2Days, ExpireTypeExpireIn7Days, ExpireTypeExpireIn30Days, ExpireTypeExpireInFuture,
}

var conditions = map[ExpireType]string{
	ExpireTypeUndefined:      "DATE(a.plan_finished_at) IS NULL",
	ExpireTypeExpired:        "DATE(a.plan_finished_at) < CURDATE()",
	ExpireTypeExpireIn1Day:   "DATE(a.plan_finished_at) = CURDATE()",
	ExpireTypeExpireIn2Days:  "DATE(a.plan_finished_at) = DATE_ADD(CURDATE(),INTERVAL 1 DAY)",
	ExpireTypeExpireIn7Days:  "DATE(a.plan_finished_at) > DATE_ADD(CURDATE(),INTERVAL 1 DAY) AND DATE(a.plan_finished_at) < DATE_ADD(CURDATE(),INTERVAL 7 DAY)",
	ExpireTypeExpireIn30Days: "DATE(a.plan_finished_at) >= DATE_ADD(CURDATE(),INTERVAL 7 DAY) AND DATE(a.plan_finished_at) < DATE_ADD(CURDATE(),INTERVAL 30 DAY)",
	ExpireTypeExpireInFuture: "DATE(a.plan_finished_at) >= DATE_ADD(CURDATE(),INTERVAL 30 DAY)",
}

func (client *DBClient) BatchUpdateIssueExpiryStatus(states []apistructs.IssueStateBelong) error {
	for _, key := range expireTypes {
		if _, ok := conditions[key]; !ok {
			continue
		}
		sql := fmt.Sprintf("UPDATE dice_issues a LEFT JOIN dice_issue_state b ON a.state = b.id SET a.expiry_status = ? WHERE a.deleted = 0 AND a.expiry_status != ? AND b.belong IN (?) AND %s", conditions[key])
		if err := client.Exec(sql, key, key, states).Error; err != nil {
			return err
		}
	}
	return nil
}

type IssueItem struct {
	dbengine.BaseModel

	PlanStartedAt  *time.Time                 // 计划开始时间
	PlanFinishedAt *time.Time                 // 计划结束时间
	ProjectID      uint64                     // 所属项目 ID
	IterationID    int64                      // 所属迭代 ID
	AppID          *uint64                    // 所属应用 ID
	RequirementID  *int64                     // 所属需求 ID
	Type           string                     // issue 类型
	Title          string                     // 标题
	Content        string                     // 内容
	State          int64                      // 状态
	Priority       apistructs.IssuePriority   // 优先级
	Complexity     apistructs.IssueComplexity // 复杂度
	Severity       apistructs.IssueSeverity   // 严重程度
	Creator        string                     // issue 创建者 ID
	Assignee       string                     // 分配到 issue 的人，即当前处理人
	Source         string                     // issue创建的来源，目前只有工单使用
	ManHour        string                     // 工时信息
	External       bool                       // 用来区分是通过ui还是bundle创建的
	Deleted        bool                       // 是否已删除
	Stage          string                     // bug阶段 or 任务类型 的值
	Owner          string                     // 负责人

	FinishTime   *time.Time // 实际结束时间
	ExpiryStatus ExpireType
	ReopenCount  int
	StartTime    *time.Time

	Name           string
	Belong         string
	ChildrenLength int
}

func (i *IssueItem) FilterPropertyRetriever(condition string) string {
	r := reflect.ValueOf(i)
	f := reflect.Indirect(r).FieldByName(condition)
	return string(f.String())
}

func (client *DBClient) ListIssueItems(req pb.IssueListRequest) ([]IssueItem, error) {
	var res []IssueItem
	sql := client.Table("dice_issues").Joins(joinState)
	sql = sql.Where("deleted = 0")
	if req.ProjectID != 0 {
		sql = sql.Where("dice_issues.project_id = ?", req.ProjectID)
	}
	if len(req.StateBelongs) > 0 {
		sql = sql.Where("dice_issue_state.belong IN (?)", req.StateBelongs)
	}
	if len(req.Assignee) > 0 {
		sql = sql.Where("assignee in (?)", req.Assignee)
	}
	if len(req.IterationIDs) > 0 {
		sql = sql.Where("iteration_id in (?)", req.IterationIDs)
	}
	if len(req.Type) > 0 {
		sql = sql.Where("type IN (?)", req.Type)
	}
	if len(req.IDs) > 0 {
		sql = sql.Where("dice_issues.id in (?)", req.IDs)
	}
	if len(req.Label) > 0 {
		sql = sql.Joins("LEFT JOIN dice_label_relations c ON dice_issues.id = c.ref_id").Where("c.label_id IN (?)", req.Label)
	}
	if err := sql.Select("dice_issues.*, dice_issue_state.name, dice_issue_state.belong").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

type IssueLabel struct {
	dbengine.BaseModel
	LabelID uint64                      // 标签 id
	RefType apistructs.ProjectLabelType // 标签作用类型, eg: issue
	RefID   string                      // 标签关联目标 id
	Name    string                      // 标签名称
	Type    apistructs.ProjectLabelType // 标签作用类型
}

var joinLabel = "LEFT JOIN dice_labels ON dice_label_relations.label_id = dice_labels.id"

func (client *DBClient) GetIssueLabelsByProjectID(projectID uint64) ([]IssueLabel, error) {
	var res []IssueLabel
	sql := client.Table("dice_label_relations").Joins(joinLabel).Where("project_id = ?", projectID)
	if err := sql.Select("dice_label_relations.*, dice_labels.name, dice_labels.type").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (client *DBClient) CountBugBySeverity(projectID uint64, iterationIDs []uint64, onlyUnclosed bool) (map[apistructs.IssueSeverity]uint64, error) {
	type Line struct {
		Total    uint64
		Severity apistructs.IssueSeverity
	}
	var results []Line

	sql := client.Model(&Issue{}).Select("COUNT(*) as `total`, `severity`").Where("`deleted` = 0").Where("`type` = ?", apistructs.IssueTypeBug)
	if projectID > 0 {
		sql = sql.Where("project_id = ?", projectID)
	}
	if len(iterationIDs) > 0 {
		sql = sql.Where("iteration_id IN (?)", iterationIDs)
	}
	if onlyUnclosed {
		// get state ids by state belong
		bugStates, err := client.GetIssuesStates(&pb.GetIssueStatesRequest{
			ProjectID:    projectID,
			IssueType:    pb.IssueTypeEnum_BUG.String(),
			StateBelongs: common.UnclosedStateBelongs,
		})
		if err != nil {
			return nil, err
		}
		var bugStateIDs []uint64
		for _, state := range bugStates {
			bugStateIDs = append(bugStateIDs, state.ID)
		}
		sql = sql.Where("state IN (?)", bugStateIDs)
	}
	sql = sql.Group("`severity`").Scan(&results)

	if err := sql.Error; err != nil {
		return nil, fmt.Errorf("failed to count bug by severity, err: %v", err)
	}

	m := make(map[apistructs.IssueSeverity]uint64)
	for _, line := range results {
		m[line.Severity] = line.Total
	}

	return m, nil
}

func (client *DBClient) BugReopenCount(projectID uint64, iterationIDs []uint64) (reopenCount, totalCount uint64, ids []uint64, err error) {
	var issues []Issue
	sql := client.Model(&Issue{}).Where("`type` = ?", apistructs.IssueTypeBug).Where("deleted = 0")
	if projectID > 0 {
		sql = sql.Where("project_id = ?", projectID)
	}
	if len(iterationIDs) > 0 {
		sql = sql.Where("iteration_id IN (?)", iterationIDs)
	}

	type Line struct {
		Sum   uint64
		Total uint64
		IDs   []uint64
	}
	var result Line
	if err := sql.Select("SUM(`reopen_count`) AS sum, COUNT(*) AS total").Scan(&result).Error; err != nil {
		return 0, 0, []uint64{}, fmt.Errorf("failed to sum bug reopen count, err: %v", err)
	}
	if err := sql.Select("id").Find(&issues).Error; err != nil {
		return 0, 0, []uint64{}, fmt.Errorf("failed to sum bug reopen count, err: %v", err)
	}
	ids = issuesToUint64(issues)
	return result.Sum, result.Total, ids, nil
}

const (
	joinRelation       = "LEFT JOIN dice_issue_relation b ON a.id = b.related_issue and b.type = ?"
	joinRelationParent = "left join dice_issue_relation d on a.id = d.issue_id and d.type = ?"
	joinStateNew       = "LEFT JOIN dice_issue_state ON a.state = dice_issue_state.id"
	joinIssueChildren  = "LEFT JOIN dice_issues a ON a.id = b.related_issue"
	joinIssueParent    = "LEFT JOIN dice_issues a ON a.id = b.issue_id"
	joinLabelRelation  = "LEFT JOIN dice_label_relations c ON a.id = c.ref_id"
	ganttOrder         = "FIELD(a.type,'REQUIREMENT','TASK','BUG')"
)

func (client *DBClient) FindIssueChildren(id uint64, req pb.PagingIssueRequest) ([]IssueItem, uint64, error) {
	sql := client.Table("dice_issue_relation b").Joins(joinIssueChildren).Joins(joinStateNew).
		Where("b.issue_id = ? AND b.type = ?", id, apistructs.IssueRelationInclusion)
	// if id == 0 {
	// 	sql = client.Table("dice_issues as a").Joins(joinRelation, apistructs.IssueRelationInclusion).Joins(joinStateNew).
	// 		Where("b.id IS NULL")
	// }
	if len(req.Type) > 0 {
		sql = sql.Where("a.type IN (?)", req.Type)
	}
	if len(req.Assignee) > 0 {
		sql = sql.Where("a.assignee in (?)", req.Assignee)
	}
	if len(req.StateBelongs) > 0 {
		sql = sql.Where("dice_issue_state.belong IN (?)", req.StateBelongs)
	}
	if len(req.State) > 0 {
		sql = sql.Where("a.state in (?)", req.State)
	}
	sql = applyCondition(sql, req)
	offset := (req.PageNo - 1) * req.PageSize
	var total uint64
	var res []IssueItem
	if err := sql.Select("DISTINCT a.*, dice_issue_state.name, dice_issue_state.belong").Order("a.type").Offset(offset).Limit(req.PageSize).Find(&res).
		Offset(0).Limit(-1).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return res, total, nil
}

func (client *DBClient) FindIssueRoot(req pb.PagingIssueRequest) ([]IssueItem, []IssueItem, uint64, error) {
	// issues without children
	sql := client.Table("dice_issues as a").Joins(joinRelation, apistructs.IssueRelationInclusion).Joins(joinStateNew).
		Joins(joinRelationParent, apistructs.IssueRelationInclusion).Where("b.id IS NULL and d.id is NULL")
	offset := (req.PageNo - 1) * req.PageSize
	if len(req.Type) > 0 {
		sql = sql.Where("a.type IN (?)", req.Type)
	}
	sql = applyCondition(sql, req)
	if len(req.StateBelongs) > 0 {
		sql = sql.Where("dice_issue_state.belong IN (?)", req.StateBelongs)
	}
	if len(req.Assignee) > 0 {
		sql = sql.Where("a.assignee in (?)", req.Assignee)
	}
	if len(req.State) > 0 {
		sql = sql.Where("a.state in (?)", req.State)
	}
	var items []IssueItem
	var totalTask uint64
	if err := sql.Select("DISTINCT a.*, dice_issue_state.name, dice_issue_state.belong").Order(ganttOrder).Offset(offset).Limit(req.PageSize).Find(&items).
		Offset(0).Limit(-1).Count(&totalTask).Error; err != nil {
		return nil, nil, 0, err
	}

	// requirements with children
	sql = client.Table("dice_issue_relation b").Joins(joinIssueParent).Joins("LEFT JOIN dice_issues d ON d.id = b.related_issue").
		Joins("LEFT JOIN dice_issue_state e ON a.state = e.id").Joins("LEFT JOIN dice_issue_state f ON d.state = f.id")
	sql = sql.Where("a.type = ?", apistructs.IssueTypeRequirement).Where("b.type = ?", apistructs.IssueRelationInclusion)
	sql = sql.Where("d.deleted = 0").Where("d.project_id = ?", req.ProjectID)
	sql = applyCondition(sql, req)
	if len(req.Assignee) > 0 {
		sql = sql.Where("a.assignee in (?) or d.assignee in (?)", req.Assignee, req.Assignee)
	}
	if len(req.StateBelongs) > 0 {
		sql = sql.Where("e.belong IN (?)", req.StateBelongs)
	}
	if len(req.State) > 0 {
		sql = sql.Where("e.id IN (?)", req.State)
	}
	var res []IssueItem
	var totalReq uint64
	if err := sql.Select("DISTINCT a.*, e.name, e.belong").Offset(offset).Limit(req.PageSize).Find(&res).
		Offset(0).Limit(-1).Count(&totalReq).Error; err != nil {
		return nil, nil, 0, err
	}
	return res, items, totalReq + totalTask, nil
}

func applyCondition(sql *gorm.DB, req pb.PagingIssueRequest) *gorm.DB {
	sql = sql.Where("a.deleted = 0").Where("a.project_id = ?", req.ProjectID)
	if len(req.IterationIDs) > 0 {
		sql = sql.Where("a.iteration_id in (?)", req.IterationIDs)
	}
	if len(req.Label) > 0 {
		sql = sql.Joins(joinLabelRelation).Where("c.label_id IN (?)", req.Label)
	}
	return sql
}

func (client *DBClient) GetIssueItem(id uint64) (IssueItem, error) {
	var issue IssueItem
	err := client.Table("dice_issues").Joins(joinState).Where("deleted = 0").Where("dice_issues.id = ?", id).Select("dice_issues.*, dice_issue_state.name, dice_issue_state.belong").First(&issue).Error
	return issue, err
}

func (client *DBClient) GetIssueParents(issueID uint64, relationType []string) ([]IssueItem, error) {
	var issues []IssueItem
	sql := client.Table("dice_issue_relation b").Joins(joinIssueParent).Joins(joinStateNew).Where("related_issue = ?", issueID)
	if len(relationType) > 0 {
		sql = sql.Where("b.type IN (?)", relationType)
	}
	if err := sql.Select("a.*, dice_issue_state.name, dice_issue_state.belong").Find(&issues).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return issues, nil
}

type timeRange struct {
	Min *time.Time
	Max *time.Time
}

func (client *DBClient) FindIssueChildrenTimeRange(id uint64) (*time.Time, *time.Time, error) {
	sql := client.Table("dice_issue_relation b").Joins(joinIssueChildren).
		Where("b.issue_id = ? AND b.type = ?", id, apistructs.IssueRelationInclusion)
	var res timeRange
	if err := sql.Select("MAX(a.`plan_finished_at`) as max, MIN(a.`plan_started_at`) as min").Find(&res).Error; err != nil {
		return nil, nil, err
	}
	return res.Min, res.Max, nil
}

// IssueForIssueStateTransMigration the type of state is string
type IssueForIssueStateTransMigration struct {
	dbengine.BaseModel

	ProjectID uint64
	Type      apistructs.IssueType
	Creator   string
}

func (client *DBClient) ListIssueForIssueStateTransMigration() ([]IssueForIssueStateTransMigration, error) {
	var issues []IssueForIssueStateTransMigration
	err := client.Table("dice_issues").Where("deleted = 0").Find(&issues).Error
	return issues, err
}

// GetHaveUndoneTaskAssigneeNum get the number of people with unfinished tasks under the iteration or project
func (client *DBClient) GetHaveUndoneTaskAssigneeNum(iterationID uint64, projectID uint64, doneStateIDS []int64) ([]string, error) {
	type Assignee struct {
		Assignee string
	}
	assignees := make([]Assignee, 0)
	sql := client.Table("dice_issues").Select("distinct(assignee) as assignee").Where("deleted = 0").Where("type = 'TASK' and assignee != ''").
		Where("state not in (?)", doneStateIDS)
	if iterationID > 0 {
		sql = sql.Where("iteration_id = ?", iterationID)
	}
	if projectID > 0 {
		sql = sql.Where("project_id = ?", projectID)
	}
	if err := sql.Find(&assignees).Error; err != nil {
		return nil, err
	}
	assigneeList := make([]string, 0)
	for _, assignee := range assignees {
		assigneeList = append(assigneeList, assignee.Assignee)
	}
	return assigneeList, nil
}

func (client *DBClient) GetSeriousBugNum(iterationID uint64) (issueIDsCount uint64, issueIDsUint64 []uint64, err error) {
	var issueIDs []Issue
	if err = client.Table("dice_issues").Select("distinct(id)").Where("iteration_id = ?", iterationID).Where("deleted = 0").Where("type = 'BUG'").
		Where("severity in (?)", []string{string(apistructs.IssueSeverityFatal), string(apistructs.IssueSeveritySerious)}).
		Find(&issueIDs).Error; err != nil {
		return 0, issueIDsUint64, err
	}
	issueIDsUint64 = issuesToUint64(issueIDs)
	issueIDsCount = uint64(len(issueIDs))
	return issueIDsCount, issueIDsUint64, err
}

func (client *DBClient) GetDemandDesignBugNum(iterationID uint64) (issueIDsCount uint64, issueIDsUint64 []uint64, err error) {
	var issueIDs []Issue
	if err = client.Table("dice_issues").Select("distinct(id)").
		Where("iteration_id = ?", iterationID).
		Where("deleted = 0").Where("type = ?", apistructs.IssueTypeBug).
		Where("stage = ?", "demandDesign").
		Find(&issueIDs).Error; err != nil {
		return 0, []uint64{}, err
	}
	issueIDsUint64 = issuesToUint64(issueIDs)
	issueIDsCount = uint64(len(issueIDs))
	return issueIDsCount, issueIDsUint64, err
}

func (client *DBClient) GetIssueNumByStates(iterationID uint64, issueType apistructs.IssueType, states []uint64) (issueIDsCount uint64, issueIDsUint64 []uint64, err error) {
	var issueIDs []Issue
	if err := client.Table("dice_issues").Select("distinct(id)").
		Where("iteration_id = ?", iterationID).
		Where("deleted = 0").Where("type = ?", issueType).
		Where("state in (?)", states).
		Find(&issueIDs).Error; err != nil {
		return 0, []uint64{}, err
	}
	issueIDsUint64 = issuesToUint64(issueIDs)
	issueIDsCount = uint64(len(issueIDs))
	return issueIDsCount, issueIDsUint64, err
}

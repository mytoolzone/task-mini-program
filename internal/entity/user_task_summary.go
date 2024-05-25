package entity

type UserTaskSummary struct {
	TotalTask           int64                 `json:"total_task_num"`                // 总任务数
	TotalDuration       int64                 `json:"total_task_duration,omitempty"` // 总任务时长
	UserTaskSummaryList []UserTaskSummaryList `json:"user_task_summary_list"`        //任务列表信息
}

type UserTaskSummaryList struct {
	TotalDuration int64 `json:"total_duration,omitempty"` // 总任务时长,默认0
	TASK_ID       int64 `json:"task_id"`                  //任务ID
	User_ID       int64 `json:"user_id"`                  //用户ID
}

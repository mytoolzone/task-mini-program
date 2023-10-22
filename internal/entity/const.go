package entity

// 定义任务的执行状态
// 任务的执行状态有：发布、审核失败、审核通过待报名, 报名结束等待运行, 运行、暂停、完成、取消
const (
	TaskStatusPublish   = "publish"
	TaskStatusAuditFail = "audit_fail"
	TaskStatusAuditPass = "audit_pass"
	TaskStatusPending   = "pending"
	TaskStatusRunning   = "running"
	TaskStatusPaused    = "paused"
	TaskStatusFinished  = "finished"
	TaskStatusCanceled  = "canceled"
)

// 用户状态 有效 冻结
const (
	UserStatusActive = "active"
	UserStatusFrozen = "frozen"
)

// 用户报名任务审核状态
const (
	UserTaskStatusApply     = "apply"
	UserTaskStatusAuditFail = "rejected"
	UserTaskStatusAuditPass = "approved"
)

// 用户报名在任务中角色 队长 队员 记录员 其他
const (
	UserTaskRoleLeader   = "leader"
	UserTaskRoleMember   = "member"
	UserTaskRoleRecorder = "recorder"
	UserTaskRoleNone     = "none"
)

// 通知阅读状态 未读 已读
const (
	NotifyStatusUnread = "unread"
	NotifyStatusRead   = "read"
)

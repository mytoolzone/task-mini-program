package entity

// 定义任务的执行状态
// 任务的执行状态有：新任务(发布待审核)、审核失败、报名(审核通过待报名), 待运行(报名结束等待运行), 运行、暂停、完成、取消
const (
	TaskStatusNew       = "new"
	TaskStatusAuditFail = "audit_fail"
	TaskStatusJoin      = "join"
	TaskStatusTorun     = "torun"   // 开始报名
	TaskStatusSign      = "sign"    // 开始报名
	TaskStatusRunning   = "running" // 报名结束
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
	UserTaskStatusAuditFail = StatusAuditReject
	UserTaskStatusAuditPass = StatusAuditApproved
)

const (
	StatusAuditReject   = "rejected"
	StatusAuditApproved = "approved"
)

// 用户角色 管理员 一般用户
const (
	UserRoleAdmin = "admin"
	UserRoleUser  = "user"
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

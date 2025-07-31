// 任务状态
export const TASK_STATUS_NEW = 'new'    -> 新创建任务 (审核中)
export const TASK_STATUS_AUDITFAIL = 'audit_fail'  -> 审核失败

export const TASK_STATUS_TORUN = 'torun' -> 任务审核通过 待运行状态

export const TASK_STATUS_JOIN = 'join' -> 报名中


export const TASK_STATUS_SIGN = 'sign' 
export const TASK_STATUS_RUNNING = 'running'
export const TASK_STATUS_PAUSED = 'paused'
export const TASK_STATUS_FINISHED = 'finished'F


[TASK_STATUS_NEW]:"待审核任务",
[TASK_STATUS_AUDITFAIL]:"审核未通过",
[TASK_STATUS_JOIN]:"待报名",
[TASK_STATUS_TORUN]:"待人员报名",
[TASK_STATUS_SIGN]:"开始签到",
[TASK_STATUS_RUNNING]:"进行中",
[TASK_STATUS_PAUSED]:"已暂停",
[TASK_STATUS_FINISHED]:"已结束",

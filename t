[1mdiff --git a/internal/controller/http/middleware/checkrole.go b/internal/controller/http/middleware/checkrole.go[m
[1mindex f0f8473..bedc3ff 100644[m
[1m--- a/internal/controller/http/middleware/checkrole.go[m
[1m+++ b/internal/controller/http/middleware/checkrole.go[m
[36m@@ -40,24 +40,24 @@[m [mfunc CheckRole(userCase usecase.User, taskCase usecase.Task) gin.HandlerFunc {[m
 		}[m
 		http_util.SetUserRole(c, role)[m
 [m
[31m-		// 任务角色查询[m
[32m+[m		[32m// 当前登录人的任务角色查询[m
 		var taskRole entity.UserTask[m
 		TaskId := c.Query("taskID")[m
[31m-		UserId := c.Query("userID")[m
[31m-		if TaskId == "" && UserId == "" {[m
[32m+[m		[32m// UserId := c.Query("userID")[m
[32m+[m		[32mif TaskId == "" {[m
 			var TaskParams struct {[m
 				TaskId int `json:"task_id"`[m
[31m-				UserId int `json:"user_id"`[m
[32m+[m				[32m// UserId int `json:"user_id"`[m
 			}[m
 			if err := c.ShouldBindJSON(&TaskParams); err == nil {[m
[31m-				taskRole, _ = taskCase.GetUserTaskRole(c, TaskParams.TaskId, TaskParams.UserId)[m
[31m-				glog.Infof("用户%d 在任务 %d 中的角色为 %s", TaskParams.UserId, TaskParams.TaskId, taskRole.Role)[m
[32m+[m				[32mtaskRole, _ = taskCase.GetUserTaskRole(c, TaskParams.TaskId, userID)[m
[32m+[m				[32mglog.Infof("用户%d 在任务 %d 中的角色为 %s", userID, TaskParams.TaskId, taskRole.Role)[m
 			}[m
 		} else {[m
 			taskIdInt, _ := strconv.Atoi(TaskId)[m
[31m-			userIdInt, _ := strconv.Atoi(UserId)[m
[31m-			taskRole, _ = taskCase.GetUserTaskRole(c, taskIdInt, userIdInt)[m
[31m-			glog.Infof("用户%d 在任务 %d 中的角色为 %s", userIdInt, taskIdInt, taskRole.Role)[m
[32m+[m			[32m// userIdInt, _ := strconv.Atoi(UserId)[m
[32m+[m			[32mtaskRole, _ = taskCase.GetUserTaskRole(c, taskIdInt, userID)[m
[32m+[m			[32mglog.Infof("用户%d 在任务 %d 中的角色为 %s", userID, taskIdInt, taskRole.Role)[m
 		}[m
 		// 判断是否需要校验权限[m
 		if checkRoles, ok := checkPathToRole[path]; ok {[m
[1mdiff --git a/internal/entity/tasks.gen.go b/internal/entity/tasks.gen.go[m
[1mindex 224a786..8e4f937 100644[m
[1m--- a/internal/entity/tasks.gen.go[m
[1m+++ b/internal/entity/tasks.gen.go[m
[36m@@ -41,12 +41,13 @@[m [mtype Task struct {[m
 	// type 表示任务的类型·。[m
 	// Enum: task,post[m
 	// Description: 任务的类型可以task正常任务，或者是post通告。[m
[31m-	Type     string `gorm:"column:type" json:"type"`[m
[31m-	Leader   int    `gorm:"column:leader" json:"leader"`[m
[31m-	Recorder int    `gorm:"column:recorder" json:"recorder"`[m
[31m-	StartAt 	time.Time	`gorm:"column:finished_at" json:"start_at"`[m
[31m-	MeetingAt 	string 	`gorm:"column:meeting_at;type:timestamp" json:"meeting_at"`[m
[31m-	Contacter 	int    `gorm:"column:contacter" json:"contacter"`[m
[32m+[m	[32mType             string    `gorm:"column:type" json:"type"`[m
[32m+[m	[32mLeader           int       `gorm:"column:leader" json:"leader"`[m
[32m+[m	[32mRecorder         int       `gorm:"column:recorder" json:"recorder"`[m
[32m+[m	[32mStartAt          time.Time `gorm:"column:finished_at" json:"start_at"`[m
[32m+[m	[32mMeetingAt        string    `gorm:"column:meeting_at;type:timestamp" json:"meeting_at"`[m
[32m+[m	[32mContacter        int       `gorm:"column:contacter" json:"contacter"`[m
[32m+[m	[32mJoinPersonsCount int64     `gorm:"-" json:"join_persons_count"`[m
 }[m
 [m
 // TableName Task's table name[m
[1mdiff --git a/internal/usecase/repo/task_postgres.go b/internal/usecase/repo/task_postgres.go[m
[1mindex 148bcc9..e77b782 100644[m
[1m--- a/internal/usecase/repo/task_postgres.go[m
[1m+++ b/internal/usecase/repo/task_postgres.go[m
[36m@@ -138,6 +138,12 @@[m [mfunc (t *TaskRepo) GetTaskList(ctx context.Context, lastId int, keyword, status[m
 	}[m
 [m
 	err := query.Order("id desc").Find(&tasks).Error[m
[32m+[m	[32mfor i, task := range tasks {[m
[32m+[m		[32m// 查询每个任务的参与人数[m
[32m+[m		[32mvar count int64[m
[32m+[m		[32m_ = t.Db.WithContext(ctx).Model(&entity.UserTask{}).Where("task_id = ? and status = ?", task.ID, entity.UserTaskStatusAuditPass).Count(&count).Error[m
[32m+[m		[32mtasks[i].JoinPersonsCount = count[m
[32m+[m	[32m}[m
 	return tasks, err[m
 }[m
 [m

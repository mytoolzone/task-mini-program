-- 创建 ENUM 类型
CREATE TYPE user_status AS ENUM ('active', 'inactive');
CREATE TYPE task_status AS ENUM ('open', 'closed');
CREATE TYPE task_user_status AS ENUM ('apply', 'approved', 'rejected');
CREATE TYPE task_user_role AS ENUM ('leader', 'member', 'recorder', 'audience');
CREATE TYPE notice_status AS ENUM ('unread', 'read');

-- 用户表
CREATE TABLE users (
           id SERIAL PRIMARY KEY,
           username VARCHAR(255) NOT NULL,
           password VARCHAR(255) NOT NULL,
           email VARCHAR(255),
           phone VARCHAR(20),
           status user_status,
           createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
           updatedAt TIMESTAMP WITH TIME ZONE,
           deletedAt TIMESTAMP WITH TIME ZONE,
           openid VARCHAR(255),
           ext JSONB
);

-- 用户角色表
CREATE TABLE user_roles (
                            id SERIAL PRIMARY KEY,
                            userId INTEGER REFERENCES users(id),
                            role VARCHAR(20) NOT NULL,
                            createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                            updatedAt TIMESTAMP WITH TIME ZONE
);

-- 任务表
CREATE TABLE tasks (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) NOT NULL,
                       createBy INTEGER REFERENCES users(id),
                       createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       updatedAt TIMESTAMP WITH TIME ZONE,
                       endAt TIMESTAMP WITH TIME ZONE,
                       deletedAt TIMESTAMP WITH TIME ZONE,
                       describe VARCHAR(1024),
                       require VARCHAR(1024),
                       location POINT,
                       status task_status
);

-- 任务执行人角色表
CREATE TABLE task_users (
                            id SERIAL PRIMARY KEY,
                            taskId INTEGER REFERENCES tasks(id),
                            userId INTEGER REFERENCES users(id),
                            createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                            status task_user_status,
                            role task_user_role
);

CREATE INDEX idx_task_users_taskId_userId ON task_users (taskId, userId);
CREATE INDEX idx_task_users_userId ON task_users (userId);

-- 任务报名审核记录表
CREATE TABLE task_user_audits (
                                  id SERIAL PRIMARY KEY,
                                  taskId INTEGER REFERENCES tasks(id),
                                  userId INTEGER REFERENCES users(id),
                                  auditUserId INTEGER REFERENCES users(id),
                                  createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                  status task_user_status,
                                  reason VARCHAR(255)
);

-- 任务执行表
CREATE TABLE task_runs (
                           id SERIAL PRIMARY KEY,
                           taskId INTEGER REFERENCES tasks(id),
                           createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                           startAt TIMESTAMP WITH TIME ZONE,
                           endAt TIMESTAMP WITH TIME ZONE,
                           duration INTERVAL,
                           status task_status
);

-- 任务执行人记录表
CREATE TABLE task_run_users (
                                id SERIAL PRIMARY KEY,
                                taskId INTEGER REFERENCES tasks(id),
                                taskRunId INTEGER REFERENCES task_runs(id),
                                userId INTEGER REFERENCES users(id),
                                duration INTERVAL,
                                createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                startAt TIMESTAMP WITH TIME ZONE,
                                endAt TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_task_run_users_taskId ON task_run_users (taskId);
CREATE INDEX idx_task_run_users_createdAt ON task_run_users (createdAt);
CREATE INDEX idx_task_run_users_user_id ON task_run_users (userId);

-- 任务执行日志表
CREATE TABLE task_run_logs (
                               id SERIAL PRIMARY KEY,
                               taskId INTEGER REFERENCES tasks(id),
                               taskRunId INTEGER REFERENCES task_runs(id),
                               userId INTEGER REFERENCES users(id),
                               content TEXT,
                               createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                               images varchar(256)[],
                               videos varchar(256)[]
);

CREATE INDEX idx_taskRunLogs_taskRunId ON task_run_logs (taskRunId);
CREATE INDEX idx_taskRunLogs_createdAt ON task_run_logs (createdAt);

-- 通知表
CREATE TABLE notices (
                         id SERIAL PRIMARY KEY,
                         userId INTEGER REFERENCES users(id),
                         content TEXT NOT NULL,
                         status notice_status,
                         createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notices_created_userId ON notices (createdAt , userId);

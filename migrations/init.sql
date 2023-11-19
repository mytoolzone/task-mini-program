-- 创建 ENUM 类型
CREATE TYPE user_status AS ENUM ('active', 'frozen');
CREATE TYPE task_status AS ENUM ('create','publish', 'audit_fail','audit_pass', 'pending', 'running', 'paused', 'finished', 'canceled', 'deleted');
CREATE TYPE task_user_status AS ENUM ('apply', 'approved', 'rejected');
CREATE TYPE task_user_role AS ENUM ('leader', 'member', 'recorder', 'none');
CREATE TYPE notice_status AS ENUM ('unread', 'read');

-- 用户表
CREATE TABLE users (
           id SERIAL PRIMARY KEY,
           username VARCHAR(255) NOT NULL,
           password VARCHAR(255) NOT NULL,
           email VARCHAR(255),
           phone VARCHAR(20),
           status user_status,
           created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP WITH TIME ZONE,
           deleted_at TIMESTAMP WITH TIME ZONE,
           openid VARCHAR(255),
           ext JSONB
);

-- 用户角色表
CREATE TABLE user_roles (
                            id SERIAL PRIMARY KEY,
                            user_id INTEGER REFERENCES users(id),
                            role VARCHAR(20) NOT NULL,
                            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP WITH TIME ZONE
);

-- 任务表
CREATE TABLE tasks (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) NOT NULL,
                       create_by INTEGER REFERENCES users(id),
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP WITH TIME ZONE,
                       finished_at TIMESTAMP WITH TIME ZONE,
                       deleted_at TIMESTAMP WITH TIME ZONE,
                       describe VARCHAR(1024),
                       require VARCHAR(1024),
                       location POINT,
                       status task_status
);

-- 任务执行人角色表
CREATE TABLE task_users (
                            id SERIAL PRIMARY KEY,
                            task_id INTEGER REFERENCES tasks(id),
                            user_id INTEGER REFERENCES users(id),
                            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                            status task_user_status,
                            role task_user_role
);

CREATE INDEX idx_task_users_task_id_userId ON task_users (task_id, user_id);
CREATE INDEX idx_task_users_userId ON task_users (user_id);

-- 任务报名审核记录表
CREATE TABLE task_user_audits (
                                  id SERIAL PRIMARY KEY,
                                  task_id INTEGER REFERENCES tasks(id),
                                  user_id INTEGER REFERENCES users(id),
                                  audit_user_id INTEGER REFERENCES users(id),
                                  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                  status task_user_status,
                                  reason VARCHAR(255)
);

-- 任务执行表
CREATE TABLE task_runs (
                           id SERIAL PRIMARY KEY,
                           task_id INTEGER REFERENCES tasks(id),
                           created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                           start_at TIMESTAMP WITH TIME ZONE,
                           endAt TIMESTAMP WITH TIME ZONE,
                           duration INTERVAL,
                           status task_status
);

-- 任务执行人记录表
CREATE TABLE task_run_users (
                id SERIAL PRIMARY KEY,
                task_id INTEGER REFERENCES tasks(id),
                task_run_id INTEGER REFERENCES task_runs(id),
                user_id INTEGER REFERENCES users(id),
                duration INTERVAL,
                created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                start_at TIMESTAMP WITH TIME ZONE,
                finished_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_task_run_users_task_id ON task_run_users (task_id);
CREATE INDEX idx_task_run_users_created_at ON task_run_users (created_at);
CREATE INDEX idx_task_run_users_user_id ON task_run_users (user_id);

-- 任务执行日志表
CREATE TABLE task_run_logs (
                               id SERIAL PRIMARY KEY,
                               task_id INTEGER REFERENCES tasks(id),
                               task_run_id INTEGER REFERENCES task_runs(id),
                               user_id INTEGER REFERENCES users(id),
                               content TEXT,
                               created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                               images varchar(256)[],
                               videos varchar(256)[]
);

CREATE INDEX idx_taskRunLogs_task_run_id ON task_run_logs (task_run_id);
CREATE INDEX idx_taskRunLogs_created_at ON task_run_logs (created_at);

-- 通知表
CREATE TABLE notices (
                         id SERIAL PRIMARY KEY,
                         user_id INTEGER REFERENCES users(id),
                         content TEXT NOT NULL,
                         status notice_status,
                         created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notices_created_user_id ON notices (created_at , user_id);

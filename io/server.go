/*
 * Copyright (c) 2019. Aberic - All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package io

import (
	"context"
	"github.com/aberic/lily/api"
)

// LilyAPIServer LilyAPIServer
type LilyAPIServer struct {
}

// GetDatabases 获取数据库集合
func (l *LilyAPIServer) ObtainDatabases(context.Context, *api.ReqDatabases) (*api.RespDatabases, error) {
	return nil, nil
}

// CreateDatabase 新建数据库
func (l *LilyAPIServer) CreateDatabase(context.Context, *api.ReqCreateDatabase) (*api.RespDatabase, error) {
	return nil, nil
}

// CreateForm 创建表
func (l *LilyAPIServer) CreateForm(context.Context, *api.ReqCreateForm) (*api.Resp, error) {
	return nil, nil
}

// CreateKey 新建主键
func (l *LilyAPIServer) CreateKey(context.Context, *api.ReqCreateKey) (*api.Resp, error) {
	return nil, nil
}

// CreateIndex 新建索引
func (l *LilyAPIServer) CreateIndex(context.Context, *api.ReqCreateIndex) (*api.Resp, error) {
	return nil, nil
}

// PutD 新增数据
func (l *LilyAPIServer) PutD(context.Context, *api.ReqPutD) (*api.RespPutD, error) {
	return nil, nil
}

// SetD 新增数据
func (l *LilyAPIServer) SetD(context.Context, *api.ReqSetD) (*api.RespSetD, error) {
	return nil, nil
}

// GetD 获取数据
func (l *LilyAPIServer) GetD(context.Context, *api.ReqGetD) (*api.RespGetD, error) {
	return nil, nil
}

// Put 新增数据
func (l *LilyAPIServer) Put(context.Context, *api.ReqPut) (*api.RespPut, error) {
	return nil, nil
}

// Set 新增数据
func (l *LilyAPIServer) Set(context.Context, *api.ReqSet) (*api.RespSet, error) {
	return nil, nil
}

// Get 获取数据
func (l *LilyAPIServer) Get(context.Context, *api.ReqGet) (*api.RespGet, error) {
	return nil, nil
}

// Insert 新增数据
func (l *LilyAPIServer) Insert(context.Context, *api.ReqInsert) (*api.RespInsert, error) {
	return nil, nil
}

// Update 更新数据
func (l *LilyAPIServer) Update(context.Context, *api.ReqUpdate) (*api.Resp, error) {
	return nil, nil
}

// Select 获取数据
func (l *LilyAPIServer) Select(context.Context, *api.ReqSelect) (*api.RespSelect, error) {
	return nil, nil
}

// Delete 删除数据
func (l *LilyAPIServer) Delete(context.Context, *api.ReqDelete) (*api.Resp, error) {
	return nil, nil
}

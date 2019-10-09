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

package lily

import (
	"context"
	"encoding/json"
	"github.com/aberic/lily/api"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/yaml.v3"
)

// APIServer APIServer
type APIServer struct {
	Conf *Conf
}

// GetConf 获取数据库引擎对象
func (l *APIServer) GetConf(ctx context.Context, req *api.ReqConf) (*api.RespConf, error) {
	return &api.RespConf{Code: api.Code_Success, Conf: l.Conf.conf2RPC()}, nil
}

// ObtainDatabases 获取数据库集合
func (l *APIServer) ObtainDatabases(ctx context.Context, req *api.ReqDatabases) (*api.RespDatabases, error) {
	return &api.RespDatabases{Code: api.Code_Success, Databases: l.formatDBs(ObtainLily().GetDatabases())}, nil
}

// ObtainForms 获取数据库表集合
func (l *APIServer) ObtainForms(ctx context.Context, req *api.ReqForms) (*api.RespForms, error) {
	return &api.RespForms{Code: api.Code_Success, Forms: l.formatFormArr(ObtainLily().GetDatabase(req.DatabaseName))}, nil
}

// CreateDatabase 新建数据库
func (l *APIServer) CreateDatabase(ctx context.Context, req *api.ReqCreateDatabase) (*api.RespDatabase, error) {
	var (
		db  Database
		err error
	)
	if db, err = ObtainLily().CreateDatabase(req.Name, req.Comment); nil != err {
		return &api.RespDatabase{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	}
	apiDB := &api.Database{
		Id:      db.getID(),
		Name:    db.getName(),
		Comment: db.getComment(),
		Forms:   l.formatForms(db),
	}
	return &api.RespDatabase{Code: api.Code_Success, Database: apiDB}, nil
}

// CreateForm 创建表
func (l *APIServer) CreateForm(ctx context.Context, req *api.ReqCreateForm) (*api.Resp, error) {
	if err := ObtainLily().CreateForm(req.DatabaseName, req.Name, req.Comment, FormatFormType(req.FormType)); nil != err {
		return &api.Resp{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &api.Resp{Code: api.Code_Success}, nil
}

// CreateKey 新建主键
func (l *APIServer) CreateKey(ctx context.Context, req *api.ReqCreateKey) (*api.Resp, error) {
	if err := ObtainLily().CreateKey(req.DatabaseName, req.FormName, req.KeyStructure); nil != err {
		return &api.Resp{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &api.Resp{Code: api.Code_Success}, nil
}

// CreateIndex 新建索引
func (l *APIServer) CreateIndex(ctx context.Context, req *api.ReqCreateIndex) (*api.Resp, error) {
	if err := ObtainLily().CreateIndex(req.DatabaseName, req.FormName, req.KeyStructure); nil != err {
		return &api.Resp{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &api.Resp{Code: api.Code_Success}, nil
}

// PutD 新增数据
func (l *APIServer) PutD(ctx context.Context, req *api.ReqPutD) (*api.RespPutD, error) {
	var (
		v       interface{}
		hashKey uint64
		err     error
	)
	if err = json.Unmarshal(req.Value, &v); nil == err { // 尝试用json解析
		goto PUT
	}
	if err = yaml.Unmarshal(req.Value, &v); nil == err { // 尝试用yaml解析
		goto PUT
	}
PUT:
	if hashKey, err = ObtainLily().PutD(req.Key, v); nil != err {
		return &api.RespPutD{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &api.RespPutD{Code: api.Code_Success, HashKey: hashKey}, nil
}

// SetD 新增数据
func (l *APIServer) SetD(ctx context.Context, req *api.ReqSetD) (*api.RespSetD, error) {
	var (
		v       interface{}
		hashKey uint64
		err     error
	)
	if err = json.Unmarshal(req.Value, &v); nil == err { // 尝试用json解析
		goto PUT
	}
	if err = yaml.Unmarshal(req.Value, &v); nil == err { // 尝试用yaml解析
		goto PUT
	}
PUT:
	if hashKey, err = ObtainLily().SetD(req.Key, v); nil != err {
		return &api.RespSetD{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &api.RespSetD{Code: api.Code_Success, HashKey: hashKey}, nil
}

// GetD 获取数据
func (l *APIServer) GetD(ctx context.Context, req *api.ReqGetD) (*api.RespGetD, error) {
	var (
		v    interface{}
		data []byte
		err  error
	)
	if v, err = ObtainLily().GetD(req.Key); nil != err {
		return &api.RespGetD{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	}
	if data, err = msgpack.Marshal(v); nil != err {
		return &api.RespGetD{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &api.RespGetD{Code: api.Code_Success, Value: data}, nil
}

// Put 新增数据
func (l *APIServer) Put(ctx context.Context, req *api.ReqPut) (*api.RespPut, error) {
	var (
		v       interface{}
		hashKey uint64
		err     error
	)
	if err = json.Unmarshal(req.Value, &v); nil == err { // 尝试用json解析
		goto PUT
	}
	if err = yaml.Unmarshal(req.Value, &v); nil == err { // 尝试用yaml解析
		goto PUT
	}
PUT:
	if hashKey, err = ObtainLily().Put(req.DatabaseName, req.FormName, req.Key, v); nil != err {
		return &api.RespPut{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &api.RespPut{Code: api.Code_Success, HashKey: hashKey}, nil
}

// Set 新增数据
func (l *APIServer) Set(ctx context.Context, req *api.ReqSet) (*api.RespSet, error) {
	var (
		v       interface{}
		hashKey uint64
		err     error
	)
	if err = json.Unmarshal(req.Value, &v); nil == err { // 尝试用json解析
		goto PUT
	}
	if err = yaml.Unmarshal(req.Value, &v); nil == err { // 尝试用yaml解析
		goto PUT
	}
PUT:
	if hashKey, err = ObtainLily().Set(req.DatabaseName, req.FormName, req.Key, v); nil != err {
		return &api.RespSet{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &api.RespSet{Code: api.Code_Success, HashKey: hashKey}, nil
}

// Get 获取数据
func (l *APIServer) Get(ctx context.Context, req *api.ReqGet) (*api.RespGet, error) {
	var (
		v    interface{}
		data []byte
		err  error
	)
	if v, err = ObtainLily().Get(req.DatabaseName, req.FormName, req.Key); nil != err {
		return &api.RespGet{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	}
	if data, err = msgpack.Marshal(v); nil != err {
		return &api.RespGet{Code: api.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &api.RespGet{Code: api.Code_Success, Value: data}, nil
}

// Insert 新增数据
func (l *APIServer) Insert(ctx context.Context, req *api.ReqInsert) (*api.RespInsert, error) {
	return nil, nil
}

// Update 更新数据
func (l *APIServer) Update(ctx context.Context, req *api.ReqUpdate) (*api.Resp, error) {
	return nil, nil
}

// Select 获取数据
func (l *APIServer) Select(ctx context.Context, req *api.ReqSelect) (*api.RespSelect, error) {
	return nil, nil
}

// Delete 删除数据
func (l *APIServer) Delete(ctx context.Context, req *api.ReqDelete) (*api.Resp, error) {
	return nil, nil
}

func (l *APIServer) formatDBs(dbs []Database) []*api.Database {
	var respDBs []*api.Database
	for _, db := range dbs {
		respDBs = append(respDBs, &api.Database{Id: db.getID(), Name: db.getName(), Comment: db.getComment(), Forms: l.formatForms(db)})
	}
	return respDBs
}

func (l *APIServer) formatForms(db Database) map[string]*api.Form {
	var fms = make(map[string]*api.Form)
	for _, form := range db.getForms() {
		fms[form.getID()] = &api.Form{
			Id:       form.getID(),
			Name:     form.getName(),
			Comment:  form.getComment(),
			FormType: FormatFormType2API(form.getFormType()),
			Indexes:  l.formatIndexes(form),
		}
	}
	return fms
}

func (l *APIServer) formatFormArr(db Database) []*api.Form {
	var fms []*api.Form
	for _, form := range db.getForms() {
		fms = append(fms, &api.Form{
			Id:       form.getID(),
			Name:     form.getName(),
			Comment:  form.getComment(),
			FormType: FormatFormType2API(form.getFormType()),
			Indexes:  l.formatIndexes(form),
		})
	}
	return fms
}

// FormatFormType 通过api表类型获取数据库表类型
func FormatFormType(ft api.FormType) string {
	switch ft {
	default:
		return FormTypeDoc
	case api.FormType_SQL:
		return FormTypeSQL
	}
}

// FormatFormType2API 通过数据库表类型获取api表类型
func FormatFormType2API(ft string) api.FormType {
	switch ft {
	default:
		return api.FormType_Doc
	case FormTypeSQL:
		return api.FormType_SQL
	}
}

func (l *APIServer) formatIndexes(fm Form) map[string]*api.Index {
	var idx = make(map[string]*api.Index)
	for _, index := range fm.getIndexes() {
		idx[index.getID()] = &api.Index{Id: index.getID(), Primary: index.isPrimary(), KeyStructure: index.getKeyStructure()}
	}
	return idx
}

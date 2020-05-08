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
	"github.com/aberic/lily/api"
)

//var ServerURL = "localhost:19877"

// GetConf 获取数据库集合
func GetConf(serverURL string) (*Conf, error) {
	res, err := getConf(serverURL, &api.ReqConf{})
	conf := &Conf{}
	conf.rpc2Conf(res.(*api.RespConf).Conf)
	return conf, err
}

// ObtainDatabases 获取数据库集合
func ObtainDatabases(serverURL string) (*api.RespDatabases, error) {
	res, err := obtainDatabases(serverURL, &api.ReqDatabases{})
	if nil != err {
		return nil, err
	}
	if nil == res {
		return &api.RespDatabases{Code: api.Code_Success, Databases: []*api.Database{}}, nil
	}
	return res.(*api.RespDatabases), nil
}

// ObtainForms 获取数据库表集合
func ObtainForms(serverURL, dbName string) (*api.RespForms, error) {
	res, err := obtainForms(serverURL, &api.ReqForms{DatabaseName: dbName})
	if nil != err {
		return nil, err
	}
	if nil == res {
		return &api.RespForms{Code: api.Code_Success, Forms: []*api.Form{}}, nil
	}
	return res.(*api.RespForms), nil
}

// CreateDatabase 创建数据库
func CreateDatabase(serverURL, name, comment string) error {
	_, err := createDatabase(serverURL, &api.ReqCreateDatabase{Name: name, Comment: comment})
	return err
}

// CreateTable 创建表
func CreateTable(serverURL, dbName, name, comment string) error {
	_, err := createForm(serverURL, &api.ReqCreateForm{DatabaseName: dbName, Name: name, Comment: comment, FormType: FormatFormType2API(FormTypeSQL)})
	return err
}

// CreateDoc 创建文档
func CreateDoc(serverURL, dbName, name, comment string) error {
	_, err := createForm(serverURL, &api.ReqCreateForm{DatabaseName: dbName, Name: name, Comment: comment, FormType: FormatFormType2API(FormTypeDoc)})
	return err
}

// PutD 新增数据
func PutD(serverURL, key, value string) (*api.RespPutD, error) {
	res, err := putD(serverURL, &api.ReqPutD{Key: key, Value: []byte(value)})
	if nil != err {
		return nil, err
	}
	return res.(*api.RespPutD), err
}

// SetD 新增数据
func SetD(serverURL, key, value string) (*api.RespSetD, error) {
	res, err := setD(serverURL, &api.ReqSetD{Key: key, Value: []byte(value)})
	if nil != err {
		return nil, err
	}
	return res.(*api.RespSetD), err
}

// GetD 获取数据
func GetD(serverURL, key string) (*api.RespGetD, error) {
	res, err := getD(serverURL, &api.ReqGetD{Key: key})
	if nil != err {
		return nil, err
	}
	return res.(*api.RespGetD), err
}

// Put 新增数据
func Put(serverURL, databaseName, formName, key, value string) (*api.RespPut, error) {
	res, err := put(serverURL, &api.ReqPut{DatabaseName: databaseName, FormName: formName, Key: key, Value: []byte(value)})
	if nil != err {
		return nil, err
	}
	return res.(*api.RespPut), err
}

// Set 新增数据
func Set(serverURL, databaseName, formName, key, value string) (*api.RespSet, error) {
	res, err := set(serverURL, &api.ReqSet{DatabaseName: databaseName, FormName: formName, Key: key, Value: []byte(value)})
	if nil != err {
		return nil, err
	}
	return res.(*api.RespSet), err
}

// Get 获取数据
func Get(serverURL, databaseName, formName, key string) (*api.RespGet, error) {
	res, err := get(serverURL, &api.ReqGet{DatabaseName: databaseName, FormName: formName, Key: key})
	if nil != err {
		return nil, err
	}
	return res.(*api.RespGet), err
}

// Select 获取数据
func Select(serverURL, databaseName, formName string, selector *api.Selector) (*api.RespSelect, error) {
	res, err := query(serverURL, &api.ReqSelect{DatabaseName: databaseName, FormName: formName, Selector: selector})
	return res.(*api.RespSelect), err
}

// Remove 删除数据
func Remove(serverURL, databaseName, formName, key string) (*api.Resp, error) {
	res, err := remove(serverURL, &api.ReqRemove{DatabaseName: databaseName, FormName: formName, Key: key})
	return res.(*api.Resp), err
}

// Delete 删除数据
func Delete(serverURL, databaseName, formName string, selector *api.Selector) (*api.Resp, error) {
	res, err := del(serverURL, &api.ReqDelete{DatabaseName: databaseName, FormName: formName, Selector: selector})
	return res.(*api.Resp), err
}

// getConf 获取数据库引擎对象
func getConf(serverURL string, req *api.ReqConf) (interface{}, error) {
	return getClient(serverURL).GetConf(context.Background(), req)
}

// obtainDatabases 获取数据库集合
func obtainDatabases(serverURL string, req *api.ReqDatabases) (interface{}, error) {
	return getClient(serverURL).ObtainDatabases(context.Background(), req)
}

// obtainForms 获取数据库表集合
func obtainForms(serverURL string, req *api.ReqForms) (interface{}, error) {
	return getClient(serverURL).ObtainForms(context.Background(), req)
}

// createDatabase 创建数据库
func createDatabase(serverURL string, req *api.ReqCreateDatabase) (interface{}, error) {
	return getClient(serverURL).CreateDatabase(context.Background(), req)
}

// createForm 创建表
func createForm(serverURL string, req *api.ReqCreateForm) (interface{}, error) {
	return getClient(serverURL).CreateForm(context.Background(), req)
}

// putD 新增数据
func putD(serverURL string, req *api.ReqPutD) (interface{}, error) {
	return getClient(serverURL).PutD(context.Background(), req)
}

// setD 新增数据
func setD(serverURL string, req *api.ReqSetD) (interface{}, error) {
	return getClient(serverURL).SetD(context.Background(), req)
}

// getD 获取数据
func getD(serverURL string, req *api.ReqGetD) (interface{}, error) {
	return getClient(serverURL).GetD(context.Background(), req)
}

// put 新增数据
func put(serverURL string, req *api.ReqPut) (interface{}, error) {
	return getClient(serverURL).Put(context.Background(), req)
}

// set 新增数据
func set(serverURL string, req *api.ReqSet) (interface{}, error) {
	return getClient(serverURL).Set(context.Background(), req)
}

// get 获取数据
func get(serverURL string, req *api.ReqGet) (interface{}, error) {
	return getClient(serverURL).Get(context.Background(), req)
}

// query 获取数据
func query(serverURL string, req *api.ReqSelect) (interface{}, error) {
	return getClient(serverURL).Select(context.Background(), req)
}

// remove 删除数据
func remove(serverURL string, req *api.ReqRemove) (interface{}, error) {
	return getClient(serverURL).Remove(context.Background(), req)
}

// del 删除数据
func del(serverURL string, req *api.ReqDelete) (interface{}, error) {
	return getClient(serverURL).Delete(context.Background(), req)
}

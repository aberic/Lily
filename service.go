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
	"google.golang.org/grpc"
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

// Insert 新增数据
func Insert(serverURL, databaseName, formName string, value interface{}) (*api.RespInsert, error) {
	// todo
	res, err := insert(serverURL, &api.ReqInsert{})
	return res.(*api.RespInsert), err
}

// Update 更新数据
func Update(serverURL string) (*api.Resp, error) {
	// todo
	res, err := update(serverURL, &api.ReqUpdate{})
	return res.(*api.Resp), err
}

// Select 获取数据
func Select(serverURL string) (*api.RespSelect, error) {
	// todo
	res, err := query(serverURL, &api.ReqSelect{})
	return res.(*api.RespSelect), err
}

// Delete 删除数据
func Delete(serverURL string) (*api.Resp, error) {
	// todo
	res, err := delete(serverURL, &api.ReqDelete{})
	return res.(*api.Resp), err
}

// getConf 获取数据库引擎对象
func getConf(serverURL string, req *api.ReqConf) (interface{}, error) {
	return rpc(serverURL, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *api.RespConf
			err    error
		)
		// 创建gRPC客户端
		c := api.NewLilyAPIClient(conn)
		// 客户端向gRPC服务端发起请求
		if result, err = c.GetConf(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// obtainDatabases 获取数据库集合
func obtainDatabases(serverURL string, req *api.ReqDatabases) (interface{}, error) {
	return rpc(serverURL, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *api.RespDatabases
			err    error
		)
		// 创建gRPC客户端
		c := api.NewLilyAPIClient(conn)
		// 客户端向gRPC服务端发起请求
		if result, err = c.ObtainDatabases(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// obtainForms 获取数据库表集合
func obtainForms(serverURL string, req *api.ReqForms) (interface{}, error) {
	return rpc(serverURL, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *api.RespForms
			err    error
		)
		// 创建gRPC客户端
		c := api.NewLilyAPIClient(conn)
		// 客户端向gRPC服务端发起请求
		if result, err = c.ObtainForms(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// createDatabase 创建数据库
func createDatabase(serverURL string, req *api.ReqCreateDatabase) (interface{}, error) {
	return rpc(serverURL, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *api.RespDatabase
			err    error
		)
		// 创建gRPC客户端
		c := api.NewLilyAPIClient(conn)
		// 客户端向gRPC服务端发起请求
		if result, err = c.CreateDatabase(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// createForm 创建表
func createForm(serverURL string, req *api.ReqCreateForm) (interface{}, error) {
	return rpc(serverURL, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *api.Resp
			err    error
		)
		// 创建gRPC客户端
		c := api.NewLilyAPIClient(conn)
		// 客户端向gRPC服务端发起请求
		if result, err = c.CreateForm(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// putD 新增数据
func putD(serverURL string, req *api.ReqPutD) (interface{}, error) {
	return rpc(serverURL, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *api.RespPutD
			err    error
		)
		// 创建gRPC客户端
		c := api.NewLilyAPIClient(conn)
		// 客户端向gRPC服务端发起请求
		if result, err = c.PutD(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// setD 新增数据
func setD(serverURL string, req *api.ReqSetD) (interface{}, error) {
	return rpc(serverURL, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *api.RespSetD
			err    error
		)
		// 创建gRPC客户端
		c := api.NewLilyAPIClient(conn)
		// 客户端向gRPC服务端发起请求
		if result, err = c.SetD(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// getD 获取数据
func getD(serverURL string, req *api.ReqGetD) (interface{}, error) {
	return rpc(serverURL, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *api.RespGetD
			err    error
		)
		// 创建gRPC客户端
		c := api.NewLilyAPIClient(conn)
		// 客户端向gRPC服务端发起请求
		if result, err = c.GetD(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// put 新增数据
func put(serverURL string, req *api.ReqPut) (interface{}, error) {
	return rpc(serverURL, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *api.RespPut
			err    error
		)
		// 创建gRPC客户端
		c := api.NewLilyAPIClient(conn)
		// 客户端向gRPC服务端发起请求
		if result, err = c.Put(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// set 新增数据
func set(serverURL string, req *api.ReqSet) (interface{}, error) {
	return rpc(serverURL, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *api.RespSet
			err    error
		)
		// 创建gRPC客户端
		c := api.NewLilyAPIClient(conn)
		// 客户端向gRPC服务端发起请求
		if result, err = c.Set(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// get 获取数据
func get(serverURL string, req *api.ReqGet) (interface{}, error) {
	return rpc(serverURL, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *api.RespGet
			err    error
		)
		// 创建gRPC客户端
		c := api.NewLilyAPIClient(conn)
		// 客户端向gRPC服务端发起请求
		if result, err = c.Get(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// insert 新增数据
func insert(serverURL string, req *api.ReqInsert) (interface{}, error) {
	return rpc(serverURL, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *api.RespInsert
			err    error
		)
		// 创建gRPC客户端
		c := api.NewLilyAPIClient(conn)
		// 客户端向gRPC服务端发起请求
		if result, err = c.Insert(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// update 更新数据
func update(serverURL string, req *api.ReqUpdate) (interface{}, error) {
	return rpc(serverURL, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *api.Resp
			err    error
		)
		// 创建gRPC客户端
		c := api.NewLilyAPIClient(conn)
		// 客户端向gRPC服务端发起请求
		if result, err = c.Update(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// query 获取数据
func query(serverURL string, req *api.ReqSelect) (interface{}, error) {
	return rpc(serverURL, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *api.RespSelect
			err    error
		)
		// 创建gRPC客户端
		c := api.NewLilyAPIClient(conn)
		// 客户端向gRPC服务端发起请求
		if result, err = c.Select(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

// delete 删除数据
func delete(serverURL string, req *api.ReqDelete) (interface{}, error) {
	return rpc(serverURL, func(conn *grpc.ClientConn) (interface{}, error) {
		var (
			result *api.Resp
			err    error
		)
		// 创建gRPC客户端
		c := api.NewLilyAPIClient(conn)
		// 客户端向gRPC服务端发起请求
		if result, err = c.Delete(context.Background(), req); nil != err {
			return nil, err
		}
		return result, nil
	})
}

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

// GetDatabases 获取数据库集合
func ObtainDatabases(clientURL string) (*api.RespDatabases, error) {
	res, err := obtainDatabases(clientURL, &api.ReqDatabases{})
	return res.(*api.RespDatabases), err
}

// obtainDatabases 获取数据库集合
func obtainDatabases(url string, req *api.ReqDatabases) (interface{}, error) {
	return rpc(url, func(conn *grpc.ClientConn) (interface{}, error) {
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

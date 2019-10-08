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
	"fmt"
	"github.com/aberic/lily/api"
	"google.golang.org/grpc"
	"net"
	"strings"
)

// RPCListener 启动rpc监听
func RPCListener(conf *Conf) {
	var (
		listener net.Listener
		err      error
	)

	fmt.Println(strings.Join([]string{"Listen announces on the local network address with port: ", conf.Port}, ""))
	if listener, err = net.Listen("tcp", strings.Join([]string{":", conf.Port}, "")); nil != err {
		panic(err)
	}
	fmt.Println("creates a gRPC server")
	server := grpc.NewServer()
	fmt.Println("register gRPC listener")
	api.RegisterLilyAPIServer(server, &APIServer{Conf: conf})
	fmt.Println("OFF")
	if err = server.Serve(listener); nil != err {
		panic(err)
	}
}

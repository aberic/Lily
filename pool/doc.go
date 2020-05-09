/*
 * Copyright (c) 2020. Aberic - All Rights Reserved.
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

// Package pool 负责监听对客户端的各种请求，接收请求，转发请求到目标模块。
//
// 每个成功连接的客户请求都会被创建或分配一个线程，该线程负责与客户端通信，
//
// 接收客户端发送的命令，传递处理后的结果信息等。
package pool

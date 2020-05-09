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

// Package cache 共享在内存中的完全一样的查询语句和对应的执行选择器。
//
// 如果完全相同的查询在缓存命中，数据库会直接读取其在缓存中关联的执行选择器并检索结果。
//
// 缓存是会话间共享的，所以为一个客户生成的结果集也能为另一个客户所用。
package cache

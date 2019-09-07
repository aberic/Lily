/*
 * Copyright (c) 2019.. Aberic - All Rights Reserved.
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

// Package lily 数据库
//
// 存储结构 {dataDir}/data/{dataName}/{formName}/{formName}.dat/idx...
//
// {dataDir}/ Lily 服务工作目录
//
// {dataDir}/data/ Lily 服务数据目录，目录下为已创建数据库目录集合
//
// {dataDir}/data/{dataName}.../ 数据库目录，目录下为已创建表目录集合
//
// {dataDir}/data/{dataName}.../{formName}.../ 表目录，目录下为表头部Hash数组对应数据目录集合以及索引目录集合
//
// {dataDir}/data/{dataName}.../{formName}.../{catalog}.../ 表头部Hash数组对应索引数据文件集合
//
// {dataDir}/data/{dataName}.../{formName}.../[0, 1, ... , 15]/ 表头部Hash数组对应数据文件集合
//
// 索引格式：5位key + 16位md5后key + 5位起始seek + 4位持续seek
//
// 索引起始seek最大值为1073741824，即每一个索引文件最大1G，新增一条索引超过1G则新开新的索引文件
//
// 索引持续seek最大值为16777216，即每一条存储value对象允许最大存储16777216个字节
package lily

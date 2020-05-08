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

// 起始语句解析内容
const (
	firstShow   = "show"
	firstUse    = "use"
	firstCreate = "create"
	firstPutD   = "putD"
	firstSetD   = "setD"
	firstGetD   = "getD"
	firstPut    = "put"
	firstSet    = "set"
	firstGet    = "get"
	firstSelect = "select"
	firstRemove = "remove"
	firstDelete = "delete"
)

// SHOW 语句解析内容
const (
	firstShowConf      = "conf"
	firstShowDatabases = "databases"
	firstShowForms     = "forms"
)

// Create 语句解析内容
const (
	firstCreateDatabase = "database"
	firstCreateTable    = "table"
	firstCreateDoc      = "doc"
)

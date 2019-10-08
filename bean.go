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

// DTODatabase 数据库对象
type DTODatabase struct {
	Name    string // Name 数据库名称，根据需求可以随时变化
	Comment string // Comment 描述
}

// DTOForm 数据库表对象
type DTOForm struct {
	Name    string // Name 数据库名称，根据需求可以随时变化
	Comment string // Comment 描述
	Type    string // Type 类型
}

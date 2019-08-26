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

package Lily

type query struct {
	Scopes     []scope     `json:"scopes"`
	Conditions []condition `json:"conditions"`
	Matches    []match     `json:"matches"`
	Skip       int32       `json:"skip"`
	Limit      int32       `json:"limit"`
}

// 范围
type scope struct {
	Param string `json:"param"`
	Start int32  `json:"start"`
	End   int32  `json:"end"`
}

// 条件
type condition struct {
	Param string `json:"param"`
	Cond  string `json:"cond"` // gt/lt/eq/dif
}

type match struct {
	Param string `json:"param"`
	Value string `json:"value"`
}

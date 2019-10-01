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

package cmd

import (
	"errors"
	"github.com/aberic/gnomon"
	"github.com/aberic/lily"
	"github.com/aberic/lily/dto"
	"github.com/modood/table"
	"strings"
)

var (
	sqlSyntaxErr = errors.New("sql syntax error")
)

func analysis(sql string) error {
	sql = gnomon.String().SingleSpace(sql)
	array := strings.Split(sql, " ")
	return first(array)
}

func first(array []string) error {
	switch strings.ToLower(array[0]) {
	default:
		return sqlSyntaxErr
	case "show":
		return show(array[1])
	}
}

func show(data string) error {
	if strings.ToLower(data) == "database" {
		dbs, err := lily.ObtainDatabases("localhost:19877")
		if nil != err {
			return err
		}
		var dbDTO []*dto.Database
		for _, db := range dbs.Databases {
			dbDTO = append(dbDTO, &dto.Database{Name: db.Name, Comment: db.Comment})
		}
		table.Output(dbDTO)
		return nil
	}
	return sqlSyntaxErr
}

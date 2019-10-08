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
	sqlSyntaxErr                   = err("sql syntax error")
	sqlDatabaseIsNilErr            = err("database is nil, you should use database first")
	sqlSyntaxParamsCountInvalidErr = syntaxErr("params count is invalid")
)

type sql struct {
	serverURL    string // serverURL 链接数据库地址，如'localhost:19877'
	databaseName string // databaseName 当前操作数据库名称
}

func err(errStr string) error {
	return errors.New(errStr)
}

func syntaxErr(errStr string) error {
	return errors.New(strings.Join([]string{"sql syntax error", errStr}, ": "))
}

func executeErr(errStr string) error {
	return errors.New(strings.Join([]string{"sql execute error", errStr}, ": "))
}

func (s *sql) analysis(sql string) error {
	sql = gnomon.String().SingleSpace(sql)
	array := strings.Split(sql, " ")
	if len(array) < 1 {
		return sqlSyntaxParamsCountInvalidErr
	}
	return s.first(array)
}

func (s *sql) first(array []string) error {
	switch strings.ToLower(array[0]) {
	default:
		return sqlSyntaxErr
	case "show":
		return s.show(array)
	case "use":
		return s.use(array)
	case "create":
		return s.create(array)
		//case "putD":
		//	return s.putD(array)
		//case "setD":
		//	return s.setD(array)
		//case "getD":
		//	return s.getD(array)
		//case "put":
		//	return s.put(array)
		//case "set":
		//	return s.set(array)
		//case "get":
		//	return s.get(array)
	}
}

// show show database
func (s *sql) show(array []string) error {
	if len(array) != 2 {
		return sqlSyntaxParamsCountInvalidErr
	}
	switch strings.ToLower(array[1]) {
	default:
		return sqlSyntaxErr
	case "conf":
		var confs []*lily.Conf
		conf, err := lily.GetConf(s.serverURL)
		if nil != err {
			return executeErr(err.Error())
		}
		confs = append(confs, conf)
		table.Output(confs)
		return nil
	case "databases":
		dbs, err := lily.ObtainDatabases(s.serverURL)
		if nil != err {
			return executeErr(err.Error())
		}
		var dbDTO []*dto.Database
		for _, db := range dbs.Databases {
			dbDTO = append(dbDTO, &dto.Database{Name: db.Name, Comment: db.Comment})
		}
		table.Output(dbDTO)
		return nil
	case "forms":
		if gnomon.String().IsEmpty(s.databaseName) {
			return sqlDatabaseIsNilErr
		}
		fms, err := lily.ObtainForms(s.serverURL, s.databaseName)
		if nil != err {
			return executeErr(err.Error())
		}
		var fmDTO []*dto.Form
		for _, fm := range fms.Forms {
			fmDTO = append(fmDTO, &dto.Form{Name: fm.Name, Comment: fm.Comment, Type: lily.FormatFormType(fm.FormType)})
		}
		table.Output(fmDTO)
		return nil
	}
}

func (s *sql) use(array []string) error {
	if len(array) != 2 {
		return sqlSyntaxParamsCountInvalidErr
	}
	dbs, err := lily.ObtainDatabases(s.serverURL)
	if nil != err {
		return executeErr(err.Error())
	}
	have := false
	for _, db := range dbs.Databases {
		if db.Name == array[1] {
			have = true
			break
		}
	}
	if have {
		s.databaseName = array[1]
		return nil
	}
	return executeErr("database not found")
}

func (s *sql) create(array []string) error {
	if len(array) < 3 {
		return sqlSyntaxParamsCountInvalidErr
	}
	switch array[1] {
	default:
		return sqlSyntaxErr
	case "database":
		return s.createDatabase(array)
	case "table":
		return s.createTable(array)
	case "doc":
		return s.createDoc(array)
		//case "key":
		//	return s.createKey(array)
	}
}

func (s *sql) createDatabase(array []string) error {
	if len(array) == 3 {
		return lily.CreateDatabase(s.serverURL, array[2], "")
	} else if len(array) > 3 {
		comment := array[3:]
		return lily.CreateDatabase(s.serverURL, array[2], strings.Join(comment, " "))
	}
	return sqlSyntaxErr
}

func (s *sql) createTable(array []string) error {
	if len(array) == 3 {
		return lily.CreateTable(s.serverURL, s.databaseName, array[2], "")
	} else if len(array) > 3 {
		comment := array[3:]
		return lily.CreateTable(s.serverURL, s.databaseName, array[2], strings.Join(comment, " "))
	}
	return sqlSyntaxErr
}

func (s *sql) createDoc(array []string) error {
	if len(array) == 3 {
		return lily.CreateDoc(s.serverURL, s.databaseName, array[2], "")
	} else if len(array) > 3 {
		comment := array[3:]
		return lily.CreateDoc(s.serverURL, s.databaseName, array[2], strings.Join(comment, " "))
	}
	return sqlSyntaxErr
}

//func (s *sql) putD(array []string) error {
//	if len(array) < 3 {
//		return sqlSyntaxParamsCountInvalidErr
//	}
//	valueStr := strings.Join(array[2:], " ")
//	lily.PutD(s.serverURL)
//}

//func put(databaseName, formName, key string, value interface{}) error {
//
//}

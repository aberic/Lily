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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aberic/gnomon"
	"github.com/aberic/lily/api"
	"github.com/modood/table"
	"github.com/vmihailenco/msgpack"
	"strings"
)

var (
	sqlSyntaxErr                   = customErr("sql syntax error")
	sqlDatabaseIsNilErr            = customErr("database is nil, you should use database first")
	sqlSyntaxParamsCountInvalidErr = syntaxErr("params count is invalid")
)

type sql struct {
	serverURL    string // serverURL 链接数据库地址，如'localhost:19877'
	databaseName string // databaseName 当前操作数据库名称
}

func customErr(errStr string) error {
	return errors.New(errStr)
}

func syntaxErr(errStr string) error {
	return errors.New(strings.Join([]string{"sql syntax error", errStr}, ": "))
}

func executeErr(errStr string) error {
	return errors.New(strings.Join([]string{"sql execute error", errStr}, ": "))
}

func (s *sql) analysis(sql string) error {
	sql = gnomon.StringSingleSpace(sql)
	array := strings.Split(sql, " ")
	if len(array) < 1 {
		return sqlSyntaxParamsCountInvalidErr
	}
	return s.first(array)
}

func (s *sql) first(array []string) error {
	switch array[0] {
	default:
		return syntaxErr(strings.Join(array, " "))
	case firstShow:
		return s.show(array)
	case firstUse:
		return s.use(array)
	case firstCreate:
		return s.create(array)
	case firstPutD:
		return s.putD(array)
	case firstSetD:
		return s.setD(array)
	case firstGetD:
		return s.getD(array)
	case firstPut:
		return s.put(array)
	case firstSet:
		return s.set(array)
	case firstGet:
		return s.get(array)
	case firstSelect:
		return s.query(array)
	case firstRemove:
		return s.remove(array)
	case firstDelete:
		return s.delete(array)
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
	case firstShowConf:
		var configs []*Conf
		conf, err := GetConf(s.serverURL)
		if nil != err {
			return executeErr(err.Error())
		}
		configs = append(configs, conf)
		table.Output(configs)
		return nil
	case firstShowDatabases:
		dbs, err := ObtainDatabases(s.serverURL)
		if nil != err {
			return executeErr(err.Error())
		}
		var dbDTO []*DTODatabase
		for _, db := range dbs.Databases {
			dbDTO = append(dbDTO, &DTODatabase{Name: db.Name, Comment: db.Comment})
		}
		table.Output(dbDTO)
		return nil
	case firstShowForms:
		if gnomon.StringIsEmpty(s.databaseName) {
			return sqlDatabaseIsNilErr
		}
		fms, err := ObtainForms(s.serverURL, s.databaseName)
		if nil != err {
			return executeErr(err.Error())
		}
		var fmDTO []*DTOForm
		for _, fm := range fms.Forms {
			fmDTO = append(fmDTO, &DTOForm{Name: fm.Name, Comment: fm.Comment, Type: FormatFormType(fm.FormType)})
		}
		if len(fmDTO) == 0 {
			fmDTO = append(fmDTO, &DTOForm{Name: "-", Comment: "-", Type: "-"})
		}
		table.Output(fmDTO)
		return nil
	}
}

func (s *sql) use(array []string) error {
	if len(array) != 2 {
		return sqlSyntaxParamsCountInvalidErr
	}
	dbs, err := ObtainDatabases(s.serverURL)
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
		return CreateDatabase(s.serverURL, array[2], "")
	} else if len(array) > 3 {
		comment := array[3:]
		return CreateDatabase(s.serverURL, array[2], strings.Join(comment, " "))
	}
	return sqlSyntaxErr
}

func (s *sql) createTable(array []string) error {
	if len(array) == 3 {
		return CreateTable(s.serverURL, s.databaseName, array[2], "")
	} else if len(array) > 3 {
		comment := array[3:]
		return CreateTable(s.serverURL, s.databaseName, array[2], strings.Join(comment, " "))
	}
	return sqlSyntaxErr
}

func (s *sql) createDoc(array []string) error {
	if len(array) == 3 {
		return CreateDoc(s.serverURL, s.databaseName, array[2], "")
	} else if len(array) > 3 {
		comment := array[3:]
		return CreateDoc(s.serverURL, s.databaseName, array[2], strings.Join(comment, " "))
	}
	return sqlSyntaxErr
}

func (s *sql) putD(array []string) error {
	if len(array) < 3 {
		return sqlSyntaxParamsCountInvalidErr
	}
	valueStr := strings.Join(array[2:], " ")
	_, err := PutD(s.serverURL, array[1], valueStr)
	if nil != err {
		return err
	}
	return nil
}

func (s *sql) setD(array []string) error {
	if len(array) < 3 {
		return sqlSyntaxParamsCountInvalidErr
	}
	valueStr := strings.Join(array[2:], " ")
	_, err := SetD(s.serverURL, array[1], valueStr)
	if nil != err {
		return err
	}
	return nil
}

func (s *sql) getD(array []string) error {
	if len(array) != 2 {
		return sqlSyntaxParamsCountInvalidErr
	}
	resp, err := GetD(s.serverURL, array[1])
	if nil != err {
		return err
	}
	var v interface{}
	if err = msgpack.Unmarshal(resp.Value, &v); nil != err {
		return err
	}
	switch v.(type) {
	default:
		fmt.Println(v)
	case map[string]interface{}:
		data, err := json.Marshal(v)
		if nil != err {
			return err
		}
		fmt.Println(string(data))
	}
	return nil
}

func (s *sql) put(array []string) error {
	if len(array) < 4 {
		return sqlSyntaxParamsCountInvalidErr
	}
	valueStr := strings.Join(array[3:], " ")
	_, err := Put(s.serverURL, s.databaseName, array[1], array[2], valueStr)
	if nil != err {
		return err
	}
	return nil
}

func (s *sql) set(array []string) error {
	if len(array) < 4 {
		return sqlSyntaxParamsCountInvalidErr
	}
	valueStr := strings.Join(array[3:], " ")
	_, err := Set(s.serverURL, s.databaseName, array[1], array[2], valueStr)
	if nil != err {
		return err
	}
	return nil
}

func (s *sql) get(array []string) error {
	if len(array) != 3 {
		return sqlSyntaxParamsCountInvalidErr
	}
	resp, err := Get(s.serverURL, s.databaseName, array[1], array[2])
	if nil != err {
		return err
	}
	var v interface{}
	if err = msgpack.Unmarshal(resp.Value, &v); nil != err {
		return err
	}
	switch v.(type) {
	default:
		fmt.Println(v)
	case map[string]interface{}:
		data, err := json.Marshal(v)
		if nil != err {
			return err
		}
		fmt.Println(string(data))
	}
	return nil
}

func (s *sql) query(array []string) error {
	if len(array) < 4 {
		return sqlSyntaxParamsCountInvalidErr
	}
	selector, err := s.selector(array)
	if nil != err {
		return err
	}
	_, err = Select(s.serverURL, s.databaseName, array[1], selector)
	return err
}

func (s *sql) remove(array []string) error {
	if len(array) != 3 {
		return sqlSyntaxParamsCountInvalidErr
	}
	_, err := Remove(s.serverURL, s.databaseName, array[1], array[2])
	return err
}

func (s *sql) delete(array []string) error {
	if len(array) < 4 {
		return sqlSyntaxParamsCountInvalidErr
	}
	selector, err := s.selector(array)
	if nil != err {
		return err
	}
	_, err = Select(s.serverURL, s.databaseName, array[1], selector)
	return err
}

func (s *sql) selector(array []string) (*api.Selector, error) {
	return nil, nil
}

//func put(databaseName, formName, key string, value interface{}) error {
//
//}

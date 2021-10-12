package data

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/chenfei531/query-api/model"
)

type SqliteDataManager struct {
	db *gorm.DB
}

type SelectParams struct {
	primaryKey  string
	containPK   bool
	field       []string
	nestedField map[string][]string
}

func NewSqliteDataManager() *SqliteDataManager {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	return &SqliteDataManager{db}
}

func getSelectParams(data interface{}, selectParams []string) (*SelectParams, error) {
	dataType := reflect.TypeOf(data)
	//assuming primary key must be ID
	//(no customized primary key name, no composite primary key)

	//TODO: add field validation
	pk := "ID"
	containPK := true
	var fields []string
	nestedFields := make(map[string][]string)

	for _, s := range selectParams {
		if strings.Contains(s, ".") {
			elems := strings.Split(s, ".")
			if len(elems) != 2 {
				return nil, errors.New("syntax error: only support 1 level of nested data")
			}
			k, v := elems[0], elems[1]
			if _, ok := nestedFields[k]; ok {
				nestedFields[k] = append(nestedFields[k], v)
			} else {
				nestedFields[k] = []string{}
				nestedFields[k] = append(nestedFields[k], elems[1])
			}
		} else {
			f, _ := dataType.FieldByName(s)
			if f.Type.Kind() != reflect.Slice {
				fields = append(fields, s)
			}
		}
	}

	return &SelectParams{primaryKey: pk, containPK: containPK, field: fields, nestedField: nestedFields}, nil
}

func (dm *SqliteDataManager) getDataByParams(data interface{}, p *model.Params) (interface{}, error) {
	dataType := reflect.TypeOf(data)
	result := reflect.New(reflect.SliceOf(dataType)).Interface()
	tx := dm.db.Model(data)
	if len(p.FilterExp) > 0 {
		tx.Where(p.FilterExp, p.FilterArgs)
	}
	if 0 != p.Offset {
		tx.Offset(p.Offset)
	}
	if 0 != p.Limit {
		tx.Limit(p.Limit)
	}
	if len(p.Sort) > 0 {
		tx.Order(p.Sort)
	}
	if len(p.Select) > 0 {
		tx.Select(strings.Join(p.Select, ", "))
	}
	error := tx.Find(result).Error
	return result, error
}

func (dm *SqliteDataManager) GetDataByNameAndParams(name string, p *model.Params) (string, error) {
	data, error := model.GetObjectByName(name)
	if nil != error {
		return "", error
	}
	sp, error := getSelectParams(data, p.Select)
	if nil != error {
		return "", error
	}

	fmt.Printf("primary key: %s\n", sp.primaryKey)
	fmt.Printf("fields: %s\n", sp.field)
	fmt.Printf("map: %s\n", sp.nestedField)

	p.Select = sp.field
	users, error := dm.getDataByParams(data, p)
	if nil != error {
		return "", error
	}
	//get user ids

	//query agents

	//fill user with agents


	/*
	agents, _ := dm.getDataByParams(model.Agent{}, &model.Params{})
	user1 := reflect.ValueOf(users).Elem().Index(0)
	agent1 := reflect.ValueOf(agents).Elem().Index(0)
	test := user1.FieldByName("Agents")
	test.Set(reflect.ValueOf(agents).Elem())
	*/
	b, _ := json.MarshalIndent(users, "", "    ")
	return string(b), nil
}

func (dm *SqliteDataManager) GetUserById(id int) model.User {
	var user model.User
	dm.db.Preload("Agents").First(&user, id)
	return user
}

func (dm *SqliteDataManager) GetAgents(offset int, limit int) []model.Agent {
	var agents []model.Agent
	dm.db.Offset(offset).Limit(limit).Find(&agents)
	return agents
}

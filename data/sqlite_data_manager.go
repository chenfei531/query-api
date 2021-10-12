package data

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

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

func getSelectParams(dataType reflect.Type, selectParams []string) (sp *SelectParams) {
	//assuming primary key must be ID
	//(no customized primary key name, no composite primary key)
	pk := "ID"
	containPK := true
	var fields []string
	nestedFields := make(map[string][]string)

	for _, s := range selectParams {
		if strings.Contains(s, ".") {
			elems := strings.Split(s, ".")
			if len(elems) != 2 {
				panic("syntax error: only support 1 level of nested data")
			}
			if val, ok := nestedFields[elems[0]]; ok {
				nestedFields[elems[0]] = append(val, elems[1])
			} else {
				panic("syntax error: sub field not found")
			}
		} else {
			f, _ := dataType.FieldByName(s)
			if f.Type.Kind() != reflect.Slice {
				fields = append(fields, s)
			} else {
				nestedFields[s] = []string{}
			}
		}
	}

	return &SelectParams{primaryKey: pk, containPK: containPK, field: fields, nestedField: nestedFields}
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
	error := tx.Find(result).Error
	return result, error
}

func (dm *SqliteDataManager) getDataByNameAndParams(name string, p *model.Params) (interface{}, error) {
	data, error := model.GetObjectByName(name)
	if nil != error {
		return nil, error
	}
	return dm.getDataByParams(data, p)
}

func (dm *SqliteDataManager) GetDataByNameAndParams(name string, p *model.Params) (string, error) {
	data, error := model.GetObjectByName(name)
	if nil != error {
		return "", error
	}
	dataType := reflect.TypeOf(data)
	sp := getSelectParams(dataType, p.Select)
	fmt.Printf("primary key: %s\n", sp.primaryKey)
	fmt.Printf("fields: %s\n", sp.field)
	fmt.Printf("map: %s\n", sp.nestedField)

	users, _ := dm.getDataByNameAndParams("User", p)
	agents, _ := dm.getDataByNameAndParams("Agent", &model.Params{})
	user1 := reflect.ValueOf(users).Elem().Index(0)
	//agent1 := reflect.ValueOf(agents).Elem().Index(0)
	test := user1.FieldByName("Agents")
	test.Set(reflect.ValueOf(agents).Elem())
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

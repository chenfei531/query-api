package data

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"fmt"

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

/*
	init params for orm query with one level of nested field
*/
func getSelectParams(data interface{}, selectParams []string) (*SelectParams, error) {
	dataType := reflect.TypeOf(data)
	//assuming primary key must be ID (no customized primary key name, no composite primary key)
	//TODO: get primary Key from metadata
	//TODO: add field validation
	pk := "ID"
	containPK := false
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
				if f.Name == pk {
					containPK = true
				}
			}
		}
	}
	if !containPK {
		fields = append(fields, pk)
	}

	return &SelectParams{primaryKey: pk, containPK: containPK, field: fields, nestedField: nestedFields}, nil
}

//Select does not work for Preloaded object
//https://github.com/go-gorm/gorm/issues/4015
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
		fmt.Printf("%s\n", p.Sort)
	}
	if len(p.Select) > 0 {
		tx.Select(strings.Join(p.Select, ", "))
	}
	error := tx.Find(result).Error
	return result, error
}

/*
	query required field of both parent and child
*/
func (dm *SqliteDataManager) GetDataByNameAndParams(name string, p *model.Params) (string, error) {
	data, error := model.GetObjectByName(name)
	if nil != error {
		return "", error
	}
	sp, error := getSelectParams(data, p.Select)
	if nil != error {
		return "", error
	}
	//get parents
	p.Select = sp.field
	results, error := dm.getDataByParams(data, p)
	if nil != error {
		return "", error
	}
	//return if no children field are required
	if len(sp.nestedField) == 0 {
		b, _ := json.MarshalIndent(results, "", "    ")
		return string(b), nil
	}

	//get parent ids
	var idList []uint64
	idIndexMap := make(map[uint64]uint)
	resultsReflect := reflect.ValueOf(results).Elem()
	for i := 0; i < resultsReflect.Len(); i++ {
		pkField := resultsReflect.Index(i).FieldByName("ID")
		id := pkField.Uint()
		idList = append(idList, id)
		idIndexMap[id] = uint(i)
		//remove primary key of parent if necessary
		if !sp.containPK {
			pkField.SetUint(0)
		}
	}

	//using parents id to query children
	for childName, selectList := range sp.nestedField {
		fkFieldName := name + sp.primaryKey
		fkDBName := name + "_" + sp.primaryKey
		childParams := model.Params{Select: sp.field, FilterExp: fkDBName + " IN (?)", FilterArgs: make([]interface{}, 0, 8)}
		for i := 0; i < len(idList); i++ {
			childParams.FilterArgs = append(childParams.FilterArgs, idList[uint(i)])
		}

		selectList = append(selectList, fkDBName)

		childParams.Select = selectList

		data, error := model.GetObjectByName(childName)
		if nil != error {
			return "", nil
		}
		children, error := dm.getDataByParams(data, &childParams)
		if nil != error {
			return "", nil
		}
		//insert children into parents
		childrenReflect := reflect.ValueOf(children).Elem()

		for i := 0; i < childrenReflect.Len(); i++ {
			childFieldName := childName + "s"
			foreignId := childrenReflect.Index(i).FieldByName(fkFieldName).Uint()
			parentIndex := int(idIndexMap[foreignId])
			parent := reflect.ValueOf(results).Elem().Index(parentIndex)
			childrenField := parent.FieldByName(childFieldName)
			buf := childrenField
			//remove unwanted fields
			childReflect := reflect.ValueOf(children).Elem().Index(i)
			childReflect.FieldByName(fkFieldName).SetUint(0)

			//append children to parents
			buf = reflect.Append(buf, childReflect)
			childrenField.Set(buf)
		}
	}
	b, _ := json.MarshalIndent(results, "", "    ")
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

package query

import (
	"encoding/json"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type GormDataReader struct {
	db *gorm.DB
}

func NewGormDataReader(db *gorm.DB) *GormDataReader {
	return &GormDataReader{db}
}

func (dm *GormDataReader) getDataByParams(data interface{}, p *Params) (interface{}, error) {
	dataType := reflect.TypeOf(data)
	result := reflect.New(reflect.SliceOf(dataType)).Interface()
	tx := dm.db.Model(data)
	if len(p.FilterExp) > 0 {
		tx.Where(p.FilterExp, p.FilterArgs...)
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

func (dm *GormDataReader) getNestedData(node *Node) (interface{}, error) {
	name := node.Name
	params := node.Params
	data, error := GetObjectByName(name)
	if nil != error {
		return "", error
	}
	results, error := dm.getDataByParams(data, &params)
	if nil != error {
		return "", error
	}

	var idList []uint64
	idIndexMap := make(map[uint64]uint)
	resultsReflect := reflect.ValueOf(results).Elem()
	for i := 0; i < resultsReflect.Len(); i++ {
		pkField := resultsReflect.Index(i).FieldByName("ID")
		id := pkField.Uint()
		idList = append(idList, id)
		idIndexMap[id] = uint(i)
		//remove primary key of parent if necessary
		if !node.ContainPK {
			pkField.SetUint(0)
		}
	}

	for _, child := range node.Children {
		fkFieldName := name + "ID"
		fkDBName := name + "_" + "ID"
		childParamsPtr := &(child.Params)
		childParamsPtr.FilterExp = fkDBName + " IN (?)"
		fkArgs := make([]interface{}, 0, 8)
		for i := 0; i < len(idList); i++ {
			fkArgs = append(fkArgs, idList[uint(i)])
		}
		childParamsPtr.FilterArgs = make([]interface{}, 0, 2)
		childParamsPtr.FilterArgs = append(childParamsPtr.FilterArgs, fkArgs)

		childParamsPtr.Select = append(childParamsPtr.Select, fkDBName)

		//b, _ := json.MarshalIndent(child, "", "    ")
		//fmt.Printf("tree:%s\n", b)
		children, _ := dm.getNestedData(child)
		childrenReflect := reflect.ValueOf(children).Elem()

		for i := 0; i < childrenReflect.Len(); i++ {
			childFieldName := child.Name + "s"
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
	return results, nil
}

func (dm *GormDataReader) GetData(root *Node) (string, error) {
	results, error := dm.getNestedData(root)
	if nil != error {
		return "", error
	}
	b, error := json.MarshalIndent(results, "", "    ")
	if nil != error {
		return "", error
	}
	return string(b), nil
}

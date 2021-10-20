package query

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type QueryWorker struct {
	queryParser    *QueryParser
	gormDataReader *GormDataReader
}

func NewQueryWorker() *QueryWorker {
	//call pkg/database/gorm.New to get db
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	gormDataReader := NewGormDataReader(db)
	queryParser := NewQueryParser()
	return &QueryWorker{queryParser: queryParser, gormDataReader: gormDataReader}
}

func (qw *QueryWorker) Query(resourceName string, queryStr string) (string, error) {
	root, error := qw.queryParser.GetQueryTree(resourceName, queryStr)
	if error != nil {
		return "", error
	}
	result, error := qw.gormDataReader.GetData(root)
	if error != nil {
		return "", error
	}
	return result, nil
}

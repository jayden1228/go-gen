package service

import (
	"go-gen/internal/pkg/database"
	"sort"
	"strings"
)

// TableNameAndComment
type TableNameAndComment struct {
	Index   int    // 索引
	Name    string // 表名
	Comment string // 注释
}

// FindDbTables 获取数据库所有表
func FindDbTables() ([]*TableNameAndComment, error) {
	// 获取表名和注释
	var nameAndComments []*TableNameAndComment
	dbName := database.GetDbName()
	if err := database.GetDB().Table("tables").Select("table_name", "table_comment").Where("table_schema = ?", dbName).Find(&nameAndComments).Error; err != nil {
		return nil, err
	}

	// 添加索引
	for idx, info := range nameAndComments {
		idx++
		info.Index = idx
	}
	//排序, 采用升序
	sort.Slice(nameAndComments, func(i, j int) bool {
		return strings.ToLower(nameAndComments[i].Name) < strings.ToLower(nameAndComments[j].Name)
	})
	return nameAndComments, nil
}

package service

import (
	"fmt"
	"go-gen/internal/constant"
	"go-gen/internal/pkg/tool"
	"os"
	"path"
)

// EntityReq 构建实体参数
type EntityReq struct {
	Index        int    // 序列号
	TableName    string // 表名称
	TableComment string // 表注释
	Path         string // 文件路径
	EntityPath   string //实体路径
	Pkg          string // 命名空间名称
	EntityPkg    string // entity实体的空间名称
	FormatList   []string
	TableDesc    []*TableDesc // 表详情
}

// TableInfo
type TableInfo struct {
	Table            string        // 表名
	NullTable        string        // 空表名称
	TableComment     string        // 表注释
	TableCommentNull string        // 表注释
	Fields           []*FieldsInfo // 表字段
}

// FieldsInfo
type FieldsInfo struct {
	Name         string
	Type         string
	NullType     string
	DbName       string
	DbOriField   string // 数据库的原生字段名称
	FormatFields string
	Remark       string
}

func BuildEntityReq() {
	//fileName := st
	//entityPath := path.Join(os.Getwd())
}

func GenerateDBEntity(req *EntityReq) error {
	// 构建目标文件路径
	pwd, _ := os.Getwd()
	filePath := path.Join(pwd, "yyy", "damn.go")
	// 填充文件头部
	header := fmt.Sprintf(constant.EntityHeader, req.EntityPkg)
	// 写入文件
	_, err := tool.WriteFile(filePath, header)
	return err
}

package service

import (
	"bytes"
	"errors"
	"fmt"
	. "go-gen/internal/constant"
	"go-gen/internal/pkg/tool"
	"log"
	"path"
	"strings"

	"html/template"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
)

// EntityReq 构建实体参数
type EntityReq struct {
	Index         int                // 序列号
	TableName     string             // 表名称
	TableComment  string             // 表注释
	Path          string             // 文件路径
	EntityPath    string             //实体路径
	Pkg           string             // 命名空间名称
	EntityPkg     string             // entity实体的空间名称
	FormatList    []string           // 标签列表 【"json", "gorm"】
	TableDesc     []*CustomTableDesc // 表详情
	ProjectModule string             // 目标项目的go module名称

}

// TableInfo
type TableInfo struct {
	ProjectModule string        // 目标项目名称
	PackageName   string        // 包名
	TableName     string        // 表名
	StructName    string        // 结构名
	TableComment  string        // 表注释
	Fields        []*FieldsInfo // 表字段
}

// FieldsInfo
type FieldsInfo struct {
	Name         string
	Type         string
	DbName       string
	DbOriField   string // 数据库的原生字段名称
	FormatFields string
	Remark       string
	PrimaryKey   bool
}

func BuildReq(formatList []string, dir string, projectModule string) ([]*EntityReq, error) {
	tables, err := FindDbTables()
	if err != nil {
		return nil, err
	}
	var reqs = make([]*EntityReq, 0)
	for idx, table := range tables {
		idx++
		tableDesc, err := GenCustomTableDesc(table.TableName)
		if err != nil {
			return nil, err
		}

		fileName := fmt.Sprintf("%s.go", table.TableName)

		req := new(EntityReq)
		req.Index = idx
		req.EntityPath = path.Join(dir, EntityDir, fileName)
		req.TableName = table.TableName
		req.TableComment = table.TableComment
		req.TableDesc = tableDesc
		req.FormatList = formatList
		//req.EntityPkg = table.TableName
		req.EntityPkg = EntityPackageName
		req.ProjectModule = projectModule
		reqs = append(reqs, req)
	}

	return reqs, nil
}

func CreateDBEntity(formatList []string, dir string, projectModule string) error {
	reqs, err := BuildReq(formatList, dir, projectModule)
	if err != nil {
		return err
	}

	if err := GenerateDbErr(dir); err != nil {
		return err
	}

	for _, req := range reqs {
		err = GenerateDBEntity(req)
		if err != nil {
			log.Fatal("CreateEntityErr>>", err)
			return err
		}
	}
	return nil
}

func GenerateDBEntity(req *EntityReq) error {
	// 声明表结构变量
	TableData := new(TableInfo)
	TableData.ProjectModule = req.ProjectModule
	TableData.PackageName = req.EntityPkg
	TableData.StructName = gstr.CamelCase(req.TableName)
	TableData.TableName = req.TableName
	TableData.TableComment = tool.AddToComment(req.TableComment, "")

	check := fmt.Sprintf("type %s struct", gstr.ToUpper(req.TableName))
	if tool.CheckFileContainsChar(req.EntityPath, check) {
		return errors.New("it already exists. Please delete it and regenerate it")
	}
	// 加载模板文件
	tplByte, err := tool.Asset(TPL_CURD)
	if err != nil {
		return err
	}
	tpl, err := template.New("entity").Parse(string(tplByte))
	if err != nil {
		return err
	}
	// 装载表字段信息
	for _, val := range req.TableDesc {
		TableData.Fields = append(TableData.Fields, &FieldsInfo{
			Name:         gstr.CamelCase(val.ColumnName),
			Type:         val.GolangType,
			DbOriField:   val.ColumnName,
			FormatFields: tool.FormatField(val.ColumnName, req.FormatList),
			Remark:       tool.AddToComment(val.ColumnComment, ""),
			PrimaryKey:   val.PrimaryKey,
		})
	}
	content := bytes.NewBuffer([]byte{})

	if err := tpl.Execute(content, TableData); err != nil {
		return err
	}
	// 表信息写入文件
	con := strings.Replace(content.String(), "&#34;", `"`, -1)
	err = gfile.PutContents(req.EntityPath, con)
	if err != nil {
		return err
	}

	tool.Gofmt(req.EntityPath)
	return nil
}

func GenerateDbErr(dir string) error {
	type tableInfo struct {
		PackageName string
	}

	errPath := path.Join(dir, EntityDir, DbErrFileName)
	if gfile.Exists(errPath) {
		return nil
	}

	// 加载模板文件
	tplByte, err := tool.Asset(TPL_EROR)
	if err != nil {
		return err
	}
	tpl, err := template.New("error").Parse(string(tplByte))
	if err != nil {
		return err
	}

	content := bytes.NewBuffer([]byte{})

	if err := tpl.Execute(content, tableInfo{PackageName: EntityPackageName}); err != nil {
		return err
	}

	err = gfile.PutContents(errPath, content.String())
	if err != nil {
		return err
	}

	tool.Gofmt(errPath)
	return nil
}

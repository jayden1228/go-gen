package service

import (
	"bytes"
	"errors"
	"fmt"
	. "go-gen/internal/constant"
	"go-gen/internal/pkg/tool"
	"log"
	"os"
	"path"
	"strings"

	"html/template"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
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
	TableDesc    []*CustomTableDesc // 表详情
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

func CreateDBEntity(formatList []string) error {
	pwd, _ := os.Getwd()
	dir := path.Join(pwd, "yyy")
	tables, err := FindDbTables()
	if err != nil {
		return err
	}
	for idx, table := range tables {
		idx++
		tableDesc, err := GenCustomTableDesc(table.TableName)
		if err != nil {
			return err
		}
		req := new(EntityReq)
		req.Index = idx
		req.EntityPath = path.Join(dir, table.TableName+".go")
		req.TableName = table.TableName
		req.TableComment = table.TableComment
		req.TableDesc = tableDesc
		req.FormatList = formatList
		req.EntityPkg = table.TableName
		err = GenerateDBEntity(req)
		if err != nil {
			log.Fatal("CreateEntityErr>>", err)
			return err
		}
	}
	return nil
}

func GenerateDBEntity(req *EntityReq) error {
	// 填充文件头部
	header := fmt.Sprintf(EntityHeader, req.EntityPkg)
	check := "github.com/go-sql-driver/mysql"
	// 写入import等头部信息
	if !tool.CheckFileContainsChar(req.EntityPath, check) {
		_ = gfile.PutContents(req.EntityPath, header)
	}

	// 声明表结构变量
	TableData := new(TableInfo)
	TableData.Table = strings.ToUpper(req.TableName)
	TableData.NullTable = DbNullPrefix + TableData.Table
	TableData.TableComment = tool.AddToComment(req.TableComment, "")

	check = fmt.Sprintf("type %s struct", strings.ToUpper(req.TableName))
	if tool.CheckFileContainsChar(req.EntityPath, check) {
		return errors.New("it already exists. Please delete it and regenerate it")
	}
	// 加载模板文件
	tplByte, err := tool.Asset(TPL_ENTITY)
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
			NullType:     val.MysqlNullType,
			DbOriField:   val.ColumnName,
			FormatFields: tool.FormatField(val.ColumnName, req.FormatList),
			Remark:       tool.AddToComment(val.ColumnComment, ""),
		})
	}
	content := bytes.NewBuffer([]byte{})

	if err := tpl.Execute(content, TableData); err != nil {
		return err
	}
	// 表信息写入文件
	con := strings.Replace(content.String(), "&#34;", `"`, -1)
	err = gfile.PutContentsAppend(req.EntityPath, con)
	if err != nil {
		return err
	}

	tool.Gofmt(req.EntityPath)
	return nil
}

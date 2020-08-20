package {{.PackageName}}
import (
	gdb "{{.ProjectModule}}/pkg/database/mysql"

	"github.com/jinzhu/copier"
)

{{.TableComment}}
type {{.StructName}} struct {
{{range $j, $item := .Fields}}{{$item.Name}}       {{$item.Type}}    {{$item.FormatFields}}        {{$item.Remark}}
{{end}}
}

// Add{{.StructName}} is a function to add a single record to {{.TableName}} table
// error - ErrInsertFailed, db save call failed
func Add{{.StructName}}(record *{{.StructName}}) (result *{{.StructName}}, RowsAffected int64, err error) {
    db := gdb.GetDB().Save(record)
	if err = db.Error; err != nil {
	    return nil, -1, ErrInsertFailed
	}
	return record, db.RowsAffected, nil
}

// Delete{{.StructName}} is a function to delete a single record from {{.TableName}} table
// error - ErrNotFound, db Find error
// error - ErrDeleteFailed, db Delete failed error
func Delete{{.StructName}}({{range $item := .Fields}} {{ if $item.PrimaryKey }} {{$item.Name}} {{$item.Type}},{{end}}{{end -}}) (rowsAffected int64, err error) {
    record := &{{.StructName}}{}
    db := gdb.GetDB().First(record, {{range $item := .Fields}} {{ if $item.PrimaryKey }} {{$item.Name -}},{{end}}{{end}})
    if db.Error != nil {
        return -1, ErrNotFound
    }

    db := gdb.GetDB().Delete(record)
    if err = db.Error; err != nil {
        return -1, ErrDeleteFailed
    }

   return db.RowsAffected, nil
}

// Update{{.StructName}} is a function to update a single record from {{.TableName}} table
// error - ErrNotFound, db record for id not found
// error - ErrUpdateFailed, db meta data copy failed or db.Save call failed
func Update{{.StructName}}({{range $item := .Fields}} {{ if $item.PrimaryKey }} {{$item.Name}} {{$item.Type}},{{end}}{{end -}}updated *{{.StructName}}) (result *{{.StructName}}, RowsAffected int64, err error) {
   result = &{{.StructName}}{}
   db := gdb.GetDB().First(result,{{range $item := .Fields}} {{ if $item.PrimaryKey }} {{$item.Name}},{{end}}{{end -}})
   if err = db.Error; err != nil {
      return nil, -1, ErrNotFound
   }

   if err = copier.Copy(result, updated); err != nil {
      return nil, -1, ErrUpdateFailed
   }

   db := gdb.GetDB().Save(result)
   if err = db.Error; err != nil  {
      return nil, -1, ErrUpdateFailed
   }

   return result, db.RowsAffected, nil
}

// Get{{.StructName}} is a function to get a single record from the {{.TableName}} table
// error - ErrNotFound, db Find error
func Get{{.StructName}}({{range $item := .Fields}} {{ if $item.PrimaryKey }} {{$item.Name}} {{$item.Type}},{{end}}{{end -}}) (record *{{.StructName}}, err error) {
	record = &{{.StructName}}{}
	if err = gdb.GetDB().First(record, {{range $item := .Fields}} {{ if $item.PrimaryKey }} {{$item.Name}},{{end}}{{end -}}).Error; err != nil {
	    err = ErrNotFound
		return record, err
	}

	return record, nil
}

// GetAll{{.StructName}} is a function to get a slice of record(s) from {{.TableName}} table
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// params - order    - db sort order column
// error - ErrNotFound, db Find error
func GetAll{{.StructName}}(page, pagesize int64, order string) (results []*{{.StructName}}, totalRows int, err error) {

	resultOrm := gdb.GetDB().Model(&{{.StructName}}{})
    resultOrm.Count(&totalRows)

	if page > 0 {
		offset := (page - 1) * pagesize
		resultOrm = resultOrm.Offset(offset).Limit(pagesize)
	} else {
		resultOrm = resultOrm.Limit(pagesize)
    }

	if order != "" {
		resultOrm = resultOrm.Order(order)
	}

	if err = resultOrm.Find(&results).Error; err != nil {
	    err = ErrNotFound
		return nil, -1, err
	}

	return results, totalRows, nil
}
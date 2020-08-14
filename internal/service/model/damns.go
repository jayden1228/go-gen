package model

import (
	gdb "go-gen/internal/pkg/database"

	"github.com/jinzhu/copier"
)

// 关于压测数据结构表
type Damns struct {
	Id      int32  `json:"id" gorm:"id"`
	Title   string `json:"title" gorm:"title"`     // 标题
	Tag     string `json:"tag" gorm:"tag"`         // 标签
	Content string `json:"content" gorm:"content"` //  内容
	Index   int32  `json:"index" gorm:"index"`     // 索引

}

// AddDamns is a function to add a single record to damns table
// error - ErrInsertFailed, db save call failed
func AddDamns(record *Damns) (result *Damns, RowsAffected int64, err error) {
	db := gdb.GetDB().Save(record)
	if err = db.Error; err != nil {
		return nil, -1, ErrInsertFailed
	}
	return record, db.RowsAffected, nil
}

// DeleteDamns is a function to delete a single record from damns table
// error - ErrNotFound, db Find error
// error - ErrDeleteFailed, db Delete failed error
func DeleteDamns(Id int32) (rowsAffected int64, err error) {
	record := &Damns{}
	db := gdb.GetDB().First(record, Id)
	if db.Error != nil {
		return -1, ErrNotFound
	}

	db = gdb.GetDB().Delete(record)
	if err = db.Error; err != nil {
		return -1, ErrDeleteFailed
	}

	return db.RowsAffected, nil
}

// UpdateDamns is a function to update a single record from damns table
// error - ErrNotFound, db record for id not found
// error - ErrUpdateFailed, db meta data copy failed or db.Save call failed
func UpdateDamns(Id int32, updated *Damns) (result *Damns, RowsAffected int64, err error) {
	result = &Damns{}
	db := gdb.GetDB().First(result, Id)
	if err = db.Error; err != nil {
		return nil, -1, ErrNotFound
	}

	if err = copier.Copy(result, updated); err != nil {
		return nil, -1, ErrUpdateFailed
	}

	db = gdb.GetDB().Save(result)
	if err = db.Error; err != nil {
		return nil, -1, ErrUpdateFailed
	}

	return result, db.RowsAffected, nil
}

// GetDamns is a function to get a single record from the damns table
// error - ErrNotFound, db Find error
func GetDamns(Id int32) (record *Damns, err error) {
	record = &Damns{}
	if err = gdb.GetDB().First(record, Id).Error; err != nil {
		err = ErrNotFound
		return record, err
	}

	return record, nil
}

// GetAllDamns is a function to get a slice of record(s) from damns table
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// params - order    - db sort order column
// error - ErrNotFound, db Find error
func GetAllDamns(page, pagesize int64, order string) (results []*Damns, totalRows int, err error) {

	resultOrm := gdb.GetDB().Model(&Damns{})
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

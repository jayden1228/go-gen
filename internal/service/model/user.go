package model

import (
	gdb "go-gen/internal/pkg/database"

	"github.com/jinzhu/copier"
)

type User struct {
	Id   int32  `json:"id" gorm:"id"`
	Name string `json:"name" gorm:"name"` // 姓名
	Age  int32  `json:"age" gorm:"age"`   // 年龄

}

// AddUser is a function to add a single record to user table
// error - ErrInsertFailed, db save call failed
func AddUser(record *User) (result *User, RowsAffected int64, err error) {
	db := gdb.GetDB().Save(record)
	if err = db.Error; err != nil {
		return nil, -1, ErrInsertFailed
	}
	return record, db.RowsAffected, nil
}

// DeleteUser is a function to delete a single record from user table
// error - ErrNotFound, db Find error
// error - ErrDeleteFailed, db Delete failed error
func DeleteUser(Id int32) (rowsAffected int64, err error) {
	record := &User{}
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

// UpdateUser is a function to update a single record from user table
// error - ErrNotFound, db record for id not found
// error - ErrUpdateFailed, db meta data copy failed or db.Save call failed
func UpdateUser(Id int32, updated *User) (result *User, RowsAffected int64, err error) {
	result = &User{}
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

// GetUser is a function to get a single record from the user table
// error - ErrNotFound, db Find error
func GetUser(Id int32) (record *User, err error) {
	record = &User{}
	if err = gdb.GetDB().First(record, Id).Error; err != nil {
		err = ErrNotFound
		return record, err
	}

	return record, nil
}

// GetAllUser is a function to get a slice of record(s) from user table
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// params - order    - db sort order column
// error - ErrNotFound, db Find error
func GetAllUser(page, pagesize int64, order string) (results []*User, totalRows int, err error) {

	resultOrm := gdb.GetDB().Model(&User{})
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

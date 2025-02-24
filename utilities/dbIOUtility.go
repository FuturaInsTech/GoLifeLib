// package utilities

// import (
// 	"errors"
// 	"time"

// 	"github.com/FuturaInsTech/GoLifeLib/initializers"
// 	"gorm.io/gorm"
// )

// type DBModel interface {
// 	Update(map[string]interface{}, map[string]interface{}) error
// 	Delete(map[string]interface{}) error
// 	First(map[string]interface{}) error
// 	Create(map[string]interface{}) error
// 	Exists(map[string]interface{}) (bool, error)
// 	Find(string, []interface{}, bool, string, int, int) (int64, error)
// 	Count(map[string]interface{}) (int64, error)
// }

// type Model struct {
// 	Db interface{}
// }

// func (model *Model) Update(fieldMap map[string]interface{}, whereCondition map[string]interface{}) error {

// 	var totalRowCount int64
// 	result := initializers.DB.Model(model.Db).Where(whereCondition).Count(&totalRowCount)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	if totalRowCount == 0 {

// 		return errors.New("record not found")

// 	}

// 	if totalRowCount > 1 {

// 		return errors.New("too many records matched")

// 	}

// 	fieldMap["updated_at"] = time.Now()

// 	result = initializers.DB.Model(model.Db).Where(whereCondition).Updates(fieldMap)

// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	return nil
// }

// func (model *Model) First(whereCondition map[string]interface{}) error {

// 	result := initializers.DB.Where(whereCondition).First(model.Db)

// 	recordNotFound := errors.Is(result.Error, gorm.ErrRecordNotFound)

// 	if recordNotFound {

// 		return errors.New("record not found")

// 	}

// 	if result.Error != nil {

// 		return result.Error

// 	}
// 	return nil
// }

// func (model *Model) Delete(whereCondition map[string]interface{}) error {

// 	var totalRowCount int64
// 	result := initializers.DB.Model(model.Db).Where(whereCondition).Count(&totalRowCount)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	if totalRowCount == 0 {

// 		return errors.New("record not found")

// 	}

// 	if totalRowCount > 1 {

// 		return errors.New("too many records matched")

// 	}

// 	result = initializers.DB.Where(whereCondition).Delete(model.Db)

// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	return nil
// }

// func (model *Model) Count(whereCondition map[string]interface{}) (int64, error) {

// 	var totalRowCount int64
// 	result := initializers.DB.Model(model.Db).Where(whereCondition).Count(&totalRowCount)
// 	if result.Error != nil {
// 		return 0, result.Error
// 	}

// 	return totalRowCount, nil
// }

// func (model *Model) Create(fieldMap map[string]interface{}) error {

// 	fieldMap["created_at"] = time.Now()

// 	result := initializers.DB.Model(model.Db).Create(fieldMap)

// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	return nil
// }

// func (model *Model) Exists(whereCondition map[string]interface{}) (bool, error) {

// 	result := initializers.DB.Where(whereCondition).First(model.Db)

// 	recordNotFound := errors.Is(result.Error, gorm.ErrRecordNotFound)

// 	if result.Error != nil && !recordNotFound {

// 		return false, result.Error

// 	}

// 	if recordNotFound {

// 		return false, nil

// 	}

// 	return true, nil
// }

// func (model *Model) Find(filteringquery string, filteringFields []interface{}, reportOnly bool, sortingList string, pageSize int, offset int) (int64, error) {

// 	var totalRowCount int64 = 0
// 	result := initializers.DB.Model(model.Db).Where(filteringquery, filteringFields...).Count(&totalRowCount)

// 	if result.Error == nil {
// 		//if for reporting purpose , pagination not required as the whole data will be downloaded in excel/pdf
// 		if reportOnly {
// 			result = initializers.DB.Order(sortingList).Where(filteringquery, filteringFields...).Find(model.Db)
// 		} else {

// 			result = initializers.DB.Order(sortingList).Limit(pageSize).Offset(offset).Where(filteringquery, filteringFields...).Find(model.Db)
// 		}

// 	}

// 	if result.Error != nil {

// 		return 0, result.Error

// 	}

// 	return totalRowCount, nil

// }

package utils

import "gorm.io/gorm"

func MustExist(tx *gorm.DB, model interface{}, field string, id interface{}) bool {
	var count int64 = 0
	err := tx.Model(model).Where(field+" = ?", id).Count(&count).Error
	if err != nil {
		return false
	}
	return count == 0
}

func FilterIfNotNil(destPtr interface{}, tx *gorm.DB, op func(query interface{}, args ...interface{}) (tx *gorm.DB), query interface{}, args ...interface{}) *gorm.DB {
	args = append(args, destPtr)
	switch destPtr.(type) {
	case *string:
		if destPtr != nil && destPtr.(*string) != nil {
			return op(query, args)
		}
	case *int:
		if destPtr != nil && destPtr.(*int) != nil {
			return op(query, args)
		}
	case *int8:
		if destPtr != nil && destPtr.(*int8) != nil {
			return op(query, args)
		}
	case *int16:
		if destPtr != nil && destPtr.(*int16) != nil {
			return op(query, args)
		}
	case *int32:
		if destPtr != nil && destPtr.(*int32) != nil {
			return op(query, args)
		}
	case *int64:
		if destPtr != nil && destPtr.(*int64) != nil {
			return op(query, args)
		}
	case *uint:
		if destPtr != nil && destPtr.(*uint) != nil {
			return op(query, args)
		}
	case *uint8:
		if destPtr != nil && destPtr.(*uint8) != nil {
			return op(query, args)
		}
	case *uint16:
		if destPtr != nil && destPtr.(*uint16) != nil {
			return op(query, args)
		}
	case *uint32:
		if destPtr != nil && destPtr.(*uint32) != nil {
			return op(query, args)
		}
	case *uint64:
		if destPtr != nil && destPtr.(*uint64) != nil {
			return op(query, args)
		}
	case *float32:
		if destPtr != nil && destPtr.(*float32) != nil {
			return op(query, args)
		}
	case *float64:
		if destPtr != nil && destPtr.(*float64) != nil {
			return op(query, args)
		}
	case *bool:
		if destPtr != nil && destPtr.(*bool) != nil {
			return op(query, args)
		}
	case *[]byte:
		if destPtr != nil && destPtr.(*[]byte) != nil {
			return op(query, args)
		}
	}
	return tx
}

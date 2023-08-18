package utilities

import (
	"errors"
	"time"

	"github.com/FuturaInsTech/GoLifeLib/initializers"
	"github.com/FuturaInsTech/GoLifeLib/models"
	"github.com/FuturaInsTech/GoLifeLib/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetVersionId(iCompany uint, lockedType types.LockedType, lockedTypeKey string) (string, error) {
	var tranLock models.TransactionLock
	result := initializers.DB.First(&tranLock, "company_id = ? and locked_type = ? and locked_type_key = ?", iCompany, lockedType, lockedTypeKey)
	if result.Error != nil {
		return "", result.Error
	}

	if !tranLock.IsValid {
		return "", errors.New("entity does not exist")
	}

	if tranLock.IsLocked {
		return "", errors.New("entity is locked")

	}
	return tranLock.VersionId, nil

}

func LockTheEntity(iCompany uint, lockedType types.LockedType, lockedTypeKey string, versionID string, iUserId uint64) error {

	var tranLock models.TransactionLock
	result := initializers.DB.First(&tranLock, "company_id = ? and locked_type = ? and locked_type_key = ?", iCompany, lockedType, lockedTypeKey)

	recordNotFound := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if recordNotFound {
		return errors.New("entity does not exist")
	}

	if result.Error != nil {
		return result.Error
	}

	if !tranLock.IsValid {
		return errors.New("entity does not exist")
	}

	if tranLock.IsLocked {
		return errors.New("entity is already locked")

	}

	if versionID != tranLock.VersionId {
		return errors.New("entity version mismatch")

	}

	tranLock.IsLocked = true
	tranLock.UpdatedID = iUserId
	tranLock.UpdatedAt = time.Now()

	result = initializers.DB.Save(&tranLock)

	if result.Error != nil {
		return result.Error
	}

	return nil

}

func CreateTheEntity(iCompany uint, lockedType types.LockedType, lockedTypeKey string) error {

	var tranLock models.TransactionLock
	result := initializers.DB.First(&tranLock, "company_id = ? and locked_type = ? and locked_type_key = ?", iCompany, lockedType, lockedTypeKey)

	recordNotFound := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if !recordNotFound && result.Error != nil {
		return result.Error
	}

	if !recordNotFound {
		return errors.New("entity already exists")
	}

	tranLock.CompanyID = iCompany
	tranLock.LockedTypeKey = lockedTypeKey
	tranLock.LockedType = lockedType
	tranLock.IsLocked = false
	tranLock.IsValid = true
	tranLock.CreatedAt = time.Now()
	tranLock.VersionId = uuid.New().String()

	result = initializers.DB.Create(&tranLock)

	if result.Error != nil {
		return result.Error
	}

	return nil

}

func UnLockTheEntity(iCompany uint, lockedType types.LockedType, lockedTypeKey string, iUserId uint64, changeVersion bool) error {

	var tranLock models.TransactionLock
	result := initializers.DB.First(&tranLock, "company_id = ? and locked_type = ? and locked_type_key = ?", iCompany, lockedType, lockedTypeKey)
	recordNotFound := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if recordNotFound {
		return errors.New("entity does not exist")
	}

	if result.Error != nil {
		return result.Error
	}

	if !tranLock.IsValid {
		return errors.New("entity does not exist")
	}

	if !tranLock.IsLocked {
		return errors.New("entity is not locked")

	}

	tranLock.IsLocked = false
	tranLock.UpdatedID = iUserId
	tranLock.UpdatedAt = time.Now()
	if changeVersion {
		tranLock.VersionId = uuid.New().String()
	}

	result = initializers.DB.Save(&tranLock)

	if result.Error != nil {
		return result.Error
	}

	return nil

}

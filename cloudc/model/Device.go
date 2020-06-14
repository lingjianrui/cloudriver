package model

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"time"
)

type Device struct {
	ID           uint32    `gorm:"primary_key;auto_increment" json:"id"`
	DeviceType   string    `gorm:"size:255;not null" json:"device_type"`
	DeviceSerial string    `gorm:"size:255;not null;unique" json:"device_serial"`
	Status       string    `gorm:"size:255;null;" json:"status"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Connection   *websocket.Conn
}

func (d *Device) Save(db *gorm.DB, did string) (*Device, error) {

	var err error
	err = db.Debug().Model(Device{}).Where("device_serial = ?", did).Take(&d).Error
	if gorm.IsRecordNotFoundError(err) {
		var cerr error
		cerr = db.Debug().Create(&d).Error
		if cerr != nil {
			return &Device{}, err
		}
	} else {
		return &Device{}, err
	}
	return d, nil
}

func (d *Device) UpdateStatus(db *gorm.DB, s string) (*Device, error) {
	db.Debug().Model(Device{}).Update("status", s)
	return d, nil
}

// THE ONLY PERSON THAT NEED TO DO THIS IS THE ADMIN, SO I HAVE COMMENTED THE ROUTES, SO SOMEONE ELSE DONT VIEW THIS DETAILS.
func (d *Device) FindAllDevice(db *gorm.DB) (*[]Device, error) {
	var err error
	devices := []Device{}
	err = db.Debug().Model(&Device{}).Limit(100).Find(&devices).Error
	if err != nil {
		return &[]Device{}, err
	}
	return &devices, err
}

func (d *Device) FindDeviceByID(db *gorm.DB, did uint32) (*Device, error) {
	var err error
	err = db.Debug().Model(Device{}).Where("id = ?", did).Take(&d).Error
	if err != nil {
		return &Device{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Device{}, errors.New("Device Not Found")
	}
	return d, err
}

package server

import (
	"go.rock.com/rock-platform/rock/server/database"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/log"
)

type InitRecord struct {
	QueryRecord  interface{}
	InsertRecord interface{}
}

// define admin account
func GetUsersInitData() []InitRecord {
	users := []InitRecord{
		InitRecord{
			QueryRecord: &models.User{
				Name: "admin",
			},
			InsertRecord: &models.User{
				Id:       1,
				Name:     "admin",
				Password: "3207ead4e092de77e022394b3204d755",
				Email:    "1031653788@qq.com",
				Salt:     "r8slTCTTHD8qVaYr",
				Token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo3LCJ1c2VybmFtZSI6ImdhbGF4aWFzX2FwaV90b2tlbiIsImRyb25lX3Rva2VuIjoiIiwicGFzc3dvcmQiOiJmMGE1NzI1NWYxNTEzMGQ3NzYwNjM5YzU4ZGU4MDc2MyIsInJvbGUiOiJhZG1pbiIsImV4cCI6NDc3MTc5NjY3MiwiaWF0IjoxNjE4MTk2NjcyLCJpc3MiOiJSb2NrIFdhbmciLCJzdWIiOiJMb2dpbiB0b2tlbiJ9.3Y00pHFvpD90yv5QmixuurNbIlJa-Dz1a1bRhO1Rp3Q", // expire time: 2121-03-19 11:04:32
				RoleId:   1,
			},
		},
	}
	return users
}

// define init role status
func GetRolesInitData() []InitRecord {
	roles := []InitRecord{
		InitRecord{
			QueryRecord: &models.Role{
				Name: "admin",
			},
			InsertRecord: &models.Role{
				Id:          1,
				Name:        "admin",
				Description: "administrator role of system",
			},
		},
		InitRecord{
			QueryRecord: &models.Role{
				Name: "system_tools_admin",
			},
			InsertRecord: &models.Role{
				Id:          2,
				Name:        "system_tools_admin",
				Description: "system tools' administrator",
			},
		},
		InitRecord{
			QueryRecord: &models.Role{
				Name: "developer",
			},
			InsertRecord: &models.Role{
				Id:          3,
				Name:        "developer",
				Description: "developer role",
			},
		},
		InitRecord{
			QueryRecord: &models.Role{
				Name: "deployer",
			},
			InsertRecord: &models.Role{
				Id:          4,
				Name:        "deployer",
				Description: "deployer role",
			},
		},
	}
	return roles
}

// if not exist then create admin account and role
func existOrInsert(e *database.DBEngine, records []InitRecord) {
	logger := log.GetLogger()
	for _, record := range records {
		var model interface{}
		switch record.QueryRecord.(type) {
		case *models.User:
			model = &models.User{}
		case *models.Role:
			model = &models.Role{}
		default:
			logger.Errorf("[DB INIT] %v model is a wrong type", record.QueryRecord)
		}

		err := e.Where(record.QueryRecord).First(model).Error
		if err != nil {
			if err.Error() == "record not found" {
				if err := e.Create(record.InsertRecord).Error; err != nil {
					logger.Errorf("[DB INIT] Error occurred when init %v , err : ", record.InsertRecord, err)
					continue
				}
				logger.Infof("[DB INIT] Init record %v successfully", record.InsertRecord)
			} else {
				logger.Errorf("[DB INIT] Error occurred when check %v existence", record.QueryRecord)
			}
		} else {
			logger.Infof("[DB INIT] User %v is already exist, skip init it", record.QueryRecord)
		}

	}
}

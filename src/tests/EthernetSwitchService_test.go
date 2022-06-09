package tests

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path"
	"rol/app/interfaces"
	"rol/app/services"
	"rol/domain"
	"rol/dtos"
	"rol/infrastructure"
	"runtime"
	"testing"
)

var testerSwitchService *GenericServiceTest[dtos.EthernetSwitchDto, dtos.EthernetSwitchCreateDto, dtos.EthernetSwitchUpdateDto, domain.EthernetSwitch]

func Test_EthernetSwitchService_Prepare(t *testing.T) {
	dbFileName := "ethernetSwitchService_test.db"
	dbConnection := sqlite.Open(dbFileName)
	testGenDb, err := gorm.Open(dbConnection, &gorm.Config{})
	if err != nil {
		t.Errorf("creating db failed: %v", err)
	}
	err = testGenDb.AutoMigrate(
		new(domain.EthernetSwitch),
	)
	if err != nil {
		t.Errorf("migration failed: %v", err)
	}

	logger := logrus.New()
	var repo interfaces.IGenericRepository[domain.EthernetSwitch]
	repo = infrastructure.NewGormGenericRepository[domain.EthernetSwitch](testGenDb, logger)
	var service interfaces.IGenericService[dtos.EthernetSwitchDto, dtos.EthernetSwitchCreateDto, dtos.EthernetSwitchUpdateDto, domain.EthernetSwitch]
	service, err = services.NewEthernetSwitchService(repo, logger)
	if err != nil {
		t.Errorf("create new service failed:  %q", err)
	}
	testerSwitchService = NewGenericServiceTest(service, repo, dbFileName)

	_, filename, _, _ := runtime.Caller(1)
	if _, err := os.Stat(path.Join(path.Dir(filename), dbFileName)); errors.Is(err, os.ErrNotExist) {
		return
	}
	err = os.Remove(dbFileName)
	if err != nil {
		t.Errorf("remove db failed:  %q", err)
	}
}

func Test_EthernetSwitchService_Create(t *testing.T) {
	createDto := dtos.EthernetSwitchCreateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
			Name:        "AutoTesting",
			Serial:      "test_serial",
			SwitchModel: 0,
			Address:     "123.123.123.123",
			Username:    "AutoUser",
		},
		//  pragma: allowlist nextline secret
		Password: "AutoPass",
	}
	err := testerSwitchService.GenericServiceCreate(createDto)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchService_GetByID(t *testing.T) {
	err := testerSwitchService.GenericServiceGetByID(testerSwitchService.InsertedID)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchService_Update(t *testing.T) {
	updateDto := dtos.EthernetSwitchUpdateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
			Name:        "AutoTestingUpdated",
			Serial:      "101",
			SwitchModel: 0,
			Address:     "123.123.123.123",
			Username:    "Test",
		},
		//  pragma: allowlist nextline secret
		Password: "Test",
	}
	err := testerSwitchService.GenericServiceUpdate(updateDto, testerSwitchService.InsertedID)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchService_Delete(t *testing.T) {
	err := testerSwitchService.GenericServiceDelete(testerSwitchService.InsertedID)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchService_Create20(t *testing.T) {
	for i := 1; i <= 20; i++ {
		createDto := dtos.EthernetSwitchCreateDto{
			EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
				Name:        fmt.Sprintf("AutoTesting_%d", i),
				Serial:      "test_serial",
				SwitchModel: 0,
				Address:     "123.123.123.123",
				Username:    "AutoUser",
			},
			//  pragma: allowlist nextline secret
			Password: "AutoPass",
		}
		err := testerSwitchService.GenericServiceCreate(createDto)
		if err != nil {
			t.Error(err)
		}
	}
}

func Test_EthernetSwitchService_GetList(t *testing.T) {
	err := testerSwitchService.GenericServiceGetList(20, 1, 10)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchService_Search(t *testing.T) {
	err := testerSwitchService.GenericServiceSearch("AutoUser")
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchService_CloseConnectionAndRemoveDb(t *testing.T) {
	err := testerSwitchService.GenericServiceCloseConnectionAndRemoveDb()
	if err != nil {
		t.Error(err)
	}
}

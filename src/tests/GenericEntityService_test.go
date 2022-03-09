package testss

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path"
	"rol/app/interfaces/generic"
	"rol/app/services"
	"rol/dtos"
	"rol/infrastructure"
	"runtime"
	"testing"
)

var testServiceFileName string
var testServiceDbConnection gorm.Dialector
var testServiceRepo generic.IGenericEntityRepository
var testService generic.IGenericEntityService
var serviceTestCreatedId uint

func Test_GenericEntityService_Prepare(t *testing.T) {
	testServiceFileName = "service_test.db"
	testServiceDbConnection = sqlite.Open(testServiceFileName)
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	testServiceRepo, _ = infrastructure.NewGormGenericEntityRepository(testServiceDbConnection, logger)
	testService, _ = services.NewGenericEntityService(testServiceRepo, logger)
	serviceTestCreatedId = 0
	_, filename, _, _ := runtime.Caller(1)
	if _, err := os.Stat(path.Join(path.Dir(filename), testServiceFileName)); errors.Is(err, os.ErrNotExist) {
		return
	}
	err := os.Remove(testServiceFileName)
	if err != nil {
		t.Errorf("remove db failed:  %q", err)
	}
}

func Test_GenericEntityService_Create(t *testing.T) {
	dto := dtos.EthernetSwitchCreateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
			Name:        "TestSwitch",
			Serial:      "TestSerial",
			SwitchModel: 146,
			Address:     "TestAddress",
			Username:    "TestUsername",
		},
		Password: "TestPassword",
	}
	var err error

	serviceTestCreatedId, err = testService.Create(&dto)
	if err != nil {
		t.Errorf("got %q, wanted %q", err, "nil")
	}
	if serviceTestCreatedId == 0 {
		t.Errorf("got %q, wanted %q", serviceTestCreatedId, " > 0")
	}
	if dto.Name != "TestSwitch" {
		t.Errorf("got name %q, wanted %q", dto.Name, "TestSwitch")
	}
}

func Test_GenericEntityService_GetById(t *testing.T) {
	dto := dtos.EthernetSwitchDto{}

	err := testService.GetById(&dto, serviceTestCreatedId)
	if err != nil {
		t.Errorf("got %q, wanted %q", err, "nil")
	}
	if dto.Name != "TestSwitch" {
		t.Errorf("got name %q, wanted %q", dto.Name, "TestSwitch")
	}
}

func Test_GenericEntityService_Update(t *testing.T) {
	dto := dtos.EthernetSwitchUpdateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{Name: "TestEdit"},
	}
	err := testService.Update(&dto, serviceTestCreatedId)
	if err != nil {
		t.Errorf("got %q, wanted %s", err, "nil")
	}
	if dto.Name != "TestEdit" {
		t.Errorf("got name %q, wanted %q", dto.Name, "TestEdit")
	}
}

func Test_GenericEntityService_GetAll(t *testing.T) {
	dtosArr := &[]*dtos.EthernetSwitchDto{}
	err := testService.GetAll(dtosArr)
	if err != nil {
		t.Errorf("got %q, wanted %q", err, "nil")
	}
	if len(*dtosArr) != 1 {
		t.Errorf("got count %d, wanted %q", len(*dtosArr), 1)
	}
}

func Test_GenericEntityService_Create20(t *testing.T) {
	for i := 2; i < 22; i++ {
		dto := dtos.EthernetSwitchCreateDto{
			EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
				Name:        fmt.Sprintf("TestSwitch_%d", i),
				Serial:      "TestSerial",
				SwitchModel: 146,
				Address:     "TestAddress",
				Username:    "TestUsername",
			},
			Password: "TestPassword",
		}
		_, _ = testService.Create(&dto)
	}

	dtosArr := &[]*dtos.EthernetSwitchDto{}
	paginator, err := testService.GetList(dtosArr, "", "", "", 1, 5)
	items := paginator.Items.(*[]*dtos.EthernetSwitchDto)
	if err != nil {
		t.Errorf("got %q, wanted %s", err, "nil")
	}
	if paginator.ItemsCount != 21 {
		t.Errorf("got %q items, wanted %q", paginator.ItemsCount, 21)
	}
	if len(*items) != 5 {
		t.Errorf("got count %q, wanted %q", len(*items), 5)
	}
}

func Test_GenericEntityService_GetListOrderBy(t *testing.T) {
	dtosArr := &[]*dtos.EthernetSwitchDto{}
	paginator, err := testService.GetList(dtosArr, "", "name", "desc", 1, 10)
	itemsPtr := paginator.Items.(*[]*dtos.EthernetSwitchDto)
	items := (*itemsPtr)[:len(*itemsPtr)]

	if err != nil {
		t.Errorf("got %q, wanted %s", err, "nil")
	}

	if items[0].Name != "TestSwitch_9" {
		t.Errorf("got id %s, wanted %s", items[0].Name, "TestSwitch_9")
	}
}

func Test_GenericEntityService_GetListOrderDirection(t *testing.T) {
	dtosArr := &[]*dtos.EthernetSwitchDto{}
	paginator, err := testService.GetList(dtosArr, "", "id", "desc", 0, 0)
	itemsPtr := paginator.Items.(*[]*dtos.EthernetSwitchDto)
	items := (*itemsPtr)[:len(*itemsPtr)]

	if err != nil {
		t.Errorf("got %q, wanted %s", err, "nil")
	}

	if items[0].Id != 21 {
		t.Errorf("got id %d, wanted %q", items[0].Id, 1)
	}
}

func Test_GenericEntityService_GetListSearch(t *testing.T) {
	dtosArr := &[]*dtos.EthernetSwitchDto{}
	paginator, err := testService.GetList(dtosArr, "TestSwitch_21", "", "", 1, 1)
	itemsPtr := paginator.Items.(*[]*dtos.EthernetSwitchDto)
	items := (*itemsPtr)[:len(*itemsPtr)]

	if err != nil {
		t.Errorf("got %q, wanted %s", err, "nil")
	}

	if len(items) != 1 {
		t.Errorf("got %d items, wanted %d", len(items), 1)
	}

	if items[0].Id != 21 {
		t.Errorf("got id %d, wanted %q", items[0].Id, 1)
	}
}

func Test_GenericEntityService_GetListBadIntValues(t *testing.T) {
	dtosArr := &[]*dtos.EthernetSwitchDto{}
	paginator, err := testService.GetList(dtosArr, "", "", "", -1, -1)
	items := paginator.Items.(*[]*dtos.EthernetSwitchDto)
	if err != nil {
		t.Errorf("got %q, wanted nil", err)
	}
	if paginator.CurrentPage != 1 {
		t.Errorf("got page %d, wanted %q", paginator.CurrentPage, 1)
	}
	if paginator.PageSize != 10 {
		t.Errorf("got pagesize %d, wanted %q", paginator.PageSize, 10)
	}
	if len(*items) != 10 {
		t.Errorf("got count %d, wanted %q", len(*items), 10)
	}
}

func Test_GenericEntityService_Delete(t *testing.T) {
	dto := dtos.EthernetSwitchDto{}
	err := testService.Delete(&dto, serviceTestCreatedId)
	if err != nil {
		t.Errorf("got eror %q, wanted %q", err, "nil")
	}
}

func Test_GenericEntityService_GetAllAfterDelete(t *testing.T) {
	dtosArr := &[]*dtos.EthernetSwitchDto{}
	err := testService.GetAll(dtosArr)
	if err != nil {
		t.Errorf("got eror %q, wanted %q", err, "nil")
	}
	if len(*dtosArr) != 20 {
		t.Errorf("got count %d, wanted %d", len(*dtosArr), 20)
	}
}

func Test_GenericEntityService_CloseConnectionAndRemoveDb(t *testing.T) {
	sqlDb, err := testServiceRepo.(*infrastructure.GormGenericEntityRepository).Db.DB()
	if err != nil {
		t.Errorf("remove db failed:  %q", err)
	}
	sqlDb.Close()
	err = os.Remove(testServiceFileName)
	if err != nil {
		t.Errorf("remove db failed:  %q", err)
	}
}

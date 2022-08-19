package tests

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path"
	"rol/app/interfaces"
	"rol/domain"
	"rol/infrastructure"
	"runtime"
	"testing"
	"time"
)

var (
	testerVLANRepository *GenericRepositoryTest[domain.EthernetSwitchVLAN]
	uPortOne             uuid.UUID
	uPortTwo             uuid.UUID
)

func Test_EthernetSwitchVLANRepository_Prepare(t *testing.T) {
	dbFileName := "ethernetSwitchVLANRepo_test.db"
	dbConnection := sqlite.Open(dbFileName)
	testGenDb, err := gorm.Open(dbConnection, &gorm.Config{})
	if err != nil {
		t.Errorf("creating db failed: %v", err)
	}
	err = testGenDb.AutoMigrate(
		new(domain.EthernetSwitchVLAN),
	)
	if err != nil {
		t.Errorf("migration failed: %v", err)
	}

	logger := logrus.New()
	var repo interfaces.IGenericRepository[domain.EthernetSwitchVLAN]
	repo = infrastructure.NewGormGenericRepository[domain.EthernetSwitchVLAN](testGenDb, logger)

	testerVLANRepository = NewGenericRepositoryTest(repo, dbFileName)

	_, filename, _, _ := runtime.Caller(1)
	if _, err := os.Stat(path.Join(path.Dir(filename), dbFileName)); errors.Is(err, os.ErrNotExist) {
		return
	}
	err = os.Remove(dbFileName)
	if err != nil {
		t.Errorf("remove db failed:  %q", err)
	}
}

func Test_EthernetSwitchVLANRepository_Insert(t *testing.T) {
	uPortOne = uuid.New()
	uPortTwo = uuid.New()
	entity := domain.EthernetSwitchVLAN{
		VlanID:           2,
		EthernetSwitchID: uuid.New(),
		UntaggedPorts:    fmt.Sprintf("%s;%s", uPortOne.String(), uPortTwo.String()),
		TaggedPorts:      uPortOne.String(),
	}
	err := testerVLANRepository.GenericRepositoryInsert(entity)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchVLANRepository_GetByID(t *testing.T) {
	err := testerVLANRepository.GenericRepositoryGetByID(testerVLANRepository.InsertedID)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchVLANRepository_Update(t *testing.T) {
	entity := domain.EthernetSwitchVLAN{
		Entity:        domain.Entity{ID: testerVLANRepository.InsertedID},
		UntaggedPorts: uPortOne.String(),
		TaggedPorts:   fmt.Sprintf("%s;%s", uPortOne.String(), uPortTwo.String()),
	}
	err := testerVLANRepository.GenericRepositoryUpdate(entity)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchVLANRepository_GetList(t *testing.T) {
	err := testerVLANRepository.GenericRepositoryGetList()
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchVLANRepository_Delete(t *testing.T) {
	err := testerVLANRepository.GenericRepositoryDelete(testerVLANRepository.InsertedID)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchVLANRepository_Insert20(t *testing.T) {
	for i := 1; i <= 20; i++ {
		entity := domain.EthernetSwitchVLAN{
			VlanID:           i,
			EthernetSwitchID: uuid.New(),
			UntaggedPorts:    uPortOne.String(),
			TaggedPorts:      uPortTwo.String(),
		}
		err := testerVLANRepository.GenericRepositoryInsert(entity)
		if err != nil {
			t.Error(err)
		}
		time.Sleep(time.Second)
	}
}

func Test_EthernetSwitchVLANRepository_Pagination(t *testing.T) {
	err := testerVLANRepository.GenericRepositoryPagination(1, 10)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchVLANRepository_Filter(t *testing.T) {
	queryBuilder := testerVLANRepository.Repository.NewQueryBuilder(testerVLANRepository.Context)
	queryBuilder.Where("VlanId", ">", 7).Where("VlanId", "<", 9)
	err := testerVLANRepository.GenericRepositoryFilter(queryBuilder)
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchVLANRepository_Sort(t *testing.T) {
	err := testerVLANRepository.GenericRepositorySort()
	if err != nil {
		t.Error(err)
	}
}

func Test_EthernetSwitchVLANRepository_CloseConnectionAndRemoveDb(t *testing.T) {
	err := testerVLANRepository.GenericRepositoryCloseConnectionAndRemoveDb()
	if err != nil {
		t.Error(err)
	}
}

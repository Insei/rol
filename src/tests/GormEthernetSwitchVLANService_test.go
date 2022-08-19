package tests

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/app/services"
	"rol/domain"
	"rol/dtos"
	"rol/infrastructure"
	"runtime"
	"testing"
)

var (
	switchVLANService    interfaces.IGenericService[dtos.EthernetSwitchVLANDto, dtos.EthernetSwitchVLANCreateDto, dtos.EthernetSwitchVLANUpdateDto, domain.EthernetSwitchVLAN]
	switchVLANRepository interfaces.IGenericRepository[domain.EthernetSwitchVLAN]
	VLANEntityID         uuid.UUID
	ethernetSwitchIDVLAN uuid.UUID
)

func Test_EthernetSwitchVLANService_Prepare(t *testing.T) {
	dbFileName := "ethernetSwitchVLANService_test.db"
	//remove old test db file
	_, filename, _, _ := runtime.Caller(1)
	if _, err := os.Stat(path.Join(path.Dir(filename), dbFileName)); err == nil {
		err = os.Remove(dbFileName)
		if err != nil {
			t.Errorf("remove db failed:  %q", err)
		}
	}

	dbConnection := sqlite.Open(dbFileName)
	testGenDb, err := gorm.Open(dbConnection, &gorm.Config{})
	if err != nil {
		t.Errorf("creating db failed: %v", err)
	}
	err = testGenDb.AutoMigrate(
		new(domain.EthernetSwitchVLAN),
		new(domain.EthernetSwitch),
	)
	if err != nil {
		t.Errorf("migration failed: %v", err)
	}

	logger := logrus.New()

	switchVLANRepository = infrastructure.NewGormGenericRepository[domain.EthernetSwitchVLAN](testGenDb, logger)
	switchPortRepository = infrastructure.NewGormGenericRepository[domain.EthernetSwitchPort](testGenDb, logger)
	switchRepo := infrastructure.NewEthernetSwitchRepository(testGenDb, logger)
	switchVLANService, err = services.NewEthernetSwitchVLANService(switchVLANRepository, switchRepo, switchPortRepository, logger)
	if err != nil {
		t.Errorf("create new switch port service failed:  %q", err)
	}
	//create switch for testing
	switchEntity := domain.EthernetSwitch{
		Name:        "AutoTesting",
		Serial:      "AutoTesting",
		SwitchModel: "unifi_switch_us-24-250w",
		Address:     "192.111.111.111",
		Username:    "AutoTesting",
		//  pragma: allowlist nextline secret
		Password: "AutoTesting",
	}
	ethernetSwitchIDVLAN, err = switchRepo.Insert(context.TODO(), switchEntity)
	if err != nil {
		t.Errorf("create switch failed: %s", err)
	}
}

func Test_EthernetSwitchVLANService_CreateVLANWithoutSwitch(t *testing.T) {
	dto := dtos.EthernetSwitchVLANCreateDto{EthernetSwitchVLANBaseDto: dtos.EthernetSwitchVLANBaseDto{
		UntaggedPorts: []uuid.UUID{uuid.New()},
		TaggedPorts:   []uuid.UUID{uuid.New()},
	}, VlanID: 99}
	service := switchVLANService.(*services.EthernetSwitchVLANService)
	_, err := service.CreateVLAN(context.TODO(), uuid.New(), dto)
	if err == nil {
		t.Errorf("nil error, expected: %s", services.ErrorSwitchExistence)
	}
}

func Test_EthernetSwitchVLANService_CreateVLAN(t *testing.T) {
	dto := dtos.EthernetSwitchVLANCreateDto{EthernetSwitchVLANBaseDto: dtos.EthernetSwitchVLANBaseDto{
		UntaggedPorts: []uuid.UUID{uuid.New()},
		TaggedPorts:   []uuid.UUID{uuid.New()},
	}, VlanID: 2}
	service := switchVLANService.(*services.EthernetSwitchVLANService)
	var err error
	VLANEntityID, err = service.CreateVLAN(context.TODO(), ethernetSwitchIDVLAN, dto)
	if err != nil {
		t.Errorf("create VLAN failed: %s", err)
	}
}

func Test_EthernetSwitchVLANService_CreateFailByNonUniqueVLANID(t *testing.T) {
	dto := dtos.EthernetSwitchVLANCreateDto{EthernetSwitchVLANBaseDto: dtos.EthernetSwitchVLANBaseDto{
		UntaggedPorts: []uuid.UUID{uuid.New()},
		TaggedPorts:   []uuid.UUID{uuid.New()},
	}, VlanID: 2}
	service := switchVLANService.(*services.EthernetSwitchVLANService)
	_, err := service.CreateVLAN(context.TODO(), ethernetSwitchIDVLAN, dto)
	if err == nil {
		t.Errorf("nil error, expected: %s", errors.ValidationErrorMessage)
	}
}

func Test_EthernetSwitchVLANService_UpdateVLAN(t *testing.T) {
	updTaggedPort := uuid.New()
	updUntaggedPort := uuid.New()
	dto := dtos.EthernetSwitchVLANUpdateDto{EthernetSwitchVLANBaseDto: dtos.EthernetSwitchVLANBaseDto{
		UntaggedPorts: []uuid.UUID{updUntaggedPort},
		TaggedPorts:   []uuid.UUID{updTaggedPort},
	}}
	service := switchVLANService.(*services.EthernetSwitchVLANService)
	err := service.UpdateVLAN(context.TODO(), ethernetSwitchIDVLAN, VLANEntityID, dto)
	if err != nil {
		t.Errorf("update port failed: %s", err)
	}

	VLAN, err := switchVLANService.GetByID(context.TODO(), VLANEntityID)
	if err != nil {
		t.Errorf("failed to get VLAN: %s", err)
	}
	if VLAN.TaggedPorts[0] != updTaggedPort {
		t.Errorf("update port failed: unexpected uuid, got '%s', expect %s", VLAN.TaggedPorts[0].String(), updTaggedPort.String())
	}
}

func Test_EthernetSwitchVLANService_GetVLANs(t *testing.T) {
	service := switchVLANService.(*services.EthernetSwitchVLANService)
	for i := 3; i < 12; i++ {
		dto := dtos.EthernetSwitchVLANCreateDto{EthernetSwitchVLANBaseDto: dtos.EthernetSwitchVLANBaseDto{
			UntaggedPorts: []uuid.UUID{uuid.New()},
			TaggedPorts:   []uuid.UUID{uuid.New()},
		}, VlanID: i}
		_, err := service.CreateVLAN(context.TODO(), ethernetSwitchIDVLAN, dto)
		if err != nil {
			t.Errorf("create VLAN failed: %s", err)
		}
	}

	VLANs, err := service.GetVLANs(context.TODO(), ethernetSwitchIDVLAN, "", "", "", 1, 10)
	if err != nil {
		t.Errorf("get VLANs failed: %s", err)
	}
	if len(*VLANs.Items) != 10 {
		t.Errorf("get ports failed: wrong number of items, got %d, expect 10", len(*VLANs.Items))
	}
}

func Test_EthernetSwitchVLANService_GetVLANByID(t *testing.T) {
	service := switchVLANService.(*services.EthernetSwitchVLANService)
	VLAN, err := service.GetVLANByID(context.TODO(), ethernetSwitchIDVLAN, VLANEntityID)
	if err != nil {
		t.Errorf("get ports failed: %s", err)
	}
	if VLAN.VlanID != 2 {
		t.Errorf("get VLAN by ID failed: unexpected VlanID, got '%d', expect 2", VLAN.VlanID)
	}
}

func Test_EthernetSwitchVLANService_CloseConnectionAndRemoveDb(t *testing.T) {
	if err := switchVLANRepository.CloseDb(); err != nil {
		t.Errorf("close db failed:  %s", err)
	}
	if err := os.Remove("ethernetSwitchVLANService_test.db"); err != nil {
		t.Errorf("remove db failed:  %s", err)
	}
}

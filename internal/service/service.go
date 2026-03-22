package service

import (
	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
	"context"

	"github.com/jmoiron/sqlx"
)

type EmobileService interface {
	Create(ctx context.Context, name string, powerPointID string,
		batteryID string, charg_sysID string, chassisID string, bodyID string, electonicsID string) (string, error)
	Update(ctx context.Context, id string, name string, powerPointID string,
		batteryID string, charg_sysID string, chassisID string, bodyID string, electonicsID string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Emobile, error)
	GetByID(ctx context.Context, id string) (*domain.Emobile, error)
}

type BatteryService interface {
	Create(ctx context.Context, name string, batteryType string,
		batteryCapacity string, batteryInfo string) (string, error)
	Update(ctx context.Context, id string, name string, batteryType string, batteryCapacity string, batteryInfo string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Battery, error)
	GetByID(ctx context.Context, id string) (*domain.Battery, error)
}

type BodyService interface {
	Create(ctx context.Context, carcassID string, doorsID string, wingsID string) (string, error)
	Update(ctx context.Context, id string, carcassID string, doorsID string, wingsID string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Body, error)
	GetByID(ctx context.Context, id string) (*domain.Body, error)
}

type CarcassService interface {
	Create(ctx context.Context, carcassName string, carcassInfo string) (string, error)
	Update(ctx context.Context, id string, carcassName string, carcassInfo string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Carcass, error)
	GetByID(ctx context.Context, id string) (*domain.Carcass, error)
}

type DoorsService interface {
	Create(ctx context.Context, doorsName string, doorInfo string) (string, error)
	Update(ctx context.Context, id string, doorsName string, doorInfo string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Doors, error)
	GetByID(ctx context.Context, id string) (*domain.Doors, error)
}

type WingsService interface {
	Create(ctx context.Context, wingsName string, wingsInfo string) (string, error)
	Update(ctx context.Context, id string, wingsName string, wingsInfo string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Wings, error)
	GetByID(ctx context.Context, id string) (*domain.Wings, error)
}

type ChargerService interface {
	Create(ctx context.Context, chargerID string, connectorID string) (string, error)
	Update(ctx context.Context, id string, chargerID string, connectorID string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Charger, error)
	GetByID(ctx context.Context, id string) (*domain.Charger, error)
}

type ChassisService interface {
	Create(ctx context.Context, frameID, suspensionID, breakSystemID string) (string, error)
	Update(ctx context.Context, id string, frameID, suspensionID, breakSystemID string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Chassis, error)
	GetByID(ctx context.Context, id string) (*domain.Chassis, error)
}

type ElectronicsService interface {
	Create(ctx context.Context, controllerID, sensorID, wiringID string) (string, error)
	Update(ctx context.Context, id string, controllerID, sensorID, wiringID string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Electronics, error)
	GetByID(ctx context.Context, id string) (*domain.Electronics, error)
}

type ControllerService interface {
	Create(ctx context.Context, controllerName, controllerInfo string) (string, error)
	Update(ctx context.Context, id string, controllerName, controllerInfo string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Controller, error)
	GetByID(ctx context.Context, id string) (*domain.Controller, error)
}

type SensorService interface {
	Create(ctx context.Context, sensorName string, sensorInfo string) (string, error)
	Update(ctx context.Context, id string, sensorName string, sensorInfo string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Sensor, error)
	GetByID(ctx context.Context, id string) (*domain.Sensor, error)
}

type WiringService interface {
	Create(ctx context.Context, wiringName string, wiringInfo string) (string, error)
	Update(ctx context.Context, id string, wiringName string, wiringInfo string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Wiring, error)
	GetByID(ctx context.Context, id string) (*domain.Wiring, error)
}

type PowerPointService interface {
	Create(ctx context.Context, engineID, invertorID, gearboxID string) (string, error)
	Update(ctx context.Context, id string, engineID, invertorID, gearboxID string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.PowerPoint, error)
	GetByID(ctx context.Context, id string) (*domain.PowerPoint, error)
}

type Service struct {
	PowerPoint PowerPointService
	//////////////////////////////////////

	Electronics ElectronicsService
	Controller  ControllerService
	Sensor      SensorService
	Wiring      WiringService
	//////////////////////////////////////////////

	ChargerSys ChargerService
	Chassis    ChassisService
	Battery    BatteryService
	Emobile    EmobileService

	/////////////////////////////////
	Carcass CarcassService
	Doors   DoorsService
	Wings   WingsService
	Body    BodyService
}

func NewService(db *sqlx.DB, repo *repository.Repository) *Service {

	carcass := NewCarcassService(repo.Carcass)
	doors := NewDoorsService(repo.Doors)
	wings := NewWingsService(repo.Wings)

	sensor := NewSensorService(repo.Sensor)
	wiring := NewWiringService(repo.Wiring)
	controller := NewControllerService(repo.Controller)

	return &Service{
		Sensor:      sensor,
		Wiring:      wiring,
		Controller:  controller,
		Electronics: NewElectronicsService(db, repo.Electronics, wiring, controller, sensor),

		Carcass: carcass,
		Doors:   doors,
		Wings:   wings,
		Body:    NewBodyService(db, repo.Body, carcass, doors, wings),
	}
}

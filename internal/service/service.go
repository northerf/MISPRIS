package service

import (
	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
	"context"

	"github.com/jmoiron/sqlx"
)

// ── Существующие интерфейсы (без изменений) ──────────────────────

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
	Create(ctx context.Context, name string, batteryType string, batteryCapacity string, batteryInfo string) (string, error)
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

type ChargerSystemService interface {
	Create(ctx context.Context, chargerID string, connectorID string) (string, error)
	Update(ctx context.Context, id string, chargerID string, connectorID string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.ChargerSystem, error)
	GetByID(ctx context.Context, id string) (*domain.ChargerSystem, error)
}

type ChargerService interface {
	Create(ctx context.Context, chargerName, chargerInfo string) (string, error)
	Update(ctx context.Context, id string, chargerName, chargerInfo string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Charger, error)
	GetByID(ctx context.Context, id string) (*domain.Charger, error)
}

type ConnectorService interface {
	Create(ctx context.Context, connectorName, connectorInfo string) (string, error)
	Update(ctx context.Context, id string, connectorName, connectorInfo string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Connector, error)
	GetByID(ctx context.Context, id string) (*domain.Connector, error)
}

type ChassisService interface {
	Create(ctx context.Context, frameID, suspensionID, breakSystemID string) (string, error)
	Update(ctx context.Context, id string, frameID, suspensionID, breakSystemID string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Chassis, error)
	GetByID(ctx context.Context, id string) (*domain.Chassis, error)
}

type FrameService interface {
	Create(ctx context.Context, frameName, frameInfo string) (string, error)
	Update(ctx context.Context, id string, frameName, frameInfo string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Frame, error)
	GetByID(ctx context.Context, id string) (*domain.Frame, error)
}

type SuspensionService interface {
	Create(ctx context.Context, suspensionName, suspensionInfo string) (string, error)
	Update(ctx context.Context, id string, suspensionName, suspensionInfo string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Suspension, error)
	GetByID(ctx context.Context, id string) (*domain.Suspension, error)
}

type BreakSystemService interface {
	Create(ctx context.Context, breakSystemName, breakSystemInfo string) (string, error)
	Update(ctx context.Context, id string, breakSystemName, breakSystemInfo string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.BreakSystem, error)
	GetByID(ctx context.Context, id string) (*domain.BreakSystem, error)
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

type EngineService interface {
	Create(ctx context.Context, name, engineType, info string) (string, error)
	Update(ctx context.Context, id, name, engineType, info string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Engine, error)
	GetByID(ctx context.Context, id string) (*domain.Engine, error)
}

type InverterService interface {
	Create(ctx context.Context, name, info string) (string, error)
	Update(ctx context.Context, id, name, info string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Inverter, error)
	GetByID(ctx context.Context, id string) (*domain.Inverter, error)
}

type GearboxService interface {
	Create(ctx context.Context, name, info string) (string, error)
	Update(ctx context.Context, id, name, info string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Gearbox, error)
	GetByID(ctx context.Context, id string) (*domain.Gearbox, error)
}

type UserService interface {
	Create(ctx context.Context, username, password, role string, IsActive bool) (string, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, id, username, password, role string, IsActive bool) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.User, error)
	Authenticate(ctx context.Context, username, password string) (*domain.User, error)
}

// ── ПР2: Перечисления ─────────────────────────────────────────────

type EnumClassService interface {
	List(ctx context.Context) ([]*domain.EnumClass, error)
	GetByID(ctx context.Context, id string) (*domain.EnumClass, error)
	Create(ctx context.Context, name, componentType string) (string, error)
	Update(ctx context.Context, id, name, componentType string) error
	Delete(ctx context.Context, id string) error
	GetValues(ctx context.Context, id string) ([]*domain.EnumPosition, error)
	ValidateValue(ctx context.Context, enumClassID, value string) (bool, error)
}

type EnumPositionService interface {
	List(ctx context.Context) ([]*domain.EnumPosition, error)
	GetByID(ctx context.Context, id string) (*domain.EnumPosition, error)
	Create(ctx context.Context, enumClassID, value, orderNum string) (string, error)
	Update(ctx context.Context, id, value, orderNum string) error
	Delete(ctx context.Context, id string) error
}

// ── ПР3: Параметры ────────────────────────────────────────────────

type ParameterService interface {
	List(ctx context.Context) ([]*domain.Parameter, error)
	GetByID(ctx context.Context, id string) (*domain.Parameter, error)
	Create(ctx context.Context, designation, name, paramType, measuringUnit, enumClassID string) (string, error)
	Update(ctx context.Context, id, designation, name, paramType, measuringUnit, enumClassID string) error
	Delete(ctx context.Context, id string) error
}

type ComponentParameterService interface {
	List(ctx context.Context) ([]*domain.ComponentParameter, error)
	GetByID(ctx context.Context, id string) (*domain.ComponentParameter, error)
	Create(ctx context.Context, componentType, parameterID, orderNum, minVal, maxVal string) (string, error)
	Update(ctx context.Context, id, orderNum, minVal, maxVal string) error
	Delete(ctx context.Context, id string) error
	GetByType(ctx context.Context, componentType string) ([]*domain.ComponentParameterFull, error)
	CopyFromType(ctx context.Context, fromType, toType string) error
}

type EmobileParameterValueService interface {
	List(ctx context.Context) ([]*domain.EmobileParameterValue, error)
	GetByID(ctx context.Context, id string) (*domain.EmobileParameterValue, error)
	Create(ctx context.Context, emobileID, componentParameterID, valReal, valInt, valStr, enumValID string) (string, error)
	Update(ctx context.Context, id, valReal, valInt, valStr, enumValID string) error
	Delete(ctx context.Context, id string) error
	GetByEmobile(ctx context.Context, emobileID string) ([]*domain.EmobileParameterValue, error)
}

// ── ПР4: Хозяйственные операции ───────────────────────────────────

type ShdService interface {
	List(ctx context.Context) ([]*domain.SHD, error)
	GetByID(ctx context.Context, id string) (*domain.SHD, error)
	Create(ctx context.Context, name, shdType, inn, description string) (string, error)
	Update(ctx context.Context, id, name, shdType, inn, description string) error
	Delete(ctx context.Context, id string) error
}

type HoClassService interface {
	List(ctx context.Context) ([]*domain.HoClass, error)
	GetByID(ctx context.Context, id string) (*domain.HoClass, error)
	Create(ctx context.Context, name, designation, parentID string, isTerminal bool) (string, error)
	Update(ctx context.Context, id, name, designation, parentID string, isTerminal bool) error
	Delete(ctx context.Context, id string) error
	GetTerminal(ctx context.Context) ([]*domain.HoClass, error)
	GetChildren(ctx context.Context, parentID string) ([]*domain.HoClass, error)
}

type HoRoleService interface {
	List(ctx context.Context) ([]*domain.HoRole, error)
	GetByID(ctx context.Context, id string) (*domain.HoRole, error)
	Create(ctx context.Context, name, description string) (string, error)
	Update(ctx context.Context, id, name, description string) error
	Delete(ctx context.Context, id string) error
}

type HoClassRoleService interface {
	List(ctx context.Context, hoClassID string) ([]*domain.HoClassRole, error)
	ListByClass(ctx context.Context, hoClassID string) ([]*domain.HoClassRole, error)
	Create(ctx context.Context, hoClassID, hoRoleID string, isRequired bool) (string, error)
	Delete(ctx context.Context, id string) error
}

type HoClassParameterService interface {
	List(ctx context.Context, hoClassID string) ([]*domain.HoClassParameter, error)
	GetByID(ctx context.Context, id string) (*domain.HoClassParameter, error)
	Create(ctx context.Context, hoClassID, parameterID string, orderNum int, minVal, maxVal float64) (string, error)
	Update(ctx context.Context, id string, orderNum int, minVal, maxVal float64) error
	Delete(ctx context.Context, id string) error
	// GetByHoClass вызывает SQL-функцию get_ho_class_parameters
	GetByHoClass(ctx context.Context, hoClassID string) ([]*domain.HoClassParameterFull, error)
	// CopyFromClass копирует параметры INSERT...SELECT
	CopyFromClass(ctx context.Context, fromClassID, toClassID string) error
}

type DocumentClassService interface {
	List(ctx context.Context) ([]*domain.DocumentClass, error)
	GetByID(ctx context.Context, id string) (*domain.DocumentClass, error)
	Create(ctx context.Context, name, code, description string) (string, error)
	Update(ctx context.Context, id, name, code, description string) error
	Delete(ctx context.Context, id string) error
}

type HoClassDocumentService interface {
	ListByClass(ctx context.Context, hoClassID string) ([]*domain.HoClassDocument, error)
	Create(ctx context.Context, hoClassID, docClassID, roleName string, isRequired bool) (string, error)
	Delete(ctx context.Context, id string) error
}

type HoInstanceService interface {
	List(ctx context.Context, hoClassID string) ([]*domain.HoInstance, error)
	GetByID(ctx context.Context, id string) (*domain.HoInstance, error)
	Create(ctx context.Context, hoClassID, docNumber, docDate string, totalAmount float64, note string) (string, error)
	Update(ctx context.Context, id, status, docNumber, docDate string, totalAmount float64, note string) error
	Delete(ctx context.Context, id string) error
	// FindByClass вызывает SQL-функцию find_ho_by_class
	FindByClass(ctx context.Context, hoClassID string) ([]*domain.HoInstanceFull, error)
}

type HoActorService interface {
	ListByHo(ctx context.Context, hoID string) ([]*domain.HoActor, error)
	// Create вызывает SQL-процедуру set_ho_actor
	Create(ctx context.Context, hoID, hoRoleID, shdID string) (string, error)
	Delete(ctx context.Context, id string) error
}

type HoParameterValueService interface {
	ListByHo(ctx context.Context, hoID string) ([]*domain.HoParameterValueFull, error)
	GetByID(ctx context.Context, id string) (*domain.HoParameterValue, error)
	// Create вызывает SQL-процедуру write_ho_par
	Create(ctx context.Context, hoID, hoClassParameterID string, valReal float64, valInt int, valStr, valDate, enumValID string) (string, error)
	Update(ctx context.Context, id, hoID, hoClassParameterID string, valReal float64, valInt int, valStr, valDate, enumValID string) error
	Delete(ctx context.Context, id string) error
}

type HoDocumentService interface {
	ListByHo(ctx context.Context, hoID string) ([]*domain.HoDocument, error)
	Create(ctx context.Context, hoID, docClassID, docNumber, docDate, note string) (string, error)
	Delete(ctx context.Context, id string) error
}

type HoPositionService interface {
	ListByHo(ctx context.Context, hoID string) ([]*domain.HoPositionFull, error)
	// Create вызывает SQL-процедуру add_ho_position
	Create(ctx context.Context, hoID, emobileID string, quantity int, unitPrice float64, note string) (string, error)
	Update(ctx context.Context, id string, quantity int, unitPrice float64, note string) error
	Delete(ctx context.Context, id string) error
}

// ── Агрегирующий Service ──────────────────────────────────────────

type Service struct {
	PowerPoint   PowerPointService
	Engine       EngineService
	Inverter     InverterService
	Gearbox      GearboxService
	Electronics  ElectronicsService
	Controller   ControllerService
	Sensor       SensorService
	Wiring       WiringService
	ChargerSystem ChargerSystemService
	Charger      ChargerService
	Connector    ConnectorService
	Chassis      ChassisService
	Frame        FrameService
	Suspension   SuspensionService
	BreakSystem  BreakSystemService
	Battery      BatteryService
	Emobile      EmobileService
	Carcass      CarcassService
	Doors        DoorsService
	Wings        WingsService
	Body         BodyService
	User         UserService

	// ПР2
	EnumClass    EnumClassService
	EnumPosition EnumPositionService

	// ПР3
	Parameter             ParameterService
	ComponentParameter    ComponentParameterService
	EmobileParameterValue EmobileParameterValueService

	// ПР4
	Shd              ShdService
	HoClass          HoClassService
	HoRole           HoRoleService
	HoClassRole      HoClassRoleService
	HoClassParameter HoClassParameterService
	DocumentClass    DocumentClassService
	HoClassDocument  HoClassDocumentService
	HoInstance       HoInstanceService
	HoActor          HoActorService
	HoParameterValue HoParameterValueService
	HoDocument       HoDocumentService
	HoPosition       HoPositionService
}

func NewService(db *sqlx.DB, repos *repository.Repository) *Service {
	// ── Шаг 1: листовые сервисы (только репозиторий, без зависимостей) ──
	engine      := NewEngineService(repos.Engine)
	inverter    := NewInverterService(repos.Inverter)
	gearbox     := NewGearboxService(repos.Gearbox)
	controller  := NewControllerService(repos.Controller)
	sensor      := NewSensorService(repos.Sensor)
	wiring      := NewWiringService(repos.Wiring)
	charger     := NewChargerService(repos.Charger)
	connector   := NewConnectorService(repos.Connector)
	frame       := NewFrameService(repos.Frame)
	suspension  := NewSuspensionService(repos.Suspension)
	breakSystem := NewBreakSystemService(repos.BreakSystem)
	battery     := NewBatteryService(repos.Battery, db) // repo первым, затем db
	carcass     := NewCarcassService(repos.Carcass)
	doors       := NewDoorsService(repos.Doors)
	wings       := NewWingsService(repos.Wings)
	user        := NewUserService(repos.User)

	// ── Шаг 2: составные сервисы (зависят от листовых + db) ─────────────
	powerPoint    := NewPowerPointService(repos.PowerPoint, db, engine, inverter, gearbox)
	electronics   := NewElectronicsService(db, repos.Electronics, wiring, controller, sensor)
	chargerSystem := NewChargerSystemService(db, repos.ChargerSystem, charger, connector)
	chassis       := NewChassisService(db, repos.Chassis, frame, suspension, breakSystem)
	body          := NewBodyService(db, repos.Body, carcass, doors, wings)
	emobile       := NewEmobileService(db, repos.Emobile, chargerSystem, body, electronics, chassis, battery, powerPoint)

	return &Service{
		PowerPoint:    powerPoint,
		Engine:        engine,
		Inverter:      inverter,
		Gearbox:       gearbox,
		Electronics:   electronics,
		Controller:    controller,
		Sensor:        sensor,
		Wiring:        wiring,
		ChargerSystem: chargerSystem,
		Charger:       charger,
		Connector:     connector,
		Chassis:       chassis,
		Frame:         frame,
		Suspension:    suspension,
		BreakSystem:   breakSystem,
		Battery:       battery,
		Emobile:       emobile,
		Carcass:       carcass,
		Doors:         doors,
		Wings:         wings,
		Body:          body,
		User:          user,

		// ПР2
		EnumClass:    NewEnumClassService(repos.EnumClass),
		EnumPosition: NewEnumPositionService(repos.EnumPosition),

		// ПР3
		Parameter:             NewParameterService(repos.Parameter),
		ComponentParameter:    NewComponentParameterService(repos.ComponentParameter),
		EmobileParameterValue: NewEmobileParameterValueService(repos.EmobileParameterValue),

		// ПР4
		Shd:              NewShdService(repos.Shd),
		HoClass:          NewHoClassService(repos.HoClass),
		HoRole:           NewHoRoleService(repos.HoRole),
		HoClassRole:      NewHoClassRoleService(repos.HoClassRole),
		HoClassParameter: NewHoClassParameterService(repos.HoClassParameter),
		DocumentClass:    NewDocumentClassService(repos.DocumentClass),
		HoClassDocument:  NewHoClassDocumentService(repos.HoClassDocument),
		HoInstance:       NewHoInstanceService(repos.HoInstance),
		HoActor:          NewHoActorService(repos.HoActor),
		HoParameterValue: NewHoParameterValueService(repos.HoParameterValue),
		HoDocument:       NewHoDocumentService(repos.HoDocument),
		HoPosition:       NewHoPositionService(repos.HoPosition),
	}
}

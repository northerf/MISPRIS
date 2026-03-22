package repository

import (
	"context"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type BodyRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Body, error)
	CreateTx(ctx context.Context, tx *sqlx.Tx, b *domain.Body) (string, error)
	Update(ctx context.Context, b *domain.Body) error
	Delete(ctx context.Context, id string) error
}

type CarcassRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Carcass, error)
	Create(ctx context.Context, c *domain.Carcass) (string, error)
	Update(ctx context.Context, c *domain.Carcass) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Carcass, error) // каталог
}

type DoorsRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Doors, error)
	Create(ctx context.Context, d *domain.Doors) (string, error)
	Update(ctx context.Context, d *domain.Doors) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Doors, error)
}

type WingsRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Wings, error)
	Create(ctx context.Context, w *domain.Wings) (string, error)
	Update(ctx context.Context, w *domain.Wings) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Wings, error)
}

//////////////////////////////////////////////////////////////////////

type ElectronicsRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Electronics, error)
	CreateTx(ctx context.Context, tx *sqlx.Tx, e *domain.Electronics) (string, error)
	Update(ctx context.Context, e *domain.Electronics) error
	Delete(ctx context.Context, id string) error
}

type ControllerRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Controller, error)
	Create(ctx context.Context, c *domain.Controller) (string, error)
	Update(ctx context.Context, c *domain.Controller) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Controller, error)
}

type SensorRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Sensor, error)
	Create(ctx context.Context, s *domain.Sensor) (string, error)
	Update(ctx context.Context, s *domain.Sensor) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Sensor, error)
}

type WiringRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Wiring, error)
	Create(ctx context.Context, w *domain.Wiring) (string, error)
	Update(ctx context.Context, w *domain.Wiring) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Wiring, error)
}

/////////////////////////////////////////////////////////////////////////

type ChargerSystemRepository interface {
	GetByID(ctx context.Context, id string) (*domain.ChargerSystem, error)
	CreateTx(ctx context.Context, tx *sqlx.Tx, cs *domain.ChargerSystem) (string, error)
	Update(ctx context.Context, cs *domain.ChargerSystem) error
	Delete(ctx context.Context, id string) error
}

type ChargerRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Charger, error)
	Create(ctx context.Context, c *domain.Charger) (string, error)
	Update(ctx context.Context, c *domain.Charger) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Charger, error)
}

type ConnectorRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Connector, error)
	Create(ctx context.Context, c *domain.Connector) (string, error)
	Update(ctx context.Context, c *domain.Connector) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Connector, error)
}

//////////////////////////////////////////////////////////////////////

type ChassisRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Chassis, error)
	CreateTx(ctx context.Context, tx *sqlx.Tx, c *domain.Chassis) (string, error)
	Update(ctx context.Context, c *domain.Chassis) error
	Delete(ctx context.Context, id string) error
}

type FrameRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Frame, error)
	Create(ctx context.Context, f *domain.Frame) (string, error)
	Update(ctx context.Context, f *domain.Frame) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Frame, error)
}

type SuspensionRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Suspension, error)
	Create(ctx context.Context, s *domain.Suspension) (string, error)
	Update(ctx context.Context, s *domain.Suspension) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Suspension, error)
}

type BreakSystemRepository interface {
	GetByID(ctx context.Context, id string) (*domain.BreakSystem, error)
	Create(ctx context.Context, b *domain.BreakSystem) (string, error)
	Update(ctx context.Context, b *domain.BreakSystem) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.BreakSystem, error)
}

/////////////////////////////////////////////////////////////////////////

type BatteryRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Battery, error)
	Create(ctx context.Context, b *domain.Battery) (string, error)
	CreateTx(ctx context.Context, tx *sqlx.Tx, b *domain.Battery) (string, error)
	Update(ctx context.Context, b *domain.Battery) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Battery, error)
}
type PowerPointRepository interface {
	GetByID(ctx context.Context, id string) (*domain.PowerPoint, error)
	Create(ctx context.Context, p *domain.PowerPoint) (string, error)
	CreateTx(ctx context.Context, tx *sqlx.Tx, p *domain.PowerPoint) (string, error)
	Update(ctx context.Context, p *domain.PowerPoint) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.PowerPoint, error)
}
type EngineRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Engine, error)
	Create(ctx context.Context, e *domain.Engine) (string, error)
	Update(ctx context.Context, e *domain.Engine) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Engine, error)
}

type InverterRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Inverter, error)
	Create(ctx context.Context, i *domain.Inverter) (string, error)
	Update(ctx context.Context, i *domain.Inverter) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Inverter, error)
}

type GearboxRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Gearbox, error)
	Create(ctx context.Context, g *domain.Gearbox) (string, error)
	Update(ctx context.Context, g *domain.Gearbox) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Gearbox, error)
}

type EmobileRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Emobile, error)
	Create(ctx context.Context, tx *sqlx.Tx, emobile *domain.Emobile) (string, error)
	Update(ctx context.Context, emobile *domain.Emobile) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Emobile, error)
}

/////////////////////////////////////////////////////////////////////////

type Repository struct {
	Emobile       EmobileRepository
	Body          BodyRepository
	Electronics   ElectronicsRepository
	Chassis       ChassisRepository
	ChargerSystem ChargerSystemRepository
	Battery       BatteryRepository
	PowerPoint    PowerPointRepository

	// body
	Carcass CarcassRepository
	Doors   DoorsRepository
	Wings   WingsRepository

	// electronics
	Controller ControllerRepository
	Sensor     SensorRepository
	Wiring     WiringRepository

	// charger system
	Charger   ChargerRepository
	Connector ConnectorRepository

	// chassis
	Frame       FrameRepository
	Suspension  SuspensionRepository
	BreakSystem BreakSystemRepository

	// power point
	Engine   EngineRepository
	Inverter InverterRepository
	Gearbox  GearboxRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Emobile:       NewEmobilePostgres(db),
		Body:          NewBodyPostgres(db),
		Electronics:   NewElectronicsPostgres(db),
		Chassis:       NewChassisPostgres(db),
		ChargerSystem: NewChargerSystemPostgres(db),
		Battery:       NewBatteryPostgres(db),
		PowerPoint:    NewPowerPointPostgres(db),

		// body
		Carcass: NewCarcassPostgres(db),
		Doors:   NewDoorsPostgres(db),
		Wings:   NewWingsPostgres(db),

		// electronics
		Controller: NewControllerPostgres(db),
		Sensor:     NewSensorPostgres(db),
		Wiring:     NewWiringPostgres(db),

		// charger_system
		Charger:   NewChargerPostgres(db),
		Connector: NewConnectorPostgres(db),

		// chassis
		Frame:       NewFramePostgres(db),
		Suspension:  NewSuspensionPostgres(db),
		BreakSystem: NewBreakSystemPostgres(db),

		// power point
		Engine:   NewEnginePostgres(db),
		Inverter: NewInverterPostgres(db),
		Gearbox:  NewGearboxPostgres(db),
	}
}

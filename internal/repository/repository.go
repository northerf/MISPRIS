package repository

import (
	"context"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type BodyRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Body, error)
	CreateTx(ctx context.Context, tx *sqlx.Tx, b *domain.Body) (int64, error)
	Update(ctx context.Context, b *domain.Body) error
	Delete(ctx context.Context, id int64) error
}

type CarcassRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Carcass, error)
	Create(ctx context.Context, c *domain.Carcass) (int64, error)
	Update(ctx context.Context, c *domain.Carcass) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*domain.Carcass, error) // каталог
}

type DoorsRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Doors, error)
	Create(ctx context.Context, d *domain.Doors) (int64, error)
	Update(ctx context.Context, d *domain.Doors) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*domain.Doors, error)
}

type WingsRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Wings, error)
	Create(ctx context.Context, w *domain.Wings) (int64, error)
	Update(ctx context.Context, w *domain.Wings) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*domain.Wings, error)
}

//////////////////////////////////////////////////////////////////////

type ElectronicsRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Electronics, error)
	CreateTx(ctx context.Context, tx *sqlx.Tx, e *domain.Electronics) (int64, error)
	Update(ctx context.Context, e *domain.Electronics) error
	Delete(ctx context.Context, id int64) error
}

type ControllerRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Controller, error)
	Create(ctx context.Context, c *domain.Controller) (int64, error)
	Update(ctx context.Context, c *domain.Controller) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*domain.Controller, error)
}

type SensorRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Sensor, error)
	Create(ctx context.Context, s *domain.Sensor) (int64, error)
	Update(ctx context.Context, s *domain.Sensor) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*domain.Sensor, error)
}

type WiringRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Wiring, error)
	Create(ctx context.Context, w *domain.Wiring) (int64, error)
	Update(ctx context.Context, w *domain.Wiring) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*domain.Wiring, error)
}

/////////////////////////////////////////////////////////////////////////

type ChargerSystemRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.ChargerSystem, error)
	CreateTx(ctx context.Context, tx *sqlx.Tx, cs *domain.ChargerSystem) (int64, error)
	Update(ctx context.Context, cs *domain.ChargerSystem) error
	Delete(ctx context.Context, id int64) error
}

type ChargerRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Charger, error)
	Create(ctx context.Context, c *domain.Charger) (int64, error)
	Update(ctx context.Context, c *domain.Charger) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*domain.Charger, error)
}

type ConnectorRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Connector, error)
	Create(ctx context.Context, c *domain.Connector) (int64, error)
	Update(ctx context.Context, c *domain.Connector) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*domain.Connector, error)
}

//////////////////////////////////////////////////////////////////////

type ChassisRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Chassis, error)
	CreateTx(ctx context.Context, tx *sqlx.Tx, c *domain.Chassis) (int64, error)
	Update(ctx context.Context, c *domain.Chassis) error
	Delete(ctx context.Context, id int64) error
}

type FrameRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Frame, error)
	Create(ctx context.Context, f *domain.Frame) (int64, error)
	Update(ctx context.Context, f *domain.Frame) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*domain.Frame, error)
}

type SuspensionRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Suspension, error)
	Create(ctx context.Context, s *domain.Suspension) (int64, error)
	Update(ctx context.Context, s *domain.Suspension) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*domain.Suspension, error)
}

type BreakSystemRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.BreakSystem, error)
	Create(ctx context.Context, b *domain.BreakSystem) (int64, error)
	Update(ctx context.Context, b *domain.BreakSystem) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*domain.BreakSystem, error)
}

/////////////////////////////////////////////////////////////////////////

type BatteryRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Battery, error)
	Create(ctx context.Context, b *domain.Battery) (int64, error)
	Update(ctx context.Context, b *domain.Battery) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*domain.Battery, error)
}

type PowerPointRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.PowerPoint, error)
	Create(ctx context.Context, p *domain.PowerPoint) (int64, error)
	Update(ctx context.Context, p *domain.PowerPoint) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*domain.PowerPoint, error)
}

type EngineRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Engine, error)
	Create(ctx context.Context, e *domain.Engine) (int64, error)
	Update(ctx context.Context, e *domain.Engine) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*domain.Engine, error)
}

type InverterRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Inverter, error)
	Create(ctx context.Context, i *domain.Inverter) (int64, error)
	Update(ctx context.Context, i *domain.Inverter) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*domain.Inverter, error)
}

type GearboxRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Gearbox, error)
	Create(ctx context.Context, g *domain.Gearbox) (int64, error)
	Update(ctx context.Context, g *domain.Gearbox) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*domain.Gearbox, error)
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

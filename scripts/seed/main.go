package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getenv("POSTGRES_HOST", "postgres"),
		getenv("POSTGRES_PORT", "5432"),
		getenv("POSTGRES_USER", "postgres"),
		getenv("POSTGRES_PASSWORD", "postgres"),
		getenv("POSTGRES_DB", "mispris"),
		getenv("SSL_MODE", "disable"),
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("seed: connect error: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	log.Println("seed: starting...")
	if err := runSeed(ctx, db); err != nil {
		log.Fatalf("seed: error: %v", err)
	}
	if err := runSeedPR2PR3(ctx, db); err != nil {
		log.Printf("seed PR2/PR3: %v (non-fatal)", err)
	}
	log.Println("seed: done ✓")
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func runSeed(ctx context.Context, db *sqlx.DB) error {
	// Идемпотентная проверка — не заполняем повторно
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM emobile`).Scan(&count); err != nil {
		return fmt.Errorf("check emobile: %w", err)
	}
	if count > 0 {
		log.Printf("seed: already has %d emobile records, skipping", count)
		return nil
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// ── Двигатели ────────────────────────────────────────────
	engines := []struct{ name, etype, info string }{
		{"PMSM-150", "AC", "Постоянный магнит, синхронный, 150 кВт, 400 Нм"},
		{"PMSM-220", "AC", "Постоянный магнит, синхронный, 220 кВт, 600 Нм"},
		{"BLDC-90",  "DC", "Бесщёточный постоянного тока, 90 кВт, 250 Нм"},
		{"BLDC-120", "DC", "Бесщёточный постоянного тока, 120 кВт, 340 Нм"},
		{"IM-180",   "AC", "Асинхронный индукционный, 180 кВт, 430 Нм"},
	}
	engineIDs := make([]int, len(engines))
	for i, e := range engines {
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO engine (engine_name,engine_type,engine_info) VALUES($1,$2,$3) RETURNING engine_id`,
			e.name, e.etype, e.info).Scan(&engineIDs[i]); err != nil {
			return fmt.Errorf("engine %s: %w", e.name, err)
		}
	}
	log.Printf("seed: engines inserted (%d)", len(engineIDs))

	// ── Инверторы ────────────────────────────────────────────
	inverters := []struct{ name, info string }{
		{"INV-SIC-150",  "Инвертор SiC MOSFET, 150 кВт, КПД 98.5%"},
		{"INV-SIC-250",  "Инвертор SiC MOSFET, 250 кВт, КПД 98.8%"},
		{"INV-IGBT-100", "Инвертор IGBT, 100 кВт, КПД 97.2%"},
	}
	invIDs := make([]int, len(inverters))
	for i, inv := range inverters {
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO inverter (inverter_name,inverter_info) VALUES($1,$2) RETURNING inverter_id`,
			inv.name, inv.info).Scan(&invIDs[i]); err != nil {
			return fmt.Errorf("inverter %s: %w", inv.name, err)
		}
	}

	// ── КПП ──────────────────────────────────────────────────
	gearboxes := []struct{ name, info string }{
		{"1-ступ. редуктор 9.7", "Передаточное число 9.7:1, масса 28 кг"},
		{"1-ступ. редуктор 8.2", "Передаточное число 8.2:1, масса 24 кг"},
		{"2-ступ. AMT 8.9/4.1",  "2-ступенчатая с автоматическим переключением"},
	}
	gbIDs := make([]int, len(gearboxes))
	for i, g := range gearboxes {
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO gearbox (gearbox_name,gearbox_info) VALUES($1,$2) RETURNING gearbox_id`,
			g.name, g.info).Scan(&gbIDs[i]); err != nil {
			return fmt.Errorf("gearbox %s: %w", g.name, err)
		}
	}

	// ── Силовые установки ────────────────────────────────────
	// [engine_idx, inverter_idx, gearbox_idx]
	ppDefs := [][3]int{{0,0,0},{1,1,2},{2,2,1},{4,1,0}}
	ppIDs := make([]int, len(ppDefs))
	for i, d := range ppDefs {
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO power_point (engine_id,inverter_id,gearbox_id) VALUES($1,$2,$3) RETURNING power_point_id`,
			engineIDs[d[0]], invIDs[d[1]], gbIDs[d[2]]).Scan(&ppIDs[i]); err != nil {
			return fmt.Errorf("power_point %d: %w", i, err)
		}
	}
	log.Printf("seed: power_points inserted (%d)", len(ppIDs))

	// ── Батареи ──────────────────────────────────────────────
	type batRow struct{ name, btype, info string; cap float64 }
	batteries := []batRow{
		{"NMC-60",     "Li-ion",     "NMC химия, 360 В номинал, 350 циклов до 80%",        60.0},
		{"NMC-80",     "Li-ion",     "NMC химия, 400 В номинал, 400 циклов до 80%",        80.0},
		{"LFP-100",    "Li-ion",     "LFP химия, 320 В номинал, 2000+ циклов",            100.0},
		{"LFP-120",    "Li-ion",     "LFP химия, 336 В номинал, улучшенный теплообмен",   120.0},
		{"Polymer-55", "Li-polymer", "Гибкие ячейки, облегчённый корпус, 55 кВтч",        55.0},
	}
	batIDs := make([]int, len(batteries))
	for i, b := range batteries {
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO battery (battery_name,battery_type,battery_capacity,battery_info) VALUES($1,$2,$3,$4) RETURNING battery_id`,
			b.name, b.btype, b.cap, b.info).Scan(&batIDs[i]); err != nil {
			return fmt.Errorf("battery %s: %w", b.name, err)
		}
	}

	// ── Зарядные устройства ──────────────────────────────────
	chargers := []struct{ name, info string }{
		{"OBC-7kW",  "Бортовое зарядное устройство, 7.4 кВт, однофазное"},
		{"OBC-11kW", "Бортовое зарядное устройство, 11 кВт, трёхфазное"},
		{"OBC-22kW", "Бортовое зарядное устройство, 22 кВт, трёхфазное"},
	}
	chIDs := make([]int, len(chargers))
	for i, c := range chargers {
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO charger (charger_name,charger_info) VALUES($1,$2) RETURNING charger_id`,
			c.name, c.info).Scan(&chIDs[i]); err != nil {
			return fmt.Errorf("charger %s: %w", c.name, err)
		}
	}

	// ── Коннекторы ───────────────────────────────────────────
	connectors := []struct{ name, info string }{
		{"Type 2 (Mennekes)", "Стандарт IEC 62196-2, до 43 кВт AC"},
		{"CCS2",              "Combined Charging System 2, AC+DC до 350 кВт"},
		{"CHAdeMO 2.0",       "CHAdeMO протокол, DC до 400 кВт"},
		{"GB/T 20234.3",      "Китайский стандарт DC, до 250 кВт"},
	}
	connIDs := make([]int, len(connectors))
	for i, c := range connectors {
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO connector (connector_name,connector_info) VALUES($1,$2) RETURNING connector_id`,
			c.name, c.info).Scan(&connIDs[i]); err != nil {
			return fmt.Errorf("connector %s: %w", c.name, err)
		}
	}

	// ── Зарядные системы [charger_idx, connector_idx] ────────
	csDefs := [][2]int{{0,0},{1,0},{2,1},{1,2}}
	csIDs := make([]int, len(csDefs))
	for i, d := range csDefs {
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO charger_system (charger_id,connector_id) VALUES($1,$2) RETURNING charger_system_id`,
			chIDs[d[0]], connIDs[d[1]]).Scan(&csIDs[i]); err != nil {
			return fmt.Errorf("charger_system %d: %w", i, err)
		}
	}

	// ── Рамы ─────────────────────────────────────────────────
	frames := []struct{ name, info string }{
		{"Стальная лонжерон.", "Лонжеронная рама, сталь S355, масса 180 кг"},
		{"Алюм. пространств.", "Алюминиевая пространственная рама, масса 95 кг"},
		{"Несущий кузов",      "Монокок из высокопрочной стали, интегрирован с кузовом"},
	}
	frIDs := make([]int, len(frames))
	for i, f := range frames {
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO frame (frame_name,frame_info) VALUES($1,$2) RETURNING frame_id`,
			f.name, f.info).Scan(&frIDs[i]); err != nil {
			return fmt.Errorf("frame %s: %w", f.name, err)
		}
	}

	// ── Подвески ─────────────────────────────────────────────
	suspensions := []struct{ name, info string }{
		{"МакФерсон перед.",   "Стойки МакФерсон спереди, стабилизатор 22 мм"},
		{"Многорычажная зад.", "Многорычажная подвеска сзади, пружины и амортизаторы"},
		{"Двойной рычаг 4x4", "Независимая двойной рычаг передняя и задняя"},
	}
	susIDs := make([]int, len(suspensions))
	for i, s := range suspensions {
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO suspension (suspension_name,suspension_info) VALUES($1,$2) RETURNING suspension_id`,
			s.name, s.info).Scan(&susIDs[i]); err != nil {
			return fmt.Errorf("suspension %s: %w", s.name, err)
		}
	}

	// ── Тормозные системы ────────────────────────────────────
	breakSystems := []struct{ name, info string }{
		{"EHB-стандарт", "Электрогидравлические тормоза, ABS+ESP+рекуперация"},
		{"EHB-спорт",    "Перфорированные диски 340 мм, 6-поршневые суппорты"},
		{"EBB-базовый",  "Электронный вакуумный усилитель, дисковые 4x4"},
	}
	bsIDs := make([]int, len(breakSystems))
	for i, b := range breakSystems {
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO break_system (break_system_name,break_info) VALUES($1,$2) RETURNING break_system_id`,
			b.name, b.info).Scan(&bsIDs[i]); err != nil {
			return fmt.Errorf("break_system %s: %w", b.name, err)
		}
	}

	// ── Шасси [frame_idx, suspension_idx, break_idx] ─────────
	chassisDefs := [][3]int{{0,0,0},{1,0,1},{2,2,2}}
	chassisIDs := make([]int, len(chassisDefs))
	for i, d := range chassisDefs {
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO chassis (frame_id,suspension_id,break_system_id) VALUES($1,$2,$3) RETURNING chassis_id`,
			frIDs[d[0]], susIDs[d[1]], bsIDs[d[2]]).Scan(&chassisIDs[i]); err != nil {
			return fmt.Errorf("chassis %d: %w", i, err)
		}
	}
	log.Printf("seed: chassis inserted (%d)", len(chassisIDs))

	// ── Каркасы ──────────────────────────────────────────────
	carcasses := []struct{ name, info string }{
		{"Седан B-класс",  "Трёхобъёмный кузов, длина 4350 мм, сталь+пластик"},
		{"Хэтчбек",        "Двухобъёмный, длина 4100 мм, усиленные стойки"},
		{"SUV кроссовер",  "Высокий кузов, клиренс 195 мм, длина 4600 мм"},
	}
	carIDs := make([]int, len(carcasses))
	for i, c := range carcasses {
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO carcass (carcass_name,carcass_info) VALUES($1,$2) RETURNING carcass_id`,
			c.name, c.info).Scan(&carIDs[i]); err != nil {
			return fmt.Errorf("carcass %s: %w", c.name, err)
		}
	}

	// ── Двери ────────────────────────────────────────────────
	doorsData := []struct{ name, info string }{
		{"4 двери стандарт", "Стандартные 4 двери, опускные стёкла"},
		{"5 дверей лифтбек", "4 двери + дверь багажника, электропривод задней"},
		{"2 двери купе",     "2 широкие двери, безрамочные стёкла"},
	}
	doorIDs := make([]int, len(doorsData))
	for i, d := range doorsData {
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO doors (doors_name,doors_info) VALUES($1,$2) RETURNING doors_id`,
			d.name, d.info).Scan(&doorIDs[i]); err != nil {
			return fmt.Errorf("doors %s: %w", d.name, err)
		}
	}

	// ── Крылья ───────────────────────────────────────────────
	wingsData := []struct{ name, info string }{
		{"Стальные штамп.", "Штампованные стальные, совмещены с арками"},
		{"Алюм. передние",  "Передние крылья алюминиевые, снижение массы"},
		{"Пластик ABS",     "Пластиковые расширители арок, съёмные"},
	}
	wingIDs := make([]int, len(wingsData))
	for i, w := range wingsData {
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO wings (wings_name,wings_info) VALUES($1,$2) RETURNING wings_id`,
			w.name, w.info).Scan(&wingIDs[i]); err != nil {
			return fmt.Errorf("wings %s: %w", w.name, err)
		}
	}

	// ── Кузова [carcass_idx, doors_idx, wings_idx] ───────────
	bodyDefs := [][3]int{{0,0,0},{1,1,1},{2,1,2}}
	bodyIDs := make([]int, len(bodyDefs))
	for i, d := range bodyDefs {
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO body (carcass_id,doors_id,wings_id) VALUES($1,$2,$3) RETURNING body_id`,
			carIDs[d[0]], doorIDs[d[1]], wingIDs[d[2]]).Scan(&bodyIDs[i]); err != nil {
			return fmt.Errorf("body %d: %w", i, err)
		}
	}

	// ── Контроллеры ──────────────────────────────────────────
	controllers := []struct{ name, info string }{
		{"VCU-v3",  "Главный контроллер автомобиля, CAN 2.0B, AUTOSAR"},
		{"BMS-v5",  "Система управления батареей, балансировка ячеек, SOC/SOH"},
		{"MCU-Pro", "Контроллер двигателя, векторное управление, 400 В"},
	}
	ctrlIDs := make([]int, len(controllers))
	for i, c := range controllers {
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO controllers (controller_name,controller_info) VALUES($1,$2) RETURNING controller_id`,
			c.name, c.info).Scan(&ctrlIDs[i]); err != nil {
			return fmt.Errorf("controller %s: %w", c.name, err)
		}
	}

	// ── Датчики ──────────────────────────────────────────────
	sensors := []struct{ name, info string }{
		{"IMU-6DOF",     "6-осевой инерциальный датчик, ±16g, CAN/SPI"},
		{"TPMS-4x",      "Датчики давления 4 колеса, BLE 5.0"},
		{"Temp-pack-8x", "8 термодатчиков батарейного блока, NTC 10кОм"},
	}
	sensorIDs := make([]int, len(sensors))
	for i, s := range sensors {
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO sensors (sensor_name,sensor_info) VALUES($1,$2) RETURNING sensor_id`,
			s.name, s.info).Scan(&sensorIDs[i]); err != nil {
			return fmt.Errorf("sensor %s: %w", s.name, err)
		}
	}

	// ── Проводка ─────────────────────────────────────────────
	wirings := []struct{ name, info string }{
		{"HV-жгут 400В",    "Высоковольтная проводка 400 В, сечение 50 мм², экранирование"},
		{"LV-жгут 12В",     "Низковольтная CAN-шина + питание ЭБУ"},
		{"Полный комплект", "HV + LV + шина данных, ресурс 300 000 км"},
	}
	wirIDs := make([]int, len(wirings))
	for i, w := range wirings {
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO wiring (wiring_name,wiring_info) VALUES($1,$2) RETURNING wiring_id`,
			w.name, w.info).Scan(&wirIDs[i]); err != nil {
			return fmt.Errorf("wiring %s: %w", w.name, err)
		}
	}

	// ── Электроника [ctrl_idx, sensor_idx, wiring_idx] ───────
	elDefs := [][3]int{{0,0,0},{1,2,2},{2,1,1}}
	elIDs := make([]int, len(elDefs))
	for i, d := range elDefs {
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO electronics (controller_id,sensor_id,wiring_id) VALUES($1,$2,$3) RETURNING electronics_id`,
			ctrlIDs[d[0]], sensorIDs[d[1]], wirIDs[d[2]]).Scan(&elIDs[i]); err != nil {
			return fmt.Errorf("electronics %d: %w", i, err)
		}
	}
	log.Printf("seed: electronics inserted (%d)", len(elIDs))

	// ── Электромобили ─────────────────────────────────────────
	// [name, pp_idx, bat_idx, cs_idx, chassis_idx, body_idx, el_idx]
	type emRow struct {
		name string
		pp, bat, cs, chs, bod, el int
	}
	emobiles := []emRow{
		{"Emobile City S",    0, 0, 0, 0, 0, 0},
		{"Emobile Pro X",     1, 2, 2, 1, 1, 1},
		{"Emobile Cargo SUV", 3, 3, 1, 2, 2, 2},
		{"Emobile Lite",      2, 4, 0, 0, 1, 0},
	}
	for _, e := range emobiles {
		var id int
		if err = tx.QueryRowContext(ctx,
			`INSERT INTO emobile (emobile_name,power_point_id,battery_id,charger_system_id,chassis_id,body_id,electronics_id)
			 VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING emobile_id`,
			e.name, ppIDs[e.pp], batIDs[e.bat], csIDs[e.cs],
			chassisIDs[e.chs], bodyIDs[e.bod], elIDs[e.el]).Scan(&id); err != nil {
			return fmt.Errorf("emobile %s: %w", e.name, err)
		}
		log.Printf("seed: emobile '%s' id=%d", e.name, id)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	// Обновляем материализованное представление
	if _, err = db.ExecContext(ctx, `REFRESH MATERIALIZED VIEW emobile_full`); err != nil {
		log.Printf("seed: warning — could not refresh emobile_full: %v", err)
	} else {
		log.Println("seed: emobile_full refreshed ✓")
	}

	return nil
}


func runSeedPR2PR3(ctx context.Context, db *sqlx.DB) error {
	// Идемпотентная проверка
	var cnt int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM enum_class`).Scan(&cnt); err != nil {
		return nil // таблица может не существовать
	}
	if cnt > 0 {
		log.Printf("seed PR2/PR3: already seeded (%d enum_class records)", cnt)
		return nil
	}

	// ── ПР2: Классы перечислений ──────────────────────────────────────────
	type enumClass struct{ name, compType string }
	classes := []enumClass{
		{"Тип двигателя",        "engine"},
		{"Тип батареи",          "battery"},
		{"Стандарт зарядки",     "charger"},
		{"Тип привода",          "emobile"},
		{"Класс защиты IP",      "battery"},
	}
	classIDs := make([]int, len(classes))
	for i, c := range classes {
		if err := db.QueryRowContext(ctx,
			`INSERT INTO enum_class (name, component_type) VALUES($1,$2) RETURNING enum_class_id`,
			c.name, c.compType).Scan(&classIDs[i]); err != nil {
			log.Printf("seed: enum_class '%s': %v", c.name, err)
			return nil
		}
	}

	// Позиции перечислений
	positions := map[int][][2]string{
		classIDs[0]: {{"AC", "0"}, {"DC", "1"}},
		classIDs[1]: {{"Li-ion", "0"}, {"Li-polymer", "1"}},
		classIDs[2]: {{"CCS2", "0"}, {"CHAdeMO", "1"}, {"Type 2", "2"}, {"GB/T", "3"}},
		classIDs[3]: {{"Передний", "0"}, {"Задний", "1"}, {"Полный", "2"}},
		classIDs[4]: {{"IP54", "0"}, {"IP65", "1"}, {"IP67", "2"}, {"IP68", "3"}},
	}
	for classID, vals := range positions {
		for _, v := range vals {
			if _, err := db.ExecContext(ctx,
				`INSERT INTO enum_position (enum_class_id, value, order_num) VALUES($1,$2,$3)`,
				classID, v[0], v[1]); err != nil {
				log.Printf("seed: enum_position: %v", err)
			}
		}
	}
	log.Printf("seed: PR2 enum_class (%d classes) seeded", len(classes))

	// ── ПР3: Параметры ─────────────────────────────────────────────────────
	type param struct{ des, name, ptype, unit string; enumClassID int }
	params := []param{
		{"range_km",   "Запас хода",              "real", "км",   0},
		{"weight_kg",  "Снаряжённая масса",        "int",  "кг",   0},
		{"max_speed",  "Максимальная скорость",    "int",  "км/ч", 0},
		{"drive_type", "Тип привода",              "enum", "",     classIDs[3]},
		{"bat_cap",    "Ёмкость батареи",          "real", "кВтч", 0},
		{"charge_std", "Стандарт зарядки",         "enum", "",     classIDs[2]},
	}
	paramIDs := make([]int, len(params))
	for i, p := range params {
		var enumPtr interface{}
		if p.enumClassID != 0 {
			enumPtr = p.enumClassID
		}
		if err := db.QueryRowContext(ctx,
			`INSERT INTO parameter (designation, name, param_type, measuring_unit, enum_class_id)
			 VALUES($1,$2,$3,NULLIF($4,''),NULLIF($5::text,'')::int) RETURNING parameter_id`,
			p.des, p.name, p.ptype, p.unit, enumPtr).Scan(&paramIDs[i]); err != nil {
			log.Printf("seed: parameter '%s': %v", p.des, err)
			return nil
		}
	}

	// Привязываем параметры к типу 'emobile'
	for i, pID := range paramIDs {
		var minV, maxV interface{}
		switch i {
		case 0: minV, maxV = 50, 800    // range_km
		case 1: minV, maxV = 800, 4000  // weight_kg
		case 2: minV, maxV = 60, 350    // max_speed
		case 4: minV, maxV = 20, 200    // bat_cap
		}
		if _, err := db.ExecContext(ctx,
			`INSERT INTO component_parameter (component_type, parameter_id, order_num, min_val, max_val)
			 VALUES('emobile',$1,$2,$3,$4) ON CONFLICT DO NOTHING`,
			pID, i, minV, maxV); err != nil {
			log.Printf("seed: component_parameter: %v", err)
		}
	}
	log.Printf("seed: PR3 parameters (%d) seeded for emobile", len(params))

	return nil
}

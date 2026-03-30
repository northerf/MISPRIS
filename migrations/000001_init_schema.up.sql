CREATE TABLE users (
                       user_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       username    VARCHAR(100) UNIQUE NOT NULL,
                       password    VARCHAR(255) NOT NULL,
                       role        TEXT NOT NULL CHECK (role IN ('admin', 'user')),
                       is_active   BOOLEAN DEFAULT TRUE,
                       created_at  TIMESTAMP DEFAULT NOW(),
                       updated_at  TIMESTAMP DEFAULT NOW(),
                       deleted_at  BOOLEAN DEFAULT FALSE
);

INSERT INTO users (username, password, role, is_active) VALUES ('daun', '$2b$10$FtD.voJMWrH4qqmjlJLwY.BF8fkw3vs7h8ngil/w063.Ks4UW3SZq', 'admin', true);
INSERT INTO users (username, password, role, is_active) VALUES ('botik', '$2b$10$FtD.voJMWrH4qqmjlJLwY.BF8fkw3vs7h8ngil/w063.Ks4UW3SZq', 'user', true);

CREATE TABLE engine (
                        engine_id    SERIAL PRIMARY KEY,
                        engine_name  VARCHAR(100) NOT NULL,
                        engine_type  VARCHAR(10) NOT NULL CHECK (engine_type IN ('AC', 'DC')),
                        engine_info  JSONB,
                        created_at   TIMESTAMP DEFAULT NOW(),
                        updated_at   TIMESTAMP DEFAULT NOW()
);

CREATE TABLE inverter (
                          inverter_id    SERIAL PRIMARY KEY,
                          inverter_name  VARCHAR(100) NOT NULL,
                          inverter_info  JSONB,
                          created_at     TIMESTAMP DEFAULT NOW(),
                          updated_at     TIMESTAMP DEFAULT NOW()
);

CREATE TABLE gearbox (
                         gearbox_id    SERIAL PRIMARY KEY,
                         gearbox_name  VARCHAR(100) NOT NULL,
                         gearbox_info  JSONB,
                         created_at    TIMESTAMP DEFAULT NOW(),
                         updated_at    TIMESTAMP DEFAULT NOW()
);

CREATE TABLE power_point (
                             power_point_id  SERIAL PRIMARY KEY,
                             engine_id       INT NOT NULL REFERENCES engine(engine_id) ON DELETE RESTRICT,
                             inverter_id     INT NOT NULL REFERENCES inverter(inverter_id) ON DELETE RESTRICT,
                             gearbox_id      INT NOT NULL REFERENCES gearbox(gearbox_id) ON DELETE RESTRICT,
                             created_at      TIMESTAMP DEFAULT NOW(),
                             updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE TABLE battery (
                         battery_id       SERIAL PRIMARY KEY,
                         battery_name     VARCHAR(100) NOT NULL,
                         battery_type     VARCHAR(20) NOT NULL CHECK (battery_type IN ('Li-ion', 'Li-polymer')),
                         battery_capacity DECIMAL(10,2),
                         battery_info     JSONB,
                         created_at       TIMESTAMP DEFAULT NOW(),
                         updated_at       TIMESTAMP DEFAULT NOW()
);

CREATE TABLE charger (
                         charger_id    SERIAL PRIMARY KEY,
                         charger_name  VARCHAR(100) NOT NULL,
                         charger_info  JSONB,
                         created_at    TIMESTAMP DEFAULT NOW(),
                         updated_at    TIMESTAMP DEFAULT NOW()
);

CREATE TABLE connector (
                           connector_id    SERIAL PRIMARY KEY,
                           connector_name  VARCHAR(100) NOT NULL,
                           connector_info  JSONB,
                           created_at      TIMESTAMP DEFAULT NOW(),
                           updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE TABLE charger_system (
                                charger_system_id  SERIAL PRIMARY KEY,
                                charger_id         INT NOT NULL REFERENCES charger(charger_id) ON DELETE RESTRICT,
                                connector_id       INT NOT NULL REFERENCES connector(connector_id) ON DELETE RESTRICT,
                                created_at         TIMESTAMP DEFAULT NOW(),
                                updated_at         TIMESTAMP DEFAULT NOW()
);

CREATE TABLE frame (
                       frame_id    SERIAL PRIMARY KEY,
                       frame_name  VARCHAR(100) NOT NULL,
                       frame_info  JSONB,
                       created_at  TIMESTAMP DEFAULT NOW(),
                       updated_at  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE suspension (
                            suspension_id    SERIAL PRIMARY KEY,
                            suspension_name  VARCHAR(100) NOT NULL,
                            suspension_info  JSONB,
                            created_at       TIMESTAMP DEFAULT NOW(),
                            updated_at       TIMESTAMP DEFAULT NOW()
);

CREATE TABLE break_system (
                              break_system_id    SERIAL PRIMARY KEY,
                              break_system_name  VARCHAR(100) NOT NULL,
                              break_info         JSONB,
                              created_at         TIMESTAMP DEFAULT NOW(),
                              updated_at         TIMESTAMP DEFAULT NOW()
);

CREATE TABLE chassis (
                         chassis_id       SERIAL PRIMARY KEY,
                         frame_id         INT NOT NULL REFERENCES frame(frame_id) ON DELETE RESTRICT,
                         suspension_id    INT NOT NULL REFERENCES suspension(suspension_id) ON DELETE RESTRICT,
                         break_system_id  INT NOT NULL REFERENCES break_system(break_system_id) ON DELETE RESTRICT,
                         created_at       TIMESTAMP DEFAULT NOW(),
                         updated_at       TIMESTAMP DEFAULT NOW()
);

CREATE TABLE carcass (
                         carcass_id    SERIAL PRIMARY KEY,
                         carcass_name  VARCHAR(100) NOT NULL,
                         carcass_info  JSONB,
                         created_at    TIMESTAMP DEFAULT NOW(),
                         updated_at    TIMESTAMP DEFAULT NOW()
);

CREATE TABLE doors (
                       doors_id    SERIAL PRIMARY KEY,
                       doors_name  VARCHAR(100) NOT NULL,
                       doors_info  JSONB,
                       created_at  TIMESTAMP DEFAULT NOW(),
                       updated_at  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE wings (
                       wings_id    SERIAL PRIMARY KEY,
                       wings_name  VARCHAR(100) NOT NULL,
                       wings_info  JSONB,
                       created_at  TIMESTAMP DEFAULT NOW(),
                       updated_at  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE body (
                      body_id      SERIAL PRIMARY KEY,
                      carcass_id   INT NOT NULL REFERENCES carcass(carcass_id) ON DELETE RESTRICT,
                      doors_id     INT NOT NULL REFERENCES doors(doors_id) ON DELETE RESTRICT,
                      wings_id     INT NOT NULL REFERENCES wings(wings_id) ON DELETE RESTRICT,
                      created_at   TIMESTAMP DEFAULT NOW(),
                      updated_at   TIMESTAMP DEFAULT NOW()
);

CREATE TABLE controllers (
                             controller_id    SERIAL PRIMARY KEY,
                             controller_name  VARCHAR(100) NOT NULL,
                             controller_info  JSONB,
                             created_at       TIMESTAMP DEFAULT NOW(),
                             updated_at       TIMESTAMP DEFAULT NOW()
);

CREATE TABLE sensors (
                         sensor_id    SERIAL PRIMARY KEY,
                         sensor_name  VARCHAR(100) NOT NULL,
                         sensor_info  JSONB,
                         created_at   TIMESTAMP DEFAULT NOW(),
                         updated_at   TIMESTAMP DEFAULT NOW()
);

CREATE TABLE wiring (
                        wiring_id    SERIAL PRIMARY KEY,
                        wiring_name  VARCHAR(100) NOT NULL,
                        wiring_info  JSONB,
                        created_at   TIMESTAMP DEFAULT NOW(),
                        updated_at   TIMESTAMP DEFAULT NOW()
);

CREATE TABLE electronics (
                             electronics_id  SERIAL PRIMARY KEY,
                             controller_id   INT NOT NULL REFERENCES controllers(controller_id) ON DELETE RESTRICT,
                             sensor_id       INT NOT NULL REFERENCES sensors(sensor_id) ON DELETE RESTRICT,
                             wiring_id       INT NOT NULL REFERENCES wiring(wiring_id) ON DELETE RESTRICT,
                             created_at      TIMESTAMP DEFAULT NOW(),
                             updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE TABLE emobile (
                         emobile_id         SERIAL PRIMARY KEY,
                         emobile_name       VARCHAR(100) NOT NULL,
                         power_point_id     INT NOT NULL REFERENCES power_point(power_point_id) ON DELETE RESTRICT,
                         battery_id         INT NOT NULL REFERENCES battery(battery_id) ON DELETE RESTRICT,
                         charger_system_id  INT NOT NULL REFERENCES charger_system(charger_system_id) ON DELETE RESTRICT,
                         chassis_id         INT NOT NULL REFERENCES chassis(chassis_id) ON DELETE RESTRICT,
                         body_id            INT NOT NULL REFERENCES body(body_id) ON DELETE RESTRICT,
                         electronics_id     INT NOT NULL REFERENCES electronics(electronics_id) ON DELETE RESTRICT,
                         created_at         TIMESTAMP DEFAULT NOW(),
                         updated_at         TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_engine_info      ON engine     USING GIN (engine_info);
CREATE INDEX idx_battery_info     ON battery    USING GIN (battery_info);
CREATE INDEX idx_charger_info     ON charger    USING GIN (charger_info);

CREATE INDEX idx_users_username ON users(username);

CREATE TABLE audit_log (
                           log_id      SERIAL PRIMARY KEY,
                           user_id     UUID REFERENCES users(user_id) ON DELETE SET NULL,
                           action      VARCHAR(10) NOT NULL CHECK (action IN ('CREATE', 'UPDATE', 'DELETE')),
                           table_name  VARCHAR(100) NOT NULL,
                           record_id   INT NOT NULL,
                           old_data    JSONB,
                           new_data    JSONB,
                           created_at  TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_audit_table    ON audit_log(table_name, record_id);
CREATE INDEX idx_audit_user     ON audit_log(user_id);
CREATE INDEX idx_audit_date     ON audit_log(created_at);


CREATE MATERIALIZED VIEW emobile_full AS
SELECT
    e.emobile_id,
    e.emobile_name,
    e.created_at,

    -- Силовая установка
    eng.engine_name,
    eng.engine_type,
    eng.engine_info,
    inv.inverter_name,
    inv.inverter_info,
    gb.gearbox_name,
    gb.gearbox_info,

    -- Батарея
    b.battery_name,
    b.battery_type,
    b.battery_capacity,
    b.battery_info,

    -- Зарядная система
    ch.charger_name,
    ch.charger_info,
    con.connector_name,
    con.connector_info,

    -- Шасси
    fr.frame_name,
    fr.frame_info,
    sus.suspension_name,
    sus.suspension_info,
    br.break_system_name,
    br.break_info,

    -- Кузов
    car.carcass_name,
    car.carcass_info,
    dr.doors_name,
    dr.doors_info,
    wi.wings_name,
    wi.wings_info,

    -- Электроника
    ctrl.controller_name,
    ctrl.controller_info,
    sen.sensor_name,
    sen.sensor_info,
    wr.wiring_name,
    wr.wiring_info

FROM emobile e
         JOIN power_point    pp   ON e.power_point_id    = pp.power_point_id
         JOIN engine         eng  ON pp.engine_id         = eng.engine_id
         JOIN inverter       inv  ON pp.inverter_id        = inv.inverter_id
         JOIN gearbox        gb   ON pp.gearbox_id         = gb.gearbox_id
         JOIN battery        b    ON e.battery_id          = b.battery_id
         JOIN charger_system cs   ON e.charger_system_id   = cs.charger_system_id
         JOIN charger        ch   ON cs.charger_id         = ch.charger_id
         JOIN connector      con  ON cs.connector_id       = con.connector_id
         JOIN chassis        chs  ON e.chassis_id          = chs.chassis_id
         JOIN frame          fr   ON chs.frame_id          = fr.frame_id
         JOIN suspension     sus  ON chs.suspension_id     = sus.suspension_id
         JOIN break_system   br   ON chs.break_system_id   = br.break_system_id
         JOIN body           bod  ON e.body_id             = bod.body_id
         JOIN carcass        car  ON bod.carcass_id        = car.carcass_id
         JOIN doors          dr   ON bod.doors_id          = dr.doors_id
         JOIN wings          wi   ON bod.wings_id          = wi.wings_id
         JOIN electronics    el   ON e.electronics_id      = el.electronics_id
         JOIN controllers    ctrl ON el.controller_id      = ctrl.controller_id
         JOIN sensors        sen  ON el.sensor_id          = sen.sensor_id
         JOIN wiring         wr   ON el.wiring_id          = wr.wiring_id;

REFRESH MATERIALIZED VIEW emobile_full;

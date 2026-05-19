-- ============================================================
--  Миграция 000002: ПР2 (Перечисления) и ПР3 (Параметры)
-- ============================================================

-- ───────────────────────────────────────────────────────────
--  ПР2: Классы перечислений и их позиции
-- ───────────────────────────────────────────────────────────

-- enum_class — справочник допустимых значений.
-- Например: «Класс защиты IP», «Стандарт зарядки», «Тип кузова».
-- component_type — к какому типу компонентов относится (battery, engine и т.д.)
CREATE TABLE IF NOT EXISTS enum_class (
    enum_class_id   SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL UNIQUE,
    component_type  VARCHAR(50)  DEFAULT NULL,
    created_at      TIMESTAMP    DEFAULT NOW(),
    updated_at      TIMESTAMP    DEFAULT NOW()
);

-- enum_position — одно конкретное значение в справочнике.
-- order_num управляет порядком вывода.
CREATE TABLE IF NOT EXISTS enum_position (
    enum_position_id SERIAL PRIMARY KEY,
    enum_class_id    INT          NOT NULL REFERENCES enum_class(enum_class_id) ON DELETE CASCADE,
    value            VARCHAR(200) NOT NULL,
    order_num        INT          NOT NULL DEFAULT 0,
    created_at       TIMESTAMP    DEFAULT NOW(),
    updated_at       TIMESTAMP    DEFAULT NOW()
);

-- ───────────────────────────────────────────────────────────
--  ПР3: Параметры и их значения для конкретных автомобилей
-- ───────────────────────────────────────────────────────────

-- parameter — описание одной характеристики.
-- Примеры: запас хода (real, км), масса (int, кг), цвет (str), класс защиты (enum).
CREATE TABLE IF NOT EXISTS parameter (
    parameter_id   SERIAL       PRIMARY KEY,
    designation    VARCHAR(50)  NOT NULL UNIQUE,
    name           VARCHAR(200) NOT NULL,
    param_type     VARCHAR(10)  NOT NULL CHECK (param_type IN ('real','int','str','enum')),
    measuring_unit VARCHAR(50)  DEFAULT NULL,
    enum_class_id  INT          DEFAULT NULL REFERENCES enum_class(enum_class_id) ON DELETE SET NULL,
    created_at     TIMESTAMP    DEFAULT NOW(),
    updated_at     TIMESTAMP    DEFAULT NOW()
);

-- component_parameter — привязка параметра к типу компонента с ограничениями.
-- component_type: 'emobile', 'battery', 'engine', ... — любой тип из системы.
-- min_val/max_val — ограничения для числовых параметров.
CREATE TABLE IF NOT EXISTS component_parameter (
    component_parameter_id SERIAL      PRIMARY KEY,
    component_type         VARCHAR(50) NOT NULL,
    parameter_id           INT         NOT NULL REFERENCES parameter(parameter_id) ON DELETE CASCADE,
    order_num              INT         NOT NULL DEFAULT 0,
    min_val                FLOAT       DEFAULT NULL,
    max_val                FLOAT       DEFAULT NULL,
    UNIQUE (component_type, parameter_id)
);

-- emobile_parameter_value — значение параметра для конкретного автомобиля.
-- Только одно из val_* заполняется (зависит от param_type параметра).
CREATE TABLE IF NOT EXISTS emobile_parameter_value (
    value_id               SERIAL PRIMARY KEY,
    emobile_id             INT    NOT NULL REFERENCES emobile(emobile_id)                     ON DELETE CASCADE,
    component_parameter_id INT    NOT NULL REFERENCES component_parameter(component_parameter_id) ON DELETE CASCADE,
    val_real               FLOAT       DEFAULT NULL,
    val_int                INT         DEFAULT NULL,
    val_str                VARCHAR(500) DEFAULT NULL,
    enum_val_id            INT          DEFAULT NULL REFERENCES enum_position(enum_position_id) ON DELETE SET NULL,
    UNIQUE (emobile_id, component_parameter_id)
);

-- Индексы для ускорения частых запросов
CREATE INDEX IF NOT EXISTS idx_enum_position_class ON enum_position(enum_class_id);
CREATE INDEX IF NOT EXISTS idx_comp_param_type     ON component_parameter(component_type);
CREATE INDEX IF NOT EXISTS idx_emobile_param_val   ON emobile_parameter_value(emobile_id);

-- ───────────────────────────────────────────────────────────
--  SQL-функции ПР2
-- ───────────────────────────────────────────────────────────

-- get_enum_values: возвращает значения перечисления в заданном порядке
CREATE OR REPLACE FUNCTION get_enum_values(p_enum_class_id INTEGER)
RETURNS TABLE(enum_position_id INTEGER, value VARCHAR, order_num INTEGER)
LANGUAGE sql AS $$
    SELECT enum_position_id, value, order_num
    FROM enum_position
    WHERE enum_class_id = p_enum_class_id
    ORDER BY order_num;
$$;

-- validate_enum_value: проверяет, существует ли значение в перечислении
CREATE OR REPLACE FUNCTION validate_enum_value(p_enum_class_id INTEGER, p_value VARCHAR)
RETURNS BOOLEAN
LANGUAGE sql AS $$
    SELECT EXISTS(
        SELECT 1 FROM enum_position
        WHERE enum_class_id = p_enum_class_id AND value = p_value
    );
$$;

-- ───────────────────────────────────────────────────────────
--  SQL-функции и процедуры ПР3
-- ───────────────────────────────────────────────────────────

-- get_component_parameters: возвращает параметры типа компонента с типами и ограничениями
CREATE OR REPLACE FUNCTION get_component_parameters(p_component_type VARCHAR)
RETURNS TABLE(
    cp_id          INTEGER,
    param_id       INTEGER,
    designation    VARCHAR,
    name           VARCHAR,
    param_type     VARCHAR,
    measuring_unit VARCHAR,
    min_val        FLOAT,
    max_val        FLOAT,
    order_num      INTEGER
)
LANGUAGE sql AS $$
    SELECT
        cp.component_parameter_id,
        p.parameter_id,
        p.designation,
        p.name,
        p.param_type,
        p.measuring_unit,
        cp.min_val,
        cp.max_val,
        cp.order_num
    FROM component_parameter cp
    JOIN parameter p ON cp.parameter_id = p.parameter_id
    WHERE cp.component_type = p_component_type
    ORDER BY cp.order_num;
$$;

-- copy_component_parameters: копирует параметры от одного типа компонента к другому.
-- Полезно для наследования параметров при создании нового типа.
CREATE OR REPLACE PROCEDURE copy_component_parameters(p_from_type VARCHAR, p_to_type VARCHAR)
LANGUAGE sql AS $$
    INSERT INTO component_parameter (component_type, parameter_id, order_num, min_val, max_val)
    SELECT p_to_type, parameter_id, order_num, min_val, max_val
    FROM component_parameter
    WHERE component_type = p_from_type
    ON CONFLICT (component_type, parameter_id) DO NOTHING;
$$;

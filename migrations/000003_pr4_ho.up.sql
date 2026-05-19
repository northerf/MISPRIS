-- ============================================================
--  Миграция 000003: ПР4 — Хозяйственные операции (ХО)
-- ============================================================

-- ───────────────────────────────────────────────────────────
--  Субъекты хозяйственной деятельности (СХД)
--  Любой участник операций: предприятие, подразделение, склад
-- ───────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS shd (
    shd_id        SERIAL PRIMARY KEY,
    name          VARCHAR(200) NOT NULL,
    shd_type      VARCHAR(30)  NOT NULL DEFAULT 'organization'
                  CHECK (shd_type IN ('organization','individual','department')),
    inn           VARCHAR(20)  DEFAULT NULL,
    description   TEXT         DEFAULT NULL,
    created_at    TIMESTAMP    DEFAULT NOW(),
    updated_at    TIMESTAMP    DEFAULT NOW()
);

-- ───────────────────────────────────────────────────────────
--  Классификатор типов ХО (дерево, как в ПР1)
--  Примеры: Отгрузка, Поступление, Списание, Перемещение
-- ───────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS ho_class (
    ho_class_id   SERIAL PRIMARY KEY,
    name          VARCHAR(200) NOT NULL,
    designation   VARCHAR(50)  DEFAULT NULL,    -- код (напр. TTN, TREB)
    parent_id     INT          DEFAULT NULL REFERENCES ho_class(ho_class_id) ON DELETE SET NULL,
    is_terminal   BOOLEAN      DEFAULT FALSE,   -- конечный тип (не группа)
    created_at    TIMESTAMP    DEFAULT NOW(),
    updated_at    TIMESTAMP    DEFAULT NOW()
);

-- ───────────────────────────────────────────────────────────
--  Классификатор ролей участников
--  Примеры: Грузоотправитель, Грузополучатель, Плательщик
-- ───────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS ho_role (
    ho_role_id    SERIAL PRIMARY KEY,
    name          VARCHAR(100) NOT NULL UNIQUE,
    description   TEXT         DEFAULT NULL
);

-- ───────────────────────────────────────────────────────────
--  Допустимые роли для каждого типа ХО
--  Определяет: у «Отгрузки» могут быть роли X, Y, Z
-- ───────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS ho_class_role (
    id            SERIAL PRIMARY KEY,
    ho_class_id   INT NOT NULL REFERENCES ho_class(ho_class_id) ON DELETE CASCADE,
    ho_role_id    INT NOT NULL REFERENCES ho_role(ho_role_id) ON DELETE CASCADE,
    is_required   BOOLEAN DEFAULT FALSE,
    UNIQUE (ho_class_id, ho_role_id)
);

-- ───────────────────────────────────────────────────────────
--  Параметры типа ХО (переиспользует parameter из ПР3)
--  Пример: у «Отгрузки» есть параметр «Сумма НДС», «Масса груза»
-- ───────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS ho_class_parameter (
    id            SERIAL PRIMARY KEY,
    ho_class_id   INT NOT NULL REFERENCES ho_class(ho_class_id) ON DELETE CASCADE,
    parameter_id  INT NOT NULL REFERENCES parameter(parameter_id) ON DELETE CASCADE,
    order_num     INT     DEFAULT 0,
    min_val       FLOAT   DEFAULT NULL,
    max_val       FLOAT   DEFAULT NULL,
    UNIQUE (ho_class_id, parameter_id)
);

-- ───────────────────────────────────────────────────────────
--  Классификатор типов документов
--  Примеры: ТТН, Счёт-фактура, Договор, Требование-накладная
-- ───────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS document_class (
    doc_class_id  SERIAL PRIMARY KEY,
    name          VARCHAR(200) NOT NULL,
    code          VARCHAR(50)  DEFAULT NULL,   -- код ОКУД
    description   TEXT         DEFAULT NULL,
    created_at    TIMESTAMP    DEFAULT NOW()
);

-- ───────────────────────────────────────────────────────────
--  Допустимые документы для типа ХО
--  Определяет: «Отгрузка» оформляется ТТН (обязательно) + Счётом (необязательно)
-- ───────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS ho_class_document (
    id            SERIAL PRIMARY KEY,
    ho_class_id   INT NOT NULL REFERENCES ho_class(ho_class_id) ON DELETE CASCADE,
    doc_class_id  INT NOT NULL REFERENCES document_class(doc_class_id) ON DELETE CASCADE,
    role_name     VARCHAR(100) DEFAULT NULL,  -- 'Основной', 'Основание', 'Приложение'
    is_required   BOOLEAN DEFAULT FALSE,
    UNIQUE (ho_class_id, doc_class_id)
);

-- ───────────────────────────────────────────────────────────
--  Экземпляр ХО — конкретная хозяйственная операция
-- ───────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS ho_instance (
    ho_id         SERIAL PRIMARY KEY,
    ho_class_id   INT  NOT NULL REFERENCES ho_class(ho_class_id) ON DELETE RESTRICT,
    doc_number    VARCHAR(100) DEFAULT NULL,
    doc_date      DATE         DEFAULT CURRENT_DATE,
    total_amount  DECIMAL(15,2) DEFAULT 0,
    status        VARCHAR(30)  DEFAULT 'draft'
                  CHECK (status IN ('draft','confirmed','cancelled')),
    note          TEXT         DEFAULT NULL,
    created_at    TIMESTAMP    DEFAULT NOW(),
    updated_at    TIMESTAMP    DEFAULT NOW()
);

-- ───────────────────────────────────────────────────────────
--  Назначение СХД на роль в конкретной ХО
-- ───────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS ho_actor (
    id            SERIAL PRIMARY KEY,
    ho_id         INT NOT NULL REFERENCES ho_instance(ho_id) ON DELETE CASCADE,
    ho_role_id    INT NOT NULL REFERENCES ho_role(ho_role_id) ON DELETE RESTRICT,
    shd_id        INT NOT NULL REFERENCES shd(shd_id) ON DELETE RESTRICT,
    UNIQUE (ho_id, ho_role_id)
);

-- ───────────────────────────────────────────────────────────
--  Значения параметров конкретной ХО
-- ───────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS ho_parameter_value (
    id                     SERIAL PRIMARY KEY,
    ho_id                  INT NOT NULL REFERENCES ho_instance(ho_id) ON DELETE CASCADE,
    ho_class_parameter_id  INT NOT NULL REFERENCES ho_class_parameter(id) ON DELETE CASCADE,
    val_real               FLOAT        DEFAULT NULL,
    val_int                INT          DEFAULT NULL,
    val_str                VARCHAR(500) DEFAULT NULL,
    val_date               DATE         DEFAULT NULL,
    enum_val_id            INT          DEFAULT NULL REFERENCES enum_position(enum_position_id) ON DELETE SET NULL,
    UNIQUE (ho_id, ho_class_parameter_id)
);

-- ───────────────────────────────────────────────────────────
--  Документы, приложенные к конкретной ХО
-- ───────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS ho_document (
    id            SERIAL PRIMARY KEY,
    ho_id         INT NOT NULL REFERENCES ho_instance(ho_id) ON DELETE CASCADE,
    doc_class_id  INT NOT NULL REFERENCES document_class(doc_class_id) ON DELETE RESTRICT,
    doc_number    VARCHAR(100) DEFAULT NULL,
    doc_date      DATE         DEFAULT NULL,
    note          TEXT         DEFAULT NULL
);

-- ───────────────────────────────────────────────────────────
--  Позиции ХО — что именно движется (изделия, количество)
-- ───────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS ho_position (
    id            SERIAL PRIMARY KEY,
    ho_id         INT NOT NULL REFERENCES ho_instance(ho_id) ON DELETE CASCADE,
    emobile_id    INT DEFAULT NULL REFERENCES emobile(emobile_id) ON DELETE RESTRICT,
    quantity      INT           NOT NULL DEFAULT 1,
    unit_price    DECIMAL(15,2) DEFAULT 0,
    total_price   DECIMAL(15,2) DEFAULT 0,
    note          TEXT          DEFAULT NULL,
    position_num  INT           DEFAULT 0
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_ho_instance_class ON ho_instance(ho_class_id);
CREATE INDEX IF NOT EXISTS idx_ho_actor_ho       ON ho_actor(ho_id);
CREATE INDEX IF NOT EXISTS idx_ho_param_val_ho   ON ho_parameter_value(ho_id);
CREATE INDEX IF NOT EXISTS idx_ho_position_ho    ON ho_position(ho_id);

-- ============================================================
--  SQL-функции и процедуры ПР4
-- ============================================================

-- Создать тип ХО в классификаторе
CREATE OR REPLACE PROCEDURE ins_ho_class(
    p_name        VARCHAR,
    p_designation VARCHAR,
    p_parent_id   INT,
    p_is_terminal BOOLEAN
)
LANGUAGE sql AS $$
    INSERT INTO ho_class (name, designation, parent_id, is_terminal, created_at, updated_at)
    VALUES (p_name, NULLIF(p_designation,''), p_parent_id, p_is_terminal, NOW(), NOW());
$$;

-- Добавить допустимую роль к типу ХО
CREATE OR REPLACE PROCEDURE add_role_to_ho_class(
    p_ho_class_id INT,
    p_ho_role_id  INT,
    p_is_required BOOLEAN
)
LANGUAGE sql AS $$
    INSERT INTO ho_class_role (ho_class_id, ho_role_id, is_required)
    VALUES (p_ho_class_id, p_ho_role_id, p_is_required)
    ON CONFLICT (ho_class_id, ho_role_id) DO UPDATE SET is_required = EXCLUDED.is_required;
$$;

-- Добавить параметр к типу ХО (использует table parameter из ПР3)
CREATE OR REPLACE PROCEDURE add_param_to_ho_class(
    p_ho_class_id  INT,
    p_parameter_id INT,
    p_order_num    INT,
    p_min_val      FLOAT,
    p_max_val      FLOAT
)
LANGUAGE sql AS $$
    INSERT INTO ho_class_parameter (ho_class_id, parameter_id, order_num, min_val, max_val)
    VALUES (p_ho_class_id, p_parameter_id, p_order_num, NULLIF(p_min_val, 0), NULLIF(p_max_val, 0))
    ON CONFLICT (ho_class_id, parameter_id) DO UPDATE
        SET order_num = EXCLUDED.order_num,
            min_val   = EXCLUDED.min_val,
            max_val   = EXCLUDED.max_val;
$$;

-- Создать экземпляр ХО и вернуть его ID
CREATE OR REPLACE FUNCTION ins_ho(
    p_ho_class_id  INT,
    p_doc_number   VARCHAR,
    p_doc_date     DATE,
    p_total_amount DECIMAL,
    p_note         VARCHAR DEFAULT NULL
)
RETURNS INTEGER
LANGUAGE plpgsql AS $$
DECLARE
    v_ho_id INT;
BEGIN
    INSERT INTO ho_instance (ho_class_id, doc_number, doc_date, total_amount, note, created_at, updated_at)
    VALUES (p_ho_class_id, NULLIF(p_doc_number,''), p_doc_date, COALESCE(p_total_amount,0), p_note, NOW(), NOW())
    RETURNING ho_id INTO v_ho_id;
    RETURN v_ho_id;
END;
$$;

-- Назначить СХД на роль в ХО (с проверкой допустимости роли)
CREATE OR REPLACE PROCEDURE set_ho_actor(
    p_ho_id      INT,
    p_ho_role_id INT,
    p_shd_id     INT
)
LANGUAGE plpgsql AS $$
DECLARE
    v_class_id INT;
BEGIN
    -- Определяем тип ХО
    SELECT ho_class_id INTO v_class_id FROM ho_instance WHERE ho_id = p_ho_id;
    -- Проверяем, допустима ли роль для этого типа
    IF NOT EXISTS (
        SELECT 1 FROM ho_class_role
        WHERE ho_class_id = v_class_id AND ho_role_id = p_ho_role_id
    ) THEN
        RAISE EXCEPTION 'Роль % недопустима для данного типа ХО', p_ho_role_id;
    END IF;
    INSERT INTO ho_actor (ho_id, ho_role_id, shd_id)
    VALUES (p_ho_id, p_ho_role_id, p_shd_id)
    ON CONFLICT (ho_id, ho_role_id) DO UPDATE SET shd_id = EXCLUDED.shd_id;
END;
$$;

-- Записать значение параметра ХО (с валидацией типа и диапазона)
CREATE OR REPLACE PROCEDURE write_ho_par(
    p_ho_id                 INT,
    p_ho_class_parameter_id INT,
    p_val_real              FLOAT   DEFAULT NULL,
    p_val_int               INT     DEFAULT NULL,
    p_val_str               VARCHAR DEFAULT NULL,
    p_val_date              DATE    DEFAULT NULL,
    p_enum_val_id           INT     DEFAULT NULL
)
LANGUAGE plpgsql AS $$
DECLARE
    v_param_type VARCHAR;
    v_min_val    FLOAT;
    v_max_val    FLOAT;
    v_enum_class INT;
BEGIN
    -- Получаем тип параметра и ограничения
    SELECT p.param_type, hcp.min_val, hcp.max_val, p.enum_class_id
    INTO v_param_type, v_min_val, v_max_val, v_enum_class
    FROM ho_class_parameter hcp
    JOIN parameter p ON hcp.parameter_id = p.parameter_id
    WHERE hcp.id = p_ho_class_parameter_id;

    -- Валидация
    IF v_param_type = 'real' THEN
        IF p_val_real IS NULL THEN
            RAISE EXCEPTION 'val_real обязателен для вещественного параметра';
        END IF;
        IF v_min_val IS NOT NULL AND p_val_real < v_min_val THEN
            RAISE EXCEPTION 'Значение % меньше минимума %', p_val_real, v_min_val;
        END IF;
        IF v_max_val IS NOT NULL AND p_val_real > v_max_val THEN
            RAISE EXCEPTION 'Значение % больше максимума %', p_val_real, v_max_val;
        END IF;
    ELSIF v_param_type = 'int' THEN
        IF p_val_int IS NULL THEN
            RAISE EXCEPTION 'val_int обязателен для целочисленного параметра';
        END IF;
        IF v_min_val IS NOT NULL AND p_val_int < v_min_val THEN
            RAISE EXCEPTION 'Значение % меньше минимума %', p_val_int, v_min_val;
        END IF;
        IF v_max_val IS NOT NULL AND p_val_int > v_max_val THEN
            RAISE EXCEPTION 'Значение % больше максимума %', p_val_int, v_max_val;
        END IF;
    ELSIF v_param_type = 'enum' THEN
        IF p_enum_val_id IS NULL THEN
            RAISE EXCEPTION 'enum_val_id обязателен для параметра-перечисления';
        END IF;
        IF NOT EXISTS (
            SELECT 1 FROM enum_position ep
            WHERE ep.enum_position_id = p_enum_val_id AND ep.enum_class_id = v_enum_class
        ) THEN
            RAISE EXCEPTION 'Значение % не входит в перечисление %', p_enum_val_id, v_enum_class;
        END IF;
    ELSIF v_param_type = 'str' THEN
        IF p_val_str IS NULL THEN
            RAISE EXCEPTION 'val_str обязателен для строкового параметра';
        END IF;
    END IF;

    INSERT INTO ho_parameter_value
        (ho_id, ho_class_parameter_id, val_real, val_int, val_str, val_date, enum_val_id)
    VALUES
        (p_ho_id, p_ho_class_parameter_id, p_val_real, p_val_int, p_val_str, p_val_date, p_enum_val_id)
    ON CONFLICT (ho_id, ho_class_parameter_id) DO UPDATE
        SET val_real    = EXCLUDED.val_real,
            val_int     = EXCLUDED.val_int,
            val_str     = EXCLUDED.val_str,
            val_date    = EXCLUDED.val_date,
            enum_val_id = EXCLUDED.enum_val_id;
END;
$$;

-- Добавить позицию (изделие) в ХО и пересчитать total_amount
CREATE OR REPLACE PROCEDURE add_ho_position(
    p_ho_id       INT,
    p_emobile_id  INT,
    p_quantity    INT,
    p_unit_price  DECIMAL
)
LANGUAGE plpgsql AS $$
DECLARE
    v_pos_num INT;
BEGIN
    SELECT COALESCE(MAX(position_num), 0) + 1 INTO v_pos_num
    FROM ho_position WHERE ho_id = p_ho_id;

    INSERT INTO ho_position (ho_id, emobile_id, quantity, unit_price, total_price, position_num)
    VALUES (p_ho_id, p_emobile_id, p_quantity, p_unit_price, p_quantity * p_unit_price, v_pos_num);

    -- Пересчитываем итоговую сумму ХО
    UPDATE ho_instance
    SET total_amount = (SELECT COALESCE(SUM(total_price), 0) FROM ho_position WHERE ho_id = p_ho_id),
        updated_at   = NOW()
    WHERE ho_id = p_ho_id;
END;
$$;

-- Найти все ХО заданного типа со сводной информацией
CREATE OR REPLACE FUNCTION find_ho_by_class(p_ho_class_id INT)
RETURNS TABLE(
    ho_id        INTEGER,
    doc_number   VARCHAR,
    doc_date     DATE,
    total_amount DECIMAL,
    status       VARCHAR,
    positions    BIGINT,
    actors       TEXT
)
LANGUAGE sql AS $$
    SELECT
        hi.ho_id,
        hi.doc_number,
        hi.doc_date,
        hi.total_amount,
        hi.status,
        COUNT(DISTINCT hp.id)  AS positions,
        STRING_AGG(DISTINCT shd.name || ' (' || hr.name || ')', ', ') AS actors
    FROM ho_instance hi
    LEFT JOIN ho_position hp  ON hp.ho_id = hi.ho_id
    LEFT JOIN ho_actor    ha  ON ha.ho_id = hi.ho_id
    LEFT JOIN shd             ON shd.shd_id = ha.shd_id
    LEFT JOIN ho_role     hr  ON hr.ho_role_id = ha.ho_role_id
    WHERE hi.ho_class_id = p_ho_class_id
    GROUP BY hi.ho_id
    ORDER BY hi.doc_date DESC, hi.ho_id DESC;
$$;

-- Получить все параметры для типа ХО
CREATE OR REPLACE FUNCTION get_ho_class_parameters(p_ho_class_id INT)
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
        hcp.id,
        p.parameter_id,
        p.designation,
        p.name,
        p.param_type,
        p.measuring_unit,
        hcp.min_val,
        hcp.max_val,
        hcp.order_num
    FROM ho_class_parameter hcp
    JOIN parameter p ON hcp.parameter_id = p.parameter_id
    WHERE hcp.ho_class_id = p_ho_class_id
    ORDER BY hcp.order_num;
$$;

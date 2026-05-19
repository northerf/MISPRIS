# MISPRIS — ПР2, ПР3, ПР4: Перечисления, Параметры изделий, Хозяйственные операции

> **Проект:** MISPRIS — Go-бэкенд каталога электромобилей (emobile)  
> **Стек:** Go (Gin), PostgreSQL, миграции golang-migrate  
> **Уже реализовано:** ПР1 (дерево классификаторов), ПР2 (перечисления), ПР3 (параметры изделий)  
> **Проектируется:** ПР4 (хозяйственные операции)

---

## Содержание

1. [ПР2 — Перечисления (Enumerations)](#пр2--перечисления)
   - 1.1 [Концептуальное объяснение: зачем нужны перечисления?](#11-концептуальное-объяснение-зачем-нужны-перечисления)
   - 1.2 [Что реально меняется в системе](#12-что-реально-меняется-в-системе)
   - 1.3 [SQL-функции ПР2](#13-sql-функции-пр2)
   - 1.4 [Практические примеры в домене электромобилей](#14-практические-примеры-в-домене-электромобилей)
2. [ПР3 — Справочник изделий с параметрами](#пр3--справочник-изделий-с-параметрами)
   - 2.1 [Что это и зачем нужно](#21-что-это-и-зачем-нужно)
   - 2.2 [Трёхслойная модель параметров](#22-трёхслойная-модель-параметров)
   - 2.3 [Механизм наследования: copy_component_parameters](#23-механизм-наследования-copy_component_parameters)
   - 2.4 [Логика валидации при записи значений](#24-логика-валидации-при-записи-значений)
   - 2.5 [Практический пример: параметр «Запас хода»](#25-практический-пример-параметр-запас-хода)
3. [ПР4 — Хозяйственные операции](#пр4--хозяйственные-операции)
   - 3.1 [Концептуальное введение](#31-концептуальное-введение)
   - 3.2 [Схема базы данных ПР4](#32-схема-базы-данных-пр4)
   - 3.3 [SQL-процедуры ПР4](#33-sql-процедуры-пр4)
   - 3.4 [Связь ПР4 с ПР1, ПР2, ПР3](#34-связь-пр4-с-пр1-пр2-пр3)
   - 3.5 [Практический пример: «Отгрузка электромобиля покупателю»](#35-практический-пример-отгрузка-электромобиля-покупателю)
   - 3.6 [Диаграмма потока данных (ASCII)](#36-диаграмма-потока-данных-ascii)
4. [Итоговая диаграмма: как ПР1–ПР4 связаны вместе](#итоговая-диаграмма-как-пр1пр4-связаны-вместе)

---

## ПР2 — Перечисления

### 1.1 Концептуальное объяснение: зачем нужны перечисления?

#### Проблема: «жёстко зашитые» значения

Представьте, что при проектировании каталога электромобилей разработчик сделал в БД колонку `engine_type` с ограничением:

```sql
-- Подход БЕЗ ПР2 — «hardcoded» значения
ALTER TABLE engine
    ADD COLUMN engine_type VARCHAR(10)
    CHECK (engine_type IN ('AC', 'DC'));
```

А во фронтенде (JavaScript) написал:

```javascript
// Тоже "жёстко зашито"
const engineTypes = ['AC', 'DC'];
```

Всё работает — пока в мире не появятся гибридные двигатели. Когда бизнес говорит: «Нам нужен тип `HYBRID`», начинается головная боль:

1. Нужно написать SQL-миграцию, чтобы изменить `CHECK` constraint.
2. Нужно изменить код бэкенда (валидацию).
3. Нужно изменить код фронтенда (массив `engineTypes`).
4. Нужно задеплоить все три изменения одновременно — иначе будут баги.

**Итого: одно бизнес-изменение = три технических изменения + деплой.**

#### Решение ПР2: мета-модель перечислений

ПР2 хранит все допустимые значения прямо в базе данных, в таблицах `enum_class` и `enum_position`:

```
enum_class                 enum_position
─────────────────          ─────────────────────────────
enum_class_id = 1          enum_position_id = 1
name = "Тип двигателя"     enum_class_id = 1, value = "AC",     order_num = 1
                           enum_position_id = 2
                           enum_class_id = 1, value = "DC",     order_num = 2
```

Теперь, чтобы добавить тип `Гибрид`, администратору достаточно выполнить одну строку:

```sql
INSERT INTO enum_position (enum_class_id, value, order_num)
VALUES (1, 'Гибрид', 3);
```

**Ноль изменений в коде. Ноль деплоев. Всё работает немедленно.**

#### Конкретный поток данных (Admin → DB → Frontend → User)

```
1. Администратор создаёт enum_class "Тип двигателя"
        ↓ POST /api/enum-class/create
2. Добавляет позиции: AC, DC, Гибрид
        ↓ POST /api/enum-position/create (3 раза)
3. Фронтенд при загрузке формы запрашивает:
        GET /api/enum-class/values{id=1}
        ← [{"id":1,"value":"AC"}, {"id":2,"value":"DC"}, {"id":3,"value":"Гибрид"}]
4. Фронтенд рисует выпадающий список <select> из этих данных
5. Пользователь выбирает "Гибрид"
6. Фронтенд отправляет enum_position_id = 3
7. Бэкенд сохраняет FK: engine.engine_type_id = 3
```

---

### 1.2 Что реально меняется в системе

#### До ПР2 — статические ограничения

| Компонент     | Что было                                                    |
|---------------|-------------------------------------------------------------|
| БД            | `CHECK (engine_type IN ('AC','DC'))` в DDL                  |
| Бэкенд        | `if engineType != "AC" && engineType != "DC" { error }`     |
| Фронтенд      | `<option value="AC">AC</option>` вшито в HTML               |
| Изменение     | Требует миграции БД + изменения кода + деплоя               |

#### После ПР2 — динамические данные

| Компонент     | Что стало                                                   |
|---------------|-------------------------------------------------------------|
| БД            | FK на `enum_position.enum_position_id`                      |
| Бэкенд        | Вызов `validate_enum_value(class_id, value)` перед сохранением |
| Фронтенд      | `GET /api/enum-class/values{id}` → динамический `<select>` |
| Изменение     | `INSERT INTO enum_position` — только данные, без кода       |

#### Структура таблиц

**`enum_class`** — «словарь словарей», каждая строка описывает один тип перечисления:

```sql
CREATE TABLE enum_class (
    enum_class_id   SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL UNIQUE,  -- "Тип двигателя", "Тип батареи"
    component_type  VARCHAR(50)  DEFAULT NULL,     -- к какому типу привязан: 'engine', 'battery'
    created_at      TIMESTAMP    DEFAULT NOW(),
    updated_at      TIMESTAMP    DEFAULT NOW()
);
```

| Колонка          | Назначение                                                              |
|------------------|-------------------------------------------------------------------------|
| `enum_class_id`  | Первичный ключ, по нему фронтенд запрашивает значения                   |
| `name`           | Человекочитаемое название справочника                                   |
| `component_type` | Опциональная привязка к типу компонента (для фильтрации в UI)           |
| `created_at`     | Аудит — когда создан                                                    |

**`enum_position`** — конкретные допустимые значения:

```sql
CREATE TABLE enum_position (
    enum_position_id SERIAL PRIMARY KEY,
    enum_class_id    INT          NOT NULL REFERENCES enum_class(enum_class_id) ON DELETE CASCADE,
    value            VARCHAR(200) NOT NULL,   -- "AC", "DC", "Гибрид"
    order_num        INT          NOT NULL DEFAULT 0,  -- порядок в выпадающем списке
    created_at       TIMESTAMP    DEFAULT NOW(),
    updated_at       TIMESTAMP    DEFAULT NOW()
);
```

| Колонка           | Назначение                                                              |
|-------------------|-------------------------------------------------------------------------|
| `enum_position_id`| PK — это и есть значение, которое хранится в FK других таблиц          |
| `enum_class_id`   | К какому справочнику относится                                          |
| `value`           | Текст значения, отображаемый пользователю                               |
| `order_num`       | **UX-важно:** значения с `order_num=0` идут первыми. Самые популярные типы батарей (например, Li-ion) ставятся выше редких — пользователь быстрее выбирает нужное |

**Пример использования `order_num`:** для справочника «Стандарт зарядки» в России чаще используется CCS2, поэтому:

```sql
INSERT INTO enum_position (enum_class_id, value, order_num) VALUES
(3, 'CCS2',    1),   -- первый — самый популярный
(3, 'Type 2',  2),
(3, 'CHAdeMO', 3),
(3, 'GB/T',    4);   -- последний — редкий
```

---

### 1.3 SQL-функции ПР2

Обе функции уже реализованы в миграции `000002_pr2_pr3.up.sql`.

#### `get_enum_values(enum_class_id)` — получить значения для выпадающего списка

```sql
CREATE OR REPLACE FUNCTION get_enum_values(p_enum_class_id INTEGER)
RETURNS TABLE(enum_position_id INTEGER, value VARCHAR, order_num INTEGER)
LANGUAGE sql AS $$
    SELECT enum_position_id, value, order_num
    FROM enum_position
    WHERE enum_class_id = p_enum_class_id
    ORDER BY order_num;
$$;
```

**Что делает:** возвращает все допустимые значения одного справочника, отсортированные по `order_num`.

**Зачем нужна как отдельная функция:** бэкенд вызывает одну хранимую функцию вместо SQL-запроса в коде Go. Это:
- скрывает логику сортировки от приложения;
- позволяет изменить логику (например, добавить фильтрацию неактивных значений) без изменения кода Go.

**Пример вызова:**
```sql
SELECT * FROM get_enum_values(1);
-- Результат:
-- enum_position_id | value  | order_num
-- 1                | AC     | 1
-- 2                | DC     | 2
-- 3                | Гибрид | 3
```

**Соответствующий HTTP-эндпоинт:** `GET /api/enum-class/values{id}` — вызывается фронтендом при загрузке любой формы с выпадающим списком.

---

#### `validate_enum_value(enum_class_id, value)` — проверить допустимость значения

```sql
CREATE OR REPLACE FUNCTION validate_enum_value(p_enum_class_id INTEGER, p_value VARCHAR)
RETURNS BOOLEAN
LANGUAGE sql AS $$
    SELECT EXISTS(
        SELECT 1 FROM enum_position
        WHERE enum_class_id = p_enum_class_id AND value = p_value
    );
$$;
```

**Что делает:** возвращает `TRUE`, если значение `p_value` существует в справочнике `p_enum_class_id`, иначе `FALSE`.

**Зачем нужна:** защита от некорректных данных на уровне БД. Даже если кто-то отправит на бэкенд POST-запрос с произвольным значением (не через UI), функция отклонит его:

```sql
SELECT validate_enum_value(1, 'Дизель');  -- FALSE — нет в справочнике
SELECT validate_enum_value(1, 'AC');      -- TRUE
```

**Соответствующий HTTP-эндпоинт:** `POST /api/enum-class/validate` — вызывается бэкендом перед сохранением любого параметра типа `enum`.

---

### 1.4 Практические примеры в домене электромобилей

| Справочник (enum_class)  | Значения (enum_position)                  | Где используется                              |
|--------------------------|-------------------------------------------|-----------------------------------------------|
| Тип двигателя            | AC, DC, Гибрид                            | `engine.engine_type_id` (FK)                  |
| Тип батареи              | Li-ion, Li-polymer, Твердотельный         | `battery.battery_type_id` (FK)                |
| Стандарт зарядки         | CCS2, CHAdeMO, Type 2, GB/T               | `connector.standard_id` (FK)                  |
| Тип привода              | Передний, Задний, Полный                  | Параметр изделия типа `emobile` (ПР3)         |
| Класс защиты IP          | IP54, IP65, IP67, IP68                    | Параметр изделия типа `battery` (ПР3)         |
| Способ доставки          | Автотранспорт, Ж/Д, Самовывоз, Авиа      | Параметр хозяйственной операции (ПР4)         |
| Статус операции ХО       | Черновик, Подтверждена, Отменена          | `ho_instance.status` (ПР4)                    |

**SQL-пример создания справочника «Тип батареи» с позициями:**

```sql
-- Шаг 1: создать справочник
INSERT INTO enum_class (name, component_type)
VALUES ('Тип батареи', 'battery')
RETURNING enum_class_id;  -- предположим, вернул id=2

-- Шаг 2: добавить значения (Li-ion — самый распространённый, поэтому order_num=1)
INSERT INTO enum_position (enum_class_id, value, order_num) VALUES
(2, 'Li-ion',        1),
(2, 'Li-polymer',    2),
(2, 'Твердотельный', 3);

-- Шаг 3: проверить результат
SELECT * FROM get_enum_values(2);
-- enum_position_id | value         | order_num
-- 4                | Li-ion        | 1
-- 5                | Li-polymer    | 2
-- 6                | Твердотельный | 3
```

---

## ПР3 — Справочник изделий с параметрами

### 2.1 Что это и зачем нужно

ПР2 отвечает на вопрос «какие значения допустимы в поле?» (откуда брать варианты для выпадающего списка).  
ПР3 отвечает на вопрос «**какие поля вообще существуют** для этого типа изделия?»

**Проблема без ПР3:**

У электромобиля, батареи и двигателя — разные характеристики:

| Компонент    | Характеристики                                              |
|--------------|-------------------------------------------------------------|
| Электромобиль| Запас хода (км), Масса (кг), Тип привода                    |
| Батарея      | Ёмкость (кВт·ч), Напряжение (В), Масса (кг), Класс IP      |
| Двигатель    | Мощность (кВт), Момент (Н·м), Макс. обороты (об/мин)       |

Без мета-модели разработчик создаёт отдельные таблицы с жёстко заданными колонками для каждого типа. Когда нужно добавить новую характеристику — снова миграция + изменение кода.

**Решение ПР3:** параметры — тоже строки в таблице. Добавить новую характеристику = `INSERT INTO parameter`. Привязать к типу изделия = `INSERT INTO component_parameter`. Никакого изменения схемы.

---

### 2.2 Трёхслойная модель параметров

```
Слой 1: parameter              — ЧТО это за характеристика (шаблон)
Слой 2: component_parameter    — КОМУ эта характеристика принадлежит (привязка)
Слой 3: emobile_parameter_value — КОНКРЕТНОЕ значение для конкретного изделия
```

#### Слой 1: `parameter` — шаблон характеристики

```sql
CREATE TABLE parameter (
    parameter_id   SERIAL       PRIMARY KEY,
    designation    VARCHAR(50)  NOT NULL UNIQUE,   -- краткий код: 'range_km', 'weight_kg'
    name           VARCHAR(200) NOT NULL,          -- "Запас хода", "Масса снаряжённая"
    param_type     VARCHAR(10)  NOT NULL CHECK (param_type IN ('real','int','str','enum')),
    measuring_unit VARCHAR(50)  DEFAULT NULL,      -- "км", "кг", "кВт·ч", NULL для enum/str
    enum_class_id  INT          DEFAULT NULL REFERENCES enum_class(enum_class_id) ON DELETE SET NULL,
    -- ^ заполняется только если param_type = 'enum', указывает на справочник ПР2
    created_at     TIMESTAMP    DEFAULT NOW(),
    updated_at     TIMESTAMP    DEFAULT NOW()
);
```

Пример: параметр «Запас хода» существует один раз в таблице `parameter`, а не дублируется для каждой модели авто.

#### Слой 2: `component_parameter` — привязка к типу компонента

```sql
CREATE TABLE component_parameter (
    component_parameter_id SERIAL      PRIMARY KEY,
    component_type         VARCHAR(50) NOT NULL,       -- 'emobile', 'battery', 'engine'
    parameter_id           INT         NOT NULL REFERENCES parameter(parameter_id) ON DELETE CASCADE,
    order_num              INT         NOT NULL DEFAULT 0,  -- порядок в форме
    min_val                FLOAT       DEFAULT NULL,    -- минимум для real/int
    max_val                FLOAT       DEFAULT NULL,    -- максимум для real/int
    UNIQUE (component_type, parameter_id)              -- нельзя привязать параметр дважды
);
```

Пример: «Запас хода» привязан к `component_type='emobile'` с `min_val=50, max_val=800`.  
Для `battery` этого параметра нет — только `component_type='battery'` записи.

#### Слой 3: `emobile_parameter_value` — конкретное значение

```sql
CREATE TABLE emobile_parameter_value (
    value_id               SERIAL PRIMARY KEY,
    emobile_id             INT    NOT NULL REFERENCES emobile(emobile_id)      ON DELETE CASCADE,
    component_parameter_id INT    NOT NULL REFERENCES component_parameter(component_parameter_id) ON DELETE CASCADE,
    val_real               FLOAT       DEFAULT NULL,  -- для param_type='real'
    val_int                INT         DEFAULT NULL,  -- для param_type='int'
    val_str                VARCHAR(500) DEFAULT NULL, -- для param_type='str'
    enum_val_id            INT          DEFAULT NULL REFERENCES enum_position(enum_position_id) ON DELETE SET NULL,
    -- ^ для param_type='enum' — FK на ПР2
    UNIQUE (emobile_id, component_parameter_id)       -- одно значение на параметр на авто
);
```

**Только одно из полей `val_*` заполнено** — то, которое соответствует `param_type` параметра.

---

### 2.3 Механизм наследования: `copy_component_parameters`

Предположим, вы создаёте новый тип `emobile_sport` (спортивный электромобиль). Он должен иметь все те же параметры, что и обычный `emobile`, плюс дополнительные.

Без механизма наследования пришлось бы вручную переписывать все `INSERT INTO component_parameter` для нового типа.

С `copy_component_parameters` — одна команда:

```sql
CREATE OR REPLACE PROCEDURE copy_component_parameters(p_from_type VARCHAR, p_to_type VARCHAR)
LANGUAGE sql AS $$
    INSERT INTO component_parameter (component_type, parameter_id, order_num, min_val, max_val)
    SELECT p_to_type, parameter_id, order_num, min_val, max_val
    FROM component_parameter
    WHERE component_type = p_from_type
    ON CONFLICT (component_type, parameter_id) DO NOTHING;
$$;
```

**Пример использования:**

```sql
-- Скопировать все параметры emobile → emobile_sport
CALL copy_component_parameters('emobile', 'emobile_sport');

-- Добавить специфичный для спорт-версии параметр (время разгона 0-100)
INSERT INTO parameter (designation, name, param_type, measuring_unit)
VALUES ('accel_0_100', 'Разгон 0-100 км/ч', 'real', 'с');

INSERT INTO component_parameter (component_type, parameter_id, min_val, max_val)
VALUES ('emobile_sport', <новый parameter_id>, 2.0, 10.0);
```

**Результат:** `emobile_sport` автоматически получил запас хода, массу, тип привода и всё остальное из `emobile` — без ручного перечисления. Нужно было добавить только уникальный параметр.

---

### 2.4 Логика валидации при записи значений

Перед сохранением `emobile_parameter_value` бэкенд (или хранимая процедура) обязан проверить:

```
Если param_type = 'real' или 'int':
    Если min_val задан → value >= min_val
    Если max_val задан → value <= max_val
    Иначе → RAISE EXCEPTION 'Значение вне допустимого диапазона'

Если param_type = 'enum':
    SELECT validate_enum_value(parameter.enum_class_id, значение)
    Если FALSE → RAISE EXCEPTION 'Значение не найдено в справочнике'

Если param_type = 'str':
    Сохраняем как есть — ограничений нет
```

**Пример (псевдокод на Go, аналогичный текущей реализации в сервисном слое):**

```go
func (s *EmobileParameterValueService) Validate(cp ComponentParameter, p Parameter, value interface{}) error {
    switch p.ParamType {
    case "real":
        v := value.(float64)
        if cp.MinVal != nil && v < *cp.MinVal {
            return fmt.Errorf("значение %.2f < допустимого минимума %.2f", v, *cp.MinVal)
        }
        if cp.MaxVal != nil && v > *cp.MaxVal {
            return fmt.Errorf("значение %.2f > допустимого максимума %.2f", v, *cp.MaxVal)
        }
    case "enum":
        valid, _ := s.repo.ValidateEnumValue(p.EnumClassID, value.(string))
        if !valid {
            return fmt.Errorf("значение '%s' не найдено в справочнике #%d", value, p.EnumClassID)
        }
    }
    return nil
}
```

---

### 2.5 Практический пример: параметр «Запас хода»

Полный сценарий от создания параметра до сохранения значения:

**Шаг 1: создать параметр-шаблон**

```sql
INSERT INTO parameter (designation, name, param_type, measuring_unit)
VALUES ('range_km', 'Запас хода', 'real', 'км')
RETURNING parameter_id;  -- например, вернул 7
```

**Шаг 2: привязать к типу компонента `emobile` с ограничениями**

```sql
INSERT INTO component_parameter (component_type, parameter_id, order_num, min_val, max_val)
VALUES ('emobile', 7, 1, 50.0, 800.0);
-- order_num=1 → отображается первым в форме (запас хода — ключевая характеристика EV)
-- min_val=50   → менее 50 км — это не серийный EV
-- max_val=800  → более 800 км — физически нереалистично при текущих технологиях
```

**Шаг 3: установить значение для конкретного электромобиля #1**

```sql
-- Находим ID привязки (component_parameter_id) для emobile + range_km
SELECT component_parameter_id FROM component_parameter
WHERE component_type = 'emobile' AND parameter_id = 7;
-- предположим, вернул component_parameter_id = 12

-- Записываем значение
INSERT INTO emobile_parameter_value (emobile_id, component_parameter_id, val_real)
VALUES (1, 12, 580.0);
-- 580 км — значение для конкретной модели Tesla-подобного EV
```

**Шаг 4: прочитать все параметры электромобиля #1**

```sql
SELECT
    p.name          AS параметр,
    p.measuring_unit AS единица,
    COALESCE(
        epv.val_real::TEXT,
        epv.val_int::TEXT,
        epv.val_str,
        ep.value
    ) AS значение
FROM emobile_parameter_value epv
JOIN component_parameter cp ON cp.component_parameter_id = epv.component_parameter_id
JOIN parameter p             ON p.parameter_id = cp.parameter_id
LEFT JOIN enum_position ep   ON ep.enum_position_id = epv.enum_val_id
WHERE epv.emobile_id = 1
ORDER BY cp.order_num;

-- Результат:
-- параметр        | единица | значение
-- Запас хода      | км      | 580
-- Масса           | кг      | 2100
-- Тип привода     | NULL    | Полный
-- Класс защиты IP | NULL    | IP67
```

**HTTP-эндпоинт для чтения:** `GET /api/emobile-parameter/byEmobile{id=1}`

---

## ПР4 — Хозяйственные операции

### 3.1 Концептуальное введение

#### Что такое «Хозяйственная операция» (ХО)?

Хозяйственная операция — это **задокументированное деловое событие**, которое что-то меняет в состоянии организации:

- **Отгрузка** 5 электромобилей покупателю ООО «ЭкоТранс»
- **Приёмка** партии батарей от поставщика АО «АккумТех»
- **Списание** 2 неисправных зарядных устройств

Каждая такая операция имеет:
- **Тип** (что произошло: отгрузка, приёмка, списание)
- **Участников** (кто отправил, кто получил, кто заплатил)
- **Роли участников** (Грузоотправитель, Грузополучатель, Плательщик)
- **Документы** (ТТН, счёт-фактура, акт)
- **Параметры** (сумма НДС, способ доставки)
- **Позиции** (конкретные изделия и количество)

#### Почему нужна гибкость?

Разные типы операций имеют разный набор участников и параметров:

| Тип ХО            | Обязательные роли                    | Параметры                        |
|--------------------|--------------------------------------|----------------------------------|
| Отгрузка           | Грузоотправитель, Грузополучатель    | Сумма НДС, Способ доставки       |
| Приёмка на склад   | Поставщик, Материально-ответственный | Дата поставки, Номер партии      |
| Списание           | Материально-ответственный            | Причина списания (enum)          |
| Межфилиальный перевод | Отправляющий филиал, Получающий   | Расстояние (км), Перевозчик      |

Если создавать отдельную таблицу под каждый тип ХО — получим десятки таблиц с хаотичной схемой. При появлении нового типа — снова миграция.

ПР4 решает эту проблему так же, как ПР2 и ПР3: **конфигурация типов ХО хранится в данных, а не в схеме**.

#### Аналогия: театр

Представьте, что ХО — это **театральная постановка**:

```
Тип ХО (ho_class)    ← это название спектакля («Отгрузка покупателю»)
Роли (ho_role)       ← это роли: Грузоотправитель, Грузополучатель
Допустимые роли      ← это сценарий: кто должен участвовать в этом спектакле
  (ho_class_role)

Экземпляр ХО         ← это конкретный показ спектакля (дата, номер)
  (ho_instance)
СХД (shd)            ← это реальные актёры: ООО «Завод», ООО «Покупатель»
Назначение на роль   ← это распределение ролей: ООО «Завод» играет Грузоотправителя
  (ho_actor)
```

Можно создать новый «спектакль» (тип ХО) без изменения «правил театра» (схемы БД).

#### Ключевые понятия

- **СХД (Субъект Хозяйственной Деятельности)** — организация или физлицо, участвующее в операции (завод, покупатель, перевозчик).
- **Роль** — в каком качестве СХД участвует в данной ХО (Грузоотправитель, Поставщик, Плательщик).
- **Тип ХО (`ho_class`)** — классификатор, организованный в дерево, как ПР1.

---

### 3.2 Схема базы данных ПР4

```sql
-- ============================================================
--  ПР4: Хозяйственные операции
--  Зависит от: parameter (ПР3), enum_position (ПР2), emobile (ПР1)
-- ============================================================

-- ───────────────────────────────────────────────────────────
--  Субъект хозяйственной деятельности (СХД)
--  Организация, подразделение или физлицо — участник операций
-- ───────────────────────────────────────────────────────────
CREATE TABLE shd (
    shd_id     SERIAL PRIMARY KEY,
    name       VARCHAR(200) NOT NULL,
    shd_type   VARCHAR(50)  NOT NULL
        CHECK (shd_type IN ('organization', 'individual', 'department')),
    -- organization — юридическое лицо (ООО, АО)
    -- individual   — физическое лицо / ИП
    -- department   — внутреннее подразделение компании
    inn        VARCHAR(12),          -- ИНН (для организаций и ИП)
    created_at TIMESTAMP DEFAULT NOW()
);

-- ───────────────────────────────────────────────────────────
--  Классификатор типов ХО (дерево, как ПР1)
--  Примеры: Отгрузка → Отгрузка покупателю, Отгрузка на склад
-- ───────────────────────────────────────────────────────────
CREATE TABLE ho_class (
    ho_class_id  SERIAL PRIMARY KEY,
    name         VARCHAR(200) NOT NULL,
    designation  VARCHAR(50),          -- краткий код: 'SHIPMENT', 'RECEIPT', 'WRITEOFF'
    parent_id    INT REFERENCES ho_class(ho_class_id) ON DELETE SET NULL,
    -- NULL → корневой узел; иначе → подтип
    is_terminal  BOOLEAN DEFAULT FALSE,
    -- TRUE → по этому типу можно создавать экземпляры ХО
    -- FALSE → промежуточный узел-группировщик (нельзя создать ХО)
    created_at   TIMESTAMP DEFAULT NOW()
);

-- ───────────────────────────────────────────────────────────
--  Классификатор ролей участников ХО
--  Роль описывает, в каком КАЧЕСТВЕ СХД участвует в операции
-- ───────────────────────────────────────────────────────────
CREATE TABLE ho_role (
    ho_role_id  SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL UNIQUE,
    -- 'Грузоотправитель', 'Грузополучатель', 'Плательщик', 'Поставщик'
    description TEXT
);

-- ───────────────────────────────────────────────────────────
--  Допустимые роли для типа ХО
--  Определяет, КАКИЕ роли участвуют в операции данного типа
-- ───────────────────────────────────────────────────────────
CREATE TABLE ho_class_role (
    id           SERIAL PRIMARY KEY,
    ho_class_id  INT NOT NULL REFERENCES ho_class(ho_class_id) ON DELETE CASCADE,
    ho_role_id   INT NOT NULL REFERENCES ho_role(ho_role_id)   ON DELETE CASCADE,
    is_required  BOOLEAN DEFAULT FALSE,
    -- TRUE → роль обязательна (нельзя подтвердить ХО без этого участника)
    -- FALSE → роль опциональна
    UNIQUE (ho_class_id, ho_role_id)
);

-- ───────────────────────────────────────────────────────────
--  Параметры типа ХО (переиспользует parameter из ПР3)
--  Например: тип «Отгрузка» имеет параметр «Сумма НДС»
-- ───────────────────────────────────────────────────────────
CREATE TABLE ho_class_parameter (
    id           SERIAL PRIMARY KEY,
    ho_class_id  INT   NOT NULL REFERENCES ho_class(ho_class_id)   ON DELETE CASCADE,
    parameter_id INT   NOT NULL REFERENCES parameter(parameter_id) ON DELETE CASCADE,
    order_num    INT   DEFAULT 0,   -- порядок отображения в форме
    min_val      FLOAT DEFAULT NULL,
    max_val      FLOAT DEFAULT NULL,
    UNIQUE (ho_class_id, parameter_id)
);

-- ───────────────────────────────────────────────────────────
--  Классификатор типов документов
--  ТТН, счёт-фактура, акт приёма-передачи и т.п.
-- ───────────────────────────────────────────────────────────
CREATE TABLE document_class (
    doc_class_id SERIAL PRIMARY KEY,
    name         VARCHAR(200) NOT NULL,
    -- 'Товарно-транспортная накладная', 'Счёт-фактура', 'Акт приёма-передачи'
    code         VARCHAR(50)
    -- Код по ОКУД (общероссийский классификатор управленческой документации)
    -- ТТН → '0504205', Счёт-фактура → '0309003'
);

-- ───────────────────────────────────────────────────────────
--  Допустимые типы документов для типа ХО
-- ───────────────────────────────────────────────────────────
CREATE TABLE ho_class_document (
    id           SERIAL PRIMARY KEY,
    ho_class_id  INT     NOT NULL REFERENCES ho_class(ho_class_id)       ON DELETE CASCADE,
    doc_class_id INT     NOT NULL REFERENCES document_class(doc_class_id) ON DELETE CASCADE,
    role_name    VARCHAR(100),
    -- 'Основной документ', 'Основание', 'Приложение'
    is_required  BOOLEAN DEFAULT FALSE
);

-- ───────────────────────────────────────────────────────────
--  Экземпляр ХО (конкретная хозяйственная операция)
-- ───────────────────────────────────────────────────────────
CREATE TABLE ho_instance (
    ho_id          SERIAL PRIMARY KEY,
    ho_class_id    INT            NOT NULL REFERENCES ho_class(ho_class_id) ON DELETE RESTRICT,
    -- RESTRICT: нельзя удалить тип ХО, если есть созданные по нему операции
    doc_number     VARCHAR(100),   -- номер документа: 'ТТН-2026-001'
    doc_date       DATE,           -- дата документа
    total_amount   DECIMAL(15,2),  -- итоговая сумма операции
    status         VARCHAR(30) DEFAULT 'draft'
        CHECK (status IN ('draft', 'confirmed', 'cancelled')),
    -- draft     → черновик, редактируется
    -- confirmed → подтверждена, изменения запрещены
    -- cancelled → отменена
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP DEFAULT NOW()
);

-- ───────────────────────────────────────────────────────────
--  Назначение СХД на роль в конкретной ХО
-- ───────────────────────────────────────────────────────────
CREATE TABLE ho_actor (
    id          SERIAL PRIMARY KEY,
    ho_id       INT NOT NULL REFERENCES ho_instance(ho_id) ON DELETE CASCADE,
    ho_role_id  INT NOT NULL REFERENCES ho_role(ho_role_id) ON DELETE RESTRICT,
    shd_id      INT NOT NULL REFERENCES shd(shd_id)         ON DELETE RESTRICT,
    UNIQUE (ho_id, ho_role_id)
    -- В одной ХО каждая роль может быть занята только одним СХД
);

-- ───────────────────────────────────────────────────────────
--  Значения параметров конкретной ХО
-- ───────────────────────────────────────────────────────────
CREATE TABLE ho_parameter_value (
    id                   SERIAL PRIMARY KEY,
    ho_id                INT NOT NULL REFERENCES ho_instance(ho_id)            ON DELETE CASCADE,
    ho_class_parameter_id INT NOT NULL REFERENCES ho_class_parameter(id)       ON DELETE CASCADE,
    val_real             FLOAT        DEFAULT NULL,  -- для param_type='real'
    val_int              INT          DEFAULT NULL,  -- для param_type='int'
    val_str              VARCHAR(500) DEFAULT NULL,  -- для param_type='str'
    val_date             DATE         DEFAULT NULL,  -- для param_type='date'
    enum_val_id          INT          DEFAULT NULL
        REFERENCES enum_position(enum_position_id) ON DELETE SET NULL,
    -- для param_type='enum' — FK на ПР2
    UNIQUE (ho_id, ho_class_parameter_id)
);

-- ───────────────────────────────────────────────────────────
--  Документы, приложенные к конкретной ХО
-- ───────────────────────────────────────────────────────────
CREATE TABLE ho_document (
    id           SERIAL PRIMARY KEY,
    ho_id        INT     NOT NULL REFERENCES ho_instance(ho_id)            ON DELETE CASCADE,
    doc_class_id INT     NOT NULL REFERENCES document_class(doc_class_id) ON DELETE RESTRICT,
    doc_number   VARCHAR(100),   -- номер конкретного документа
    doc_date     DATE,           -- дата документа
    file_path    VARCHAR(500)    -- путь к скан-копии файла
);

-- ───────────────────────────────────────────────────────────
--  Позиции ХО: конкретные изделия, участвующие в операции
-- ───────────────────────────────────────────────────────────
CREATE TABLE ho_position (
    id            SERIAL PRIMARY KEY,
    ho_id         INT            NOT NULL REFERENCES ho_instance(ho_id)  ON DELETE CASCADE,
    emobile_id    INT            REFERENCES emobile(emobile_id)          ON DELETE RESTRICT,
    -- FK на ПР1: конкретный электромобиль из каталога
    quantity      INT            NOT NULL DEFAULT 1
        CHECK (quantity > 0),
    unit_price    DECIMAL(15,2),  -- цена за единицу
    total_price   DECIMAL(15,2)   GENERATED ALWAYS AS (quantity * unit_price) STORED,
    -- вычисляемое поле: автоматически = quantity * unit_price
    note          TEXT,           -- примечание к позиции
    position_num  INT DEFAULT 0   -- порядковый номер позиции в документе
);

-- ───────────────────────────────────────────────────────────
--  Индексы для производительности
-- ───────────────────────────────────────────────────────────
CREATE INDEX IF NOT EXISTS idx_ho_instance_class  ON ho_instance(ho_class_id);
CREATE INDEX IF NOT EXISTS idx_ho_instance_status ON ho_instance(status);
CREATE INDEX IF NOT EXISTS idx_ho_actor_ho        ON ho_actor(ho_id);
CREATE INDEX IF NOT EXISTS idx_ho_position_ho     ON ho_position(ho_id);
CREATE INDEX IF NOT EXISTS idx_ho_param_val_ho    ON ho_parameter_value(ho_id);
```

---

### 3.3 SQL-процедуры ПР4

#### 1. `ins_ho_class` — создать тип ХО в классификаторе

```sql
-- Создаёт новый тип ХО в дереве классификатора.
-- p_parent_id = NULL → корневой узел (например, «Движение ТМЦ»)
-- p_is_terminal = TRUE → по этому типу можно создавать экземпляры ХО
-- p_is_terminal = FALSE → промежуточный узел-группировщик
CREATE OR REPLACE PROCEDURE ins_ho_class(
    p_name        VARCHAR,
    p_designation VARCHAR,
    p_parent_id   INT,
    p_is_terminal BOOLEAN
)
LANGUAGE plpgsql AS $$
BEGIN
    -- Если указан родитель, проверить что он существует
    IF p_parent_id IS NOT NULL THEN
        IF NOT EXISTS (SELECT 1 FROM ho_class WHERE ho_class_id = p_parent_id) THEN
            RAISE EXCEPTION 'Родительский тип ХО с id=% не найден', p_parent_id;
        END IF;
    END IF;

    -- Проверить уникальность обозначения
    IF p_designation IS NOT NULL AND EXISTS (
        SELECT 1 FROM ho_class WHERE designation = p_designation
    ) THEN
        RAISE EXCEPTION 'Тип ХО с обозначением "%" уже существует', p_designation;
    END IF;

    INSERT INTO ho_class (name, designation, parent_id, is_terminal)
    VALUES (p_name, p_designation, p_parent_id, p_is_terminal);
END;
$$;
```

**Пример:**
```sql
-- Создать корневой узел
CALL ins_ho_class('Движение ТМЦ', 'INVENTORY_MOVEMENT', NULL, FALSE);

-- Создать терминальный тип (подтип)
CALL ins_ho_class('Отгрузка покупателю', 'SHIPMENT_CUSTOMER', 1, TRUE);
```

---

#### 2. `add_param_to_ho_class` — добавить параметр к типу ХО

```sql
-- Привязывает существующий параметр (из ПР3) к типу ХО.
-- Переиспользует таблицу parameter: один параметр может быть
-- у нескольких типов ХО с разными ограничениями min/max.
CREATE OR REPLACE PROCEDURE add_param_to_ho_class(
    p_ho_class_id  INT,
    p_parameter_id INT,
    p_order_num    INT,
    p_min_val      FLOAT,
    p_max_val      FLOAT
)
LANGUAGE plpgsql AS $$
DECLARE
    v_param_type VARCHAR;
BEGIN
    -- Проверить существование типа ХО
    IF NOT EXISTS (SELECT 1 FROM ho_class WHERE ho_class_id = p_ho_class_id) THEN
        RAISE EXCEPTION 'Тип ХО с id=% не найден', p_ho_class_id;
    END IF;

    -- Проверить существование параметра
    IF NOT EXISTS (SELECT 1 FROM parameter WHERE parameter_id = p_parameter_id) THEN
        RAISE EXCEPTION 'Параметр с id=% не найден', p_parameter_id;
    END IF;

    -- min/max имеют смысл только для числовых типов
    SELECT param_type INTO v_param_type
    FROM parameter WHERE parameter_id = p_parameter_id;

    IF v_param_type NOT IN ('real', 'int') AND (p_min_val IS NOT NULL OR p_max_val IS NOT NULL) THEN
        RAISE WARNING 'min_val/max_val игнорируются для параметра типа %', v_param_type;
    END IF;

    INSERT INTO ho_class_parameter (ho_class_id, parameter_id, order_num, min_val, max_val)
    VALUES (p_ho_class_id, p_parameter_id, p_order_num, p_min_val, p_max_val)
    ON CONFLICT (ho_class_id, parameter_id)
    DO UPDATE SET order_num = EXCLUDED.order_num,
                  min_val   = EXCLUDED.min_val,
                  max_val   = EXCLUDED.max_val;
END;
$$;
```

---

#### 3. `add_role_to_ho_class` — добавить допустимую роль к типу ХО

```sql
-- Определяет, что для типа ХО требуется участник с данной ролью.
-- is_required=TRUE → нельзя подтвердить ХО без этого участника.
CREATE OR REPLACE PROCEDURE add_role_to_ho_class(
    p_ho_class_id INT,
    p_ho_role_id  INT,
    p_is_required BOOLEAN
)
LANGUAGE plpgsql AS $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM ho_class WHERE ho_class_id = p_ho_class_id) THEN
        RAISE EXCEPTION 'Тип ХО с id=% не найден', p_ho_class_id;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM ho_role WHERE ho_role_id = p_ho_role_id) THEN
        RAISE EXCEPTION 'Роль с id=% не найдена', p_ho_role_id;
    END IF;

    INSERT INTO ho_class_role (ho_class_id, ho_role_id, is_required)
    VALUES (p_ho_class_id, p_ho_role_id, p_is_required)
    ON CONFLICT (ho_class_id, ho_role_id)
    DO UPDATE SET is_required = EXCLUDED.is_required;
END;
$$;
```

---

#### 4. `ins_ho` — создать экземпляр ХО

```sql
-- Создаёт конкретную хозяйственную операцию.
-- Возвращает ho_id созданного экземпляра.
-- Важно: создаётся только в статусе 'draft' — операция не подтверждена.
CREATE OR REPLACE FUNCTION ins_ho(
    p_ho_class_id  INT,
    p_doc_number   VARCHAR,
    p_doc_date     DATE,
    p_total_amount DECIMAL
)
RETURNS INT
LANGUAGE plpgsql AS $$
DECLARE
    v_ho_id       INT;
    v_is_terminal BOOLEAN;
BEGIN
    -- Проверить, что тип ХО существует и является терминальным
    SELECT is_terminal INTO v_is_terminal
    FROM ho_class WHERE ho_class_id = p_ho_class_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Тип ХО с id=% не найден', p_ho_class_id;
    END IF;

    IF NOT v_is_terminal THEN
        RAISE EXCEPTION 'Тип ХО id=% является группировочным (is_terminal=FALSE). '
                        'Экземпляры создаются только для терминальных типов.', p_ho_class_id;
    END IF;

    -- Создать экземпляр в статусе 'draft'
    INSERT INTO ho_instance (ho_class_id, doc_number, doc_date, total_amount, status)
    VALUES (p_ho_class_id, p_doc_number, p_doc_date, p_total_amount, 'draft')
    RETURNING ho_id INTO v_ho_id;

    RETURN v_ho_id;
END;
$$;
```

**Пример:**
```sql
SELECT ins_ho(2, 'ТТН-2026-001', '2026-04-06', 6000000.00);
-- Вернёт: 1  (ho_id нового черновика)
```

---

#### 5. `set_ho_actor` — назначить СХД на роль в ХО

```sql
-- Назначает организацию (СХД) на роль в конкретной ХО.
-- Проверяет, что данная роль допустима для типа этой ХО.
-- Повторный вызов с тем же ho_id+ho_role_id заменяет предыдущее назначение.
CREATE OR REPLACE PROCEDURE set_ho_actor(
    p_ho_id       INT,
    p_ho_role_id  INT,
    p_shd_id      INT
)
LANGUAGE plpgsql AS $$
DECLARE
    v_ho_class_id INT;
    v_status      VARCHAR;
BEGIN
    -- Получить тип и статус ХО
    SELECT ho_class_id, status INTO v_ho_class_id, v_status
    FROM ho_instance WHERE ho_id = p_ho_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'ХО с id=% не найдена', p_ho_id;
    END IF;

    -- Нельзя изменять подтверждённую операцию
    IF v_status = 'confirmed' THEN
        RAISE EXCEPTION 'ХО id=% уже подтверждена. Изменения запрещены.', p_ho_id;
    END IF;

    -- Проверить, что роль допустима для данного типа ХО
    IF NOT EXISTS (
        SELECT 1 FROM ho_class_role
        WHERE ho_class_id = v_ho_class_id AND ho_role_id = p_ho_role_id
    ) THEN
        RAISE EXCEPTION 'Роль id=% не предусмотрена для типа ХО id=%',
                        p_ho_role_id, v_ho_class_id;
    END IF;

    -- Проверить существование СХД
    IF NOT EXISTS (SELECT 1 FROM shd WHERE shd_id = p_shd_id) THEN
        RAISE EXCEPTION 'СХД с id=% не найден', p_shd_id;
    END IF;

    -- Назначить (UPSERT)
    INSERT INTO ho_actor (ho_id, ho_role_id, shd_id)
    VALUES (p_ho_id, p_ho_role_id, p_shd_id)
    ON CONFLICT (ho_id, ho_role_id)
    DO UPDATE SET shd_id = EXCLUDED.shd_id;
END;
$$;
```

---

#### 6. `write_ho_par` — записать значение параметра ХО

```sql
-- Проверяет тип параметра и ограничения, затем сохраняет значение.
-- Заполняется только соответствующее поле val_*; остальные остаются NULL.
CREATE OR REPLACE PROCEDURE write_ho_par(
    p_ho_id                INT,
    p_ho_class_parameter_id INT,
    p_val_real             FLOAT   DEFAULT NULL,
    p_val_int              INT     DEFAULT NULL,
    p_val_str              VARCHAR DEFAULT NULL,
    p_val_date             DATE    DEFAULT NULL,
    p_enum_val_id          INT     DEFAULT NULL
)
LANGUAGE plpgsql AS $$
DECLARE
    v_param_type  VARCHAR;
    v_enum_class  INT;
    v_min_val     FLOAT;
    v_max_val     FLOAT;
    v_ho_class_id INT;
    v_status      VARCHAR;
BEGIN
    -- Получить статус ХО
    SELECT ho_class_id, status INTO v_ho_class_id, v_status
    FROM ho_instance WHERE ho_id = p_ho_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'ХО с id=% не найдена', p_ho_id;
    END IF;

    IF v_status = 'confirmed' THEN
        RAISE EXCEPTION 'ХО id=% подтверждена. Изменение параметров запрещено.', p_ho_id;
    END IF;

    -- Получить тип параметра и ограничения
    SELECT p.param_type, p.enum_class_id, hcp.min_val, hcp.max_val
    INTO v_param_type, v_enum_class, v_min_val, v_max_val
    FROM ho_class_parameter hcp
    JOIN parameter p ON p.parameter_id = hcp.parameter_id
    WHERE hcp.id = p_ho_class_parameter_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Параметр ho_class_parameter id=% не найден', p_ho_class_parameter_id;
    END IF;

    -- Валидация по типу
    IF v_param_type = 'real' THEN
        IF p_val_real IS NULL THEN
            RAISE EXCEPTION 'Для параметра типа real необходимо передать p_val_real';
        END IF;
        IF v_min_val IS NOT NULL AND p_val_real < v_min_val THEN
            RAISE EXCEPTION 'Значение % меньше допустимого минимума %', p_val_real, v_min_val;
        END IF;
        IF v_max_val IS NOT NULL AND p_val_real > v_max_val THEN
            RAISE EXCEPTION 'Значение % превышает допустимый максимум %', p_val_real, v_max_val;
        END IF;

    ELSIF v_param_type = 'int' THEN
        IF p_val_int IS NULL THEN
            RAISE EXCEPTION 'Для параметра типа int необходимо передать p_val_int';
        END IF;
        IF v_min_val IS NOT NULL AND p_val_int < v_min_val THEN
            RAISE EXCEPTION 'Значение % меньше допустимого минимума %', p_val_int, v_min_val;
        END IF;
        IF v_max_val IS NOT NULL AND p_val_int > v_max_val THEN
            RAISE EXCEPTION 'Значение % превышает допустимый максимум %', p_val_int, v_max_val;
        END IF;

    ELSIF v_param_type = 'enum' THEN
        IF p_enum_val_id IS NULL THEN
            RAISE EXCEPTION 'Для параметра типа enum необходимо передать p_enum_val_id';
        END IF;
        -- Проверить, что выбранная позиция принадлежит правильному справочнику (ПР2)
        IF NOT EXISTS (
            SELECT 1 FROM enum_position
            WHERE enum_position_id = p_enum_val_id
              AND enum_class_id = v_enum_class
        ) THEN
            RAISE EXCEPTION 'enum_position id=% не принадлежит справочнику id=%',
                            p_enum_val_id, v_enum_class;
        END IF;

    ELSIF v_param_type = 'str' THEN
        IF p_val_str IS NULL THEN
            RAISE EXCEPTION 'Для параметра типа str необходимо передать p_val_str';
        END IF;
    END IF;

    -- Сохранить значение (UPSERT)
    INSERT INTO ho_parameter_value
        (ho_id, ho_class_parameter_id, val_real, val_int, val_str, val_date, enum_val_id)
    VALUES
        (p_ho_id, p_ho_class_parameter_id,
         p_val_real, p_val_int, p_val_str, p_val_date, p_enum_val_id)
    ON CONFLICT (ho_id, ho_class_parameter_id)
    DO UPDATE SET
        val_real    = EXCLUDED.val_real,
        val_int     = EXCLUDED.val_int,
        val_str     = EXCLUDED.val_str,
        val_date    = EXCLUDED.val_date,
        enum_val_id = EXCLUDED.enum_val_id;
END;
$$;
```

---

#### 7. `add_ho_position` — добавить позицию (изделие) в ХО

```sql
-- Добавляет конкретный электромобиль из каталога (ПР1) в состав ХО.
-- total_price вычисляется автоматически (GENERATED ALWAYS AS).
CREATE OR REPLACE PROCEDURE add_ho_position(
    p_ho_id        INT,
    p_emobile_id   INT,
    p_quantity     INT,
    p_unit_price   DECIMAL,
    p_note         TEXT    DEFAULT NULL
)
LANGUAGE plpgsql AS $$
DECLARE
    v_status     VARCHAR;
    v_pos_num    INT;
BEGIN
    -- Проверить статус ХО
    SELECT status INTO v_status FROM ho_instance WHERE ho_id = p_ho_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'ХО с id=% не найдена', p_ho_id;
    END IF;

    IF v_status = 'confirmed' THEN
        RAISE EXCEPTION 'ХО id=% подтверждена. Добавление позиций запрещено.', p_ho_id;
    END IF;

    -- Проверить, что изделие существует в каталоге
    IF NOT EXISTS (SELECT 1 FROM emobile WHERE emobile_id = p_emobile_id) THEN
        RAISE EXCEPTION 'Электромобиль с id=% не найден в каталоге', p_emobile_id;
    END IF;

    IF p_quantity <= 0 THEN
        RAISE EXCEPTION 'Количество должно быть больше нуля, получено: %', p_quantity;
    END IF;

    -- Определить следующий порядковый номер позиции
    SELECT COALESCE(MAX(position_num), 0) + 1 INTO v_pos_num
    FROM ho_position WHERE ho_id = p_ho_id;

    INSERT INTO ho_position (ho_id, emobile_id, quantity, unit_price, note, position_num)
    VALUES (p_ho_id, p_emobile_id, p_quantity, p_unit_price, p_note, v_pos_num);

    -- Пересчитать итоговую сумму ХО
    UPDATE ho_instance
    SET total_amount = (
        SELECT COALESCE(SUM(total_price), 0)
        FROM ho_position
        WHERE ho_id = p_ho_id
    ),
    updated_at = NOW()
    WHERE ho_id = p_ho_id;
END;
$$;
```

---

#### 8. `find_ho_by_class` — найти ХО по типу

```sql
-- Возвращает список операций указанного типа с кратким описанием участников.
-- actor_roles — агрегированная строка вида "Грузоотправитель: ООО Завод; Грузополучатель: ООО Клиент"
CREATE OR REPLACE FUNCTION find_ho_by_class(p_ho_class_id INT)
RETURNS TABLE (
    ho_id        INT,
    doc_number   VARCHAR,
    doc_date     DATE,
    total_amount DECIMAL,
    status       VARCHAR,
    actor_roles  TEXT
)
LANGUAGE sql AS $$
    SELECT
        hi.ho_id,
        hi.doc_number,
        hi.doc_date,
        hi.total_amount,
        hi.status,
        STRING_AGG(
            hr.name || ': ' || s.name,
            '; ' ORDER BY hr.name
        ) AS actor_roles
    FROM ho_instance hi
    LEFT JOIN ho_actor  ha ON ha.ho_id     = hi.ho_id
    LEFT JOIN ho_role   hr ON hr.ho_role_id = ha.ho_role_id
    LEFT JOIN shd        s ON s.shd_id      = ha.shd_id
    WHERE hi.ho_class_id = p_ho_class_id
    GROUP BY hi.ho_id, hi.doc_number, hi.doc_date, hi.total_amount, hi.status
    ORDER BY hi.doc_date DESC, hi.ho_id DESC;
$$;
```

**Пример вызова:**
```sql
SELECT * FROM find_ho_by_class(2);
-- ho_id | doc_number    | doc_date   | total_amount | status    | actor_roles
-- 1     | ТТН-2026-001  | 2026-04-06 | 6000000.00   | confirmed | Грузополучатель: ООО ЭкоТранс; Грузоотправитель: АО ЭлектроМобиль
-- 2     | ТТН-2026-002  | 2026-04-07 | 3200000.00   | draft     | Грузоотправитель: АО ЭлектроМобиль
```

---

### 3.4 Связь ПР4 с ПР1, ПР2, ПР3

| Связь                     | Как реализована                                                          |
|---------------------------|--------------------------------------------------------------------------|
| **ПР4 ← ПР1**             | `ho_position.emobile_id` → FK на `emobile.emobile_id`: изделия из каталога попадают в позиции ХО |
| **ПР4 ← ПР2**             | `ho_parameter_value.enum_val_id` → FK на `enum_position.enum_position_id`: параметры типа `enum` используют справочники ПР2 |
| **ПР4 ← ПР3**             | `ho_class_parameter.parameter_id` → FK на `parameter.parameter_id`: типы ХО переиспользуют те же шаблоны параметров, что и компоненты изделий |
| **ПР4 ← ПР1 (структура)** | `ho_class` организован в дерево с `parent_id`, как классификатор изделий в ПР1 |

**Принцип повторного использования ПР3 в ПР4:**

Параметр «Сумма НДС» создаётся один раз в таблице `parameter`:
```sql
INSERT INTO parameter (designation, name, param_type, measuring_unit)
VALUES ('vat_amount', 'Сумма НДС', 'real', 'руб.');
```

Затем тот же параметр привязывается и к компоненту, и к типу ХО:
```sql
-- Привязка к компоненту (ПР3): НДС у электромобиля
INSERT INTO component_parameter (component_type, parameter_id) VALUES ('emobile', <id>);

-- Привязка к типу ХО (ПР4): НДС в операции отгрузки
INSERT INTO ho_class_parameter (ho_class_id, parameter_id) VALUES (2, <id>);
```

---

### 3.5 Практический пример: «Отгрузка электромобиля покупателю»

Полный сценарий от настройки до оформления конкретной операции:

#### Этап 1: Настройка (выполняется один раз администратором)

**Шаг 1.1 — Создать тип ХО**
```sql
CALL ins_ho_class('Операции с ТМЦ', NULL, NULL, FALSE);
-- → ho_class_id = 1 (корневой узел, группировщик)

CALL ins_ho_class('Отгрузка покупателю', 'SHIPMENT_CUSTOMER', 1, TRUE);
-- → ho_class_id = 2 (терминальный, можно создавать ХО)
```

**Шаг 1.2 — Создать роли участников**
```sql
INSERT INTO ho_role (name, description) VALUES
('Грузоотправитель', 'Организация, отправляющая груз'),
('Грузополучатель',  'Организация, получающая груз'),
('Плательщик',       'Организация, оплачивающая поставку');
-- → ho_role_id: 1, 2, 3
```

**Шаг 1.3 — Привязать роли к типу ХО**
```sql
CALL add_role_to_ho_class(2, 1, TRUE);   -- Грузоотправитель — обязательный
CALL add_role_to_ho_class(2, 2, TRUE);   -- Грузополучатель — обязательный
CALL add_role_to_ho_class(2, 3, FALSE);  -- Плательщик — опциональный
```

**Шаг 1.4 — Создать параметры и справочники**
```sql
-- Создать справочник "Способ доставки" (ПР2)
INSERT INTO enum_class (name) VALUES ('Способ доставки') RETURNING enum_class_id;
-- → enum_class_id = 10

INSERT INTO enum_position (enum_class_id, value, order_num) VALUES
(10, 'Автотранспорт', 1),
(10, 'Ж/Д',          2),
(10, 'Авиа',         3),
(10, 'Самовывоз',    4);

-- Создать параметры (ПР3)
INSERT INTO parameter (designation, name, param_type, measuring_unit)
VALUES ('vat_amount', 'Сумма НДС', 'real', 'руб.')
RETURNING parameter_id;  -- → 20

INSERT INTO parameter (designation, name, param_type, enum_class_id)
VALUES ('delivery_method', 'Способ доставки', 'enum', 10)
RETURNING parameter_id;  -- → 21
```

**Шаг 1.5 — Привязать параметры к типу ХО**
```sql
CALL add_param_to_ho_class(2, 20, 1, 0.0, NULL);   -- Сумма НДС, order=1
CALL add_param_to_ho_class(2, 21, 2, NULL, NULL);   -- Способ доставки, order=2
-- → ho_class_parameter.id: 1 (НДС), 2 (доставка)
```

**Шаг 1.6 — Зарегистрировать типы документов**
```sql
INSERT INTO document_class (name, code) VALUES
('Товарно-транспортная накладная', '0504205'),
('Счёт-фактура', '0309003');
-- → doc_class_id: 1, 2

INSERT INTO ho_class_document (ho_class_id, doc_class_id, role_name, is_required) VALUES
(2, 1, 'Основной документ', TRUE),   -- ТТН обязательна
(2, 2, 'Основание',          FALSE);  -- счёт-фактура — опционально
```

**Шаг 1.7 — Зарегистрировать участников (СХД)**
```sql
INSERT INTO shd (name, shd_type, inn) VALUES
('АО ЭлектроМобиль',  'organization', '7700000001'),  -- наш завод, shd_id=1
('ООО ЭкоТранс',      'organization', '7700000002');  -- покупатель, shd_id=2
```

---

#### Этап 2: Работа (оформление конкретной отгрузки)

**Шаг 2.1 — Создать экземпляр ХО**
```sql
SELECT ins_ho(2, 'ТТН-2026-001', '2026-04-06', NULL);
-- Вернёт ho_id = 1
-- NULL в total_amount: сумма будет посчитана автоматически при добавлении позиций
```

**Шаг 2.2 — Назначить участников на роли**
```sql
CALL set_ho_actor(1, 1, 1);  -- ХО #1, роль "Грузоотправитель", СХД "АО ЭлектроМобиль"
CALL set_ho_actor(1, 2, 2);  -- ХО #1, роль "Грузополучатель",  СХД "ООО ЭкоТранс"
```

**Шаг 2.3 — Заполнить параметры ХО**
```sql
-- Сумма НДС = 1 000 000 руб.
CALL write_ho_par(1, 1, p_val_real := 1000000.00);

-- Способ доставки = "Автотранспорт" (enum_position_id = 1)
CALL write_ho_par(1, 2, p_enum_val_id := 1);
```

**Шаг 2.4 — Добавить изделия в позиции ХО**
```sql
-- 2 электромобиля #1 по 3 000 000 руб. каждый
CALL add_ho_position(1, 1, 2, 3000000.00, 'Партия 2026-04 / синий металлик');
-- total_price автоматически = 2 * 3000000 = 6000000
-- ho_instance.total_amount автоматически пересчитан = 6000000
```

**Шаг 2.5 — Прикрепить документы**
```sql
INSERT INTO ho_document (ho_id, doc_class_id, doc_number, doc_date, file_path)
VALUES (1, 1, 'ТТН-2026-001', '2026-04-06', '/docs/2026/ttn_2026_001.pdf');
```

**Шаг 2.6 — Подтвердить операцию**
```sql
UPDATE ho_instance
SET status = 'confirmed', updated_at = NOW()
WHERE ho_id = 1;
-- После этого изменения в ХО запрещены (set_ho_actor и write_ho_par выбросят EXCEPTION)
```

**Проверка результата:**
```sql
SELECT * FROM find_ho_by_class(2);
-- ho_id | doc_number   | doc_date   | total_amount | status    | actor_roles
-- 1     | ТТН-2026-001 | 2026-04-06 | 6000000.00   | confirmed | Грузополучатель: ООО ЭкоТранс; Грузоотправитель: АО ЭлектроМобиль
```

---

### 3.6 Диаграмма потока данных (ASCII)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  НАСТРОЙКА (однократно — администратор)                                     │
│                                                                             │
│  ho_class (дерево типов ХО)                                                 │
│    │                                                                        │
│    ├─── ho_class_role ──────────────► ho_role                               │
│    │      is_required                 (Грузоотправитель,                    │
│    │                                   Грузополучатель...)                  │
│    ├─── ho_class_parameter ──────────► parameter ◄── (из ПР3!)             │
│    │      min_val, max_val              │                                   │
│    │                                   └──► enum_class ◄── (из ПР2!)       │
│    └─── ho_class_document ──────────► document_class                        │
│             role_name, is_required      (ТТН, Счёт-фактура...)             │
└─────────────────────────────────────────────────────────────────────────────┘
                          │ конфигурация используется при
                          ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  РАБОТА (регулярно — оператор)                                              │
│                                                                             │
│  ho_instance ◄── тип из ho_class                                            │
│    │  (doc_number, doc_date, status)                                        │
│    │                                                                        │
│    ├─── ho_actor ─────────────────────► shd                                │
│    │      │                              (АО Завод, ООО Клиент...)         │
│    │      └──► ho_role (какая роль)                                        │
│    │                                                                        │
│    ├─── ho_parameter_value                                                  │
│    │      │  val_real / val_int / val_str / val_date                        │
│    │      └──► enum_val_id ──────────────► enum_position (ПР2)            │
│    │                                                                        │
│    ├─── ho_document                                                         │
│    │      └──► document_class (тип документа)                              │
│    │                                                                        │
│    └─── ho_position                                                         │
│           │  quantity, unit_price, total_price (вычисляемое)               │
│           └──► emobile ───────────────────────► (ПР1: каталог изделий)    │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Итоговая диаграмма: как ПР1–ПР4 связаны вместе

```
╔══════════════════════════════════════════════════════════════════════════════╗
║                           MISPRIS — архитектура данных                      ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║   ПР1: КЛАССИФИКАТОР ИЗДЕЛИЙ          ПР2: ПЕРЕЧИСЛЕНИЯ                     ║
║   ─────────────────────────           ─────────────────                     ║
║   emobile_class (дерево)              enum_class                            ║
║     └─ emobile (изделия)                └─ enum_position                   ║
║          │                                    │                             ║
║          │ ◄────────────────────────────────── │                            ║
║          │  parameter.enum_class_id указывает  │                            ║
║          │  на справочник допустимых значений   │                            ║
║          │                                    │                             ║
║   ПР3: ПАРАМЕТРЫ ИЗДЕЛИЙ              │                             ║
║   ─────────────────────────           │                             ║
║   parameter (шаблоны) ◄───────────────┘                             ║
║     │  (designation, param_type,      │ param_type='enum' → FK      ║
║     │   measuring_unit)               │ на enum_position            ║
║     │                                                               ║
║     ├─► component_parameter           ║
║     │     (component_type,            ║
║     │      min_val, max_val)          ║
║     │          │ привязка к типу      ║
║     │          ▼                      ║
║     └─► emobile_parameter_value ──────────────────────────────────── ║
║              (конкретные значения)    │ enum_val_id → enum_position  ║
║                    │                  │ (ПР2 даёт допустимые варианты)║
║                    │ emobile_id → emobile (ПР1)                      ║
║                                                                       ║
║   ПР4: ХОЗЯЙСТВЕННЫЕ ОПЕРАЦИИ                                        ║
║   ────────────────────────────                                        ║
║   ho_class (дерево типов)                                             ║
║     │  (как emobile_class в ПР1,                                     ║
║     │   но для бизнес-событий)                                        ║
║     │                                                                 ║
║     ├─► ho_class_role ──► ho_role                                    ║
║     │     (допустимые роли)           │                              ║
║     │                                 │                              ║
║     ├─► ho_class_parameter ──────────► parameter (из ПР3!)          ║
║     │     (параметры типа ХО)         │ переиспользуем тот же        ║
║     │                                 │ каталог параметров           ║
║     └─► ho_class_document ──► document_class                         ║
║                                                                       ║
║   ho_instance (конкретная операция)                                   ║
║     ├─► ho_actor ──────────────────────────────► shd                 ║
║     │     (кто участвует, в какой роли)                              ║
║     │                                                                 ║
║     ├─► ho_parameter_value                                            ║
║     │     └─► enum_val_id ─────────────────────► enum_position (ПР2)║
║     │                                                                 ║
║     └─► ho_position ────────────────────────────► emobile (ПР1!)    ║
║           (что отгружается / принимается)         изделие из каталога║
╚══════════════════════════════════════════════════════════════════════════════╝

Стрелки зависимостей:
  ПР4 ──uses──► ПР3 (parameter): параметры ХО — те же шаблоны, что и параметры изделий
  ПР4 ──uses──► ПР2 (enum_position): параметры типа 'enum' ссылаются на справочники
  ПР4 ──uses──► ПР1 (emobile): позиции ХО содержат изделия из каталога
  ПР3 ──uses──► ПР2 (enum_class): параметр типа 'enum' ссылается на справочник
```

### Краткая таблица пересечений

| Откуда | Куда | Связь | Пример |
|--------|------|-------|--------|
| ПР3 `parameter` | ПР2 `enum_class` | `parameter.enum_class_id` | Параметр «Тип привода» ссылается на справочник «Тип привода» |
| ПР3 `emobile_parameter_value` | ПР2 `enum_position` | `val.enum_val_id` | Конкретный авто имеет «Тип привода = Полный» |
| ПР3 `emobile_parameter_value` | ПР1 `emobile` | `val.emobile_id` | Значение принадлежит конкретному изделию |
| ПР4 `ho_class_parameter` | ПР3 `parameter` | `hcp.parameter_id` | Тип ХО «Отгрузка» имеет параметр «Сумма НДС» |
| ПР4 `ho_parameter_value` | ПР2 `enum_position` | `hpv.enum_val_id` | ХО №1 имеет «Способ доставки = Автотранспорт» |
| ПР4 `ho_position` | ПР1 `emobile` | `pos.emobile_id` | В отгрузке №1 участвует электромобиль с id=1 |

---

*Документ описывает архитектуру MISPRIS v2 — Go/PostgreSQL бэкенд каталога электромобилей с поддержкой динамических перечислений (ПР2), параметрического описания изделий (ПР3) и хозяйственных операций (ПР4).*

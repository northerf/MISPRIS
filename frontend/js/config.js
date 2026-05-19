/* config.js */
const API_BASE = 'http://localhost:8080/api';

function escapeHtml(s){
    if(!s&&s!==0)return'';
    return String(s).replace(/[&<>\"']/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','\"':'&quot;',"'":'&#39;'}[c]));
}
function shortId(id){return id!==undefined&&id!==null?String(id).slice(0,8):'—'}

const PG_TABLE_RU={
    engine:'Двигатель',inverter:'Инвертор',gearbox:'КПП',
    frame:'Рама',suspension:'Подвеска',break_system:'Тормозная система',
    charger:'Зарядное устройство',connector:'Коннектор',
    controller:'Контроллер',sensor:'Датчик',wiring:'Проводка',
    carcass:'Каркас',doors:'Двери',wings:'Крылья',battery:'Батарея',
    power_point:'Силовая установка',chassis:'Шасси',
    charger_system:'Зарядная система',electronics:'Электроника',
    body:'Кузов',emobile:'Автомобиль',users:'Пользователь',
};
function friendlyError(raw){
    if(!raw)return'Неизвестная ошибка';
    const s=String(raw);
    const fkMatch=s.match(/on table[s]?\s+"?(\w+)"?.+on table[s]?\s+"?(\w+)"?/i);
    if((s.includes('23001')||s.includes('23503')||s.includes('violates'))&&fkMatch){
        const src=PG_TABLE_RU[fkMatch[1]]||fkMatch[1];
        const dst=PG_TABLE_RU[fkMatch[2]]||fkMatch[2];
        return`Нельзя удалить/изменить «${src}» — используется в «${dst}». Сначала удалите связанную запись.`;
    }
    if(s.includes('23001')||s.includes('23503')||s.includes('violates'))
        return'Нельзя выполнить: запись используется в другой таблице.';
    if(s.includes('23505')||s.includes('duplicate key'))
        return'Такая запись уже существует.';
    if(s.includes('22P02')||s.includes('invalid input syntax'))
        return'Некорректный формат данных. Выберите значение из списка.';
    if(s.includes('cannot unmarshal'))
        return'Ошибка типов данных. Обновите страницу и попробуйте снова.';
    if(s.includes('42703')||s.includes('does not exist'))
        return'Ошибка схемы БД. Сообщите разработчику.';
    if(s.includes('401')||s.toLowerCase().includes('unauthorized'))return'Необходима авторизация.';
    if(s.includes('403')||s.toLowerCase().includes('forbidden'))return'Нет прав на это действие.';
    if(s.includes('404'))return'Запись не найдена.';
    if(s.includes('500'))return'Внутренняя ошибка сервера.';
    if(s.toLowerCase().includes('failed to fetch')||s.toLowerCase().includes('networkerror'))
        return'Сервер недоступен. Проверьте подключение.';
    return s;
}

// ============================================================
//  ЗАГРУЗКА ЗНАЧЕНИЙ ПЕРЕЧИСЛЕНИЯ (ПР2)
//  БЕЗ кэша — каждый раз идёт в БД, чтобы сразу видеть
//  значения добавленные через «Справочники».
// ============================================================
async function loadEnumCacheForType(componentType){
    if(!componentType) return [];
    try{
        // Ищем enum_class для данного типа компонента
        const classes = await apiRequest(`${API_BASE}/enum-class/list`);
        const cls = (classes||[]).find(c => c.component_type === componentType);
        if(!cls) return [];
        // Загружаем значения упорядоченные по order_num
        const vals = await apiRequest(`${API_BASE}/enum-class/values${cls.enum_class_id}`);
        return vals || [];
    }catch(e){
        console.warn(`[enum load] ${componentType}:`, e.message);
        return [];
    }
}

/*
 * subFields:
 *   id            — id HTML-элемента
 *   dbField       — ключ из ответа API
 *   label         — метка поля
 *   type          — 'text' | 'number' | 'textarea' | 'select' | 'select_enum'
 *   enumComponent — для 'select_enum': component_type в enum_class
 *   options       — fallback-список если справочник ПР2 пуст
 */
const CATEGORY_MAP=[
    /* ── Батареи ────────────────────────────────────────── */
    {key:'batteries',label:'Батареи',group:'battery',getWord:'Battery',
     idField:'battery_id',nameField:'battery_name',
     infoFn:i=>[i.battery_type,
                i.battery_capacity?`${i.battery_capacity} кВтч`:'',
                i.battery_info].filter(Boolean).join(' · '),
     subFields:[
         {id:'name',             dbField:'battery_name',     label:'Название',       type:'text'},
         {id:'battery_type',     dbField:'battery_type',     label:'Тип',            type:'select_enum', enumComponent:'battery', options:['Li-ion','Li-polymer']},
         {id:'battery_capacity', dbField:'battery_capacity', label:'Ёмкость (кВтч)', type:'number'},
         {id:'info',             dbField:'battery_info',     label:'Информация',     type:'textarea'},
     ]},
    /* ── Двигатели ───────────────────────────────────────── */
    {key:'engines',label:'Двигатели',group:'engine',getWord:'Engine',
     idField:'engine_id',nameField:'engine_name',
     infoFn:i=>[i.engine_type,i.engine_info].filter(Boolean).join(' · '),
     subFields:[
         {id:'name',        dbField:'engine_name', label:'Название',   type:'text'},
         {id:'engine_type', dbField:'engine_type', label:'Тип',        type:'select_enum', enumComponent:'engine', options:['AC','DC']},
         {id:'info',        dbField:'engine_info', label:'Информация', type:'textarea'},
     ]},
    /* ── Инверторы ───────────────────────────────────────── */
    {key:'inverters',label:'Инверторы',group:'inverter',getWord:'Inverter',
     idField:'inverter_id',nameField:'inverter_name',
     infoFn:i=>i.inverter_info||'',
     subFields:[
         {id:'name',dbField:'inverter_name',label:'Название',  type:'text'},
         {id:'info',dbField:'inverter_info',label:'Информация',type:'textarea'},
     ]},
    /* ── КПП ─────────────────────────────────────────────── */
    {key:'gearboxes',label:'КПП',group:'gearbox',getWord:'Gearbox',
     idField:'gearbox_id',nameField:'gearbox_name',
     infoFn:i=>i.gearbox_info||'',
     subFields:[
         {id:'name',dbField:'gearbox_name',label:'Название',  type:'text'},
         {id:'info',dbField:'gearbox_info',label:'Информация',type:'textarea'},
     ]},
    /* ── Рамы ────────────────────────────────────────────── */
    {key:'frames',label:'Рамы',group:'frame',getWord:'Frame',
     idField:'frame_id',nameField:'frame_name',
     infoFn:i=>i.frame_info||'',
     subFields:[
         {id:'name',dbField:'frame_name',label:'Название',  type:'text'},
         {id:'info',dbField:'frame_info',label:'Информация',type:'textarea'},
     ]},
    /* ── Подвески ────────────────────────────────────────── */
    {key:'suspensions',label:'Подвески',group:'suspension',getWord:'Suspension',
     idField:'suspension_id',nameField:'suspension_name',
     infoFn:i=>i.suspension_info||'',
     subFields:[
         {id:'name',dbField:'suspension_name',label:'Название',  type:'text'},
         {id:'info',dbField:'suspension_info',label:'Информация',type:'textarea'},
     ]},
    /* ── Тормозные системы ───────────────────────────────── */
    {key:'breakSystems',label:'Тормозные системы',group:'break-system',getWord:'BreakSystem',
     idField:'break_system_id',nameField:'break_system_name',
     infoFn:i=>i.break_system_info||'',
     subFields:[
         {id:'name',dbField:'break_system_name',label:'Название',  type:'text'},
         {id:'info',dbField:'break_system_info',label:'Информация',type:'textarea'},
     ]},
    /* ── Зарядные устройства ─────────────────────────────── */
    {key:'chargers',label:'Зарядные устройства',group:'charger',getWord:'Charger',
     idField:'charger_id',nameField:'charger_name',
     infoFn:i=>i.charger_info||'',
     subFields:[
         {id:'name',dbField:'charger_name',label:'Название',  type:'text'},
         {id:'info',dbField:'charger_info',label:'Информация',type:'textarea'},
     ]},
    /* ── Коннекторы ──────────────────────────────────────── */
    {key:'connectors',label:'Коннекторы',group:'connector',getWord:'Connector',
     idField:'connector_id',nameField:'connector_name',
     infoFn:i=>i.connector_info||'',
     subFields:[
         {id:'name',dbField:'connector_name',label:'Название',  type:'text'},
         {id:'info',dbField:'connector_info',label:'Информация',type:'textarea'},
     ]},
    /* ── Контроллеры ─────────────────────────────────────── */
    {key:'controllers',label:'Контроллеры',group:'controller',getWord:'Controller',
     idField:'controller_id',nameField:'controller_name',
     infoFn:i=>i.controller_info||'',
     subFields:[
         {id:'name',dbField:'controller_name',label:'Название',  type:'text'},
         {id:'info',dbField:'controller_info',label:'Информация',type:'textarea'},
     ]},
    /* ── Датчики ─────────────────────────────────────────── */
    {key:'sensors',label:'Датчики',group:'sensor',getWord:'Sensor',
     idField:'sensor_id',nameField:'sensor_name',
     infoFn:i=>i.sensor_info||'',
     subFields:[
         {id:'name',dbField:'sensor_name',label:'Название',  type:'text'},
         {id:'info',dbField:'sensor_info',label:'Информация',type:'textarea'},
     ]},
    /* ── Проводка ────────────────────────────────────────── */
    {key:'wirings',label:'Проводка',group:'wiring',getWord:'Wiring',
     idField:'wiring_id',nameField:'wiring_name',
     infoFn:i=>i.wiring_info||'',
     subFields:[
         {id:'name',dbField:'wiring_name',label:'Название',  type:'text'},
         {id:'info',dbField:'wiring_info',label:'Информация',type:'textarea'},
     ]},
    /* ── Каркасы ─────────────────────────────────────────── */
    {key:'carcasses',label:'Каркасы',group:'carcass',getWord:'Carcass',
     idField:'carcass_id',nameField:'carcass_name',
     infoFn:i=>i.carcass_info||'',
     subFields:[
         {id:'name',dbField:'carcass_name',label:'Название',  type:'text'},
         {id:'info',dbField:'carcass_info',label:'Информация',type:'textarea'},
     ]},
    /* ── Двери ───────────────────────────────────────────── */
    {key:'doors',label:'Двери',group:'doors',getWord:'Doors',
     idField:'doors_id',nameField:'doors_name',
     infoFn:i=>i.doors_info||'',
     subFields:[
         {id:'name',dbField:'doors_name',label:'Название',  type:'text'},
         {id:'info',dbField:'doors_info',label:'Информация',type:'textarea'},
     ]},
    /* ── Крылья ──────────────────────────────────────────── */
    {key:'wings',label:'Крылья',group:'wings',getWord:'Wings',
     idField:'wings_id',nameField:'wings_name',
     infoFn:i=>i.wings_info||'',
     subFields:[
         {id:'name',dbField:'wings_name',label:'Название',  type:'text'},
         {id:'info',dbField:'wings_info',label:'Информация',type:'textarea'},
     ]},
    /* ── Составные ───────────────────────────────────────── */
    {key:'powerPoints',label:'Силовые установки',group:'power-point',getWord:'PowerPoint',composite:true,
     idField:'power_point_id',nameField:null,
     infoFn:i=>`#${shortId(i.engine_id)}`},
    {key:'chassis',label:'Шасси',group:'chassis',getWord:'Chassis',composite:true,
     idField:'chassis_id',nameField:null,
     infoFn:i=>`#${shortId(i.frame_id)}`},
    {key:'chargerSystems',label:'Зарядные системы',group:'charger-system',getWord:'ChargSystem',composite:true,
     idField:'charger_system_id',nameField:null,
     infoFn:i=>`#${shortId(i.charger_id)}`},
    {key:'electronics',label:'Электроника',group:'electronics',getWord:'Electronics',composite:true,
     idField:'electronics_id',nameField:null,
     infoFn:i=>`#${shortId(i.controller_id)}`},
    {key:'bodies',label:'Кузов',group:'body',getWord:'Body',composite:true,
     idField:'body_id',nameField:null,
     infoFn:i=>`#${shortId(i.carcass_id)}`},
    {key:'emobiles',label:'Автомобили',group:'emobile',getWord:'Emobile',composite:true,
     idField:'emobile_id',nameField:'emobile_name',
     infoFn:i=>`#${shortId(i.power_point_id)}`},
];
function getCatConfig(key){return CATEGORY_MAP.find(c=>c.key===key)||null}

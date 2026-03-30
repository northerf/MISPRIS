/* config.js — исправлен под реальные JSON-теги Go-схем */
const API_BASE = 'http://localhost:8080/api';

function escapeHtml(str) {
    if (!str && str !== 0) return '';
    return String(str).replace(/[&<>"']/g, ch => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[ch]));
}
function shortId(id) { return id ? String(id).slice(0,8) : '—'; }

const CATEGORY_MAP = [
    // ── Листовые ────────────────────────────────────────────────────────
    { key:'batteries', label:'Батареи', group:'battery', getWord:'Battery',
      idField:'battery_id', nameField:'battery_name',
      infoFn: i=>[i.battery_type, i.battery_capacity?`${i.battery_capacity} кВтч`:''].filter(Boolean).join(' · '),
      subFields:[
          {id:'name',             label:'Название',       type:'text'},
          {id:'battery_type',     label:'Тип',            type:'text'},
          {id:'battery_capacity', label:'Ёмкость (кВтч)', type:'number'},
          {id:'battery_info',     label:'Информация',     type:'textarea'}
      ]},

    { key:'engines', label:'Двигатели', group:'engine', getWord:'Engine',
      idField:'engine_id', nameField:'engine_name',
      infoFn: i=>i.engine_type||'',
      subFields:[
          {id:'name',        label:'Название',    type:'text'},
          {id:'engine_type', label:'Тип (AC/DC)', type:'text'},
          {id:'info',        label:'Информация',  type:'textarea'}
      ]},

    { key:'inverters', label:'Инверторы', group:'inverter', getWord:'Inverter',
      idField:'inverter_id', nameField:'inverter_name', infoFn: i=>i.inverter_info||'',
      subFields:[
          {id:'name', label:'Название',   type:'text'},
          {id:'info', label:'Информация', type:'textarea'}
      ]},

    { key:'gearboxes', label:'КПП', group:'gearbox', getWord:'Gearbox',
      idField:'gearbox_id', nameField:'gearbox_name', infoFn: i=>i.gearbox_info||'',
      subFields:[
          {id:'name', label:'Название',   type:'text'},
          {id:'info', label:'Информация', type:'textarea'}
      ]},

    { key:'frames', label:'Рамы', group:'frame', getWord:'Frame',
      idField:'frame_id', nameField:'frame_name', infoFn: i=>i.frame_info||'',
      subFields:[
          {id:'name', label:'Название',   type:'text'},
          {id:'info', label:'Информация', type:'textarea'}
      ]},

    { key:'suspensions', label:'Подвески', group:'suspension', getWord:'Suspension',
      idField:'suspension_id', nameField:'suspension_name', infoFn: i=>i.suspension_info||'',
      subFields:[
          {id:'name', label:'Название',   type:'text'},
          {id:'info', label:'Информация', type:'textarea'}
      ]},

    { key:'breakSystems', label:'Тормозные системы', group:'break-system', getWord:'BreakSystem',
      idField:'break_system_id', nameField:'break_system_name', infoFn: i=>i.break_info||'',
      subFields:[
          {id:'name', label:'Название',   type:'text'},
          {id:'info', label:'Информация', type:'textarea'}
      ]},

    { key:'chargers', label:'Зарядные устройства', group:'charger', getWord:'Charger',
      idField:'charger_id', nameField:'charger_name', infoFn: i=>i.charger_info||'',
      subFields:[
          {id:'name', label:'Название',   type:'text'},
          {id:'info', label:'Информация', type:'textarea'}
      ]},

    { key:'connectors', label:'Коннекторы', group:'connector', getWord:'Connector',
      idField:'connector_id', nameField:'connector_name', infoFn: i=>i.connector_info||'',
      subFields:[
          {id:'name', label:'Название',   type:'text'},
          {id:'info', label:'Информация', type:'textarea'}
      ]},

    { key:'controllers', label:'Контроллеры', group:'controller', getWord:'Controller',
      idField:'controller_id', nameField:'controller_name', infoFn: i=>i.controller_info||'',
      subFields:[
          {id:'name', label:'Название',   type:'text'},
          {id:'info', label:'Информация', type:'textarea'}
      ]},

    { key:'sensors', label:'Датчики', group:'sensor', getWord:'Sensor',
      idField:'sensor_id', nameField:'sensor_name', infoFn: i=>i.sensor_info||'',
      subFields:[
          {id:'name', label:'Название',   type:'text'},
          {id:'info', label:'Информация', type:'textarea'}
      ]},

    { key:'wirings', label:'Проводка', group:'wiring', getWord:'Wiring',
      idField:'wiring_id', nameField:'wiring_name', infoFn: i=>i.wiring_info||'',
      subFields:[
          {id:'name', label:'Название',   type:'text'},
          {id:'info', label:'Информация', type:'textarea'}
      ]},

    { key:'carcasses', label:'Каркасы', group:'carcass', getWord:'Carcass',
      idField:'carcass_id', nameField:'carcass_name', infoFn: i=>i.carcass_info||'',
      subFields:[
          {id:'name', label:'Название',   type:'text'},
          {id:'info', label:'Информация', type:'textarea'}
      ]},

    { key:'doors', label:'Двери', group:'doors', getWord:'Doors',
      idField:'doors_id', nameField:'doors_name', infoFn: i=>i.doors_info||'',
      subFields:[
          {id:'name', label:'Название',   type:'text'},
          {id:'info', label:'Информация', type:'textarea'}
      ]},

    { key:'wings', label:'Крылья', group:'wings', getWord:'Wings',
      idField:'wings_id', nameField:'wings_name', infoFn: i=>i.wings_info||'',
      subFields:[
          {id:'name', label:'Название',   type:'text'},
          {id:'info', label:'Информация', type:'textarea'}
      ]},

    // ── Составные ────────────────────────────────────────────────────────
    { key:'powerPoints', label:'Силовые установки', group:'power-point', getWord:'PowerPoint', composite:true,
      idField:'power_point_id', nameField:null,
      infoFn: i=>`Двиг:${shortId(i.engine_id)} · Инв:${shortId(i.inverter_id)} · КПП:${shortId(i.gearbox_id)}`},

    { key:'chassis', label:'Шасси', group:'chassis', getWord:'Chassis', composite:true,
      idField:'chassis_id', nameField:null,
      infoFn: i=>`Рама:${shortId(i.frame_id)} · Подв:${shortId(i.suspension_id)} · Торм:${shortId(i.break_system_id)}`},

    { key:'chargerSystems', label:'Зарядные системы', group:'charger-system', getWord:'ChargSystem', composite:true,
      idField:'charger_system_id', nameField:null,
      infoFn: i=>`Заряд:${shortId(i.charger_id)} · Конн:${shortId(i.connector_id)}`},

    { key:'electronics', label:'Электроника', group:'electronics', getWord:'Electronics', composite:true,
      idField:'electronics_id', nameField:null,
      infoFn: i=>`Контр:${shortId(i.controller_id)} · Сенс:${shortId(i.sensor_id)} · Пров:${shortId(i.wiring_id)}`},

    { key:'bodies', label:'Кузов', group:'body', getWord:'Body', composite:true,
      idField:'body_id', nameField:null,
      infoFn: i=>`Карк:${shortId(i.carcass_id)} · Дв:${shortId(i.doors_id)} · Кр:${shortId(i.wings_id)}`},

    { key:'emobiles', label:'Автомобили', group:'emobile', getWord:'Emobile', composite:true,
      idField:'emobile_id', nameField:'emobile_name',
      infoFn: i=>`СУ:${shortId(i.power_point_id)} · Бат:${shortId(i.battery_id)}`},
];

function getCatConfig(key) { return CATEGORY_MAP.find(c => c.key === key) || null; }

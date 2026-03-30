/* forms.js */
function createModal(html, wide=false) {
    const o=document.createElement('div'); o.className='modal';
    o.innerHTML=`<div class="modal-content${wide?' modal-wide':''}">${html}</div>`;
    o.addEventListener('click',e=>{if(e.target===o)o.remove();});
    document.body.appendChild(o); return o;
}
function showError(msg){
    const m=createModal(`<div class="modal-title" style="color:#f87171"><i class="fas fa-exclamation-circle"></i> Ошибка</div><p>${escapeHtml(msg)}</p><div class="modal-footer"><button class="btn btn-secondary" id="_ok">OK</button></div>`);
    m.querySelector('#_ok').onclick=()=>m.remove();
}

function confirmDelete(id,categoryKey,name){
    const m=createModal(`<div class="modal-title" style="color:#f87171"><i class="fas fa-trash-alt"></i> Удаление</div>
        <p>Вы уверены, что хотите удалить <strong>${escapeHtml(name)}</strong>?</p>
        <div class="modal-footer"><button class="btn btn-secondary" id="_c">Отмена</button><button class="btn btn-danger" id="_d">Удалить</button></div>`);
    m.querySelector('#_c').onclick=()=>m.remove();
    m.querySelector('#_d').onclick=async()=>{
        try{await api[categoryKey].delete(id);await loadCatalogData();}catch(e){showError(e.message);}
        m.remove();
    };
}

async function openPartForm(editId, categoryKey){
    if(!categoryKey){_openCategoryPicker();return;}
    const cat=getCatConfig(categoryKey);
    if(!cat) return;
    if(cat.composite){_openCompositeForm(cat,editId);return;}
    let existing={};
    if(editId){try{const l=await api[cat.key].list();existing=l.find(i=>i[cat.idField]===editId)||{};}catch(_){}}
    const fieldsHtml=cat.subFields.map(f=>{
        const v=escapeHtml(existing[f.id]||existing[cat.nameField]||'');
        const inp=f.type==='textarea'
            ?`<textarea id="_f_${f.id}" class="input-dark" rows="2">${f.id==='name'?escapeHtml(existing[cat.nameField]||''):escapeHtml(existing[f.id]||'')}</textarea>`
            :`<input id="_f_${f.id}" type="${f.type}" class="input-dark" value="${f.id==='name'?escapeHtml(existing[cat.nameField]||''):escapeHtml(existing[f.id]||'')}" placeholder="${escapeHtml(f.label)}">`;
        return `<div class="form-group"><label>${escapeHtml(f.label)}</label>${inp}</div>`;
    }).join('');
    const m=createModal(`
        <div class="modal-title"><i class="fas fa-${editId?'edit':'plus-circle'}" style="color:var(--accent)"></i> ${editId?'Редактировать':'Новая деталь'} — ${escapeHtml(cat.label)}</div>
        ${fieldsHtml}
        <div id="_fe" class="text-danger" style="font-size:0.85rem;min-height:1.2em"></div>
        <div class="modal-footer"><button class="btn btn-secondary" id="_cc">Отмена</button><button class="btn btn-primary" id="_ss"><i class="fas fa-save"></i> Сохранить</button></div>`);
    m.querySelector('#_cc').onclick=()=>m.remove();
    m.querySelector('#_ss').onclick=async()=>{
        const payload={};
        cat.subFields.forEach(f=>{const el=m.querySelector(`#_f_${f.id}`);payload[f.id]=el?el.value.trim():'';});
        if(!payload['name']){m.querySelector('#_fe').textContent='Заполните название';return;}
        try{
            if(editId) await api[cat.key].update(editId,payload);
            else       await api[cat.key].create(payload);
            await loadCatalogData(); m.remove();
        }catch(e){m.querySelector('#_fe').textContent=`Ошибка: ${e.message}`;}
    };
}

function _openCategoryPicker(){
    const leaf=CATEGORY_MAP.filter(c=>!c.composite), comp=CATEGORY_MAP.filter(c=>c.composite);
    const g=arr=>arr.map(c=>`<button class="btn btn-secondary _pc" data-key="${c.key}" style="text-align:left;justify-content:flex-start">${escapeHtml(c.label)}</button>`).join('');
    const m=createModal(`
        <div class="modal-title"><i class="fas fa-list" style="color:var(--accent)"></i> Тип детали</div>
        <div class="section-label">Листовые</div>
        <div style="display:grid;grid-template-columns:1fr 1fr;gap:0.4rem;margin-bottom:1rem">${g(leaf)}</div>
        <div class="section-label">Составные</div>
        <div style="display:grid;grid-template-columns:1fr 1fr;gap:0.4rem;margin-bottom:1rem">${g(comp)}</div>
        <div class="modal-footer"><button class="btn btn-secondary" id="_cx">Отмена</button></div>`,true);
    m.querySelector('#_cx').onclick=()=>m.remove();
    m.querySelectorAll('._pc').forEach(b=>b.addEventListener('click',()=>{m.remove();openPartForm(null,b.dataset.key);}));
}

// Поля составных сущностей — используем реальные JSON-теги Go-схем
const COMPOSITE_DEPS={
    powerPoints:[
        {f:'engine_id',      c:'engines',      l:'Двигатель'},
        {f:'inverter_id',    c:'inverters',    l:'Инвертор'},
        {f:'gearbox_id',     c:'gearboxes',    l:'КПП'},
    ],
    chassis:[
        {f:'frame_id',       c:'frames',       l:'Рама'},
        {f:'suspension_id',  c:'suspensions',  l:'Подвеска'},
        {f:'break_system_id',c:'breakSystems', l:'Тормозная система'},
    ],
    chargerSystems:[
        {f:'charger_id',     c:'chargers',     l:'Зарядное устройство'},
        {f:'connector_id',   c:'connectors',   l:'Коннектор'},
    ],
    electronics:[
        {f:'controller_id',  c:'controllers',  l:'Контроллер'},
        {f:'sensor_id',      c:'sensors',      l:'Датчик'},
        {f:'wiring_id',      c:'wirings',      l:'Проводка'},
    ],
    bodies:[
        {f:'carcass_id',     c:'carcasses',    l:'Каркас'},
        {f:'doors_id',       c:'doors',        l:'Двери'},
        {f:'wings_id',       c:'wings',        l:'Крылья'},
    ],
    emobiles:[
        {f:'name',              c:null,            l:'Название'},
        {f:'power_point_id',    c:'powerPoints',   l:'Силовая установка'},
        {f:'battery_id',        c:'batteries',     l:'Батарея'},
        {f:'charger_system_id', c:'chargerSystems',l:'Зарядная система'},
        {f:'chassis_id',        c:'chassis',       l:'Шасси'},
        {f:'body_id',           c:'bodies',        l:'Кузов'},
        {f:'electronics_id',    c:'electronics',   l:'Электроника'},
    ],
};

async function _openCompositeForm(cat,editId){
    const deps=COMPOSITE_DEPS[cat.key]; if(!deps) return;
    const listsMap={};
    await Promise.all(deps.map(async dep=>{
        if(!dep.c) return;
        try{listsMap[dep.f]=await api[dep.c].list();listsMap[`_c_${dep.f}`]=getCatConfig(dep.c);}catch(_){listsMap[dep.f]=[];}
    }));
    let existing={};
    if(editId){try{const l=await api[cat.key].list();existing=l.find(i=>i[cat.idField]===editId)||{};}catch(_){}}
    const fieldsHtml=deps.map(dep=>{
        if(!dep.c) return `<div class="form-group"><label>${escapeHtml(dep.l)}</label><input id="_f_${dep.f}" type="text" class="input-dark" value="${escapeHtml(existing[dep.f]||existing[cat.nameField]||'')}"></div>`;
        const items=listsMap[dep.f]||[], dc=listsMap[`_c_${dep.f}`];
        const opts=items.map(i=>{const iId=i[dc.idField],iN=dc.nameField?i[dc.nameField]:`#${shortId(iId)}`;return `<option value="${escapeHtml(iId)}" ${existing[dep.f]==iId?'selected':''}>${escapeHtml(iN)}</option>`;}).join('');
        return `<div class="form-group"><label>${escapeHtml(dep.l)}</label><select id="_f_${dep.f}" class="input-dark"><option value="">— выберите —</option>${opts}</select></div>${items.length===0?'<p class="text-muted" style="font-size:0.8rem">Список пуст — сначала создайте эту деталь.</p>':''}`;
    }).join('');
    const m=createModal(`
        <div class="modal-title"><i class="fas fa-puzzle-piece" style="color:var(--accent)"></i> ${editId?'Редактировать':'Создать'} — ${escapeHtml(cat.label)}</div>
        ${fieldsHtml}
        <div id="_fe" class="text-danger" style="font-size:0.85rem;min-height:1.2em"></div>
        <div class="modal-footer"><button class="btn btn-secondary" id="_cc">Отмена</button><button class="btn btn-primary" id="_ss"><i class="fas fa-save"></i> Сохранить</button></div>`,true);
    m.querySelector('#_cc').onclick=()=>m.remove();
    m.querySelector('#_ss').onclick=async()=>{
        const payload={};
        deps.forEach(dep=>{const el=m.querySelector(`#_f_${dep.f}`);payload[dep.f]=el?el.value.trim():'';});
        const miss=deps.find(d=>!payload[d.f]);
        if(miss){m.querySelector('#_fe').textContent=`Заполните поле «${miss.l}»`;return;}
        try{
            if(editId) await api[cat.key].update(editId,payload);
            else       await api[cat.key].create(payload);
            await loadCatalogData(); m.remove();
        }catch(e){m.querySelector('#_fe').textContent=`Ошибка: ${e.message}`;}
    };
}

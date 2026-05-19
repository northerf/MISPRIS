/* forms.js */
function createModal(html,wide=false){
    const o=document.createElement('div');
    o.className='modal';
    o.setAttribute('role','presentation');
    o.innerHTML=`<div class="modal-content${wide?' modal-wide':''}" role="dialog" aria-modal="true" tabindex="-1">${html}</div>`;
    const close=()=>o.remove();
    o.addEventListener('click',e=>{if(e.target===o)close();});
    o.addEventListener('keydown',e=>{if(e.key==='Escape')close();});
    document.body.appendChild(o);
    const focusTarget=o.querySelector('input,select,textarea,button,[tabindex]:not([tabindex="-1"])')||o.querySelector('.modal-content');
    setTimeout(()=>focusTarget&&focusTarget.focus(),0);
    return o;
}
function showToast(msg, ok=true){
    const t=document.createElement('div');
    t.className=`toast ${ok?'toast-ok':'toast-error'}`;
    t.setAttribute('role',ok?'status':'alert');
    t.textContent=msg;
    document.body.appendChild(t);
    setTimeout(()=>t.remove(),4200);
}
function showError(msg){
    const m=createModal(`
        <div class="modal-title" style="color:#f87171">
            <i class="fas fa-exclamation-circle"></i> Ошибка</div>
        <p style="line-height:1.6;color:var(--text)">${escapeHtml(msg)}</p>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_ok">OK</button></div>`);
    m.querySelector('#_ok').onclick=()=>m.remove();
}
function setFormError(container,msg){
    const el=container.querySelector('.error-banner');if(!el)return;
    if(msg){el.querySelector('span').textContent=msg;el.classList.remove('hidden');}
    else el.classList.add('hidden');
}

/* ── Удаление ──────────────────────────────────────────────── */
function confirmDelete(id,categoryKey,name){
    const m=createModal(`
        <div class="modal-title" style="color:#f87171">
            <i class="fas fa-trash-alt"></i> Удаление</div>
        <p>Вы уверены, что хотите удалить <strong>${escapeHtml(name)}</strong>?</p>
        <div class="error-banner hidden"><i class="fas fa-exclamation-circle"></i><span></span></div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_c">Отмена</button>
            <button class="btn btn-danger"    id="_d">
                <i class="fas fa-trash-alt"></i> Удалить</button>
        </div>`);
    m.querySelector('#_c').onclick=()=>m.remove();
    m.querySelector('#_d').onclick=async()=>{
        const btn=m.querySelector('#_d');
        btn.disabled=true;btn.innerHTML='<i class="fas fa-spinner fa-spin"></i>';
        setFormError(m,'');
        try{
            await api[categoryKey].delete(id);
            await loadCatalogData();m.remove();
        }catch(e){
            setFormError(m,e.message);
            btn.disabled=false;
            btn.innerHTML='<i class="fas fa-trash-alt"></i> Удалить';
        }
    };
}

/* ── HTML поля (async из-за select_enum) ──────────────────── */
async function buildFieldHtml(f, val){
    const v = escapeHtml(val ?? '');

    if(f.type === 'textarea')
        return `<textarea id="_f_${f.id}" class="input-dark" rows="2" style="resize:none">${v}</textarea>`;

    // ── ПР2: select_enum — ВСЕГДА идёт в БД, никакого кэша ──
    if(f.type === 'select_enum'){
        const enumVals = await loadEnumCacheForType(f.enumComponent || '');
        if(enumVals && enumVals.length > 0){
            const opts = enumVals.map(ev =>
                `<option value="${escapeHtml(ev.value)}"
                         ${val === ev.value ? 'selected' : ''}>${escapeHtml(ev.value)}</option>`
            ).join('');
            return `<select id="_f_${f.id}" class="input-dark">
                        <option value="">— выберите —</option>${opts}
                    </select>
                    <small style="color:var(--muted);font-size:.75rem;margin-top:.15rem;display:block">
                        📋 Из справочника ПР2 (${f.enumComponent})
                    </small>`;
        }
        // Fallback: статика, если справочник не заведён
        const opts = (f.options||[]).map(o =>
            `<option value="${escapeHtml(o)}" ${val===o?'selected':''}>${escapeHtml(o)}</option>`
        ).join('');
        return `<select id="_f_${f.id}" class="input-dark">
                    <option value="">— выберите —</option>${opts}
                </select>
                <small style="color:#f59e0b;font-size:.75rem;margin-top:.15rem;display:block">
                    ⚠ Справочник «${f.enumComponent}» не найден в ПР2 — используются значения по умолчанию
                </small>`;
    }

    if(f.type === 'select'){
        const opts=(f.options||[]).map(o=>
            `<option value="${escapeHtml(o)}" ${val===o?'selected':''}>${escapeHtml(o)}</option>`
        ).join('');
        return `<select id="_f_${f.id}" class="input-dark">
                    <option value="">— выберите —</option>${opts}</select>`;
    }

    return `<input id="_f_${f.id}" type="${f.type}" class="input-dark"
                value="${v}" placeholder="${escapeHtml(f.label)}">`;
}

/* ── Форма листовой детали ─────────────────────────────────── */
async function openPartForm(editId, categoryKey){
    if(!categoryKey){_openCategoryPicker();return;}
    const cat=getCatConfig(categoryKey);if(!cat)return;
    if(cat.composite){_openCompositeForm(cat,editId);return;}

    let existing={};
    if(editId){
        try{
            const list=await api[cat.key].list();
            existing=list.find(i=>String(i[cat.idField])===String(editId))||{};
        }catch(_){}
    }

    // Строим поля асинхронно (select_enum требует API-вызова)
    const fieldHtmlParts = await Promise.all(cat.subFields.map(async f => {
        const val = f.id==='name'
            ? (existing[cat.nameField]??'')
            : (existing[f.dbField]??'');
        const fieldHtml = await buildFieldHtml(f, val);
        return `<div class="form-group">
                    <label>${escapeHtml(f.label)}</label>
                    ${fieldHtml}
                </div>`;
    }));

    const m=createModal(`
        <div class="modal-title">
            <i class="fas fa-${editId?'edit':'plus-circle'}"
               style="color:var(--accent)"></i>
            ${editId?'Редактировать':'Новая деталь'} — ${escapeHtml(cat.label)}</div>
        ${fieldHtmlParts.join('')}
        <div class="error-banner hidden">
            <i class="fas fa-exclamation-circle"></i><span></span></div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_cc">Отмена</button>
            <button class="btn btn-primary"   id="_ss">
                <i class="fas fa-save"></i> Сохранить</button>
        </div>`);

    m.querySelector('#_cc').onclick=()=>m.remove();
    m.querySelector('#_ss').onclick=async()=>{
        const payload={};
        cat.subFields.forEach(f=>{
            const el=m.querySelector(`#_f_${f.id}`);
            payload[f.id]=el?el.value.trim():'';
        });
        if(!payload['name']){setFormError(m,'Заполните название');return;}
        const btn=m.querySelector('#_ss');
        btn.disabled=true;btn.innerHTML='<i class="fas fa-spinner fa-spin"></i>';
        setFormError(m,'');
        try{
            if(editId) await api[cat.key].update(editId,payload);
            else       await api[cat.key].create(payload);
            await loadCatalogData();m.remove();
        }catch(e){
            setFormError(m,e.message);
            btn.disabled=false;
            btn.innerHTML='<i class="fas fa-save"></i> Сохранить';
        }
    };
}

/* ── Выбор типа детали ─────────────────────────────────────── */
function _openCategoryPicker(){
    const leaf=CATEGORY_MAP.filter(c=>!c.composite);
    const comp=CATEGORY_MAP.filter(c=> c.composite);
    const g=arr=>arr.map(c=>
        `<button class="btn btn-secondary _pc" data-key="${c.key}"
             style="text-align:left;justify-content:flex-start">
             ${escapeHtml(c.label)}</button>`).join('');
    const m=createModal(`
        <div class="modal-title">
            <i class="fas fa-list" style="color:var(--accent)"></i> Тип детали</div>
        <div class="section-label">Листовые</div>
        <div style="display:grid;grid-template-columns:1fr 1fr;gap:.4rem;margin-bottom:1rem">
            ${g(leaf)}</div>
        <div class="section-label">Составные</div>
        <div style="display:grid;grid-template-columns:1fr 1fr;gap:.4rem;margin-bottom:1rem">
            ${g(comp)}</div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_cx">Отмена</button>
        </div>`,true);
    m.querySelector('#_cx').onclick=()=>m.remove();
    m.querySelectorAll('._pc').forEach(b=>
        b.addEventListener('click',()=>{m.remove();openPartForm(null,b.dataset.key);}));
}

/* ── Зависимости составных деталей ────────────────────────── */
const COMPOSITE_DEPS={
    powerPoints:[
        {f:'engine_id',   c:'engines',     l:'Двигатель'},
        {f:'inverter_id', c:'inverters',   l:'Инвертор'},
        {f:'gearbox_id',  c:'gearboxes',   l:'КПП'},
    ],
    chassis:[
        {f:'frame_id',       c:'frames',       l:'Рама'},
        {f:'suspension_id',  c:'suspensions',  l:'Подвеска'},
        {f:'break_system_id',c:'breakSystems', l:'Тормозная система'},
    ],
    chargerSystems:[
        {f:'charger_id',   c:'chargers',   l:'Зарядное устройство'},
        {f:'connector_id', c:'connectors', l:'Коннектор'},
    ],
    electronics:[
        {f:'controller_id',c:'controllers',l:'Контроллер'},
        {f:'sensor_id',    c:'sensors',    l:'Датчик'},
        {f:'wiring_id',    c:'wirings',    l:'Проводка'},
    ],
    bodies:[
        {f:'carcass_id',c:'carcasses',l:'Каркас'},
        {f:'doors_id',  c:'doors',    l:'Двери'},
        {f:'wings_id',  c:'wings',    l:'Крылья'},
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
    const deps=COMPOSITE_DEPS[cat.key];if(!deps)return;
    const listsMap={};
    await Promise.all(deps.map(async dep=>{
        if(!dep.c)return;
        try{
            listsMap[dep.f]         =await api[dep.c].list();
            listsMap[`_c_${dep.f}`] =getCatConfig(dep.c);
        }catch(_){listsMap[dep.f]=[];}
    }));

    let existing={};
    if(editId){
        try{
            const list=await api[cat.key].list();
            existing=list.find(i=>String(i[cat.idField])===String(editId))||{};
        }catch(_){}
    }

    const fieldsHtml=deps.map(dep=>{
        if(!dep.c){
            const v=escapeHtml(existing[dep.f]||existing[cat.nameField]||'');
            return`<div class="form-group"><label>${escapeHtml(dep.l)}</label>
                   <input id="_f_${dep.f}" type="text" class="input-dark" value="${v}"></div>`;
        }
        const items=listsMap[dep.f]||[];
        const dc   =listsMap[`_c_${dep.f}`];
        const opts =items.map(i=>{
            const iId=i[dc.idField];
            const iN =dc.nameField?i[dc.nameField]:`${dc.label} #${shortId(iId)}`;
            return`<option value="${escapeHtml(iId)}"
                           ${String(existing[dep.f])===String(iId)?'selected':''}>
                       ${escapeHtml(iN)}</option>`;
        }).join('');
        const hint=items.length===0
            ?`<p style="font-size:.8rem;color:var(--muted);margin-top:.25rem">
               ⚠ Список пуст — сначала создайте эту деталь.</p>`:'';
        return`<div class="form-group"><label>${escapeHtml(dep.l)}</label>
               <select id="_f_${dep.f}" class="input-dark">
                   <option value="">— выберите —</option>${opts}
               </select>${hint}</div>`;
    }).join('');

    const m=createModal(`
        <div class="modal-title">
            <i class="fas fa-puzzle-piece" style="color:var(--accent)"></i>
            ${editId?'Редактировать':'Создать'} — ${escapeHtml(cat.label)}</div>
        ${fieldsHtml}
        <div class="error-banner hidden">
            <i class="fas fa-exclamation-circle"></i><span></span></div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_cc">Отмена</button>
            <button class="btn btn-primary"   id="_ss">
                <i class="fas fa-save"></i> Сохранить</button>
        </div>`,true);

    m.querySelector('#_cc').onclick=()=>m.remove();
    m.querySelector('#_ss').onclick=async()=>{
        const payload={};
        deps.forEach(dep=>{
            const el=m.querySelector(`#_f_${dep.f}`);
            payload[dep.f]=el?el.value.trim():'';
        });
        const miss=deps.find(d=>!payload[d.f]);
        if(miss){setFormError(m,`Заполните поле «${miss.l}»`);return;}
        const btn=m.querySelector('#_ss');
        btn.disabled=true;btn.innerHTML='<i class="fas fa-spinner fa-spin"></i>';
        setFormError(m,'');
        try{
            let newId;
            if(editId){
                await api[cat.key].update(editId,payload);
            }else{
                const resp = await api[cat.key].create(payload);
                newId = resp?.id || resp?.ID;
            }
            await loadCatalogData();
            m.remove();
            // ── ПР3: после создания автомобиля открываем форму параметров ──
            if(!editId && cat.key === 'emobiles' && newId){
                openEmobileParamsForm(newId, payload['name'] || 'новый автомобиль');
            }
        }catch(e){
            setFormError(m,e.message);
            btn.disabled=false;
            btn.innerHTML='<i class="fas fa-save"></i> Сохранить';
        }
    };
}

/* ── ПР3: Параметры конкретного автомобиля ─────────────────── */
async function openEmobileParamsForm(emobileId, emobileName){
    // Показываем загрузку
    const loadingModal = createModal(`
        <div class="modal-title">
            <i class="fas fa-sliders-h" style="color:var(--accent)"></i>
            Параметры: ${escapeHtml(emobileName)}
        </div>
        <div style="text-align:center;padding:2rem;color:var(--muted)">
            <i class="fas fa-spinner fa-spin"></i> Загрузка параметров…
        </div>`, true);

    let compParams = [];
    try{
        // Загружаем параметры типа 'emobile' через SQL-функцию get_component_parameters
        const url = `${API_BASE}/component-parameter/byTypeemobile`;
        compParams = await apiRequest(url) || [];
    }catch(e){
        loadingModal.remove();
        showError(`Ошибка загрузки параметров: ${e.message}`);
        return;
    }

    // Если параметры не настроены — показываем подсказку
    if(!compParams || compParams.length === 0){
        loadingModal.remove();
        const hint = createModal(`
            <div class="modal-title">
                <i class="fas fa-sliders-h" style="color:var(--accent)"></i>
                Параметры (ПР3)
            </div>
            <div style="padding:.5rem 0;color:var(--muted);line-height:1.7">
                <p>Для автомобиля <strong>${escapeHtml(emobileName)}</strong> ещё не настроены параметры ПР3.</p>
                <p style="margin-top:.8rem">Чтобы добавить параметры:</p>
                <ol style="margin-left:1.2rem;margin-top:.4rem">
                    <li>Нажмите кнопку <strong>Параметры</strong> в главном меню</li>
                    <li>Вкладка <strong>Описания параметров</strong> → создать параметр<br>
                        (например: «Запас хода», тип real, единица км)</li>
                    <li>Вкладка <strong>Параметры компонентов</strong> → привязать параметр,<br>
                        тип компонента: <code>emobile</code></li>
                    <li>После этого кнопка <i class="fas fa-sliders-h"></i> покажет форму заполнения</li>
                </ol>
            </div>
            <div class="modal-footer">
                <button class="btn btn-secondary" id="_phClose">Понятно</button>
                <button class="btn btn-primary" id="_phParams">Открыть Параметры</button>
            </div>`);
        hint.querySelector('#_phClose').onclick = () => hint.remove();
        hint.querySelector('#_phParams').onclick = () => { hint.remove(); openParametersManager(); };
        return;
    }

    // Загружаем уже заполненные значения
    let existingVals = {};
    try{
        const vals = await apiRequest(`${API_BASE}/emobile-parameter/byEmobile${emobileId}`) || [];
        vals.forEach(v => {
            // Ключ — component_parameter_id
            existingVals[String(v.component_parameter_id)] = v;
        });
    }catch(_){}

    loadingModal.remove();

    // Строим поля формы для каждого параметра
    const buildParamField = async (cp) => {
        const cpId = String(cp.cp_id || cp.component_parameter_id);
        const existing = existingVals[cpId] || {};
        const unitStr = cp.measuring_unit ? ` (${cp.measuring_unit})` : '';
        const rangeStr = (cp.min_val || cp.max_val)
            ? `<span style="font-size:.72rem;color:var(--faint)"> [${cp.min_val}–${cp.max_val}]</span>`
            : '';

        if(cp.param_type === 'real' || cp.param_type === 'int'){
            const curVal = cp.param_type === 'real'
                ? (existing.val_real !== undefined && existing.val_real !== 0 ? existing.val_real : '')
                : (existing.val_int  !== undefined && existing.val_int  !== 0 ? existing.val_int  : '');
            return `<div class="form-group">
                <label>${escapeHtml(cp.name + unitStr)}${rangeStr}</label>
                <input id="_pp_${cpId}" type="number" step="${cp.param_type==='int'?'1':'any'}"
                       class="input-dark" value="${curVal}"
                       placeholder="${cp.min_val||''}–${cp.max_val||''}">
            </div>`;
        }
        if(cp.param_type === 'str'){
            return `<div class="form-group">
                <label>${escapeHtml(cp.name)}</label>
                <input id="_pp_${cpId}" type="text" class="input-dark"
                       value="${escapeHtml(existing.val_str||'')}">
            </div>`;
        }
        if(cp.param_type === 'enum'){
            // Загружаем параметр чтобы узнать его enum_class_id
            let enumOpts = '';
            try{
                const param = await apiRequest(`${API_BASE}/parameter/getParameter${cp.param_id}`);
                if(param && param.enum_class_id){
                    const vals = await apiRequest(`${API_BASE}/enum-class/values${param.enum_class_id}`) || [];
                    enumOpts = vals.map(ev =>
                        `<option value="${ev.enum_position_id}"
                                 ${String(existing.enum_val_id) === String(ev.enum_position_id) ? 'selected' : ''}>
                             ${escapeHtml(ev.value)}</option>`
                    ).join('');
                }
            }catch(e){ console.warn('[pr3 enum field]', e.message); }
            return `<div class="form-group">
                <label>${escapeHtml(cp.name)}</label>
                <select id="_pp_${cpId}" class="input-dark">
                    <option value="">— выберите —</option>${enumOpts}
                </select>
            </div>`;
        }
        return '';
    };

    const fieldHtmls = await Promise.all(compParams.map(buildParamField));

    const m = createModal(`
        <div class="modal-title">
            <i class="fas fa-sliders-h" style="color:var(--accent)"></i>
            Параметры: ${escapeHtml(emobileName)} (ПР3)
        </div>
        <p style="font-size:.82rem;color:var(--muted);margin-bottom:.8rem">
            Характеристики автомобиля. Оставьте пустым чтобы не менять.
        </p>
        ${fieldHtmls.join('')}
        <div class="error-banner hidden"><i class="fas fa-exclamation-circle"></i><span></span></div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_ppCancel">Отмена</button>
            <button class="btn btn-primary" id="_ppSave">
                <i class="fas fa-save"></i> Сохранить параметры
            </button>
        </div>`, true);

    m.querySelector('#_ppCancel').onclick = () => m.remove();
    m.querySelector('#_ppSave').onclick = async () => {
        const btn = m.querySelector('#_ppSave');
        btn.disabled = true;
        btn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Сохранение…';
        setFormError(m, '');

        const errors = [];
        for(const cp of compParams){
            const cpId = String(cp.cp_id || cp.component_parameter_id);
            const el = m.querySelector(`#_pp_${cpId}`);
            if(!el || el.value === '') continue;  // пропускаем пустые

            const body = {
                emobile_id:             String(emobileId),
                component_parameter_id: cpId,
            };
            if(cp.param_type === 'real')      body.val_real  = el.value;
            else if(cp.param_type === 'int')  body.val_int   = el.value;
            else if(cp.param_type === 'str')  body.val_str   = el.value;
            else if(cp.param_type === 'enum') body.enum_val_id = el.value;

            try{
                const existRecord = existingVals[cpId];
                if(existRecord && existRecord.value_id){
                    // Обновляем существующее значение
                    await api.emobileParam.update(String(existRecord.value_id), body);
                }else{
                    // Создаём новое
                    await api.emobileParam.create(body);
                }
            }catch(e){
                errors.push(`${cp.name}: ${e.message}`);
            }
        }

        if(errors.length){
            setFormError(m, errors.join('; '));
            btn.disabled = false;
            btn.innerHTML = '<i class="fas fa-save"></i> Сохранить параметры';
        }else{
            m.remove();
        }
    };
}

/* ── Пользователь ──────────────────────────────────────────── */
function openCreateUserForm(){
    const m=createModal(`
        <div class="modal-title">
            <i class="fas fa-user-plus" style="color:var(--accent)"></i>
            Новый пользователь</div>
        <div class="form-group"><label>Логин</label>
            <input id="_cu" type="text" class="input-dark"
                placeholder="username" autocomplete="off"></div>
        <div class="form-group"><label>Пароль</label>
            <input id="_cp" type="password" class="input-dark"
                placeholder="password" autocomplete="new-password"></div>
        <div class="form-group"><label>Роль</label>
            <select id="_cr" class="input-dark">
                <option value="user">user — обычный пользователь</option>
                <option value="admin">admin — администратор</option>
            </select></div>
        <div class="error-banner hidden">
            <i class="fas fa-exclamation-circle"></i><span></span></div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_ccc">Отмена</button>
            <button class="btn btn-primary"   id="_css">
                <i class="fas fa-save"></i> Создать</button>
        </div>`);
    m.querySelector('#_ccc').onclick=()=>m.remove();
    m.querySelector('#_css').onclick=async()=>{
        const username=m.querySelector('#_cu').value.trim();
        const password=m.querySelector('#_cp').value.trim();
        const role    =m.querySelector('#_cr').value;
        if(!username||!password){setFormError(m,'Заполните все поля');return;}
        const btn=m.querySelector('#_css');
        btn.disabled=true;btn.innerHTML='<i class="fas fa-spinner fa-spin"></i>';
        setFormError(m,'');
        try{
            await api.users.create({username,password,role});
            m.remove();
        }catch(e){
            setFormError(m,e.message);
            btn.disabled=false;
            btn.innerHTML='<i class="fas fa-save"></i> Создать';
        }
    };
}

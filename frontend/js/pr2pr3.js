/* pr2pr3.js — управление перечислениями (ПР2) и параметрами (ПР3) */

// ══════════════════════════════════════════════════════════════════
//  ОБЩИЕ УТИЛИТЫ (используют функции из forms.js)
// ══════════════════════════════════════════════════════════════════

function showToast(msg, ok=true){
    const t=document.createElement('div');
    t.className=`toast ${ok?'toast-ok':'toast-error'}`;
    t.setAttribute('role',ok?'status':'alert');
    t.textContent=msg;
    document.body.appendChild(t);
    setTimeout(()=>t.remove(),4200);
}

// ══════════════════════════════════════════════════════════════════
//  ПР2 — ПЕРЕЧИСЛЕНИЯ
//  Открывается через кнопку «Справочники» в тулбаре
// ══════════════════════════════════════════════════════════════════

async function openEnumsManager(){
    const modal=createModal(`
        <div class="modal-title">
            <i class="fas fa-list-ul" style="color:var(--accent)"></i>
            Справочники перечислений (ПР2)
        </div>
        <div style="display:flex;gap:.5rem;margin-bottom:.8rem;flex-wrap:wrap">
            <button class="btn btn-primary" id="_ecNew">
                <i class="fas fa-plus"></i> Новый справочник</button>
            <button class="btn btn-secondary" id="_ecRefresh">
                <i class="fas fa-sync"></i> Обновить</button>
        </div>
        <div id="_ecList" style="min-height:120px">
            <i class="fas fa-spinner fa-spin"></i> Загрузка…
        </div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_ecClose">Закрыть</button>
        </div>`,true);

    modal.querySelector('#_ecClose').onclick=()=>modal.remove();
    modal.querySelector('#_ecNew').onclick=()=>openEnumClassForm(null, ()=>refreshEnumList(modal));
    modal.querySelector('#_ecRefresh').onclick=()=>refreshEnumList(modal);

    await refreshEnumList(modal);
}

async function refreshEnumList(modal){
    const el=modal.querySelector('#_ecList');
    el.innerHTML='<i class="fas fa-spinner fa-spin"></i> Загрузка…';
    try{
        const list=await api.enumClass.list()||[];
        if(!list.length){el.innerHTML='<p style="color:var(--muted)">Нет справочников</p>';return;}
        el.innerHTML=list.map(ec=>`
            <div style="border:1px solid var(--border);border-radius:.4rem;padding:.6rem .8rem;
                        margin-bottom:.4rem;display:flex;align-items:center;gap:.5rem">
                <span style="flex:1;font-weight:500">${escapeHtml(ec.name)}</span>
                ${ec.component_type?`<span class="badge badge-default">${escapeHtml(ec.component_type)}</span>`:''}
                <i class="fas fa-eye action-icon" style="color:#60a5fa" title="Позиции"
                   data-ecid="${ec.enum_class_id}" data-ecname="${escapeHtml(ec.name)}"></i>
                <i class="fas fa-edit action-icon" style="color:#a78bfa" title="Редактировать"
                   data-edit-ec="${ec.enum_class_id}"></i>
                <i class="fas fa-trash-alt action-icon" style="color:#f87171" title="Удалить"
                   data-del-ec="${ec.enum_class_id}" data-del-name="${escapeHtml(ec.name)}"></i>
            </div>`).join('');

        el.querySelectorAll('[data-ecid]').forEach(btn=>
            btn.addEventListener('click',()=>
                openEnumPositionsManager(btn.dataset.ecid, btn.dataset.ecname)));
        el.querySelectorAll('[data-edit-ec]').forEach(btn=>
            btn.addEventListener('click',()=>
                openEnumClassForm(btn.dataset.editEc, ()=>refreshEnumList(modal))));
        el.querySelectorAll('[data-del-ec]').forEach(btn=>
            btn.addEventListener('click',async()=>{
                if(!confirm(`Удалить справочник «${btn.dataset.delName}»?`))return;
                try{
                    await api.enumClass.delete(btn.dataset.delEc);
                    showToast('Справочник удалён');
                    await refreshEnumList(modal);
                }catch(e){showToast(e.message,false);}
            }));
    }catch(e){el.innerHTML=`<p style="color:#f87171">${escapeHtml(e.message)}</p>`;}
}

function openEnumClassForm(id, onSave){
    const m=createModal(`
        <div class="modal-title">
            <i class="fas fa-${id?'edit':'plus-circle'}" style="color:var(--accent)"></i>
            ${id?'Редактировать':'Новый'} справочник
        </div>
        <div class="form-group">
            <label>Название *</label>
            <input id="_ecn" class="input-dark" placeholder="Стандарт зарядки">
        </div>
        <div class="form-group">
            <label>Тип компонента</label>
            <select id="_ect" class="input-dark">
                <option value="">— общий —</option>
                <option value="emobile">Автомобиль</option>
                <option value="battery">Батарея</option>
                <option value="engine">Двигатель</option>
                <option value="chassis">Шасси</option>
                <option value="charger">Зарядное устройство</option>
                <option value="electronics">Электроника</option>
                <option value="body">Кузов</option>
            </select>
        </div>
        <div class="error-banner hidden"><i class="fas fa-exclamation-circle"></i><span></span></div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_cc">Отмена</button>
            <button class="btn btn-primary" id="_cs">
                <i class="fas fa-save"></i> Сохранить</button>
        </div>`);

    if(id){
        api.enumClass.getById(id).then(ec=>{
            m.querySelector('#_ecn').value=ec.name||'';
            m.querySelector('#_ect').value=ec.component_type||'';
        }).catch(()=>{});
    }

    m.querySelector('#_cc').onclick=()=>m.remove();
    m.querySelector('#_cs').onclick=async()=>{
        const name=m.querySelector('#_ecn').value.trim();
        if(!name){setFormError(m,'Введите название');return;}
        const btn=m.querySelector('#_cs');
        btn.disabled=true;btn.innerHTML='<i class="fas fa-spinner fa-spin"></i>';
        setFormError(m,'');
        try{
            const data={name, component_type:m.querySelector('#_ect').value};
            if(id) await api.enumClass.update(id,data);
            else   await api.enumClass.create(data);
            showToast(id?'Справочник обновлён':'Справочник создан');
            m.remove();
            if(onSave) onSave();
        }catch(e){
            setFormError(m,e.message);
            btn.disabled=false;btn.innerHTML='<i class="fas fa-save"></i> Сохранить';
        }
    };
}

async function openEnumPositionsManager(enumClassId, enumClassName){
    const modal=createModal(`
        <div class="modal-title">
            <i class="fas fa-tags" style="color:var(--accent)"></i>
            Значения: ${escapeHtml(enumClassName)}
        </div>
        <div style="display:flex;gap:.5rem;margin-bottom:.8rem;flex-wrap:wrap">
            <button class="btn btn-primary" id="_epNew">
                <i class="fas fa-plus"></i> Добавить значение</button>
            <button class="btn btn-secondary" id="_epVals">
                <i class="fas fa-sort-numeric-up"></i> По порядку (SQL)</button>
            <button class="btn btn-outline" id="_epValidate" style="margin-left:auto;border:1px solid var(--border)">
                Проверить значение</button>
        </div>
        <div id="_epList" style="min-height:100px">
            <i class="fas fa-spinner fa-spin"></i> Загрузка…
        </div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_epClose">Закрыть</button>
        </div>`,true);

    modal.querySelector('#_epClose').onclick=()=>modal.remove();
    modal.querySelector('#_epNew').onclick=()=>
        openEnumPositionForm(null, enumClassId, ()=>refreshPositionList(modal, enumClassId));
    modal.querySelector('#_epVals').onclick=async()=>{
        try{
            const vals=await api.enumClass.getValues(enumClassId);
            const info=vals.map(v=>`#${v.order_num}: ${v.value}`).join('\n');
            alert('Значения по порядку:\n\n'+info);
        }catch(e){showToast(e.message,false);}
    };
    modal.querySelector('#_epValidate').onclick=()=>openValidateForm(enumClassId);

    await refreshPositionList(modal, enumClassId);
}

async function refreshPositionList(modal, enumClassId){
    const el=modal.querySelector('#_epList');
    el.innerHTML='<i class="fas fa-spinner fa-spin"></i>';
    try{
        const list=(await api.enumPosition.list()||[])
            .filter(p=>String(p.enum_class_id)===String(enumClassId))
            .sort((a,b)=>a.order_num-b.order_num);
        if(!list.length){el.innerHTML='<p style="color:var(--muted)">Нет значений</p>';return;}
        el.innerHTML=list.map(p=>`
            <div style="border:1px solid var(--border);border-radius:.4rem;padding:.5rem .8rem;
                        margin-bottom:.3rem;display:flex;align-items:center;gap:.5rem">
                <span style="width:2rem;text-align:center;color:var(--muted);font-size:.8rem">
                    ${p.order_num}</span>
                <span style="flex:1">${escapeHtml(p.value)}</span>
                <i class="fas fa-edit action-icon" style="color:#a78bfa" title="Редактировать"
                   data-edit-ep="${p.enum_position_id}"></i>
                <i class="fas fa-trash-alt action-icon" style="color:#f87171" title="Удалить"
                   data-del-ep="${p.enum_position_id}" data-val="${escapeHtml(p.value)}"></i>
            </div>`).join('');

        el.querySelectorAll('[data-edit-ep]').forEach(btn=>
            btn.addEventListener('click',()=>
                openEnumPositionForm(btn.dataset.editEp, enumClassId,
                    ()=>refreshPositionList(modal, enumClassId))));
        el.querySelectorAll('[data-del-ep]').forEach(btn=>
            btn.addEventListener('click',async()=>{
                if(!confirm(`Удалить «${btn.dataset.val}»?`))return;
                try{
                    await api.enumPosition.delete(btn.dataset.delEp);
                    showToast('Значение удалено');
                    await refreshPositionList(modal, enumClassId);
                }catch(e){showToast(e.message,false);}
            }));
    }catch(e){el.innerHTML=`<p style="color:#f87171">${escapeHtml(e.message)}</p>`;}
}

function openEnumPositionForm(id, enumClassId, onSave){
    const m=createModal(`
        <div class="modal-title">
            <i class="fas fa-${id?'edit':'plus'}" style="color:var(--accent)"></i>
            ${id?'Редактировать':'Новое'} значение
        </div>
        <div class="form-group">
            <label>Значение *</label>
            <input id="_epv" class="input-dark" placeholder="CCS2">
        </div>
        <div class="form-group">
            <label>Порядок</label>
            <input id="_epo" class="input-dark" type="number" value="0">
        </div>
        <div class="error-banner hidden"><i class="fas fa-exclamation-circle"></i><span></span></div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_epc">Отмена</button>
            <button class="btn btn-primary" id="_eps">
                <i class="fas fa-save"></i> Сохранить</button>
        </div>`);

    if(id){
        api.enumPosition.getById(id).then(p=>{
            m.querySelector('#_epv').value=p.value||'';
            m.querySelector('#_epo').value=p.order_num||0;
        }).catch(()=>{});
    }

    m.querySelector('#_epc').onclick=()=>m.remove();
    m.querySelector('#_eps').onclick=async()=>{
        const value=m.querySelector('#_epv').value.trim();
        if(!value){setFormError(m,'Введите значение');return;}
        const btn=m.querySelector('#_eps');
        btn.disabled=true;btn.innerHTML='<i class="fas fa-spinner fa-spin"></i>';
        setFormError(m,'');
        try{
            const data={
                enum_class_id: enumClassId,
                value,
                order_num: m.querySelector('#_epo').value||'0',
            };
            if(id) await api.enumPosition.update(id,data);
            else   await api.enumPosition.create(data);
            showToast(id?'Значение обновлено':'Значение добавлено');
            m.remove();
            if(onSave) onSave();
        }catch(e){
            setFormError(m,e.message);
            btn.disabled=false;btn.innerHTML='<i class="fas fa-save"></i> Сохранить';
        }
    };
}

function openValidateForm(enumClassId){
    const m=createModal(`
        <div class="modal-title">
            <i class="fas fa-check-circle" style="color:var(--accent)"></i>
            Проверить значение
        </div>
        <div class="form-group">
            <label>Значение для проверки</label>
            <input id="_vv" class="input-dark" placeholder="CCS2">
        </div>
        <div id="_vr" style="margin:.6rem 0;font-size:.9rem"></div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_vc">Закрыть</button>
            <button class="btn btn-primary" id="_vb">Проверить</button>
        </div>`);

    m.querySelector('#_vc').onclick=()=>m.remove();
    m.querySelector('#_vb').onclick=async()=>{
        const val=m.querySelector('#_vv').value.trim();
        if(!val)return;
        try{
            const res=await api.enumClass.validate({enum_class_id:enumClassId, value:val});
            const el=m.querySelector('#_vr');
            if(res.valid){
                el.innerHTML='<span style="color:#4ade80">✓ Значение допустимо</span>';
            }else{
                el.innerHTML='<span style="color:#f87171">✗ Значение не найдено</span>';
            }
        }catch(e){
            m.querySelector('#_vr').innerHTML=`<span style="color:#f87171">${escapeHtml(e.message)}</span>`;
        }
    };
}


// ══════════════════════════════════════════════════════════════════
//  ПР3 — ПАРАМЕТРЫ ИЗДЕЛИЙ
//  Открывается через кнопку «Параметры» в тулбаре
// ══════════════════════════════════════════════════════════════════

async function openParametersManager(){
    const modal=createModal(`
        <div class="modal-title">
            <i class="fas fa-sliders-h" style="color:var(--accent)"></i>
            Параметры изделий (ПР3)
        </div>

        <!-- Вкладки -->
        <div style="display:flex;gap:.3rem;margin-bottom:.8rem;border-bottom:1px solid var(--border);padding-bottom:.4rem">
            <button class="btn btn-primary _ptab" data-tab="params" style="font-size:.82rem">
                Описания параметров</button>
            <button class="btn btn-secondary _ptab" data-tab="comp" style="font-size:.82rem">
                Параметры компонентов</button>
            <button class="btn btn-secondary _ptab" data-tab="vals" style="font-size:.82rem">
                Значения автомобилей</button>
        </div>

        <!-- Панели -->
        <div id="_ptabs">
            <div id="_pt-params">
                <div style="display:flex;gap:.4rem;margin-bottom:.6rem">
                    <button class="btn btn-primary" id="_pNew">
                        <i class="fas fa-plus"></i> Новый параметр</button>
                    <button class="btn btn-secondary" id="_pRef">
                        <i class="fas fa-sync"></i> Обновить</button>
                </div>
                <div id="_pList"><i class="fas fa-spinner fa-spin"></i></div>
            </div>
            <div id="_pt-comp" style="display:none">
                <div style="display:flex;gap:.4rem;margin-bottom:.6rem;flex-wrap:wrap">
                    <button class="btn btn-primary" id="_cpNew">
                        <i class="fas fa-plus"></i> Привязать параметр</button>
                    <button class="btn btn-secondary" id="_cpRef">
                        <i class="fas fa-sync"></i> Обновить</button>
                    <button class="btn btn-secondary" id="_cpCopy" style="margin-left:auto">
                        <i class="fas fa-copy"></i> Скопировать тип</button>
                </div>
                <div style="margin-bottom:.5rem">
                    <input id="_cpTypeFilter" class="input-dark" style="max-width:240px"
                           placeholder="Фильтр по типу: emobile, battery…">
                    <button class="btn btn-secondary" style="margin-left:.3rem" id="_cpFilterBtn">
                        Фильтр (SQL)</button>
                </div>
                <div id="_cpList"><i class="fas fa-spinner fa-spin"></i></div>
            </div>
            <div id="_pt-vals" style="display:none">
                <div style="display:flex;gap:.4rem;margin-bottom:.6rem;flex-wrap:wrap">
                    <button class="btn btn-primary" id="_epNew">
                        <i class="fas fa-plus"></i> Новое значение</button>
                    <button class="btn btn-secondary" id="_epRef">
                        <i class="fas fa-sync"></i> Обновить</button>
                </div>
                <div style="margin-bottom:.5rem">
                    <input id="_epEmobileFilter" class="input-dark" style="max-width:220px"
                           placeholder="ID автомобиля для фильтра">
                    <button class="btn btn-secondary" style="margin-left:.3rem" id="_epFilterBtn">
                        По автомобилю</button>
                </div>
                <div id="_epvList"><i class="fas fa-spinner fa-spin"></i></div>
            </div>
        </div>

        <div class="modal-footer">
            <button class="btn btn-secondary" id="_pmClose">Закрыть</button>
        </div>`,true);

    modal.querySelector('#_pmClose').onclick=()=>modal.remove();

    // Переключение вкладок
    modal.querySelectorAll('._ptab').forEach(tab=>{
        tab.addEventListener('click',()=>{
            modal.querySelectorAll('._ptab').forEach(t=>{
                t.className=t===tab?'btn btn-primary _ptab':'btn btn-secondary _ptab';
                t.style.fontSize='.82rem';
            });
            modal.querySelectorAll('#_pt-params,#_pt-comp,#_pt-vals').forEach(p=>p.style.display='none');
            modal.querySelector(`#_pt-${tab.dataset.tab}`).style.display='block';
        });
    });

    // ── Параметры ────────────────────────────────────────────────
    const refreshParams=async()=>{
        const el=modal.querySelector('#_pList');
        el.innerHTML='<i class="fas fa-spinner fa-spin"></i>';
        try{
            const list=await api.parameter.list()||[];
            if(!list.length){el.innerHTML='<p style="color:var(--muted)">Нет параметров</p>';return;}
            el.innerHTML=list.map(p=>`
                <div style="border:1px solid var(--border);border-radius:.4rem;padding:.5rem .8rem;
                            margin-bottom:.3rem;display:flex;align-items:center;gap:.5rem">
                    <code style="min-width:80px;font-size:.8rem;color:var(--accent)">${escapeHtml(p.designation)}</code>
                    <span style="flex:1">${escapeHtml(p.name)}</span>
                    <span class="badge badge-default">${escapeHtml(p.param_type)}</span>
                    ${p.measuring_unit?`<span style="font-size:.78rem;color:var(--muted)">${escapeHtml(p.measuring_unit)}</span>`:''}
                    <i class="fas fa-edit action-icon" style="color:#a78bfa" title="Редактировать"
                       data-edit-p="${p.parameter_id}"></i>
                    <i class="fas fa-trash-alt action-icon" style="color:#f87171" title="Удалить"
                       data-del-p="${p.parameter_id}" data-pname="${escapeHtml(p.name)}"></i>
                </div>`).join('');

            el.querySelectorAll('[data-edit-p]').forEach(btn=>
                btn.addEventListener('click',()=>openParameterForm(btn.dataset.editP, refreshParams)));
            el.querySelectorAll('[data-del-p]').forEach(btn=>
                btn.addEventListener('click',async()=>{
                    if(!confirm(`Удалить параметр «${btn.dataset.pname}»?`))return;
                    try{await api.parameter.delete(btn.dataset.delP);showToast('Параметр удалён');await refreshParams();}
                    catch(e){showToast(e.message,false);}
                }));
        }catch(e){el.innerHTML=`<p style="color:#f87171">${escapeHtml(e.message)}</p>`;}
    };
    modal.querySelector('#_pNew').onclick=()=>openParameterForm(null, refreshParams);
    modal.querySelector('#_pRef').onclick=refreshParams;
    await refreshParams();

    // ── Параметры компонентов ────────────────────────────────────
    const refreshCompParams=async(filterType=null)=>{
        const el=modal.querySelector('#_cpList');
        el.innerHTML='<i class="fas fa-spinner fa-spin"></i>';
        try{
            let list;
            if(filterType){
                // SQL-функция get_component_parameters
                list=await api.componentParameter.byType(filterType)||[];
            }else{
                list=await api.componentParameter.list()||[];
            }
            if(!list.length){el.innerHTML='<p style="color:var(--muted)">Нет данных</p>';return;}
            el.innerHTML=list.map(cp=>`
                <div style="border:1px solid var(--border);border-radius:.4rem;padding:.5rem .8rem;
                            margin-bottom:.3rem;display:flex;align-items:center;gap:.5rem;flex-wrap:wrap">
                    <span class="badge badge-default" style="min-width:70px;text-align:center">
                        ${escapeHtml(cp.component_type||cp.cp_id||'')}</span>
                    <span style="flex:1;font-size:.88rem">
                        ${cp.name?escapeHtml(cp.name):`параметр #${cp.parameter_id}`}</span>
                    ${(cp.min_val||cp.max_val)?`<span style="font-size:.78rem;color:var(--muted)">[${cp.min_val}–${cp.max_val}]</span>`:''}
                    ${!filterType?`
                    <i class="fas fa-edit action-icon" style="color:#a78bfa" title="Редактировать"
                       data-edit-cp="${cp.component_parameter_id}"></i>
                    <i class="fas fa-trash-alt action-icon" style="color:#f87171" title="Удалить"
                       data-del-cp="${cp.component_parameter_id}"></i>`:''}
                </div>`).join('');

            if(!filterType){
                el.querySelectorAll('[data-edit-cp]').forEach(btn=>
                    btn.addEventListener('click',()=>openCompParamForm(btn.dataset.editCp, ()=>refreshCompParams())));
                el.querySelectorAll('[data-del-cp]').forEach(btn=>
                    btn.addEventListener('click',async()=>{
                        if(!confirm('Удалить привязку?'))return;
                        try{await api.componentParameter.delete(btn.dataset.delCp);showToast('Удалено');await refreshCompParams();}
                        catch(e){showToast(e.message,false);}
                    }));
            }
        }catch(e){el.innerHTML=`<p style="color:#f87171">${escapeHtml(e.message)}</p>`;}
    };

    modal.querySelector('#_cpNew').onclick=()=>openCompParamForm(null, ()=>refreshCompParams());
    modal.querySelector('#_cpRef').onclick=()=>refreshCompParams();
    modal.querySelector('#_cpFilterBtn').onclick=()=>{
        const t=modal.querySelector('#_cpTypeFilter').value.trim();
        refreshCompParams(t||null);
    };
    modal.querySelector('#_cpCopy').onclick=()=>openCopyTypeForm();
    await refreshCompParams();

    // ── Значения параметров ──────────────────────────────────────
    const refreshParamVals=async(emobileId=null)=>{
        const el=modal.querySelector('#_epvList');
        el.innerHTML='<i class="fas fa-spinner fa-spin"></i>';
        try{
            const list=emobileId
                ? (await api.emobileParam.byEmobile(emobileId)||[])
                : (await api.emobileParam.list()||[]);
            if(!list.length){el.innerHTML='<p style="color:var(--muted)">Нет значений</p>';return;}
            el.innerHTML=list.map(v=>`
                <div style="border:1px solid var(--border);border-radius:.4rem;padding:.5rem .8rem;
                            margin-bottom:.3rem;display:flex;align-items:center;gap:.5rem;flex-wrap:wrap">
                    <span style="font-size:.78rem;color:var(--muted);min-width:50px">
                        авт.#${shortId(v.emobile_id)}</span>
                    <span style="flex:1;font-size:.88rem">
                        cp#${v.component_parameter_id}: 
                        ${v.val_real?v.val_real+' ':''} 
                        ${v.val_int?v.val_int+' ':''} 
                        ${escapeHtml(v.val_str||v.enum_val_id||'')}
                    </span>
                    <i class="fas fa-edit action-icon" style="color:#a78bfa" title="Редактировать"
                       data-edit-epv="${v.value_id}"></i>
                    <i class="fas fa-trash-alt action-icon" style="color:#f87171" title="Удалить"
                       data-del-epv="${v.value_id}"></i>
                </div>`).join('');

            el.querySelectorAll('[data-edit-epv]').forEach(btn=>
                btn.addEventListener('click',()=>openEmobileParamValueForm(btn.dataset.editEpv, ()=>refreshParamVals())));
            el.querySelectorAll('[data-del-epv]').forEach(btn=>
                btn.addEventListener('click',async()=>{
                    if(!confirm('Удалить значение?'))return;
                    try{await api.emobileParam.delete(btn.dataset.delEpv);showToast('Удалено');await refreshParamVals();}
                    catch(e){showToast(e.message,false);}
                }));
        }catch(e){el.innerHTML=`<p style="color:#f87171">${escapeHtml(e.message)}</p>`;}
    };

    modal.querySelector('#_epNew').onclick=()=>openEmobileParamValueForm(null, ()=>refreshParamVals());
    modal.querySelector('#_epRef').onclick=()=>refreshParamVals();
    modal.querySelector('#_epFilterBtn').onclick=()=>{
        const id=modal.querySelector('#_epEmobileFilter').value.trim();
        refreshParamVals(id||null);
    };
    await refreshParamVals();
}

// ── Форма параметра ───────────────────────────────────────────────

function openParameterForm(id, onSave){
    const addToEmobileHtml = id ? '' : `
        <div class="form-group" style="margin-top:.4rem">
            <label style="display:flex;align-items:center;gap:.5rem;cursor:pointer;user-select:none">
                <input type="checkbox" id="_p_emobile" checked
                       style="width:16px;height:16px;accent-color:var(--accent)">
                <span>Добавить к параметрам автомобилей (emobile)</span>
            </label>
        </div>`;
    const m=createModal(`
        <div class="modal-title">
            <i class="fas fa-${id?'edit':'plus'}" style="color:var(--accent)"></i>
            ${id?'Редактировать':'Новый'} параметр
        </div>
        <div class="form-group">
            <label>Обозначение (код) *</label>
            <input id="_pcod" class="input-dark" placeholder="range_km">
        </div>
        <div class="form-group">
            <label>Полное название *</label>
            <input id="_pnam" class="input-dark" placeholder="Запас хода">
        </div>
        <div class="form-group">
            <label>Тип *</label>
            <select id="_ptyp" class="input-dark">
                <option value="real">real — вещественный</option>
                <option value="int">int — целый</option>
                <option value="str">str — строковый</option>
                <option value="enum">enum — перечисление</option>
            </select>
        </div>
        <div class="form-group">
            <label>Единица измерения</label>
            <input id="_punit" class="input-dark" placeholder="км, кг, кВт…">
        </div>
        <div class="form-group">
            <label>ID перечисления (для типа enum)</label>
            <input id="_pecid" class="input-dark" type="number" placeholder="ID enum_class">
        </div>
        ${addToEmobileHtml}
        <div class="error-banner hidden"><i class="fas fa-exclamation-circle"></i><span></span></div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_pcc">Отмена</button>
            <button class="btn btn-primary" id="_pss">
                <i class="fas fa-save"></i> Сохранить</button>
        </div>`);

    if(id){
        api.parameter.getById(id).then(p=>{
            m.querySelector('#_pcod').value=p.designation||'';
            m.querySelector('#_pnam').value=p.name||'';
            m.querySelector('#_ptyp').value=p.param_type||'real';
            m.querySelector('#_punit').value=p.measuring_unit||'';
            m.querySelector('#_pecid').value=p.enum_class_id||'';
        }).catch(()=>{});
    }

    m.querySelector('#_pcc').onclick=()=>m.remove();
    m.querySelector('#_pss').onclick=async()=>{
        const des=m.querySelector('#_pcod').value.trim();
        const name=m.querySelector('#_pnam').value.trim();
        if(!des||!name){setFormError(m,'Заполните обозначение и название');return;}
        const btn=m.querySelector('#_pss');
        btn.disabled=true;btn.innerHTML='<i class="fas fa-spinner fa-spin"></i>';
        setFormError(m,'');
        try{
            const data={
                designation:des, name,
                param_type:m.querySelector('#_ptyp').value,
                measuring_unit:m.querySelector('#_punit').value,
                enum_class_id:m.querySelector('#_pecid').value,
            };
            if(id){
                await api.parameter.update(id,data);
                showToast('Параметр обновлён');
            } else {
                const resp = await api.parameter.create(data);
                // Если чекбокс отмечен — привязываем к типу emobile
                const linkCb = m.querySelector('#_p_emobile');
                if(linkCb && linkCb.checked && resp && resp.id){
                    await api.componentParameter.create({
                        component_type:'emobile',
                        parameter_id:String(resp.id),
                        order_num:0,
                    }).catch(()=>{});
                }
                showToast('Параметр создан');
            }
            m.remove();if(onSave)onSave();
        }catch(e){
            setFormError(m,e.message);
            btn.disabled=false;btn.innerHTML='<i class="fas fa-save"></i> Сохранить';
        }
    };
}

// ── Форма привязки параметра к типу компонента ─────────────────────

function openCompParamForm(id, onSave){
    const m=createModal(`
        <div class="modal-title">
            <i class="fas fa-link" style="color:var(--accent)"></i>
            ${id?'Редактировать':'Привязать'} параметр к типу
        </div>
        <div class="form-group">
            <label>Тип компонента *</label>
            <select id="_cpt" class="input-dark">
                <option value="emobile">emobile — Автомобиль</option>
                <option value="battery">battery — Батарея</option>
                <option value="engine">engine — Двигатель</option>
                <option value="chassis">chassis — Шасси</option>
                <option value="charger">charger — Зарядное устройство</option>
                <option value="electronics">electronics — Электроника</option>
                <option value="body">body — Кузов</option>
                <option value="frame">frame — Рама</option>
                <option value="suspension">suspension — Подвеска</option>
            </select>
        </div>
        <div class="form-group">
            <label>ID параметра *</label>
            <input id="_cpp" class="input-dark" type="number" placeholder="ID parameter">
        </div>
        <div class="form-group">
            <label>Порядок отображения</label>
            <input id="_cpo" class="input-dark" type="number" value="0">
        </div>
        <div style="display:flex;gap:.5rem">
            <div class="form-group" style="flex:1">
                <label>Минимум</label>
                <input id="_cpmin" class="input-dark" type="number" placeholder="0">
            </div>
            <div class="form-group" style="flex:1">
                <label>Максимум</label>
                <input id="_cpmax" class="input-dark" type="number" placeholder="1000">
            </div>
        </div>
        <div class="error-banner hidden"><i class="fas fa-exclamation-circle"></i><span></span></div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_cpc">Отмена</button>
            <button class="btn btn-primary" id="_cps">
                <i class="fas fa-save"></i> Сохранить</button>
        </div>`);

    if(id){
        api.componentParameter.getById(id).then(cp=>{
            m.querySelector('#_cpt').value=cp.component_type||'emobile';
            m.querySelector('#_cpp').value=cp.parameter_id||'';
            m.querySelector('#_cpo').value=cp.order_num||0;
            m.querySelector('#_cpmin').value=cp.min_val||'';
            m.querySelector('#_cpmax').value=cp.max_val||'';
        }).catch(()=>{});
    }

    m.querySelector('#_cpc').onclick=()=>m.remove();
    m.querySelector('#_cps').onclick=async()=>{
        const paramId=m.querySelector('#_cpp').value.trim();
        if(!paramId){setFormError(m,'Укажите ID параметра');return;}
        const btn=m.querySelector('#_cps');
        btn.disabled=true;btn.innerHTML='<i class="fas fa-spinner fa-spin"></i>';
        setFormError(m,'');
        try{
            const data={
                component_type: m.querySelector('#_cpt').value,
                parameter_id:   paramId,
                order_num:      m.querySelector('#_cpo').value||'0',
                min_val:        m.querySelector('#_cpmin').value||'0',
                max_val:        m.querySelector('#_cpmax').value||'0',
            };
            if(id) await api.componentParameter.update(id,{order_num:data.order_num,min_val:data.min_val,max_val:data.max_val});
            else   await api.componentParameter.create(data);
            showToast('Сохранено');
            m.remove();if(onSave)onSave();
        }catch(e){
            setFormError(m,e.message);
            btn.disabled=false;btn.innerHTML='<i class="fas fa-save"></i> Сохранить';
        }
    };
}

// ── Форма копирования параметров (SQL-процедура) ───────────────────

function openCopyTypeForm(){
    const m=createModal(`
        <div class="modal-title">
            <i class="fas fa-copy" style="color:var(--accent)"></i>
            Скопировать параметры (SQL-процедура)
        </div>
        <p style="font-size:.85rem;color:var(--muted)">
            Копирует все параметры от одного типа компонента к другому.<br>
            Использует хранимую процедуру <code>copy_component_parameters</code>.
        </p>
        <div class="form-group">
            <label>Источник (from_type) *</label>
            <input id="_cft" class="input-dark" placeholder="emobile">
        </div>
        <div class="form-group">
            <label>Получатель (to_type) *</label>
            <input id="_ctt" class="input-dark" placeholder="emobile_sport">
        </div>
        <div class="error-banner hidden"><i class="fas fa-exclamation-circle"></i><span></span></div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_cfcc">Отмена</button>
            <button class="btn btn-primary" id="_cfss">
                <i class="fas fa-copy"></i> Скопировать</button>
        </div>`);

    m.querySelector('#_cfcc').onclick=()=>m.remove();
    m.querySelector('#_cfss').onclick=async()=>{
        const from=m.querySelector('#_cft').value.trim();
        const to=m.querySelector('#_ctt').value.trim();
        if(!from||!to){setFormError(m,'Заполните оба поля');return;}
        const btn=m.querySelector('#_cfss');
        btn.disabled=true;btn.innerHTML='<i class="fas fa-spinner fa-spin"></i>';
        setFormError(m,'');
        try{
            await api.componentParameter.copyFromType({from_type:from, to_type:to});
            showToast(`Параметры скопированы: ${from} → ${to}`);
            m.remove();
        }catch(e){
            setFormError(m,e.message);
            btn.disabled=false;btn.innerHTML='<i class="fas fa-copy"></i> Скопировать';
        }
    };
}

// ── Форма значения параметра для автомобиля ────────────────────────

function openEmobileParamValueForm(id, onSave){
    const m=createModal(`
        <div class="modal-title">
            <i class="fas fa-${id?'edit':'plus'}" style="color:var(--accent)"></i>
            ${id?'Редактировать':'Новое'} значение параметра
        </div>
        <div style="display:flex;gap:.5rem">
            <div class="form-group" style="flex:1">
                <label>ID автомобиля *</label>
                <input id="_evem" class="input-dark" placeholder="emobile_id">
            </div>
            <div class="form-group" style="flex:1">
                <label>ID привязки параметра *</label>
                <input id="_evcp" class="input-dark" placeholder="component_parameter_id">
            </div>
        </div>
        <div style="display:flex;gap:.5rem;flex-wrap:wrap">
            <div class="form-group" style="flex:1">
                <label>val_real (вещественный)</label>
                <input id="_evr" class="input-dark" type="number" step="any" placeholder="320.5">
            </div>
            <div class="form-group" style="flex:1">
                <label>val_int (целый)</label>
                <input id="_evi" class="input-dark" type="number" placeholder="1850">
            </div>
        </div>
        <div style="display:flex;gap:.5rem;flex-wrap:wrap">
            <div class="form-group" style="flex:1">
                <label>val_str (строковый)</label>
                <input id="_evs" class="input-dark" placeholder="Германия">
            </div>
            <div class="form-group" style="flex:1">
                <label>enum_val_id (ID позиции)</label>
                <input id="_eve" class="input-dark" type="number" placeholder="ID enum_position">
            </div>
        </div>
        <div class="error-banner hidden"><i class="fas fa-exclamation-circle"></i><span></span></div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_evc">Отмена</button>
            <button class="btn btn-primary" id="_evs2">
                <i class="fas fa-save"></i> Сохранить</button>
        </div>`);

    if(id){
        api.emobileParam.getById(id).then(v=>{
            m.querySelector('#_evem').value=v.emobile_id||'';
            m.querySelector('#_evcp').value=v.component_parameter_id||'';
            m.querySelector('#_evr').value=v.val_real||'';
            m.querySelector('#_evi').value=v.val_int||'';
            m.querySelector('#_evs').value=v.val_str||'';
            m.querySelector('#_eve').value=v.enum_val_id||'';
        }).catch(()=>{});
    }

    m.querySelector('#_evc').onclick=()=>m.remove();
    m.querySelector('#_evs2').onclick=async()=>{
        const emId=m.querySelector('#_evem').value.trim();
        const cpId=m.querySelector('#_evcp').value.trim();
        if(!emId||!cpId){setFormError(m,'Укажите ID автомобиля и привязки');return;}
        const btn=m.querySelector('#_evs2');
        btn.disabled=true;btn.innerHTML='<i class="fas fa-spinner fa-spin"></i>';
        setFormError(m,'');
        try{
            const data={
                emobile_id:emId, component_parameter_id:cpId,
                val_real:m.querySelector('#_evr').value||'0',
                val_int:m.querySelector('#_evi').value||'0',
                val_str:m.querySelector('#_evs').value,
                enum_val_id:m.querySelector('#_eve').value||'',
            };
            if(id) await api.emobileParam.update(id,{val_real:data.val_real,val_int:data.val_int,val_str:data.val_str,enum_val_id:data.enum_val_id});
            else   await api.emobileParam.create(data);
            showToast('Значение сохранено');
            m.remove();if(onSave)onSave();
        }catch(e){
            setFormError(m,e.message);
            btn.disabled=false;btn.innerHTML='<i class="fas fa-save"></i> Сохранить';
        }
    };
}

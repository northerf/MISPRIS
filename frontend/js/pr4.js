/* pr4.js — ПР4: Хозяйственные операции (SHD, HO Class, HO Roles, HO Instances) */

// ══════════════════════════════════════════════════════════════════
//  УТИЛИТЫ
// ══════════════════════════════════════════════════════════════════

// showToast already defined in pr2pr3.js; provide fallback if loaded standalone
if(typeof showToast === 'undefined'){
    window.showToast = function(msg, ok=true){
        const t=document.createElement('div');
        t.className=`toast ${ok?'toast-ok':'toast-error'}`;
        t.setAttribute('role',ok?'status':'alert');
        t.textContent=msg;
        document.body.appendChild(t);
        setTimeout(()=>t.remove(),4200);
    };
}

function _shdTypeName(t){
    return {organization:'Юридическое лицо',individual:'Физическое лицо',department:'Подразделение'}[t]||t||'—';
}

function _statusBadge(s){
    const map={
        draft:'#475569',new:'#1d4ed8',active:'#15803d',confirmed:'#0e7490',
        cancelled:'#991b1b',closed:'#374151'
    };
    const labels={draft:'Черновик',confirmed:'Подтверждена',cancelled:'Отменена',new:'Новая',active:'Активна',closed:'Закрыта'};
    const col=map[s]||'#334155';
    return `<span style="background:${col};color:#e2e8f0;border-radius:.3rem;padding:.15rem .45rem;font-size:.75rem;font-weight:600">${escapeHtml(labels[s]||s||'Черновик')}</span>`;
}

function _isHoAdmin(){
    return currentRole==='admin';
}

function _formatHoDate(value){
    if(!value)return '—';
    return String(value).slice(0,10);
}

function _dateInputValue(value){
    if(!value)return '';
    return String(value).slice(0,10);
}

function _formatHoMoney(value){
    if(value===null||value===undefined||value==='')return '—';
    return Number(value).toLocaleString('ru-RU');
}

function _loadingHtml(){
    return '<div style="padding:1.2rem;text-align:center;color:var(--muted)"><i class="fas fa-spinner fa-spin"></i> Загрузка…</div>';
}

function _emptyHtml(msg='Нет данных'){
    return `<div style="padding:1rem;text-align:center;color:var(--muted);font-size:.85rem"><i class="fas fa-inbox"></i> ${escapeHtml(msg)}</div>`;
}

// Build a simple <select> from array, with optional empty option
function _buildSelect(id, items, keyField, labelField, selectedVal='', emptyLabel='— выберите —'){
    let opts=`<option value="">${escapeHtml(emptyLabel)}</option>`;
    for(const it of items){
        const v=it[keyField]||'';
        const l=it[labelField]||String(v);
        opts+=`<option value="${escapeHtml(String(v))}" ${String(v)===String(selectedVal)?'selected':''}>${escapeHtml(l)}</option>`;
    }
    return `<select id="${id}" class="input-dark">${opts}</select>`;
}

// ══════════════════════════════════════════════════════════════════
//  РАЗДЕЛ A — КОНФИГУРАЦИЯ (admin)
// ══════════════════════════════════════════════════════════════════

async function openHoConfigManager(){
    const modal=createModal(`
        <div class="modal-title">
            <i class="fas fa-cog" style="color:var(--accent)"></i>
            ПР4: Настройка Хозяйственных Операций
        </div>
        <div style="display:flex;gap:.3rem;margin-bottom:.8rem;flex-wrap:wrap;border-bottom:1px solid var(--border);padding-bottom:.6rem">
            <button class="btn btn-secondary _cfg-tab active" data-tab="hoclass">
                <i class="fas fa-sitemap"></i> Типы ХО</button>
            <button class="btn btn-secondary _cfg-tab" data-tab="roles">
                <i class="fas fa-users"></i> Роли</button>
            <button class="btn btn-secondary _cfg-tab" data-tab="shd">
                <i class="fas fa-building"></i> СХД</button>
            <button class="btn btn-secondary _cfg-tab" data-tab="docs">
                <i class="fas fa-file-alt"></i> Документы</button>
        </div>
        <div id="_cfg-body" style="min-height:200px">${_loadingHtml()}</div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_cfg-close">Закрыть</button>
        </div>`,true);

    modal.querySelector('#_cfg-close').onclick=()=>modal.remove();

    modal.querySelectorAll('._cfg-tab').forEach(btn=>{
        btn.addEventListener('click',()=>{
            modal.querySelectorAll('._cfg-tab').forEach(b=>b.classList.remove('active'));
            btn.classList.add('active');
            _renderCfgTab(modal, btn.dataset.tab);
        });
    });

    await _renderCfgTab(modal,'hoclass');
}

async function _renderCfgTab(modal, tab){
    const body=modal.querySelector('#_cfg-body');
    body.innerHTML=_loadingHtml();
    try{
        if(tab==='hoclass') await _renderHoClassTab(modal, body);
        else if(tab==='roles') await _renderHoRolesTab(modal, body);
        else if(tab==='shd')  await _renderShdTab(modal, body);
        else if(tab==='docs') await _renderDocClassTab(modal, body);
    }catch(e){
        body.innerHTML=`<p style="color:#f87171">${escapeHtml(e.message)}</p>`;
    }
}

// ── Типы ХО ────────────────────────────────────────────────────────

async function _renderHoClassTab(modal, body){
    const list=await api.hoClass.list()||[];

    // Build tree structure (parent_id → children)
    const byParent={};
    for(const c of list){
        const pid=c.parent_id||null;
        (byParent[pid]||(byParent[pid]=[])).push(c);
    }

    function renderNode(nodes, depth=0){
        if(!nodes||!nodes.length) return '';
        return nodes.map(c=>`
            <div style="margin-left:${depth*1.2}rem;border:1px solid var(--border);
                        border-radius:.4rem;padding:.5rem .7rem;margin-bottom:.35rem;
                        background:rgba(255,255,255,.02)">
                <div style="display:flex;align-items:center;gap:.4rem;flex-wrap:wrap">
                    <span style="font-weight:600;flex:1;min-width:120px">
                        ${escapeHtml(c.name||'—')}
                        ${c.designation?`<span style="color:var(--muted);font-size:.78rem;margin-left:.3rem">[${escapeHtml(c.designation)}]</span>`:''}
                    </span>
                    ${c.is_terminal?'<span style="font-size:.72rem;background:#1e3a5f;color:#93c5fd;border-radius:.25rem;padding:.1rem .35rem">конечный</span>':''}
                    <button class="btn btn-secondary" style="font-size:.72rem;padding:.2rem .5rem" data-hc-roles="${c.ho_class_id}" data-hc-name="${escapeHtml(c.name)}">
                        <i class="fas fa-users"></i> Роли</button>
                    <button class="btn btn-secondary" style="font-size:.72rem;padding:.2rem .5rem" data-hc-params="${c.ho_class_id}" data-hc-name="${escapeHtml(c.name)}">
                        <i class="fas fa-sliders-h"></i> Параметры</button>
                    <button class="btn btn-secondary" style="font-size:.72rem;padding:.2rem .5rem" data-hc-docs="${c.ho_class_id}" data-hc-name="${escapeHtml(c.name)}">
                        <i class="fas fa-file-alt"></i> Документы</button>
                    <i class="fas fa-plus-circle action-icon" style="color:#4ade80" title="Добавить дочерний" data-hc-child="${c.ho_class_id}"></i>
                    <i class="fas fa-edit action-icon" style="color:#60a5fa" title="Редактировать" data-hc-edit="${c.ho_class_id}"></i>
                    <i class="fas fa-trash-alt action-icon" style="color:#f87171" title="Удалить" data-hc-del="${c.ho_class_id}" data-hc-dname="${escapeHtml(c.name)}"></i>
                </div>
            </div>
            ${renderNode(byParent[c.ho_class_id]||[], depth+1)}`
        ).join('');
    }

    const rootNodes=byParent[null]||byParent['null']||
        list.filter(c=>!c.parent_id);

    body.innerHTML=`
        <div style="display:flex;gap:.5rem;margin-bottom:.8rem">
            <button class="btn btn-primary" id="_hc-new">
                <i class="fas fa-plus"></i> Новый тип ХО</button>
            <button class="btn btn-secondary" id="_hc-refresh">
                <i class="fas fa-sync"></i> Обновить</button>
        </div>
        <div id="_hc-tree" style="max-height:420px;overflow-y:auto">
            ${list.length?renderNode(rootNodes):_emptyHtml('Нет типов ХО')}
        </div>`;

    const tree=body.querySelector('#_hc-tree');

    body.querySelector('#_hc-new').onclick=()=>
        openHoClassForm(null,null,()=>_renderHoClassTab(modal,body));
    body.querySelector('#_hc-refresh').onclick=()=>_renderHoClassTab(modal,body);

    tree.querySelectorAll('[data-hc-child]').forEach(el=>
        el.addEventListener('click',()=>
            openHoClassForm(null,el.dataset.hcChild,()=>_renderHoClassTab(modal,body))));
    tree.querySelectorAll('[data-hc-edit]').forEach(el=>
        el.addEventListener('click',()=>
            openHoClassForm(el.dataset.hcEdit,null,()=>_renderHoClassTab(modal,body))));
    tree.querySelectorAll('[data-hc-del]').forEach(el=>
        el.addEventListener('click',async()=>{
            if(!confirm(`Удалить тип ХО «${el.dataset.hcDname}»?`))return;
            try{
                await api.hoClass.delete(el.dataset.hcDel);
                showToast('Тип ХО удалён');
                await _renderHoClassTab(modal,body);
            }catch(e){showToast(e.message,false);}
        }));
    tree.querySelectorAll('[data-hc-roles]').forEach(el=>
        el.addEventListener('click',()=>
            openHoClassRolesManager(el.dataset.hcRoles,el.dataset.hcName)));
    tree.querySelectorAll('[data-hc-params]').forEach(el=>
        el.addEventListener('click',()=>
            openHoClassParamsManager(el.dataset.hcParams,el.dataset.hcName)));
    tree.querySelectorAll('[data-hc-docs]').forEach(el=>
        el.addEventListener('click',()=>
            openHoClassDocsManager(el.dataset.hcDocs,el.dataset.hcName)));
}

function openHoClassForm(id, parentId, onSave){
    const isEdit=!!id;
    const m=createModal(`
        <div class="modal-title">
            <i class="fas fa-${isEdit?'edit':'plus-circle'}" style="color:var(--accent)"></i>
            ${isEdit?'Редактировать':'Новый'} тип ХО
        </div>
        <div class="form-group">
            <label>Название *</label>
            <input id="_hcn" class="input-dark" placeholder="Закупка оборудования">
        </div>
        <div class="form-group">
            <label>Обозначение</label>
            <input id="_hcd" class="input-dark" placeholder="ЗО">
        </div>
        <div class="form-group">
            <label>Родительский тип (ID)</label>
            <input id="_hcp" class="input-dark" placeholder="UUID или пусто" value="${escapeHtml(parentId||'')}">
        </div>
        <div class="form-group" style="display:flex;align-items:center;gap:.5rem">
            <input type="checkbox" id="_hct">
            <label for="_hct" style="margin:0;cursor:pointer">Конечный (терминальный)</label>
        </div>
        <div id="_hcErr"></div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_hcCancel">Отмена</button>
            <button class="btn btn-primary"   id="_hcSave">
                <i class="fas fa-check"></i> ${isEdit?'Сохранить':'Создать'}</button>
        </div>`);

    m.querySelector('#_hcCancel').onclick=()=>m.remove();

    if(isEdit){
        api.hoClass.getById(id).then(c=>{
            if(!c)return;
            m.querySelector('#_hcn').value=c.name||'';
            m.querySelector('#_hcd').value=c.designation||'';
            m.querySelector('#_hcp').value=c.parent_id||'';
            m.querySelector('#_hct').checked=!!c.is_terminal;
        }).catch(()=>{});
    }

    m.querySelector('#_hcSave').onclick=async()=>{
        const name=m.querySelector('#_hcn').value.trim();
        if(!name){setFormError(m,'Название обязательно');return;}
        const data={
            name,
            designation:m.querySelector('#_hcd').value.trim()||null,
            parent_id:m.querySelector('#_hcp').value.trim()||null,
            is_terminal:m.querySelector('#_hct').checked,
        };
        try{
            if(isEdit) await api.hoClass.update(id,data);
            else       await api.hoClass.create(data);
            showToast(isEdit?'Тип ХО обновлён':'Тип ХО создан');
            m.remove();
            onSave&&onSave();
        }catch(e){setFormError(m,e.message);}
    };
}

// ── Роли ───────────────────────────────────────────────────────────

async function _renderHoRolesTab(modal, body){
    const list=await api.hoRole.list()||[];
    body.innerHTML=`
        <div style="display:flex;gap:.5rem;margin-bottom:.8rem">
            <button class="btn btn-primary" id="_hr-new">
                <i class="fas fa-plus"></i> Новая роль</button>
            <button class="btn btn-secondary" id="_hr-refresh">
                <i class="fas fa-sync"></i> Обновить</button>
        </div>
        <div id="_hr-list">
            ${list.length?list.map(r=>`
                <div style="border:1px solid var(--border);border-radius:.4rem;padding:.5rem .8rem;
                            margin-bottom:.35rem;display:flex;align-items:center;gap:.5rem">
                    <span style="flex:1;font-weight:500">${escapeHtml(r.name||'—')}</span>
                    ${r.description?`<span class="badge badge-default">${escapeHtml(r.description)}</span>`:''}
                    <i class="fas fa-edit action-icon" style="color:#60a5fa" data-hr-edit="${r.ho_role_id}"></i>
                    <i class="fas fa-trash-alt action-icon" style="color:#f87171"
                       data-hr-del="${r.ho_role_id}" data-hr-dname="${escapeHtml(r.name)}"></i>
                </div>`).join(''):_emptyHtml('Нет ролей')}
        </div>`;

    body.querySelector('#_hr-new').onclick=()=>
        openHoRoleForm(null,()=>_renderHoRolesTab(modal,body));
    body.querySelector('#_hr-refresh').onclick=()=>_renderHoRolesTab(modal,body);

    const lst=body.querySelector('#_hr-list');
    lst.querySelectorAll('[data-hr-edit]').forEach(el=>
        el.addEventListener('click',()=>openHoRoleForm(el.dataset.hrEdit,()=>_renderHoRolesTab(modal,body))));
    lst.querySelectorAll('[data-hr-del]').forEach(el=>
        el.addEventListener('click',async()=>{
            if(!confirm(`Удалить роль «${el.dataset.hrDname}»?`))return;
            try{
                await api.hoRole.delete(el.dataset.hrDel);
                showToast('Роль удалена');
                await _renderHoRolesTab(modal,body);
            }catch(e){showToast(e.message,false);}
        }));
}

function openHoRoleForm(id, onSave){
    const isEdit=!!id;
    const m=createModal(`
        <div class="modal-title">
            <i class="fas fa-${isEdit?'edit':'plus-circle'}" style="color:var(--accent)"></i>
            ${isEdit?'Редактировать':'Новая'} роль ХО
        </div>
        <div class="form-group">
            <label>Название *</label>
            <input id="_hrn" class="input-dark" placeholder="Грузоотправитель">
        </div>
        <div class="form-group">
            <label>Описание</label>
            <input id="_hrd" class="input-dark" placeholder="Организация, отправляющая груз">
        </div>
        <div id="_hrErr"></div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_hrCancel">Отмена</button>
            <button class="btn btn-primary"   id="_hrSave">
                <i class="fas fa-check"></i> ${isEdit?'Сохранить':'Создать'}</button>
        </div>`);

    m.querySelector('#_hrCancel').onclick=()=>m.remove();

    if(isEdit){
        api.hoRole.getById(id).then(r=>{
            if(!r)return;
            m.querySelector('#_hrn').value=r.name||'';
            m.querySelector('#_hrd').value=r.description||'';
        }).catch(()=>{});
    }

    m.querySelector('#_hrSave').onclick=async()=>{
        const name=m.querySelector('#_hrn').value.trim();
        if(!name){setFormError(m,'Название обязательно');return;}
        const data={name,description:m.querySelector('#_hrd').value.trim()||''};
        try{
            if(isEdit) await api.hoRole.update(id,data);
            else       await api.hoRole.create(data);
            showToast(isEdit?'Роль обновлена':'Роль создана');
            m.remove();
            onSave&&onSave();
        }catch(e){setFormError(m,e.message);}
    };
}

// ── СХД ────────────────────────────────────────────────────────────

async function _renderShdTab(modal, body){
    const list=await api.shd.list()||[];
    body.innerHTML=`
        <div style="display:flex;gap:.5rem;margin-bottom:.8rem">
            <button class="btn btn-primary" id="_shd-new">
                <i class="fas fa-plus"></i> Новый СХД</button>
            <button class="btn btn-secondary" id="_shd-refresh">
                <i class="fas fa-sync"></i> Обновить</button>
        </div>
        <div id="_shd-list">
            ${list.length?list.map(s=>`
                <div style="border:1px solid var(--border);border-radius:.4rem;padding:.5rem .8rem;
                            margin-bottom:.35rem;display:flex;align-items:center;gap:.5rem;flex-wrap:wrap">
                    <span style="flex:1;font-weight:500">${escapeHtml(s.name||'—')}</span>
                    <span class="badge badge-default">${_shdTypeName(s.shd_type)}</span>
                    ${s.inn?`<span style="font-size:.78rem;color:var(--muted)">ИНН: ${escapeHtml(s.inn)}</span>`:''}
                    <i class="fas fa-edit action-icon" style="color:#60a5fa" data-shd-edit="${s.shd_id}"></i>
                    <i class="fas fa-trash-alt action-icon" style="color:#f87171"
                       data-shd-del="${s.shd_id}" data-shd-dname="${escapeHtml(s.name)}"></i>
                </div>`).join(''):_emptyHtml('Нет субъектов')}
        </div>`;

    body.querySelector('#_shd-new').onclick=()=>
        openShdForm(null,()=>_renderShdTab(modal,body));
    body.querySelector('#_shd-refresh').onclick=()=>_renderShdTab(modal,body);

    const lst=body.querySelector('#_shd-list');
    lst.querySelectorAll('[data-shd-edit]').forEach(el=>
        el.addEventListener('click',()=>openShdForm(el.dataset.shdEdit,()=>_renderShdTab(modal,body))));
    lst.querySelectorAll('[data-shd-del]').forEach(el=>
        el.addEventListener('click',async()=>{
            if(!confirm(`Удалить СХД «${el.dataset.shdDname}»?`))return;
            try{
                await api.shd.delete(el.dataset.shdDel);
                showToast('СХД удалён');
                await _renderShdTab(modal,body);
            }catch(e){showToast(e.message,false);}
        }));
}

function openShdForm(id, onSave){
    const isEdit=!!id;
    const m=createModal(`
        <div class="modal-title">
            <i class="fas fa-${isEdit?'edit':'plus-circle'}" style="color:var(--accent)"></i>
            ${isEdit?'Редактировать':'Новый'} субъект хозяйственной деятельности
        </div>
        <div class="form-group">
            <label>Название *</label>
            <input id="_sn" class="input-dark" placeholder="ООО «Пример»">
        </div>
        <div class="form-group">
            <label>Тип</label>
            <select id="_st" class="input-dark">
                <option value="organization">Юридическое лицо</option>
                <option value="individual">Физическое лицо</option>
                <option value="department">Подразделение</option>
            </select>
        </div>
        <div class="form-group">
            <label>ИНН</label>
            <input id="_si" class="input-dark" placeholder="1234567890">
        </div>
        <div id="_sErr"></div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_sCancel">Отмена</button>
            <button class="btn btn-primary"   id="_sSave">
                <i class="fas fa-check"></i> ${isEdit?'Сохранить':'Создать'}</button>
        </div>`);

    m.querySelector('#_sCancel').onclick=()=>m.remove();

    if(isEdit){
        api.shd.getById(id).then(s=>{
            if(!s)return;
            m.querySelector('#_sn').value=s.name||'';
            m.querySelector('#_st').value=s.shd_type||'organization';
            m.querySelector('#_si').value=s.inn||'';
        }).catch(()=>{});
    }

    m.querySelector('#_sSave').onclick=async()=>{
        const name=m.querySelector('#_sn').value.trim();
        if(!name){setFormError(m,'Название обязательно');return;}
        const data={
            name,
            shd_type:m.querySelector('#_st').value,
            inn:m.querySelector('#_si').value.trim()||null,
        };
        try{
            if(isEdit) await api.shd.update(id,data);
            else       await api.shd.create(data);
            showToast(isEdit?'СХД обновлён':'СХД создан');
            m.remove();
            onSave&&onSave();
        }catch(e){setFormError(m,e.message);}
    };
}

// Standalone SHD manager (accessible from catalog toolbar if needed)
async function openShdManager(){
    const modal=createModal(`
        <div class="modal-title">
            <i class="fas fa-building" style="color:var(--accent)"></i>
            Субъекты хозяйственной деятельности
        </div>
        <div style="display:flex;gap:.5rem;margin-bottom:.8rem">
            <button class="btn btn-primary" id="_sm-new">
                <i class="fas fa-plus"></i> Новый СХД</button>
            <button class="btn btn-secondary" id="_sm-refresh">
                <i class="fas fa-sync"></i> Обновить</button>
        </div>
        <div id="_sm-body">${_loadingHtml()}</div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_sm-close">Закрыть</button>
        </div>`,true);

    const refreshSm=async()=>{
        const body=modal.querySelector('#_sm-body');
        body.innerHTML=_loadingHtml();
        try{
            const list=await api.shd.list()||[];
            if(!list.length){body.innerHTML=_emptyHtml('Нет субъектов');return;}
            body.innerHTML=list.map(s=>`
                <div style="border:1px solid var(--border);border-radius:.4rem;padding:.5rem .8rem;
                            margin-bottom:.35rem;display:flex;align-items:center;gap:.5rem;flex-wrap:wrap">
                    <span style="flex:1;font-weight:500">${escapeHtml(s.name||'—')}</span>
                    <span class="badge badge-default">${_shdTypeName(s.shd_type)}</span>
                    ${s.inn?`<span style="font-size:.78rem;color:var(--muted)">ИНН: ${escapeHtml(s.inn)}</span>`:''}
                    <i class="fas fa-edit action-icon" style="color:#60a5fa" data-shd-edit="${s.shd_id}"></i>
                    <i class="fas fa-trash-alt action-icon" style="color:#f87171"
                       data-shd-del="${s.shd_id}" data-shd-dname="${escapeHtml(s.name)}"></i>
                </div>`).join('');
            body.querySelectorAll('[data-shd-edit]').forEach(el=>
                el.addEventListener('click',()=>openShdForm(el.dataset.shdEdit,refreshSm)));
            body.querySelectorAll('[data-shd-del]').forEach(el=>
                el.addEventListener('click',async()=>{
                    if(!confirm(`Удалить СХД «${el.dataset.shdDname}»?`))return;
                    try{await api.shd.delete(el.dataset.shdDel);showToast('СХД удалён');await refreshSm();}
                    catch(e){showToast(e.message,false);}
                }));
        }catch(e){body.innerHTML=`<p style="color:#f87171">${escapeHtml(e.message)}</p>`;}
    };

    modal.querySelector('#_sm-close').onclick=()=>modal.remove();
    modal.querySelector('#_sm-new').onclick=()=>openShdForm(null,refreshSm);
    modal.querySelector('#_sm-refresh').onclick=refreshSm;
    await refreshSm();
}

// ── Классы документов ──────────────────────────────────────────────

async function _renderDocClassTab(modal, body){
    const list=await api.documentClass.list()||[];
    body.innerHTML=`
        <div style="display:flex;gap:.5rem;margin-bottom:.8rem">
            <button class="btn btn-primary" id="_dc-new">
                <i class="fas fa-plus"></i> Новый класс документа</button>
            <button class="btn btn-secondary" id="_dc-refresh">
                <i class="fas fa-sync"></i> Обновить</button>
        </div>
        <div id="_dc-list">
            ${list.length?list.map(d=>`
                <div style="border:1px solid var(--border);border-radius:.4rem;padding:.5rem .8rem;
                            margin-bottom:.35rem;display:flex;align-items:center;gap:.5rem">
                    <span style="flex:1;font-weight:500">${escapeHtml(d.name||'—')}</span>
                    ${d.code?`<span class="badge badge-default">${escapeHtml(d.code)}</span>`:''}
                    <i class="fas fa-edit action-icon" style="color:#60a5fa" data-dc-edit="${d.doc_class_id}"></i>
                    <i class="fas fa-trash-alt action-icon" style="color:#f87171"
                       data-dc-del="${d.doc_class_id}" data-dc-dname="${escapeHtml(d.name)}"></i>
                </div>`).join(''):_emptyHtml('Нет классов документов')}
        </div>`;

    body.querySelector('#_dc-new').onclick=()=>
        openDocClassForm(null,()=>_renderDocClassTab(modal,body));
    body.querySelector('#_dc-refresh').onclick=()=>_renderDocClassTab(modal,body);

    const lst=body.querySelector('#_dc-list');
    lst.querySelectorAll('[data-dc-edit]').forEach(el=>
        el.addEventListener('click',()=>openDocClassForm(el.dataset.dcEdit,()=>_renderDocClassTab(modal,body))));
    lst.querySelectorAll('[data-dc-del]').forEach(el=>
        el.addEventListener('click',async()=>{
            if(!confirm(`Удалить класс документа «${el.dataset.dcDname}»?`))return;
            try{
                await api.documentClass.delete(el.dataset.dcDel);
                showToast('Класс документа удалён');
                await _renderDocClassTab(modal,body);
            }catch(e){showToast(e.message,false);}
        }));
}

function openDocClassForm(id, onSave){
    const isEdit=!!id;
    const m=createModal(`
        <div class="modal-title">
            <i class="fas fa-file-alt" style="color:var(--accent)"></i>
            ${isEdit?'Редактировать':'Новый'} класс документа
        </div>
        <div class="form-group">
            <label>Название *</label>
            <input id="_dcn" class="input-dark" placeholder="Счёт-фактура">
        </div>
        <div class="form-group">
            <label>Код документа</label>
            <input id="_dcd" class="input-dark" placeholder="0309003">
        </div>
        <div id="_dcErr"></div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_dcCancel">Отмена</button>
            <button class="btn btn-primary"   id="_dcSave">
                <i class="fas fa-check"></i> ${isEdit?'Сохранить':'Создать'}</button>
        </div>`);

    m.querySelector('#_dcCancel').onclick=()=>m.remove();
    if(isEdit){
        api.documentClass.list().then(list=>{
            const d=(list||[]).find(x=>String(x.doc_class_id)===String(id));
            if(!d)return;
            m.querySelector('#_dcn').value=d.name||'';
            m.querySelector('#_dcd').value=d.code||'';
        }).catch(()=>{});
    }
    m.querySelector('#_dcSave').onclick=async()=>{
        const name=m.querySelector('#_dcn').value.trim();
        if(!name){setFormError(m,'Название обязательно');return;}
        const data={name,code:m.querySelector('#_dcd').value.trim()||'',description:''};
        try{
            if(isEdit) await api.documentClass.update(id,data);
            else       await api.documentClass.create(data);
            showToast(isEdit?'Класс документа обновлён':'Класс документа создан');
            m.remove();
            onSave&&onSave();
        }catch(e){setFormError(m,e.message);}
    };
}

// ── Роли типа ХО ──────────────────────────────────────────────────

async function openHoClassRolesManager(hoClassId, className){
    const modal=createModal(`
        <div class="modal-title">
            <i class="fas fa-users" style="color:var(--accent)"></i>
            Роли типа ХО: ${escapeHtml(className)}
        </div>
        <div id="_hcr-list" style="margin-bottom:.8rem;min-height:60px">${_loadingHtml()}</div>
        <div style="border-top:1px solid var(--border);padding-top:.8rem;margin-top:.4rem">
            <div style="font-weight:600;margin-bottom:.5rem;font-size:.85rem">Добавить роль:</div>
            <div style="display:flex;gap:.5rem;flex-wrap:wrap;align-items:flex-end">
                <div class="form-group" style="flex:1;min-width:160px;margin:0">
                    <select id="_hcr-sel" class="input-dark">
                        <option value="">Загрузка ролей…</option>
                    </select>
                </div>
                <div style="display:flex;align-items:center;gap:.3rem">
                    <input type="checkbox" id="_hcr-req">
                    <label for="_hcr-req" style="font-size:.82rem;cursor:pointer;margin:0">Обязательная</label>
                </div>
                <button class="btn btn-primary" id="_hcr-add">
                    <i class="fas fa-plus"></i> Добавить</button>
            </div>
        </div>
        <div id="_hcrErr"></div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_hcr-close">Закрыть</button>
        </div>`);

    modal.querySelector('#_hcr-close').onclick=()=>modal.remove();

    // Load available roles into select
    api.hoRole.list().then(roles=>{
        const sel=modal.querySelector('#_hcr-sel');
        roles=roles||[];
        sel.innerHTML=`<option value="">— выберите роль —</option>`+
            roles.map(r=>`<option value="${escapeHtml(String(r.ho_role_id))}">${escapeHtml(r.name)}</option>`).join('');
    }).catch(()=>{});

    const refreshList=async()=>{
        const el=modal.querySelector('#_hcr-list');
        el.innerHTML=_loadingHtml();
        try{
            const list=await api.hoClassRole.list(hoClassId)||[];
            if(!list.length){el.innerHTML=_emptyHtml('Роли не назначены');return;}
            el.innerHTML=list.map(r=>`
                <div style="border:1px solid var(--border);border-radius:.4rem;padding:.4rem .7rem;
                            margin-bottom:.3rem;display:flex;align-items:center;gap:.5rem">
                    <span style="flex:1">${escapeHtml(r.role_name||r.ho_role_id)}</span>
                    ${r.is_required?'<span style="font-size:.72rem;background:#1e3a5f;color:#93c5fd;border-radius:.25rem;padding:.1rem .35rem">обязательная</span>':''}
                    <i class="fas fa-trash-alt action-icon" style="color:#f87171"
                       data-hcr-del="${r.id}"></i>
                </div>`).join('');
            el.querySelectorAll('[data-hcr-del]').forEach(btn=>
                btn.addEventListener('click',async()=>{
                    try{
                        await api.hoClassRole.delete(btn.dataset.hcrDel);
                        showToast('Роль удалена из типа ХО');
                        await refreshList();
                    }catch(e){showToast(e.message,false);}
                }));
        }catch(e){el.innerHTML=`<p style="color:#f87171">${escapeHtml(e.message)}</p>`;}
    };

    modal.querySelector('#_hcr-add').onclick=async()=>{
        const roleId=modal.querySelector('#_hcr-sel').value;
        if(!roleId){showToast('Выберите роль',false);return;}
        const data={
            ho_class_id:hoClassId,
            ho_role_id:roleId,
            is_required:modal.querySelector('#_hcr-req').checked,
        };
        try{
            await api.hoClassRole.create(data);
            showToast('Роль добавлена');
            modal.querySelector('#_hcr-sel').value='';
            modal.querySelector('#_hcr-req').checked=false;
            await refreshList();
        }catch(e){showToast(e.message,false);}
    };

    await refreshList();
}

// ── Параметры типа ХО ──────────────────────────────────────────────

async function openHoClassParamsManager(hoClassId, className){
    const modal=createModal(`
        <div class="modal-title">
            <i class="fas fa-sliders-h" style="color:var(--accent)"></i>
            Параметры типа ХО: ${escapeHtml(className)}
        </div>
        <div style="display:flex;gap:.5rem;margin-bottom:.6rem">
            <button class="btn btn-secondary" id="_hcp-refresh">
                <i class="fas fa-sync"></i> Обновить</button>
            <button class="btn btn-secondary" id="_hcp-copy" title="Скопировать параметры из другого типа ХО">
                <i class="fas fa-copy"></i> Скопировать из класса</button>
        </div>
        <div id="_hcp-list" style="margin-bottom:.8rem;min-height:60px">${_loadingHtml()}</div>
        <div style="border-top:1px solid var(--border);padding-top:.8rem">
            <div style="font-weight:600;margin-bottom:.5rem;font-size:.85rem">Добавить параметр:</div>
            <div style="display:flex;gap:.5rem;flex-wrap:wrap;align-items:flex-end">
                <div class="form-group" style="flex:2;min-width:160px;margin:0">
                    <label style="font-size:.78rem">ID параметра</label>
                    <input id="_hcpp-id" class="input-dark" placeholder="UUID параметра">
                </div>
                <div class="form-group" style="flex:1;min-width:80px;margin:0">
                    <label style="font-size:.78rem">Порядок</label>
                    <input id="_hcpp-ord" class="input-dark" type="number" placeholder="0" value="0">
                </div>
                <div class="form-group" style="flex:1;min-width:70px;margin:0">
                    <label style="font-size:.78rem">Мин</label>
                    <input id="_hcpp-min" class="input-dark" type="number" placeholder="—">
                </div>
                <div class="form-group" style="flex:1;min-width:70px;margin:0">
                    <label style="font-size:.78rem">Макс</label>
                    <input id="_hcpp-max" class="input-dark" type="number" placeholder="—">
                </div>
                <button class="btn btn-primary" id="_hcpp-add">
                    <i class="fas fa-plus"></i></button>
            </div>
        </div>
        <div id="_hcpErr"></div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_hcp-close">Закрыть</button>
        </div>`,true);

    modal.querySelector('#_hcp-close').onclick=()=>modal.remove();

    const refreshList=async()=>{
        const el=modal.querySelector('#_hcp-list');
        el.innerHTML=_loadingHtml();
        try{
            const list=await api.hoClass.parameters(hoClassId)||[];
            if(!list.length){el.innerHTML=_emptyHtml('Параметры не назначены');return;}
            el.innerHTML=`<table style="width:100%;border-collapse:collapse;font-size:.83rem">
                <thead><tr style="color:var(--muted);border-bottom:1px solid var(--border)">
                    <th style="text-align:left;padding:.25rem .4rem">Параметр</th>
                    <th style="text-align:center;padding:.25rem .4rem">Порядок</th>
                    <th style="text-align:center;padding:.25rem .4rem">Мин</th>
                    <th style="text-align:center;padding:.25rem .4rem">Макс</th>
                    <th></th>
                </tr></thead>
                <tbody>${list.map(p=>`<tr style="border-bottom:1px solid #1e293b">
                    <td style="padding:.3rem .4rem;font-weight:500">${escapeHtml(p.name||p.parameter_id||'—')}</td>
                    <td style="text-align:center;padding:.3rem .4rem">${p.order_num??'—'}</td>
                    <td style="text-align:center;padding:.3rem .4rem">${p.min_val??'—'}</td>
                    <td style="text-align:center;padding:.3rem .4rem">${p.max_val??'—'}</td>
                    <td style="text-align:right;padding:.3rem .4rem">
                        <i class="fas fa-trash-alt action-icon" style="color:#f87171"
                           data-hcpp-del="${p.cp_id||p.id}"></i>
                    </td></tr>`).join('')}
                </tbody></table>`;
            el.querySelectorAll('[data-hcpp-del]').forEach(btn=>
                btn.addEventListener('click',async()=>{
                    if(!confirm('Удалить этот параметр?'))return;
                    try{
                        await api.hoClassParam.delete(btn.dataset.hcppDel);
                        showToast('Параметр удалён');
                        await refreshList();
                    }catch(e){showToast(e.message,false);}
                }));
        }catch(e){el.innerHTML=`<p style="color:#f87171">${escapeHtml(e.message)}</p>`;}
    };

    modal.querySelector('#_hcp-refresh').onclick=refreshList;

    modal.querySelector('#_hcpp-add').onclick=async()=>{
        const paramId=modal.querySelector('#_hcpp-id').value.trim();
        if(!paramId){showToast('Укажите ID параметра',false);return;}
        const data={
            ho_class_id:hoClassId,
            parameter_id:paramId,
            order_num:parseInt(modal.querySelector('#_hcpp-ord').value)||0,
            min_val:modal.querySelector('#_hcpp-min').value!==''?parseFloat(modal.querySelector('#_hcpp-min').value):0,
            max_val:modal.querySelector('#_hcpp-max').value!==''?parseFloat(modal.querySelector('#_hcpp-max').value):0,
        };
        try{
            await api.hoClassParam.create(data);
            showToast('Параметр добавлен');
            modal.querySelector('#_hcpp-id').value='';
            await refreshList();
        }catch(e){showToast(e.message,false);}
    };

    modal.querySelector('#_hcp-copy').onclick=()=>{
        const cm=createModal(`
            <div class="modal-title">
                <i class="fas fa-copy" style="color:var(--accent)"></i>
                Скопировать параметры из класса ХО
            </div>
            <div class="form-group">
                <label>ID исходного класса ХО</label>
                <input id="_copy-from" class="input-dark" placeholder="UUID исходного класса">
            </div>
            <div id="_copyErr"></div>
            <div class="modal-footer">
                <button class="btn btn-secondary" id="_copy-cancel">Отмена</button>
                <button class="btn btn-primary"   id="_copy-ok">
                    <i class="fas fa-copy"></i> Скопировать</button>
            </div>`);
        cm.querySelector('#_copy-cancel').onclick=()=>cm.remove();
        cm.querySelector('#_copy-ok').onclick=async()=>{
            const fromId=cm.querySelector('#_copy-from').value.trim();
            if(!fromId){setFormError(cm,'Укажите ID исходного класса');return;}
            try{
                await api.hoClassParam.copyFromClass({from_class_id:fromId,to_class_id:hoClassId});
                showToast('Параметры скопированы');
                cm.remove();
                await refreshList();
            }catch(e){setFormError(cm,e.message);}
        };
    };

    await refreshList();
}

// ── Документы типа ХО ──────────────────────────────────────────────

async function openHoClassDocsManager(hoClassId, className){
    const modal=createModal(`
        <div class="modal-title">
            <i class="fas fa-file-alt" style="color:var(--accent)"></i>
            Документы типа ХО: ${escapeHtml(className)}
        </div>
        <div id="_hcd-list" style="margin-bottom:.8rem;min-height:60px">${_loadingHtml()}</div>
        <div style="border-top:1px solid var(--border);padding-top:.8rem">
            <div style="font-weight:600;margin-bottom:.5rem;font-size:.85rem">Добавить тип документа:</div>
            <div style="display:flex;gap:.5rem;align-items:flex-end">
                <div class="form-group" style="flex:1;margin:0">
                    <select id="_hcd-sel" class="input-dark">
                        <option value="">Загрузка…</option>
                    </select>
                </div>
                <button class="btn btn-primary" id="_hcd-add">
                    <i class="fas fa-plus"></i> Добавить</button>
            </div>
        </div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_hcd-close">Закрыть</button>
        </div>`);

    modal.querySelector('#_hcd-close').onclick=()=>modal.remove();

    api.documentClass.list().then(dcs=>{
        const sel=modal.querySelector('#_hcd-sel');
        dcs=dcs||[];
        sel.innerHTML=`<option value="">— выберите класс документа —</option>`+
            dcs.map(d=>`<option value="${escapeHtml(String(d.doc_class_id))}">${escapeHtml(d.name)}</option>`).join('');
    }).catch(()=>{});

    const refreshList=async()=>{
        const el=modal.querySelector('#_hcd-list');
        el.innerHTML=_loadingHtml();
        try{
            const list=await api.hoClassDoc.list(hoClassId)||[];
            if(!list.length){el.innerHTML=_emptyHtml('Типы документов не назначены');return;}
            el.innerHTML=list.map(d=>`
                <div style="border:1px solid var(--border);border-radius:.4rem;padding:.4rem .7rem;
                            margin-bottom:.3rem;display:flex;align-items:center;gap:.5rem">
                    <span style="flex:1">${escapeHtml(d.doc_class_name||d.doc_class_id)}</span>
                    <i class="fas fa-trash-alt action-icon" style="color:#f87171"
                       data-hcdd-del="${d.id}"></i>
                </div>`).join('');
            el.querySelectorAll('[data-hcdd-del]').forEach(btn=>
                btn.addEventListener('click',async()=>{
                    try{
                        await api.hoClassDoc.delete(btn.dataset.hcddDel);
                        showToast('Тип документа удалён');
                        await refreshList();
                    }catch(e){showToast(e.message,false);}
                }));
        }catch(e){el.innerHTML=`<p style="color:#f87171">${escapeHtml(e.message)}</p>`;}
    };

    modal.querySelector('#_hcd-add').onclick=async()=>{
        const docClassId=modal.querySelector('#_hcd-sel').value;
        if(!docClassId){showToast('Выберите тип документа',false);return;}
        try{
            await api.hoClassDoc.create({ho_class_id:hoClassId,doc_class_id:docClassId});
            showToast('Тип документа добавлен');
            modal.querySelector('#_hcd-sel').value='';
            await refreshList();
        }catch(e){showToast(e.message,false);}
    };

    await refreshList();
}

// ══════════════════════════════════════════════════════════════════
//  РАЗДЕЛ B — ОПЕРАЦИИ (главный экран)
// ══════════════════════════════════════════════════════════════════

async function openHoOperationsView(){
    const modal=createModal(`
        <div class="modal-title">
            <i class="fas fa-exchange-alt" style="color:var(--accent)"></i>
            Хозяйственные операции (ПР4)
        </div>
        <div style="display:flex;gap:.3rem;margin-bottom:.8rem;flex-wrap:wrap;
                    border-bottom:1px solid var(--border);padding-bottom:.6rem">
            <button class="btn btn-secondary _ho-tab active" data-tab="all">
                <i class="fas fa-list"></i> Все операции</button>
            <button class="btn btn-secondary _ho-tab" data-tab="search">
                <i class="fas fa-search"></i> Поиск по типу</button>
            ${currentRole==='admin'?`<button class="btn btn-secondary _ho-tab" data-tab="config">
                <i class="fas fa-cog"></i> Конфигурация</button>`:''}
        </div>
        <div id="_ho-body" style="min-height:240px">${_loadingHtml()}</div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_ho-close">Закрыть</button>
        </div>`,true);

    modal.querySelector('#_ho-close').onclick=()=>modal.remove();

    modal.querySelectorAll('._ho-tab').forEach(btn=>{
        btn.addEventListener('click',()=>{
            modal.querySelectorAll('._ho-tab').forEach(b=>b.classList.remove('active'));
            btn.classList.add('active');
            _renderHoTab(modal, btn.dataset.tab);
        });
    });

    await _renderHoTab(modal,'all');
}

async function _renderHoTab(modal, tab){
    const body=modal.querySelector('#_ho-body');
    body.innerHTML=_loadingHtml();
    try{
        if(tab==='all')    await _renderAllHoTab(modal,body);
        else if(tab==='search') _renderHoSearchTab(modal,body);
        else if(tab==='config') body.innerHTML='',openHoConfigManager();
    }catch(e){
        body.innerHTML=`<p style="color:#f87171">${escapeHtml(e.message)}</p>`;
    }
}

async function _renderAllHoTab(modal, body){
    // Load ho_class list for filter
    let hoClasses=[];
    try{hoClasses=await api.hoClass.list()||[];}catch(_){}
    const hoClassMap={};
    hoClasses.forEach(c=>{hoClassMap[String(c.ho_class_id)]=c.name;});

    body.innerHTML=`
        <div style="display:flex;gap:.5rem;margin-bottom:.8rem;flex-wrap:wrap;align-items:flex-end">
            <div class="form-group" style="margin:0;flex:1;min-width:200px">
                <label style="font-size:.78rem">Фильтр по типу ХО</label>
                ${_buildSelect('_ho-filter', hoClasses, 'ho_class_id', 'name', '', 'Все типы')}
            </div>
            <button class="btn btn-secondary" id="_ho-reload">
                <i class="fas fa-sync"></i> Обновить</button>
            ${currentRole==='admin'?`<button class="btn btn-primary" id="_ho-create">
                <i class="fas fa-plus"></i> Новая операция</button>`:''}
        </div>
        <div id="_ho-table">${_loadingHtml()}</div>`;

    const loadTable=async()=>{
        const el=body.querySelector('#_ho-table');
        el.innerHTML=_loadingHtml();
        const filterClassId=body.querySelector('#_ho-filter').value||null;
        try{
            const list=await api.ho.list(filterClassId)||[];
            if(!list.length){el.innerHTML=_emptyHtml('Нет хозяйственных операций');return;}
            el.innerHTML=`<div style="overflow-x:auto">
                <table style="width:100%;border-collapse:collapse;font-size:.83rem">
                    <thead><tr style="color:var(--muted);border-bottom:1px solid var(--border)">
                        <th style="text-align:left;padding:.3rem .4rem">ID</th>
                        <th style="text-align:left;padding:.3rem .4rem">Номер</th>
                        <th style="text-align:left;padding:.3rem .4rem">Дата</th>
                        <th style="text-align:left;padding:.3rem .4rem">Тип</th>
                        <th style="text-align:right;padding:.3rem .4rem">Сумма</th>
                        <th style="text-align:center;padding:.3rem .4rem">Статус</th>
                        <th style="text-align:center;padding:.3rem .4rem">Действия</th>
                    </tr></thead>
                    <tbody>${list.map(h=>`<tr style="border-bottom:1px solid #1e293b">
                        <td style="padding:.3rem .4rem;color:var(--muted);font-size:.75rem">${escapeHtml(shortId(h.ho_id))}</td>
                        <td style="padding:.3rem .4rem;font-weight:500">${escapeHtml(h.doc_number||'—')}</td>
                        <td style="padding:.3rem .4rem">${escapeHtml(_formatHoDate(h.doc_date))}</td>
                        <td style="padding:.3rem .4rem">${escapeHtml(h.ho_class_name||hoClassMap[String(h.ho_class_id)]||h.ho_class_id||'—')}</td>
                        <td style="padding:.3rem .4rem;text-align:right">${_formatHoMoney(h.total_amount)}</td>
                        <td style="padding:.3rem .4rem;text-align:center">${_statusBadge(h.status)}</td>
                        <td style="padding:.3rem .4rem;text-align:center;white-space:nowrap">
                            <i class="fas fa-eye action-icon" style="color:#60a5fa" title="Детали" data-ho-view="${h.ho_id}"></i>
                            ${currentRole==='admin'?`
                            <i class="fas fa-edit action-icon" style="color:#a78bfa" title="Редактировать" data-ho-edit="${h.ho_id}" data-ho-class="${h.ho_class_id}"></i>
                            <i class="fas fa-trash-alt action-icon" style="color:#f87171" title="Удалить" data-ho-del="${h.ho_id}" data-ho-dnum="${escapeHtml(h.doc_number||h.ho_id)}"></i>`:''}
                        </td>
                    </tr>`).join('')}
                    </tbody>
                </table></div>`;
            el.querySelectorAll('[data-ho-view]').forEach(btn=>
                btn.addEventListener('click',()=>openHoDetailsModal(btn.dataset.hoView)));
            el.querySelectorAll('[data-ho-edit]').forEach(btn=>
                btn.addEventListener('click',()=>openEditHoModal(btn.dataset.hoEdit,btn.dataset.hoClass,loadTable)));
            el.querySelectorAll('[data-ho-del]').forEach(btn=>
                btn.addEventListener('click',async()=>{
                    if(!confirm(`Удалить операцию «${btn.dataset.hoDnum}»?`))return;
                    try{
                        await api.ho.delete(btn.dataset.hoDel);
                        showToast('Операция удалена');
                        await loadTable();
                    }catch(e){showToast(e.message,false);}
                }));
        }catch(e){el.innerHTML=`<p style="color:#f87171">${escapeHtml(e.message)}</p>`;}
    };

    body.querySelector('#_ho-reload').onclick=loadTable;
    body.querySelector('#_ho-filter').onchange=loadTable;
    const createBtn=body.querySelector('#_ho-create');
    if(createBtn)createBtn.onclick=async()=>{
        // Ask for ho_class selection first
        let hoClassList=[];
        try{hoClassList=await api.hoClass.terminal()||[];}catch(_){
            try{hoClassList=await api.hoClass.list()||[];}catch(_){}
        }
        const sm=createModal(`
            <div class="modal-title">
                <i class="fas fa-plus-circle" style="color:var(--accent)"></i>
                Выберите тип операции
            </div>
            <div class="form-group">
                <label>Тип ХО *</label>
                ${_buildSelect('_pick-hc', hoClassList, 'ho_class_id', 'name')}
            </div>
            <div class="modal-footer">
                <button class="btn btn-secondary" id="_pick-cancel">Отмена</button>
                <button class="btn btn-primary" id="_pick-ok">Далее →</button>
            </div>`);
        sm.querySelector('#_pick-cancel').onclick=()=>sm.remove();
        sm.querySelector('#_pick-ok').onclick=()=>{
            const v=sm.querySelector('#_pick-hc').value;
            if(!v){showToast('Выберите тип',false);return;}
            sm.remove();
            openCreateHoModal(v,loadTable);
        };
    };

    await loadTable();
}

function _renderHoSearchTab(modal, body){
    body.innerHTML=`
        <div style="margin-bottom:.8rem">
            <div style="font-size:.85rem;color:var(--muted);margin-bottom:.6rem">
                Поиск хозяйственных операций по типу ХО (SQL-функция findByClass)
            </div>
            <div style="display:flex;gap:.5rem;align-items:flex-end">
                <div class="form-group" style="flex:1;margin:0">
                    <label style="font-size:.78rem">ID типа ХО</label>
                    <input id="_hos-id" class="input-dark" placeholder="UUID типа ХО">
                </div>
                <button class="btn btn-primary" id="_hos-search">
                    <i class="fas fa-search"></i> Найти</button>
            </div>
        </div>
        <div id="_hos-results"></div>`;

    body.querySelector('#_hos-search').onclick=async()=>{
        const id=body.querySelector('#_hos-id').value.trim();
        if(!id){showToast('Укажите ID типа ХО',false);return;}
        const res=body.querySelector('#_hos-results');
        res.innerHTML=_loadingHtml();
        try{
            const list=await api.ho.findByClass(id)||[];
            if(!list.length){res.innerHTML=_emptyHtml('Операции не найдены');return;}
            res.innerHTML=`<div style="overflow-x:auto">
                <table style="width:100%;border-collapse:collapse;font-size:.83rem">
                    <thead><tr style="color:var(--muted);border-bottom:1px solid var(--border)">
                        <th style="text-align:left;padding:.3rem .4rem">ID</th>
                        <th style="text-align:left;padding:.3rem .4rem">Номер</th>
                        <th style="text-align:left;padding:.3rem .4rem">Дата</th>
                        <th style="text-align:right;padding:.3rem .4rem">Сумма</th>
                        <th style="text-align:center;padding:.3rem .4rem">Статус</th>
                        <th></th>
                    </tr></thead>
                    <tbody>${list.map(h=>`<tr style="border-bottom:1px solid #1e293b">
                        <td style="padding:.3rem .4rem;color:var(--muted);font-size:.75rem">${escapeHtml(shortId(h.ho_id))}</td>
                        <td style="padding:.3rem .4rem;font-weight:500">${escapeHtml(h.doc_number||'—')}</td>
                        <td style="padding:.3rem .4rem">${escapeHtml(_formatHoDate(h.doc_date))}</td>
                        <td style="padding:.3rem .4rem;text-align:right">${_formatHoMoney(h.total_amount)}</td>
                        <td style="padding:.3rem .4rem;text-align:center">${_statusBadge(h.status)}</td>
                        <td style="padding:.3rem .4rem;text-align:center">
                            <i class="fas fa-eye action-icon" style="color:#60a5fa" data-hos-view="${h.ho_id}" title="Детали"></i>
                        </td>
                    </tr>`).join('')}
                    </tbody>
                </table></div>`;
            res.querySelectorAll('[data-hos-view]').forEach(btn=>
                btn.addEventListener('click',()=>openHoDetailsModal(btn.dataset.hosView)));
        }catch(e){res.innerHTML=`<p style="color:#f87171">${escapeHtml(e.message)}</p>`;}
    };
}

// ══════════════════════════════════════════════════════════════════
//  СОЗДАНИЕ ХО — Визард
// ══════════════════════════════════════════════════════════════════

async function openCreateHoModal(hoClassId, onCreated){
    // Load class roles + parameters concurrently
    let roles=[], params=[], shdList=[];
    try{[roles,params,shdList]=await Promise.all([
        api.hoClassRole.list(hoClassId).catch(()=>[]),
        api.hoClass.parameters(hoClassId).catch(()=>[]),
        api.shd.list().catch(()=>[]),
    ]);}catch(_){}
    roles=roles||[];
    params=params||[];
    shdList=shdList||[];

    let step=1;
    const totalSteps=4;
    const state={basic:{},actors:{},paramVals:{},positions:[]};

    const modal=createModal(`
        <div class="modal-title">
            <i class="fas fa-plus-circle" style="color:var(--accent)"></i>
            Новая хозяйственная операция
        </div>
        <div id="_wiz-progress" style="display:flex;gap:.3rem;margin-bottom:1rem">
            ${[1,2,3,4].map(i=>`
                <div data-step="${i}" style="flex:1;height:.35rem;border-radius:.2rem;
                     background:${i===1?'var(--accent)':'#334155'};transition:background .2s"></div>`).join('')}
        </div>
        <div id="_wiz-label" style="font-size:.78rem;color:var(--muted);margin-bottom:.8rem">
            Шаг 1 из 4: Основные данные
        </div>
        <div id="_wiz-body" style="min-height:180px"></div>
        <div id="_wizErr"></div>
        <div class="modal-footer" style="justify-content:space-between">
            <button class="btn btn-secondary" id="_wiz-back" disabled>
                <i class="fas fa-chevron-left"></i> Назад</button>
            <button class="btn btn-secondary" id="_wiz-cancel">Отмена</button>
            <button class="btn btn-primary"   id="_wiz-next">
                Далее <i class="fas fa-chevron-right"></i></button>
        </div>`,true);

    modal.querySelector('#_wiz-cancel').onclick=()=>modal.remove();

    const labels=['Основные данные','Участники (акторы)','Значения параметров','Позиции (товары/услуги)'];

    const updateProgress=()=>{
        modal.querySelectorAll('[data-step]').forEach(el=>{
            const s=parseInt(el.dataset.step);
            el.style.background=s<=step?'var(--accent)':'#334155';
        });
        modal.querySelector('#_wiz-label').textContent=`Шаг ${step} из ${totalSteps}: ${labels[step-1]}`;
        modal.querySelector('#_wiz-back').disabled=step===1;
        modal.querySelector('#_wiz-next').innerHTML=step===totalSteps
            ?'<i class="fas fa-check"></i> Создать'
            :'Далее <i class="fas fa-chevron-right"></i>';
    };

    const renderStep=()=>{
        const body=modal.querySelector('#_wiz-body');
        if(step===1){
            body.innerHTML=`
                <div class="form-group">
                    <label>Номер документа</label>
                    <input id="_w-docnum" class="input-dark" placeholder="№ ДОК-001" value="${escapeHtml(state.basic.doc_number||'')}">
                </div>
                <div class="form-group">
                    <label>Дата документа</label>
                    <input id="_w-docdate" class="input-dark" type="date" value="${escapeHtml(state.basic.doc_date||'')}">
                </div>
                <div class="form-group">
                    <label>Примечание</label>
                    <textarea id="_w-note" class="input-dark" rows="2">${escapeHtml(state.basic.note||'')}</textarea>
                </div>`;
        } else if(step===2){
            if(!roles.length){
                body.innerHTML=_emptyHtml('Для этого типа ХО роли не определены');
                return;
            }
            body.innerHTML=roles.map(r=>`
                <div class="form-group">
                    <label>${escapeHtml(r.role_name||r.ho_role_id)}
                        ${r.is_required?'<span style="color:#f87171">*</span>':''}</label>
                    ${_buildSelect('_w-actor-'+r.ho_role_id, shdList, 'shd_id', 'name',
                        state.actors[r.ho_role_id]||'',
                        r.is_required?'— обязательно —':'— не указан —')}
                </div>`).join('');
        } else if(step===3){
            if(!params.length){
                body.innerHTML=_emptyHtml('Параметры не определены для этого типа ХО');
                return;
            }
            // params — HoClassParameterFull: cp_id, param_id, name, param_type, min_val, max_val, measuring_unit
            body.innerHTML=params.map(p=>{
                const cpId=String(p.cp_id||p.id||'');
                const val=state.paramVals[cpId]||'';
                const rangeStr=(p.min_val||p.max_val)
                    ?`<span style="font-size:.76rem;color:var(--muted)"> [${p.min_val}–${p.max_val}]</span>`
                    :'';
                const unitStr=p.measuring_unit?` (${p.measuring_unit})`
                    :'';
                return`<div class="form-group">
                    <label>${escapeHtml((p.name||cpId)+unitStr)}${rangeStr}</label>
                    <input id="_w-pval-${escapeHtml(cpId)}" class="input-dark"
                           type="${p.param_type==='int'?'number':'text'}"
                           step="${p.param_type==='int'?'1':'any'}"
                           placeholder="${p.param_type==='real'||p.param_type==='int'?'Числовое значение':'Строковое значение'}"
                           value="${escapeHtml(String(val))}">
                </div>`;
            }).join('');
        } else if(step===4){
            body.innerHTML=`
                <div id="_pos-list" style="margin-bottom:.6rem">
                    ${state.positions.length?state.positions.map((pos,i)=>`
                        <div style="border:1px solid var(--border);border-radius:.35rem;padding:.4rem .6rem;
                                    margin-bottom:.3rem;display:flex;gap:.5rem;align-items:center;flex-wrap:wrap">
                            <span style="flex:1;font-size:.83rem">Emobile: ${escapeHtml(pos.emobile_id||'—')}</span>
                            <span style="font-size:.83rem">Кол: ${pos.quantity}</span>
                            <span style="font-size:.83rem">Цена: ${pos.unit_price}</span>
                            <i class="fas fa-trash-alt action-icon" style="color:#f87171" data-pos-del="${i}"></i>
                        </div>`).join(''):_emptyHtml('Нет позиций')}
                </div>
                <div style="border-top:1px solid var(--border);padding-top:.6rem">
                    <div style="display:flex;gap:.4rem;flex-wrap:wrap;align-items:flex-end">
                        <div class="form-group" style="flex:2;min-width:120px;margin:0">
                            <label style="font-size:.78rem">ID Emobile</label>
                            <input id="_pos-em" class="input-dark" placeholder="UUID">
                        </div>
                        <div class="form-group" style="flex:1;min-width:70px;margin:0">
                            <label style="font-size:.78rem">Кол-во</label>
                            <input id="_pos-qty" class="input-dark" type="number" placeholder="1" value="1" min="0.001">
                        </div>
                        <div class="form-group" style="flex:1;min-width:70px;margin:0">
                            <label style="font-size:.78rem">Цена</label>
                            <input id="_pos-price" class="input-dark" type="number" placeholder="0.00" value="0">
                        </div>
                        <button class="btn btn-secondary" id="_pos-add">
                            <i class="fas fa-plus"></i> Добавить</button>
                    </div>
                </div>`;

            body.querySelector('#_pos-add').onclick=()=>{
                const emId=body.querySelector('#_pos-em').value.trim();
                const qty=parseFloat(body.querySelector('#_pos-qty').value)||1;
                const price=parseFloat(body.querySelector('#_pos-price').value)||0;
                if(!emId){showToast('Укажите ID Emobile',false);return;}
                state.positions.push({emobile_id:emId,quantity:qty,unit_price:price});
                renderStep();
            };
            body.querySelectorAll('[data-pos-del]').forEach(btn=>
                btn.addEventListener('click',()=>{
                    state.positions.splice(parseInt(btn.dataset.posDel),1);
                    renderStep();
                }));
        }
    };

    const collectStep=()=>{
        if(step===1){
            state.basic.doc_number=modal.querySelector('#_w-docnum').value.trim()||null;
            state.basic.doc_date=modal.querySelector('#_w-docdate').value||null;
            state.basic.note=modal.querySelector('#_w-note').value.trim()||null;
        } else if(step===2){
            roles.forEach(r=>{
                const sel=modal.querySelector('#_w-actor-'+r.ho_role_id);
                if(sel) state.actors[r.ho_role_id]=sel.value||null;
            });
        } else if(step===3){
            params.forEach(p=>{
                const cpId=String(p.cp_id||p.id||'');
                const inp=modal.querySelector('#_w-pval-'+cpId);
                if(inp) state.paramVals[cpId]=inp.value!==''?inp.value:null;
            });
        }
    };

    modal.querySelector('#_wiz-back').onclick=()=>{
        collectStep();
        step--;
        updateProgress();
        renderStep();
        setFormError(modal,'');
    };

    modal.querySelector('#_wiz-next').onclick=async()=>{
        collectStep();
        setFormError(modal,'');

        if(step<totalSteps){
            step++;
            updateProgress();
            renderStep();
            return;
        }

        // Final submit
        const btn=modal.querySelector('#_wiz-next');
        btn.disabled=true;
        btn.innerHTML='<i class="fas fa-spinner fa-spin"></i>';
        try{
            // Step 1: create the HO
            const hoData={
                ho_class_id:hoClassId,
                doc_number:state.basic.doc_number,
                doc_date:state.basic.doc_date||null,
                note:state.basic.note,
                status:'draft',
            };
            const created=await api.ho.create(hoData);
            const hoId=created?.ho_id||created?.id;
            if(!hoId) throw new Error('Не получен ID новой операции');

            // Step 2: actors
            for(const [roleId,shdId] of Object.entries(state.actors)){
                if(!shdId) continue;
                await api.hoActor.create({ho_id:hoId,ho_role_id:roleId,shd_id:shdId}).catch(()=>{});
            }

            // Step 3: param values — cpId = ho_class_parameter_id
            for(const [cpId,val] of Object.entries(state.paramVals)){
                if(val===null||val===undefined||val==='') continue;
                // Определяем тип по param definition
                const paramDef=params.find(p=>String(p.cp_id||p.id||'')===String(cpId));
                const ptype=paramDef?paramDef.param_type:'real';
                const body={ho_id:hoId,ho_class_parameter_id:cpId};
                if(ptype==='int')        body.val_int=Math.round(parseFloat(val)||0);
                else if(ptype==='str')   body.val_str=String(val);
                else                     body.val_real=parseFloat(val)||0;
                await api.hoParamVal.create(body).catch(()=>{});
            }

            // Step 4: positions
            for(const pos of state.positions){
                await api.hoPosition.create({ho_id:hoId,...pos}).catch(()=>{});
            }

            showToast('Операция создана');
            modal.remove();
            onCreated&&onCreated();
        }catch(e){
            setFormError(modal,e.message);
            btn.disabled=false;
            btn.innerHTML='<i class="fas fa-check"></i> Создать';
        }
    };

    updateProgress();
    renderStep();
}

// ══════════════════════════════════════════════════════════════════
//  РЕДАКТИРОВАНИЕ ХО (базовые поля)
// ══════════════════════════════════════════════════════════════════

function openEditHoModal(hoId, hoClassId, onSaved){
    const m=createModal(`
        <div class="modal-title">
            <i class="fas fa-edit" style="color:var(--accent)"></i>
            Редактировать операцию
        </div>
        <div class="form-group">
            <label>Номер документа</label>
            <input id="_eh-num" class="input-dark" placeholder="№ ДОК-001">
        </div>
        <div class="form-group">
            <label>Дата документа</label>
            <input id="_eh-date" class="input-dark" type="date">
        </div>
        <div class="form-group">
            <label>Статус</label>
            <select id="_eh-status" class="input-dark">
                <option value="draft">Черновик</option>
                <option value="confirmed">Подтверждена</option>
                <option value="cancelled">Отменена</option>
            </select>
        </div>
        <div class="form-group">
            <label>Примечание</label>
            <textarea id="_eh-note" class="input-dark" rows="2"></textarea>
        </div>
        <div id="_ehErr"></div>
        <div class="modal-footer">
            <button class="btn btn-secondary" id="_eh-cancel">Отмена</button>
            <button class="btn btn-primary"   id="_eh-save">
                <i class="fas fa-check"></i> Сохранить</button>
        </div>`);

    m.querySelector('#_eh-cancel').onclick=()=>m.remove();

    api.ho.getById(hoId).then(h=>{
        if(!h)return;
        m.querySelector('#_eh-num').value=h.doc_number||'';
        m.querySelector('#_eh-date').value=_dateInputValue(h.doc_date);
        m.querySelector('#_eh-status').value=h.status||'draft';
        m.querySelector('#_eh-note').value=h.note||'';
    }).catch(()=>{});

    m.querySelector('#_eh-save').onclick=async()=>{
        const data={
            doc_number:m.querySelector('#_eh-num').value.trim()||null,
            doc_date:m.querySelector('#_eh-date').value||null,
            status:m.querySelector('#_eh-status').value,
            note:m.querySelector('#_eh-note').value.trim()||null,
        };
        try{
            await api.ho.update(hoId,data);
            showToast('Операция обновлена');
            m.remove();
            onSaved&&onSaved();
        }catch(e){setFormError(m,e.message);}
    };
}

// ══════════════════════════════════════════════════════════════════
//  ДЕТАЛИ ХО
// ══════════════════════════════════════════════════════════════════

async function openHoDetailsModal(hoId){
    const isAdmin=_isHoAdmin();
    const modal=createModal(`
        <div class="modal-title">
            <i class="fas fa-eye" style="color:var(--accent)"></i>
            Детали хозяйственной операции
        </div>
        <div id="_hod-info" style="margin-bottom:.8rem;padding:.6rem;border:1px solid var(--border);border-radius:.4rem">
            ${_loadingHtml()}
        </div>
        <div style="display:flex;gap:.3rem;margin-bottom:.6rem;flex-wrap:wrap;
                    border-bottom:1px solid var(--border);padding-bottom:.5rem">
            <button class="btn btn-secondary _hod-stab active" data-stab="actors">
                <i class="fas fa-users"></i> Акторы</button>
            <button class="btn btn-secondary _hod-stab" data-stab="params">
                <i class="fas fa-sliders-h"></i> Параметры</button>
            <button class="btn btn-secondary _hod-stab" data-stab="positions">
                <i class="fas fa-boxes"></i> Позиции</button>
            <button class="btn btn-secondary _hod-stab" data-stab="documents">
                <i class="fas fa-file-alt"></i> Документы</button>
        </div>
        <div id="_hod-sub" style="min-height:120px">${_loadingHtml()}</div>
        <div id="_hodErr"></div>
        <div class="modal-footer" style="justify-content:space-between">
            ${isAdmin?`<div style="display:flex;gap:.4rem">
                <button class="btn btn-primary"   id="_hod-confirm" title="Подтвердить операцию">
                    <i class="fas fa-check"></i> Подтвердить</button>
                <button class="btn btn-danger"    id="_hod-cancel-op" title="Отменить операцию">
                    <i class="fas fa-times"></i> Отменить оп.</button>
            </div>`:'<div></div>'}
            <button class="btn btn-secondary" id="_hod-close">Закрыть</button>
        </div>`,true);

    modal.querySelector('#_hod-close').onclick=()=>modal.remove();

    // Load & show basic info
    const loadInfo=async()=>{
        const el=modal.querySelector('#_hod-info');
        try{
            const h=await api.ho.getById(hoId);
            if(!h){el.innerHTML='<span style="color:var(--muted)">Не найдено</span>';return;}
            let className=h.ho_class_name||h.ho_class_id||'—';
            if(h.ho_class_id){
                try{
                    const hc=await api.hoClass.getById(h.ho_class_id);
                    className=hc?.name||className;
                }catch(_){}
            }
            el.innerHTML=`
                <div style="display:flex;flex-wrap:wrap;gap:.8rem;align-items:center">
                    <div><span style="color:var(--muted);font-size:.78rem">Тип ХО</span><br>
                        <strong>${escapeHtml(className)}</strong></div>
                    <div><span style="color:var(--muted);font-size:.78rem">Номер</span><br>
                        <strong>${escapeHtml(h.doc_number||'—')}</strong></div>
                    <div><span style="color:var(--muted);font-size:.78rem">Дата</span><br>
                        <strong>${escapeHtml(_formatHoDate(h.doc_date))}</strong></div>
                    <div><span style="color:var(--muted);font-size:.78rem">Сумма</span><br>
                        <strong>${_formatHoMoney(h.total_amount)}</strong></div>
                    <div>${_statusBadge(h.status)}</div>
                    ${h.note?`<div style="width:100%;font-size:.82rem;color:var(--muted)"><i class="fas fa-note-sticky"></i> ${escapeHtml(h.note)}</div>`:''}
                </div>`;
        }catch(e){el.innerHTML=`<span style="color:#f87171">${escapeHtml(e.message)}</span>`;}
    };

    if(isAdmin){
        modal.querySelector('#_hod-confirm').onclick=async()=>{
            if(!confirm('Подтвердить операцию?'))return;
            try{
                const h=await api.ho.getById(hoId);
                await api.ho.update(hoId,{
                    status:'confirmed',
                    doc_number:h?.doc_number||'',
                    doc_date:_dateInputValue(h?.doc_date),
                    total_amount:h?.total_amount||0,
                    note:h?.note||'',
                });
                showToast('Операция подтверждена');
                await loadInfo();
            }catch(e){showToast(e.message,false);}
        };
        modal.querySelector('#_hod-cancel-op').onclick=async()=>{
            if(!confirm('Отменить операцию?'))return;
            try{
                const h=await api.ho.getById(hoId);
                await api.ho.update(hoId,{
                    status:'cancelled',
                    doc_number:h?.doc_number||'',
                    doc_date:_dateInputValue(h?.doc_date),
                    total_amount:h?.total_amount||0,
                    note:h?.note||'',
                });
                showToast('Операция отменена');
                await loadInfo();
            }catch(e){showToast(e.message,false);}
        };
    }

    const renderSubTab=async(tab)=>{
        const sub=modal.querySelector('#_hod-sub');
        sub.innerHTML=_loadingHtml();
        try{
            if(tab==='actors')    await _renderHodActors(sub,hoId);
            else if(tab==='params')     await _renderHodParams(sub,hoId);
            else if(tab==='positions')  await _renderHodPositions(sub,hoId);
            else if(tab==='documents')  await _renderHodDocuments(sub,hoId);
        }catch(e){sub.innerHTML=`<p style="color:#f87171">${escapeHtml(e.message)}</p>`;}
    };

    modal.querySelectorAll('._hod-stab').forEach(btn=>{
        btn.addEventListener('click',()=>{
            modal.querySelectorAll('._hod-stab').forEach(b=>b.classList.remove('active'));
            btn.classList.add('active');
            renderSubTab(btn.dataset.stab);
        });
    });

    await loadInfo();
    await renderSubTab('actors');
}

// ── Акторы операции ────────────────────────────────────────────────

async function _renderHodActors(el, hoId){
    const isAdmin=_isHoAdmin();
    let shdList=[];
    let roleList=[];
    try{[shdList,roleList]=await Promise.all([api.shd.list(),api.hoRole.list()]);}catch(_){}
    shdList=shdList||[];
    roleList=roleList||[];
    const shdMap={};shdList.forEach(s=>shdMap[s.shd_id]=s.name);
    const roleMap={};roleList.forEach(r=>roleMap[r.ho_role_id]=r.name);

    const refresh=async()=>{
        const list=await api.hoActor.list(hoId)||[];
        el.innerHTML=`
            <div style="margin-bottom:.5rem">
                ${list.length?list.map(a=>`
                    <div style="border:1px solid var(--border);border-radius:.35rem;padding:.35rem .6rem;
                                margin-bottom:.3rem;display:flex;align-items:center;gap:.5rem">
                        <span style="min-width:100px;color:var(--muted);font-size:.8rem">
                            ${escapeHtml(roleMap[a.ho_role_id]||a.ho_role_id||'—')}</span>
                        <span style="flex:1">${escapeHtml(shdMap[a.shd_id]||a.shd_id||'—')}</span>
                        ${isAdmin?`<i class="fas fa-trash-alt action-icon" style="color:#f87171" data-ha-del="${a.id}"></i>`:''}
                    </div>`).join(''):_emptyHtml('Акторы не назначены')}
            </div>
            ${isAdmin?`<div style="border-top:1px solid var(--border);padding-top:.5rem;
                        display:flex;gap:.4rem;flex-wrap:wrap;align-items:flex-end">
                <div class="form-group" style="flex:1;min-width:120px;margin:0">
                    <label style="font-size:.78rem">Роль</label>
                    ${_buildSelect('_ha-role', roleList, 'ho_role_id', 'name')}
                </div>
                <div class="form-group" style="flex:1;min-width:140px;margin:0">
                    <label style="font-size:.78rem">СХД</label>
                    ${_buildSelect('_ha-shd', shdList, 'shd_id', 'name')}
                </div>
                <button class="btn btn-primary" id="_ha-add">
                    <i class="fas fa-plus"></i></button>
            </div>`:''}`;

        if(isAdmin){
            el.querySelectorAll('[data-ha-del]').forEach(btn=>
                btn.addEventListener('click',async()=>{
                    try{await api.hoActor.delete(btn.dataset.haDel);showToast('Актор удалён');await refresh();}
                    catch(e){showToast(e.message,false);}
                }));
            el.querySelector('#_ha-add').onclick=async()=>{
                const roleId=el.querySelector('#_ha-role').value;
                const shdId=el.querySelector('#_ha-shd').value;
                if(!roleId||!shdId){showToast('Выберите роль и СХД',false);return;}
                try{
                    await api.hoActor.create({ho_id:hoId,ho_role_id:roleId,shd_id:shdId});
                    showToast('Актор добавлен');
                    await refresh();
                }catch(e){showToast(e.message,false);}
            };
        }
    };
    await refresh();
}

// ── Значения параметров операции ────────────────────────────────────

async function _renderHodParams(el, hoId){
    const isAdmin=_isHoAdmin();
    const refresh=async()=>{
        // list — HoParameterValueFull: id, ho_id, ho_class_parameter_id,
        //         param_name, param_type, measuring_unit, val_real, val_int, val_str, val_date, enum_val_id
        const [list, enumPositions]=await Promise.all([
            api.hoParamVal.list(hoId).catch(()=>[]),
            api.enumPosition.list().catch(()=>[]),
        ]);
        const enumMap={};
        (enumPositions||[]).forEach(p=>{enumMap[String(p.enum_position_id)]=p.value;});

        const rowHtml=list.map(p=>{
            // Определяем текущее отображаемое значение
            let displayVal='—', numVal='';
            if(p.param_type==='int'){
                displayVal=p.val_int!==0?String(p.val_int):'—';
                numVal=p.val_int!==0?String(p.val_int):'';
            } else if(p.param_type==='str'){
                displayVal=p.val_str||'—'; numVal=p.val_str||'';
            } else if(p.param_type==='date'){
                displayVal=_formatHoDate(p.val_date); numVal=_dateInputValue(p.val_date);
            } else if(p.param_type==='enum'){
                displayVal=enumMap[String(p.enum_val_id)]||p.enum_val_id||'—'; numVal=p.enum_val_id||'';
            } else { // real (default)
                displayVal=p.val_real!==0?String(p.val_real):'—';
                numVal=p.val_real!==0?String(p.val_real):'';
            }
            const unit=p.measuring_unit?` ${p.measuring_unit}`
                :'';
            return`<div style="border:1px solid var(--border);border-radius:.35rem;padding:.35rem .6rem;
                            margin-bottom:.3rem;display:flex;align-items:center;gap:.5rem;flex-wrap:wrap">
                <span style="min-width:140px;color:var(--muted);font-size:.8rem">
                    ${escapeHtml(p.param_name||p.ho_class_parameter_id||'—')}</span>
                <span style="flex:1;font-weight:500">${escapeHtml(displayVal+unit)}</span>
                ${isAdmin?`<input type="${p.param_type==='str'?'text':p.param_type==='date'?'date':'number'}" class="input-dark _hpv-val" data-pvid="${p.id}"
                       data-ptype="${escapeHtml(p.param_type||'real')}"
                       style="width:110px;padding:.25rem .4rem;font-size:.82rem"
                       value="${escapeHtml(numVal)}">
                <button class="btn btn-secondary _hpv-save" data-pvid="${p.id}"
                        style="font-size:.75rem;padding:.2rem .5rem">
                    <i class="fas fa-check"></i></button>
                <i class="fas fa-trash-alt action-icon" style="color:#f87171" data-hpv-del="${p.id}"></i>`:''}
            </div>`;
        }).join('');

        el.innerHTML=`
            <div style="margin-bottom:.5rem">
                ${list.length?rowHtml:_emptyHtml('Значения параметров не заданы')}
            </div>
            ${isAdmin?`<div style="border-top:1px solid var(--border);padding-top:.5rem;
                        display:flex;gap:.4rem;flex-wrap:wrap;align-items:flex-end">
                <div class="form-group" style="flex:2;min-width:180px;margin:0">
                    <label style="font-size:.78rem">ID параметра (типа ХО)</label>
                    <input id="_hpv-pid" class="input-dark" placeholder="ho_class_parameter_id">
                </div>
                <div class="form-group" style="flex:1;min-width:90px;margin:0">
                    <label style="font-size:.78rem">Значение</label>
                    <input id="_hpv-val" class="input-dark" type="number" placeholder="0">
                </div>
                <button class="btn btn-primary" id="_hpv-add"><i class="fas fa-plus"></i></button>
            </div>`:''}`;

        if(isAdmin){
            // Сохранение измененного значения
            el.querySelectorAll('._hpv-save').forEach(btn=>{
                btn.addEventListener('click',async()=>{
                    const pvId=btn.dataset.pvid;
                    const inp=el.querySelector(`._hpv-val[data-pvid="${pvId}"]`);
                    const ptype=inp?inp.dataset.ptype:'real';
                    const rawVal=inp?inp.value:'';
                    const body={};
                    if(ptype==='int')      body.val_int=Math.round(parseFloat(rawVal)||0);
                    else if(ptype==='str') body.val_str=rawVal;
                    else if(ptype==='date')body.val_date=rawVal;
                    else if(ptype==='enum')body.enum_val_id=rawVal;
                    else                   body.val_real=parseFloat(rawVal)||0;
                    try{await api.hoParamVal.update(pvId,body);showToast('Значение обновлено');await refresh();}
                    catch(e){showToast(e.message,false);}
                });
            });
            // Удаление
            el.querySelectorAll('[data-hpv-del]').forEach(btn=>{
                btn.addEventListener('click',async()=>{
                    try{await api.hoParamVal.delete(btn.dataset.hpvDel);showToast('Значение удалено');await refresh();}
                    catch(e){showToast(e.message,false);}
                });
            });
            // Добавление нового значения
            el.querySelector('#_hpv-add').onclick=async()=>{
                const cpId=el.querySelector('#_hpv-pid').value.trim();
                const rawVal=el.querySelector('#_hpv-val').value;
                if(!cpId){showToast('Укажите ho_class_parameter_id',false);return;}
                try{
                    await api.hoParamVal.create({
                        ho_id:hoId,
                        ho_class_parameter_id:cpId,
                        val_real:parseFloat(rawVal)||0,
                    });
                    showToast('Значение добавлено');
                    await refresh();
                }catch(e){showToast(e.message,false);}
            };
        }
    };
    await refresh();
}

// ── Позиции операции ────────────────────────────────────────────────

async function _renderHodPositions(el, hoId){
    const isAdmin=_isHoAdmin();
    const refresh=async()=>{
        const list=await api.hoPosition.list(hoId)||[];
        el.innerHTML=`
            <div style="margin-bottom:.5rem">
                ${list.length?`<div style="overflow-x:auto"><table style="width:100%;border-collapse:collapse;font-size:.82rem">
                    <thead><tr style="color:var(--muted);border-bottom:1px solid var(--border)">
                        <th style="text-align:left;padding:.25rem .4rem">Emobile</th>
                        <th style="text-align:center;padding:.25rem .4rem">Кол-во</th>
                        <th style="text-align:right;padding:.25rem .4rem">Цена</th>
                        <th style="text-align:right;padding:.25rem .4rem">Сумма</th>
                        ${isAdmin?'<th></th>':''}
                    </tr></thead>
                    <tbody>${list.map(p=>`<tr style="border-bottom:1px solid #1e293b">
                        <td style="padding:.25rem .4rem">${escapeHtml(p.emobile_name||p.emobile_id||'—')}</td>
                        <td style="text-align:center;padding:.25rem .4rem">${p.quantity??'—'}</td>
                        <td style="text-align:right;padding:.25rem .4rem">${_formatHoMoney(p.unit_price)}</td>
                        <td style="text-align:right;padding:.25rem .4rem">${_formatHoMoney(p.total_price)}</td>
                        ${isAdmin?`<td style="text-align:right;padding:.25rem .4rem">
                            <i class="fas fa-trash-alt action-icon" style="color:#f87171" data-hpp-del="${p.id}"></i>
                        </td>`:''}</tr>`).join('')}
                    </tbody></table></div>`
                :_emptyHtml('Позиции не добавлены')}
            </div>
            ${isAdmin?`<div style="border-top:1px solid var(--border);padding-top:.5rem;
                        display:flex;gap:.4rem;flex-wrap:wrap;align-items:flex-end">
                <div class="form-group" style="flex:2;min-width:120px;margin:0">
                    <label style="font-size:.78rem">ID Emobile</label>
                    <input id="_hpp-em" class="input-dark" placeholder="UUID">
                </div>
                <div class="form-group" style="flex:1;min-width:70px;margin:0">
                    <label style="font-size:.78rem">Кол-во</label>
                    <input id="_hpp-qty" class="input-dark" type="number" placeholder="1" value="1">
                </div>
                <div class="form-group" style="flex:1;min-width:70px;margin:0">
                    <label style="font-size:.78rem">Цена</label>
                    <input id="_hpp-price" class="input-dark" type="number" placeholder="0">
                </div>
                <button class="btn btn-primary" id="_hpp-add"><i class="fas fa-plus"></i> Добавить</button>
            </div>`:''}`;

        if(isAdmin){
            el.querySelectorAll('[data-hpp-del]').forEach(btn=>
                btn.addEventListener('click',async()=>{
                    if(!confirm('Удалить позицию?'))return;
                    try{await api.hoPosition.delete(btn.dataset.hppDel);showToast('Позиция удалена');await refresh();}
                    catch(e){showToast(e.message,false);}
                }));
            el.querySelector('#_hpp-add').onclick=async()=>{
                const emId=el.querySelector('#_hpp-em').value.trim();
                const qty=parseFloat(el.querySelector('#_hpp-qty').value)||1;
                const price=parseFloat(el.querySelector('#_hpp-price').value)||0;
                if(!emId){showToast('Укажите ID Emobile',false);return;}
                try{
                    await api.hoPosition.create({ho_id:hoId,emobile_id:emId,quantity:qty,unit_price:price});
                    showToast('Позиция добавлена');
                    await refresh();
                }catch(e){showToast(e.message,false);}
            };
        }
    };
    await refresh();
}

// ── Документы операции ────────────────────────────────────────────

async function _renderHodDocuments(el, hoId){
    const isAdmin=_isHoAdmin();
    let docClassList=[];
    try{docClassList=await api.documentClass.list()||[];}catch(_){}
    const dcMap={};docClassList.forEach(d=>dcMap[d.doc_class_id]=d.name);

    const refresh=async()=>{
        const list=await api.hoDoc.list(hoId)||[];
        el.innerHTML=`
            <div style="margin-bottom:.5rem">
                ${list.length?list.map(d=>`
                    <div style="border:1px solid var(--border);border-radius:.35rem;padding:.35rem .6rem;
                                margin-bottom:.3rem;display:flex;align-items:center;gap:.5rem">
                        <span style="flex:1">${escapeHtml(dcMap[d.doc_class_id]||d.doc_class_id||'—')}</span>
                        ${d.doc_number?`<span style="font-size:.78rem;color:var(--muted)">${escapeHtml(d.doc_number)}</span>`:''}
                        ${d.doc_date?`<span style="font-size:.78rem;color:var(--muted)">${escapeHtml(_formatHoDate(d.doc_date))}</span>`:''}
                        ${isAdmin?`<i class="fas fa-trash-alt action-icon" style="color:#f87171" data-hd-del="${d.id}"></i>`:''}
                    </div>`).join(''):_emptyHtml('Документы не прикреплены')}
            </div>
            ${isAdmin?`<div style="border-top:1px solid var(--border);padding-top:.5rem;
                        display:flex;gap:.4rem;flex-wrap:wrap;align-items:flex-end">
                <div class="form-group" style="flex:1;min-width:140px;margin:0">
                    <label style="font-size:.78rem">Класс документа</label>
                    ${_buildSelect('_hd-cls', docClassList, 'doc_class_id', 'name')}
                </div>
                <div class="form-group" style="flex:1;min-width:100px;margin:0">
                    <label style="font-size:.78rem">Номер документа</label>
                    <input id="_hd-num" class="input-dark" placeholder="№ —">
                </div>
                <button class="btn btn-primary" id="_hd-add"><i class="fas fa-plus"></i> Добавить</button>
            </div>`:''}`;

        if(isAdmin){
            el.querySelectorAll('[data-hd-del]').forEach(btn=>
                btn.addEventListener('click',async()=>{
                    try{await api.hoDoc.delete(btn.dataset.hdDel);showToast('Документ удалён');await refresh();}
                    catch(e){showToast(e.message,false);}
                }));
            el.querySelector('#_hd-add').onclick=async()=>{
                const dcId=el.querySelector('#_hd-cls').value;
                const num=el.querySelector('#_hd-num').value.trim();
                if(!dcId){showToast('Выберите класс документа',false);return;}
                try{
                    await api.hoDoc.create({ho_id:hoId,doc_class_id:dcId,doc_number:num||null});
                    showToast('Документ добавлен');
                    await refresh();
                }catch(e){showToast(e.message,false);}
            };
        }
    };
    await refresh();
}

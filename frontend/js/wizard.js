/* wizard.js */
async function openWizard(){
    const STEPS=[
        {key:'powerPoint',    label:'Силовая установка', listKey:'powerPoints',    idField:'power_point_id',    nameField:null},
        {key:'battery',       label:'Батарея',           listKey:'batteries',      idField:'battery_id',        nameField:'battery_name'},
        {key:'chargerSystem', label:'Зарядная система',  listKey:'chargerSystems', idField:'charger_system_id', nameField:null},
        {key:'chassis',       label:'Шасси',             listKey:'chassis',        idField:'chassis_id',        nameField:null},
        {key:'body',          label:'Кузов',             listKey:'bodies',         idField:'body_id',           nameField:null},
        {key:'electronics',   label:'Электроника',       listKey:'electronics',    idField:'electronics_id',    nameField:null},
        {key:'carName',       label:'Название',          listKey:null},
    ];
    const lists={};
    try{
        await Promise.all(STEPS.filter(s=>s.listKey).map(async s=>{
            lists[s.key]=await api[s.listKey].list();
        }));
    }catch(e){showError(`Ошибка загрузки: ${e.message}`);return;}

    let stepIdx=0;const sel={};let modal=null;

    const render=()=>{
        if(modal)modal.remove();
        const step=STEPS[stepIdx];const isLast=stepIdx===STEPS.length-1;
        const dots=STEPS.map((s,i)=>
            `<span class="step-dot ${i<stepIdx?'done':i===stepIdx?'active':''}">${escapeHtml(s.label)}</span>`
        ).join('');
        let body;
        if(isLast){
            body=`<div class="form-group"><label>Название автомобиля</label>
                  <input id="_wn" type="text" class="input-dark"
                      value="${escapeHtml(sel.carName||'')}" placeholder="Например: Emobile X1"></div>`;
        }else{
            const items=lists[step.key]||[];
            const catCfg=getCatConfig(step.listKey);
            const opts=items.map(i=>{
                const iId=i[step.idField];
                const iN=step.nameField
                    ?i[step.nameField]
                    :(catCfg?`${catCfg.label} #${shortId(iId)}`:`#${shortId(iId)}`);
                return`<option value="${escapeHtml(iId)}" ${sel[step.key]==iId?'selected':''}>${escapeHtml(iN)}</option>`;
            }).join('');
            body=`<div class="form-group"><label>${escapeHtml(step.label)}</label>
                  <select id="_ws" class="input-dark">
                      <option value="">— выберите —</option>${opts}</select></div>
                  ${items.length===0
                    ?`<p style="font-size:.85rem;color:var(--muted)">⚠ Список пуст — сначала создайте эту деталь.</p>`:''}`;
        }
        modal=createModal(`
            <div class="modal-title">
                <i class="fas fa-car-side" style="color:var(--accent)"></i> Сборка автомобиля</div>
            <div class="step-indicator">${dots}</div>
            ${body}
            <div class="error-banner hidden"><i class="fas fa-exclamation-circle"></i><span></span></div>
            <div class="modal-footer" style="justify-content:space-between">
                <button class="btn btn-secondary" id="_wp" ${stepIdx===0?'disabled':''}>← Назад</button>
                <button class="btn btn-secondary" id="_wc">Отмена</button>
                <button class="btn btn-primary"   id="_wn2">
                    ${isLast?'<i class="fas fa-check"></i> Создать':'Далее →'}</button>
            </div>`,true);

        modal.querySelector('#_wc').onclick=()=>modal.remove();
        modal.querySelector('#_wp').onclick=()=>{stepIdx--;render();};
        modal.querySelector('#_wn2').onclick=async()=>{
            setFormError(modal,'');
            if(isLast){
                const carName=modal.querySelector('#_wn')?.value.trim();
                if(!carName){setFormError(modal,'Введите название автомобиля');return;}
                sel.carName=carName;
                const btn=modal.querySelector('#_wn2');
                btn.disabled=true;btn.innerHTML='<i class="fas fa-spinner fa-spin"></i> Создаю…';
                try{
                    /* БАГ 4 FIX: ID отправляем как string (Go struct UpdateEmobileInput.power_point_id: string) */
                    await api.emobiles.create({
                        name:               sel.carName,
                        power_point_id:     String(sel.powerPoint),
                        battery_id:         String(sel.battery),
                        charger_system_id:  String(sel.chargerSystem),
                        chassis_id:         String(sel.chassis),
                        body_id:            String(sel.body),
                        electronics_id:     String(sel.electronics),
                    });
                    await loadCatalogData();modal.remove();
                }catch(e){
                    setFormError(modal,e.message);
                    btn.disabled=false;btn.innerHTML='<i class="fas fa-check"></i> Создать';
                }
            }else{
                const val=modal.querySelector('#_ws')?.value;
                if(!val){setFormError(modal,'Выберите значение');return;}
                sel[step.key]=val;stepIdx++;render();
            }
        };
    };
    render();
}

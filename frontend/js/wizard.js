/* wizard.js */
async function openWizard(){
    const STEPS=[
        {key:'powerPoint',    label:'Силовая установка',  listKey:'powerPoints',    idField:'power_point_id',    nameField:null},
        {key:'battery',       label:'Батарея',            listKey:'batteries',      idField:'battery_id',        nameField:'battery_name'},
        {key:'chargerSystem', label:'Зарядная система',   listKey:'chargerSystems', idField:'charger_system_id', nameField:null},
        {key:'chassis',       label:'Шасси',              listKey:'chassis',        idField:'chassis_id',        nameField:null},
        {key:'body',          label:'Кузов',              listKey:'bodies',         idField:'body_id',           nameField:null},
        {key:'electronics',   label:'Электроника',        listKey:'electronics',    idField:'electronics_id',    nameField:null},
        {key:'carName',       label:'Название',           listKey:null},
    ];
    let lists={};
    try{
        await Promise.all(STEPS.filter(s=>s.listKey).map(async s=>{lists[s.key]=await api[s.listKey].list();}));
    }catch(e){showError(`Ошибка загрузки: ${e.message}`);return;}
    let stepIdx=0; const sel={}; let modal=null;
    const render=()=>{
        if(modal) modal.remove();
        const step=STEPS[stepIdx], isLast=stepIdx===STEPS.length-1;
        const dots=STEPS.map((s,i)=>`<span class="step-dot ${i<stepIdx?'done':i===stepIdx?'active':''}">${escapeHtml(s.label)}</span>`).join('');
        let body=isLast
            ?`<div class="form-group"><label>Название автомобиля</label><input id="_wn" type="text" class="input-dark" value="${escapeHtml(sel.carName||'')}" placeholder="Например: Emobile X1"></div>`
            :(()=>{
                const items=lists[step.key]||[];
                const opts=items.map(i=>{
                    const iId=i[step.idField];
                    const iN=step.nameField?i[step.nameField]:`${step.label} …${shortId(iId)}`;
                    return `<option value="${escapeHtml(iId)}" ${sel[step.key]===iId?'selected':''}>${escapeHtml(iN)}</option>`;
                }).join('');
                return `<div class="form-group"><label>${escapeHtml(step.label)}</label><select id="_ws" class="input-dark"><option value="">— выберите —</option>${opts}</select></div>${items.length===0?'<p class="text-muted">Список пуст — сначала создайте эту деталь.</p>':''}`;
            })();
        modal=createModal(`
            <div class="modal-title"><i class="fas fa-car-side" style="color:var(--accent)"></i> Сборка автомобиля</div>
            <div class="step-indicator">${dots}</div>
            ${body}
            <div id="_we" class="text-danger" style="font-size:0.85rem;min-height:1.2em"></div>
            <div class="modal-footer" style="justify-content:space-between">
                <button class="btn btn-secondary" id="_wp" ${stepIdx===0?'disabled':''}>← Назад</button>
                <button class="btn btn-secondary" id="_wc">Отмена</button>
                <button class="btn btn-primary" id="_wn2">${isLast?'<i class="fas fa-check"></i> Создать':'Далее →'}</button>
            </div>`,true);
        modal.querySelector('#_wc').onclick=()=>modal.remove();
        modal.querySelector('#_wp').onclick=()=>{stepIdx--;render();};
        modal.querySelector('#_wn2').onclick=async()=>{
            if(isLast){
                const nameEl=modal.querySelector('#_wn');
                const carName=nameEl?nameEl.value.trim():'';
                if(!carName){modal.querySelector('#_we').textContent='Введите название';return;}
                sel.carName=carName;
                try{
                    // Поле называется "name" согласно CreateEmobileInput
                    await api.emobiles.create({
                        name:             sel.carName,
                        power_point_id:   sel.powerPoint,
                        battery_id:       sel.battery,
                        charger_system_id:sel.chargerSystem,
                        chassis_id:       sel.chassis,
                        body_id:          sel.body,
                        electronics_id:   sel.electronics,
                    });
                    await loadCatalogData(); modal.remove();
                }catch(e){modal.querySelector('#_we').textContent=`Ошибка: ${e.message}`;}
            }else{
                const val=modal.querySelector('#_ws')?.value;
                if(!val){modal.querySelector('#_we').textContent='Выберите значение';return;}
                sel[step.key]=val; stepIdx++; render();
            }
        };
    };
    render();
}

/* catalog.js — с добавлением кнопок ПР2, ПР3 и ПР4 */
let allPartsData=[],categoryFilter='all',searchText='';

async function loadCatalogData(){
    renderCatalog(false,true);
    const raw={};
    const loadErrors=[];
    for(const cat of CATEGORY_MAP){
        try{raw[cat.key]=await api[cat.key].list()||[];}
        catch(e){
            raw[cat.key]=[];
            loadErrors.push(`${cat.label}: ${e.message}`);
            console.error(`[catalog] ${cat.label}:`,e.message);
        }
    }
    if(loadErrors.length){
        const warn=document.createElement('div');
        warn.style.cssText='position:fixed;bottom:1rem;right:1rem;z-index:9999;'+
            'background:#7f1d1d;color:#fca5a5;border:1px solid #ef4444;'+
            'border-radius:.5rem;padding:.75rem 1rem;font-size:.8rem;max-width:360px;'+
            'box-shadow:0 4px 16px rgba(0,0,0,.5)';
        warn.innerHTML='<strong>⚠ Не удалось загрузить:</strong><br>'+loadErrors.join('<br>');
        document.body.appendChild(warn);
        setTimeout(()=>warn.remove(),12000);
    }
    const nameFor={};
    for(const cat of CATEGORY_MAP){
        nameFor[cat.idField]={};
        for(const item of raw[cat.key]){
            nameFor[cat.idField][item[cat.idField]]=cat.nameField
                ?item[cat.nameField]
                :`${cat.label} #${shortId(item[cat.idField])}`;
        }
    }

    // ── ПР̣3: параметры автомобилей ────────────────────────
    let paramDefsById={};      // cp_id → {name, param_type, measuring_unit}
    let paramValsByEmobile={}; // emobile_id → [EmobileParameterValue]
    try{
        const [defs,vals]=await Promise.all([
            api.componentParameter.byType('emobile').catch(()=>[]),
            api.emobileParam.list().catch(()=>[]),
        ]);
        (defs||[]).forEach(d=>{ paramDefsById[String(d.cp_id)]=d; });
        (vals||[]).forEach(v=>{
            const eid=String(v.emobile_id);
            (paramValsByEmobile[eid]=paramValsByEmobile[eid]||[]).push(v);
        });
    }catch(_){}

    const getInfo={
        powerPoints:   i=>[
            {l:'Двиг',v:nameFor['engine_id']?.[i.engine_id]            ||'?'},
            {l:'Инв', v:nameFor['inverter_id']?.[i.inverter_id]         ||'?'},
            {l:'КПП', v:nameFor['gearbox_id']?.[i.gearbox_id]           ||'?'},
        ],
        chassis:       i=>[
            {l:'Рама', v:nameFor['frame_id']?.[i.frame_id]              ||'?'},
            {l:'Подв', v:nameFor['suspension_id']?.[i.suspension_id]    ||'?'},
            {l:'Торм', v:nameFor['break_system_id']?.[i.break_system_id]||'?'},
        ],
        chargerSystems:i=>[
            {l:'Заряд',v:nameFor['charger_id']?.[i.charger_id]          ||'?'},
            {l:'Конн', v:nameFor['connector_id']?.[i.connector_id]      ||'?'},
        ],
        electronics:   i=>[
            {l:'Контр',v:nameFor['controller_id']?.[i.controller_id]    ||'?'},
            {l:'Сенс', v:nameFor['sensor_id']?.[i.sensor_id]            ||'?'},
            {l:'Пров', v:nameFor['wiring_id']?.[i.wiring_id]            ||'?'},
        ],
        bodies:        i=>[
            {l:'Карк',v:nameFor['carcass_id']?.[i.carcass_id]           ||'?'},
            {l:'Дв',  v:nameFor['doors_id']?.[i.doors_id]               ||'?'},
            {l:'Кр',  v:nameFor['wings_id']?.[i.wings_id]               ||'?'},
        ],
        emobiles: i=>{
            const baseLines=[
                {l:'СУ',     v:nameFor['power_point_id']?.[i.power_point_id]      ||'?'},
                {l:'Батарея',v:nameFor['battery_id']?.[i.battery_id]              ||'?'},
                {l:'Зарядка',v:nameFor['charger_system_id']?.[i.charger_system_id]||'?'},
                {l:'Шасси',  v:nameFor['chassis_id']?.[i.chassis_id]              ||'?'},
                {l:'Кузов',  v:nameFor['body_id']?.[i.body_id]                    ||'?'},
                {l:'Электро',v:nameFor['electronics_id']?.[i.electronics_id]      ||'?'},
            ];
            // Приписываем значения параметров ПР̣3 (если есть)
            const pvList=paramValsByEmobile[String(i.emobile_id)]||[];
            const paramLines=pvList.map(pv=>{
                const def=paramDefsById[String(pv.component_parameter_id)];
                const name=def?def.name:'?';
                const unit=def&&def.measuring_unit?` ${def.measuring_unit}`:'';
                let val='—';
                if(def){
                    if(def.param_type==='int')       val=pv.val_int!==0?String(pv.val_int):'—';
                    else if(def.param_type==='str')  val=pv.val_str||'—';
                    else if(def.param_type==='enum') val=pv.enum_val_id||'—';
                    else                             val=pv.val_real!==0?String(pv.val_real):'—';
                }
                return {l:name, v:val+unit};
            });
            return [...baseLines,...paramLines];
        },
    };
    const flat=[];
    for(const cat of CATEGORY_MAP){
        const resolver=getInfo[cat.key]||null;
        for(const item of raw[cat.key]){
            const name=cat.nameField
                ?(item[cat.nameField]||`${cat.label} #${shortId(item[cat.idField])}`)
                :`${cat.label} …${shortId(item[cat.idField])}`;
            const infoLines=resolver?resolver(item):null;
            const infoStr=resolver?infoLines.map(r=>`${r.l}: ${r.v}`).join(' · '):cat.infoFn(item);
            flat.push({id:item[cat.idField],name,category:cat.label,
                       categoryKey:cat.key,infoLines,infoStr,rawData:item});
        }
    }
    allPartsData=flat;
    renderCatalog(false,false);
}

function renderInfoCell(p){
    if(!p.infoLines)
        return`<td class="info-cell">${escapeHtml(p.infoStr)}</td>`;
    const lines=p.infoLines.map(r=>
        `<span class="info-line">
            <span class="info-label">${escapeHtml(r.l)}</span>
            <span class="info-value" title="${escapeHtml(r.v)}">${escapeHtml(r.v)}</span>
         </span>`
    ).join('');
    return`<td class="info-cell"><div class="info-lines">${lines}</div></td>`;
}

function _buildChipsHtml(){
    const allLabels=CATEGORY_MAP.map(c=>c.label);
    return['all',...allLabels].map(l=>
        `<span class="filter-chip ${categoryFilter===l?'active':''}" data-cat="${escapeHtml(l)}">
            ${l==='all'?'Все':escapeHtml(l)}</span>`).join('');
}

function _buildRowsHtml(isAdmin){
    const q=searchText.toLowerCase();
    const rows=allPartsData.filter(p=>{
        const mc=categoryFilter==='all'||p.category===categoryFilter;
        return mc&&(!q||p.name.toLowerCase().includes(q)||String(p.id).includes(q));
    });
    if(!rows.length)return`<tr><td colspan="${isAdmin?5:4}">
        <div class="empty-state"><i class="fas fa-inbox"></i> Нет данных</div></td></tr>`;
    return rows.map(p=>`<tr>
        <td><span class="text-mono">${escapeHtml(shortId(p.id))}</span></td>
        <td style="font-weight:500">${escapeHtml(p.name)}</td>
        <td><span class="badge badge-default">${escapeHtml(p.category)}</span></td>
        ${renderInfoCell(p)}
        ${isAdmin?`<td style="white-space:nowrap">
            <i class="fas fa-edit action-icon" style="color:#60a5fa"
               data-action="edit"   data-id="${escapeHtml(p.id)}"
               data-cat="${escapeHtml(p.categoryKey)}" title="Редактировать"></i>
            ${p.categoryKey==='emobiles'?`<i class="fas fa-sliders-h action-icon" style="color:#a78bfa"
               data-action="params" data-id="${escapeHtml(p.id)}"
               data-name="${escapeHtml(p.name)}" title="Параметры (ПР3)"></i>`:''}
            <i class="fas fa-trash-alt action-icon" style="color:#f87171"
               data-action="delete" data-id="${escapeHtml(p.id)}"
               data-cat="${escapeHtml(p.categoryKey)}"
               data-name="${escapeHtml(p.name)}" title="Удалить"></i>
        </td>`:''}`
    ).join('');
}

function _buildCatalogMetrics(){
    const carCount=allPartsData.filter(p=>p.categoryKey==='emobiles').length;
    const compositeCount=allPartsData.filter(p=>getCatConfig(p.categoryKey)?.composite).length;
    const categoryCount=new Set(allPartsData.map(p=>p.categoryKey)).size;
    return [
        {icon:'fa-database', value:allPartsData.length, label:'записей в справочнике'},
        {icon:'fa-car-side', value:carCount, label:'готовых автомобилей'},
        {icon:'fa-puzzle-piece', value:compositeCount, label:'составных узлов'},
        {icon:'fa-layer-group', value:categoryCount||CATEGORY_MAP.length, label:'разделов каталога'},
    ].map(m=>`
        <div class="metric-card">
            <div style="display:flex;align-items:center;justify-content:space-between;gap:.5rem">
                <div>
                    <div class="metric-value">${escapeHtml(m.value)}</div>
                    <div class="metric-label">${escapeHtml(m.label)}</div>
                </div>
                <i class="fas ${m.icon}" style="color:var(--accent);font-size:1.25rem"></i>
            </div>
        </div>`).join('');
}

function _buildQuickActions(isAdmin){
    const adminActions=[
        {a:'createCar', icon:'fa-car', title:'Собрать автомобиль', text:'Пошаговый сценарий выбора узлов и создания изделия.'},
        {a:'addPart', icon:'fa-plus', title:'Добавить деталь', text:'Создание листовых и составных элементов справочника.'},
        {a:'enums', icon:'fa-list-ul', title:'Справочники', text:'ПР2: допустимые значения и порядок вывода.'},
        {a:'params', icon:'fa-sliders-h', title:'Параметры', text:'ПР3: характеристики изделий и значения автомобилей.'},
        {a:'ho', icon:'fa-exchange-alt', title:'Хоз. операции', text:'ПР4: оформление, поиск и просмотр операций.'},
        {a:'createUser', icon:'fa-user-plus', title:'Пользователи', text:'Создание учетных записей и назначение ролей.'},
    ];
    const userActions=[
        {a:'ho', icon:'fa-exchange-alt', title:'Хоз. операции', text:'Просмотр операций и документов без изменения настроек.'},
        {a:'guide', icon:'fa-route', title:'Карта интерфейса', text:'Информационная архитектура и основные пользовательские сценарии.'},
        {a:'focusSearch', icon:'fa-search', title:'Найти изделие', text:'Быстрый поиск по названию или идентификатору.'},
        {a:'resetFilters', icon:'fa-filter-circle-xmark', title:'Сбросить фильтры', text:'Вернуться ко всему каталогу.'},
    ];
    return (isAdmin?adminActions:userActions).map(item=>`
        <button class="action-card" type="button" data-catalog-action="${item.a}">
            <i class="fas ${item.icon}"></i>
            <strong>${escapeHtml(item.title)}</strong>
            <span>${escapeHtml(item.text)}</span>
        </button>`).join('');
}

function _buildSiteNav(){
    const items=[
        {page:'about', icon:'fa-circle-info', label:'О компании'},
        {page:'services', icon:'fa-screwdriver-wrench', label:'Услуги'},
        {page:'delivery', icon:'fa-truck-fast', label:'Доставка и цена'},
        {page:'contacts', icon:'fa-address-book', label:'Контакты'},
    ];
    return items.map(item=>`
        <button class="btn btn-secondary" type="button" data-site-page="${item.page}">
            <i class="fas ${item.icon}"></i> ${escapeHtml(item.label)}
        </button>`).join('');
}

function renderCatalog(partial=true,loading=false){
    const app=document.getElementById('app');
    if(!app)return;
    const isAdmin=currentRole==='admin';
    if(partial&&app.querySelector('.table-wrap')){
        app.querySelector('tbody').innerHTML=loading
            ?`<tr><td colspan="${isAdmin?5:4}" style="text-align:center;padding:2rem;color:var(--muted)">
               <i class="fas fa-spinner fa-spin"></i> Загрузка…</td></tr>`
            :_buildRowsHtml(isAdmin);
        const fc=app.querySelector('#filterChips');
        if(fc)fc.innerHTML=_buildChipsHtml();
        _attachActions(app,isAdmin);
        return;
    }
    const roleBadge=isAdmin
        ?`<span class="badge badge-admin"><i class="fas fa-shield-alt"></i> Администратор</span>`
        :`<span class="badge badge-user"><i class="fas fa-user"></i> Пользователь</span>`;
    app.innerHTML=`
    <div class="app-shell">
        <header class="topbar">
            <div class="topbar-inner">
                <div class="brand">
                    <span class="brand-icon"><i class="fas fa-charging-station"></i></span>
                    <span>Emobile <span class="brand-sub">Справочник изделий и хозяйственных операций</span></span>
                </div>
                <div class="header-actions">
                    ${roleBadge}
                    <button class="btn btn-secondary" data-catalog-action="guide">
                        <i class="fas fa-route"></i> Карта интерфейса</button>
                    <button class="btn btn-secondary" id="logoutBtn" title="Выйти">
                        <i class="fas fa-sign-out-alt"></i> Выйти</button>
                </div>
            </div>
            <nav class="site-nav" aria-label="Разделы сайта">
                ${_buildSiteNav()}
            </nav>
        </header>
        <main class="container">
            <section class="hero-card">
                <div>
                    <div class="hero-kicker">${isAdmin?'Рабочее место администратора':'Рабочее место пользователя'}</div>
                    <h1 class="hero-title">Каталог электромобилей</h1>
                    <p class="hero-text">
                        ${isAdmin
                            ?'Управляйте составом изделий, справочниками допустимых значений, параметрами автомобилей и хозяйственными операциями в одном интерфейсе.'
                            :'Просматривайте состав автомобилей, фильтруйте справочник и открывайте хозяйственные операции без риска изменить настройки системы.'}
                    </p>
                    <div class="hero-actions">
                        ${isAdmin?`
                            <button class="btn btn-primary" data-catalog-action="createCar">
                                <i class="fas fa-car"></i> Создать автомобиль</button>
                            <button class="btn btn-secondary" data-catalog-action="addPart">
                                <i class="fas fa-plus"></i> Добавить деталь</button>`:''}
                        <button class="btn btn-secondary" data-catalog-action="ho">
                            <i class="fas fa-exchange-alt"></i> Хоз. операции</button>
                    </div>
                </div>
                <aside class="role-card">
                    <h2>${isAdmin?'Основной сценарий администратора':'Основной сценарий пользователя'}</h2>
                    <p>${isAdmin
                        ?'1. Настроить справочники и параметры. 2. Создать детали и собрать автомобиль. 3. Оформить хозяйственную операцию.'
                        :'1. Найти автомобиль. 2. Посмотреть состав и характеристики. 3. Проверить операции и документы по изделию.'}</p>
                </aside>
            </section>
            <section class="metric-grid">${_buildCatalogMetrics()}</section>
            <section class="quick-grid">${_buildQuickActions(isAdmin)}</section>
            <section class="catalog-panel">
                <div class="catalog-panel-head">
                    <div>
                        <div class="catalog-panel-title">Содержимое справочника</div>
                        <div class="catalog-panel-sub">Единая структура таблиц, фильтров и действий помогает быстро ориентироваться в системе.</div>
                    </div>
                    ${isAdmin?'<span class="badge badge-admin"><i class="fas fa-pen"></i> Редактирование включено</span>':'<span class="badge badge-user"><i class="fas fa-eye"></i> Только просмотр</span>'}
                </div>
                <div class="search-row">
                    <input type="text" id="searchInput" class="input-dark" style="max-width:360px"
                        placeholder="Поиск по ID или названию…" value="${escapeHtml(searchText)}">
                    <button class="btn btn-secondary" data-catalog-action="resetFilters">
                        <i class="fas fa-filter-circle-xmark"></i> Сбросить</button>
                </div>
                <div class="filter-chips" id="filterChips">${_buildChipsHtml()}</div>
                <div class="table-wrap">
                    <table>
                        <thead><tr>
                            <th>ID</th><th>Название</th><th>Категория</th><th>Информация</th>
                            ${isAdmin?'<th>Действия</th>':''}
                        </tr></thead>
                        <tbody>${loading
                            ?`<tr><td colspan="${isAdmin?5:4}" style="text-align:center;padding:2rem;color:var(--muted)">
                               <i class="fas fa-spinner fa-spin"></i> Загрузка…</td></tr>`
                            :_buildRowsHtml(isAdmin)}</tbody>
                    </table>
                </div>
            </section>
        </main>
    </div>`;
    app.querySelector('#searchInput').addEventListener('input',e=>{
        searchText=e.target.value;renderCatalog(true);});
    app.querySelector('#filterChips').addEventListener('click',e=>{
        const c=e.target.closest('.filter-chip');if(!c)return;
        categoryFilter=c.dataset.cat;renderCatalog(true);});
    app.querySelector('#logoutBtn').addEventListener('click',()=>{
        currentRole=null;saveSession(null);
        allPartsData=[];categoryFilter='all';searchText='';
        showLogin();});
    app.querySelectorAll('[data-catalog-action]').forEach(el=>el.addEventListener('click',()=>{
        const action=el.dataset.catalogAction;
        if(action==='createCar'&&isAdmin)openWizard();
        else if(action==='addPart'&&isAdmin)openPartForm(null,null);
        else if(action==='createUser'&&isAdmin)openCreateUserForm();
        else if(action==='enums'&&isAdmin)openEnumsManager();
        else if(action==='params'&&isAdmin)openParametersManager();
        else if(action==='ho')openHoOperationsView();
        else if(action==='guide')openInterfaceGuide();
        else if(action==='focusSearch')app.querySelector('#searchInput')?.focus();
        else if(action==='resetFilters'){
            categoryFilter='all';searchText='';
            renderCatalog(false,false);
        }
    }));
    app.querySelectorAll('[data-site-page]').forEach(el=>
        el.addEventListener('click',()=>openSiteInfoPage(el.dataset.sitePage)));
    if(isAdmin){
        _attachActions(app,isAdmin);
    }
}

function _attachActions(app,isAdmin){
    if(!isAdmin)return;
    app.querySelectorAll('[data-action="edit"]').forEach(el=>
        el.addEventListener('click',()=>openPartForm(el.dataset.id,el.dataset.cat)));
    app.querySelectorAll('[data-action="delete"]').forEach(el=>
        el.addEventListener('click',()=>confirmDelete(el.dataset.id,el.dataset.cat,el.dataset.name)));
    // ПР3: кнопка параметров автомобиля
    app.querySelectorAll('[data-action="params"]').forEach(el=>
        el.addEventListener('click',()=>openEmobileParamsForm(el.dataset.id, el.dataset.name)));
}

function openInterfaceGuide(){
    createModal(`
        <div class="modal-title">
            <i class="fas fa-route" style="color:var(--accent)"></i>
            Карта интерфейса и сценарии
        </div>
        <div class="guide-grid">
            <div class="guide-card">
                <h3>Информационная архитектура</h3>
                <p>Главный экран разделен на каталог изделий, параметры, справочники и хозяйственные операции. Пользователь видит просмотр, администратор дополнительно получает настройку и редактирование.</p>
            </div>
            <div class="guide-card">
                <h3>Task flow администратора</h3>
                <ol style="margin-left:1rem">
                    <li>Создать справочники ПР2.</li>
                    <li>Настроить параметры ПР3.</li>
                    <li>Собрать автомобиль и оформить ХО.</li>
                </ol>
            </div>
            <div class="guide-card">
                <h3>Task flow пользователя</h3>
                <ol style="margin-left:1rem">
                    <li>Найти изделие через поиск или фильтр.</li>
                    <li>Посмотреть состав и параметры.</li>
                    <li>Открыть связанные операции.</li>
                </ol>
            </div>
        </div>
        <div style="margin-top:1rem;color:var(--muted);font-size:.85rem">
            Интерфейс использует единые карточки, крупные кнопки, контрастные состояния наведения, понятные ошибки и всплывающие уведомления.
        </div>
        <div class="modal-footer">
            <button class="btn btn-primary" onclick="this.closest('.modal').remove()">Понятно</button>
        </div>`,true);
}

function openSiteInfoPage(page){
    const pages={
        about:{
            icon:'fa-circle-info',
            title:'О компании',
            lead:'Emobile Catalog — учебная информационная система для ведения справочника электромобилей, комплектующих и хозяйственных операций.',
            cards:[
                ['Назначение системы','Компания ведет каталог изделий: батареи, двигатели, зарядные системы, шасси, кузова, электронику и готовые электромобили. Данные используются для подбора комплектации и оформления операций.'],
                ['Роли пользователей','Пользователь просматривает каталог, состав изделий и операции. Администратор управляет справочниками, параметрами, пользователями и настройкой хозяйственных операций.'],
                ['Связь с курсовой','Раздел показывает предметную область, целевую аудиторию и назначение приложения, что требуется в части проектирования пользовательского интерфейса.'],
                ['Преимущества','Единый справочник, фильтрация, параметрическое описание автомобилей, история операций и понятный интерфейс для разных ролей.'],
            ],
        },
        services:{
            icon:'fa-screwdriver-wrench',
            title:'Услуги',
            lead:'Раздел описывает функции, которые система предоставляет клиентам и сотрудникам компании.',
            cards:[
                ['Подбор электромобиля','Поиск и фильтрация моделей по названию, категории, составу узлов и параметрам автомобиля.'],
                ['Конфигурирование изделия','Администратор может собрать автомобиль из готовых узлов: силовой установки, батареи, зарядной системы, шасси, кузова и электроники.'],
                ['Ведение справочников','Поддержка перечислений ПР2: типы батарей, стандарты зарядки, типы приводов и другие допустимые значения.'],
                ['Оформление операций','ПР4 позволяет учитывать отгрузку, приемку, документы, участников, позиции и сумму операции.'],
            ],
        },
        delivery:{
            icon:'fa-truck-fast',
            title:'Доставка и цена',
            lead:'Стоимость в системе формируется из цены позиций хозяйственной операции, а доставка фиксируется как параметр операции или документ.',
            cards:[
                ['Цена автомобиля','В хозяйственной операции каждая позиция содержит изделие, количество и цену за единицу. Итоговая сумма рассчитывается по позициям.'],
                ['Способ доставки','Для операций можно завести параметр ПР3/ПР4 “Способ доставки” и связать его со справочником значений: самовывоз, автотранспорт, Ж/Д, авиа.'],
                ['Документы','К операции прикрепляются типы документов: счет, накладная, акт приема-передачи. Это показывает полный путь от заказа до оформления.'],
                ['Статусы','Операция проходит состояния “Черновик”, “Подтверждена”, “Отменена”, чтобы пользователь видел актуальность цены и доставки.'],
            ],
        },
        contacts:{
            icon:'fa-address-book',
            title:'Контакты',
            lead:'Контактный раздел нужен пользователю как точка связи с компанией и администратором справочника.',
            custom:`
                <div class="info-section">
                    <div class="info-card">
                        <h3>Учебная компания Emobile</h3>
                        <div class="contact-list">
                            <div class="contact-item"><i class="fas fa-location-dot"></i><span>Санкт-Петербург, учебный демонстрационный офис</span></div>
                            <div class="contact-item"><i class="fas fa-phone"></i><span>+7 (812) 000-00-00</span></div>
                            <div class="contact-item"><i class="fas fa-envelope"></i><span>support@emobile.local</span></div>
                            <div class="contact-item"><i class="fas fa-clock"></i><span>Пн-Пт, 10:00-18:00</span></div>
                        </div>
                    </div>
                    <div class="info-card">
                        <h3>К кому обращаться</h3>
                        <p>Пользователь обращается за консультацией по каталогу и операциям. Администратор отвечает за наполнение справочников, параметры изделий и учетные записи.</p>
                    </div>
                </div>`,
        },
    };
    const data=pages[page]||pages.about;
    const body=data.custom||`<div class="info-section">${
        data.cards.map(([title,text])=>`
            <div class="info-card">
                <h3>${escapeHtml(title)}</h3>
                <p>${escapeHtml(text)}</p>
            </div>`).join('')
    }</div>`;
    createModal(`
        <div class="modal-title">
            <i class="fas ${data.icon}" style="color:var(--accent)"></i>
            ${escapeHtml(data.title)}
        </div>
        <p style="color:var(--muted);margin-bottom:1rem">${escapeHtml(data.lead)}</p>
        ${body}
        <div class="modal-footer">
            <button class="btn btn-primary" onclick="this.closest('.modal').remove()">Закрыть</button>
        </div>`,true);
}

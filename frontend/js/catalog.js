/* ===================================================
   catalog.js — загрузка данных и рендер каталога
   =================================================== */

let allPartsData   = [];
let categoryFilter = 'all';
let searchText     = '';

// ---------- Загрузка ----------
async function loadCatalogData() {
    const flat = [];
    for (const cat of CATEGORY_MAP) {
        try {
            const items = await api[cat.key].list();
            if (!Array.isArray(items)) continue;
            for (const item of items) {
                const name = cat.nameField
                    ? (item[cat.nameField] || `${cat.label} #${shortId(item[cat.idField])}`)
                    : `${cat.label} …${shortId(item[cat.idField])}`;
                flat.push({
                    id:          item[cat.idField],
                    name,
                    category:    cat.label,
                    categoryKey: cat.key,
                    info:        cat.infoFn(item),
                    rawData:     item,
                });
            }
        } catch (e) {
            console.warn(`[catalog] Ошибка загрузки «${cat.label}»:`, e.message);
        }
    }
    allPartsData = flat;
    renderCatalog();
}

// ---------- Рендер ----------
function renderCatalog() {
    const appDiv = document.getElementById('app');
    if (!appDiv) return;

    const filtered = allPartsData.filter(p => {
        const matchCat    = categoryFilter === 'all' || p.category === categoryFilter;
        const q           = searchText.toLowerCase();
        const matchSearch = !q || p.name.toLowerCase().includes(q) || p.id.toLowerCase().includes(q);
        return matchCat && matchSearch;
    });

    // Чипы фильтра
    const usedLabels  = [...new Set(allPartsData.map(p => p.category))];
    const filterChips = ['all', ...usedLabels].map(label => {
        const active = categoryFilter === label ? 'active' : '';
        const display = label === 'all' ? 'Все' : label;
        return `<span class="filter-chip ${active}" data-cat="${escapeHtml(label)}">${escapeHtml(display)}</span>`;
    }).join('');

    // Строки таблицы
    const adminCols = currentRole === 'admin';
    const rows = filtered.length
        ? filtered.map(p => `
            <tr>
                <td><span class="text-mono">${escapeHtml(shortId(p.id))}</span></td>
                <td style="font-weight:500">${escapeHtml(p.name)}</td>
                <td><span class="badge badge-default">${escapeHtml(p.category)}</span></td>
                <td class="text-muted">${escapeHtml(p.info)}</td>
                ${adminCols ? `<td style="white-space:nowrap">
                    <i class="fas fa-edit action-icon" style="color:#60a5fa"
                       data-action="edit" data-id="${escapeHtml(p.id)}" data-cat="${escapeHtml(p.categoryKey)}"
                       title="Редактировать"></i>
                    <i class="fas fa-trash-alt action-icon" style="color:#f87171"
                       data-action="delete" data-id="${escapeHtml(p.id)}" data-cat="${escapeHtml(p.categoryKey)}"
                       data-name="${escapeHtml(p.name)}" title="Удалить"></i>
                </td>` : ''}
            </tr>`).join('')
        : `<tr><td colspan="${adminCols ? 5 : 4}">
               <div class="empty-state"><i class="fas fa-inbox"></i>Нет данных</div>
           </td></tr>`;

    const roleBadge = currentRole === 'admin'
        ? `<span class="badge badge-admin"><i class="fas fa-shield-alt"></i> Администратор</span>`
        : `<span class="badge badge-user"><i class="fas fa-user"></i> Пользователь</span>`;

    const adminToolbar = adminCols ? `
        <div class="toolbar">
            <button class="btn btn-primary" id="createCarBtn">
                <i class="fas fa-car"></i> Создать автомобиль
            </button>
            <button class="btn btn-secondary" id="addPartBtn">
                <i class="fas fa-plus"></i> Добавить деталь
            </button>
        </div>` : '';

    appDiv.innerHTML = `
        <div class="container">
            <div class="page-header">
                <h1 class="page-title"><i class="fas fa-charging-station"></i> Каталог Emobile</h1>
                <div style="display:flex;align-items:center;gap:0.75rem">
                    ${roleBadge}
                    <button class="btn btn-secondary" id="logoutBtn" title="Выйти">
                        <i class="fas fa-sign-out-alt"></i>
                    </button>
                </div>
            </div>

            ${adminToolbar}

            <div class="toolbar">
                <input type="text" id="searchInput" class="input-dark"
                       style="max-width:320px" placeholder="Поиск по ID или названию…"
                       value="${escapeHtml(searchText)}">
            </div>

            <div class="filter-chips" id="filterChips">${filterChips}</div>

            <div class="table-wrap">
                <table>
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>Название</th>
                            <th>Категория</th>
                            <th>Информация</th>
                            ${adminCols ? '<th>Действия</th>' : ''}
                        </tr>
                    </thead>
                    <tbody>${rows}</tbody>
                </table>
            </div>
        </div>`;

    // ---- события ----
    document.getElementById('searchInput').addEventListener('input', e => {
        searchText = e.target.value;
        renderCatalog();
    });

    document.getElementById('filterChips').addEventListener('click', e => {
        const chip = e.target.closest('.filter-chip');
        if (!chip) return;
        categoryFilter = chip.dataset.cat;
        renderCatalog();
    });

    document.getElementById('logoutBtn').addEventListener('click', () => {
        currentRole = null;
        allPartsData = [];
        showLogin();
    });

    if (adminCols) {
        document.getElementById('createCarBtn').addEventListener('click', openWizard);
        document.getElementById('addPartBtn').addEventListener('click', () => openPartForm(null, null));

        document.querySelectorAll('[data-action="edit"]').forEach(el => {
            el.addEventListener('click', () => openPartForm(el.dataset.id, el.dataset.cat));
        });
        document.querySelectorAll('[data-action="delete"]').forEach(el => {
            el.addEventListener('click', () => confirmDelete(el.dataset.id, el.dataset.cat, el.dataset.name));
        });
    }
}

/* ===================================================
   api.js — HTTP-клиент, адаптированный под handler.go
   Паттерн маршрутов (Gin):
     GET    /api/{group}/list
     GET    /api/{group}/get{Word}{id}
     POST   /api/{group}/create
     PUT    /api/{group}/{updatePrefix}{id}
     DELETE /api/{group}/delete{id}
   =================================================== */

async function apiRequest(url, options = {}) {
    const res = await fetch(url, {
        ...options,
        headers: { 'Content-Type': 'application/json', ...options.headers },
    });
    if (!res.ok) {
        let msg = `HTTP ${res.status}`;
        try { const e = await res.json(); msg = e.message || e.error || msg; } catch (_) {}
        throw new Error(msg);
    }
    if (res.status === 204) return null;
    return res.json();
}

/**
 * Создаёт набор методов CRUD для ресурса на основе конфига из CATEGORY_MAP.
 * @param {string} group         — префикс группы в URL (/api/{group}/...)
 * @param {string} getWord       — слово в GET /:id маршруте  (/get{Word}{id})
 * @param {string} [updatePrefix] — префикс PUT-маршрута (по умолчанию «update»)
 */
function makeResource(group, getWord, updatePrefix = 'update') {
    const base = `${API_BASE}/${group}`;
    return {
        list:    ()          => apiRequest(`${base}/list`),
        getById: (id)        => apiRequest(`${base}/get${getWord}${id}`),
        create:  (data)      => apiRequest(`${base}/create`,
                                    { method: 'POST', body: JSON.stringify(data) }),
        update:  (id, data)  => apiRequest(`${base}/${updatePrefix}${id}`,
                                    { method: 'PUT',  body: JSON.stringify(data) }),
        delete:  (id)        => apiRequest(`${base}/delete${id}`,
                                    { method: 'DELETE' }),
    };
}

// Строим api-объект из CATEGORY_MAP — никаких дублирующих объявлений
const api = { auth: {} };

// Заполняем динамически после загрузки config.js (CATEGORY_MAP уже доступен)
// Это выполнится сразу, т.к. config.js подключён раньше
CATEGORY_MAP.forEach(cat => {
    api[cat.key] = makeResource(
        cat.group,
        cat.getWord,
        cat.updatePrefix   // undefined → функция подставит 'update'
    );
});

// Auth — отдельный маршрут: POST /api/auth/login
api.auth.login = (username, password) =>
    apiRequest(`${API_BASE}/auth/login`, {
        method: 'POST',
        body: JSON.stringify({ username, password }),
    });

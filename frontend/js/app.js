/* ===================================================
   app.js — точка входа: инициализация и страница логина
   =================================================== */

let currentRole = null;   // 'admin' | 'user' | null

function showLogin() {
    const appDiv = document.getElementById('app');
    appDiv.innerHTML = `
        <div class="login-wrap">
            <div class="login-card">
                <h2><i class="fas fa-bolt"></i> MISPRIS</h2>
                <p class="text-muted" style="text-align:center;margin-bottom:1.5rem;font-size:0.9rem">
                    Управление сборкой электромобилей
                </p>
                <form id="loginForm">
                    <div class="form-group">
                        <label>Логин</label>
                        <input type="text" id="loginUsername" class="input-dark"
                               autocomplete="username" required>
                    </div>
                    <div class="form-group">
                        <label>Пароль</label>
                        <input type="password" id="loginPassword" class="input-dark"
                               autocomplete="current-password" required>
                    </div>
                    <button type="submit" class="btn btn-primary" style="width:100%;justify-content:center;padding:0.65rem">
                        <i class="fas fa-sign-in-alt"></i> Войти
                    </button>
                    <div id="loginError" class="text-danger" style="text-align:center;margin-top:0.75rem;font-size:0.88rem"></div>
                </form>
            </div>
        </div>
    `;

    document.getElementById('loginForm').addEventListener('submit', async e => {
        e.preventDefault();
        const username = document.getElementById('loginUsername').value.trim();
        const password = document.getElementById('loginPassword').value;
        const errEl    = document.getElementById('loginError');
        errEl.textContent = '';
        try {
            const data = await api.auth.login(username, password);
            currentRole = data.role;
            categoryFilter = 'all';
            searchText     = '';
            await loadCatalogData();
        } catch (err) {
            errEl.textContent = 'Неверный логин или пароль';
        }
    });
}

// Инициализация
document.addEventListener('DOMContentLoaded', showLogin);

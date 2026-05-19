/* auth.js */
let currentRole=null;
function saveSession(role){try{sessionStorage.setItem('role',role||'');}catch(_){}currentRole=role;}
function loadSession(){try{currentRole=sessionStorage.getItem('role')||null;}catch(_){}}

function showLogin(){
    const app=document.getElementById('app');
    app.innerHTML=`
    <div class="login-wrap"><div class="login-card">
        <div class="login-logo"><i class="fas fa-charging-station"></i> Emobile</div>
        <div class="login-sub">Каталог электромобилей</div>
        <div class="form-group"><label>Логин</label>
            <input id="lU" type="text" class="input-dark" placeholder="username" autocomplete="username"></div>
        <div class="form-group"><label>Пароль</label>
            <input id="lP" type="password" class="input-dark" placeholder="••••••••" autocomplete="current-password"></div>
        <div id="lE" class="error-banner hidden"><i class="fas fa-exclamation-circle"></i><span></span></div>
        <div style="margin-top:1rem">
            <button class="btn btn-primary" id="lBtn" style="width:100%">
                <i class="fas fa-sign-in-alt"></i> Войти</button></div>
    </div></div>`;
    const uEl=app.querySelector('#lU'),pEl=app.querySelector('#lP'),
          eEl=app.querySelector('#lE'),btn=app.querySelector('#lBtn');
    function setErr(m){eEl.querySelector('span').textContent=m;eEl.classList.toggle('hidden',!m);}
    async function doLogin(){
        const u=uEl.value.trim(),p=pEl.value.trim();
        if(!u||!p){setErr('Заполните логин и пароль');return;}
        btn.disabled=true;btn.innerHTML='<i class="fas fa-spinner fa-spin"></i> Вход…';setErr('');
        try{
            const res=await api.auth.login(u,p);
            saveSession(res?.role||res?.data?.role||'user');
            await loadCatalogData();
        }catch(e){setErr(e.message);}
        finally{btn.disabled=false;btn.innerHTML='<i class="fas fa-sign-in-alt"></i> Войти';}
    }
    btn.addEventListener('click',doLogin);
    pEl.addEventListener('keydown',e=>{if(e.key==='Enter')doLogin();});
    uEl.focus();
}

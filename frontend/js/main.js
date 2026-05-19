/* main.js */
(async()=>{
    loadSession();
    if(currentRole){await loadCatalogData();}
    else{showLogin();}
})();

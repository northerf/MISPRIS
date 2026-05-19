/* api.js — полный файл с добавлением ПР2, ПР3 и ПР4 */
async function apiRequest(url,opts={}){
    let res;
    try{
        res=await fetch(url,{...opts,headers:{'Content-Type':'application/json',...opts.headers}});
    }catch(e){throw new Error(friendlyError(e.message||'NetworkError'));}
    if(res.status===204)return null;
    let body=null;
    try{body=await res.json();}catch(_){}
    if(!res.ok){
        const raw=body?.message||body?.error||body?.detail||`HTTP ${res.status}`;
        throw new Error(friendlyError(raw));
    }
    return body;
}
function makeResource(group,getWord){
    const b=`${API_BASE}/${group}`;
    return{
        list:    ()     =>apiRequest(`${b}/list`),
        getById: id     =>apiRequest(`${b}/get${getWord}${id}`),
        create:  data   =>apiRequest(`${b}/create`,{method:'POST',body:JSON.stringify(data)}),
        update:  (id,d) =>apiRequest(`${b}/update${id}`,{method:'PUT',body:JSON.stringify(d)}),
        delete:  id     =>apiRequest(`${b}/delete${id}`,{method:'DELETE'}),
    };
}
const api={auth:{}};
CATEGORY_MAP.forEach(cat=>{api[cat.key]=makeResource(cat.group,cat.getWord);});
api.auth.login=(username,password)=>
    apiRequest(`${API_BASE}/auth/login`,{method:'POST',body:JSON.stringify({username,password})});
api.users={
    list:  ()   =>apiRequest(`${API_BASE}/users/list`),
    create:data =>apiRequest(`${API_BASE}/users/create`,{method:'POST',body:JSON.stringify(data)}),
    delete:id   =>apiRequest(`${API_BASE}/users/delete${id}`,{method:'DELETE'}),
};

/* ── ПР2: Перечисления ──────────────────────────────────────── */
api.enumClass={
    list:    ()      =>apiRequest(`${API_BASE}/enum-class/list`),
    getById: id      =>apiRequest(`${API_BASE}/enum-class/getEnumClass${id}`),
    create:  data    =>apiRequest(`${API_BASE}/enum-class/create`,{method:'POST',body:JSON.stringify(data)}),
    update:  (id,d)  =>apiRequest(`${API_BASE}/enum-class/update${id}`,{method:'PUT',body:JSON.stringify(d)}),
    delete:  id      =>apiRequest(`${API_BASE}/enum-class/delete${id}`,{method:'DELETE'}),
    // Получить значения в порядке order_num (SQL-функция get_enum_values)
    getValues:  id   =>apiRequest(`${API_BASE}/enum-class/values${id}`),
    // Проверить допустимость значения (SQL-функция validate_enum_value)
    validate: data   =>apiRequest(`${API_BASE}/enum-class/validate`,{method:'POST',body:JSON.stringify(data)}),
};
api.enumPosition={
    list:    ()      =>apiRequest(`${API_BASE}/enum-position/list`),
    getById: id      =>apiRequest(`${API_BASE}/enum-position/getEnumPosition${id}`),
    create:  data    =>apiRequest(`${API_BASE}/enum-position/create`,{method:'POST',body:JSON.stringify(data)}),
    update:  (id,d)  =>apiRequest(`${API_BASE}/enum-position/update${id}`,{method:'PUT',body:JSON.stringify(d)}),
    delete:  id      =>apiRequest(`${API_BASE}/enum-position/delete${id}`,{method:'DELETE'}),
    // Изменить порядок позиции
    reorder: (id,d)  =>apiRequest(`${API_BASE}/enum-position/reorder${id}`,{method:'POST',body:JSON.stringify(d)}),
};

/* ── ПР3: Параметры ─────────────────────────────────────────── */
api.parameter={
    list:    ()      =>apiRequest(`${API_BASE}/parameter/list`),
    getById: id      =>apiRequest(`${API_BASE}/parameter/getParameter${id}`),
    create:  data    =>apiRequest(`${API_BASE}/parameter/create`,{method:'POST',body:JSON.stringify(data)}),
    update:  (id,d)  =>apiRequest(`${API_BASE}/parameter/update${id}`,{method:'PUT',body:JSON.stringify(d)}),
    delete:  id      =>apiRequest(`${API_BASE}/parameter/delete${id}`,{method:'DELETE'}),
};
api.componentParameter={
    list:    ()      =>apiRequest(`${API_BASE}/component-parameter/list`),
    getById: id      =>apiRequest(`${API_BASE}/component-parameter/getComponentParameter${id}`),
    create:  data    =>apiRequest(`${API_BASE}/component-parameter/create`,{method:'POST',body:JSON.stringify(data)}),
    update:  (id,d)  =>apiRequest(`${API_BASE}/component-parameter/update${id}`,{method:'PUT',body:JSON.stringify(d)}),
    delete:  id      =>apiRequest(`${API_BASE}/component-parameter/delete${id}`,{method:'DELETE'}),
    // Получить параметры по типу компонента (SQL-функция get_component_parameters)
    byType:  type    =>apiRequest(`${API_BASE}/component-parameter/byType${type}`),
    // Скопировать параметры от одного типа к другому (SQL-процедура copy_component_parameters)
    copyFromType: d  =>apiRequest(`${API_BASE}/component-parameter/copyFromType`,{method:'POST',body:JSON.stringify(d)}),
};
api.emobileParam={
    list:    ()      =>apiRequest(`${API_BASE}/emobile-parameter/list`),
    getById: id      =>apiRequest(`${API_BASE}/emobile-parameter/getEmobileParameter${id}`),
    create:  data    =>apiRequest(`${API_BASE}/emobile-parameter/create`,{method:'POST',body:JSON.stringify(data)}),
    update:  (id,d)  =>apiRequest(`${API_BASE}/emobile-parameter/update${id}`,{method:'PUT',body:JSON.stringify(d)}),
    delete:  id      =>apiRequest(`${API_BASE}/emobile-parameter/delete${id}`,{method:'DELETE'}),
    // Все параметры конкретного автомобиля
    byEmobile: id    =>apiRequest(`${API_BASE}/emobile-parameter/byEmobile${id}`),
};

/* ── ПР4: Хозяйственные операции ────────────────────────── */
api.shd={
    list:    ()     =>apiRequest(`${API_BASE}/shd/list`),
    getById: id     =>apiRequest(`${API_BASE}/shd/getShd${id}`),
    create:  data   =>apiRequest(`${API_BASE}/shd/create`,{method:'POST',body:JSON.stringify(data)}),
    update:  (id,d) =>apiRequest(`${API_BASE}/shd/update${id}`,{method:'PUT',body:JSON.stringify(d)}),
    delete:  id     =>apiRequest(`${API_BASE}/shd/delete${id}`,{method:'DELETE'}),
};
api.hoClass={
    list:        ()   =>apiRequest(`${API_BASE}/ho-class/list`),
    getById:     id   =>apiRequest(`${API_BASE}/ho-class/getHoClass${id}`),
    create:      d    =>apiRequest(`${API_BASE}/ho-class/create`,{method:'POST',body:JSON.stringify(d)}),
    update:     (id,d)=>apiRequest(`${API_BASE}/ho-class/update${id}`,{method:'PUT',body:JSON.stringify(d)}),
    delete:      id   =>apiRequest(`${API_BASE}/ho-class/delete${id}`,{method:'DELETE'}),
    terminal:    ()   =>apiRequest(`${API_BASE}/ho-class/terminal`),
    children:    id   =>apiRequest(`${API_BASE}/ho-class/children${id}`),
    parameters:  id   =>apiRequest(`${API_BASE}/ho-class/parameters${id}`),
};
api.hoRole={
    list:    ()     =>apiRequest(`${API_BASE}/ho-role/list`),
    getById: id     =>apiRequest(`${API_BASE}/ho-role/getHoRole${id}`),
    create:  data   =>apiRequest(`${API_BASE}/ho-role/create`,{method:'POST',body:JSON.stringify(data)}),
    update:  (id,d) =>apiRequest(`${API_BASE}/ho-role/update${id}`,{method:'PUT',body:JSON.stringify(d)}),
    delete:  id     =>apiRequest(`${API_BASE}/ho-role/delete${id}`,{method:'DELETE'}),
};
api.hoClassRole={
    list:   (hoClassId)=>apiRequest(`${API_BASE}/ho-class-role/list${hoClassId?'?ho_class='+hoClassId:''}`),
    create: data       =>apiRequest(`${API_BASE}/ho-class-role/create`,{method:'POST',body:JSON.stringify(data)}),
    delete: id         =>apiRequest(`${API_BASE}/ho-class-role/delete${id}`,{method:'DELETE'}),
};
api.hoClassParam={
    list:         (hoClassId)=>apiRequest(`${API_BASE}/ho-class-parameter/list${hoClassId?'?ho_class='+hoClassId:''}`),
    getById:      id         =>apiRequest(`${API_BASE}/ho-class-parameter/getHoClassParameter${id}`),
    create:       data       =>apiRequest(`${API_BASE}/ho-class-parameter/create`,{method:'POST',body:JSON.stringify(data)}),
    update:       (id,d)     =>apiRequest(`${API_BASE}/ho-class-parameter/update${id}`,{method:'PUT',body:JSON.stringify(d)}),
    delete:       id         =>apiRequest(`${API_BASE}/ho-class-parameter/delete${id}`,{method:'DELETE'}),
    copyFromClass:data       =>apiRequest(`${API_BASE}/ho-class-parameter/copyFromClass`,{method:'POST',body:JSON.stringify(data)}),
};
api.documentClass={
    list:    ()     =>apiRequest(`${API_BASE}/document-class/list`),
    create:  data   =>apiRequest(`${API_BASE}/document-class/create`,{method:'POST',body:JSON.stringify(data)}),
    update:  (id,d) =>apiRequest(`${API_BASE}/document-class/update${id}`,{method:'PUT',body:JSON.stringify(d)}),
    delete:  id     =>apiRequest(`${API_BASE}/document-class/delete${id}`,{method:'DELETE'}),
};
api.hoClassDoc={
    list:   (hoClassId)=>apiRequest(`${API_BASE}/ho-class-document/list${hoClassId?'?ho_class='+hoClassId:''}`),
    create: data       =>apiRequest(`${API_BASE}/ho-class-document/create`,{method:'POST',body:JSON.stringify(data)}),
    delete: id         =>apiRequest(`${API_BASE}/ho-class-document/delete${id}`,{method:'DELETE'}),
};
api.ho={
    list:        (hoClassId)=>apiRequest(`${API_BASE}/ho/list${hoClassId?'?ho_class='+hoClassId:''}`),
    getById:     id         =>apiRequest(`${API_BASE}/ho/getHo${id}`),
    create:      data       =>apiRequest(`${API_BASE}/ho/create`,{method:'POST',body:JSON.stringify(data)}),
    update:      (id,d)     =>apiRequest(`${API_BASE}/ho/update${id}`,{method:'PUT',body:JSON.stringify(d)}),
    delete:      id         =>apiRequest(`${API_BASE}/ho/delete${id}`,{method:'DELETE'}),
    findByClass: id         =>apiRequest(`${API_BASE}/ho/findByClass${id}`),
};
api.hoActor={
    list:   hoId =>apiRequest(`${API_BASE}/ho-actor/list${hoId?'?ho='+hoId:''}`),
    create: data =>apiRequest(`${API_BASE}/ho-actor/create`,{method:'POST',body:JSON.stringify(data)}),
    delete: id   =>apiRequest(`${API_BASE}/ho-actor/delete${id}`,{method:'DELETE'}),
};
api.hoParamVal={
    list:    hoId   =>apiRequest(`${API_BASE}/ho-param-value/list${hoId?'?ho='+hoId:''}`),
    getById: id     =>apiRequest(`${API_BASE}/ho-param-value/getHoParamValue${id}`),
    create:  data   =>apiRequest(`${API_BASE}/ho-param-value/create`,{method:'POST',body:JSON.stringify(data)}),
    update:  (id,d) =>apiRequest(`${API_BASE}/ho-param-value/update${id}`,{method:'PUT',body:JSON.stringify(d)}),
    delete:  id     =>apiRequest(`${API_BASE}/ho-param-value/delete${id}`,{method:'DELETE'}),
};
api.hoDoc={
    list:   hoId =>apiRequest(`${API_BASE}/ho-document/list${hoId?'?ho='+hoId:''}`),
    create: data =>apiRequest(`${API_BASE}/ho-document/create`,{method:'POST',body:JSON.stringify(data)}),
    delete: id   =>apiRequest(`${API_BASE}/ho-document/delete${id}`,{method:'DELETE'}),
};
api.hoPosition={
    list:    hoId   =>apiRequest(`${API_BASE}/ho-position/list${hoId?'?ho='+hoId:''}`),
    create:  data   =>apiRequest(`${API_BASE}/ho-position/create`,{method:'POST',body:JSON.stringify(data)}),
    update:  (id,d) =>apiRequest(`${API_BASE}/ho-position/update${id}`,{method:'PUT',body:JSON.stringify(d)}),
    delete:  id     =>apiRequest(`${API_BASE}/ho-position/delete${id}`,{method:'DELETE'}),
};

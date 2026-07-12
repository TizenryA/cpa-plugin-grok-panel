package main

const htmlPage = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Grok 面板 v1.1.0</title>
<style>
:root{--bg:#ecebe5;--card:#f7f6f0;--ink:#242522;--muted:#6f716b;--line:#454740;--soft:#dddcd4;--soft2:#d1d0c7;--green:#5f735f;--red:#8b5f59;--yellow:#8a7f5d;--blue:#5d6f78;--violet:#6f6678}
*{box-sizing:border-box}
html,body{min-height:100vh}
body{margin:0;background:var(--bg);color:var(--ink);font-family:'SF Mono','Cascadia Code','JetBrains Mono',Consolas,monospace;font-size:0.95vw;line-height:1.45;padding:2vw}
button,input,select{font-family:inherit;color:var(--ink)}
button{background:var(--card);border:0.12vw solid var(--line);padding:0.7vh 0.9vw;min-height:3.8vh;cursor:pointer;font-size:0.78vw;text-transform:uppercase;letter-spacing:0.04vw;border-radius:0}
button:hover:not(:disabled),button.armed{background:var(--line);color:var(--card)}
button:disabled{cursor:not-allowed;opacity:0.46;background:var(--soft)}
input,select{background:var(--card);border:0.12vw solid var(--line);padding:0.65vh 0.7vw;min-height:3.8vh;font-size:0.78vw;border-radius:0}
input:focus,select:focus,button:focus{outline:0.18vw solid var(--blue);outline-offset:0.18vw}
input[type=checkbox]{appearance:none;width:1.1vw;height:1.1vw;min-height:1.1vw;padding:0;margin:0;background:var(--card);vertical-align:middle}
input[type=checkbox]:checked{background:var(--line)}
.shell{width:96vw;margin:0 auto}
.topline{display:flex;justify-content:space-between;align-items:stretch;gap:1vw;margin-bottom:1.2vh}
.brand{display:flex;align-items:center;gap:0.8vw;min-height:4vh}
h1{font-size:1.8vw;line-height:1;margin:0;font-weight:800;letter-spacing:0.08vw}
h2{font-size:1.05vw;margin:0 0 1vh 0;font-weight:800;letter-spacing:0.05vw;border-bottom:0.12vw solid var(--line);padding-bottom:0.8vh;text-transform:uppercase}
.status-dot{width:0.8vw;height:0.8vw;background:var(--green);display:inline-block;border:0.08vw solid var(--line)}
.status-dot.err{background:var(--red)}
.muted{color:var(--muted);font-size:0.74vw}
.top-actions{display:flex;gap:0.7vw;align-items:center;flex-wrap:wrap;justify-content:flex-end}
.feedback{border:0.12vw solid var(--line);background:var(--card);padding:0.8vh 1vw;min-height:4vh;margin-bottom:1vh;color:var(--muted)}
.feedback.ok{border-color:var(--green);color:var(--green)}
.feedback.warn{border-color:var(--yellow);color:var(--yellow)}
.feedback.error{border-color:var(--red);color:var(--red)}
.stats-grid{display:grid;grid-template-columns:23vw 23vw 23vw 23vw;gap:1vw;margin-bottom:1.5vh}
.stat-card,.panel{background:var(--card);border:0.12vw solid var(--line);padding:1.1vh 1vw;border-radius:0}
.stat-card{min-height:13vh}
.stat-label{font-size:0.68vw;color:var(--muted);text-transform:uppercase;letter-spacing:0.06vw}
.stat-value{font-size:1.6vw;line-height:1.1;font-weight:800;margin-top:0.7vh;word-break:break-word}
.stat-sub{font-size:0.68vw;color:var(--muted);margin-top:0.55vh}
.bar-container{width:18vw;height:0.9vh;background:var(--soft2);margin-top:0.7vh;border:0.08vw solid var(--line)}
.bar-container.small{width:11vw;height:0.85vh;margin-top:0.35vh}
.bar-fill{height:0.7vh;background:var(--green)}
.bar-fill.warn{background:var(--yellow)}
.bar-fill.danger{background:var(--red)}
.panel{margin-bottom:1.5vh}
.form-grid{display:grid;grid-template-columns:18vw 18vw 18vw 18vw 18vw;gap:0.8vw;align-items:end}
.field{display:flex;flex-direction:column;gap:0.45vh}
.field label,.checkline{font-size:0.72vw;color:var(--muted)}
.number-input{width:18vw}
.checkline{display:flex;align-items:center;gap:0.45vw;min-height:3.8vh;border:0.08vw solid var(--soft2);padding:0.6vh 0.6vw;background:var(--bg)}
.help-text{font-size:0.68vw;color:var(--muted);margin-top:0.9vh}
.chart-row{display:flex;align-items:flex-end;gap:0.25vw;height:16vh;padding:1vh 0;overflow-x:auto;overflow-y:hidden;background:var(--bg);border:0.08vw solid var(--soft2)}
.chart-bar{flex:0 0 1.3vw;min-height:0.35vh;background:var(--green);opacity:0.72;position:relative;border:0.06vw solid var(--line)}
.chart-bar.fail{background:var(--red)}
.chart-bar:hover{opacity:1}
.chart-bar:hover::after{content:attr(data-tip);position:absolute;bottom:15vh;left:0;background:var(--line);color:var(--card);padding:0.45vh 0.55vw;font-size:0.64vw;white-space:normal;z-index:8;width:14vw;border:0.08vw solid var(--card)}
.chart-empty{color:var(--muted);padding:5vh 1vw;font-size:0.8vw}
.filter-grid{display:grid;grid-template-columns:23vw 13vw 13vw 13vw 13vw 17vw;gap:0.8vw;margin-bottom:0.9vh;align-items:end}
.search-box{width:23vw}
.select-filter{width:13vw}
.sort-filter{width:17vw}
.batchbar{display:flex;align-items:center;gap:0.7vw;flex-wrap:wrap;margin-bottom:0.9vh;background:var(--bg);border:0.08vw solid var(--soft2);padding:0.8vh 0.8vw}
.batchbar .spacer{flex:1}
.table-wrap{max-height:52vh;overflow:auto;border:0.12vw solid var(--line);background:var(--card)}
table{border-collapse:collapse;width:142vw;min-width:142vw;background:var(--card)}
th,td{text-align:left;padding:0.7vh 0.6vw;border-bottom:0.08vw solid var(--soft2);font-size:0.75vw;vertical-align:top}
th{position:sticky;top:0;z-index:4;background:var(--line);color:var(--card);font-weight:800;text-transform:uppercase;font-size:0.65vw;letter-spacing:0.04vw}
tr:hover{background:var(--bg)}
.email-cell{width:27vw;word-break:break-all}
.actions-cell{display:flex;gap:0.45vw;flex-wrap:wrap}
.tag{display:inline-block;padding:0.25vh 0.4vw;border:0.08vw solid var(--line);font-size:0.64vw;font-weight:800;letter-spacing:0.03vw;text-transform:uppercase;min-width:4.8vw;text-align:center;background:var(--soft)}
.tag.active,.tag.healthy,.tag.standard{background:var(--green);color:var(--card)}
.tag.disabled,.tag.invalid{background:var(--red);color:var(--card)}
.tag.warn,.tag.heavy{background:var(--yellow);color:var(--card)}
.tag.unknown{background:var(--soft2);color:var(--ink)}
.tag.super{background:var(--violet);color:var(--card)}
.tag.info{background:var(--blue);color:var(--card)}
.cell-sub{font-size:0.62vw;color:var(--muted);margin-top:0.35vh;word-break:break-word}
.red-text{color:var(--red)}
.summary-row{display:flex;gap:1.2vw;flex-wrap:wrap;font-size:0.72vw;color:var(--muted);margin-top:0.9vh}
.summary-row b{color:var(--ink)}
.row-invalid{background:#f0e6e2}
.row-disabled{background:#ebe4e0}
.row-warn{background:#eee9dc}
@media (orientation:portrait),(hover:none) and (pointer:coarse){body{font-size:3.2vw;padding:3vw}.shell{width:94vw}.topline{flex-direction:column}.brand{align-items:flex-start;flex-direction:column;gap:1vh}.top-actions{justify-content:stretch;flex-direction:column;align-items:stretch}h1{font-size:6vw}h2{font-size:3.6vw}.muted{font-size:2.8vw}.status-dot{width:2.8vw;height:2.8vw}button,input,select{font-size:2.8vw;min-height:5.2vh;padding:0.9vh 2vw}input[type=checkbox]{width:4vw;height:4vw;min-height:4vw}.feedback{font-size:2.8vw;padding:1vh 2vw}.stats-grid{grid-template-columns:94vw}.stat-card{min-height:12vh}.stat-label,.stat-sub,.field label,.checkline,.help-text,.summary-row{font-size:2.6vw}.stat-value{font-size:5.2vw}.form-grid,.filter-grid{grid-template-columns:94vw}.number-input,.search-box,.select-filter,.sort-filter{width:94vw}.checkline{gap:2vw;padding:1vh 2vw}.top-actions button,.batchbar button{width:94vw}.batchbar{flex-direction:column;align-items:stretch;gap:1vh}.bar-container{width:80vw;height:1.2vh}.bar-container.small{width:24vw}.bar-fill{height:0.95vh}.chart-row{height:18vh}.chart-bar{flex-basis:3vw}.chart-bar:hover::after{font-size:2.4vw;width:42vw;bottom:16vh}.table-wrap{max-height:55vh}table{width:210vw;min-width:210vw}th,td{font-size:2.6vw;padding:0.9vh 1.4vw}th{font-size:2.3vw}.tag{font-size:2.3vw;min-width:13vw;padding:0.4vh 1vw}.cell-sub{font-size:2.2vw}.actions-cell button{font-size:2.3vw;min-height:4.8vh}}
</style>
</head>
<body>
<div class="shell">
<header class="topline">
<div class="brand"><h1>GROK PANEL</h1><span class="status-dot" id="statusDot"></span><span class="muted" id="lastUpdate">等待数据</span></div>
<div class="top-actions"><button id="refreshBtn">刷新数据</button><button id="checkVisibleTopBtn">手动检查可见</button></div>
</header>
<div id="feedback" class="feedback">就绪：所有操作使用同源插件端点，不包含管理密钥。</div>
<section class="stats-grid">
<div class="stat-card"><div class="stat-label">Grok 文件数</div><div class="stat-value" id="statTotal">--</div><div class="stat-sub" id="statTotalSub"></div></div>
<div class="stat-card"><div class="stat-label">活跃账号</div><div class="stat-value" id="statActive">--</div><div class="stat-sub" id="statActiveSub"></div></div>
<div class="stat-card"><div class="stat-label">总请求数</div><div class="stat-value" id="statRequests">--</div><div class="stat-sub" id="statRequestsSub"></div></div>
<div class="stat-card"><div class="stat-label">估算 Token</div><div class="stat-value" id="statTokens">--</div><div class="stat-sub" id="statTokensSub"></div></div>
<div class="stat-card"><div class="stat-label">总容量</div><div class="stat-value" id="statCapacity">--</div><div class="stat-sub" id="statCapacitySub"></div></div>
<div class="stat-card"><div class="stat-label">使用率</div><div class="stat-value" id="statUsage">--</div><div class="stat-sub" id="statUsageSub"></div><div class="bar-container"><div class="bar-fill" id="usageBar"></div></div></div>
<div class="stat-card"><div class="stat-label">账号类型</div><div class="stat-value" id="statTypes">--</div><div class="stat-sub" id="statTypesSub"></div></div>
<div class="stat-card"><div class="stat-label">健康概览</div><div class="stat-value" id="statHealth">--</div><div class="stat-sub" id="statHealthSub"></div></div>
<div class="stat-card"><div class="stat-label">待清理</div><div class="stat-value" id="statInvalid">--</div><div class="stat-sub" id="statInvalidSub"></div></div>
<div class="stat-card"><div class="stat-label">已选择</div><div class="stat-value" id="statSelected">0</div><div class="stat-sub" id="statSelectedSub"></div></div>
</section>
<section class="panel">
<h2>设置</h2>
<div class="form-grid">
<div class="field"><label for="tokenLimit">每账号上限 Token</label><input class="number-input" type="number" id="tokenLimit" value="2000000" min="1"></div>
<div class="field"><label for="tokensPerReq">每请求估算 Token</label><input class="number-input" type="number" id="tokensPerReq" value="5000" min="1"></div>
<div class="field"><label for="failThreshold">失败阈值</label><input class="number-input" type="number" id="failThreshold" value="3" min="1"></div>
<label class="checkline"><input type="checkbox" id="autoCheck"> 自动检查 默认关</label>
<label class="checkline"><input type="checkbox" id="autoDelete"> 自动删除无效 默认关</label>
<label class="checkline"><input type="checkbox" id="protectSuper" checked> 保护 super 默认开</label>
<label class="checkline"><input type="checkbox" id="protectHeavy" checked> 保护 heavy 默认开</label>
<label class="checkline"><input type="checkbox" id="protectUnknown" checked> 保护 unknown 默认开</label>
</div>
<div class="help-text">删除和清理需要再次点击同一按钮确认；若后端端点或 CPA 管理授权不可用，会在这里显示明确错误。</div>
</section>
<section class="panel"><h2>请求趋势</h2><div class="chart-row" id="chartRow"></div></section>
<section class="panel"><h2>账号明细</h2>
<div class="filter-grid">
<div class="field"><label for="searchBox">搜索</label><input type="text" class="search-box" id="searchBox" placeholder="邮箱、状态、类型"></div>
<div class="field"><label for="statusFilter">状态</label><select class="select-filter" id="statusFilter"><option value="all">全部</option><option value="active">活跃</option><option value="disabled">禁用</option><option value="other">其他</option><option value="unknown">未知</option></select></div>
<div class="field"><label for="typeFilter">类型</label><select class="select-filter" id="typeFilter"><option value="all">全部</option><option value="standard">standard</option><option value="super">super</option><option value="heavy">heavy</option><option value="unknown">unknown</option></select></div>
<div class="field"><label for="healthFilter">健康</label><select class="select-filter" id="healthFilter"><option value="all">全部</option><option value="healthy">健康</option><option value="warn">警告</option><option value="invalid">无效</option><option value="disabled">禁用</option><option value="unknown">未知</option></select></div>
<div class="field"><label for="usageFilter">用量</label><select class="select-filter" id="usageFilter"><option value="all">全部</option><option value="unused">未使用</option><option value="low">低于一半</option><option value="warn">一半以上</option><option value="high">高于八成</option></select></div>
<div class="field"><label for="sortFilter">排序</label><select class="sort-filter" id="sortFilter"><option value="success_desc">成功请求降序</option><option value="failed_desc">失败降序</option><option value="usage_desc">用量降序</option><option value="health_asc">健康优先</option><option value="type_asc">类型</option><option value="email_asc">邮箱</option></select></div>
</div>
<div class="batchbar">
<label class="checkline"><input type="checkbox" id="selectVisible"> 选择可见</label>
<button id="checkSelectedBtn">检查选中</button>
<button id="batchDeleteBtn">删除选中</button>
<button id="cleanupInvalidBtn">清理无效</button>
<span class="muted" id="selectionInfo">未选择</span><span class="spacer"></span><span class="muted" id="filterInfo">无过滤</span>
</div>
<div class="table-wrap"><table><thead><tr><th>选</th><th>#</th><th>邮箱</th><th>类型</th><th>健康</th><th>状态</th><th>成功</th><th>失败</th><th>估算 Token</th><th>用量</th><th>操作</th></tr></thead><tbody id="tableBody"></tbody></table></div>
<div class="summary-row"><span>可见 <b id="sumCount">0</b></span><span>已使用 <b id="sumUsed">0</b></span><span>剩余 <b id="sumAvail">0</b></span><span>成功率 <b id="sumRate">0%</b></span><span>无效 <b id="sumInvalid">0</b></span><span>可清理 <b id="sumCleanup">0</b></span></div>
</section>
</div>
<script>
/*
Frontend v1.1.0 same-origin endpoint contract for a matching backend.
No management key is embedded in this page; CPA iframe/session auth must be supplied by the host.
GET  ./data                         -> current stats shape used by v1.0 plus optional account_type, type, health, unavailable fields.
POST ./accounts/check               -> body {emails:[string], threshold:number}; returns {results:[{email, health, detail, account_type}]}.
POST ./accounts/delete              -> body {emails:[string], threshold:number, protect:{super:boolean, heavy:boolean, unknown:boolean}}.
POST ./accounts/cleanup-invalid     -> body {emails:[string], threshold:number, protect:{super:boolean, heavy:boolean, unknown:boolean}}.
If these mutation endpoints are absent, the UI reports that operations are unavailable instead of using any hardcoded key.
*/
var settingsKey='grok-panel-v1.1.0-settings';
var allData=[];
var lastData=null;
var selected={};
var confirmUntil={};
var healthOverrides={};
var typeOverrides={};
var busy=false;
var autoCheckBusy=false;
var lastAutoCheckAt=0;
var refreshHandle=null;
var settings=loadSettings();
function byId(id){return document.getElementById(id)}
function defaults(){return{tokenLimit:2000000,tokensPerReq:5000,threshold:3,autoCheck:false,autoDelete:false,protectSuper:true,protectHeavy:true,protectUnknown:true}}
function loadSettings(){var base=defaults();try{var raw=localStorage.getItem(settingsKey);if(raw){var parsed=JSON.parse(raw);Object.keys(base).forEach(function(k){if(parsed[k]!==undefined)base[k]=parsed[k]})}}catch(e){}return base}
function saveSettings(){try{localStorage.setItem(settingsKey,JSON.stringify(settings))}catch(e){}}
function initSettings(){byId('tokenLimit').value=settings.tokenLimit;byId('tokensPerReq').value=settings.tokensPerReq;byId('failThreshold').value=settings.threshold;byId('autoCheck').checked=!!settings.autoCheck;byId('autoDelete').checked=!!settings.autoDelete;byId('protectSuper').checked=!!settings.protectSuper;byId('protectHeavy').checked=!!settings.protectHeavy;byId('protectUnknown').checked=!!settings.protectUnknown}
function readSettings(evt){settings.tokenLimit=Math.max(1,parseInt(byId('tokenLimit').value,10)||2000000);settings.tokensPerReq=Math.max(1,parseInt(byId('tokensPerReq').value,10)||5000);settings.threshold=Math.max(1,parseInt(byId('failThreshold').value,10)||3);settings.autoCheck=!!byId('autoCheck').checked;settings.autoDelete=!!byId('autoDelete').checked;settings.protectSuper=!!byId('protectSuper').checked;settings.protectHeavy=!!byId('protectHeavy').checked;settings.protectUnknown=!!byId('protectUnknown').checked;saveSettings();if(evt&&evt.target&&evt.target.id==='autoDelete'&&settings.autoDelete)setFeedback('自动删除已开启：只处理未受保护且判定无效的账号。','warn');if(evt&&evt.target&&evt.target.id==='autoCheck'&&settings.autoCheck)setFeedback('自动检查已开启：刷新后会调用同源检查端点。','warn');renderAll()}
function bindEvents(){byId('refreshBtn').addEventListener('click',function(){fetchData(true,false)});byId('checkVisibleTopBtn').addEventListener('click',manualCheckVisible);byId('checkSelectedBtn').addEventListener('click',manualCheckSelected);byId('batchDeleteBtn').addEventListener('click',requestBatchDelete);byId('cleanupInvalidBtn').addEventListener('click',requestCleanupInvalid);byId('selectVisible').addEventListener('change',toggleVisibleSelection);['searchBox','statusFilter','typeFilter','healthFilter','usageFilter','sortFilter'].forEach(function(id){byId(id).addEventListener('input',function(){renderAll()});byId(id).addEventListener('change',function(){renderAll()})});['tokenLimit','tokensPerReq','failThreshold','autoCheck','autoDelete','protectSuper','protectHeavy','protectUnknown'].forEach(function(id){byId(id).addEventListener('input',readSettings);byId(id).addEventListener('change',readSettings)});byId('tableBody').addEventListener('click',handleTableClick);byId('tableBody').addEventListener('change',handleTableChange);window.addEventListener('resize',function(){renderAll()})}
function fmt(n){n=Number(n)||0;if(n>=1000000000)return(n/1000000000).toFixed(2)+'B';if(n>=1000000)return(n/1000000).toFixed(2)+'M';if(n>=1000)return(n/1000).toFixed(1)+'K';return String(n)}
function fmtTime(){var d=new Date();return d.toLocaleTimeString('zh-CN',{hour12:false})+' '+d.toLocaleDateString('zh-CN')}
function esc(v){return String(v===undefined||v===null?'':v).replace(/[&<>"']/g,function(c){return{'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]})}
function apiBase(){return window.location.pathname.replace(/\/+$/,'')}
function fixedApiUrl(endpoint){return apiBase()+'/'+String(endpoint).replace(/^\/+/, '')}
function xorDecode(bytes,key){var out=new Uint8Array(bytes.length);for(var i=0;i<bytes.length;i++)out[i]=bytes[i]^key[i%key.length];return new TextDecoder().decode(out)}
function readCPAConnection(){try{var raw=localStorage.getItem('cli-proxy-auth');if(!raw)return null;if(raw.indexOf('enc::v1::')===0){var bin=atob(raw.slice(9));var bytes=new Uint8Array(bin.length);for(var i=0;i<bin.length;i++)bytes[i]=bin.charCodeAt(i);var salt='cli-proxy-api-webui::secure-storage|'+window.location.host+'|'+navigator.userAgent;raw=xorDecode(bytes,new TextEncoder().encode(salt))}var parsed=JSON.parse(raw);var state=parsed&&parsed.state?parsed.state:parsed;return state&&state.managementKey?{apiBase:String(state.apiBase||'').replace(/\/$/,''),managementKey:String(state.managementKey)}:null}catch(e){return null}}
function managementFetch(path,options){var conn=readCPAConnection();if(!conn)throw new Error('当前管理中心没有保存管理密钥。请退出后勾选“记住密码”重新登录，再打开插件。');options=options||{};options.headers=options.headers||{};options.headers.Authorization='Bearer '+conn.managementKey;if(!options.headers.accept)options.headers.accept='application/json';return fetch((conn.apiBase||window.location.origin)+'/v0/management'+path,options)}
function setFeedback(msg,type){var el=byId('feedback');el.className='feedback '+(type||'');el.textContent=msg}
function setBusy(flag){busy=!!flag;document.body.classList.toggle('busy',busy);updateToolbarState();renderTable()}
function parseJsonText(text,endpoint){try{return text?JSON.parse(text):{}}catch(e){var low=String(text||'').toLowerCase();if(low.indexOf('<!doctype')>=0||low.indexOf('<html')>=0)throw new Error('操作端点 '+endpoint+' 未启用：当前后端返回了面板页面，请升级插件后端 v1.1.0 或注册该管理路由。');throw new Error('操作端点 '+endpoint+' 返回非 JSON：'+String(text||'').slice(0,90))}}
async function managementPluginPost(path,payload){var resp=await managementFetch('/plugins/grok-panel/'+String(path).replace(/^\/+/,''),{method:'POST',headers:{'content-type':'application/json'},body:JSON.stringify(payload||{})});var text=await resp.text();var data=parseJsonText(text,path);if(!resp.ok)throw new Error('HTTP '+resp.status+'：'+messageFromData(data,text));return data||{}}
async function runPluginChecks(emails){var indices=[];emails.forEach(function(email){var x=accountByEmail(email);var idx=x&&String(x.auth_index||x.authIndex||'').trim();if(idx)indices.push(idx)});var records=[];for(var i=0;i<indices.length;i++){var data=await managementPluginPost('checks',{auth_index:indices[i]});if(Array.isArray(data.records))records=records.concat(data.records)}return{records:records}}
async function deleteAuthNames(names){var resp=await managementFetch('/auth-files',{method:'DELETE',headers:{'content-type':'application/json'},body:JSON.stringify({names:names})});var text=await resp.text();var data=parseJsonText(text,'auth-files');if(!resp.ok)throw new Error('HTTP '+resp.status+'：'+messageFromData(data,text));return data||{}}
function messageFromData(data,text){if(data&&data.error&&data.error.message)return data.error.message;if(data&&data.message)return data.message;if(data&&data.error)return String(data.error);return String(text||'操作失败').slice(0,120)}
async function fetchData(showFeedback,skipAuto){if(showFeedback)setFeedback('正在刷新数据...','');byId('statusDot').classList.remove('err');try{var resp=await fetch(fixedApiUrl('data'),{credentials:'same-origin',headers:{'accept':'application/json'}});var text=await resp.text();if(resp.status===401||resp.status===403)throw new Error('数据接口授权不可用：请检查 CPA 管理会话。');if(!resp.ok)throw new Error('HTTP '+resp.status);var data=parseJsonText(text,'data');if(data&&data.ok===true&&data.result!==undefined)data=data.result;lastData=normalizeData(data);allData=lastData.files||[];allData.forEach(function(x,i){x._rowKey=makeRowKey(x,i)});pruneSelection();renderStats(lastData);renderChart(lastData.recent_buckets||[]);renderTable();byId('lastUpdate').textContent='更新于 '+fmtTime();if(showFeedback)setFeedback('数据已刷新。','ok');if(settings.autoCheck&&!skipAuto)maybeAutoCheck()}catch(e){byId('statusDot').classList.add('err');setFeedback('连接失败：'+e.message,'error')}}
function normalizeData(data){data=data||{};if(!Array.isArray(data.files))data.files=[];return data}
function makeRowKey(x,i){var email=getEmail(x);return email?email.toLowerCase():'row-'+i}
function getEmail(x){return String((x&&(x.email||x.account||x.name||x.id))||'').trim()}
function getAuthName(x){return String((x&&(x.name||x.auth_index||x.authIndex))||'').trim()}
function getStatus(x){return String((x&&x.status)||'').trim()}
function tokenLimit(){return Math.max(1,parseInt(settings.tokenLimit,10)||2000000)}
function tokensPerReq(){return Math.max(1,parseInt(settings.tokensPerReq,10)||5000)}
function failThreshold(){return Math.max(1,parseInt(settings.threshold,10)||3)}
function usagePct(x){var tl=tokenLimit();var et=(Number(x&&x.success)||0)*tokensPerReq();return tl>0?Math.max(0,Math.min(999,et/tl*100)):0}
function isMobileView(){return window.matchMedia&&window.matchMedia('(orientation:portrait),(hover:none) and (pointer:coarse)').matches}
function meterWidth(pct,kind){var span=kind==='row'?(isMobileView()?24:11):(isMobileView()?80:18);var clamped=Math.max(0,Math.min(100,Number(pct)||0));return(clamped*span/100).toFixed(2)+'vw'}
function classifyType(x){var key=makeRowKey(x,0);var raw=typeOverrides[key]||x.tier||x.account_type||x.accountType||x.account_kind||x.accountKind||x.plan||x.type||x.label||'';raw=String(raw||'').trim();var low=raw.toLowerCase();if(!raw||low==='unknown'||low==='unk')return{key:'unknown',label:'unknown'};if(low.indexOf('super')>=0||low.indexOf('premium')>=0||low.indexOf('paid')>=0||low.indexOf('pro')>=0||low.indexOf('max')>=0)return{key:'super',label:raw};if(low.indexOf('heavy')>=0||low.indexOf('bulk')>=0||low.indexOf('team')>=0)return{key:'heavy',label:raw};return{key:'standard',label:raw}}
function isHeavyAccount(x){return classifyType(x).key==='heavy'}
function statusKey(x){var s=getStatus(x).toLowerCase();if(x&&x.disabled)return'disabled';if(!s)return'unknown';if(s.indexOf('disable')>=0||s.indexOf('off')>=0)return'disabled';if(s.indexOf('active')>=0||s.indexOf('ok')>=0||s.indexOf('ready')>=0||s.indexOf('available')>=0)return'active';return'other'}
function mapHealth(raw){var low=String(raw||'').toLowerCase();if(!low)return null;if(low.indexOf('disabled')>=0||low.indexOf('off')>=0)return{key:'disabled',label:'禁用',detail:raw};if(low.indexOf('invalid')>=0||low.indexOf('expired')>=0||low.indexOf('revoked')>=0||low.indexOf('unavailable')>=0||low.indexOf('error')>=0||low.indexOf('dead')>=0)return{key:'invalid',label:'无效',detail:raw};if(low.indexOf('warn')>=0||low.indexOf('limited')>=0||low.indexOf('rate')>=0||low.indexOf('fail')>=0)return{key:'warn',label:'警告',detail:raw};if(low.indexOf('healthy')>=0||low.indexOf('active')>=0||low.indexOf('ok')>=0||low.indexOf('valid')>=0)return{key:'healthy',label:'健康',detail:raw};if(low.indexOf('unknown')>=0)return{key:'unknown',label:'未知',detail:raw};return null}
function deriveHealth(x){var key=makeRowKey(x,0);if(healthOverrides[key])return healthOverrides[key];var mapped=mapHealth(x.health||x.account_health||x.accountHealth||x.health_status||x.healthStatus);if(mapped)return mapped;if(x&&x.disabled)return{key:'disabled',label:'禁用',detail:'CPA 已禁用'};var s=getStatus(x).toLowerCase();if(statusKey(x)==='active')return{key:'healthy',label:'健康',detail:s||'CPA active'};return{key:'unknown',label:'未检查',detail:'需要执行健康检查'}}
function protectReason(x){var t=classifyType(x);if(settings.protectSuper&&t.key==='super')return'保护 super';if(settings.protectHeavy&&isHeavyAccount(x))return'保护 heavy';if(settings.protectUnknown&&t.key==='unknown')return'保护 unknown';return''}
function isInvalidCandidate(x){var h=deriveHealth(x).key;return h==='invalid'||h==='disabled'}
function cleanupCandidates(list){return(list||getFilteredData()).filter(function(x){return getEmail(x)&&isInvalidCandidate(x)&&!protectReason(x)})}
function renderStats(d){d=d||{};var total=d.total_files||allData.length||0;var active=d.active_files||0;var dis=d.disabled_files||0;var ts=d.total_success||0;var tf=d.total_failed||0;if(!d.total_files&&allData.length){active=allData.filter(function(x){return statusKey(x)==='active'}).length;dis=allData.filter(function(x){return statusKey(x)==='disabled'}).length;ts=sum(allData,'success');tf=sum(allData,'failed')}var et=ts*tokensPerReq();var cap=total*tokenLimit();var pct=cap>0?et/cap*100:0;byId('statTotal').textContent=total;byId('statTotalSub').textContent=dis+' 个已禁用';byId('statActive').textContent=active;byId('statActiveSub').textContent=dis+' 个禁用';byId('statRequests').textContent=fmt(ts);byId('statRequestsSub').textContent='失败 '+fmt(tf)+' | 率 '+(ts+tf>0?(ts/(ts+tf)*100).toFixed(1):0)+'%';byId('statTokens').textContent=fmt(et);byId('statTokensSub').textContent=tokensPerReq()+' token/请求';byId('statCapacity').textContent=fmt(cap);byId('statCapacitySub').textContent=tokenLimit()+' token/账号';byId('statUsage').textContent=pct.toFixed(2)+'%';byId('statUsageSub').textContent=fmt(et)+' / '+fmt(cap);var bar=byId('usageBar');bar.style.width=meterWidth(pct,'stat');bar.className='bar-fill'+(pct>80?' danger':pct>50?' warn':'');var typeCounts={standard:0,super:0,heavy:0,unknown:0};var healthCounts={healthy:0,warn:0,invalid:0,disabled:0,unknown:0};allData.forEach(function(x){var t=classifyType(x).key;typeCounts[t]=(typeCounts[t]||0)+1;var h=deriveHealth(x).key;healthCounts[h]=(healthCounts[h]||0)+1});var invalidAll=allData.filter(isInvalidCandidate);var cleanAll=cleanupCandidates(allData);byId('statTypes').textContent='S '+typeCounts.super+' / H '+typeCounts.heavy;byId('statTypesSub').textContent='standard '+typeCounts.standard+' | unknown '+typeCounts.unknown;byId('statHealth').textContent=healthCounts.healthy+' 健康';byId('statHealthSub').textContent='警告 '+healthCounts.warn+' | 无效 '+healthCounts.invalid+' | 未知 '+healthCounts.unknown;byId('statInvalid').textContent=invalidAll.length;byId('statInvalidSub').textContent=cleanAll.length+' 可清理 | '+(invalidAll.length-cleanAll.length)+' 受保护';byId('statSelected').textContent=getSelectedEmails().length;byId('statSelectedSub').textContent=getFilteredData().length+' 个当前可见'}
function sum(arr,key){return arr.reduce(function(n,x){return n+(Number(x&&x[key])||0)},0)}
function renderChart(buckets){var c=byId('chartRow');c.innerHTML='';if(!buckets||!buckets.length){c.innerHTML='<div class="chart-empty">暂无趋势数据</div>';return}buckets=buckets.slice().sort(function(a,b){return a.time<b.time?-1:1});var max=1;buckets.forEach(function(b){max=Math.max(max,(Number(b.success)||0)+(Number(b.failed)||0))});buckets.forEach(function(b){var total=(Number(b.success)||0)+(Number(b.failed)||0);var h=Math.max(0.35,total/max*14);var bar=document.createElement('div');bar.className='chart-bar'+((Number(b.failed)||0)>(Number(b.success)||0)?' fail':'');bar.style.height=h.toFixed(2)+'vh';bar.setAttribute('data-tip',String(b.time||'')+' | 成功:'+fmt(b.success)+' 失败:'+fmt(b.failed));c.appendChild(bar)})}
function getFilteredData(){var s=byId('searchBox')?byId('searchBox').value.toLowerCase().trim():'';var sf=byId('statusFilter')?byId('statusFilter').value:'all';var tf=byId('typeFilter')?byId('typeFilter').value:'all';var hf=byId('healthFilter')?byId('healthFilter').value:'all';var uf=byId('usageFilter')?byId('usageFilter').value:'all';var arr=allData.filter(function(x){var email=getEmail(x).toLowerCase();var st=getStatus(x).toLowerCase();var typ=classifyType(x);var health=deriveHealth(x);var pct=usagePct(x);if(s&&email.indexOf(s)<0&&st.indexOf(s)<0&&typ.label.toLowerCase().indexOf(s)<0&&health.label.indexOf(s)<0)return false;if(sf!=='all'&&statusKey(x)!==sf)return false;if(tf!=='all'&&typ.key!==tf)return false;if(hf!=='all'&&health.key!==hf)return false;if(uf==='unused'&&(Number(x.success)||0)!==0)return false;if(uf==='low'&&pct>=50)return false;if(uf==='warn'&&(pct<50||pct>=80))return false;if(uf==='high'&&pct<80)return false;return true});var sort=byId('sortFilter')?byId('sortFilter').value:'success_desc';arr.sort(function(a,b){if(sort==='failed_desc')return(Number(b.failed)||0)-(Number(a.failed)||0);if(sort==='usage_desc')return usagePct(b)-usagePct(a);if(sort==='health_asc')return healthRank(a)-healthRank(b);if(sort==='type_asc')return classifyType(a).key<classifyType(b).key?-1:1;if(sort==='email_asc')return getEmail(a)<getEmail(b)?-1:1;return(Number(b.success)||0)-(Number(a.success)||0)});return arr}
function healthRank(x){var order={invalid:0,disabled:1,warn:2,unknown:3,healthy:4};return order[deriveHealth(x).key]===undefined?9:order[deriveHealth(x).key]}
function renderTable(){var tbody=byId('tableBody');if(!tbody)return;var f=getFilteredData();tbody.innerHTML='';f.forEach(function(x,i){var email=getEmail(x);var key=makeRowKey(x,i);var su=Number(x.success)||0;var fa=Number(x.failed)||0;var et=su*tokensPerReq();var pct=usagePct(x);var bc=pct>80?'danger':pct>50?'warn':'';var t=classifyType(x);var h=deriveHealth(x);var prot=protectReason(x);var armed=isArmed('delete:'+key);var disabled=busy||!email||!!prot;var tr=document.createElement('tr');tr.className='row-'+h.key;var protectedText=prot?'<div class="cell-sub">'+esc(prot)+'</div>':'';var heavyText=(isHeavyAccount(x)&&t.key!=='heavy')?'<div class="cell-sub">heavy usage</div>':'';var healthDetail=h.detail?'<div class="cell-sub">'+esc(h.detail)+'</div>':'';tr.innerHTML='<td><input type="checkbox" class="row-select" data-key="'+esc(key)+'" data-email="'+esc(email)+'" '+(selected[key]?'checked':'')+' '+(!email||busy?'disabled':'')+'></td><td>'+esc(i+1)+'</td><td class="email-cell">'+esc(email||'?')+protectedText+'</td><td><span class="tag '+esc(t.key)+'">'+esc(t.key)+'</span><div class="cell-sub">'+esc(t.label)+'</div>'+heavyText+'</td><td><span class="tag '+esc(h.key)+'">'+esc(h.label)+'</span>'+healthDetail+'</td><td><span class="tag '+esc(statusKey(x))+'">'+esc(statusKey(x))+'</span><div class="cell-sub">'+esc(getStatus(x)||'-')+'</div></td><td>'+fmt(su)+'</td><td>'+(fa>0?'<span class="red-text">'+fmt(fa)+'</span>':'0')+'</td><td>'+fmt(et)+'</td><td><div class="cell-sub">'+pct.toFixed(1)+'%</div><div class="bar-container small"><div class="bar-fill '+bc+'" style="width:'+meterWidth(pct,'row')+'"></div></div></td><td><div class="actions-cell"><button data-act="check" data-key="'+esc(key)+'" '+(busy||!email?'disabled':'')+'>检查</button><button data-act="delete" data-key="'+esc(key)+'" '+(disabled?'disabled':'')+' class="'+(armed?'armed':'')+'">'+(armed?'确认删除':'删除')+'</button></div></td>';tbody.appendChild(tr)});if(!f.length)tbody.innerHTML='<tr><td colspan="11" class="muted">无匹配结果</td></tr>';renderSummary(f);updateToolbarState(f);if(lastData)renderStatsOnlySelection()}
function renderStatsOnlySelection(){byId('statSelected').textContent=getSelectedEmails().length;byId('statSelectedSub').textContent=getFilteredData().length+' 个当前可见'}
function renderSummary(f){var used=f.filter(function(x){return(Number(x.success)||0)>0}).length;var avail=f.filter(function(x){return(Number(x.success)||0)===0&&!x.disabled}).length;var ts=sum(f,'success');var tf=sum(f,'failed');var invalid=f.filter(isInvalidCandidate).length;var clean=cleanupCandidates(f).length;byId('sumCount').textContent=f.length;byId('sumUsed').textContent=used;byId('sumAvail').textContent=avail;byId('sumRate').textContent=(ts+tf>0?(ts/(ts+tf)*100).toFixed(1):0)+'%';byId('sumInvalid').textContent=invalid;byId('sumCleanup').textContent=clean;byId('filterInfo').textContent=f.length+' / '+allData.length+' 可见，'+clean+' 可清理'}
function renderAll(){if(lastData){renderStats(lastData);renderChart(lastData.recent_buckets||[])}renderTable()}
function handleTableClick(evt){var btn=evt.target.closest('button[data-act]');if(!btn)return;var key=btn.getAttribute('data-key');if(btn.getAttribute('data-act')==='check')checkOne(key);if(btn.getAttribute('data-act')==='delete')requestSingleDelete(key)}
function handleTableChange(evt){if(!evt.target.classList.contains('row-select'))return;var key=evt.target.getAttribute('data-key');var email=evt.target.getAttribute('data-email');if(evt.target.checked)selected[key]=email;else delete selected[key];renderAll()}
function toggleVisibleSelection(){var on=byId('selectVisible').checked;getFilteredData().forEach(function(x,i){var email=getEmail(x);if(!email)return;var key=makeRowKey(x,i);if(on)selected[key]=email;else delete selected[key]});renderAll()}
function pruneSelection(){var valid={};allData.forEach(function(x,i){valid[makeRowKey(x,i)]=true});Object.keys(selected).forEach(function(k){if(!valid[k])delete selected[k]})}
function getSelectedEmails(){var out=[];Object.keys(selected).forEach(function(k){if(selected[k])out.push(selected[k])});return unique(out)}
function unique(arr){var seen={};var out=[];(arr||[]).forEach(function(v){v=String(v||'').trim();var k=v.toLowerCase();if(v&&!seen[k]){seen[k]=true;out.push(v)}});return out}
function updateToolbarState(f){if(!byId('selectionInfo'))return;f=f||getFilteredData();var count=getSelectedEmails().length;var visibleWithEmail=f.filter(function(x){return!!getEmail(x)});var selectedVisible=visibleWithEmail.filter(function(x,i){return!!selected[makeRowKey(x,i)]}).length;var cb=byId('selectVisible');cb.checked=visibleWithEmail.length>0&&selectedVisible===visibleWithEmail.length;cb.indeterminate=selectedVisible>0&&selectedVisible<visibleWithEmail.length;byId('selectionInfo').textContent=count?('已选择 '+count+' 个'):'未选择';byId('checkSelectedBtn').disabled=busy||count===0;byId('batchDeleteBtn').disabled=busy||count===0;byId('batchDeleteBtn').textContent=isArmed('delete:batch')?'确认删除选中':'删除选中';if(isArmed('delete:batch'))byId('batchDeleteBtn').classList.add('armed');else byId('batchDeleteBtn').classList.remove('armed');var clean=cleanupCandidates(f).length;byId('cleanupInvalidBtn').disabled=busy||clean===0;byId('cleanupInvalidBtn').textContent=isArmed('cleanup:visible')?'确认清理无效':'清理无效 '+clean;if(isArmed('cleanup:visible'))byId('cleanupInvalidBtn').classList.add('armed');else byId('cleanupInvalidBtn').classList.remove('armed');byId('checkVisibleTopBtn').disabled=busy||visibleWithEmail.length===0;byId('refreshBtn').disabled=busy}
function isArmed(key){return(confirmUntil[key]||0)>Date.now()}
function armOrProceed(key,msg){if(isArmed(key)){delete confirmUntil[key];return true}confirmUntil[key]=Date.now()+6500;setFeedback(msg,'warn');window.setTimeout(function(){if(confirmUntil[key]&&confirmUntil[key]<Date.now()){delete confirmUntil[key];renderAll()}},6800);renderAll();return false}
function accountByKey(key){for(var i=0;i<allData.length;i++){if(makeRowKey(allData[i],i)===key)return allData[i]}return null}
function accountByEmail(email){var target=String(email||'').toLowerCase();for(var i=0;i<allData.length;i++){if(getEmail(allData[i]).toLowerCase()===target)return allData[i]}return null}
function checkOne(key){var x=accountByKey(key);var email=x?getEmail(x):'';if(!email){setFeedback('无法检查：账号缺少邮箱或 ID。','error');return}runHealthCheck([email],'手动检查 1 个',false)}
function manualCheckVisible(){var emails=unique(getFilteredData().map(getEmail));if(!emails.length){setFeedback('没有可检查的可见账号。','warn');return}runHealthCheck(emails,'手动检查可见 '+emails.length+' 个',false)}
function manualCheckSelected(){var emails=getSelectedEmails();if(!emails.length){setFeedback('请先选择账号。','warn');return}runHealthCheck(emails,'手动检查选中 '+emails.length+' 个',false)}
async function runHealthCheck(emails,label,autoMode){emails=unique(emails);if(!emails.length)return;if(busy&&!autoMode)return;setBusy(true);try{var data=await runPluginChecks(emails);var results=applyCheckResult(data);setFeedback(label+' 完成：收到 '+results.length+' 条检查结果。','ok');if(settings.autoDelete)await runAutoCleanup();renderAll();fetchData(false,true)}catch(e){setFeedback(label+' 失败：'+e.message,'error')}finally{setBusy(false)}}
function applyCheckResult(data){var results=[];if(data){if(Array.isArray(data.records))results=data.records;else if(Array.isArray(data.results))results=data.results;else if(Array.isArray(data.files))results=data.files;else if(Array.isArray(data.accounts))results=data.accounts;else if(Array.isArray(data.checked))results=data.checked;else if(typeof data==='object'){Object.keys(data).forEach(function(k){if(data[k]&&typeof data[k]==='object'){var v=data[k];if(!v.email)v.email=k;results.push(v)}})}}results.forEach(function(r){var email=String(r.email||r.account||r.id||'').trim();if(!email)return;var key=email.toLowerCase();var mapped=mapHealth(r.health||r.status||r.result||'')||{key:'unknown',label:'未知',detail:r.reason||r.detail||r.message||'checked'};mapped.detail=String(r.reason||r.detail||r.message||mapped.detail||'checked');healthOverrides[key]=mapped;var tier=r.classification&&r.classification.tier;if(tier||r.account_type||r.accountType||r.type)typeOverrides[key]=tier||r.account_type||r.accountType||r.type;var x=accountByEmail(email);if(x){x.delete_eligible=r.delete_eligible===true;x.invalid_streak=Number(r.invalid_streak)||0;x.health=r.health||x.health;x.tier=tier||x.tier}});return results}
function requestSingleDelete(key){var x=accountByKey(key);if(!x){setFeedback('无法删除：找不到该账号。','error');return}var email=getEmail(x);var prot=protectReason(x);if(prot){setFeedback('删除已阻止：'+email+' 受 '+prot+' 保护。','warn');return}if(!armOrProceed('delete:'+key,'再次点击该行“确认删除”将删除 '+email+'。'))return;deleteAccounts([email],'删除 '+email)}
function requestBatchDelete(){var emails=getSelectedEmails();if(!emails.length){setFeedback('请先选择账号。','warn');return}if(!armOrProceed('delete:batch','再次点击“确认删除选中”将删除未受保护的选中账号。'))return;deleteAccounts(emails,'批量删除')}
async function deleteAccounts(emails,label){emails=unique(emails);var names=[];var blocked=[];emails.forEach(function(email){var x=accountByEmail(email);var prot=x?protectReason(x):(settings.protectUnknown?'保护 unknown':'');var name=x?getAuthName(x):'';if(prot)blocked.push(email+'('+prot+')');else if(name)names.push(name);else blocked.push(email+'(缺少 auth 文件名)')});names=unique(names);if(!names.length){setFeedback(label+' 已取消：没有未受保护且具备文件名的账号。'+(blocked.length?' 已阻止 '+blocked.length+' 个。':''),'warn');return}setBusy(true);try{var data=await deleteAuthNames(names);emails.forEach(function(email){delete selected[email.toLowerCase()]});setFeedback(label+' 完成：已请求删除 '+names.length+' 个 auth 文件。'+mutationSummary(data)+(blocked.length?' 阻止 '+blocked.length+' 个。':''),'ok');await fetchData(false,true)}catch(e){setFeedback(label+' 失败：'+e.message,'error')}finally{setBusy(false)}}
function requestCleanupInvalid(){var candidates=cleanupCandidates(getFilteredData());if(!candidates.length){setFeedback('当前筛选范围没有未受保护的无效账号可清理。','warn');return}if(!armOrProceed('cleanup:visible','再次点击“确认清理无效”将清理当前筛选范围内 '+candidates.length+' 个未受保护的无效账号。'))return;cleanupInvalid(candidates,'清理无效')}
async function cleanupInvalid(candidates,label){var emails=unique(candidates.filter(function(x){return x.delete_eligible===true&&!protectReason(x)}).map(getEmail));if(!emails.length){setFeedback('没有达到连续 '+failThreshold()+' 次明确 401/403 且未受保护的账号。','warn');return}await deleteAccounts(emails,label)}
function currentProtectSettings(){return{super:!!settings.protectSuper,heavy:!!settings.protectHeavy,unknown:!!settings.protectUnknown}}
function mutationSummary(data){if(!data)return'';var parts=[];if(data.deleted!==undefined)parts.push('已删除 '+data.deleted);if(data.skipped!==undefined)parts.push('跳过 '+data.skipped);if(data.checked!==undefined)parts.push('检查 '+data.checked);if(data.message)parts.push(data.message);return parts.length?'（'+parts.join('，')+'）':''}
function maybeAutoCheck(){if(autoCheckBusy)return;var now=Date.now();if(now-lastAutoCheckAt<300000)return;var emails=unique(allData.filter(function(x){return getEmail(x)&&!x.disabled}).map(getEmail));if(!emails.length)return;lastAutoCheckAt=now;autoCheckBusy=true;runHealthCheck(emails,'自动检查 '+emails.length+' 个',true).finally(function(){autoCheckBusy=false})}
async function runAutoCleanup(){if(!settings.autoDelete)return;var candidates=cleanupCandidates(allData);if(!candidates.length){setFeedback('自动检查完成：没有未受保护的无效账号。','ok');return}await cleanupInvalid(candidates,'自动删除无效')}
function init(){initSettings();bindEvents();fetchData(true,false);refreshHandle=window.setInterval(function(){fetchData(false,false)},30000)}
init();
</script>
</body></html>`

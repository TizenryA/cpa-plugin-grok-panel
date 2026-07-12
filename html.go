package main

const htmlPage = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Grok 面板</title>
<style>
:root{--bg:#f5f5f0;--card:#fff;--border:#1a1a1a;--text:#1a1a1a;--muted:#888;--bar:#e8e8e3;--green:#2d7d32;--red:#c62828;--yellow:#f9a825}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--text);font-family:'SF Mono','Cascadia Code','JetBrains Mono',Consolas,monospace;font-size:14px;line-height:1.6;padding:1.5vw}
h1{font-size:1.4rem;font-weight:700;margin-bottom:.5rem}
h2{font-size:1rem;font-weight:700;margin-bottom:.5rem;padding-bottom:.3rem;border-bottom:2px solid var(--border)}
.container{max-width:1200px;margin:0 auto}
.header{display:flex;justify-content:space-between;align-items:center;margin-bottom:1rem;flex-wrap:wrap;gap:.5rem}
.header-left{display:flex;align-items:center;gap:.8rem}
.status-dot{width:.7rem;height:.7rem;background:var(--green);display:inline-block}
.status-dot.err{background:var(--red)}
.refresh-btn{background:var(--border);color:var(--card);border:2px solid var(--border);padding:.4rem 1rem;cursor:pointer;font-family:inherit;font-size:.8rem;font-weight:600}
.refresh-btn:hover{background:var(--card);color:var(--border)}
.config-bar{background:var(--card);border:2px solid var(--border);padding:.8rem 1rem;display:flex;align-items:center;gap:1rem;flex-wrap:wrap;margin-bottom:1.5rem}
.config-bar label{font-size:.75rem;color:var(--muted)}
.config-bar input{background:var(--bg);border:2px solid var(--border);padding:.3rem .5rem;font-family:inherit;font-size:.8rem;width:6rem}
.config-bar input:focus{outline:none;border-color:#1565c0}
.stats-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(160px,1fr));gap:.8rem;margin-bottom:1.5rem}
.stat-card{background:var(--card);border:2px solid var(--border);padding:1rem}
.stat-label{font-size:.7rem;color:var(--muted);text-transform:uppercase;letter-spacing:.05em}
.stat-value{font-size:1.5rem;font-weight:700;margin-top:.2rem}
.stat-sub{font-size:.7rem;color:var(--muted);margin-top:.2rem}
.section{margin-bottom:1.5rem}
table{width:100%;border-collapse:collapse;background:var(--card)}
th,td{text-align:left;padding:.5rem .6rem;border-bottom:1px solid var(--bar);font-size:.8rem}
th{background:var(--border);color:var(--card);font-weight:600;text-transform:uppercase;font-size:.7rem;letter-spacing:.03em;position:sticky;top:0}
tr:hover{background:var(--bg)}
.email-cell{font-family:inherit;word-break:break-all}
.tag{display:inline-block;padding:.1rem .4rem;font-size:.65rem;font-weight:600;border:1px solid var(--border)}
.tag.active{background:var(--green);color:var(--card);border-color:var(--green)}
.tag.disabled{background:var(--red);color:var(--card);border-color:var(--red)}
.bar-container{width:100%;height:.5rem;background:var(--bar);position:relative;margin-top:.2rem}
.bar-fill{height:100%;background:var(--green);transition:width .3s}
.bar-fill.warn{background:var(--yellow)}
.bar-fill.danger{background:var(--red)}
.chart-row{display:flex;align-items:flex-end;gap:2px;height:5rem;padding:.5rem 0;overflow-x:auto}
.chart-bar{flex:1;min-width:3px;background:var(--green);opacity:.7;position:relative}
.chart-bar.fail{background:var(--red)}
.chart-bar:hover{opacity:1}
.chart-bar:hover::after{content:attr(data-tip);position:absolute;bottom:100%;left:50%;transform:translateX(-50%);background:var(--border);color:var(--card);padding:.2rem .4rem;font-size:.65rem;white-space:nowrap;z-index:10}
.search-box{background:var(--card);border:2px solid var(--border);padding:.4rem .6rem;font-family:inherit;font-size:.8rem;width:100%;margin-bottom:.5rem}
.search-box:focus{outline:none;border-color:#1565c0}
.loading{text-align:center;padding:2rem;color:var(--muted)}
.error-msg{background:var(--red);color:var(--card);padding:.8rem;border:2px solid var(--border);margin-bottom:1rem;font-size:.8rem}
.summary-row{display:flex;gap:1.5rem;flex-wrap:wrap;font-size:.75rem;color:var(--muted);margin-top:.5rem}
.summary-row b{color:var(--text)}
@media(max-width:600px){.stats-grid{grid-template-columns:repeat(2,1fr)}.stat-value{font-size:1.1rem}th,td{font-size:.7rem;padding:.3rem}.config-bar{font-size:.7rem}}
</style>
</head>
<body>
<div class="container">
<div class="header">
<div class="header-left"><h1>GROK 面板</h1><span class="status-dot" id="statusDot"></span><span style="font-size:.75rem;color:var(--muted)" id="lastUpdate">--</span></div>
<button class="refresh-btn" onclick="fetchData()">刷新</button>
</div>
<div id="errorBox"></div>
<div class="stats-grid">
<div class="stat-card"><div class="stat-label">Grok 文件数</div><div class="stat-value" id="statTotal">--</div><div class="stat-sub" id="statTotalSub"></div></div>
<div class="stat-card"><div class="stat-label">活跃账号</div><div class="stat-value" id="statActive">--</div><div class="stat-sub" id="statActiveSub"></div></div>
<div class="stat-card"><div class="stat-label">总请求数</div><div class="stat-value" id="statRequests">--</div><div class="stat-sub" id="statRequestsSub"></div></div>
<div class="stat-card"><div class="stat-label">估算 Token</div><div class="stat-value" id="statTokens">--</div><div class="stat-sub" id="statTokensSub"></div></div>
<div class="stat-card"><div class="stat-label">总容量</div><div class="stat-value" id="statCapacity">--</div><div class="stat-sub" id="statCapacitySub"></div></div>
<div class="stat-card"><div class="stat-label">使用率</div><div class="stat-value" id="statUsage">--</div><div class="stat-sub" id="statUsageSub"></div><div class="bar-container"><div class="bar-fill" id="usageBar" style="width:0%"></div></div></div>
</div>
<div class="config-bar">
<label>每账号上限(Token)</label><input type="number" id="tokenLimit" value="2000000">
<label>每请求估算(Token)</label><input type="number" id="tokensPerReq" value="5000">
</div>
<div class="section"><h2>请求趋势</h2><div class="chart-row" id="chartRow"></div></div>
<div class="section"><h2>账号明细</h2>
<input type="text" class="search-box" id="searchBox" placeholder="搜索邮箱..." oninput="renderTable()">
<div style="max-height:60vh;overflow:auto;border:2px solid var(--border)">
<table><thead><tr><th>#</th><th>邮箱</th><th>状态</th><th>成功</th><th>失败</th><th>估算Token</th><th>用量</th></tr></thead>
<tbody id="tableBody"></tbody></table></div>
<div class="summary-row"><span>共 <b id="sumCount">0</b> 个</span><span>已使用 <b id="sumUsed">0</b> 个</span><span>剩余 <b id="sumAvail">0</b> 个</span><span>成功率 <b id="sumRate">0%</b></span></div>
</div></div>
<script>
var allData=[];
function fmt(n){if(n>=1e9)return(n/1e9).toFixed(2)+'B';if(n>=1e6)return(n/1e6).toFixed(2)+'M';if(n>=1e3)return(n/1e3).toFixed(1)+'K';return String(n)}
function fmtTime(){var d=new Date();return d.toLocaleTimeString('zh-CN',{hour12:false})+' '+d.toLocaleDateString('zh-CN')}
async function fetchData(){
document.getElementById('errorBox').innerHTML='';
document.getElementById('statusDot').classList.remove('err');
try{
var base=window.location.pathname.replace(/\/$/,'');
var resp=await fetch(base+'/data');
if(!resp.ok)throw new Error('HTTP '+resp.status);
var data=await resp.json();
allData=data.files||[];
allData.sort(function(a,b){return(b.success||0)-(a.success||0)});
renderStats(data);
renderChart(data.recent_buckets||[]);
renderTable();
document.getElementById('lastUpdate').textContent='更新于 '+fmtTime();
}catch(e){
document.getElementById('errorBox').innerHTML='<div class="error-msg">连接失败: '+e.message+'</div>';
document.getElementById('statusDot').classList.add('err');
}
}
function renderStats(d){
var tl=parseInt(document.getElementById('tokenLimit').value)||2e6;
var tpr=parseInt(document.getElementById('tokensPerReq').value)||5e3;
var total=d.total_files||0,active=d.active_files||0,dis=d.disabled_files||0;
var ts=d.total_success||0,tf=d.total_failed||0;
var et=ts*tpr,cap=total*tl,pct=cap>0?(et/cap*100):0;
document.getElementById('statTotal').textContent=total;
document.getElementById('statTotalSub').textContent=dis+' 个已禁用';
document.getElementById('statActive').textContent=active;
document.getElementById('statActiveSub').textContent=dis+' 个禁用';
document.getElementById('statRequests').textContent=fmt(ts);
document.getElementById('statRequestsSub').textContent='失败 '+fmt(tf)+' | 率 '+(ts+tf>0?(ts/(ts+tf)*100).toFixed(1):0)+'%';
document.getElementById('statTokens').textContent=fmt(et);
document.getElementById('statTokensSub').textContent=tpr+' token/请求';
document.getElementById('statCapacity').textContent=fmt(cap);
document.getElementById('statCapacitySub').textContent=tl+' token/账号';
document.getElementById('statUsage').textContent=pct.toFixed(2)+'%';
document.getElementById('statUsageSub').textContent=fmt(et)+' / '+fmt(cap);
var bar=document.getElementById('usageBar');
bar.style.width=Math.min(pct,100)+'%';
bar.className='bar-fill'+(pct>80?' danger':pct>50?' warn':'');
document.getElementById('sumCount').textContent=total;
document.getElementById('sumUsed').textContent=allData.filter(function(f){return(f.success||0)>0}).length;
document.getElementById('sumAvail').textContent=allData.filter(function(f){return(f.success||0)===0&&!f.disabled}).length;
var rate=ts+tf>0?(ts/(ts+tf)*100).toFixed(1):0;
document.getElementById('sumRate').textContent=rate+'%';
}
function renderChart(buckets){
var c=document.getElementById('chartRow');c.innerHTML='';
if(!buckets||!buckets.length)return;
buckets.sort(function(a,b){return a.time<b.time?-1:1});
var max=Math.max.apply(Math,buckets.map(function(b){return b.success+b.failed}).concat([1]));
buckets.forEach(function(b){
var total=b.success+b.failed;
var h=(total/max*100);
var bar=document.createElement('div');
bar.className='chart-bar'+(b.failed>b.success?' fail':'');
bar.style.height=Math.max(2,h)+'%';
bar.setAttribute('data-tip',b.time+' | 成功:'+b.success+' 失败:'+b.failed);
c.appendChild(bar);
});
}
function renderTable(){
var tbody=document.getElementById('tableBody');
var s=document.getElementById('searchBox').value.toLowerCase();
var tl=parseInt(document.getElementById('tokenLimit').value)||2e6;
var tpr=parseInt(document.getElementById('tokensPerReq').value)||5e3;
var f=allData.filter(function(x){return!s||(x.email||'').toLowerCase().indexOf(s)>=0});
tbody.innerHTML='';
f.forEach(function(x,i){
var su=x.success||0,fa=x.failed||0,et=su*tpr,pct=Math.min(100,(et/tl*100));
var bc=pct>80?'danger':pct>50?'warn':'';
var tag=x.disabled?'<span class="tag disabled">禁用</span>':'<span class="tag active">活跃</span>';
var tr=document.createElement('tr');
tr.innerHTML='<td>'+(i+1)+'</td><td class="email-cell">'+(x.email||'?')+'</td><td>'+tag+'</td><td>'+su+'</td><td>'+(fa>0?'<span style="color:var(--red)">'+fa+'</span>':'0')+'</td><td>'+fmt(et)+'</td><td style="min-width:8rem"><div style="font-size:.7rem;color:var(--muted)">'+pct.toFixed(1)+'%</div><div class="bar-container"><div class="bar-fill '+bc+'" style="width:'+pct+'%"></div></div></td>';
tbody.appendChild(tr);
});
if(!f.length)tbody.innerHTML='<tr><td colspan="7" style="text-align:center;color:var(--muted);padding:2rem">无匹配结果</td></tr>';
}
fetchData();
setInterval(fetchData,30000);
</script>
</body></html>`

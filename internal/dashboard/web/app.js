const $ = (s) => document.querySelector(s);
const state = { matches: [], teams: [], team: 'vasco', status: '', search: '', selected: 1, tab: 'lineup', side: 'home', favoritesOnly: false, favorites: [] };
try { const saved=JSON.parse(localStorage.getItem('noc-fut-favorites') || '[]'); if(Array.isArray(saved)) state.favorites=saved.filter(Number.isInteger); } catch {}
const esc = (s) => String(s).replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]));
const crest = (t) => `<span class="crest ${t.id==='vasco'?'vasco':''}" style="--club:${esc(t.color)}" aria-hidden="true">${t.id==='vasco'?'✠':esc(t.short)}</span>`;
const statusText = m => m.status==='live'?`${m.minute}′ · Em jogo · demo`:m.status==='finished'?'Encerrado · demo':`${m.time} · A seguir`;
const scoreText = m => m.status==='scheduled'?m.time:`${m.homeScore} : ${m.awayScore}`;
function favorite(id) { state.favorites=state.favorites.includes(id)?state.favorites.filter(x=>x!==id):[...state.favorites,id];try{localStorage.setItem('noc-fut-favorites',JSON.stringify(state.favorites));}catch{} render(); }
function selectMatch(id) { state.selected=id;state.tab='lineup';const m=state.matches.find(x=>x.id===id);state.side=m?.away.id===state.team?'away':'home';render(); }
function filtered() { return state.matches.filter(m=>(!state.team||m.home.id===state.team||m.away.id===state.team)&&(!state.status||m.status===state.status)&&(!state.favoritesOnly||state.favorites.includes(m.id))&&(`${m.home.name} ${m.away.name}`.toLocaleLowerCase('pt-BR').includes(state.search.toLocaleLowerCase('pt-BR')))); }
function render() {
  const matches=filtered();
  $('#total-count').textContent=state.matches.length;
  $('#live-count').textContent=state.matches.filter(m=>m.status==='live').length;
  $('#next-count').textContent=state.matches.filter(m=>m.status==='scheduled').length;
  $('#fav-count').textContent=state.favorites.length;$('#summary-favs').textContent=state.favorites.length;
  $('#list-label').textContent=state.favoritesOnly?'SEUS FAVORITOS':state.team?`PARTIDAS · ${state.teams.find(t=>t.id===state.team)?.name||''}`:'TODAS AS PARTIDAS';
  $('#result-count').textContent=`${matches.length} ${matches.length===1?'JOGO':'JOGOS'}`;
  $('#nav-games').classList.toggle('active',!state.favoritesOnly);$('#nav-favorites').classList.toggle('active',state.favoritesOnly);
  document.querySelectorAll('[data-status]').forEach(b=>b.classList.toggle('active',b.dataset.status===state.status));
  if(!matches.some(m=>m.id===state.selected)){state.selected=matches[0]?.id??null;const m=matches[0];state.side=m?.away.id===state.team?'away':'home';}
  $('#match-list').innerHTML=matches.length?matches.map(m=>`<article class="match-card ${m.id===state.selected?'selected':''}" tabindex="0" role="button" aria-pressed="${m.id===state.selected}" data-match="${m.id}" aria-label="${esc(m.home.name)} contra ${esc(m.away.name)}"><div class="match-meta"><span class="${m.status==='live'?'live':''}">${statusText(m)}</span><button class="star ${state.favorites.includes(m.id)?'on':''}" data-favorite="${m.id}" aria-label="${state.favorites.includes(m.id)?'Remover dos':'Adicionar aos'} favoritos: ${esc(m.home.name)} e ${esc(m.away.name)}" aria-pressed="${state.favorites.includes(m.id)}">${state.favorites.includes(m.id)?'★':'☆'}</button></div><div class="teams-row"><div class="team">${crest(m.home)}<span>${esc(m.home.name)}</span></div><div class="score ${m.status==='scheduled'?'scheduled':''}">${scoreText(m)}</div><div class="team away"><span>${esc(m.away.name)}</span>${crest(m.away)}</div></div><div class="stadium">⌖ ${esc(m.stadium)} <span>· Amostra fictícia</span></div></article>`).join(''):`<div class="empty">Nenhuma partida com esses filtros.<button id="clear-filters">Limpar filtros</button></div>`;
  document.querySelectorAll('[data-match]').forEach(el=>{el.onclick=()=>selectMatch(Number(el.dataset.match));el.onkeydown=e=>{if(e.target!==el)return;if(e.key==='Enter'||e.key===' '){e.preventDefault();selectMatch(Number(el.dataset.match));}};});
  document.querySelectorAll('[data-favorite]').forEach(el=>el.onclick=e=>{e.stopPropagation();favorite(Number(el.dataset.favorite));});
  if($('#clear-filters'))$('#clear-filters').onclick=()=>{state.team='';state.status='';state.search='';state.favoritesOnly=false;$('#team-filter').value='';$('#search').value='';render();};
  renderDetail();
}
function renderDetail() {
  const m=state.matches.find(x=>x.id===state.selected);if(!m){$('#detail-panel').innerHTML='<div class="empty">Selecione outro filtro para encontrar uma partida.</div>';return;}
  $('#detail-panel').innerHTML=`<div class="detail-top"><div class="detail-caption"><span>BRASILEIRÃO · SÉRIE A</span><b>PARTIDA DEMONSTRATIVA</b></div><div class="detail-scoreboard"><div class="detail-team">${crest(m.home)}<strong>${esc(m.home.name)}</strong><small>MANDANTE</small></div><div><div class="detail-score">${m.status==='scheduled'?'— : —':scoreText(m)}</div><div class="state">${statusText(m)}</div></div><div class="detail-team">${crest(m.away)}<strong>${esc(m.away.name)}</strong><small>VISITANTE</small></div></div><div class="stadium-center">⌖ ${esc(m.stadium)} · Brasil</div></div><div class="detail-tabs" role="group" aria-label="Informações da partida"><button data-tab="lineup" class="${state.tab==='lineup'?'active':''}">Escalação</button><button data-tab="stats" class="${state.tab==='stats'?'active':''}">Estatísticas</button><button data-tab="events" class="${state.tab==='events'?'active':''}">Lances</button></div><div class="detail-content" id="detail-content"></div>`;
  document.querySelectorAll('[data-tab]').forEach(b=>b.onclick=()=>{state.tab=b.dataset.tab;renderDetail();});
  const content=$('#detail-content');
  if(state.tab==='lineup') {
    const players=m.lineups[state.side];
    content.innerHTML=`<div class="lineup-heading"><div class="side-switch"><button data-side="home" class="${state.side==='home'?'active':''}">${esc(m.home.short)}</button><button data-side="away" class="${state.side==='away'?'active':''}">${esc(m.away.short)}</button></div><span class="formation">4 — 3 — 3</span></div><div class="pitch" aria-label="Escalação fictícia de ${esc(m[state.side].name)}"><div class="pitch-lines"></div><div class="halfway"></div><div class="center-circle"></div><div class="penalty top"></div><div class="penalty bottom"></div><div class="goalbox top"></div><div class="goalbox bottom"></div>${players.map((p,i)=>`<div class="player ${i===0?'goalkeeper':''}" style="left:${p.x}%;top:${p.y}%" title="${p.number} · ${esc(p.name)}"><div class="jersey">${p.number}</div><span class="player-name">${esc(p.name)}</span></div>`).join('')}</div><p class="lineup-note">ⓘ Escalação ilustrativa: nomes e posições fictícios, sem vínculo com o elenco atual. Escudos estilizados.</p>`;
    document.querySelectorAll('[data-side]').forEach(b=>b.onclick=()=>{state.side=b.dataset.side;renderDetail();});
  } else if(state.tab==='stats') {
    content.innerHTML=m.status==='scheduled'?'<div class="empty">As estatísticas ficam disponíveis após o início da partida.<br>Esta é uma demonstração local.</div>':m.stats.map(s=>`<div class="stats-row"><div class="stats-label"><b>${s.home}${s.label==='Posse de bola'?'%':''}</b><span>${esc(s.label)}</span><b>${s.away}${s.label==='Posse de bola'?'%':''}</b></div><div class="stats-track"><i style="flex:${s.home}"></i><i style="flex:${s.away}"></i></div></div>`).join('')+'<p class="lineup-note">Valores fictícios para demonstração da interface.</p>';
  } else content.innerHTML=m.events.length?m.events.slice().reverse().map(e=>`<div class="event"><b>${e.minute}′</b><div>${e.type==='goal'?'⚽':e.type==='yellow'?'🟨':'↔'} ${esc(e.text)}<small>${esc(e.team)} · evento fictício</small></div></div>`).join(''):'<div class="empty">Nenhum lance nesta partida demonstrativa.</div>';
}
async function load() {
  $('#refresh').disabled=true;
  try {
    const responses=await Promise.all([fetch('/api/matches'),fetch('/api/teams')]);if(responses.some(r=>!r.ok))throw Error('Falha na API local');
    const [data,teams]=await Promise.all(responses.map(r=>r.json()));state.matches=data.matches;state.teams=teams;
    $('#team-filter').innerHTML='<option value="">Todos os times</option>'+teams.slice().sort((a,b)=>a.name.localeCompare(b.name,'pt-BR')).map(t=>`<option value="${esc(t.id)}">${esc(t.name)}</option>`).join('');$('#team-filter').value=state.team;
    render();$('#updated').textContent='Dados locais · carregados às '+new Date().toLocaleTimeString('pt-BR',{hour:'2-digit',minute:'2-digit'});
  } catch { $('#match-list').innerHTML='<div class="empty">Não foi possível carregar a API local.<br>Confirme se o servidor está ativo e tente recarregar.</div>';$('#detail-panel').innerHTML='<div class="empty">Aguardando conexão com o servidor local.</div>'; }
  finally {$('#refresh').disabled=false;}
}
$('#team-filter').onchange=e=>{state.team=e.target.value;state.side='home';state.tab='lineup';render();};
$('#search').oninput=e=>{state.search=e.target.value;render();};
document.querySelectorAll('[data-status]').forEach(b=>b.onclick=()=>{state.status=b.dataset.status;render();});
$('#nav-games').onclick=()=>{state.favoritesOnly=false;render();};
$('#nav-favorites').onclick=()=>{state.favoritesOnly=true;state.team='';state.status='';state.search='';$('#team-filter').value='';$('#search').value='';render();};
$('#refresh').onclick=load;
load();

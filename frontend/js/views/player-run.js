import { API } from '../api.js';
import { state, navigate, toast, confirmModal } from '../app.js';
import { LeafletMap, haversine } from '../map.js';

const BOSS_EMOJIS = ['🐉', '💀', '🧙', '👹', '🦇', '🕷️', '👻', '🧟', '🐍', '🦂'];

let map = null;

export const playerRunView = {
  runId: null,
  dungeonId: null,
  dungeon: null,
  run: null,
  steps: [],
  playerLat: null,
  playerLon: null,

  async init(data) {
    this.runId = data.runId;
    this.dungeonId = data.dungeonId;
    this.playerLat = null;
    this.playerLon = null;

    map = new LeafletMap('pr-map');
    map.onClick = (lat, lon) => this.movePlayer(lat, lon);

    document.getElementById('pr-attack').onclick = () => this.attempt();
    document.getElementById('pr-quit').onclick = () => navigate('player-list');
    document.getElementById('pr-abandon').onclick = () => this.abandon();
    document.getElementById('pr-victory').style.display = 'none';
    document.getElementById('pr-victory-ok').onclick = () => {
      document.getElementById('pr-victory').style.display = 'none';
    };

    await this.loadData();
  },

  destroy() {
    if (map) { map.destroy(); map = null; }
  },

  async loadData() {
    try {
      const [dRes, rRes] = await Promise.all([
        API.getDungeon(this.dungeonId),
        API.getRun(this.runId),
      ]);
      this.dungeon = dRes.data;
      this.run = rRes.data;
      this.run.killedSteps = this.run.killedSteps || [];
      this.steps = this.dungeon.bossSteps || [];

      document.getElementById('pr-title').textContent = `🏰 ${this.dungeon.title}`;
      this.refresh();
      if (this.steps.length > 0) map.fitToItems();
      this.loadGold();
    } catch (e) {
      console.error('loadData', e);
      this.msg('Erreur: ' + (e.message || JSON.stringify(e)), true);
    }
  },

  async reloadRun() {
    try {
      const res = await API.getRun(this.runId);
      this.run = res.data;
      this.run.killedSteps = this.run.killedSteps || [];
      this.refresh();
    } catch (e) { console.error('reloadRun', e); }
  },

  async loadGold() {
    try {
      const res = await API.getPlayer(state.playerId);
      document.getElementById('pr-gold').textContent = res.data?.gold || 0;
    } catch (_) {}
  },

  refresh() {
    const s = this.run.state;
    const killedCount = this.run.killedSteps.length;
    const emoji = s === 'active' ? '🟢' : s === 'completed' ? '🏆' : '⛔';
    document.getElementById('pr-state').innerHTML =
      `${emoji} ${s.toUpperCase()} · Step ${this.run.currentStep}/${this.steps.length} · ${killedCount} boss tues`;

    // Hide abandon button if completed
    document.getElementById('pr-abandon').style.display = s === 'completed' ? 'none' : 'block';

    const stepDiv = document.getElementById('pr-steps');
    if (!stepDiv) return;

    // Calculate pending rewards
    let pendingGold = 0;
    const pendingItems = {};
    this.run.killedSteps.forEach(ks => {
      if (ks.rewardsGiven) {
        pendingGold += ks.rewardsGiven.gold || 0;
        (ks.rewardsGiven.items || []).forEach(it => {
          pendingItems[it.itemId] = (pendingItems[it.itemId] || 0) + it.qty;
        });
      }
    });

    const colorAccents = ['var(--acc-magenta)', 'var(--acc-cyan)', 'var(--acc-yellow)', 'var(--acc-orange)', 'var(--acc-purple)'];
    stepDiv.innerHTML = `
      <div class="form-box" style="margin-bottom:15px; border-color:var(--acc-yellow); padding:10px;">
        <div style="font-size:10px; font-weight:bold; color:var(--acc-yellow); text-transform:uppercase; margin-bottom:5px;">💰 Butin accumulé (en attente)</div>
        <div style="font-size:18px; font-weight:900; color:white;">${pendingGold} G</div>
        <div style="display:flex; flex-wrap:wrap; gap:4px; margin-top:5px;">
           ${Object.entries(pendingItems).map(([id, qty]) => `<span class="badge" style="font-size:9px; background:rgba(0,245,212,0.2); border:1px solid var(--acc-cyan);">${qty}x 📦</span>`).join('')}
        </div>
      </div>
    ` + this.steps.map((st, i) => {
      const dead = this.run.killedSteps.some(k => k.bossStepId === st.id);
      const cur = st.order === this.run.currentStep && s === 'active';
      const icon = dead ? '💀' : (st.emoji || BOSS_EMOJIS[i % BOSS_EMOJIS.length]);
      const accent = colorAccents[i % colorAccents.length];
      
      let cardStyle = `border-color: ${accent}; margin-bottom: 12px;`;
      if (cur) cardStyle += ` border-width: 4px; box-shadow: 0 0 20px ${accent}; transform: scale(1.05);`;
      if (dead) cardStyle += ` opacity: 0.6; filter: grayscale(0.5);`;

      return `
        <div class="card" style="padding:16px; cursor:default; ${cardStyle}">
          <div style="display:flex; align-items:center; gap:12px;">
            <div style="width:45px; height:45px; background:rgba(255,255,255,0.1); border-radius:50%; display:flex; align-items:center; justify-content:center; flex-shrink:0; font-size: 24px;">
              ${icon}
            </div>
            <div style="flex:1">
              <div style="font-weight: 900; font-family: var(--font-heading); text-transform: uppercase; color: white;">${st.name}</div>
              <div class="card-sub" style="font-size: 11px; font-weight: bold; color: ${accent};">⚔️ NIV. ${st.difficulty} · 💰 ${st.goldReward}G</div>
            </div>
          </div>
        </div>`;
    }).join('');

    const mapItems = this.steps.map((st, i) => {
      const dead = this.run.killedSteps.some(k => k.bossStepId === st.id);
      const cur = st.order === this.run.currentStep && s === 'active';
      return {
        location: st.location,
        emoji: dead ? '✅' : (st.emoji || BOSS_EMOJIS[i % BOSS_EMOJIS.length]),
        label: st.name,
        fontSize: cur ? 34 : 26,
        circleColor: cur ? 'rgba(255, 58, 242, 0.2)' : dead ? 'rgba(0, 245, 212, 0.1)' : 'rgba(255, 230, 0, 0.1)',
        strokeColor: cur ? 'var(--acc-magenta)' : dead ? 'var(--acc-cyan)' : 'var(--acc-yellow)',
      };
    });
    map.setItems(mapItems);

    if (s === 'completed') this.showVictory();
  },

  movePlayer(lat, lon) {
    this.playerLat = lat;
    this.playerLon = lon;
    map.setPlayer(lat, lon);

    if (this.run.state === 'active') {
      document.getElementById('pr-attack').disabled = false;
    }

    const cur = this.steps.find(st => st.order === this.run.currentStep);
    if (cur?.location) {
      const dist = haversine(lat, lon, cur.location.lat, cur.location.lon);
      const ok = dist <= (cur.location.radiusMeters || 500);
      this.msg(`📍 Distance: ${Math.round(dist)}m / rayon: ${cur.location.radiusMeters}m ${ok ? '✅ A portee !' : '❌ Trop loin'}`, !ok);
    }
  },

  async attempt() {
    if (this.playerLat == null) { this.msg('❌ Clique sur la carte !', true); return; }
    const cur = this.steps.find(st => st.order === this.run.currentStep);
    if (!cur) { this.msg('❌ Pas de boss', true); return; }

    const btn = document.getElementById('pr-attack');
    btn.disabled = true;
    btn.textContent = '⏳ ...';

    try {
      const res = await API.attemptBoss(this.runId, cur.id, this.playerLat, this.playerLon);
      const r = res.data;
      
      this.spawnParticles(cur);
      await this.reloadRun();
      this.loadGold();

      if (r.runCompleted) {
        this.msg('🏆 DONJON NETTOYÉ !');
        this.showVictory(r.rewards); // Pass final boss rewards to popup
      } else {
        let lootMsg = `🔥 VICTOIRE !!! +${r.rewards.gold}💰`;
        if (r.rewards.items?.length) lootMsg += ` (+${r.rewards.items.length} objets)`;
        toast(lootMsg);
        this.msg(lootMsg);
      }
    } catch (e) {
      const m = e.message || JSON.stringify(e);
      if (m.includes('NOT_IN_RANGE')) {
        const dist = haversine(this.playerLat, this.playerLon, cur.location.lat, cur.location.lon);
        this.msg(`📍 NOT_IN_RANGE ! ${Math.round(dist)}m / ${cur.location.radiusMeters}m — rapproche-toi !`, true);
      } else if (m.includes('WRONG_STEP_ORDER')) {
        this.msg('⚠️ WRONG_STEP_ORDER', true);
      } else {
        this.msg('❌ ' + m, true);
      }
    }
    btn.disabled = false;
    btn.textContent = '⚔️ ATTAQUER !';
  },

  spawnParticles(step) {
    if (!step?.location) return;
    const mapDiv = document.getElementById('pr-map');
    const rect = mapDiv.getBoundingClientRect();
    const point = map.map.latLngToContainerPoint([step.location.lat, step.location.lon]);
    const particles = ['💰', '⭐', '✨', '💎', '🎉'];

    for (let i = 0; i < 12; i++) {
      const span = document.createElement('span');
      span.textContent = particles[Math.floor(Math.random() * particles.length)];
      span.style.cssText = `position:fixed; left:${rect.left + point.x}px; top:${rect.top + point.y}px; font-size:20px; pointer-events:none; z-index:200; transition: all .8s ease-out;`;
      document.body.appendChild(span);
      requestAnimationFrame(() => {
        span.style.left = `${rect.left + point.x + (Math.random() - .5) * 200}px`;
        span.style.top = `${rect.top + point.y + (Math.random() - .5) * 200}px`;
        span.style.opacity = '0';
        span.style.transform = 'scale(0.2)';
      });
      setTimeout(() => span.remove(), 900);
    }
  },

  async showVictory(finalRewards = null) {
    const v = document.getElementById('pr-victory');
    const textEl = document.getElementById('pr-victory-text');
    const rewardsEl = document.getElementById('pr-victory-rewards');
    
    let totalGold = 0;
    const totalItems = {};
    
    // Aggregate all loot from run
    this.run.killedSteps.forEach(ks => {
      if (ks.rewardsGiven) {
        totalGold += ks.rewardsGiven.gold || 0;
        (ks.rewardsGiven.items || []).forEach(it => {
          if (!totalItems[it.itemId]) {
            totalItems[it.itemId] = { qty: 0, details: it.itemDetails };
          }
          totalItems[it.itemId].qty += it.qty;
        });
      }
    });

    textEl.innerHTML = `<span style="color:var(--acc-yellow);">${this.run.killedSteps.length} BOSS VAINCUS !</span><br>Bravo, tu as survécu à l'épreuve. ✨`;
    
    rewardsEl.innerHTML = `
      <div style="font-size: 12px; color:white; opacity:0.6; text-transform:uppercase; margin-bottom:10px; font-weight:900;">Butin total accumulé :</div>
      <div style="display:flex; justify-content:space-between; font-size: 24px; font-weight:900; color:var(--acc-yellow); margin-bottom: 20px;">
         <span>💰 OR TOTAL</span>
         <span>${totalGold} G</span>
      </div>
      <div style="display:flex; flex-direction:column; gap:8px;">
         ${Object.entries(totalItems).length ? Object.entries(totalItems).map(([id, data]) => {
           const info = data.details || { name: 'Objet Inconnu', rarity: 'common', type: 'unknown' };
           const accent = info.rarity === 'uncommon' ? 'var(--acc-cyan)' : info.rarity === 'rare' ? 'var(--acc-purple)' : info.rarity === 'epic' ? 'var(--acc-magenta)' : info.rarity === 'legendary' ? 'var(--acc-yellow)' : '#8899aa';

           return `
             <div class="card" style="padding:10px; margin:0; border-color:${accent}; display:flex; align-items:center; gap:12px; background:rgba(0,0,0,0.3);">
               <span style="font-size:24px;">🎁</span>
               <div style="flex:1;">
                  <div style="font-size:10px; opacity:0.5; font-weight:900; text-transform:uppercase;">${info.type}</div>
                  <div style="font-weight:900; font-size:14px; color:white;" title="ID: ${id}">${info.name}</div>
               </div>
               <span class="badge" style="background:${accent}; color:var(--bg); border-radius:4px;">${data.qty}x</span>
             </div>`;
         }).join('') : '<div style="color:#8899aa; font-style:italic; font-size:12px; text-align:center;">Aucun objet trouvé dans les décombres...</div>'}
      </div>
    `;

    v.style.display = 'flex';
  },

  async abandon() {
    const ok = await confirmModal(
      "Abandonner l'aventure ?",
      "💔 ATTENTION : En abandonnant maintenant, tu ne récupéreras que 50% du butin accumulé. Es-tu sûr de vouloir fuir ?",
      "🚩"
    );
    if (!ok) return;

    try {
      await API.abandonRun(this.runId);
      navigate('player-list');
    } catch (e) { this.msg('Erreur: ' + (e.message || ''), true); }
  },

  msg(txt, isError = false) {
    const el = document.getElementById('pr-msg');
    el.innerHTML = txt;
    el.className = `status-msg ${isError ? 'err' : 'ok'}`;
  },
};

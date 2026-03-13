import { API } from '../api.js';
import { state, navigate, toast } from '../app.js';
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
    const killed = this.run.killedSteps.length;
    const emoji = s === 'active' ? '🟢' : s === 'completed' ? '🏆' : '⛔';
    document.getElementById('pr-state').innerHTML =
      `${emoji} ${s.toUpperCase()} · Step ${this.run.currentStep}/${this.steps.length} · ${killed} boss tues`;

    const stepDiv = document.getElementById('pr-steps');
    stepDiv.innerHTML = this.steps.map((st, i) => {
      const dead = this.run.killedSteps.some(k => k.bossStepId === st.id);
      const cur = st.order === this.run.currentStep && s === 'active';
      const icon = dead ? '✅' : BOSS_EMOJIS[i % BOSS_EMOJIS.length];
      const style = cur ? 'border-color:var(--accent);' : dead ? 'opacity:.5;' : '';
      return `<div class="card" style="padding:8px; cursor:default; ${style}">
        <span>${icon} ${st.name}</span>
        <div class="card-sub">⚔️${st.difficulty} · 💰${st.goldReward}g · 📍${st.location?.radiusMeters || 500}m</div>
      </div>`;
    }).join('');

    const mapItems = this.steps.map((st, i) => {
      const dead = this.run.killedSteps.some(k => k.bossStepId === st.id);
      const cur = st.order === this.run.currentStep && s === 'active';
      return {
        location: st.location,
        emoji: dead ? '✅' : BOSS_EMOJIS[i % BOSS_EMOJIS.length],
        label: st.name,
        fontSize: cur ? 34 : 26,
        circleColor: cur ? 'rgba(233,69,96,0.15)' : dead ? 'rgba(46,204,113,0.1)' : 'rgba(85,85,85,0.1)',
        strokeColor: cur ? 'rgba(233,69,96,0.7)' : dead ? 'rgba(46,204,113,0.5)' : 'rgba(85,85,85,0.4)',
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
      let txt = `✅ VICTOIRE ! +${r.rewards.gold}💰`;
      if (r.rewards.items?.length) txt += ' ' + r.rewards.items.map(i => `${i.qty}x 🎁`).join(', ');
      if (r.runCompleted) txt += ' — 🏆 DONJON TERMINE !';
      this.msg(txt);
      this.spawnParticles(cur);
      await this.reloadRun();
      this.loadGold();
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
    btn.textContent = '⚔️ Attaquer le boss !';
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

  showVictory() {
    const v = document.getElementById('pr-victory');
    v.style.display = '';
    document.getElementById('pr-victory-text').textContent = `${this.run.killedSteps.length} boss vaincus — Bravo ! 🎉`;
  },

  async abandon() {
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

import { API } from '../api.js';
import { navigate, toast } from '../app.js';
import { LeafletMap } from '../map.js';

const BOSS_EMOJIS = ['🐉', '💀', '🧙', '👹', '🦇', '🕷️', '👻', '🧟', '🐍', '🦂'];

let map = null;

export const gmEditView = {
  dungeonId: null,
  dungeon: null,
  steps: [],

  async init(data) {
    this.dungeonId = data.dungeonId;
    this.steps = [];
    this.dungeon = null;

    map = new LeafletMap('gme-map');
    map.onRightClick = (lat, lon) => this.addStep(lat, lon);

    document.getElementById('gme-save').onclick = () => this.save();
    document.getElementById('gme-publish').onclick = () => this.publish();
    document.getElementById('gme-back').onclick = () => navigate('gm-list');

    await this.load();
  },

  destroy() {
    if (map) { map.destroy(); map = null; }
  },

  async load() {
    try {
      const res = await API.getDungeon(this.dungeonId);
      this.dungeon = res.data;
      this.steps = this.dungeon.bossSteps || [];
      document.getElementById('gme-title').value = this.dungeon.title || '';
      document.getElementById('gme-desc').value = this.dungeon.description || '';
      document.getElementById('gme-zone').value = this.dungeon.area?.name || '';
      this.renderSteps();
    } catch (e) {
      this.msg('Erreur: ' + (e.message || ''), true);
    }
  },

  renderSteps() {
    document.getElementById('gme-count').textContent = `(${this.steps.length})`;
    const list = document.getElementById('gme-steps');

    if (this.steps.length === 0) {
      list.innerHTML = '<p style="color:#8899aa; font-size:12px; text-align:center; padding:10px;">Clic droit sur la carte pour ajouter un boss</p>';
    } else {
      list.innerHTML = this.steps.map((s, i) => `
        <div class="card" style="cursor:default; padding:8px;">
          <div style="display:flex; justify-content:space-between;">
            <span>${BOSS_EMOJIS[i % BOSS_EMOJIS.length]} ${s.name}</span>
            <span style="font-size:11px; color:#8899aa;">#${s.order}</span>
          </div>
          <div class="card-sub">⚔️ ${s.difficulty} · 💰 ${s.goldReward}g · 📍 ${s.location?.radiusMeters || 500}m</div>
        </div>
      `).join('');
    }

    const items = this.steps.map((s, i) => ({
      location: s.location,
      emoji: BOSS_EMOJIS[i % BOSS_EMOJIS.length],
      label: `${i + 1}`,
      circleColor: 'rgba(233,69,96,0.15)',
      strokeColor: 'rgba(233,69,96,0.7)',
    }));
    map.setItems(items);
    if (this.steps.length > 0) map.fitToItems();
  },

  async addStep(lat, lon) {
    const order = this.steps.length + 1;
    try {
      const res = await API.createStep(this.dungeonId, {
        name: `Boss ${order}`,
        location: { lat, lon, radiusMeters: 500 },
        zoneDescription: '',
        difficulty: Math.min(order * 2, 10),
        goldReward: order * 100,
        lootTable: [],
      });
      this.steps.push(res.data);
      this.renderSteps();
      this.msg(`Boss ${order} ajoute !`);
    } catch (e) {
      this.msg('Erreur: ' + (e.message || ''), true);
    }
  },

  async save() {
    try {
      const bounds = map.getBounds();
      await API.updateDungeon(this.dungeonId, {
        title: document.getElementById('gme-title').value,
        description: document.getElementById('gme-desc').value,
        area: {
          name: document.getElementById('gme-zone').value,
          boundingBox: bounds,
        },
      });
      this.msg('Sauvegarde !');
    } catch (e) {
      this.msg('Erreur: ' + (e.message || ''), true);
    }
  },

  async publish() {
    try {
      await API.publishDungeon(this.dungeonId);
      toast('Donjon publie !');
      this.msg('Publie ! 🎉');
    } catch (e) {
      this.msg('Erreur: ' + (e.message || ''), true);
    }
  },

  msg(txt, isError = false) {
    const el = document.getElementById('gme-msg');
    el.textContent = txt;
    el.className = `status-msg ${isError ? 'err' : 'ok'}`;
    if (!isError) setTimeout(() => { el.textContent = ''; }, 3000);
  },
};

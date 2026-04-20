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
    map.onMove = (id, lat, lon) => this.updateStepPosition(id, lat, lon);

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
      
      const publishBtn = document.getElementById('gme-publish');
      if (this.dungeon.status === 'published') {
        publishBtn.style.display = 'none';
      } else {
        publishBtn.style.display = 'inline-block';
      }
      
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
      const colorAccents = ['var(--acc-magenta)', 'var(--acc-cyan)', 'var(--acc-yellow)', 'var(--acc-orange)', 'var(--acc-purple)'];
      list.innerHTML = this.steps.map((s, i) => {
        const accent = colorAccents[i % colorAccents.length];
        return `
          <div class="card" style="cursor:default; padding:20px; margin-bottom:20px; border-color: ${accent}; border-style: ${i % 2 === 0 ? 'solid' : 'dashed'}; transform: rotate(${i % 2 === 0 ? '1deg' : '-1deg'});">
            <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:12px;">
              <div style="display:flex; gap:12px; align-items:center; flex: 1;">
                <input type="text" value="${s.emoji || BOSS_EMOJIS[i % BOSS_EMOJIS.length]}" 
                  style="width:45px; height:45px; background:rgba(255,255,255,0.1); border:none; text-align:center; font-size:24px; border-radius: 50%; line-height:45px; padding:0; flex-shrink:0; cursor: pointer;"
                  onblur="gmEditView.updateStepFieldValue('${s.id}', 'emoji', this.value)">
                <input type="text" value="${s.name}" 
                  style="background:transparent; border:none; color:white; font-weight:900; flex:1; min-width:0; font-size: 18px; font-family: var(--font-heading); text-transform: uppercase; outline: none;"
                  onblur="gmEditView.updateStepFieldValue('${s.id}', 'name', this.value)">
              </div>
              <button class="btn btn-sm" style="background:var(--acc-orange); border: none; padding:0; border-radius: 50%; width: 40px; height: 40px; display: flex; align-items: center; justify-content: center; flex-shrink:0; margin-left:8px;" onclick="gmEditView.deleteStep('${s.id}')">🗑️</button>
            </div>
            
            <div class="form-box" style="padding: 12px; border-radius: var(--radius-btn); border-width: 2px; margin-bottom: 12px; border-color: var(--acc-cyan);">
              <div style="display:flex; gap:10px; font-size:14px; align-items:center; color:white; font-weight: bold;">
                 <span style="text-transform: uppercase; font-size: 12px; color: var(--acc-cyan);">Rayon d'action:</span>
                 <input type="number" value="${s.location?.radiusMeters || 500}" 
                   style="width:80px; background:transparent; border:none; color:var(--acc-yellow); padding:0; font-weight: 900; font-size: 18px;"
                   onblur="gmEditView.updateStepFieldValue('${s.id}', 'radiusMeters', this.value)">
                 <span style="color: var(--acc-cyan);">m</span>
              </div>
            </div>

            <div style="display: flex; justify-content: space-between; font-weight: 900; font-size: 14px; text-transform: uppercase;">
              <span style="color: var(--acc-magenta);">⚔️ Niv. ${s.difficulty}</span>
              <span style="color: var(--acc-yellow);">💰 ${s.goldReward}G</span>
            </div>
          </div>
        `;
      }).join('');
    }

    const items = this.steps.map((s, i) => ({
      id: s.id,
      location: s.location,
      emoji: s.emoji || BOSS_EMOJIS[i % BOSS_EMOJIS.length],
      draggable: true,
      label: `${i + 1}`,
      circleColor: 'rgba(0, 245, 212, 0.1)',
      strokeColor: 'rgba(255, 58, 242, 0.8)',
    }));
    map.setItems(items);
    // Only fit to items on first load or when explicitly needed
  },

  async updateStepPosition(stepId, lat, lon) {
    try {
      const step = this.steps.find(s => s.id === stepId);
      if (!step) return;
      await API.updateStep(this.dungeonId, stepId, {
        ...step,
        location: { ...step.location, lat, lon }
      });
      step.location.lat = lat;
      step.location.lon = lon;
      this.msg('Position mise à jour');
      this.renderSteps();
    } catch (e) {
      this.msg('Erreur de déplacement: ' + (e.message || ''), true);
    }
  },

  async updateStepFieldValue(stepId, field, value) {
    try {
      const step = this.steps.find(s => s.id === stepId);
      if (!step) return;
      
      // Don't update if value hasn't changed
      let newValue = value;
      if (field === 'radiusMeters') {
        newValue = parseFloat(value);
        if (step.location.radiusMeters === newValue) return;
      } else {
        if (step[field] === value) return;
      }

      this.msg(`Mise à jour ${field}...`);
      
      const updateData = { ...step };
      if (field === 'radiusMeters') {
        updateData.location = { ...step.location, radiusMeters: newValue };
      } else {
        updateData[field] = value;
      }

      console.log('UPDATING BOSS STEP:', { dungeonId: this.dungeonId, stepId, updateData });
      await API.updateStep(this.dungeonId, stepId, updateData);
      
      if (field === 'radiusMeters') {
        step.location.radiusMeters = newValue;
      } else {
        step[field] = value;
      }
      
      this.msg('Boss mis à jour');
      // Force a full refresh of map items with updated data
      this.renderSteps();
    } catch (e) {
      this.msg('Erreur: ' + (e.message || ''), true);
    }
  },

  async deleteStep(stepId) {
    if (!confirm('Supprimer ce boss ?')) return;
    try {
      await API.deleteStep(this.dungeonId, stepId);
      this.steps = this.steps.filter(s => s.id !== stepId);
      this.renderSteps();
      this.msg('Boss supprimé');
    } catch (e) {
      this.msg('Erreur de suppression: ' + (e.message || ''), true);
    }
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

import { API } from '../api.js';
import { navigate, toast } from '../app.js';
import { LeafletMap } from '../map.js';

const BOSS_EMOJIS = ['🐉', '💀', '🧙', '👹', '🦇', '🕷️', '👻', '🧟', '🐍', '🦂'];

let map = null;

export const gmEditView = {
  dungeonId: null,
  dungeon: null,
  steps: [],
  availableItems: [],

  async init(data) {
    this.dungeonId = data.dungeonId;
    this.steps = [];
    this.dungeon = null;
    this.availableItems = [];

    map = new LeafletMap('gme-map');
    map.onRightClick = (lat, lon) => this.addStep(lat, lon);
    map.onMove = (id, lat, lon) => this.updateStepPosition(id, lat, lon);

    document.getElementById('gme-save').onclick = () => this.save();
    document.getElementById('gme-publish').onclick = () => this.publish();
    document.getElementById('gme-back').onclick = () => navigate('gm-list');
    document.getElementById('gme-show-map').onclick = () => {
      this.toggleMap(true);
      if (this.steps.length > 0) map.fitToItems();
    };
    document.getElementById('gme-hide-map').onclick = () => this.toggleMap(false);

    // DELEGATION FOR BOSS STEPS (CSP COMPLIANCE)
    const list = document.getElementById('gme-steps');
    list.onclick = (e) => this.handleStepClick(e);
    list.oninput = (e) => this.handleStepInput(e); // Also handle typing
    list.onfocusout = (e) => this.handleStepInput(e); // Replace onblur

    await this.fetchItems();
    await this.load();
  },

  toggleMap(visible, center = null) {
    const modal = document.getElementById('gme-map-modal');
    modal.style.display = visible ? 'flex' : 'none';
    if (visible && map) {
      setTimeout(() => {
        map.map.invalidateSize();
        if (center) {
          map.map.flyTo(center, 16);
        }
      }, 150);
    }
  },

  async fetchItems() {
    try {
      const res = await API.getItems();
      this.availableItems = res.data || [];
    } catch (e) {
      console.error('fetchItems', e);
    }
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
      document.getElementById('gme-completion-gold').value = this.dungeon.completionGoldReward || 0;
      
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
      list.innerHTML = '<div style="grid-column: 1/-1; padding:40px; text-align:center; background:rgba(255,255,255,0.05); border-radius:20px; border:2px dashed var(--acc-purple); color:#8899aa;">Aucun boss pour le moment. Cliquez sur "Gérer la carte" pour en ajouter.</div>';
    } else {
      const colorAccents = ['var(--acc-magenta)', 'var(--acc-cyan)', 'var(--acc-yellow)', 'var(--acc-orange)', 'var(--acc-purple)'];
      list.innerHTML = this.steps.map((s, i) => {
        const accent = colorAccents[i % colorAccents.length];
        
        const lootHtml = (s.lootTable || []).map((l, li) => {
          const itemName = this.availableItems.find(it => it.id === l.itemId)?.name || 'Objet inconnu';
          return `
            <div style="display:flex; gap:8px; align-items:center; font-size:12px; background:rgba(255,255,255,0.05); padding:6px; margin-bottom:4px; border-radius:8px; border: 1px solid rgba(255,255,255,0.1);">
              <span style="flex:1; overflow:hidden; white-space:nowrap; text-overflow:ellipsis;" title="${itemName}">${itemName}</span>
              <span style="color:var(--acc-cyan); font-weight:bold;">${Math.round(l.dropRate * 100)}%</span>
              <button class="btn btn-sm gme-btn-remove-loot" data-step-id="${s.id}" data-index="${li}" style="padding:0 6px; height:22px; line-height:22px; background:transparent; box-shadow:none; border:none;">❌</button>
            </div>
          `;
        }).join('');

        const itemOptions = this.availableItems.map(it => `<option value="${it.id}">${it.name}</option>`).join('');

        return `
          <div class="card" style="padding:24px; border-color: ${accent}; border-radius: 20px; transition: 0.2s;">
            <div style="display:flex; justify-content:space-between; align-items:flex-start; margin-bottom:16px;">
              <div style="display:flex; gap:16px; align-items:center; flex: 1;">
                <input type="text" value="${s.emoji || BOSS_EMOJIS[i % BOSS_EMOJIS.length]}" 
                  class="gme-input" data-step-id="${s.id}" data-field="emoji"
                  style="width:50px; height:50px; padding:0; background:rgba(255,255,255,0.1); border:none; text-align:center; font-size:28px; border-radius: 12px; line-height:50px; cursor: pointer;">
                <div style="flex:1">
                    <input type="text" value="${s.name}" 
                      class="gme-input" data-step-id="${s.id}" data-field="name"
                      style="background:transparent; border:none; color:white; font-weight:900; width:100%; font-size: 20px; text-transform: uppercase; padding:0; margin:0;">
                    <div style="font-size:10px; color:rgba(255,255,255,0.3); font-family:monospace;">ID: ${s.id}</div>
                    <div style="font-size:12px; color:${accent}; font-weight:bold;">Gardien #${i+1}</div>
                </div>
              </div>
              <button class="btn btn-sm gme-btn-delete" data-step-id="${s.id}" style="background:var(--acc-orange); border: none; padding:0; border-radius: 50%; width: 36px; height: 36px; display: flex; align-items: center; justify-content: center;">🗑️</button>
            </div>
            
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin-bottom: 20px;">
                <div class="form-box gme-btn-map" data-step-id="${s.id}" data-lat="${s.location.lat}" data-lon="${s.location.lon}" style="padding: 12px; border-color:${accent}; cursor:pointer; background:rgba(255,255,255,0.02);">
                    <div style="text-transform: uppercase; font-size: 10px; color:${accent}; margin-bottom: 4px; font-weight:900;">📍 Localisation</div>
                    <div style="font-size: 13px; color:white;">Modifier sur la carte</div>
                </div>
                <div class="form-box" style="padding: 12px; background:rgba(255,255,255,0.02); border-color:var(--acc-yellow);">
                    <div style="text-transform: uppercase; font-size: 10px; color:var(--acc-yellow); margin-bottom: 4px; font-weight:900;">💰 Or du Boss</div>
                    <input type="number" value="${s.goldReward}" 
                        class="gme-input" data-step-id="${s.id}" data-field="goldReward"
                        style="width:100%; background:transparent; border:none; color:white; font-weight: 900; font-size: 18px; padding:0;">
                </div>
            </div>

            <div class="form-box" style="padding: 15px; border-color:rgba(0, 245, 212, 0.3); background:rgba(0, 245, 212, 0.05); border-radius:15px;">
                <div style="display:flex; justify-content:space-between; align-items:center;">
                    <span style="text-transform: uppercase; font-size: 11px; color: var(--acc-cyan); font-weight: 900;">📏 Rayon d'action (m)</span>
                    <input type="number" value="${s.location?.radiusMeters || 500}" 
                        class="gme-input" data-step-id="${s.id}" data-field="radiusMeters"
                        style="width:100px; background:rgba(255,255,255,0.1); border:2px solid var(--acc-cyan); border-radius:10px; color:white; font-weight: 900; font-size: 18px; text-align:center; padding:4px;">
                </div>
            </div>

            <div class="form-box" style="padding: 15px; border-color:rgba(255, 255, 255, 0.1); background:rgba(255, 255, 255, 0.05); margin-top:12px;">
                <div style="text-transform: uppercase; font-size: 11px; font-weight: 900; color:var(--acc-cyan); margin-bottom:12px;">🎁 Table de Butin</div>
                <div id="loot-list-${s.id}">${lootHtml}</div>
                <div style="display:flex; gap:8px; margin-top:12px; border-top:1px solid rgba(255,255,255,0.1); padding-top:12px;">
                    <select id="loot-add-id-${s.id}" style="flex:1; background:var(--bg); color:white; border: 2px solid var(--acc-purple); font-size:12px; padding:8px 15px; border-radius:15px;">
                        <option value="">+ Ajouter...</option>
                        ${itemOptions}
                    </select>
                    <button class="btn btn-sm gme-btn-add-loot" data-step-id="${s.id}" style="padding:0 15px; font-size:11px; box-shadow:none; background:var(--acc-cyan); color:var(--bg); border:none;">ADD</button>
                </div>
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
  },

  async addLoot(stepId) {
    const itemId = document.getElementById(`loot-add-id-${stepId}`).value;
    if (!itemId) return;
    
    const step = this.steps.find(s => s.id === stepId);
    if (!step) return;
    
    const lootTable = step.lootTable || [];
    lootTable.push({ itemId, dropRate: 1.0, minQty: 1, maxQty: 1 });
    
    await this.updateStepFieldValue(stepId, 'lootTable', lootTable);
  },

  async removeLoot(stepId, index) {
    const step = this.steps.find(s => s.id === stepId);
    if (!step) return;
    
    const lootTable = step.lootTable || [];
    lootTable.splice(index, 1);
    
    await this.updateStepFieldValue(stepId, 'lootTable', lootTable);
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
    if (this._lastVal?.[stepId]?.[field] === value) return; // Prevent double trigger
    if (!this._lastVal) this._lastVal = {};
    if (!this._lastVal[stepId]) this._lastVal[stepId] = {};
    this._lastVal[stepId][field] = value;

    console.log(`[SYNC] updateStepFieldValue triggered for field: ${field}, value: ${value}, stepId: ${stepId}`);
    try {
      const step = this.steps.find(s => s.id === stepId);
      if (!step) {
          console.error(`[SYNC] ERROR: Step with ID ${stepId} not found in this.steps! Available IDs:`, this.steps.map(s => s.id));
          return;
      }
      
      let newValue = value;
      if (field === 'radiusMeters') {
        newValue = parseFloat(value);
      } else if (field === 'goldReward' || field === 'difficulty') {
        newValue = parseInt(value) || 0;
      }

      // 1. UPDATE LOCAL STATE IMMEDIATELY (Prevent race conditions)
      if (field === 'radiusMeters') {
        step.location.radiusMeters = newValue;
      } else {
        step[field] = newValue;
      }

      this.msg(`MAJ ${field}...`);
      
      // 2. SEND CURRENT (LOCAL) STATE TO SERVER
      console.log(`[SYNC] Sending update for boss ${stepId}, field ${field}:`, newValue);
      const res = await API.updateStep(this.dungeonId, stepId, step);
      console.log(`[SYNC] Server response for ${field}:`, res.message || 'OK');
      
      this.msg('Mis à jour');
      // this.renderSteps(); // Removed to avoid losing focus during typing if delegation is used
    } catch (e) {
      this.msg('Erreur: ' + (e.message || ''), true);
      console.error(e);
      await this.load(); 
    }
  },

  handleStepClick(e) {
    const target = e.target.closest('[data-step-id]');
    if (!target) return;
    const stepId = target.dataset.stepId;

    if (e.target.closest('.gme-btn-delete')) {
      this.deleteStep(stepId);
    } else if (e.target.closest('.gme-btn-add-loot')) {
      this.addLoot(stepId);
    } else if (e.target.closest('.gme-btn-map')) {
        const lat = parseFloat(target.dataset.lat);
        const lon = parseFloat(target.dataset.lon);
        this.toggleMap(true, [lat, lon]);
    } else if (e.target.closest('.gme-btn-remove-loot')) {
        const index = parseInt(e.target.closest('.gme-btn-remove-loot').dataset.index);
        this.removeLoot(stepId, index);
    }
  },

  handleStepInput(e) {
    const target = e.target;
    if (!target.classList.contains('gme-input')) return;
    
    const stepId = target.dataset.stepId;
    const field = target.dataset.field;
    const value = target.value;
    
    this.updateStepFieldValue(stepId, field, value);
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
      const dungeonData = {
        title: document.getElementById('gme-title').value,
        description: document.getElementById('gme-desc').value,
        completionGoldReward: parseInt(document.getElementById('gme-completion-gold').value) || 0,
        area: {
          name: document.getElementById('gme-zone').value,
          boundingBox: bounds,
        },
      };

      await API.updateDungeonFull(this.dungeonId, {
        dungeon: dungeonData,
        bossSteps: this.steps
      });

      this.msg('Sauvegarde globale réussie !');
      return true;
    } catch (e) {
      this.msg('Erreur: ' + (e.message || ''), true);
      return false;
    }
  },

  async publish() {
    try {
      this.msg('Sauvegarde avant publication...');
      const saved = await this.save(); 
      if (!saved) {
        throw new Error("La sauvegarde a échoué. Publication annulée.");
      }

      await API.publishDungeon(this.dungeonId);
      toast('Donjon publie ! 🚀✨');
      this.msg('Publie ! 🎉');
      await this.load();
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

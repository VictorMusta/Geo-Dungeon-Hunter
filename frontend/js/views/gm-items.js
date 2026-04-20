import { API } from '../api.js';
import { navigate, toast } from '../app.js';

export const gmItemsView = {
  items: [],
  editingId: null,

  async init() {
    this.editingId = null;
    this.resetForm();
    
    document.getElementById('gmi-back').onclick = () => navigate('gm-list');
    document.getElementById('gmi-save').onclick = () => this.save();
    document.getElementById('gmi-cancel').onclick = () => this.cancelEdit();
    
    // DELEGATION FOR ITEMS LIST (CSP COMPLIANCE)
    const list = document.getElementById('gmi-list');
    list.onclick = (e) => this.handleItemClick(e);

    await this.load();
  },

  async load() {
    try {
      const res = await API.getItems();
      this.items = res.data || [];
      this.render();
    } catch (e) {
      toast('Erreur lors du chargement des objets', true);
    }
  },

  resetForm() {
    this.editingId = null;
    document.getElementById('gmi-form-title').textContent = "🆕 Forger un nouvel objet";
    document.getElementById('gmi-save').textContent = "🛠️ Forger l'objet";
    document.getElementById('gmi-cancel').style.display = 'none';
    
    document.getElementById('gmi-id').value = '';
    document.getElementById('gmi-name').value = '';
    document.getElementById('gmi-type').value = 'weapon';
    document.getElementById('gmi-rarity').value = 'common';
    document.getElementById('gmi-value').value = '100';
    document.getElementById('gmi-desc').value = '';
    document.getElementById('gmi-tradable').checked = true;
  },

  async save() {
    const itemData = {
      name: document.getElementById('gmi-name').value,
      type: document.getElementById('gmi-type').value,
      rarity: document.getElementById('gmi-rarity').value,
      baseValue: parseInt(document.getElementById('gmi-value').value),
      description: document.getElementById('gmi-desc').value,
      tradable: document.getElementById('gmi-tradable').checked,
    };

    if (!itemData.name) return toast('Le nom est requis', true);

    try {
      if (this.editingId) {
        await API.updateItem(this.editingId, itemData);
        toast('Objet reforgé avec succès ! ✨');
      } else {
        await API.createItem(itemData);
        toast('Objet forgé avec succès ! ✨');
      }
      this.resetForm();
      await this.load();
    } catch (e) {
      toast('Erreur lors de la forge : ' + (e.message || ''), true);
    }
  },

  edit(item) {
    this.editingId = item.id;
    document.getElementById('gmi-form-title').textContent = "⚒️ Modifier l'objet";
    document.getElementById('gmi-save').textContent = "💾 Enregistrer les modifications";
    document.getElementById('gmi-cancel').style.display = 'inline-block';

    document.getElementById('gmi-id').value = item.id;
    document.getElementById('gmi-name').value = item.name;
    document.getElementById('gmi-type').value = item.type;
    document.getElementById('gmi-rarity').value = item.rarity;
    document.getElementById('gmi-value').value = item.baseValue;
    document.getElementById('gmi-desc').value = item.description || '';
    document.getElementById('gmi-tradable').checked = item.tradable;
    
    document.querySelector('#view-gm-items .panel-page').scrollTop = 0;
  },

  cancelEdit() {
    this.resetForm();
  },

  async delete(id) {
    if (!confirm('Es-tu sûr de vouloir détruire cet objet ? Cette action est irréversible.')) return;
    try {
      await API.deleteItem(id);
      toast('Objet détruit ! 🔥');
      await this.load();
    } catch (e) {
      toast('Erreur lors de la destruction', true);
    }
  },

  render() {
    const list = document.getElementById('gmi-list');
    if (this.items.length === 0) {
      list.innerHTML = `<div style="padding:40px; text-align:center; opacity:0.5;">Le grimoire est vide...</div>`;
      return;
    }

    const rarityColors = {
      common: '#8899aa',
      uncommon: '#2ecc71',
      rare: '#3498db',
      epic: '#9b59b6',
      legendary: '#f1c40f'
    };

    list.innerHTML = this.items.map(it => `
      <div class="card" style="padding:15px; border-left: 5px solid ${rarityColors[it.rarity] || '#fff'};">
        <div style="display:flex; justify-content:space-between; align-items:center;">
          <div>
            <div style="display:flex; align-items:center; gap:8px;">
               <span style="font-weight:900; color:white; font-size:16px;">${it.name}</span>
               <span style="font-size:10px; text-transform:uppercase; color:${rarityColors[it.rarity]}; font-weight:bold;">${it.rarity}</span>
            </div>
            <div style="font-size:11px; opacity:0.7;">${it.type} • ${it.baseValue} Or • ${it.tradable ? 'Échangeable' : 'Lié'}</div>
          </div>
          <div style="display:flex; gap:8px;">
            <button class="btn btn-sm gmi-btn-edit" data-id="${it.id}" style="padding:0 12px; height:30px; font-size:12px; background:rgba(255,255,255,0.1); box-shadow:none; border:none;">✎</button>
            <button class="btn btn-sm gmi-btn-delete" data-id="${it.id}" style="padding:0 12px; height:30px; font-size:12px; background:rgba(233, 69, 96, 0.2); color:#e94560; box-shadow:none; border:none;">🗑️</button>
          </div>
        </div>
      </div>
    `).join('');
  },

  handleItemClick(e) {
    const btn = e.target.closest('button');
    if (!btn) return;

    const id = btn.dataset.id;
    if (btn.classList.contains('gmi-btn-edit')) {
      const item = this.items.find(it => it.id === id);
      if (item) this.edit(item);
    } else if (btn.classList.contains('gmi-btn-delete')) {
      this.delete(id);
    }
  }
};

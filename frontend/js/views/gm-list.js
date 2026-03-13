import { API } from '../api.js';
import { state, navigate } from '../app.js';

export const gmListView = {
  async init() {
    document.getElementById('gml-back').onclick = () => navigate('menu');
    document.getElementById('gml-new').onclick = () => this.createDungeon();
    await this.load();
  },

  async load() {
    const list = document.getElementById('gml-list');
    list.innerHTML = '<p style="color:#8899aa; text-align:center; padding:40px;">Chargement...</p>';
    try {
      const res = await API.getDungeons();
      const dungeons = res.data || [];
      if (dungeons.length === 0) {
        list.innerHTML = '<p style="color:#8899aa; text-align:center; padding:40px;">Aucun donjon. Cree le premier !</p>';
        return;
      }
      list.innerHTML = dungeons.map(d => `
        <div class="card" data-id="${d.id}">
          <div style="display:flex; justify-content:space-between; align-items:center;">
            <span class="card-title">🏰 ${d.title}</span>
            <span class="badge badge-${d.status}">${d.status}</span>
          </div>
          <div class="card-sub">${d.description || 'Pas de description'}</div>
          <div class="card-sub">📍 ${d.area?.name || '—'}</div>
        </div>
      `).join('');
      list.querySelectorAll('.card').forEach(c =>
        c.addEventListener('click', () => navigate('gm-edit', { dungeonId: c.dataset.id }))
      );
    } catch (e) {
      list.innerHTML = `<p style="color:#e94560; text-align:center; padding:40px;">Erreur: ${e.message || 'API indisponible'}</p>`;
    }
  },

  async createDungeon() {
    try {
      const res = await API.createDungeon({
        title: 'Nouveau Donjon',
        description: 'A configurer...',
        createdBy: state.playerId,
        area: { name: 'Paris' },
      });
      navigate('gm-edit', { dungeonId: res.data.id });
    } catch (e) {
      console.error(e);
    }
  }
};

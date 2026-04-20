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
      const res = await API.getMJDungeons(state.playerId);
      const dungeons = res.data || [];
      if (dungeons.length === 0) {
        list.innerHTML = '<p style="color:#8899aa; text-align:center; padding:40px;">Aucun donjon. Cree le premier !</p>';
        return;
      }
      const colorAccents = ['var(--acc-magenta)', 'var(--acc-cyan)', 'var(--acc-yellow)', 'var(--acc-orange)', 'var(--acc-purple)'];
      list.innerHTML = dungeons.map((d, i) => {
        const accent = colorAccents[i % colorAccents.length];
        const rotate = (i % 2 === 0 ? 'rotate(1deg)' : 'rotate(-1deg)');
        const offset = (i % 2 === 0 ? 'translateY(10px)' : 'translateY(-10px)');
        
        return `
          <div class="card" data-id="${d.id}" style="border-color: ${accent}; transform: ${rotate} ${offset}; margin-bottom: 30px;">
            <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom: 12px;">
              <span class="card-title shadow-magenta" style="font-size: 20px; font-weight: 900;">🏰 ${d.title}</span>
              <span class="badge badge-${d.status}">${d.status}</span>
            </div>
            <div class="card-sub" style="font-weight: bold; color: white; opacity: 0.9; margin-bottom: 8px;">${d.description || 'Pas de description'}</div>
            <div class="card-sub" style="color: var(--acc-cyan); font-weight: 900; letter-spacing: 1px;">📍 ${d.area?.name || '—'}</div>
            <div style="position: absolute; top: -10px; right: -10px; font-size: 40px; transform: rotate(15deg); opacity: 0.1; pointer-events: none;">✨</div>
          </div>
        `;
      }).join('');
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

import { API } from '../api.js';
import { state, navigate } from '../app.js';

export const playerListView = {
  async init() {
    document.getElementById('pl-back').onclick = () => navigate('menu');
    document.getElementById('pl-runs').onclick = () => this.toggleRuns();
    document.getElementById('pl-runs-panel').style.display = 'none';
    await this.load();
  },

  async load() {
    const list = document.getElementById('pl-list');
    list.innerHTML = '<p style="color:#8899aa; text-align:center; padding:40px;">Chargement...</p>';
    try {
      const res = await API.getDungeons();
      const dungeons = res.data || [];
      if (dungeons.length === 0) {
        list.innerHTML = '<p style="color:#8899aa; text-align:center; padding:40px;">Aucun donjon publie</p>';
        return;
      }
      list.innerHTML = dungeons.map(d => `
        <div class="card" data-id="${d.id}">
          <div style="display:flex; justify-content:space-between; align-items:center;">
            <span class="card-title">🏰 ${d.title}</span>
          </div>
          <div class="card-sub">${d.description || ''}</div>
          <div class="card-sub">📍 ${d.area?.name || '?'}</div>
          <button class="btn btn-sm play-btn" data-id="${d.id}" style="margin-top:8px;">🎮 Jouer</button>
        </div>
      `).join('');

      list.querySelectorAll('.play-btn').forEach(btn =>
        btn.addEventListener('click', (e) => { e.stopPropagation(); this.startRun(btn.dataset.id); })
      );
    } catch (e) {
      list.innerHTML = `<p style="color:#e94560; text-align:center; padding:40px;">Erreur: ${e.message || 'API indisponible'}</p>`;
    }
  },

  async startRun(dungeonId) {
    try {
      const res = await API.createRun(dungeonId, state.playerId);
      navigate('player-run', { runId: res.data.id, dungeonId });
    } catch (e) {
      if (e.message?.includes('already has an active run')) {
        try {
          const runsRes = await API.getRuns(state.playerId);
          const active = (runsRes.data || []).find(r => r.dungeonId === dungeonId && r.state === 'active');
          if (active) { navigate('player-run', { runId: active.id, dungeonId }); return; }
        } catch (_) {}
      }
      console.error('Start run error', e);
    }
  },

  async toggleRuns() {
    const panel = document.getElementById('pl-runs-panel');
    if (panel.style.display === 'none') {
      panel.style.display = '';
      const list = document.getElementById('pl-runs-list');
      try {
        const res = await API.getRuns(state.playerId);
        const runs = res.data || [];
        if (runs.length === 0) { list.innerHTML = '<p class="card-sub" style="padding:8px;">Aucun run</p>'; return; }
        list.innerHTML = runs.map(r => `
          <div class="card" style="padding:10px;">
            <div style="display:flex; justify-content:space-between;">
              <span>${r.state === 'active' ? '🟢' : r.state === 'completed' ? '✅' : '⛔'} Run</span>
              <span class="badge badge-${r.state}">${r.state}</span>
            </div>
            <div class="card-sub">Step ${r.currentStep} · ${r.killedSteps?.length || 0} boss tues</div>
            ${r.state === 'active' ? `<button class="btn btn-sm resume-btn" data-rid="${r.id}" data-did="${r.dungeonId}" style="margin-top:6px;">Reprendre</button>` : ''}
          </div>
        `).join('');
        list.querySelectorAll('.resume-btn').forEach(btn =>
          btn.addEventListener('click', () => navigate('player-run', { runId: btn.dataset.rid, dungeonId: btn.dataset.did }))
        );
      } catch (e) {
        list.innerHTML = '<p style="color:#e94560;">Erreur</p>';
      }
    } else {
      panel.style.display = 'none';
    }
  }
};

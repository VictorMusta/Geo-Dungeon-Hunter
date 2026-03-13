import { API } from '../api.js';
import { navigate } from '../app.js';

const MEDALS = ['🥇', '🥈', '🥉'];

export const leaderboardView = {
  tab: 'completions',

  init() {
    document.getElementById('lb-back').onclick = () => navigate('menu');

    document.querySelectorAll('#view-leaderboard .tab').forEach(t => {
      t.addEventListener('click', () => {
        document.querySelectorAll('#view-leaderboard .tab').forEach(b => b.classList.remove('active'));
        t.classList.add('active');
        this.tab = t.dataset.key;
        this.load();
      });
    });

    this.load();
  },

  async load() {
    const container = document.getElementById('lb-table');
    container.innerHTML = '<p style="color:#8899aa; text-align:center; padding:40px;">Chargement...</p>';
    try {
      const res = await API.getLeaderboard(this.tab, 20);
      const entries = res.data || [];
      if (entries.length === 0) {
        container.innerHTML = '<p style="color:#8899aa; text-align:center; padding:40px;">Aucune donnee</p>';
        return;
      }
      const fmt = (e) => {
        if (this.tab === 'gold') return `💰 ${e.gold || e.score || 0}`;
        if (this.tab === 'speed') return `⚡ ${((e.bestTime || e.score || 0) / 1000).toFixed(1)}s`;
        return `🏰 ${e.completions || e.score || 0}`;
      };
      container.innerHTML = `
        <table class="lb-table">
          <thead><tr><th>#</th><th>Joueur</th><th class="score">Score</th></tr></thead>
          <tbody>
            ${entries.map((e, i) => `
              <tr style="${i < 3 ? 'font-weight:700;' : ''}">
                <td>${i < 3 ? MEDALS[i] : i + 1}</td>
                <td>${e.displayName || e.playerId?.slice(0, 8) || '???'}</td>
                <td class="score">${fmt(e)}</td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      `;
    } catch (e) {
      container.innerHTML = `<p style="color:#e94560; text-align:center; padding:40px;">Erreur: ${e.message || 'API indisponible'}</p>`;
    }
  },
};

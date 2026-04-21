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
        return `🏰 ${e.completions || e.score || 0}`;
      };
      const colorAccents = ['var(--acc-magenta)', 'var(--acc-cyan)', 'var(--acc-yellow)'];
      container.innerHTML = `
        <table class="lb-table">
          <thead>
            <tr style="background: var(--acc-purple); color: white;">
              <th style="border: none; padding: 16px;">RANG</th>
              <th style="border: none; padding: 16px;">HÉROS</th>
              <th class="score" style="border: none; padding: 16px;">RÉSULTAT</th>
            </tr>
          </thead>
          <tbody>
            ${entries.map((e, i) => {
              const accent = i < 3 ? colorAccents[i] : 'transparent';
              const rowBg = i % 2 === 0 ? 'rgba(255,255,255,0.03)' : 'transparent';
              return `
              <tr style="background: ${rowBg}; border-bottom: 2px dashed rgba(255,255,255,0.1);">
                <td style="padding: 16px; font-weight: 900; color: ${accent || 'white'}; font-size: ${i < 3 ? '24px' : '18px'};">
                  ${i < 3 ? MEDALS[i] : (i + 1)}
                </td>
                <td style="padding: 16px; font-weight: 700; font-family: var(--font-heading); text-transform: uppercase;">
                  ${e.displayName || e.playerId?.slice(0, 8) || '???'}
                </td>
                <td class="score" style="padding: 16px; font-family: var(--font-display); color: ${i < 3 ? accent : 'var(--acc-cyan)'};">
                  ${fmt(e)}
                </td>
              </tr>
            `;}).join('')}
          </tbody>
        </table>
      `;
    } catch (e) {
      container.innerHTML = `<p style="color:#e94560; text-align:center; padding:40px;">Erreur: ${e.message || 'API indisponible'}</p>`;
    }
  },
};

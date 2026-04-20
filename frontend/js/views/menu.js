import { API } from '../api.js';
import { state, navigate } from '../app.js';

export const menuView = {
  init() {
    const nameInput = document.getElementById('menu-name');
    const btnCreate = document.getElementById('menu-create');
    const statusMsg = document.getElementById('menu-status');
    const roles = document.getElementById('menu-roles');

    if (state.playerId) {
      nameInput.value = state.playerName || '';
      nameInput.disabled = true;
      btnCreate.textContent = '✅ Connecte';
      btnCreate.disabled = true;
      statusMsg.textContent = `ID: ${state.playerId}`;
      statusMsg.className = 'status-msg ok';
      roles.style.display = '';
    } else {
      nameInput.value = '';
      nameInput.disabled = false;
      btnCreate.textContent = 'Creer mon personnage';
      btnCreate.disabled = false;
      statusMsg.textContent = '';
      roles.style.display = 'none';
    }

    btnCreate.onclick = async () => {
      const name = nameInput.value.trim();
      if (!name) { statusMsg.textContent = '❌ Entre un pseudo !'; statusMsg.className = 'status-msg err'; return; }
      btnCreate.disabled = true;
      btnCreate.textContent = '...';
      try {
        const res = await API.createPlayer(name);
        state.playerId = res.data.id;
        state.playerName = name;
        state.token = res.token; // Capture the JWT from the response
        
        statusMsg.textContent = `✅ Cree ! ID: ${res.data.id.slice(0, 8)}...`;
        statusMsg.className = 'status-msg ok';
        nameInput.disabled = true;
        btnCreate.textContent = '✅ Connecte';
        roles.style.display = '';
      } catch (e) {
        statusMsg.textContent = `❌ ${e.message || 'API indisponible'}`;
        statusMsg.className = 'status-msg err';
        btnCreate.textContent = 'Creer mon personnage';
        btnCreate.disabled = false;
      }
    };

    document.getElementById('menu-gm').onclick = () => { if (state.playerId) navigate('gm-list'); };
    document.getElementById('menu-player').onclick = () => { if (state.playerId) navigate('player-list'); };
    document.getElementById('menu-inv').onclick = () => { if (state.playerId) navigate('inventory'); };
    document.getElementById('menu-auc').onclick = () => { if (state.playerId) navigate('auction'); };
    document.getElementById('menu-lb').onclick = () => navigate('leaderboard');
  }
};

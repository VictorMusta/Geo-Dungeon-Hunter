import { API } from '../api.js';
import { state, navigate } from '../app.js';

export const menuView = {
  init() {
    const nameInput = document.getElementById('menu-name');
    const btnCreate = document.getElementById('menu-create');
    const statusMsg = document.getElementById('menu-status');
    const loginBox = document.getElementById('menu-login-box');
    const roles = document.getElementById('menu-roles');

    if (state.playerId) {
      if (loginBox) loginBox.style.display = 'none';
      roles.style.display = '';
    } else {
      if (loginBox) loginBox.style.display = 'block';
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
        
        statusMsg.textContent = `✅ Cree !`;
        statusMsg.className = 'status-msg ok';
        
        // Refresh the view and update global status bar
        navigate('menu');
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

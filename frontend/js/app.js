const views = {};
let currentView = null;

export const state = {
  playerId: localStorage.getItem('dq_playerId'),
  playerName: localStorage.getItem('dq_playerName'),
  token: localStorage.getItem('dq_token'),
  role: null,
};

export function registerView(name, mod) {
  views[name] = mod;
}

export function navigate(name, data = {}) {
  // Sync state to localStorage for persistence
  if (state.playerId) {
    localStorage.setItem('dq_playerId', state.playerId);
    localStorage.setItem('dq_playerName', state.playerName);
    localStorage.setItem('dq_token', state.token);
  } else {
    localStorage.removeItem('dq_playerId');
    localStorage.removeItem('dq_playerName');
    localStorage.removeItem('dq_token');
  }

  document.querySelectorAll('.view').forEach(v => v.classList.remove('active'));
  const el = document.getElementById(`view-${name}`);
  if (el) el.classList.add('active');

  if (currentView && views[currentView]?.destroy) views[currentView].destroy();
  currentView = name;
  if (views[name]?.init) views[name].init(data);

  updateUserStatus();
}

export function logout() {
  state.playerId = null;
  state.playerName = null;
  state.token = null;
  navigate('menu');
}

export function updateUserStatus() {
  const bar = document.getElementById('user-status-bar');
  const nameEl = document.getElementById('usb-name');
  if (state.playerId && currentView === 'menu') {
    bar.style.display = 'block';
    nameEl.textContent = state.playerName || 'Anonyme';
  } else {
    bar.style.display = 'none';
  }
}

// Bind logout once
document.getElementById('usb-logout').onclick = logout;

export function toast(msg, isError = false) {
  const el = document.createElement('div');
  el.className = `toast ${isError ? 'toast-error' : 'toast-success'}`;
  el.textContent = msg;
  document.body.appendChild(el);
  setTimeout(() => el.remove(), 3000);
}

/**
 * Custom modern confirmation modal
 * @param {string} title 
 * @param {string} msg 
 * @param {string} icon 
 * @returns {Promise<boolean>}
 */
export function confirmModal(title, msg, icon = '❓') {
  return new Promise((resolve) => {
    const modal = document.getElementById('modal-confirm');
    const titleEl = document.getElementById('confirm-title');
    const msgEl = document.getElementById('confirm-msg');
    const iconEl = document.getElementById('confirm-icon');
    const okBtn = document.getElementById('confirm-ok');
    const cancelBtn = document.getElementById('confirm-cancel');

    titleEl.textContent = title;
    msgEl.textContent = msg;
    iconEl.textContent = icon;

    modal.style.display = 'flex';

    const cleanUp = (val) => {
      modal.style.display = 'none';
      resolve(val);
    };

    okBtn.onclick = () => cleanUp(true);
    cancelBtn.onclick = () => cleanUp(false);
    
    // Close on overlay click
    modal.onclick = (e) => {
      if (e.target === modal) cleanUp(false);
    };
  });
}

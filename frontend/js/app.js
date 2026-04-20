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
  if (state.playerId) localStorage.setItem('dq_playerId', state.playerId);
  if (state.playerName) localStorage.setItem('dq_playerName', state.playerName);
  if (state.token) localStorage.setItem('dq_token', state.token);

  document.querySelectorAll('.view').forEach(v => v.classList.remove('active'));
  const el = document.getElementById(`view-${name}`);
  if (el) el.classList.add('active');

  if (currentView && views[currentView]?.destroy) views[currentView].destroy();
  currentView = name;
  if (views[name]?.init) views[name].init(data);
}

export function toast(msg, isError = false) {
  const el = document.createElement('div');
  el.className = `toast ${isError ? 'toast-error' : 'toast-success'}`;
  el.textContent = msg;
  document.body.appendChild(el);
  setTimeout(() => el.remove(), 3000);
}

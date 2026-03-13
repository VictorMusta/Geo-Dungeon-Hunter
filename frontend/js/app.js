const views = {};
let currentView = null;

export const state = {
  playerId: null,
  playerName: null,
  role: null,
};

export function registerView(name, mod) {
  views[name] = mod;
}

export function navigate(name, data = {}) {
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

const BASE_URL = 'http://localhost:8080';

async function request(method, path, body = null) {
  const token = localStorage.getItem('dq_token');
  const opts = {
    method,
    headers: { 
      'Content-Type': 'application/json',
      ...(token ? { 'Authorization': `Bearer ${token}` } : {})
    },
  };
  if (body) opts.body = JSON.stringify(body);
  const res = await fetch(`${BASE_URL}${path}`, opts);
  const data = await res.json();
  
  // Store token if present in response
  if (data.token) {
    localStorage.setItem('dq_token', data.token);
  }

  if (!res.ok) throw { status: res.status, ...data };
  return data;
}

export const API = {
  createPlayer: (name) => request('POST', '/v1/players', { display_name: name, gold: 0 }),
  getPlayer: (id) => request('GET', `/v1/players/${id}`),
  getPlayers: () => request('GET', '/v1/players'),

  createItem: (item) => request('POST', '/v1/items', item),
  getItems: () => request('GET', '/v1/items'),

  createDungeon: (d) => request('POST', '/v1/mj/dungeons', d),
  getMJDungeons: (mjId) => request('GET', `/v1/mj/dungeons?mjId=${mjId}`),
  updateDungeon: (id, d) => request('PUT', `/v1/mj/dungeons/${id}`, d),
  publishDungeon: (id) => request('POST', `/v1/mj/dungeons/${id}/publish`),
  createStep: (did, s) => request('POST', `/v1/mj/dungeons/${did}/steps`, s),
  updateStep: (did, sid, s) => request('PUT', `/v1/mj/dungeons/${did}/steps/${sid}`, s),
  deleteStep: (did, sid) => request('DELETE', `/v1/mj/dungeons/${did}/steps/${sid}`),

  getDungeons: () => request('GET', '/v1/dungeons'),
  getDungeon: (id) => request('GET', `/v1/dungeons/${id}`),

  createRun: (dungeonId, playerId) => request('POST', '/v1/runs', { dungeonId, playerId }),
  getRuns: (playerId) => request('GET', `/v1/runs?playerId=${playerId}`),
  getRun: (id) => request('GET', `/v1/runs/${id}`),
  abandonRun: (id) => request('POST', `/v1/runs/${id}/abandon`),
  attemptBoss: (runId, stepId, lat, lon) =>
    request('POST', `/v1/runs/${runId}/steps/${stepId}/attempt`, { lat, lon }),

  getInventory: (playerId) => request('GET', `/v1/inventory?playerId=${playerId}`),

  getLeaderboard: (type, limit = 10) => request('GET', `/v1/leaderboard?type=${type}&limit=${limit}`),
};

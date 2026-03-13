import { registerView, navigate } from './app.js';
import { menuView } from './views/menu.js';
import { gmListView } from './views/gm-list.js';
import { gmEditView } from './views/gm-edit.js';
import { playerListView } from './views/player-list.js';
import { playerRunView } from './views/player-run.js';
import { leaderboardView } from './views/leaderboard.js';

registerView('menu', menuView);
registerView('gm-list', gmListView);
registerView('gm-edit', gmEditView);
registerView('player-list', playerListView);
registerView('player-run', playerRunView);
registerView('leaderboard', leaderboardView);

navigate('menu');

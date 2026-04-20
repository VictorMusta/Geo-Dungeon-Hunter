import { registerView, navigate } from './app.js';
import { menuView } from './views/menu.js';
import { gmListView } from './views/gm-list.js';
import { gmEditView } from './views/gm-edit.js';
import { gmItemsView } from './views/gm-items.js';
import { playerListView } from './views/player-list.js';
import { playerRunView } from './views/player-run.js';
import { leaderboardView } from './views/leaderboard.js';
import { inventoryView } from './views/inventory.js';
import { auctionView } from './views/auction.js';

registerView('menu', menuView);
registerView('gm-list', gmListView);
registerView('gm-edit', gmEditView);
registerView('gm-items', gmItemsView);
registerView('player-list', playerListView);
registerView('player-run', playerRunView);
registerView('leaderboard', leaderboardView);
registerView('inventory', inventoryView);
registerView('auction', auctionView);

// Expose views to window scope for inline HTML handlers
window.gmEditView = gmEditView;
window.gmListView = gmListView;
window.gmItemsView = gmItemsView;
window.playerListView = playerListView;
window.playerRunView = playerRunView;
window.leaderboardView = leaderboardView;

navigate('menu');

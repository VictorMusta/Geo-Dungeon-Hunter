import { API } from '../api.js';
import { state, navigate, toast } from '../app.js';

export const inventoryView = {
  async init() {
    document.getElementById('inv-back').onclick = () => navigate('menu');
    await this.loadData();
  },

  async loadData() {
    const listEl = document.getElementById('inv-list');
    listEl.innerHTML = '<p class="status-msg">Chargement...</p>';
    
    try {
      // Load gold
      const pRes = await API.getPlayer(state.playerId);
      document.getElementById('inv-gold').textContent = pRes.data.gold;

      // Load inventory
      const iRes = await API.getInventory(state.playerId);
      const items = iRes.data.items || [];

      if (items.length === 0) {
        listEl.innerHTML = '<p class="subtitle shadow-cyan" style="text-align:center; margin-top:40px;">Ton sac est vide... 🕸️</p>';
        return;
      }

      listEl.innerHTML = '';
      items.forEach(item => {
        const card = document.createElement('div');
        card.className = 'card';
        card.innerHTML = `
          <div style="display:flex; justify-content:space-between; align-items:center;">
            <div>
              <h3 class="text-gradient shadow-magenta">${item.itemId}</h3>
              <p style="color: var(--acc-cyan); font-weight:bold;">Quantité : ${item.qty}</p>
            </div>
            <button class="btn btn-sm sell-btn" data-id="${item.itemId}" data-qty="${item.qty}">🏷️ Vendre</button>
          </div>
        `;
        listEl.appendChild(card);
      });

      // Bind sell buttons
      document.querySelectorAll('.sell-btn').forEach(btn => {
        btn.onclick = (e) => this.promptSell(e.target.dataset.id, parseInt(e.target.dataset.qty));
      });

    } catch (e) {
      listEl.innerHTML = `<p class="status-msg err">Erreur de chargement de l'inventaire</p>`;
    }
  },

  async promptSell(itemId, maxQty) {
    const qtyStr = prompt(`Combien de "${itemId}" veux-tu vendre ? (Max: ${maxQty})`, "1");
    if (!qtyStr) return;
    const qty = parseInt(qtyStr);
    if (isNaN(qty) || qty <= 0 || qty > maxQty) {
      toast("Quantité invalide", true);
      return;
    }

    const priceStr = prompt(`Fixe un prix unitaire (en Gold) pour chaque "${itemId}" :`, "10");
    if (!priceStr) return;
    const price = parseInt(priceStr);
    if (isNaN(price) || price <= 0) {
      toast("Prix invalide", true);
      return;
    }

    try {
      await API.createAuction({
        sellerId: state.playerId,
        itemId: itemId,
        qty: qty,
        pricePerUnit: price
      });
      toast("✅ Objet mis en vente au marché !");
      this.loadData();
    } catch (e) {
      toast(`❌ Erreur: ${e.message}`, true);
    }
  }
};

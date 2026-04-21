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
        const itemInfo = item.item || { name: 'Objet Inconnu', rarity: 'common', type: 'unknown', description: '???' };
        const card = document.createElement('div');
        card.className = 'card';
        card.style.borderColor = `var(--acc-${itemInfo.rarity === 'uncommon' ? 'cyan' : itemInfo.rarity === 'rare' ? 'purple' : itemInfo.rarity === 'epic' ? 'magenta' : itemInfo.rarity === 'legendary' ? 'yellow' : 'purple'})`; // Fallback to purple/gray
        
        if (itemInfo.rarity === 'common') card.style.borderColor = '#8899aa';

        card.innerHTML = `
          <div style="display:flex; justify-content:space-between; align-items:flex-start;">
            <div style="flex:1;">
              <h3 class="text-gradient shadow-magenta mb-2" title="ID: ${item.itemId}" style="cursor:help; display:inline-block;">
                ${itemInfo.name}
              </h3>
              <div style="display:flex; gap:8px; margin-bottom:8px; align-items:center;">
                <span class="badge rarity-${itemInfo.rarity}" style="background:rgba(255,255,255,0.05); border:1px solid currentColor;">
                  ${itemInfo.rarity.toUpperCase()}
                </span>
                <span style="font-size:12px; opacity:0.6; text-transform:uppercase; font-weight:bold;">
                  ${itemInfo.type}
                </span>
              </div>
              <p style="font-size:14px; color:var(--muted-text); margin-bottom:12px;">${itemInfo.description || 'Aucune description disponible.'}</p>
              <p style="color: var(--acc-cyan); font-weight:900; font-size:18px;">📦 Quantité : ${item.qty}</p>
            </div>
            <button class="btn btn-sm sell-btn" data-id="${item.itemId}" data-qty="${item.qty}" data-name="${itemInfo.name}">🏷️ Vendre</button>
          </div>
        `;
        listEl.appendChild(card);
      });

      // Bind sell buttons
      document.querySelectorAll('.sell-btn').forEach(btn => {
        btn.onclick = (e) => this.promptSell(e.target.dataset.id, parseInt(e.target.dataset.qty), e.target.dataset.name);
      });

    } catch (e) {
      listEl.innerHTML = `<p class="status-msg err">Erreur de chargement de l'inventaire</p>`;
    }
  },

  async promptSell(itemId, maxQty, itemName) {
    const modal = document.getElementById('modal-sell');
    const nameEl = document.getElementById('sell-item-name');
    const maxQtyEl = document.getElementById('sell-max-qty');
    const qtyInput = document.getElementById('sell-qty');
    const priceInput = document.getElementById('sell-price');
    const totalEl = document.getElementById('sell-total');
    const confirmBtn = document.getElementById('sell-confirm');
    const cancelBtn = document.getElementById('sell-cancel');

    nameEl.textContent = `Vendre : ${itemName || itemId}`;
    maxQtyEl.textContent = maxQty;
    qtyInput.value = 1;
    qtyInput.max = maxQty;
    priceInput.value = 10;
    
    const updateTotal = () => {
      const q = parseInt(qtyInput.value) || 0;
      const p = parseInt(priceInput.value) || 0;
      totalEl.textContent = q * p;
    };

    qtyInput.oninput = updateTotal;
    priceInput.oninput = updateTotal;
    updateTotal();

    modal.style.display = 'flex';

    cancelBtn.onclick = () => {
      modal.style.display = 'none';
    };

    confirmBtn.onclick = async () => {
      const qty = parseInt(qtyInput.value);
      const price = parseInt(priceInput.value);

      if (isNaN(qty) || qty <= 0 || qty > maxQty) {
        toast("Quantité invalide", true);
        return;
      }
      if (isNaN(price) || price <= 0) {
        toast("Prix invalide", true);
        return;
      }

      try {
        confirmBtn.disabled = true;
        confirmBtn.textContent = "Placement...";
        
        await API.createAuction({
          sellerId: state.playerId,
          itemId: itemId,
          qty: qty,
          pricePerUnit: price
        });
        
        modal.style.display = 'none';
        toast("✅ Objet mis en vente au marché !");
        await this.loadData();
      } catch (e) {
        toast(`❌ Erreur: ${e.message}`, true);
      } finally {
        confirmBtn.disabled = false;
        confirmBtn.textContent = "Confirmer 🚀";
      }
    };
  }
};

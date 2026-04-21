import { API } from '../api.js';
import { state, navigate, toast, confirmModal } from '../app.js';

export const auctionView = {
  auctions: [],
  activeTab: 'market',

  async init() {
    document.getElementById('auc-back').onclick = () => navigate('menu');
    
    document.querySelectorAll('#view-auction .tab').forEach(t => {
      t.onclick = (e) => {
        document.querySelectorAll('#view-auction .tab').forEach(b => b.classList.remove('active'));
        e.target.classList.add('active');
        this.activeTab = e.target.dataset.key;
        this.renderList();
      };
    });

    await this.loadData();
  },

  async loadData() {
    const listEl = document.getElementById('auc-list');
    listEl.innerHTML = '<p class="status-msg">Chargement...</p>';
    
    try {
      const pRes = await API.getPlayer(state.playerId);
      document.getElementById('auc-gold').textContent = pRes.data.gold;

      const aRes = await API.getAuctions();
      this.auctions = (aRes.data || []).filter(a => a.status === 'active');
      this.renderList();
    } catch (e) {
      listEl.innerHTML = `<p class="status-msg err">Erreur de chargement du marché</p>`;
    }
  },

  renderList() {
    const listEl = document.getElementById('auc-list');
    listEl.innerHTML = '';

    const filtered = this.auctions.filter(a => {
      if (this.activeTab === 'market') return a.sellerId !== state.playerId;
      return a.sellerId === state.playerId;
    });

    if (filtered.length === 0) {
      listEl.innerHTML = `<p class="subtitle shadow-cyan" style="text-align:center; margin-top:40px;">Rien à afficher ici... 🎐</p>`;
      return;
    }

    filtered.forEach(a => {
      const isMine = this.activeTab === 'mine';
      const card = document.createElement('div');
      card.className = 'card';
      card.innerHTML = `
        <div style="display:flex; justify-content:space-between; align-items:center;">
          <div>
            <h3 class="text-gradient shadow-cyan">${a.itemId}</h3>
            <p style="color: var(--fg);">Quantité: <b>${a.qty}</b> <span style="opacity:0.5;">(Vendu par: ${a.sellerId.slice(0,6)})</span></p>
            <p style="color: var(--acc-yellow); font-weight:bold;">Prix unité: ${a.pricePerUnit} 💰 / Total: ${a.pricePerUnit * a.qty} 💰</p>
          </div>
          <button class="btn btn-sm ${isMine ? 'btn-secondary' : ''}" style="${isMine ? 'color:var(--acc-magenta); border-color:var(--acc-magenta);' : ''}" id="btn-auc-${a.id}">
            ${isMine ? '❌ Annuler' : '🛒 Acheter'}
          </button>
        </div>
      `;
      listEl.appendChild(card);

      document.getElementById(`btn-auc-${a.id}`).onclick = () => {
        if (isMine) this.cancelAuction(a.id);
        else this.buyAuction(a.id, a.pricePerUnit * a.qty);
      };
    });
  },

  async buyAuction(id, totalPrice) {
    const ok = await confirmModal(
      "Confirmation d'achat", 
      `Souhaites-tu vraiment acquérir cet objet pour la modique somme de ${totalPrice} 💰 ?`,
      "🛒"
    );
    if (!ok) return;

    try {
      await API.buyAuction(id, state.playerId, 1);
      toast("✅ Achat réussi !");
      this.loadData();
    } catch (e) {
      toast(`❌ Erreur: ${e.message}`, true);
    }
  },

  async cancelAuction(id) {
    const ok = await confirmModal(
      "Annuler la vente",
      "Es-tu certain de vouloir retirer cet objet du marché ? Il retournera dans ton sac.",
      "❌"
    );
    if (!ok) return;

    try {
      await API.cancelAuction(id, state.playerId);
      toast("✅ Vente annulée, obligé retourné dans le sac.");
      this.loadData();
    } catch (e) {
      toast(`❌ Erreur: ${e.message}`, true);
    }
  }
};

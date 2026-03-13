const DEFAULT_CENTER = [48.8566, 2.3522];
const DEFAULT_ZOOM = 13;

export class LeafletMap {
  constructor(divId, options = {}) {
    this.divId = divId;
    this.markers = [];
    this.circles = [];
    this.playerMarker = null;
    this.onClick = null;
    this.onRightClick = null;

    const center = options.center || DEFAULT_CENTER;
    const zoom = options.zoom || DEFAULT_ZOOM;

    this.map = L.map(divId, {
      center,
      zoom,
      zoomControl: true,
    });

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors',
      maxZoom: 19,
    }).addTo(this.map);

    this.map.on('click', (e) => {
      if (this.onClick) this.onClick(e.latlng.lat, e.latlng.lng);
    });

    this.map.on('contextmenu', (e) => {
      if (this.onRightClick) this.onRightClick(e.latlng.lat, e.latlng.lng);
    });

    setTimeout(() => this.map.invalidateSize(), 200);
  }

  _emojiIcon(emoji, size = 28) {
    return L.divIcon({
      className: 'emoji-marker',
      html: `<span style="font-size:${size}px;">${emoji}</span>`,
      iconSize: [size + 4, size + 4],
      iconAnchor: [(size + 4) / 2, (size + 4) / 2],
    });
  }

  setItems(items) {
    this.markers.forEach(m => this.map.removeLayer(m));
    this.circles.forEach(c => this.map.removeLayer(c));
    this.markers = [];
    this.circles = [];

    for (const item of items) {
      if (!item.location) continue;
      const lat = item.location.lat;
      const lon = item.location.lon;
      const radius = item.location.radiusMeters || 500;

      const circle = L.circle([lat, lon], {
        radius,
        color: item.strokeColor || 'rgba(233,69,96,0.7)',
        fillColor: item.circleColor || 'rgba(233,69,96,0.15)',
        fillOpacity: 0.3,
        weight: 2,
      }).addTo(this.map);
      this.circles.push(circle);

      const emoji = item.emoji || '❓';
      const marker = L.marker([lat, lon], {
        icon: this._emojiIcon(emoji, item.fontSize || 28),
      }).addTo(this.map);

      if (item.label) {
        marker.bindTooltip(item.label, {
          permanent: true,
          direction: 'bottom',
          offset: [0, 12],
          className: 'emoji-tooltip',
        });
      }
      this.markers.push(marker);
    }
  }

  setPlayer(lat, lon) {
    if (this.playerMarker) {
      this.playerMarker.setLatLng([lat, lon]);
    } else {
      this.playerMarker = L.marker([lat, lon], {
        icon: this._emojiIcon('🚶', 30),
        zIndexOffset: 1000,
      }).addTo(this.map);
      this.playerMarker.bindTooltip('MOI', {
        permanent: true,
        direction: 'top',
        offset: [0, -14],
        className: 'player-tooltip',
      });
    }
  }

  fitToItems() {
    const allLatLngs = [];
    this.markers.forEach(m => allLatLngs.push(m.getLatLng()));
    this.circles.forEach(c => allLatLngs.push(c.getLatLng()));
    if (this.playerMarker) allLatLngs.push(this.playerMarker.getLatLng());

    if (allLatLngs.length > 0) {
      const bounds = L.latLngBounds(allLatLngs);
      this.map.fitBounds(bounds.pad(0.2));
    }
  }

  getBounds() {
    const b = this.map.getBounds();
    return {
      minLat: b.getSouth(),
      maxLat: b.getNorth(),
      minLon: b.getWest(),
      maxLon: b.getEast(),
    };
  }

  invalidateSize() {
    this.map.invalidateSize();
  }

  destroy() {
    this.map.remove();
  }
}

export function haversine(lat1, lon1, lat2, lon2) {
  const R = 6371000;
  const dLat = (lat2 - lat1) * Math.PI / 180;
  const dLon = (lon2 - lon1) * Math.PI / 180;
  const a = Math.sin(dLat / 2) ** 2 +
            Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) *
            Math.sin(dLon / 2) ** 2;
  return R * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
}

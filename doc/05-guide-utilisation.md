# Guide d'utilisation de l'API

[< Retour au sommaire](README.md)

---

## Prerequis

1. **Go** 1.25+ installe
2. **MongoDB Atlas** configure (voir [MongoDB_Atlas_Guide_Complet.md](MongoDB_Atlas_Guide_Complet.md))
3. Fichier `.env` a la racine :

```env
API_VERSION="1.0.0"
MODE=DEVELOP
DB_HOST="mongodb+srv://user:password@cluster.mongodb.net/?appName=Cluster0"
TOKEN_KEY=""
API_PORT=":8080"
ALLOW_ORIGIN="*"
LOG_FORMAT="HUMAN"
```

## Lancer le serveur

```bash
go run cmd/api/*.go
```

Le serveur demarre sur `http://localhost:8080`. Verifier avec :

```bash
curl http://localhost:8080/ping
# -> {"status":200,"messageType":"ping.Done","message":"pong"}
```

---

## Tutoriel pas a pas : une partie complete

### Etape 1 - Creer des joueurs

```bash
# Creer le Game Master
curl -X POST http://localhost:8080/v1/players \
  -H "Content-Type: application/json" \
  -d '{"display_name": "MaitreDuDonjon", "gold": 0}'

# Creer un joueur
curl -X POST http://localhost:8080/v1/players \
  -H "Content-Type: application/json" \
  -d '{"display_name": "GuerrierFou", "gold": 100}'
```

Reponse :
```json
{
  "meta": { "object_name": "Player", "total_count": 1, "count": 1, "offSet": 0 },
  "data": {
    "id": "a1b2c3d4-...",
    "display_name": "GuerrierFou",
    "gold": 100,
    "created_at": "2026-03-12T10:00:00Z"
  }
}
```

### Etape 2 - Creer le catalogue d'items

```bash
curl -X POST http://localhost:8080/v1/items \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Epee de Flammes",
    "type": "weapon",
    "rarity": "rare",
    "description": "Une epee qui brule d un feu eternel",
    "stats": {"attack": 45, "fire_damage": 20},
    "tradable": true,
    "baseValue": 500
  }'
```

### Etape 3 - Creer un donjon

```bash
curl -X POST http://localhost:8080/v1/mj/dungeons \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Le Donjon de Montmartre",
    "description": "Parcourez les ruelles de Montmartre et affrontez ses gardiens",
    "createdBy": "MJ_ID_HERE",
    "area": {
      "name": "Montmartre, Paris",
      "boundingBox": {
        "minLat": 48.882, "maxLat": 48.892,
        "minLon": 2.333, "maxLon": 2.345
      }
    }
  }'
# -> Retourne l'ID du donjon, par ex: "dng-001"
```

### Etape 4 - Ajouter des boss steps

```bash
# Boss 1 : Le Gardien du Sacre-Coeur
curl -X POST http://localhost:8080/v1/mj/dungeons/dng-001/steps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Le Gardien du Sacre-Coeur",
    "location": { "lat": 48.8867, "lon": 2.3431, "radiusMeters": 50 },
    "zoneDescription": "Au pied de la Basilique",
    "difficulty": 3,
    "goldReward": 100,
    "lootTable": [
      { "itemId": "SWORD_ID", "dropRate": 0.3, "minQty": 1, "maxQty": 1 }
    ]
  }'

# Boss 2 : La Sorciere de la Place du Tertre
curl -X POST http://localhost:8080/v1/mj/dungeons/dng-001/steps \
  -H "Content-Type: application/json" \
  -d '{
    "name": "La Sorciere de la Place du Tertre",
    "location": { "lat": 48.8863, "lon": 2.3409, "radiusMeters": 30 },
    "zoneDescription": "Entre les chevalets des peintres",
    "difficulty": 7,
    "goldReward": 250,
    "lootTable": [
      { "itemId": "RING_ID", "dropRate": 0.1, "minQty": 1, "maxQty": 1 }
    ]
  }'
```

### Etape 5 - Publier le donjon

```bash
curl -X POST http://localhost:8080/v1/mj/dungeons/dng-001/publish
# -> 200 OK "dungeon published"
```

### Etape 6 - Demarrer un run

```bash
curl -X POST http://localhost:8080/v1/runs \
  -H "Content-Type: application/json" \
  -d '{ "dungeonId": "dng-001", "playerId": "PLAYER_ID" }'
# -> 201 Created, run avec state="active", currentStep=1
```

### Etape 7 - Tenter un boss

```bash
curl -X POST http://localhost:8080/v1/runs/RUN_ID/steps/STEP1_ID/attempt \
  -H "Content-Type: application/json" \
  -d '{ "lat": 48.8868, "lon": 2.3430 }'
```

Reponse reussie :
```json
{
  "meta": { "object_name": "AttemptResult", "total_count": 1, "count": 1 },
  "data": {
    "success": true,
    "rewards": {
      "gold": 100,
      "items": [{ "itemId": "SWORD_ID", "qty": 1 }]
    },
    "runCompleted": false
  }
}
```

Reponse si trop loin :
```json
{ "status": 409, "messageType": "attempt.Error", "message": "NOT_IN_RANGE" }
```

### Etape 8 - Voir l'inventaire

```bash
curl "http://localhost:8080/v1/inventory?playerId=PLAYER_ID"
```

### Etape 9 - Vendre un item a l'Auction House

```bash
curl -X POST http://localhost:8080/v1/auction/listings \
  -H "Content-Type: application/json" \
  -d '{
    "sellerId": "PLAYER_ID",
    "itemId": "SWORD_ID",
    "qty": 1,
    "pricePerUnit": 300
  }'
```

### Etape 10 - Consulter le leaderboard

```bash
curl "http://localhost:8080/v1/leaderboard?type=completions&limit=10"
curl "http://localhost:8080/v1/leaderboard?type=gold&limit=5"
curl "http://localhost:8080/v1/leaderboard?type=speed&dungeonId=dng-001&limit=3"
```

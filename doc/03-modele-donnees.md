# Modele de donnees

[< Retour au sommaire](README.md)

---

## Diagramme entite-relation

```mermaid
erDiagram
    PLAYER {
        string customID PK
        string displayName
        int64 gold
        datetime createdAt
        datetime updatedAt
    }

    DUNGEON {
        string customID PK
        string title
        string description
        string createdBy FK
        string status
        object area
        datetime createdAt
        datetime updatedAt
    }

    BOSS_STEP {
        string customID PK
        string dungeonId FK
        int order
        string name
        object location
        int difficulty
        int64 goldReward
        array lootTable
        datetime createdAt
    }

    RUN {
        string customID PK
        string dungeonId FK
        string playerId FK
        string state
        int currentStep
        array killedSteps
        datetime startedAt
        datetime endedAt
    }

    ITEM_DEF {
        string customID PK
        string name
        string type
        string rarity
        string description
        object stats
        boolean tradable
        int64 baseValue
    }

    INVENTORY {
        string playerId FK
        string itemId FK
        int64 qty
        datetime updatedAt
    }

    LISTING {
        string customID PK
        string sellerId FK
        string itemId FK
        int qty
        int64 pricePerUnit
        string status
        datetime expiresAt
    }

    TRADE {
        string customID PK
        string buyerId FK
        string sellerId FK
        string listingId FK
        int qty
        int64 totalPrice
        datetime createdAt
    }

    PLAYER ||--o{ RUN : "lance"
    PLAYER ||--o{ INVENTORY : "possede"
    PLAYER ||--o{ LISTING : "vend"
    PLAYER ||--o{ TRADE : "achete/vend"
    DUNGEON ||--o{ BOSS_STEP : "contient"
    DUNGEON ||--o{ RUN : "est joue via"
    ITEM_DEF ||--o{ INVENTORY : "reference"
    ITEM_DEF ||--o{ LISTING : "concerne"
    LISTING ||--o{ TRADE : "genere"
```

## Collections MongoDB

| Collection | Cle primaire | Index supplementaires |
|-----------|-------------|----------------------|
| `player` | `customID` | -- |
| `dungeon` | `customID` | `status` |
| `boss_step` | `customID` | `dungeonId` + `order` |
| `run` | `customID` | `playerId` + `dungeonId` + `state` |
| `item` | `customID` | -- |
| `inventory` | composite `(playerId, itemId)` | -- |
| `listing` | `customID` | `status` |
| `trade` | `customID` | -- |

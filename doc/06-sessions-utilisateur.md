# Sessions utilisateur

[< Retour au sommaire](README.md)

---

## Session Game Master : creation d'un donjon

```mermaid
sequenceDiagram
    actor MJ as Game Master

    MJ->>API: POST /v1/players {display_name: "DarkLord42"}
    API-->>MJ: 201 {id: "mj-001"}

    MJ->>API: POST /v1/items {name: "Dague Maudite", type: "weapon", rarity: "epic", ...}
    API-->>MJ: 201 {id: "item-dague"}

    MJ->>API: POST /v1/items {name: "Potion de Vie", type: "consumable", rarity: "common", ...}
    API-->>MJ: 201 {id: "item-potion"}

    MJ->>API: POST /v1/mj/dungeons {title: "Catacombes de Paris", createdBy: "mj-001", ...}
    API-->>MJ: 201 {id: "dng-cata", status: "draft"}

    MJ->>API: POST /v1/mj/dungeons/dng-cata/steps {name: "Gardien de l'Entree", difficulty: 2, ...}
    API-->>MJ: 201 {id: "step-1", order: 1}

    MJ->>API: POST /v1/mj/dungeons/dng-cata/steps {name: "Le Roi des Os", difficulty: 8, ...}
    API-->>MJ: 201 {id: "step-2", order: 2}

    MJ->>API: PUT /v1/mj/dungeons/dng-cata/steps/reorder {order: ["step-2", "step-1"]}
    API-->>MJ: 200 "steps reordered"

    Note over MJ: Change d'avis, remet l'ordre initial
    MJ->>API: PUT /v1/mj/dungeons/dng-cata/steps/reorder {order: ["step-1", "step-2"]}
    API-->>MJ: 200 OK

    MJ->>API: POST /v1/mj/dungeons/dng-cata/publish
    API-->>MJ: 200 "dungeon published"

    Note over MJ: Le donjon est maintenant visible par les joueurs
```

## Session Joueur : parcours d'un donjon

```mermaid
sequenceDiagram
    actor P as Joueur

    P->>API: GET /v1/dungeons
    API-->>P: Liste des donjons publies

    P->>API: GET /v1/dungeons/dng-cata
    API-->>P: Detail + 2 boss steps

    P->>API: POST /v1/runs {dungeonId: "dng-cata", playerId: "p-001"}
    API-->>P: 201 {id: "run-001", state: "active", currentStep: 1}

    Note over P: Se deplace vers le Boss 1
    P->>API: POST /v1/runs/run-001/steps/step-1/attempt {lat: ..., lon: ...}
    API-->>P: 200 {rewards: {gold: 50, items: [...]}, runCompleted: false}

    Note over P: Se deplace vers le Boss 2
    P->>API: POST /v1/runs/run-001/steps/step-2/attempt {lat: ..., lon: ...}
    API-->>P: 200 {rewards: {gold: 200, items: [...]}, runCompleted: true}

    P->>API: GET /v1/inventory?playerId=p-001
    API-->>P: {items: [{itemId: "item-dague", qty: 1}, ...]}

    P->>API: POST /v1/auction/listings {sellerId: "p-001", itemId: "item-dague", qty: 1, pricePerUnit: 400}
    API-->>P: 201 Listing creee

    P->>API: GET /v1/leaderboard?type=completions
    API-->>P: Classement
```

## Session Auction House : echange entre joueurs

```mermaid
sequenceDiagram
    actor Seller as Joueur A (Vendeur)
    actor Buyer as Joueur B (Acheteur)

    Seller->>API: POST /v1/auction/listings<br/>{sellerId: "a", itemId: "epee", qty: 1, pricePerUnit: 200}
    API-->>Seller: 201 {id: "list-001", status: "active"}
    Note over API: Items retires de l'inventaire du vendeur

    Buyer->>API: GET /v1/auction/listings
    API-->>Buyer: [{id: "list-001", itemId: "epee", pricePerUnit: 200, ...}]

    Buyer->>API: POST /v1/auction/listings/list-001/buy<br/>{buyerId: "b", qty: 1}
    Note over API: Transaction atomique
    API-->>Buyer: 200 "purchase successful"
    Note over Buyer: Gold -200, +1 Epee
    Note over Seller: Gold +200
```

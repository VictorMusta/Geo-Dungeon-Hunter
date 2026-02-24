# 🏰 Dungeon Real-World API -- Student Project

## 📌 Project Overview

This project consists of building a REST API in Go for a real-world
dungeon game.

A **Game Master (MJ)** can create dungeons located in real geographical
areas.\
Players must physically travel to locations, defeat bosses in order, and
earn rewards.

Your mission: **implement the remaining backend logic following the
specifications below.**

------------------------------------------------------------------------

# 🧩 Domain Model

## 1️⃣ Dungeon

A dungeon is a campaign created by a Game Master, linked to a real-world
geographic area.

### Dungeon

-   id
-   title
-   description
-   createdBy (mjId)
-   area (name, optional bounding box)
-   bosses\[\] (ordered steps)
-   status (draft / published / archived)
-   createdAt
-   updatedAt

------------------------------------------------------------------------

## 2️⃣ Boss / Step

Each boss has a physical location and a difficulty level.

### BossStep

-   id
-   dungeonId
-   order (1..n)
-   name
-   location:
    -   lat
    -   lon
    -   radiusMeters ⚠️ (used to validate player presence)
-   zoneDescription
-   difficulty (1..10 or enum)
-   rewards (fixed or loot table)
-   bossState (optional -- for respawn logic)

⚠️ `radiusMeters` is mandatory to validate "on-site" presence.

------------------------------------------------------------------------

## 3️⃣ Run (Session)

A Run represents one player playing one dungeon.

### Run

-   id
-   dungeonId
-   playerId
-   state (active / completed / abandoned)
-   currentStep (index or bossStepId)
-   killedSteps\[\]:
    -   bossStepId
    -   killedAt
    -   proof/meta
-   startedAt
-   endedAt

------------------------------------------------------------------------

## 4️⃣ Player / Account

### Player

-   id
-   displayName
-   wallet (gold)
-   inventory (items + quantities)
-   createdAt

------------------------------------------------------------------------

## 5️⃣ Items

### ItemDef (catalog)

-   id
-   type (weapon / artifact / consumable)
-   rarity
-   stats (JSON)

### InventoryEntry

-   playerId
-   itemId
-   qty

------------------------------------------------------------------------

## 6️⃣ Auction House

### Listing

-   id
-   sellerId
-   itemId
-   qty
-   pricePerUnit
-   status (active / sold / cancelled / expired)
-   createdAt
-   expiresAt

### Trade (optional but recommended)

-   buyerId
-   sellerId
-   listingId
-   qty
-   totalPrice
-   createdAt

------------------------------------------------------------------------

# 📍 Geolocation & Anti-Cheat

## MVP Rules

-   Client sends `lat/lon` when attempting a boss.
-   Server validates: `distance <= radiusMeters`.
-   Log:
    -   attemptId
    -   deviceTime
    -   optional GPS accuracy

## Future (V2)

-   Play Integrity / DeviceCheck
-   Mock location detection
-   Signed location proofs

------------------------------------------------------------------------

# 🌐 REST Endpoints (MVP)

## Auth

POST /v1/auth/register\
POST /v1/auth/login\
GET /v1/me

------------------------------------------------------------------------

## Game Master (Dungeon Management)

POST /v1/mj/dungeons\
PUT /v1/mj/dungeons/{id}\
POST /v1/mj/dungeons/{id}/publish\
POST /v1/mj/dungeons/{id}/steps\
PUT /v1/mj/dungeons/{id}/steps/{stepId}\
PUT /v1/mj/dungeons/{id}/steps/reorder

------------------------------------------------------------------------

## Player

GET /v1/dungeons\
GET /v1/dungeons/{id}\
POST /v1/runs\
GET /v1/runs\
GET /v1/runs/{id}

------------------------------------------------------------------------

## Boss Attempt

POST /v1/runs/{id}/steps/{stepId}/attempt

Body: { "lat": 0.0, "lon": 0.0 }

Responses:

409 NOT_IN_RANGE\
409 WRONG_STEP_ORDER\
200 OK + rewards

------------------------------------------------------------------------

## Inventory & Economy

GET /v1/inventory\
POST /v1/auction/listings\
GET /v1/auction/listings\
POST /v1/auction/listings/{id}/buy\
POST /v1/auction/listings/{id}/cancel

------------------------------------------------------------------------

# ⚙️ Business Rules (Mandatory)

1.  Ordered progression: a player can only defeat `currentStep`.
2.  Location validation: use Haversine distance formula.
3.  Atomic rewards: boss kill must update wallet + inventory +
    progression in one transaction.
4.  Atomic auction transactions: debit buyer, credit seller, transfer
    items in one DB transaction.
5.  Idempotency: boss attempt endpoint must prevent double rewards
    (retry-safe).

------------------------------------------------------------------------

# 🎯 Student Objectives

You must:

-   Implement the database schema
-   Implement services and controllers
-   Handle transactions properly
-   Implement distance calculation logic
-   Ensure idempotent boss attempts
-   Add validation & proper HTTP error codes

------------------------------------------------------------------------

# 🏆 Bonus Ideas

-   Leaderboard
-   Multiplayer runs
-   Guild system
-   Dynamic loot tables
-   Rate limiting
-   Event-based world bosses

------------------------------------------------------------------------

# 🧠 Technical Recommendations

-   Use Mongo-DB 
-   Use transactions for economy and boss rewards
-   Use clean architecture or controller/service separation
-   Write unit tests for:
    -   distance validation
    -   progression rules
    -   transaction safety

------------------------------------------------------------------------

Happy coding 🚀

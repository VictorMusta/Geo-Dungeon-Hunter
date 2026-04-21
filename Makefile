.PHONY: help install up down back front test test-integration tidy logs zip

# --- AIDE ---
help: ## Affiche cette aide
ifeq ($(OS),Windows_NT)
	@findstr /R /C:"^[a-zA-Z_-]*:.*##" $(MAKEFILE_LIST) | sort
else
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
endif

# --- INSTALLATION ---
install: ## Installe toutes les dépendances (Go + Node)
	@echo "📦 Installation des dépendances Go..."
	go mod tidy
	@echo "📦 Installation des dépendances Frontend..."
	cd frontend && npm install

# --- INFRASTRUCTURE (DOCKER) ---
up: ## Démarre la base de données et les outils dev (Mongo Express, Swagger)
	@echo "🚀 Démarrage de l'infrastructure Docker..."
	docker compose up -d
	@echo "✨ MongoDB: localhost:27017"
	@echo "📈 Mongo Express: http://localhost:8081"
	@echo "📖 Swagger UI: http://localhost:8082"

down: ## Arrête l'infrastructure Docker
	@echo "🛑 Arrêt de l'infrastructure..."
	docker compose down

clear-db: ## Supprime TOUTES les données de la base (Docker Volumes)
	@echo "⚠️ Suppression des données de la base..."
	docker compose down -v
	@echo "🚀 Redémarrage et initialisation de la base..."
	docker compose up -d mongodb
	@timeout /t 5 > nul
	@docker exec dungeons-db mongosh --eval "rs.initiate({_id:'rs0', members:[{_id:0, host:'localhost:27017'}]})" > nul 2>&1 || echo "⚠️ Note: Réinitialisation déjà effectuée ou en cours."
	@docker compose down
	@echo "✅ Base de données réinitialisée et prête."

# --- DEVELOPPEMENT ---
dev: ## Lance TOUT le projet (DB, Backend, Frontend) dans un seul terminal
	@echo "🚀 Démarrage global du projet..."
	docker compose up -d
	npx -y concurrently --kill-others \
		"make back" \
		"make front" \
		--prefix-colors "cyan,magenta"

back: ## Lance le backend en mode live-reload (Air)
	@echo "🔥 Lancement du backend avec Air..."
	air

front: ## Lance le serveur de développement Frontend
	@echo "💻 Lancement du Frontend..."
	cd frontend && npm start

# --- TESTS ---
test: ## Exécute les tests unitaires
	@echo "🧪 Tests unitaires..."
	go test ./app/functions/... -v

test-integration: ## Exécute les tests d'intégration (Nécessite Docker UP)
	@echo "🧪 Tests d'intégration..."
	go test ./tests/... -v

# --- MAINTENANCE ---
fmt: ## Formate le code Go
	go fmt ./...

tidy: ## Nettoie et optimise les modules Go
	go mod tidy

# --- LIVRAISON ---
zip: ## Crée l'archive .zip pour le rendu (livre un projet propre et prêt à l'emploi)
	@echo "📦 Préparation de l'archive de livraison..."
	@if exist dist rmdir /s /q dist
	@if exist delivery.zip del delivery.zip
	@mkdir dist
	@echo "📂 Copie des fichiers sources..."
	@xcopy /E /I /Y app dist\app > nul
	@xcopy /E /I /Y cmd dist\cmd > nul
	@xcopy /E /I /Y frontend dist\frontend > nul
	@xcopy /E /I /Y tests dist\tests > nul
	@xcopy /E /I /Y doc dist\doc > nul
	@copy go.mod dist\ > nul
	@copy go.sum dist\ > nul
	@copy Makefile dist\ > nul
	@copy README.md dist\ > nul
	@copy docker-compose.yml dist\ > nul
	@copy .env dist\ > nul
	@copy .air.toml dist\ > nul
	@echo "🧹 Nettoyage des dossiers inutiles..."
	@if exist dist\frontend\node_modules rmdir /s /q dist\frontend\node_modules
	@if exist dist\tmp rmdir /s /q dist\tmp
	@if exist dist\graphify-out rmdir /s /q dist\graphify-out
	@echo "🤐 Compression en cours (via tar)..."
	@tar -a -c -f delivery.zip -C dist .
	@rmdir /s /q dist
	@echo "✅ Archive 'delivery.zip' créée avec succès ! Elle contient tout le nécessaire pour le correcteur."

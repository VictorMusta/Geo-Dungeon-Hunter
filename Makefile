.PHONY: help install up down back front test test-integration tidy logs

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

logs: ## Affiche les logs Docker
	docker compose logs -f

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

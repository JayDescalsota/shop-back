#!/bin/bash
# service.sh - Bash runner script for backend microservices (Linux).
# Use: chmod +x service.sh && ./service.sh <command> [service]

set -euo pipefail

SERVICES=(
    gateway auth users vehicles bookings repair
    inventory parts-marketplace payments payroll
    lookup notifications search staff
)

info()  { echo -e "\033[36m[INFO]\033[0m $*"; }
ok()    { echo -e "\033[32m[SUCCESS]\033[0m $*"; }
err()   { echo -e "\033[31m[ERROR]\033[0m $*" >&2; }
warn()  { echo -e "\033[33m[WARN]\033[0m $*"; }

usage() {
    cat <<EOF
Usage: ./service.sh <command> [service]

Commands:
  setup                  Start database containers and run migrations
  start                  Start all services via Docker Compose
  start <service>        Start a specific service locally (build + run)
  stop                   Stop all running services
  restart                Restart all running services
  build                  Build all services (Go binaries + Docker images)
  build <service>        Build a specific service (Go binary + Docker image)
  rebuild [service]      Rebuild all services (Go + Docker + restart), or specific one
  logs [service]         Follow logs for all or a specific service
  status                 Show status of all services

Examples:
  ./service.sh build
  ./service.sh build auth
  ./service.sh start auth
  ./service.sh rebuild auth
EOF
}

CMD="${1:-help}"
SVC="${2:-}"

# ----- build one -----
build_service() {
    local svc="$1"
    info "Building '$svc'..."
    mkdir -p build
    go build -ldflags="-w -s" -o "build/$svc" "./services/$svc"
    ok "Built build/$svc"
    info "Building Docker image for '$svc'..."
    docker compose -p autolab build "$svc"
    ok "Docker image built for '$svc'."
}

# ----- build all -----
build_all() {
    info "Building all services..."
    mkdir -p build
    for svc in "${SERVICES[@]}"; do
        info "  Building '$svc'..."
        go build -ldflags="-w -s" -o "build/$svc" "./services/$svc"
        ok "    Built build/$svc"
    done
    info "Building all Docker images..."
    docker compose -p autolab build
    ok "All Docker images built."
}

# ----- start one locally -----
start_service() {
    local svc="$1"
    local binary="build/$svc"
    if [ ! -f "$binary" ]; then
        warn "Binary not found. Building first..."
        build_service "$svc"
    fi
    local pid
    pid=$(pgrep -x "$svc" 2>/dev/null || true)
    if [ -n "$pid" ]; then
        info "'$svc' is already running (PID $pid). Restarting..."
        kill "$pid" 2>/dev/null || true
        sleep 1
    fi
    info "Starting '$svc' in background..."
    local logfile="build/$svc.log"
    nohup "$binary" > "$logfile" 2>&1 &
    ok "'$svc' started (PID $!). Logs: $logfile"
}

# ----- rebuild one -----
rebuild_service() {
    local svc="$1"
    if ! docker compose -p autolab config --services 2>/dev/null | grep -qx "$svc"; then
        err "Service '$svc' not found in docker-compose.yml."
        exit 1
    fi
    info "Rebuilding '$svc'..."
    mkdir -p build
    go build -ldflags="-w -s" -o "build/$svc" "./services/$svc"
    ok "Built build/$svc"
    info "Stopping old '$svc' container..."
    docker compose -p autolab stop "$svc" 2>/dev/null || true
    docker compose -p autolab rm -f "$svc" 2>/dev/null || true
    docker compose -p autolab build "$svc"
    docker compose -p autolab up -d --no-deps "$svc"
    ok "'$svc' rebuilt and restarted."
}

# ----- rebuild all -----
rebuild_all() {
    info "Rebuilding all services..."
    mkdir -p build
    for svc in "${SERVICES[@]}"; do
        info "  Building '$svc'..."
        go build -ldflags="-w -s" -o "build/$svc" "./services/$svc"
    done
    info "Building Docker images..."
    docker compose -p autolab build
    info "Stopping and starting all services..."
    docker compose -p autolab down
    docker compose -p autolab up -d
    ok "All services rebuilt and restarted."
}

# ----- setup (db + migrations) -----
do_setup() {
    info "Starting database infrastructure..."
    docker compose -p autolab up -d postgres redis
    info "Waiting for PostgreSQL..."
    for i in $(seq 1 30); do
        if docker compose -p autolab exec -T postgres pg_isready -U postgres -d postgres >/dev/null 2>&1; then
            ok "Database is ready!"
            break
        fi
        sleep 1
    done
    info "Running migrations..."
    find services -path "*/migrations/*.sql" | sort | while read -r f; do
        info "  $(basename "$f")..."
        docker compose -p autolab exec -T postgres psql -U postgres -d postgres < "$f"
    done
    ok "Setup complete."
}

# ----- main dispatch -----
case "$CMD" in
    help)
        usage
        ;;
    build)
        if [ -n "$SVC" ]; then
            # validate
            found=0; for s in "${SERVICES[@]}"; do [ "$s" = "$SVC" ] && found=1; done
            if [ "$found" -eq 0 ]; then
                err "Unknown service '$SVC'. Available: ${SERVICES[*]}"
                exit 1
            fi
            build_service "$SVC"
        else
            build_all
        fi
        ;;
    start)
        if [ -n "$SVC" ]; then
            found=0; for s in "${SERVICES[@]}"; do [ "$s" = "$SVC" ] && found=1; done
            if [ "$found" -eq 0 ]; then
                err "Unknown service '$SVC'. Available: ${SERVICES[*]}"
                exit 1
            fi
            start_service "$SVC"
        else
            info "Starting all services via Docker Compose..."
            docker compose -p autolab up -d
            ok "All services started."
        fi
        ;;
    stop)
        info "Stopping all services..."
        docker compose -p autolab down
        ok "All services stopped."
        ;;
    restart)
        info "Restarting all services..."
        docker compose -p autolab restart
        ok "All services restarted."
        ;;
    rebuild)
        if [ -z "$SVC" ] || [ "$SVC" = "all" ]; then
            rebuild_all
        else
            rebuild_service "$SVC"
        fi
        ;;
    logs)
        if [ -n "$SVC" ]; then
            docker compose -p autolab logs -f "$SVC"
        else
            docker compose -p autolab logs -f
        fi
        ;;
    status)
        docker compose -p autolab ps
        ;;
    setup)
        do_setup
        ;;
    *)
        usage
        exit 1
        ;;
esac

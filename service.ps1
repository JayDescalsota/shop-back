# service.ps1 - PowerShell runner script for backend microservices.

function Write-Info ($msg) { Write-Host "[INFO] $msg" -ForegroundColor Cyan }
function Write-Success ($msg) { Write-Host "[SUCCESS] $msg" -ForegroundColor Green }
function Write-ErrorMsg ($msg) { Write-Host "[ERROR] $msg" -ForegroundColor Red }
function Write-Warn ($msg) { Write-Host "[WARN] $msg" -ForegroundColor Yellow }

$SERVICES = @(
    'gateway', 'auth', 'users', 'vehicles', 'bookings', 'repair',
    'inventory', 'parts-marketplace', 'payments', 'payroll',
    'lookup', 'notifications', 'search', 'staff'
)

function Show-Help {
    Write-Host "Usage: .\service.ps1 <command> [service]"
    Write-Host ""
    Write-Host "Commands:"
    Write-Host "  setup                  Start database containers and run migrations"
    Write-Host "  start                  Start all services via Docker Compose"
    Write-Host "  start <service>        Start a specific service locally (build + run)"
    Write-Host "  stop                   Stop all running services"
    Write-Host "  restart                Restart all running services"
    Write-Host "  build                  Build all services (Go binaries + Docker images)"
    Write-Host "  build <service>        Build a specific service (Go binary + Docker image)"
    Write-Host "  rebuild [service]      Rebuild all services (Go + Docker + restart), or specific one"
    Write-Host "  logs [service]         Follow logs for all or a specific service"
    Write-Host "  status                 Show status of all services"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\service.ps1 build"
    Write-Host "  .\service.ps1 build auth"
    Write-Host "  .\service.ps1 start auth"
    Write-Host "  .\service.ps1 rebuild auth"
}

if ($args.Count -eq 0 -or $args[0] -eq 'help') { Show-Help; exit 0 }

$Command = $args[0]
$Service = if ($args.Count -gt 1) { $args[1] } else { '' }

switch ($Command) {
    "build" {
        if ($Service -ne '') {
            if ($SERVICES -notcontains $Service) {
                Write-ErrorMsg "Unknown service '$Service'. Available: $($SERVICES -join ', ')"
                exit 1
            }
            Write-Info "Building '$Service'..."
            $null = New-Item -ItemType Directory -Path ".\build" -Force
            $Output = ".\build\$Service.exe"
            go build -ldflags="-w -s" -o $Output .\services\$Service
            Write-Success "Built $Output"
            Write-Info "Building Docker image for '$Service'..."
            docker compose -p autolab build $Service
            Write-Success "Docker image built for '$Service'."
        } else {
            Write-Info "Building all services..."
            $null = New-Item -ItemType Directory -Path ".\build" -Force
            foreach ($svc in $SERVICES) {
                Write-Info "  Building '$svc'..."
                $Output = ".\build\$svc.exe"
                go build -ldflags="-w -s" -o $Output .\services\$svc
                Write-Success "    Built $Output"
            }
            Write-Info "Building all Docker images..."
            docker compose -p autolab build
            Write-Success "All Docker images built."
        }
    }

    "start" {
        if ($Service -ne '') {
            if ($SERVICES -notcontains $Service) {
                Write-ErrorMsg "Unknown service '$Service'. Available: $($SERVICES -join ', ')"
                exit 1
            }
            $Binary = ".\build\$Service.exe"
            if (-not (Test-Path $Binary)) {
                Write-Warn "Binary not found. Building first..."
                & $PSCommandPath build $Service
            }
            $Existing = Get-Process -Name "$Service" -ErrorAction SilentlyContinue
            if ($Existing) {
                Write-Info "'$Service' is already running (PID $($Existing.Id)). Restarting..."
                Stop-Process -Id $Existing.Id -Force
                Start-Sleep -Seconds 1
            }
            Write-Info "Starting '$Service' in background..."
            $logFile = ".\build\$Service.log"
            $process = Start-Process -FilePath (Resolve-Path $Binary).Path -NoNewWindow -PassThru -RedirectStandardOutput $logFile
            Write-Success "'$Service' started (PID $($process.Id)). Logs: $logFile"
        } else {
            Write-Info "Starting all services via Docker Compose..."
            docker compose -p autolab up -d
            Write-Success "All services started."
        }
    }

    "stop" {
        Write-Info "Stopping all services..."
        docker compose -p autolab down
        Write-Success "All services stopped."
    }

    "restart" {
        Write-Info "Restarting all services..."
        docker compose -p autolab restart
        Write-Success "All services restarted."
    }

    "setup" {
        Write-Info "Starting database infrastructure..."
        docker compose -p autolab up -d postgres redis
        Write-Info "Waiting for PostgreSQL..."
        $ready = $false
        for ($i = 1; $i -le 30; $i++) {
            docker compose -p autolab exec -T postgres pg_isready -U postgres -d postgres > $null 2>&1
            if ($LastExitCode -eq 0) { $ready = $true; break }
            Start-Sleep -Seconds 1
        }
        if (-not $ready) { Write-ErrorMsg "PostgreSQL failed to start."; exit 1 }
        Write-Success "Database is ready!"
        Write-Info "Running migrations..."
        Get-ChildItem -Path "services" -Filter "*.sql" -Recurse | Where-Object { $_.FullName -match "\\migrations\\" } | Sort-Object Name | ForEach-Object {
            Write-Info "  $($_.Name)..."
            Get-Content $_.FullName -Raw | docker compose -p autolab exec -T postgres psql -U postgres -d postgres
        }
        Write-Success "Setup complete."
    }

    "rebuild" {
        if ($Service -eq '' -or $Service -eq 'all') {
            Write-Info "Rebuilding all services..."
            $null = New-Item -ItemType Directory -Path ".\build" -Force
            foreach ($svc in $SERVICES) {
                Write-Info "  Building '$svc'..."
                $Output = ".\build\$svc.exe"
                go build -ldflags="-w -s" -o $Output .\services\$svc
            }
            Write-Info "Building Docker images..."
            docker compose -p autolab build
            Write-Info "Stopping and starting all services..."
            docker compose -p autolab down
            docker compose -p autolab up -d
            Write-Success "All services rebuilt and restarted."
        } else {
            $available = docker compose -p autolab config --services
            if ($available -notcontains $Service) { Write-ErrorMsg "Service '$Service' not found in docker-compose.yml."; exit 1 }
            Write-Info "Rebuilding '$Service'..."
            $null = New-Item -ItemType Directory -Path ".\build" -Force
            $Output = ".\build\$Service.exe"
            go build -ldflags="-w -s" -o $Output .\services\$Service
            Write-Success "Built $Output"
            Write-Info "Stopping old '$Service' container..."
            docker compose -p autolab stop $Service
            docker compose -p autolab rm -f $Service
            docker compose -p autolab build $Service
            docker compose -p autolab up -d --no-deps $Service
            Write-Success "'$Service' rebuilt and restarted."
        }
    }

    "logs" {
        if ($Service -ne '') { docker compose -p autolab logs -f $Service }
        else { docker compose -p autolab logs -f }
    }

    "status" {
        docker compose -p autolab ps
    }

    default {
        Show-Help; exit 1
    }
}

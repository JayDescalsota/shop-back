# service.ps1 - PowerShell runner script for backend microservices and databases.
# Compatible with native Windows PowerShell.

function Write-Info ($msg) {
    Write-Host "[INFO] $msg" -ForegroundColor Cyan
}

function Write-Success ($msg) {
    Write-Host "[SUCCESS] $msg" -ForegroundColor Green
}

function Write-ErrorMsg ($msg) {
    Write-Host "[ERROR] $msg" -ForegroundColor Red
}

function Write-Warn ($msg) {
    Write-Host "[WARN] $msg" -ForegroundColor Yellow
}

function Show-Help {
    Write-Host "Usage: .\service.ps1 <command> [options]"
    Write-Host ""
    Write-Host "Commands:"
    Write-Host "  setup                    Start database containers and run migrations"
    Write-Host "  start                    Start all services (infra + microservices)"
    Write-Host "  stop                     Stop all running services"
    Write-Host "  restart                  Restart all running services"
    Write-Host "  rebuild <service>        Rebuild and restart a specific service"
    Write-Host "  rebuild -service <svc>   Rebuild and restart a specific service"
    Write-Host "  logs [service]           Follow logs for all or a specific service"
    Write-Host "  status                   Show status of all services"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\service.ps1 setup"
    Write-Host "  .\service.ps1 start"
    Write-Host "  .\service.ps1 rebuild auth"
    Write-Host "  .\service.ps1 rebuild -service auth"
}

# Ensure docker is installed
$dockerCmd = Get-Command docker -ErrorAction SilentlyContinue
if ($null -eq $dockerCmd) {
    Write-ErrorMsg "Docker is not installed or not in PATH."
    exit 1
}

if ($args.Count -eq 0) {
    Show-Help
    exit 0
}

$Command = $args[0]
$RemainingArgs = $args | Select-Object -Skip 1

switch ($Command) {
    "setup" {
        Write-Info "Starting database infrastructure (postgres, redis)..."
        docker compose up -d postgres redis
        
        Write-Info "Waiting for PostgreSQL database to be ready..."
        $ready = $false
        for ($i = 1; $i -le 30; $i++) {
            # Execute pg_isready inside postgres container
            docker compose exec -T postgres pg_isready -U postgres -d postgres > $null 2>&1
            if ($LastExitCode -eq 0) {
                $ready = $true
                break
            }
            Start-Sleep -Seconds 1
        }
        
        if ($ready) {
            Write-Success "Database is ready!"
            Write-Info "Running SQL migrations..."
            
            $migrationsFound = $false
            if (Test-Path "services") {
                # Find all SQL files recursively in 'services/*/migrations'
                # Sort them by name to execute in correct order
                $migrations = Get-ChildItem -Path "services" -Filter "*.sql" -Recurse | 
                              Where-Object { $_.FullName -match "\\migrations\\" } |
                              Sort-Object Name
                
                foreach ($migration in $migrations) {
                    $migrationsFound = $true
                    Write-Info "Executing migration: $($migration.FullName)"
                    # Read all contents and pipe it into docker compose exec
                    # Using raw content ensures it's read exactly as-is
                    Get-Content $migration.FullName -Raw | docker compose exec -T postgres psql -U postgres -d postgres
                }
            }
            
            if ($migrationsFound) {
                Write-Success "All migrations executed successfully."
            } else {
                Write-Warn "No migration files found under services/*/migrations/*.sql."
            }
            Write-Success "Database setup complete."
        } else {
            Write-ErrorMsg "PostgreSQL database failed to start or become ready in time."
            exit 1
        }
    }
    
    "start" {
        Write-Info "Starting all services (infra + microservices)..."
        docker compose up -d
        Write-Success "All services started successfully in the background."
        Write-Info "Use '.\service.ps1 logs' to view logs or '.\service.ps1 status' to view container status."
    }
    
    "stop" {
        Write-Info "Stopping all services..."
        docker compose down
        Write-Success "All services stopped."
    }
    
    "restart" {
        Write-Info "Restarting all services..."
        docker compose restart
        Write-Success "All services restarted."
    }
    
    "rebuild" {
        $serviceName = ""
        if ($RemainingArgs.Count -gt 0) {
            if ($RemainingArgs[0] -eq "-service") {
                if ($RemainingArgs.Count -gt 1) {
                    $serviceName = $RemainingArgs[1]
                }
            } else {
                $serviceName = $RemainingArgs[0]
            }
        }
        
        if ([string]::IsNullOrEmpty($serviceName)) {
            Write-ErrorMsg "Please specify a service to rebuild. Example: .\service.ps1 rebuild auth"
            exit 1
        }
        
        # Verify that the service is defined in docker-compose.yml
        $availableServices = docker compose config --services
        if ($availableServices -notcontains $serviceName) {
            Write-ErrorMsg "Service '$serviceName' not found in docker-compose.yml."
            Write-Info "Available services are: $($availableServices -join ', ')"
            exit 1
        }
        
        Write-Info "Rebuilding service '$serviceName'..."
        docker compose build $serviceName
        
        Write-Info "Restarting service '$serviceName'..."
        docker compose up -d --no-deps $serviceName
        Write-Success "Service '$serviceName' successfully rebuilt and restarted."
    }
    
    "logs" {
        # Feed all remaining arguments to logs command (e.g. .\service.ps1 logs auth)
        docker compose logs -f $RemainingArgs
    }
    
    "status" {
        docker compose ps
    }
    
    Default {
        Show-Help
    }
}

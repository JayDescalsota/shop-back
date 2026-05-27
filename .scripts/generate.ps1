param(
    [string]$Service = 'all'
)

$originalDir = Get-Location
$servicesDir = "services"
if (Test-Path $servicesDir) {
    $services = Get-ChildItem -Path $servicesDir -Directory | ForEach-Object { $_.Name }
} else {
    $services = @()
}

if ($Service -eq 'all') {
    foreach ($svc in $services) {
        Write-Host "Generating $svc..."
        Set-Location $originalDir
        Set-Location services\$svc
        go run github.com/99designs/gqlgen generate
    }
} else {
    if (-not $services.Contains($Service)) {
        Write-Error "Service '$Service' not found. Available: $($services -join ', ')"
        exit 1
    }
    Set-Location $originalDir
    Set-Location services\$Service
    go run github.com/99designs/gqlgen generate
}

Set-Location $originalDir

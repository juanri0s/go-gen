rm scripts/auth0-exercise.exe
go imports -w .
go build -o scripts/auth0-exercise.exe ./cmd/app

$env:PORT = 8080
$env:GITHUB_AUTH_TOKEN = ""

./scripts/auth0-exercise.exe
rm scripts/auth0-exercise.exe
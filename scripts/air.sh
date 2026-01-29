cd ..
go tool air --build.cmd "go build -o bin/GymApp cmd/App/main.go" --build.entrypoint "./bin/GymApp"

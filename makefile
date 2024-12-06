run:
	go build -o cmd/cmd.exe cmd/main.go && ./cmd/cmd.exe

migrate:
	go build -o cmd/cmd.exe cmd/main.go && ./cmd/cmd.exe --runMigrations
run:
	go build -o cmd/cmd.exe cmd/main.go && ./cmd/cmd.exe

run_with_migrations:
	go build -o cmd/cmd.exe cmd/main.go && ./cmd/cmd.exe --runMigrations

speed-test:
	go run tests/speed-test.go
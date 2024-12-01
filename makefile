run:
	cd cmd && go build && ./cmd.exe

migrate:
	cd cmd && go build && ./cmd.exe --runMigrations
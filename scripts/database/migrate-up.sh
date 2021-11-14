source ./scripts/database/get-env.sh

migrate -path ./db/migrations -database "${DATABASE_URL}" -verbose up
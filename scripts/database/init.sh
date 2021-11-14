source ./scripts/database/get-env.sh

docker exec cryptomath-captcha-postgres createdb --username="${POSTGRES_USER}" "${POSTGRES_DB}"
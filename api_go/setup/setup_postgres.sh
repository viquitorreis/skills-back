#!/bin/bash

# Permissão antes de rodar:
# chmod +x setup_postgres.sh

# Aguardar até que o PostgreSQL esteja pronto para aceitar conexões
until docker exec -it skills pg_isready -U postgres -h localhost -p 5432; do
  echo "Aguardando o PostgreSQL iniciar..."
  sleep 1
done

# Executar o comando SQL para criar a extensão citext
docker exec -it skills psql -U postgres -d postgres -c "CREATE EXTENSION IF NOT EXISTS citext;"

echo "Extensão citext criada com sucesso!"

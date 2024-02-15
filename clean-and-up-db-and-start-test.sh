docker rm -f $(docker ps -a -q)
docker volume rm $(docker volume ls -q)
docker compose up -d db
wait 5000
cd ~/Documentos/rinha-de-backend-2024-q1/
./executar-teste-local.sh
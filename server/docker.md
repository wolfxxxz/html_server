version: '3'

services:
  postgres:
    image: postgres:alpine3.18
    container_name: translator_postgres 
    ports:
      - 5435:5437
#    volumes:
#      - ${HOME}/pgdata/:/var/lib/postgresql/data  
    environment:
      POSTGRES_PASSWORD: 1 

  server:
    build:
      dockerfile: Dockerfile
    depends_on:
      - postgres
    env_file:
      - .env
    container_name: server-server
    network_mode: host
    environment:
      - TZ=Europe/Kiev
#    restart: unless-stopped
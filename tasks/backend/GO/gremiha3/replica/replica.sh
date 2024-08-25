#!/bin/bash
rm -r /var/lib/postgresql/data/*
PGPASSWORD="pgpass" pg_basebackup --host=pg-master -D /var/lib/postgresql/data -U repluser -Fp -Xs -P -R

# Устанавливаем права для постгрес
chown -R postgres:postgres /var/lib/postgresql/data

shutdown now
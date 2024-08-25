#!/bin/bash
# Добавляем файл с настройками postgresql.conf
echo "wal_level = replica" >> /var/lib/postgresql/data/postgresql.conf && \
echo "max_wal_senders = 2" >> /var/lib/postgresql/data/postgresql.conf && \
echo "max_replication_slots = 2" >> /var/lib/postgresql/data/postgresql.conf && \
echo "hot_standby = on" >> /var/lib/postgresql/data/postgresql.conf && \
echo "hot_standby_feedback = on" >> /var/lib/postgresql/data/postgresql.conf

# Добавляем запись в pg_hba.conf
echo "host replication all 0.0.0.0/0 trust" >> /var/lib/postgresql/data/pg_hba.conf

# Устанавливаем права для постгрес
chown -R postgres:postgres /var/lib/postgresql/data

shutdown now
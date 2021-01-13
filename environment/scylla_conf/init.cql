CREATE ROLE vaultadmin WITH SUPERUSER = true AND LOGIN = true and PASSWORD = 'vaultpass';

CREATE KEYSPACE IF NOT EXISTS ks_new WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
CREATE TABLE IF NOT EXISTS ks_new.ts_raw (
    tsid text,
    time timestamp,
    value float,
    PRIMARY KEY (tsid, time)
) WITH CLUSTERING ORDER BY (time ASC);
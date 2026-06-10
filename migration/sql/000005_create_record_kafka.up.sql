CREATE TABLE record_kafka (
    id VARCHAR(255) PRIMARY KEY,
    npsn VARCHAR(255),
    sn VARCHAR(255),
    nisn VARCHAR(255),
    timestamp TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_record_kafka_npsn
ON record_kafka(npsn);

CREATE INDEX idx_record_kafka_sn
ON record_kafka(sn);

CREATE INDEX idx_record_kafka_nisn
ON record_kafka(nisn);

CREATE INDEX idx_record_kafka_deleted_at
ON record_kafka(deleted_at);

ALTER TABLE record_kafka
ADD CONSTRAINT fk_record_area
FOREIGN KEY (npsn)
REFERENCES area_kafka(npsn)
ON UPDATE CASCADE;

ALTER TABLE record_kafka
ADD CONSTRAINT fk_record_device
FOREIGN KEY (sn)
REFERENCES device_kafka(sn)
ON UPDATE CASCADE;

ALTER TABLE record_kafka
ADD CONSTRAINT fk_record_user
FOREIGN KEY (nisn)
REFERENCES user_kafka(nisn)
ON UPDATE CASCADE;
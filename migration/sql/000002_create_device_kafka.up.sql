CREATE TABLE device_kafka (
    sn VARCHAR(255) PRIMARY KEY,
    npsn VARCHAR(255) NOT NULL,
    create_at TIMESTAMP,
    brand VARCHAR(255),
    timezone INTEGER,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_device_kafka_npsn
ON device_kafka(npsn);

CREATE INDEX idx_device_kafka_deleted_at
ON device_kafka(deleted_at);

ALTER TABLE device_kafka
ADD CONSTRAINT fk_device_area
FOREIGN KEY (npsn)
REFERENCES area_kafka(npsn)
ON UPDATE CASCADE;
CREATE TABLE area_kafka (
    npsn VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    create_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_area_kafka_deleted_at
ON area_kafka(deleted_at);
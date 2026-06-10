CREATE TABLE user_kafka (
    nisn VARCHAR(255) PRIMARY KEY,
    npsn VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    created_at TIMESTAMP,
    biometric_id VARCHAR(255),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_user_kafka_npsn
ON user_kafka(npsn);

CREATE INDEX idx_user_kafka_deleted_at
ON user_kafka(deleted_at);

ALTER TABLE user_kafka
ADD CONSTRAINT fk_user_area
FOREIGN KEY (npsn)
REFERENCES area_kafka(npsn)
ON UPDATE CASCADE;

ALTER TABLE user_kafka
ADD CONSTRAINT fk_user_biometric
FOREIGN KEY (biometric_id)
REFERENCES biometric_kafka(id)
ON UPDATE CASCADE;
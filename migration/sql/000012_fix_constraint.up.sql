-- device_kafka -> area_kafka

ALTER TABLE device_kafka
DROP CONSTRAINT IF EXISTS fk_device_area;

ALTER TABLE device_kafka
ADD CONSTRAINT fk_device_area
FOREIGN KEY (npsn)
REFERENCES area_kafka(npsn)
ON UPDATE CASCADE
ON DELETE CASCADE;


-- user_kafka -> area_kafka

ALTER TABLE user_kafka
DROP CONSTRAINT IF EXISTS fk_user_area;

ALTER TABLE user_kafka
ADD CONSTRAINT fk_user_area
FOREIGN KEY (npsn)
REFERENCES area_kafka(npsn)
ON UPDATE CASCADE
ON DELETE CASCADE;


-- user_kafka -> biometric_kafka

ALTER TABLE user_kafka
DROP CONSTRAINT IF EXISTS fk_user_biometric;

ALTER TABLE user_kafka
ADD CONSTRAINT fk_user_biometric
FOREIGN KEY (biometric_id)
REFERENCES biometric_kafka(id)
ON UPDATE CASCADE
ON DELETE SET NULL;


-- record_kafka -> area_kafka

ALTER TABLE record_kafka
DROP CONSTRAINT IF EXISTS fk_record_area;

ALTER TABLE record_kafka
ADD CONSTRAINT fk_record_area
FOREIGN KEY (npsn)
REFERENCES area_kafka(npsn)
ON UPDATE CASCADE
ON DELETE CASCADE;


-- record_kafka -> device_kafka

ALTER TABLE record_kafka
DROP CONSTRAINT IF EXISTS fk_record_device;

ALTER TABLE record_kafka
ADD CONSTRAINT fk_record_device
FOREIGN KEY (sn)
REFERENCES device_kafka(sn)
ON UPDATE CASCADE
ON DELETE CASCADE;


-- record_kafka -> user_kafka

ALTER TABLE record_kafka
DROP CONSTRAINT IF EXISTS fk_record_user;

ALTER TABLE record_kafka
ADD CONSTRAINT fk_record_user
FOREIGN KEY (nisn)
REFERENCES user_kafka(nisn)
ON UPDATE CASCADE
ON DELETE CASCADE;
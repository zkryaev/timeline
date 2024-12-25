
CREATE OR REPLACE FUNCTION soft_delete_user()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE users SET is_delete = TRUE WHERE org_id = OLD.org_id;
    DELETE FROM users_verify WHERE user_id = OLD.user_id;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_soft_delete_user
BEFORE DELETE ON users
FOR EACH ROW
EXECUTE FUNCTION soft_delete_user();

CREATE OR REPLACE FUNCTION soft_delete_org()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE orgs SET is_delete = TRUE WHERE org_id = OLD.org_id;
    UPDATE services SET is_delete = TRUE WHERE org_id = OLD.org_id;
    UPDATE workers SET is_delete = TRUE WHERE org_id = OLD.org_id;
    DELETE FROM orgs_verify WHERE org_id = OLD.org_id;
    DELETE FROM timetables WHERE org_id = OLD.org_id;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_soft_delete_org
BEFORE DELETE ON orgs
FOR EACH ROW
EXECUTE FUNCTION soft_delete_org();

CREATE OR REPLACE FUNCTION soft_delete_worker()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE workers SET is_delete = TRUE WHERE worker_id = OLD.worker_id;
    UPDATE worker_schedules SET is_delete = TRUE WHERE worker_id = OLD.worker_id;
    UPDATE worker_services SET is_delete = TRUE WHERE worker_id = OLD.worker_id;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER soft_delete_worker
BEFORE DELETE ON workers
FOR EACH ROW
EXECUTE FUNCTION soft_delete_worker();

CREATE OR REPLACE FUNCTION soft_delete_service()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE services SET is_delete = TRUE WHERE service_id = OLD.service_id;
    UPDATE worker_services SET is_delete = TRUE WHERE service_id = OLD.service_id;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER soft_delete_service
BEFORE DELETE ON services
FOR EACH ROW
EXECUTE FUNCTION soft_delete_service();

CREATE OR REPLACE FUNCTION check_media_user_limit()
RETURNS TRIGGER AS $$
DECLARE
    record_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO record_count FROM media_users WHERE user_id = NEW.user_id;
    IF record_count >= 1 THEN
        RAISE EXCEPTION 'The maximum number of records has been exceeded for user_id %', NEW.user_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_media_user
BEFORE INSERT ON media_users
FOR EACH ROW
EXECUTE FUNCTION check_media_user_limit();

CREATE OR REPLACE FUNCTION check_media_org_limit()
RETURNS TRIGGER AS $$
DECLARE
    record_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO record_count FROM media_orgs WHERE org_id = NEW.org_id;
    IF record_count >= 6 THEN
        RAISE EXCEPTION 'The maximum number of records has been exceeded for org_id %', NEW.org_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_media_orgs
BEFORE INSERT ON media_orgs
FOR EACH ROW
EXECUTE FUNCTION check_media_org_limit();

CREATE OR REPLACE FUNCTION check_media_worker_limit()
RETURNS TRIGGER AS $$
DECLARE
    record_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO record_count FROM media_workers WHERE worker_id = NEW.worker_id;
    IF record_count >= 1 THEN
        RAISE EXCEPTION 'The maximum number of records has been exceeded for worker_id %', NEW.worker_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_media_workers
BEFORE INSERT ON media_workers
FOR EACH ROW
EXECUTE FUNCTION check_media_worker_limit();
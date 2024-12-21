
# Мягкое удаление записей
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

# DELETE orgs -> workers, services
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

# DELETE workers -> worker_schedules, worker_services
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

# Мягкое удаление работников (workers) ведет к мягкому удалению в worker_schedules, worker_services
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
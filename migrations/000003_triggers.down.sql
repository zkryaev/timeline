DROP TRIGGER IF EXISTS trigger_soft_delete_user ON users;

DROP TRIGGER IF EXISTS trigger_soft_delete_org ON orgs;

DROP TRIGGER IF EXISTS soft_delete_worker ON workers;

DROP TRIGGER IF EXISTS soft_delete_service ON services;

DROP TRIGGER IF EXISTS before_insert_showcase ON showcase;
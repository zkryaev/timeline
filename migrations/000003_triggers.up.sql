CREATE OR REPLACE FUNCTION check_media_org_limit()
RETURNS TRIGGER AS $$
DECLARE
    record_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO record_count FROM showcase WHERE org_id = NEW.org_id;
    IF record_count >= 6 THEN
        RAISE EXCEPTION 'The maximum number of records has been exceeded for org_id %', NEW.org_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_showcase
BEFORE INSERT ON showcase
FOR EACH ROW
EXECUTE FUNCTION check_media_org_limit();
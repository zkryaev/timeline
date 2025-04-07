CREATE OR REPLACE FUNCTION check_media_org_limit()
RETURNS TRIGGER AS $$
DECLARE
    media_cnt INTEGER;
    banner_cnt INTEGER;
BEGIN
    IF NEW.org_id IS NULL THEN
        RAISE EXCEPTION 'org_id cannot be NULL';  -- Или обработайте NULL как вам нужно
    END IF;

    SELECT
        COUNT(org_id),
        SUM(CASE WHEN type = 'banner' THEN 1 ELSE 0 END)
    INTO
        media_cnt,
        banner_cnt
    FROM showcase
    WHERE org_id = NEW.org_id::INT;

    IF NEW.type = 'banner' AND banner_cnt >= 1 THEN
        RAISE EXCEPTION 'the banner already exists for org_id %', NEW.org_id;
    END IF;

    IF media_cnt > 5 THEN
        RAISE EXCEPTION 'the image limit has been reached for org_id %', NEW.org_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_showcase
BEFORE INSERT ON showcase
FOR EACH ROW
EXECUTE FUNCTION check_media_org_limit();
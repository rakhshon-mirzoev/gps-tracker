DO $$
BEGIN
    -- 1. Пересоздаём индекс: исключаем approved/rejected из "активных" визитов
    IF EXISTS (
        SELECT 1 FROM pg_indexes WHERE indexname = 'uniq_active_visit'
    ) THEN
        DROP INDEX uniq_active_visit;
    END IF;

    -- 2. Обновляем CHECK-constraint: добавляем статус 'revision'
    IF EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'visits_review_status_check'
          AND conrelid = 'visits'::regclass
    ) THEN
        ALTER TABLE visits DROP CONSTRAINT visits_review_status_check;
    END IF;

    ALTER TABLE visits
        ADD CONSTRAINT visits_review_status_check
            CHECK (review_status IN ('new', 'pending', 'approved', 'rejected', 'revision', 'schedule_deleted'));
END$$;

CREATE UNIQUE INDEX uniq_active_visit
    ON visits(staff_id)
    WHERE check_out_at IS NULL
      AND review_status NOT IN ('approved', 'rejected');

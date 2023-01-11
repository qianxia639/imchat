ALTER TABLE IF EXISTS "contacts" DROP CONSTRAINT IF EXISTS "contacts_owner_id_fkey";

ALTER TABLE IF EXISTS "contacts" DROP CONSTRAINT IF EXISTS "contacts_target_id_fkey";

DROP TABLE IF EXISTS "contacts";
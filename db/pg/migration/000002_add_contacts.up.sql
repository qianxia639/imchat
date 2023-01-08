CREATE TABLE "contact" (
  "id" bigserial PRIMARY KEY,
  "owner_id" bigint NOT NULL,
  "target_id" bigint NOT NULL,
  "type" smallint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);


COMMENT ON COLUMN "contact"."id" IS '主键Id';

COMMENT ON COLUMN "contact"."owner_id" IS '谁的关系';

COMMENT ON COLUMN "contact"."target_id" IS '对应的谁';

COMMENT ON COLUMN "contact"."type" IS '对应类型, 1: 好友, 2: 群组';

COMMENT ON COLUMN "contact"."created_at" IS '创建时间';

ALTER TABLE "contact" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("id");

ALTER TABLE "contact" ADD FOREIGN KEY ("target_id") REFERENCES "users" ("id");
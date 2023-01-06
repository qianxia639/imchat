CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar(20) UNIQUE NOT NULL,
  "email" varchar(30) UNIQUE NOT NULL,
  "nickname" varchar(20) UNIQUE NOT NULL,
  "password" varchar(100) NOT NULL,
  "gender" smallint NOT NULL DEFAULT 3,
  "avatar" varchar(255) NOT NULL DEFAULT ('default.jpg'),
  "register_time" timestamptz NOT NULL DEFAULT (now())
);

COMMENT ON COLUMN "users"."id" IS '主键Id';

COMMENT ON COLUMN "users"."username" IS '用户名';

COMMENT ON COLUMN "users"."email" IS '用户邮箱';

COMMENT ON COLUMN "users"."nickname" IS '用户昵称';

COMMENT ON COLUMN "users"."password" IS '密码';

COMMENT ON COLUMN "users"."gender" IS '1: 男, 2: 女, 3: 未知';

COMMENT ON COLUMN "users"."avatar" IS '用户头像';

COMMENT ON COLUMN "users"."register_time" IS '注册时间';
CREATE TABLE "messages" (
  "id" bigserial PRIMARY KEY,
  "from_id" bigint NOT NULL,
  "to_id" bigint NOT NULL,
  "transfer_type" smallint NOT NULL DEFAULT 0,
  "message_type" smallint NOT NULL DEFAULT 0,
  "content" varchar NOT NULL,
  "pic" varchar,
  "url" varchar,
  "send_time" bigint,
  "recv_time" bigint
);

CREATE INDEX ON "messages" ("from_id");

COMMENT ON COLUMN "messages"."id" IS '主键Id';

COMMENT ON COLUMN "messages"."from_id" IS '发送者';

COMMENT ON COLUMN "messages"."to_id" IS '接收者';

COMMENT ON COLUMN "messages"."transfer_type" IS '消息传输类型, 0: 私聊, 1: 群聊, 2: 公告, 3: 心跳';

COMMENT ON COLUMN "messages"."message_type" IS '消息内容类型, 0:  文字, 1:表情符号, 2: 图片, 3: 文件, 4: 语音, 5: 视频 , 6: 语音通话, 7: 视频通话';

COMMENT ON COLUMN "messages"."content" IS '消息内容';

COMMENT ON COLUMN "messages"."pic" IS '缩略图';

COMMENT ON COLUMN "messages"."url" IS '文件或图片地址';

COMMENT ON COLUMN "messages"."send_time" IS '发送时间';

COMMENT ON COLUMN "messages"."recv_time" IS '接收时间';
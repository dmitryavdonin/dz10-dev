CREATE TABLE "billing" (
    "id" serial primary key,
    "user_id" integer,
    "balance" integer,
    "created_at" timestamp not null,
    "modified_at" timestamp not null
)
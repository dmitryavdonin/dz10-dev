CREATE TABLE "billing_account" (
    "id" serial primary key,
    "user_id" integer,
    "balance" integer,
    "created_at" timestamp not null,
    "modified_at" timestamp not null
);

CREATE TABLE "billing_transaction" (
    "id" serial primary key,
    "user_id" integer,
    "order_id" integer,
    "operation" varchar,
    "amount" int,
    "status" varchar,
    "reason" varchar,
    "created_at" timestamp not null,
    "modified_at" timestamp not null
);
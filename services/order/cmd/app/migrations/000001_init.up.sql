CREATE TABLE "order" (
    "id" serial primary key,
    "user_id" integer not null,
    "price" integer,
    "status" varchar,
    "reason" varchar,    
    "created_at" timestamp not null,
    "modified_at" timestamp not null
);
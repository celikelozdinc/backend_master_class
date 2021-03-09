CREATE TABLE "students" (
    "id" bigserial PRIMARY KEY,
    "name" varchar NOT NULL,
    "number" bigint NOT NULL,
    "nation" varchar NOT NULL,
    "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "grades" (
    "id" bigserial PRIMARY KEY,
    "student_id" bigint NOT NULL,
    "grade" bigint NOT NULL,
    "created_at" timestamptz DEFAULT (now())
);

ALTER TABLE "grades" ADD FOREIGN KEY ("student_id") REFERENCES "students" ("id");

CREATE INDEX ON "students" ("name");

CREATE INDEX ON "grades" ("student_id");
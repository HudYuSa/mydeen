CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
    IF NOT EXISTS "masters" (
        "master_id" UUID NOT NULL DEFAULT (uuid_generate_v4()),
        "email" VARCHAR(50) UNIQUE,
        "password" VARCHAR(255) NOT NULL,
        "verified" BOOLEAN DEFAULT false,
        "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT "masters_pkey" PRIMARY KEY("master_id")
    );

CREATE TABLE
    IF NOT EXISTS "verification_codes" (
        "verification_code_id" UUID NOT NULL DEFAULT (uuid_generate_v4()),
        "master_id" UUID NOT NULL,
        "totp_code" VARCHAR(255) NOT NULL,
        "expire_date" TIMESTAMP NOT NULL,
        "used" BOOLEAN DEFAULT false,
        "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT "verification_code_pkey" PRIMARY KEY ("verification_code_id"),
        CONSTRAINT "fk_master" FOREIGN KEY("master_id") REFERENCES "masters"("master_id") ON DELETE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS "invitations" (
        "invitation_id" UUID NOT NULL DEFAULT (uuid_generate_v4()),
        "master_id" UUID NOT NULL,
        "code" VARCHAR(255) NOT NULL,
        "expire_date" TIMESTAMP NOT NULL,
        "used" BOOLEAN DEFAULT false,
        "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT "invitations_pkey" PRIMARY KEY ("invitation_id"),
        CONSTRAINT "fk_master" FOREIGN KEY("master_id") REFERENCES "masters"("master_id") ON DELETE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS "admins" (
        "admin_id" UUID NOT NULL DEFAULT (uuid_generate_v4()),
        "invitation_id" UUID NOT NULL,
        "username" varchar(50) NOT NULL,
        "email" varchar(50) UNIQUE,
        "password" varchar(255) NOT NULL,
        "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT "admins_pkey" PRIMARY KEY ("admin_id"),
        CONSTRAINT "fk_invitation" FOREIGN KEY("invitation_id") REFERENCES "invitations"("invitation_id") ON DELETE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS "events" (
        "event_id" UUID NOT NULL DEFAULT (uuid_generate_v4()),
        "admin_id" UUID NOT NULL,
        "event_name" VARCHAR(50) NOT NULL,
        "status" VARCHAR(50) NOT NULL DEFAULT "scheduled" "moderation" BOOLEAN DEFAULT false,
        "max_questions" INTEGER NOT NULL DEFAULT 1,
        "max_question_length" INTEGER NOT NULL DEFAULT 160,
        "event_code" VARCHAR(50) NOT NULL,
        "start_date" TIMESTAMP NOT NULL,
        "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT "events_pkey" PRIMARY KEY ("event_id"),
        CONSTRAINT "fk_admin" FOREIGN KEY("admin_id") REFERENCES "admins"("admin_id") ON DELETE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS "users" (
        "user_id" UUID NOT NULL DEFAULT (uuid_generate_v4()),
        "username" VARCHAR(255),
        "email" VARCHAR(255),
        "ip_address" INET NOT NULL,
        CONSTRAINT "users_pkey" PRIMARY KEY ("user_id")
    );

CREATE TABLE
    IF NOT EXISTS "questions"(
        "question_id" UUID NOT NULL DEFAULT (uuid_generate_v4()),
        "event_id" UUID NOT NULL,
        "user_id" UUID NOT NULL,
        "content" TEXT NOT NULL,
        "starred" BOOLEAN DEFAULT false,
        "approved" BOOLEAN DEFAULT false,
        "answered" BOOLEAN DEFAULT false,
        "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT "question_pkey" PRIMARY KEY ("question_id"),
        CONSTRAINT "fk_event" FOREIGN KEY("event_id") REFERENCES "events"("event_id") ON DELETE CASCADE,
        CONSTRAINT "fk_user" FOREIGN KEY("user_id") REFERENCES "users"("user_id") ON DELETE CASCADE
    );
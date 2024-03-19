CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "masters"(
    "master_id" uuid NOT NULL DEFAULT (uuid_generate_v4()),
    "email" varchar(50) UNIQUE,
    "password" varchar(255) NOT NULL,
    "verified" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "masters_pkey" PRIMARY KEY ("master_id")
);

CREATE TABLE IF NOT EXISTS "verification_codes"(
    "verification_code_id" uuid NOT NULL DEFAULT (uuid_generate_v4()),
    "master_id" uuid NOT NULL,
    "code" varchar(255) NOT NULL,
    "expire_date" timestamp NOT NULL,
    "used" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "verification_codes_pkey" PRIMARY KEY ("verification_code_id"),
    CONSTRAINT "fk_master" FOREIGN KEY ("master_id") REFERENCES "masters"("master_id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "master_otps"(
    "master_otp_id" uuid NOT NULL DEFAULT (uuid_generate_v4()),
    "master_id" uuid NOT NULL,
    "code" varchar(10) NOT NULL,
    "expire_date" timestamp NOT NULL,
    "used" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "master_otps_pkey" PRIMARY KEY ("master_otp_id"),
    CONSTRAINT "fk_master" FOREIGN KEY ("master_id") REFERENCES "masters"("master_id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "invitations"(
    "invitation_id" uuid NOT NULL DEFAULT (uuid_generate_v4()),
    "master_id" uuid,
    "code" varchar(255) NOT NULL,
    "expire_date" timestamp NOT NULL,
    "used" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "invitations_pkey" PRIMARY KEY ("invitation_id"),
    CONSTRAINT "fk_master" FOREIGN KEY ("master_id") REFERENCES "masters"("master_id") ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS "admins"(
    "admin_id" uuid NOT NULL DEFAULT (uuid_generate_v4()),
    "invitation_id" uuid,
    "username" varchar(50) NOT NULL,
    "email" varchar(50) UNIQUE,
    "password" varchar(255) NOT NULL,
    "admin_code" varchar(50) UNIQUE,
    "enable_2fa" boolean NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "admins_pkey" PRIMARY KEY ("admin_id"),
    CONSTRAINT "fk_invitation" FOREIGN KEY ("invitation_id") REFERENCES "invitations"("invitation_id") ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS "admin_otps"(
    "admin_otp_id" uuid NOT NULL DEFAULT (uuid_generate_v4()),
    "admin_id" uuid NOT NULL,
    "code" varchar(10) NOT NULL,
    "expire_date" timestamp NOT NULL,
    "used" boolean NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "admin_otps_pkey" PRIMARY KEY ("admin_otp_id"),
    CONSTRAINT "fk_admin" FOREIGN KEY ("admin_id") REFERENCES "admins"("admin_id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "events"(
    "event_id" uuid NOT NULL DEFAULT (uuid_generate_v4()),
    "admin_id" uuid NOT NULL,
    "event_name" varchar(250) NOT NULL,
    "status" varchar(50) NOT NULL,
    "moderation" boolean,
    "max_questions" integer NOT NULL,
    "max_question_length" integer NOT NULL,
    "event_code" varchar(50) UNIQUE,
    "start_date" timestamp NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "events_pkey" PRIMARY KEY ("event_id"),
    CONSTRAINT "fk_admin" FOREIGN KEY ("admin_id") REFERENCES "admins"("admin_id") ON DELETE CASCADE,
    CONSTRAINT "valid_status" CHECK ("status" IN ('scheduled', 'live', 'finished')),
    CONSTRAINT "valid_max_questions" CHECK ("max_questions" IN (1, 3, 5)),
    CONSTRAINT "valid_max_question_length" CHECK ("max_question_length" IN (160, 240, 360, 480))
);

CREATE TABLE IF NOT EXISTS "questions"(
    "question_id" uuid NOT NULL DEFAULT (uuid_generate_v4()),
    "event_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "username" varchar(50),
    "content" text NOT NULL,
    "starred" boolean,
    "approved" boolean,
    "answered" boolean,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "question_pkey" PRIMARY KEY ("question_id"),
    CONSTRAINT "fk_event" FOREIGN KEY ("event_id") REFERENCES "events"("event_id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "likes"(
    "like_id" uuid NOT NULL DEFAULT (uuid_generate_v4()),
    "question_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    CONSTRAINT "like_pkey" PRIMARY KEY ("like_id"),
    CONSTRAINT "fk_question" FOREIGN KEY ("question_id") REFERENCES "questions"("question_id") ON DELETE CASCADE,
    CONSTRAINT "unique_like_user_question" UNIQUE ("question_id", "user_id")
);


SET search_path=public;

CREATE TABLE public.policies (

);
ALTER TABLE public.policies OWNER TO watchdog;

CREATE TABLE public.repositories (
  "id" uuid NOT NULL,
  "created_by" uuid,
  "last_analysis" TIMESTAMP,
  "repository_url" VARCHAR(255)
);
ALTER TABLE public.repositories OWNER TO watchdog;
ALTER TABLE public.repositories ADD CONSTRAINT repositories_repository_url_key UNIQUE (repository_url);

CREATE TABLE public.repositories_issues (
  "id" uuid NOT NULL,
  "author" VARCHAR(255),
  "commit_hash" VARCHAR(255) NOT NULL,
  "created_at" TIMESTAMP,
  "email" VARCHAR(255),
  "file" VARCHAR(255),
  "line" VARCHAR(255),
  "line_number" INTEGER,
  "analysis_id" uuid NOT NULL,
  "rule" VARCHAR(255),
  "severity" VARCHAR(255) NOT NULL
);
ALTER TABLE public.repositories_issues OWNER TO watchdog;
ALTER TABLE public.repositories_issues ADD CONSTRAINT repositories_issues_analysis_id_commit_hash_key UNIQUE (analysis_id, commit_hash);

CREATE TABLE public.repositories_analyzes (
  "id" uuid NOT NULL,
  "created_at" TIMESTAMP,
  "duration" BIGINT,
  "finished_at" TIMESTAMP,
  "repository_id" uuid NOT NULL,
  "severity" VARCHAR(255) NOT NULL,
  "started_at" TIMESTAMP,
  "state" VARCHAR(255) NOT NULL,
  "total_issues" BIGINT DEFAULT 0
);
ALTER TABLE public.repositories_analyzes OWNER TO watchdog;

CREATE TABLE public.rules (
  "id" SERIAL NOT NULL,
  "created_at" TIMESTAMP,
  "display_name" CHARACTER varying(128),
  "description" CHARACTER varying(1024),
  "enabled" BOOLEAN DEFAULT true NOT NULL,
  "file" VARCHAR(255) NOT NULL,
  "severity" VARCHAR(255) NOT NULL,
  "tags" CHARACTER varying(512) DEFAULT ''
);
ALTER TABLE public.rules OWNER TO watchdog;

CREATE TABLE public.rules_allowed_entries (
  "id" SERIAL NOT NULL,
  "rule_id" uuid NOT NULL
);
ALTER TABLE public.rules_allowed_entries OWNER TO watchdog;

CREATE TABLE public.rules_entropies (
  "id" SERIAL NOT NULL,
  "rule_id" uuid NOT NULL
);
ALTER TABLE public.rules_entropies OWNER TO watchdog;

CREATE TABLE public.users (
  "id" uuid NOT NULL,
  "created_at" TIMESTAMP,
  "created_by" uuid,
  "email" VARCHAR(128),
  "first_name" VARCHAR(64),
  "last_name" VARCHAR(64),
  "locked" BOOLEAN DEFAULT false NOT NULL,
  "locked_at" TIMESTAMP,
  "password" VARCHAR(128),
  "state" SMALLINT DEFAULT 1 NOT NULL,
  "updated_at" TIMESTAMP,
  "updated_by" uuid
);
ALTER TABLE public.users OWNER TO watchdog;
CREATE UNIQUE INDEX index_users_email ON users USING btree (email);

ALTER TABLE ONLY public.repositories ADD CONSTRAINT repositories_primary_key PRIMARY KEY (id);
ALTER TABLE ONLY public.repositories_issues ADD CONSTRAINT repositories_issues_primary_key PRIMARY KEY (id);
ALTER TABLE ONLY public.repositories_analyzes ADD CONSTRAINT repositories_analyzes_primary_key PRIMARY KEY (id);
ALTER TABLE ONLY public.users ADD CONSTRAINT users_primary_key PRIMARY KEY (id);

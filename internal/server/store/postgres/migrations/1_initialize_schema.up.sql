SET search_path=public;

CREATE TABLE public.access_tokens (
  "id" BIGINT NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  "expires_at" DATE,
  "last_used_at" TIMESTAMP WITH TIME ZONE,
  "name" CHARACTER varying NOT NULL,
  "revoked" BOOLEAN DEFAULT false,
  "scopes" CHARACTER varying DEFAULT '--- []'::CHARACTER varying NOT NULL,
  "token" CHARACTER varying NOT NULL,
  "user_id" uuid NOT NULL
);
ALTER TABLE public.access_tokens OWNER TO watchdog;
CREATE SEQUENCE access_tokens_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER SEQUENCE access_tokens_id_seq OWNED BY access_tokens.id;

CREATE TABLE public.analytics (
	"container_id" uuid NOT NULL,
	"container_type" TEXT,
	"container_count" BIGINT,
	"window_type" TEXT,
	"window_size" INTEGER,
	"window_value" TIMESTAMP
);
ALTER TABLE public.analytics OWNER TO watchdog;
ALTER TABLE public.analytics ADD CONSTRAINT analytics_unique_index UNIQUE (container_id, window_value);

CREATE TABLE public.analytics_rollups (
  "rollup_name" TEXT PRIMARY KEY,
  "container_type" TEXT NOT NULL,
  "container_id_sequence_name" TEXT NOT NULL,
  "last_aggregated_id" BIGINT DEFAULT 0
);
ALTER TABLE public.analytics_rollups OWNER TO watchdog;

CREATE TABLE public.categories (
  "id" BIGINT NOT NULL,
  "description" TEXT,
  "extension" VARCHAR(255) NOT NULL,
  "level" INTEGER DEFAULT 0 NOT NULL,
  "lft" BIGINT DEFAULT 0 NOT NULL,
  "rgt" BIGINT DEFAULT 0 NOT NULL,
  "parent_id" INTEGER,
  "title" VARCHAR(255) NOT NULL,
  "value" VARCHAR(255)
);
ALTER TABLE public.categories OWNER TO watchdog;
CREATE SEQUENCE categories_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER SEQUENCE categories_id_seq OWNED BY categories.id;

CREATE TABLE public.integrations (
  "id" SERIAL NOT NULL,
  "api_token" VARCHAR(255),
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  "created_by" uuid NOT NULL,
  "enable_monitoring" BOOLEAN DEFAULT false NOT NULL,
  "instance_name" VARCHAR(255),
  "instance_type" VARCHAR(255),
  "instance_url" VARCHAR(255),
  "synced_at" TIMESTAMP WITH TIME ZONE,
  "syncing_error" TEXT DEFAULT ''
);
ALTER TABLE public.integrations OWNER TO watchdog;
ALTER TABLE public.integrations ADD CONSTRAINT integrations_id_key UNIQUE (id);

CREATE TABLE public.integrations_webhooks (
  "id" SERIAL NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  "group_id" BIGINT,
  "integration_id" BIGINT NOT NULL,
  "token" CHARACTER varying NOT NULL,
  "url" VARCHAR(255),
  "webhook_id" BIGINT NOT NULL
);
ALTER TABLE public.integrations_webhooks OWNER TO watchdog;
ALTER TABLE public.integrations_webhooks ADD CONSTRAINT integrations_webhooks_id_key UNIQUE (id);
ALTER TABLE ONLY public.integrations_webhooks ADD CONSTRAINT integrations_webhooks_foreign_key FOREIGN KEY (integration_id) REFERENCES public.integrations(id) ON DELETE CASCADE;

CREATE TABLE public.jobs (
  "id" SERIAL NOT NULL,
  "args" JSON NOT NULL DEFAULT '[]'::json,
  "error_count" INTEGER NOT NULL DEFAULT 0,
  "last_error" TEXT,
  "priority" SMALLINT NOT NULL DEFAULT 100,
  "queue" TEXT NOT NULL DEFAULT '',
  "started_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "type" TEXT NOT NULL
);
ALTER TABLE public.jobs OWNER TO watchdog;

CREATE TABLE public.policies (
  "id" SERIAL NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  "created_by" uuid,
  "description" CHARACTER varying(1024),
  "enabled" BOOLEAN DEFAULT false NOT NULL,
  "display_name" VARCHAR(255),
  "name" VARCHAR(255),
  "severity" VARCHAR(32) NOT NULL,
  "type" VARCHAR(255) NOT NULL
);
ALTER TABLE public.policies OWNER TO watchdog;
ALTER TABLE public.policies ADD CONSTRAINT policies_id_key UNIQUE (id);

CREATE TABLE public.policies_conditions (
  "id" SERIAL NOT NULL,
  "pattern" CHARACTER varying(1024),
  "policy_id" BIGINT NOT NULL,
  "rejection_message" TEXT,
  "skip" BOOLEAN,
  "type" VARCHAR(255)
);
ALTER TABLE public.policies_conditions OWNER TO watchdog;

CREATE TABLE public.repositories (
  "id" uuid NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  "created_by" uuid,
  "enable_monitoring" BOOLEAN DEFAULT false NOT NULL,
  "integration_id" BIGINT,
  "last_analysis" TIMESTAMP,
  "repository_url" VARCHAR(255),
  "visibility" VARCHAR(255)
);
ALTER TABLE public.repositories OWNER TO watchdog;
ALTER TABLE public.repositories ADD CONSTRAINT repositories_repository_url_key UNIQUE (repository_url);
ALTER TABLE ONLY public.repositories ADD CONSTRAINT repositories_integrations_foreign_key FOREIGN KEY (integration_id) REFERENCES public.integrations(id) ON DELETE CASCADE;

CREATE TABLE public.repositories_analyzes (
  "id" uuid NOT NULL,
  "created_at" TIMESTAMP,
  "created_by" uuid,
  "duration" BIGINT,
  "finished_at" TIMESTAMP,
  "last_commit_hash" VARCHAR(255),
  "last_commit_date" TIMESTAMP,
  "repository_id" uuid NOT NULL,
  "severity" VARCHAR(255) NOT NULL,
  "started_at" TIMESTAMP,
  "state" VARCHAR(255) NOT NULL,
  "state_message" TEXT DEFAULT '',
  "total_issues" BIGINT DEFAULT 0,
  "trigger" VARCHAR(255) NOT NULL
);
ALTER TABLE public.repositories_analyzes OWNER TO watchdog;

CREATE TABLE public.repositories_issues (
  "id" uuid NOT NULL,
  "analysis_id" uuid NOT NULL,
  "author" VARCHAR(255),
  "commit_hash" VARCHAR(255) NOT NULL,
  "condition_type" VARCHAR(255) NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  "email" VARCHAR(255),
  "offender_object" VARCHAR(255) DEFAULT '',
  "offender_operand" VARCHAR(255) DEFAULT '',
  "offender_operator" VARCHAR(255) DEFAULT '',
  "offender_value" TEXT DEFAULT '',
  "policy_id" BIGINT NOT NULL,
  "policy_type" TEXT,
  "repository_id" uuid NOT NULL,
  "severity" VARCHAR(255) NOT NULL
);
ALTER TABLE public.repositories_issues OWNER TO watchdog;
ALTER TABLE public.repositories_issues ADD CONSTRAINT repositories_issues_analysis_id_commit_hash_key UNIQUE (id, analysis_id, commit_hash);

CREATE TABLE public.repositories_leaks (
  "id" SERIAL NOT NULL,
  "analysis_id" uuid NOT NULL,
  "author_email" VARCHAR(255),
  "author_name" VARCHAR(255),
  "commit_hash" VARCHAR(255) NOT NULL,
  "content" TEXT,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  "file" VARCHAR(255),
  "line" TEXT,
  "line_number" INTEGER,
  "offender" TEXT,
  "repository_id" uuid NOT NULL,
  "rule_id" BIGINT NOT NULL,
  "secret_hash" VARCHAR(255),
  "secret_revoked" BOOLEAN DEFAULT false NOT NULL,
  "severity" VARCHAR(255) NOT NULL
);
ALTER TABLE public.repositories_leaks OWNER TO watchdog;
ALTER TABLE public.repositories_leaks ADD CONSTRAINT secret_hash UNIQUE (secret_hash);

CREATE TABLE public.repositories_statistics (
  "id" BIGINT NOT NULL,
  "analysis_count" BIGINT DEFAULT 0 NOT NULL,
  "issue_count" BIGINT DEFAULT 0 NOT NULL,
  "leak_count" BIGINT DEFAULT 0 NOT NULL,
  "repository_id" uuid NOT NULL
);
ALTER TABLE public.repositories_statistics OWNER TO watchdog;

CREATE TABLE public.rules (
  "id" SERIAL NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  "created_by" uuid,
  "display_name" CHARACTER varying(128) DEFAULT '',
  "description" CHARACTER varying(1024) DEFAULT '',
  "enabled" BOOLEAN DEFAULT true NOT NULL,
  "file" VARCHAR(255) DEFAULT '',
  "name" CHARACTER varying(128),
  "pattern" VARCHAR(255) NOT NULL,
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

CREATE TABLE public.settings (
  "id" BIGINT NOT NULL,
  "container_id" BIGINT,
  "container_type" VARCHAR(255) NOT NULL,
  "setting_key" VARCHAR(255) NOT NULL,
  "setting_type" VARCHAR(255) NOT NULL,
  "setting_value" VARCHAR(255) NOT NULL
);
ALTER TABLE public.settings OWNER TO watchdog;
CREATE SEQUENCE settings_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER SEQUENCE settings_id_seq OWNED BY settings.id;

CREATE TABLE public.users (
  "id" uuid NOT NULL,
  "created_at" TIMESTAMP NOT NULL,
  "created_by" uuid,
  "email" VARCHAR(128) NOT NULL,
  "first_name" VARCHAR(64) NOT NULL,
  "last_login" TIMESTAMP,
  "last_name" VARCHAR(64) NOT NULL,
  "locked" BOOLEAN DEFAULT false NOT NULL,
  "locked_at" TIMESTAMP,
  "password" VARCHAR(128) NOT NULL,
  "provider" VARCHAR(128) DEFAULT 'local',
  "roles" VARCHAR(255),
  "state" SMALLINT DEFAULT 1 NOT NULL,
  "updated_at" TIMESTAMP,
  "updated_by" uuid,
  "username" VARCHAR(128)
);
ALTER TABLE public.users OWNER TO watchdog;
ALTER TABLE public.users ADD CONSTRAINT users_email_key UNIQUE (email);
ALTER TABLE ONLY public.users ADD CONSTRAINT users_primary_key PRIMARY KEY (id);

CREATE TABLE public.whitelist (
  "id" BIGINT NOT NULL,
  "commit" TEXT,
  "container_id" uuid NULL,
	"container_type" TEXT,
  "description" TEXT,
  "files" TEXT,
  "paths" TEXT,
  "repositories" TEXT
);
ALTER TABLE public.whitelist OWNER TO watchdog;
CREATE SEQUENCE whitelist_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER SEQUENCE whitelist_id_seq OWNED BY whitelist.id;

ALTER TABLE ONLY public.jobs ADD CONSTRAINT jobs_primary_key PRIMARY KEY (id, "queue", priority, started_at);
ALTER TABLE ONLY public.policies ADD CONSTRAINT policies_primary_key PRIMARY KEY (id);
ALTER TABLE ONLY public.policies ADD CONSTRAINT policies_created_by_foreign_key FOREIGN KEY ("created_by") REFERENCES public.users(id);
ALTER TABLE ONLY public.policies_conditions ADD CONSTRAINT policies_conditions_primary_key PRIMARY KEY (id);
ALTER TABLE ONLY public.policies_conditions ADD CONSTRAINT policies_conditions_foreign_key FOREIGN KEY (policy_id) REFERENCES public.policies(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.repositories ADD CONSTRAINT repositories_primary_key PRIMARY KEY (id);
ALTER TABLE ONLY public.repositories ADD CONSTRAINT repositories_created_by_foreign_key FOREIGN KEY ("created_by") REFERENCES public.users(id);
ALTER TABLE ONLY public.repositories_analyzes ADD CONSTRAINT repositories_analyzes_primary_key PRIMARY KEY (id);
ALTER TABLE ONLY public.repositories_analyzes ADD CONSTRAINT repositories_analyzes_foreign_key FOREIGN KEY (repository_id) REFERENCES public.repositories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.repositories_analyzes ADD CONSTRAINT repositories_analyzes_created_by_foreign_key FOREIGN KEY ("created_by") REFERENCES public.users(id);
ALTER TABLE ONLY public.repositories_issues ADD CONSTRAINT repositories_issues_primary_key PRIMARY KEY (id);
ALTER TABLE ONLY public.repositories_issues ADD CONSTRAINT repositories_issues_foreign_key FOREIGN KEY (repository_id) REFERENCES public.repositories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.repositories_issues ADD CONSTRAINT repositories_issues_analyzes_foreign_key FOREIGN KEY (analysis_id) REFERENCES public.repositories_analyzes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.repositories_leaks ADD CONSTRAINT repositories_leaks_foreign_key FOREIGN KEY (repository_id) REFERENCES public.repositories(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.rules ADD CONSTRAINT rules_primary_key PRIMARY KEY (id);
ALTER TABLE ONLY public.rules ADD CONSTRAINT rules_created_by_foreign_key FOREIGN KEY ("created_by") REFERENCES public.users(id);

CREATE OR REPLACE FUNCTION analytics_aggregation(OUT start_id BIGINT, OUT end_id BIGINT)
RETURNS record AS $$
BEGIN
  SELECT window_start, window_end INTO start_id, end_id
  FROM incremental_rollup_window('analytics');

  IF start_id > end_id THEN RETURN; END IF;
	INSERT INTO "analytics"
		SELECT
			"repositories_leaks"."repository_id" AS container_id,
			'leak' container_type,
			count(*) AS container_count,
			'day' window_type,
			'1'::INTEGER windows_size,
			date_trunc('day', "repositories_leaks"."created_at") AS window_value
		FROM "repositories_leaks"
    WHERE "repositories_leaks"."id" BETWEEN start_id AND end_id
		GROUP BY container_id, window_value
	ON CONFLICT (container_id, window_value)
	DO UPDATE
	SET container_count = analytics.container_count + excluded.container_count;
	RETURN;
END;
$$LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION incremental_rollup_window(name TEXT, OUT window_start BIGINT, OUT window_end BIGINT)
RETURNS record AS $$
DECLARE
    table_to_lock regclass;
BEGIN
    SELECT container_type, last_aggregated_id+1, pg_sequence_last_value(container_id_sequence_name)
    INTO table_to_lock, window_start, window_end
    FROM analytics_rollups
    WHERE rollup_name = name FOR UPDATE;
    IF NOT FOUND THEN
        RAISE 'rollup ''%'' is not in the rollups table', rollup_name;
    END IF;
    IF window_end IS NULL THEN
        window_end := 0;
        RETURN;
    END IF;
    BEGIN
        EXECUTE format('LOCK %s IN EXCLUSIVE MODE', table_to_lock);
        RAISE 'release table lock';
    EXCEPTION WHEN OTHERS THEN
    END;
    UPDATE analytics_rollups SET last_aggregated_id = window_end WHERE rollup_name = name;
END;
$$LANGUAGE plpgsql;

INSERT INTO analytics_rollups (rollup_name, container_type, container_id_sequence_name)
VALUES ('analytics', 'repositories_leaks','repositories_leaks_id_seq');

SELECT analytics_aggregation();

--- Seeding categories
INSERT INTO public.categories("id", "lft", "rgt", "level", "title", "value", "extension") VALUES
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Branch', 'branch', 'handler_type'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Commit', 'commit', 'handler_type'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'File', 'file', 'handler_type'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Jira', 'jira', 'handler_type'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Security', 'security', 'handler_type'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Tag', 'tag', 'handler_type'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Email', 'email', 'condition_type'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Extension', 'extension', 'condition_type'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'IP Address', 'ip', 'condition_type'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Issue', 'issue', 'condition_type'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Length', 'length', 'condition_type'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Pattern', 'pattern', 'condition_type'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Protected', 'protected', 'condition_type'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Secret', 'secret', 'condition_type'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Semver', 'semver', 'condition_type'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Signature', 'signature', 'condition_type'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Size', 'size', 'condition_type'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Low', '0', 'issue_severity'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Medium', '1', 'issue_severity'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'High', '2', 'issue_severity'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Blocker', 'BLOCKER', 'rule_severity'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Critical', 'CRITICAL', 'rule_severity'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Info', 'INFO', 'rule_severity'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Major', 'MAJOR', 'rule_severity'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Minor', 'MINOR', 'rule_severity'),
  (NEXTVAL('categories_id_seq'), 1, 2, 0, 'Gitlab', 'gitlab', 'integration_type');

--- Seeding policies
INSERT INTO public.policies("id", "description", "display_name", "name", "severity", "type") VALUES (
  NEXTVAL('policies_id_seq'),
  'Ensure no secrets are committed on your repositories',
  'Secrets detection',
  'secrets',
  'BLOCKER',
  'security'
);
INSERT INTO public.policies_conditions("policy_id", "type") VALUES (
  CURRVAL('policies_id_seq'),
  'secret'
);
---
INSERT INTO public.policies("id", "description", "display_name", "name", "severity", "type") VALUES (
  NEXTVAL('policies_id_seq'),
  'Ensure a .gitignore file is present in your repositories',
  '.gitignore file',
  'gitignore',
  'CRITICAL',
  'file'
);
INSERT INTO public.policies_conditions("policy_id", "type", "pattern") VALUES (
  CURRVAL('policies_id_seq'),
  'exist',
  '.gitignore'
);
---
INSERT INTO public.policies("id", "description", "display_name", "name", "severity", "type") VALUES (
  NEXTVAL('policies_id_seq'),
  'Enforce branch naming',
  'Git Flow',
  'git_flow',
  'MINOR',
  'branch'
);
INSERT INTO public.policies_conditions("policy_id", "type", "pattern", "rejection_message") VALUES (
  CURRVAL('policies_id_seq'),
  'pattern',
  '(feature|release|hotfix)\/[a-z\d-_.]+',
  'Branch `{{ .Branch }}` must match Gitflow naming convention'
);
---
INSERT INTO public.policies("id", "description", "display_name", "name", "severity", "type") VALUES (
  NEXTVAL('policies_id_seq'),
  'Ensure that commit messages meet the conventional commit format.',
  'Conventional commit',
  'conventional_commit',
  'MINOR',
  'commit'
);
INSERT INTO public.policies_conditions("policy_id", "type", "pattern", "rejection_message") VALUES (
  CURRVAL('policies_id_seq'),
  'pattern',
  '(?m)^(build|ci|docs|feat|fix|perf|refactor|style|test)\([a-z]+\):\s([a-z\.\-\s]+)',
  'Message must be formatted like type(scope): subject'
);
---
INSERT INTO public.policies("id", "description", "display_name", "name", "severity", "type") VALUES (
  NEXTVAL('policies_id_seq'),
  'Ensure no files with sensitive extensions are commited (ex: .key, .cert)',
  'File Extensions',
  'file_extensions',
  'MAJOR',
  'file'
);
INSERT INTO public.policies_conditions("policy_id", "type", "pattern") VALUES (
  CURRVAL('policies_id_seq'),
  'extension',
  'key'
);
---
INSERT INTO public.policies("id", "description", "display_name", "name", "severity", "type") VALUES (
  NEXTVAL('policies_id_seq'),
  'Ensure no large files are commited',
  'File Size',
  'file_size',
  'MAJOR',
  'file'
);
INSERT INTO public.policies_conditions("policy_id", "type", "pattern", "rejection_message") VALUES (
  CURRVAL('policies_id_seq'),
  'size',
  'lt 1mb',
  'File {{ .Object }} size {{ .Value }} greater or equal than {{ .Operand }}'
);

--- Seeding rules
INSERT INTO public.rules("id", "description", "display_name", "name", "pattern", "severity", "tags") VALUES
(
  NEXTVAL('rules_id_seq'),
  '',
  'Asymmetric Private Key',
  'ASYMMETRIC_PRIVATE_KEY',
  '(\-){5}BEGIN[[:blank:]]*?(RSA|OPENSSH|DSA|EC|PGP)?[[:blank:]]*?PRIVATE[[:blank:]]KEY[[:blank:]]*?(BLOCK)?(\-){5}.*',
  'CRITICAL',
  '["Key", "Asymmetric Private Key"]'
), (
  NEXTVAL('rules_id_seq'),
  '',
  'AWS Keys',
  'AWS_ACCESS_KEY',
  '(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}',
  'MAJOR',
  '["Cloud Provider"]'
), (
  NEXTVAL('rules_id_seq'),
  '',
  'SendGrid API Key',
  'SENDGRID',
  'SG\.[\w_]{16,32}\.[\w_]{16,64}',
  'MAJOR',
  '["Key", "Send Grid"]'
), (
  NEXTVAL('rules_id_seq'),
  '',
  'Stripe API key',
  'STRIPE',
  '(?i)stripe(.{0,20})?[''\"][sk|rk]_live_[0-9a-zA-Z]{24}',
  'MAJOR',
  '["Key", "Stripe"]'
);

--- Seeding settings
INSERT INTO public.settings("id", "container_id", "container_type", "setting_key", "setting_type", "setting_value") VALUES
  (NEXTVAL('settings_id_seq'), NULL, 'global', 'enable_signup', 'boolean', 'true'),
  (NEXTVAL('settings_id_seq'), NULL, 'global', 'enable_oauth_signup', 'boolean', 'false');

--- Seeding allowed items
INSERT INTO public.whitelist("id", "paths") VALUES
  (NEXTVAL('whitelist_id_seq'), '["node_modules", "vendor"]');

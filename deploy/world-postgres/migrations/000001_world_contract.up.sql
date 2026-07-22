BEGIN;
CREATE SCHEMA IF NOT EXISTS world AUTHORIZATION aese_world_app;
CREATE TABLE world.schema_migrations (version bigint PRIMARY KEY, applied_at timestamptz NOT NULL DEFAULT now());
CREATE TABLE world.runs (
  tenant_id text NOT NULL, world_run_id text NOT NULL, branch_id text NOT NULL DEFAULT 'main',
  contract jsonb NOT NULL, created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (tenant_id, world_run_id, branch_id)
);
CREATE TABLE world.events (
  tenant_id text NOT NULL, world_run_id text NOT NULL, branch_id text NOT NULL,
  sequence bigint NOT NULL CHECK (sequence >= 0), event_id text NOT NULL, idempotency_key text NOT NULL,
  contract jsonb NOT NULL, recorded_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (tenant_id, world_run_id, branch_id, sequence),
  UNIQUE (tenant_id, event_id), UNIQUE (tenant_id, world_run_id, branch_id, idempotency_key),
  FOREIGN KEY (tenant_id, world_run_id, branch_id) REFERENCES world.runs(tenant_id, world_run_id, branch_id)
);
CREATE TABLE world.snapshots (
  tenant_id text NOT NULL, world_run_id text NOT NULL, branch_id text NOT NULL,
  through_sequence bigint NOT NULL CHECK (through_sequence >= 0), state_hash text NOT NULL, contract jsonb NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (tenant_id, world_run_id, branch_id, through_sequence),
  FOREIGN KEY (tenant_id, world_run_id, branch_id) REFERENCES world.runs(tenant_id, world_run_id, branch_id)
);
CREATE TABLE world.knowledge (tenant_id text NOT NULL, world_run_id text NOT NULL, branch_id text NOT NULL, knowledge_id text NOT NULL, contract jsonb NOT NULL, PRIMARY KEY (tenant_id, knowledge_id));
CREATE TABLE world.discrepancies (tenant_id text NOT NULL, world_run_id text NOT NULL, branch_id text NOT NULL, discrepancy_id text NOT NULL, contract jsonb NOT NULL, PRIMARY KEY (tenant_id, discrepancy_id));
CREATE TABLE world.bridge_checkpoints (tenant_id text NOT NULL, world_run_id text NOT NULL, branch_id text NOT NULL, last_iaos_cursor bigint NOT NULL CHECK (last_iaos_cursor >= 0), updated_at timestamptz NOT NULL DEFAULT now(), PRIMARY KEY (tenant_id, world_run_id, branch_id));
ALTER TABLE world.runs ENABLE ROW LEVEL SECURITY;
ALTER TABLE world.runs FORCE ROW LEVEL SECURITY;
ALTER TABLE world.events ENABLE ROW LEVEL SECURITY;
ALTER TABLE world.events FORCE ROW LEVEL SECURITY;
ALTER TABLE world.snapshots ENABLE ROW LEVEL SECURITY;
ALTER TABLE world.snapshots FORCE ROW LEVEL SECURITY;
ALTER TABLE world.knowledge ENABLE ROW LEVEL SECURITY;
ALTER TABLE world.knowledge FORCE ROW LEVEL SECURITY;
ALTER TABLE world.discrepancies ENABLE ROW LEVEL SECURITY;
ALTER TABLE world.discrepancies FORCE ROW LEVEL SECURITY;
ALTER TABLE world.bridge_checkpoints ENABLE ROW LEVEL SECURITY;
ALTER TABLE world.bridge_checkpoints FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON world.runs USING (tenant_id = current_setting('app.tenant_id', true)) WITH CHECK (tenant_id = current_setting('app.tenant_id', true));
CREATE POLICY tenant_isolation ON world.events USING (tenant_id = current_setting('app.tenant_id', true)) WITH CHECK (tenant_id = current_setting('app.tenant_id', true));
CREATE POLICY tenant_isolation ON world.snapshots USING (tenant_id = current_setting('app.tenant_id', true)) WITH CHECK (tenant_id = current_setting('app.tenant_id', true));
CREATE POLICY tenant_isolation ON world.knowledge USING (tenant_id = current_setting('app.tenant_id', true)) WITH CHECK (tenant_id = current_setting('app.tenant_id', true));
CREATE POLICY tenant_isolation ON world.discrepancies USING (tenant_id = current_setting('app.tenant_id', true)) WITH CHECK (tenant_id = current_setting('app.tenant_id', true));
CREATE POLICY tenant_isolation ON world.bridge_checkpoints USING (tenant_id = current_setting('app.tenant_id', true)) WITH CHECK (tenant_id = current_setting('app.tenant_id', true));
GRANT USAGE ON SCHEMA world TO aese_world_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON world.runs, world.events, world.snapshots, world.knowledge, world.discrepancies, world.bridge_checkpoints TO aese_world_app;
GRANT SELECT ON world.schema_migrations TO aese_world_app;
INSERT INTO world.schema_migrations(version) VALUES (1);
COMMIT;

DO $$ BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'aese_world_app') THEN
    CREATE ROLE aese_world_app LOGIN PASSWORD 'aese_world_local_only';
  END IF;
END $$;
GRANT CONNECT ON DATABASE aese_world TO aese_world_app;

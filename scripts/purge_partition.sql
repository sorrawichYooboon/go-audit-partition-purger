CREATE OR REPLACE PROCEDURE purge_old_audit_partition(partition_name TEXT)
LANGUAGE plpgsql
AS $$
BEGIN
    EXECUTE format('ALTER TABLE audit_logs DETACH PARTITION %I', partition_name);
    EXECUTE format('DROP TABLE IF EXISTS %I', partition_name);

    RAISE NOTICE 'Partition % has been successfully purged.', partition_name;
EXCEPTION
    WHEN undefined_table THEN
        RAISE NOTICE 'Partition % does not exist.', partition_name;
    WHEN OTHERS THEN
        RAISE EXCEPTION 'Failed to purge partition %: %', partition_name, SQLERRM;
END;
$$;

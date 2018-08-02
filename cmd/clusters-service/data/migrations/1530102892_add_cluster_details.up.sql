ALTER TABLE clusters
ADD COLUMN name text NOT NULL,
ADD COLUMN region text,
ADD COLUMN master_nodes int,
ADD COLUMN infra_nodes int, 
ADD COLUMN compute_nodes int, 
ADD COLUMN memory int, 
ADD COLUMN cpu_cores int,
ADD COLUMN storage int, 
ADD COLUMN state text;


ALTER TABLE projects 
ADD COLUMN env_variables JSONB DEFAULT '[]'::jsonb;

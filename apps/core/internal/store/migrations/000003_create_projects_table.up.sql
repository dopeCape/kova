CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    name VARCHAR(50) NOT NULL,
    user_id TEXT NOT NULL,
    repo_id BIGINT NOT NULL,
    repo_name VARCHAR(255) NOT NULL,
    repo_full_name VARCHAR(255) NOT NULL,
    repo_url TEXT NOT NULL,
    repo_branch VARCHAR(255) DEFAULT 'main',
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_projects_user_id 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT unique_user_project_name 
        UNIQUE(user_id, name),
        
    CONSTRAINT projects_name_length 
        CHECK (char_length(name) >= 1),
    CONSTRAINT projects_name_format 
        CHECK (name ~ '^[a-zA-Z0-9_-]+$'),
    CONSTRAINT projects_repo_id_positive 
        CHECK (repo_id > 0),
    CONSTRAINT projects_status_valid 
        CHECK (status IN ('active', 'inactive', 'archived')),
    CONSTRAINT projects_repo_url_format 
        CHECK (repo_url ~ '^https://github\.com/')
);

CREATE INDEX IF NOT EXISTS idx_projects_user_id ON projects(user_id);
CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(user_id, name);
CREATE INDEX IF NOT EXISTS idx_projects_repo_id ON projects(repo_id);
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
CREATE INDEX IF NOT EXISTS idx_projects_created_at ON projects(created_at DESC);

CREATE TRIGGER update_projects_updated_at 
    BEFORE UPDATE ON projects 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE projects (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    name VARCHAR(50) NOT NULL,
    user_id TEXT NOT NULL,
    repo_id BIGINT NOT NULL,
    repo_name VARCHAR(255) NOT NULL,
    repo_full_name VARCHAR(255) NOT NULL,
    repo_url TEXT NOT NULL,
    repo_branch VARCHAR(255) DEFAULT 'main',
    status VARCHAR(20) DEFAULT 'active',
    env_variables JSONB DEFAULT '[]'::jsonb,
    deployment_status VARCHAR(20) DEFAULT 'pending',
    domain TEXT,
    port INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    UNIQUE(user_id, name),
    CHECK (char_length(name) >= 1),
    CHECK (name ~ '^[a-zA-Z0-9_-]+$'),
    CHECK (repo_id > 0),
    CHECK (status IN ('active', 'inactive', 'archived')),
    CHECK (deployment_status IN ('pending', 'building', 'deploying', 'deployed', 'failed')),
    CHECK (repo_url ~ '^https://github\.com/'),
    CHECK (port IS NULL OR (port >= 8000 AND port <= 9000))
);

-- Create indexes for performance
CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_projects_name ON projects(user_id, name);
CREATE INDEX idx_projects_repo_id ON projects(repo_id);
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_deployment_status ON projects(deployment_status);
CREATE INDEX idx_projects_port ON projects(port);
CREATE INDEX idx_projects_created_at ON projects(created_at DESC);

-- Create trigger for updated_at
CREATE TRIGGER update_projects_updated_at 
    BEFORE UPDATE ON projects 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

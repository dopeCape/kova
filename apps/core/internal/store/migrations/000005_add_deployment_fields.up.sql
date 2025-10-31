ALTER TABLE projects
ADD COLUMN deployment_status VARCHAR(20) DEFAULT 'pending',
ADD COLUMN domain TEXT,
ADD COLUMN port INTEGER;

-- Add constraints
ALTER TABLE projects ADD CONSTRAINT check_deployment_status CHECK (
  deployment_status IN (
    'pending',
    'building',
    'deploying',
    'deployed',
    'failed'
  )
);

-- Add index for deployment status queries
CREATE INDEX idx_projects_deployment_status ON projects (deployment_status);

-- Add index for port (to quickly find used ports)
CREATE INDEX idx_projects_port ON projects (port);

CREATE TABLE accounts (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    user_id TEXT NOT NULL,
    github_username VARCHAR(255) NOT NULL,
    github_id BIGINT UNIQUE NOT NULL,
    avatar_url TEXT,
    access_token TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_accounts_github_id ON accounts(github_id);
CREATE INDEX idx_accounts_github_username ON accounts(github_username);
CREATE INDEX idx_accounts_created_at ON accounts(created_at);

CREATE TRIGGER update_accounts_updated_at 
    BEFORE UPDATE ON accounts 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

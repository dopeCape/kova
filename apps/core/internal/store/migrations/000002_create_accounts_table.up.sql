CREATE TABLE IF NOT EXISTS accounts (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    user_id TEXT NOT NULL,
    github_username VARCHAR(255) NOT NULL,
    github_id BIGINT UNIQUE NOT NULL,
    avatar_url TEXT,
    access_token TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraint
    CONSTRAINT fk_accounts_user_id 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE,
    
    -- Unique constraint to prevent duplicate GitHub accounts per user
    CONSTRAINT unique_user_github_account 
        UNIQUE(user_id, github_id),
    
    CONSTRAINT accounts_github_username_length 
        CHECK (char_length(github_username) >= 1),
    CONSTRAINT accounts_github_username_format 
        CHECK (github_username ~ '^[a-zA-Z0-9]([a-zA-Z0-9-])*[a-zA-Z0-9]$|^[a-zA-Z0-9]$'),
    CONSTRAINT accounts_github_id_positive 
        CHECK (github_id > 0),
    CONSTRAINT accounts_avatar_url_format 
        CHECK (avatar_url IS NULL OR avatar_url ~ '^https?://'),
    CONSTRAINT accounts_access_token_length 
        CHECK (char_length(access_token) >= 40)
);

CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_accounts_github_id ON accounts(github_id);
CREATE INDEX IF NOT EXISTS idx_accounts_github_username ON accounts(github_username);
CREATE INDEX IF NOT EXISTS idx_accounts_created_at ON accounts(created_at DESC);

CREATE TRIGGER update_accounts_updated_at 
    BEFORE UPDATE ON accounts 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

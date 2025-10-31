-- name: CreateProject :one
INSERT INTO projects (name, user_id, repo_id, repo_name, repo_full_name, repo_url, repo_branch, status, env_variables, deployment_status, domain, port)
VALUES ($1, $2, $3, $4, $5, $6, $7, 'active', $8, 'pending', $9, $10)
RETURNING id, name, user_id, repo_id, repo_name, repo_full_name, repo_url, repo_branch, status, env_variables, deployment_status, domain, port, created_at, updated_at;

-- name: GetProjectByID :one
SELECT id, name, user_id, repo_id, repo_name, repo_full_name, repo_url, repo_branch, status, env_variables, deployment_status, domain, port, created_at, updated_at
FROM projects
WHERE id = $1;

-- name: GetProjectByUserIDAndName :one
SELECT id, name, user_id, repo_id, repo_name, repo_full_name, repo_url, repo_branch, status, env_variables, deployment_status, domain, port, created_at, updated_at
FROM projects
WHERE user_id = $1 AND name = $2;

-- name: GetProjectsByUserID :many
SELECT id, name, user_id, repo_id, repo_name, repo_full_name, repo_url, repo_branch, status, env_variables, deployment_status, domain, port, created_at, updated_at
FROM projects
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetProjectsByUserIDAndStatus :many
SELECT id, name, user_id, repo_id, repo_name, repo_full_name, repo_url, repo_branch, status, env_variables, deployment_status, domain, port, created_at, updated_at
FROM projects
WHERE user_id = $1 AND status = $2
ORDER BY created_at DESC;

-- name: GetProjectsByRepoID :many
SELECT id, name, user_id, repo_id, repo_name, repo_full_name, repo_url, repo_branch, status, env_variables, deployment_status, domain, port, created_at, updated_at
FROM projects
WHERE repo_id = $1
ORDER BY created_at DESC;

-- name: UpdateProject :one
UPDATE projects
SET name = $2, repo_branch = $3, status = $4, domain = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, user_id, repo_id, repo_name, repo_full_name, repo_url, repo_branch, status, env_variables, deployment_status, domain, port, created_at, updated_at;

-- name: UpdateProjectStatus :one
UPDATE projects
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, user_id, repo_id, repo_name, repo_full_name, repo_url, repo_branch, status, env_variables, deployment_status, domain, port, created_at, updated_at;

-- name: UpdateProjectDeploymentStatus :one
UPDATE projects
SET deployment_status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, user_id, repo_id, repo_name, repo_full_name, repo_url, repo_branch, status, env_variables, deployment_status, domain, port, created_at, updated_at;

-- name: UpdateProjectBranch :one
UPDATE projects
SET repo_branch = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, user_id, repo_id, repo_name, repo_full_name, repo_url, repo_branch, status, env_variables, deployment_status, domain, port, created_at, updated_at;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1;

-- name: DeleteProjectsByUserID :exec
DELETE FROM projects
WHERE user_id = $1;

-- name: ListProjects :many
SELECT id, name, user_id, repo_id, repo_name, repo_full_name, repo_url, repo_branch, status, env_variables, deployment_status, domain, port, created_at, updated_at
FROM projects
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListProjectsByStatus :many
SELECT id, name, user_id, repo_id, repo_name, repo_full_name, repo_url, repo_branch, status, env_variables, deployment_status, domain, port, created_at, updated_at
FROM projects
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountProjects :one
SELECT COUNT(*) FROM projects;

-- name: CountProjectsByUserID :one
SELECT COUNT(*) FROM projects WHERE user_id = $1;

-- name: CountProjectsByStatus :one
SELECT COUNT(*) FROM projects WHERE status = $1;

-- name: ProjectExistsByUserIDAndName :one
SELECT EXISTS(SELECT 1 FROM projects WHERE user_id = $1 AND name = $2);

-- name: ProjectExistsByID :one
SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1);

-- name: SearchProjectsByUserID :many
SELECT id, name, user_id, repo_id, repo_name, repo_full_name, repo_url, repo_branch, status, env_variables, deployment_status, domain, port, created_at, updated_at
FROM projects
WHERE user_id = $1 AND (
    name ILIKE '%' || $2 || '%' 
    OR repo_name ILIKE '%' || $2 || '%'
    OR repo_full_name ILIKE '%' || $2 || '%'
)
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: SearchProjects :many
SELECT id, name, user_id, repo_id, repo_name, repo_full_name, repo_url, repo_branch, status, env_variables, deployment_status, domain, port, created_at, updated_at
FROM projects
WHERE name ILIKE '%' || $1 || '%' 
   OR repo_name ILIKE '%' || $1 || '%'
   OR repo_full_name ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetActiveProjectsByUserID :many
SELECT id, name, user_id, repo_id, repo_name, repo_full_name, repo_url, repo_branch, status, env_variables, deployment_status, domain, port, created_at, updated_at
FROM projects
WHERE user_id = $1 AND status = 'active'
ORDER BY created_at DESC;

-- name: ArchiveProject :one
UPDATE projects
SET status = 'archived', updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, user_id, repo_id, repo_name, repo_full_name, repo_url, repo_branch, status, env_variables, deployment_status, domain, port, created_at, updated_at;

-- name: ActivateProject :one
UPDATE projects
SET status = 'active', updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, user_id, repo_id, repo_name, repo_full_name, repo_url, repo_branch, status, env_variables, deployment_status, domain, port, created_at, updated_at;

-- name: GetUsedPorts :many
SELECT port FROM projects WHERE port IS NOT NULL ORDER BY port ASC;

-- name: UpdateProjectPort :exec
UPDATE projects
SET port = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetUsers :many
SELECT 
    id, role_id, wallet_address, email, name, surname, is_default, created_at
FROM 
    t_users
WHERE
    (sqlc.narg(id)::UUID IS NULL OR id = sqlc.narg(id)::UUID) AND
    (sqlc.narg(email)::TEXT IS NULL OR email = sqlc.narg(email)::TEXT) AND
    (sqlc.narg(wallet_address)::TEXT IS NULL OR wallet_address = sqlc.narg(wallet_address)::TEXT)
LIMIT @lim OFFSET @off;

-- name: GetUserByID :one
SELECT 
    id, role_id, wallet_address, email, name, surname, is_default, created_at
FROM 
    t_users
WHERE 
    id = @user_id;

-- name: GetUserByWalletAddress :one
SELECT 
    id, role_id, wallet_address, email, name, surname, is_default, created_at
FROM 
    t_users 
WHERE 
    wallet_address = @wallet_address;

-- name: GetUserByEmail :one
SELECT 
    id, role_id, wallet_address, email, name, surname, is_default, created_at
FROM 
    t_users 
WHERE 
    wallet_address = @wallet_address;

-- name: GetDefaultUser :one
SELECT 
    id, role_id, wallet_address, email, name, surname, is_default, created_at
FROM 
    t_users 
WHERE 
    is_default = true;

-- name: CreateUser :one
INSERT INTO t_users 
    (role_id, name, surname, email, wallet_address, is_default, created_at)
VALUES 
    (@role_id, @name, @surname, @email, @wallet_address, @is_default, NOW())
RETURNING id;

-- name: UpdateUser :exec
UPDATE
    t_users
SET
    name = COALESCE(@name, name),
    surname = COALESCE(@surname, surname),
    email = COALESCE(@email, email)
WHERE
    id = @user_id;

-- name: ChangeUserRole :exec
UPDATE
    t_users
SET
    role_id = COALESCE(@role_id, role_id)
WHERE
    id = @user_id;

-- name: DeleteUser :exec
DELETE FROM 
    t_users
WHERE
    id = @user_id;

-- name: CountUserByWalletAddress :one
SELECT COUNT(*) 
FROM t_users 
WHERE wallet_address = @wallet_address;

-- name: IsDefaultUserExists :one
SELECT 
    EXISTS (
        SELECT 1 
        FROM t_users u
        WHERE u.is_default = true
    ) AS is_default_user_exists;
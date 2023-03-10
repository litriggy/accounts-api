-- name: GetOauth :one
SELECT * FROM oauth
WHERE oauth_id = ? AND oauth_type = ?;

-- name: GetUserByOauth :one
SELECT * FROM users
WHERE id in (
    SELECT user_id FROM oauth
    WHERE oauth_type = ? AND oauth_id = ?
);

-- name: CreateOauth :execresult
INSERT INTO oauth (
    `user_id`, `oauth_id`, `oauth_type`, `version`
) VALUES (
    ?, ?, ?, ?
);

-- name: DeleteOauthByUserID :exec
DELETE FROM oauth
WHERE user_id = ?;

-- name: DeleteOauthByOauth :exec
DELETE FROM oauth
WHERE oauth_id = ? AND oauth_type = ?;

-- name: CreateUser :execresult
INSERT INTO users (
    `nickname`, `email`, `type`
) VALUES (
    ?, ?, ?
);

-- name: GetUser :one
SELECT * FROM users
WHERE id = ?;

-- name: UpdateNickname :exec
UPDATE users
SET nickname = ?
WHERE id = ?;

-- name: FindUserByEmail :many
SELECT `id`, `nickname`, `email`  FROM users
WHERE LOWER(`email`) = LOWER(?) AND `is_locked` = 0;

-- name: FindUserByNickname :many
SELECT `id`, `nickname`, `email` FROM users
WHERE LOWER(`nickname`) = LOWER(?) AND `is_locked` = 0;

-- name: CreateService :exec
INSERT INTO services (
    `name`, `symbol`, `decimals`, `image`, `is_native`, `contract_addr`
) VALUES (
    ?, ?, ?, ?, ?, ? 
);

-- name: DeleteService :exec
DELETE FROM services
WHERE id = ?;

-- name: GetBalances :many
SELECT bal.amount, s.name, s.symbol, s.contract_addr, s.is_native, s.image, s.decimals FROM user_balances AS bal
LEFT JOIN services AS s
ON bal.service_id = s.id
WHERE bal.user_id = ?;

-- name: UpdateBalance :exec
INSERT INTO user_balances (
    `user_id`, `service_id`, `amount`
) VALUES (
    ?, ?, ?
) ON DUPLICATE KEY UPDATE
amount = amount + VALUES(amount);

-- name: GetBalance :one
SELECT bal.amount, s.name, s.symbol, s.contract_addr, s.is_native, s.image, s.decimals FROM user_balances AS bal
LEFT JOIN services AS s
ON bal.service_id = s.id
WHERE bal.user_id = ? AND bal.service_id = ?;

-- name: CreateTx :execresult
INSERT INTO transactions (
    `from_id`, `to_id`, `service_id`, `memo`, `total_amount`
) VALUES (
    ?, ?, ?, ?, ?
);

-- name: CreateTxDetails :execresult
INSERT INTO transaction_detail (
    `transaction_id`, `from`, `to`, `amount`, `is_onchain`, `txhash`
) VALUES (
    ?, ?, ?, ?, ?, ?
);

-- name: UpdateTxHash :exec
UPDATE transaction_detail
SET txhash = ?
WHERE `transaction_id` = ? AND `from` = ?;

-- name: UpdateTxTotalAmount :exec
UPDATE transactions
SET total_amount = ?
WHERE `id` = ?;

-- name: GetServiceData :one
SELECT is_native, contract_addr, net_type FROM services
WHERE id = ?;

-- name: GetAllServices :many
SELECT * FROM services;

-- name: GetUserServices :many
SELECT us.service_id AS service_id, ub.amount AS amount, s.name, s.symbol, s.decimals, s.image, s.is_native, s.contract_addr, s.net_type 
FROM user_service AS us
LEFT JOIN services AS s
ON us.service_id = s.id
LEFT JOIN user_balances AS ub
ON ub.service_id = s.id AND ub.user_id = us.user_id
WHERE us.user_id = ?;

-- name: UserAddService :exec
INSERT INTO user_service (
    user_id, service_id
) VALUES (
    ?, ?
);

-- name: UserRemoveService :exec
DELETE FROM user_service
WHERE user_id = ? AND service_id = ?;

-- name: CreateUserWallet :exec
INSERT INTO user_wallets (
    `user_id`, `wallet_addr`, `sec_pk`, `is_integrated`
) VALUES (
    ?, ?, ?, ?
);

-- name: GetUserWallets :many
SELECT wallet_addr, is_integrated FROM user_wallets
WHERE user_id = ?;

-- name: GetPKFromWallet :one
SELECT sec_pk FROM user_wallets
WHERE wallet_addr = ?;

-- name: CheckUserWallet :one
SELECT * FROM user_wallets
WHERE user_id = ? AND wallet_addr = ?;

-- name: CreateSecondPassword :exec
INSERT INTO user_pw (
    user_id, sec_pw
) VALUES (
    ?, ?
)
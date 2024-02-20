-- name: getUser
SELECT
    *
FROM
    USERS
WHERE
    username = $1;
-- LeetCode 262 - Trips and Users
--
-- Pseudo code:
--   JOIN Trips with Users twice (client + driver) filtering banned = 'No'
--   GROUP BY day; cancellation rate = non-completed trips / total trips
--   ROUND to 2 decimal places

SELECT
    t.Request_at AS Day,
    ROUND(
        SUM(CASE WHEN t.Status != 'completed' THEN 1 ELSE 0 END) / COUNT(*),
        2
    ) AS 'Cancellation Rate'
FROM Trips t
JOIN Users u1 ON t.Client_Id = u1.Users_Id AND u1.Banned = 'No'
JOIN Users u2 ON t.Driver_Id = u2.Users_Id AND u2.Banned = 'No'
WHERE t.Request_at BETWEEN '2013-10-01' AND '2013-10-03'
GROUP BY t.Request_at
ORDER BY t.Request_at;

-- LeetCode 262 - Trips and Users
--
-- Problem:
--   The Trips table holds all taxi trips. The Users table holds all users
--   with a Banned flag. Find the cancellation rate (rounded to 2 decimals)
--   for each day between '2013-10-01' and '2013-10-03', excluding trips
--   where the client OR driver is banned.
--
-- Example:
--   Trips: id=1, client=1, driver=10, status='completed',  date='2013-10-01'
--          id=2, client=2, driver=11, status='cancelled_by_driver', date='2013-10-01'
--          id=3, client=3, driver=12, status='completed',  date='2013-10-02'
--   Users: id=1 (not banned), id=2 (banned), id=3 (not banned)
--          id=10,11,12 (not banned)
--
--   Output:
--   Day          | Cancellation Rate
--   2013-10-01   | 0.00   (trip 2 excluded because client 2 is banned)
--   2013-10-02   | 0.00
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

SELECT * FROM moz_origins
WHERE frecency IS NOT NULL
ORDER BY frecency DESC
LIMIT 10;

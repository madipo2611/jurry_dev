ALTER TABLE posts DROP CONSTRAINT IF EXISTS posts_userid_fkey;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS posts;
-- Удаляем индексы
DROP INDEX IF EXISTS idx_notifications_user;
DROP INDEX IF EXISTS idx_dm_recipient;
DROP INDEX IF EXISTS idx_dm_sender;
DROP INDEX IF EXISTS idx_comment_likes_comment;
DROP INDEX IF EXISTS idx_comments_post;
DROP INDEX IF EXISTS idx_post_likes_post;
DROP INDEX IF EXISTS idx_subscriptions_target;
DROP INDEX IF EXISTS idx_subscriptions_subscriber;
DROP INDEX IF EXISTS idx_posts_userID;

-- Удаляем таблицы
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS direct_messages;
DROP TABLE IF EXISTS comment_likes;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS post_likes;
DROP TABLE IF EXISTS subscriptions;

-- Откатываем изменения в posts и users
ALTER TABLE posts
DROP COLUMN IF EXISTS video,
    DROP COLUMN IF EXISTS likes_count,
    DROP COLUMN IF EXISTS comments_count,
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS privacy_level;

ALTER TABLE users
DROP COLUMN IF EXISTS email,
    DROP COLUMN IF EXISTS avatar,
    DROP COLUMN IF EXISTS bio;
-- Добавляем новые столбцы в users
ALTER TABLE users ADD COLUMN email VARCHAR(100) UNIQUE;
ALTER TABLE users ADD COLUMN avatar VARCHAR(150);
ALTER TABLE users ADD COLUMN bio TEXT;

-- Добавляем новые столбцы в posts
ALTER TABLE posts ADD COLUMN likes_count INTEGER NOT NULL DEFAULT 0;
ALTER TABLE posts ADD COLUMN comments_count INTEGER NOT NULL DEFAULT 0;
ALTER TABLE posts ADD COLUMN privacy_level SMALLINT NOT NULL DEFAULT 1;

-- Таблица подписок
CREATE TABLE subscriptions (
                               id SERIAL PRIMARY KEY,
                               subscriberID INTEGER NOT NULL,
                               targetID INTEGER NOT NULL,
                               created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                               status VARCHAR(20) DEFAULT 'active',
                               FOREIGN KEY (subscriberID) REFERENCES users(id) ON DELETE CASCADE,
                               FOREIGN KEY (targetID) REFERENCES users(id) ON DELETE CASCADE,
                               UNIQUE (subscriberID, targetID)
);

-- Таблица лайков постов
CREATE TABLE post_likes (
                            id SERIAL PRIMARY KEY,
                            userID INTEGER NOT NULL,
                            postID INTEGER NOT NULL,
                            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            FOREIGN KEY (userID) REFERENCES users(id) ON DELETE CASCADE,
                            FOREIGN KEY (postID) REFERENCES posts(id) ON DELETE CASCADE,
                            UNIQUE (userID, postID)
);

-- Таблица комментариев
CREATE TABLE comments (
                          id SERIAL PRIMARY KEY,
                          userID INTEGER NOT NULL,
                          postID INTEGER NOT NULL,
                          parentID INTEGER,
                          text TEXT NOT NULL,
                          likes_count INTEGER NOT NULL DEFAULT 0,
                          created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP,
                          is_deleted BOOLEAN DEFAULT FALSE,
                          FOREIGN KEY (userID) REFERENCES users(id) ON DELETE CASCADE,
                          FOREIGN KEY (postID) REFERENCES posts(id) ON DELETE CASCADE,
                          FOREIGN KEY (parentID) REFERENCES comments(id) ON DELETE CASCADE
);

-- Таблица лайков комментариев
CREATE TABLE comment_likes (
                               id SERIAL PRIMARY KEY,
                               userID INTEGER NOT NULL,
                               commentID INTEGER NOT NULL,
                               created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                               FOREIGN KEY (userID) REFERENCES users(id) ON DELETE CASCADE,
                               FOREIGN KEY (commentID) REFERENCES comments(id) ON DELETE CASCADE,
                               UNIQUE (userID, commentID)
);

-- Таблица сообщений
CREATE TABLE direct_messages (
                                 id SERIAL PRIMARY KEY,
                                 senderID INTEGER NOT NULL,
                                 recipientID INTEGER NOT NULL,
                                 message TEXT NOT NULL,
                                 is_read BOOLEAN DEFAULT FALSE,
                                 created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                 FOREIGN KEY (senderID) REFERENCES users(id) ON DELETE CASCADE,
                                 FOREIGN KEY (recipientID) REFERENCES users(id) ON DELETE CASCADE
);

-- Таблица уведомлений
CREATE TABLE notifications (
                               id SERIAL PRIMARY KEY,
                               userID INTEGER NOT NULL,
                               fromUserID INTEGER NOT NULL,
                               type VARCHAR(20) NOT NULL,
                               entityID INTEGER,
                               is_read BOOLEAN DEFAULT FALSE,
                               created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                               FOREIGN KEY (userID) REFERENCES users(id) ON DELETE CASCADE,
                               FOREIGN KEY (fromUserID) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_posts_userID ON posts(userID);
CREATE INDEX idx_subscriptions_subscriber ON subscriptions(subscriberID);
CREATE INDEX idx_subscriptions_target ON subscriptions(targetID);
CREATE INDEX idx_post_likes_post ON post_likes(postID);
CREATE INDEX idx_comments_post ON comments(postID);
CREATE INDEX idx_comment_likes_comment ON comment_likes(commentID);
CREATE INDEX idx_dm_sender ON direct_messages(senderID);
CREATE INDEX idx_dm_recipient ON direct_messages(recipientID);
CREATE INDEX idx_notifications_user ON notifications(userID);
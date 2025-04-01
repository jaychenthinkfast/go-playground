// 切换到 playground 数据库
db = db.getSiblingDB('playground');

// 创建 shares 集合
db.createCollection('shares');

// 创建索引
db.shares.createIndex({ "shareId": 1 }, { unique: true });
db.shares.createIndex({ "expires_at": 1 }, { expireAfterSeconds: 0 });

// 创建视图统计集合
db.createCollection('share_stats');
db.share_stats.createIndex({ "shareId": 1 }, { unique: true }); 
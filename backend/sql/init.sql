CREATE DATABASE IF NOT EXISTS cocktail_bar DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE cocktail_bar;

DROP TABLE IF EXISTS purchase_items;
DROP TABLE IF EXISTS purchases;
DROP TABLE IF EXISTS waste_records;
DROP TABLE IF EXISTS special_creations;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS recipe_ingredients;
DROP TABLE IF EXISTS recipes;
DROP TABLE IF EXISTS ingredients;
DROP TABLE IF EXISTS spirits;

CREATE TABLE spirits (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    brand VARCHAR(100),
    origin VARCHAR(100),
    alcohol_content DECIMAL(4,1),
    volume_ml INT NOT NULL DEFAULT 700,
    unit VARCHAR(20) NOT NULL DEFAULT '瓶',
    stock_quantity INT NOT NULL DEFAULT 0,
    min_stock INT NOT NULL DEFAULT 5,
    cost_price DECIMAL(10,2) NOT NULL,
    selling_price_per_ml DECIMAL(10,4),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_category (category),
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='基酒库存表';

CREATE TABLE ingredients (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    stock_quantity DECIMAL(10,2) NOT NULL DEFAULT 0,
    min_stock DECIMAL(10,2) NOT NULL DEFAULT 0,
    cost_price DECIMAL(10,2) NOT NULL,
    supplier VARCHAR(100),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_category (category),
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='配料表';

CREATE TABLE recipes (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    glass_type VARCHAR(50),
    serving_ml INT,
    price DECIMAL(10,2) NOT NULL,
    cost DECIMAL(10,2),
    preparation_method TEXT,
    garnish VARCHAR(200),
    taste_profile VARCHAR(200),
    difficulty VARCHAR(20),
    is_signature BOOLEAN DEFAULT FALSE,
    image_url VARCHAR(500),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_category (category),
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='调酒配方表';

CREATE TABLE recipe_ingredients (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    recipe_id BIGINT NOT NULL,
    ingredient_type VARCHAR(20) NOT NULL,
    ingredient_id BIGINT NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_recipe (recipe_id),
    INDEX idx_ingredient (ingredient_type, ingredient_id),
    FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='配方配料关联表';

CREATE TABLE orders (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    order_no VARCHAR(50) NOT NULL UNIQUE,
    table_no VARCHAR(20),
    customer_count INT DEFAULT 1,
    total_amount DECIMAL(10,2) NOT NULL,
    discount DECIMAL(10,2) DEFAULT 0,
    actual_amount DECIMAL(10,2) NOT NULL,
    payment_method VARCHAR(30),
    status VARCHAR(20) NOT NULL DEFAULT 'completed',
    remark TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_order_no (order_no),
    INDEX idx_created_at (created_at),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='客单记录表';

CREATE TABLE order_items (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT NOT NULL,
    recipe_id BIGINT NOT NULL,
    recipe_name VARCHAR(100) NOT NULL,
    quantity INT NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    remark TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_order (order_id),
    INDEX idx_recipe (recipe_id),
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (recipe_id) REFERENCES recipes(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单明细表';

CREATE TABLE waste_records (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    ingredient_type VARCHAR(20) NOT NULL,
    ingredient_id BIGINT NOT NULL,
    ingredient_name VARCHAR(100) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    reason VARCHAR(200) NOT NULL,
    cost DECIMAL(10,2) NOT NULL,
    operator VARCHAR(50),
    remark TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_created_at (created_at),
    INDEX idx_ingredient (ingredient_type, ingredient_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='原料损耗表';

CREATE TABLE special_creations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    creator VARCHAR(50),
    inspiration TEXT,
    taste_profile VARCHAR(200),
    glass_type VARCHAR(50),
    serving_ml INT,
    price DECIMAL(10,2),
    preparation_method TEXT,
    garnish VARCHAR(200),
    ingredients_text TEXT,
    image_url VARCHAR(500),
    status VARCHAR(20) DEFAULT 'draft',
    tasting_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status),
    INDEX idx_creator (creator)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='特调新品存档表';

CREATE TABLE purchases (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    purchase_no VARCHAR(50) NOT NULL UNIQUE,
    supplier VARCHAR(100),
    total_amount DECIMAL(10,2) NOT NULL,
    purchase_date DATE NOT NULL,
    operator VARCHAR(50),
    remark TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_purchase_no (purchase_no),
    INDEX idx_purchase_date (purchase_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='采购台账表';

CREATE TABLE purchase_items (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    purchase_id BIGINT NOT NULL,
    ingredient_type VARCHAR(20) NOT NULL,
    ingredient_id BIGINT NOT NULL,
    ingredient_name VARCHAR(100) NOT NULL,
    quantity DECIMAL(10,2) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    batch_no VARCHAR(100),
    expiry_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_purchase (purchase_id),
    INDEX idx_ingredient (ingredient_type, ingredient_id),
    FOREIGN KEY (purchase_id) REFERENCES purchases(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='采购明细表';

INSERT INTO spirits (name, category, brand, origin, alcohol_content, volume_ml, unit, stock_quantity, min_stock, cost_price, selling_price_per_ml, description) VALUES
('Tanqueray Gin', '金酒', 'Tanqueray', '英国', 47.3, 700, '瓶', 12, 3, 280.00, 0.4000, '添加利伦敦干金酒，经典金酒代表'),
('Bombay Sapphire', '金酒', 'Bombay', '英国', 40.0, 750, '瓶', 10, 3, 220.00, 0.2933, '孟买蓝宝石金酒，十种植物精华'),
('Absolut Vodka', '伏特加', 'Absolut', '瑞典', 40.0, 700, '瓶', 15, 5, 120.00, 0.1714, '绝对伏特加原味'),
('Grey Goose', '伏特加', 'Grey Goose', '法国', 40.0, 750, '瓶', 8, 2, 380.00, 0.5067, '灰雁伏特加，法国高端伏特加'),
('Johnnie Walker Black', '威士忌', 'Johnnie Walker', '苏格兰', 40.0, 700, '瓶', 10, 3, 320.00, 0.4571, '尊尼获加黑牌12年调配威士忌'),
('Macallan 12 Year', '威士忌', 'Macallan', '苏格兰', 40.0, 700, '瓶', 5, 1, 680.00, 0.9714, '麦卡伦12年单一麦芽威士忌'),
('Havana Club 7 Year', '朗姆酒', 'Havana Club', '古巴', 40.0, 700, '瓶', 8, 2, 180.00, 0.2571, '哈瓦那俱乐部7年黑朗姆'),
('Bacardi Superior', '朗姆酒', 'Bacardi', '波多黎各', 37.5, 750, '瓶', 12, 4, 95.00, 0.1267, '百加得白朗姆酒'),
('Patron Silver', '龙舌兰', 'Patron', '墨西哥', 40.0, 750, '瓶', 6, 2, 450.00, 0.6000, '培恩银龙舌兰，高端龙舌兰'),
('Jose Cuervo Gold', '龙舌兰', 'Jose Cuervo', '墨西哥', 38.0, 750, '瓶', 10, 3, 150.00, 0.2000, '豪帅金标龙舌兰'),
('Cointreau', '力娇酒', 'Cointreau', '法国', 40.0, 700, '瓶', 8, 3, 210.00, 0.3000, '君度橙味力娇酒'),
('Baileys Irish Cream', '力娇酒', 'Baileys', '爱尔兰', 17.0, 700, '瓶', 6, 2, 180.00, 0.2571, '百利甜酒'),
('Martini Bianco', '味美思', 'Martini', '意大利', 15.0, 1000, '瓶', 5, 2, 130.00, 0.1300, '马天尼白味美思'),
('Campari', '利口酒', 'Campari', '意大利', 25.0, 1000, '瓶', 6, 2, 160.00, 0.1600, '金巴利利口酒'),
('Angostura Bitters', '苦精', 'Angostura', '委内瑞拉', 44.7, 200, '瓶', 10, 5, 80.00, 0.4000, '安格斯特拉苦精');

INSERT INTO ingredients (name, category, unit, stock_quantity, min_stock, cost_price, supplier, description) VALUES
('新鲜柠檬', '水果', '个', 50.00, 20.00, 3.50, '本地果蔬供应商', '用于调酒的新鲜黄柠檬'),
('新鲜青柠', '水果', '个', 40.00, 15.00, 4.00, '本地果蔬供应商', '用于调酒的新鲜青柠'),
('西柚汁', '果汁', 'ml', 2000.00, 500.00, 0.08, '果汁供应商', '鲜榨西柚汁'),
('橙汁', '果汁', 'ml', 3000.00, 1000.00, 0.06, '果汁供应商', '鲜榨橙汁'),
('蔓越莓汁', '果汁', 'ml', 2500.00, 800.00, 0.10, '果汁供应商', '进口蔓越莓汁'),
('菠萝汁', '果汁', 'ml', 2000.00, 500.00, 0.07, '果汁供应商', '进口菠萝汁'),
('苹果汁', '果汁', 'ml', 2000.00, 500.00, 0.05, '果汁供应商', '鲜榨苹果汁'),
('红石榴糖浆', '糖浆', 'ml', 1000.00, 300.00, 0.12, '调酒原料供应商', '红石榴糖浆'),
('单糖浆', '糖浆', 'ml', 2000.00, 500.00, 0.03, '自制', '自制1:1单糖浆'),
('蜂蜜糖浆', '糖浆', 'ml', 500.00, 200.00, 0.15, '自制', '蜂蜜糖浆'),
('焦糖糖浆', '糖浆', 'ml', 800.00, 200.00, 0.10, '调酒原料供应商', '焦糖糖浆'),
('香草糖浆', '糖浆', 'ml', 800.00, 200.00, 0.10, '调酒原料供应商', '香草糖浆'),
('苦艾酒', '辅料', 'ml', 500.00, 100.00, 0.20, '调酒原料供应商', '用于调制鸡尾酒'),
('蛋清', '辅料', '个', 30.00, 10.00, 1.50, '本地蛋商', '新鲜蛋清'),
('安格斯特拉苦精', '辅料', 'ml', 200.00, 50.00, 0.40, '调酒原料供应商', '苦精'),
('橙味苦精', '辅料', 'ml', 200.00, 50.00, 0.35, '调酒原料供应商', '橙味苦精'),
('苏打水', '汽水', 'ml', 5000.00, 2000.00, 0.01, '本地供应商', '苏打水'),
('汤力水', '汽水', 'ml', 4000.00, 1000.00, 0.02, '本地供应商', '汤力水'),
('干姜水', '汽水', 'ml', 3000.00, 1000.00, 0.02, '本地供应商', '干姜水'),
('可乐', '汽水', 'ml', 6000.00, 2000.00, 0.01, '本地供应商', '可口可乐'),
('薄荷', '香草', '枝', 50.00, 20.00, 0.50, '本地果蔬供应商', '新鲜薄荷叶'),
('罗勒', '香草', '枝', 30.00, 10.00, 0.80, '本地果蔬供应商', '新鲜罗勒叶'),
('迷迭香', '香草', '枝', 30.00, 10.00, 0.60, '本地果蔬供应商', '新鲜迷迭香'),
('盐', '调料', 'g', 500.00, 100.00, 0.01, '本地供应商', '细盐'),
('黑胡椒', '调料', 'g', 200.00, 50.00, 0.05, '本地供应商', '黑胡椒粉'),
('方糖', '调料', '块', 200.00, 50.00, 0.20, '本地供应商', '方糖');

INSERT INTO recipes (name, category, glass_type, serving_ml, price, cost, preparation_method, garnish, taste_profile, difficulty, is_signature, description) VALUES
('金汤力', '经典鸡尾酒', 'Highball', 250, 68.00, 15.00, '1. 杯子装满冰块\n2. 加入金酒45ml\n3. 用汤力水补满\n4. 轻轻搅拌\n5. 用青柠角装饰', '青柠角', '清爽、草本、微苦', '简单', FALSE, '最经典的金酒鸡尾酒，清爽解暑'),
('莫吉托', '经典鸡尾酒', 'Highball', 300, 78.00, 18.00, '1. 杯中放入薄荷叶和青柠角\n2. 加入单糖浆，轻压薄荷叶\n3. 加入朗姆酒和青柠汁\n4. 加满碎冰\n5. 用苏打水补满\n6. 轻轻搅拌', '薄荷叶、青柠角', '清爽、薄荷、酸甜', '简单', FALSE, '古巴经典，夏日必备'),
('长岛冰茶', '经典鸡尾酒', 'Highball', 350, 98.00, 35.00, '1. 杯中装满冰\n2. 加入朗姆、伏特加、龙舌兰、金酒各15ml\n3. 加入君度15ml，柠檬汁25ml\n4. 加入单糖浆\n5. 用可乐补满\n6. 轻轻搅拌', '柠檬片', '酸甜、酒味浓郁', '中等', FALSE, '传说中的失身酒，后劲十足'),
('威士忌酸', '经典鸡尾酒', 'Coupe', 180, 88.00, 28.00, '1. 雪克壶中加入威士忌60ml\n2. 加入柠檬汁25ml，单糖浆20ml\n3. 加入蛋清（可选）\n4. 加冰摇匀\n5. 双重过滤倒入杯中\n6. 滴入苦精装饰', '柠檬皮卷、安格斯特拉苦精', '酸甜、平衡、顺滑', '中等', FALSE, '经典酸酒，蛋清让口感更顺滑'),
('曼哈顿', '经典鸡尾酒', 'Coupe', 160, 98.00, 45.00, '1. 调酒杯中加入威士忌60ml\n2. 加入红味美思30ml，安格斯特拉苦精2滴\n3. 加冰搅拌30秒\n4. 过滤倒入冰镇杯\n5. 用马拉斯奇诺樱桃装饰', '马拉斯奇诺樱桃', '醇厚、香甜、复杂', '简单', FALSE, '鸡尾酒之王，优雅经典'),
('马天尼', '经典鸡尾酒', 'Martini', 150, 108.00, 48.00, '1. 调酒杯中加入金酒60ml，干味美思15ml\n2. 加冰搅拌20秒（越久越干）\n3. 过滤倒入冰镇马天尼杯\n4. 用橄榄或柠檬皮装饰', '橄榄、柠檬皮', '干冽、强劲、优雅', '简单', TRUE, '鸡尾酒中的传奇，007的最爱'),
('自由古巴', '经典鸡尾酒', 'Highball', 300, 68.00, 12.00, '1. 杯中装满冰\n2. 加入朗姆酒50ml\n3. 挤入青柠汁\n4. 用可乐补满\n5. 轻轻搅拌', '青柠角', '甜润、可乐香、轻松', '简单', FALSE, '源自古巴，朗姆配可乐'),
('玛格丽特', '经典鸡尾酒', 'Margarita', 180, 88.00, 25.00, '1. 杯口用青柠润湿，沾盐边\n2. 雪克壶加入龙舌兰50ml，君度25ml，青柠汁25ml\n3. 加冰摇匀\n4. 过滤倒入加冰的杯中', '盐边、青柠片', '酸甜、龙舌兰香、微咸', '简单', FALSE, '龙舌兰经典，盐边是灵魂'),
('特调-东方茉莉', '特调鸡尾酒', 'Coupe', 160, 128.00, 38.00, '1. 茉莉花茶浸泡的金酒50ml\n2. 加入蜂蜜糖浆15ml，柠檬汁20ml\n3. 加冰摇匀\n4. 过滤倒入冰镇杯\n5. 茉莉花瓣漂浮装饰', '茉莉花瓣', '花香、清甜、优雅', '中等', TRUE, '本店招牌特调，东方茶香与金酒的完美结合'),
('特调-烟熏古典', '特调鸡尾酒', 'Rock', 200, 158.00, 65.00, '1. 方糖用苦精和水化开\n2. 加入烟熏威士忌60ml\n3. 加入冰块搅拌\n4. 橙皮喷香后放入\n5. 烟熏罩注入烟熏风味', '橙皮卷、马拉斯奇诺樱桃', '烟熏、醇厚、复杂', '困难', TRUE, '本店招牌，使用烟熏工艺的高端特调');

INSERT INTO recipe_ingredients (recipe_id, ingredient_type, ingredient_id, amount, unit) VALUES
(1, 'spirit', 1, 45.00, 'ml'),
(1, 'ingredient', 17, 200.00, 'ml'),
(1, 'ingredient', 2, 1.00, '个'),
(2, 'spirit', 7, 50.00, 'ml'),
(2, 'ingredient', 2, 2.00, '个'),
(2, 'ingredient', 8, 15.00, 'ml'),
(2, 'ingredient', 20, 5.00, '枝'),
(2, 'ingredient', 16, 150.00, 'ml'),
(3, 'spirit', 7, 15.00, 'ml'),
(3, 'spirit', 3, 15.00, 'ml'),
(3, 'spirit', 10, 15.00, 'ml'),
(3, 'spirit', 1, 15.00, 'ml'),
(3, 'spirit', 11, 15.00, 'ml'),
(3, 'ingredient', 1, 25.00, 'ml'),
(3, 'ingredient', 8, 15.00, 'ml'),
(3, 'ingredient', 19, 200.00, 'ml'),
(4, 'spirit', 5, 60.00, 'ml'),
(4, 'ingredient', 1, 25.00, 'ml'),
(4, 'ingredient', 8, 20.00, 'ml'),
(4, 'ingredient', 13, 1.00, '个'),
(4, 'ingredient', 14, 2.00, '滴'),
(5, 'spirit', 5, 60.00, 'ml'),
(5, 'ingredient', 12, 30.00, 'ml'),
(5, 'ingredient', 14, 2.00, '滴'),
(6, 'spirit', 1, 60.00, 'ml'),
(6, 'ingredient', 12, 15.00, 'ml'),
(7, 'spirit', 7, 50.00, 'ml'),
(7, 'ingredient', 2, 1.00, '个'),
(7, 'ingredient', 19, 250.00, 'ml'),
(8, 'spirit', 10, 50.00, 'ml'),
(8, 'spirit', 11, 25.00, 'ml'),
(8, 'ingredient', 2, 25.00, 'ml'),
(8, 'ingredient', 23, 1.00, 'g'),
(9, 'spirit', 1, 50.00, 'ml'),
(9, 'ingredient', 1, 20.00, 'ml'),
(9, 'ingredient', 9, 15.00, 'ml'),
(10, 'spirit', 6, 60.00, 'ml'),
(10, 'ingredient', 14, 3.00, '滴'),
(10, 'ingredient', 25, 1.00, '块');

INSERT INTO orders (order_no, table_no, customer_count, total_amount, discount, actual_amount, payment_method, status, remark) VALUES
('ORD202401010001', 'A1', 2, 236.00, 0.00, 236.00, '微信支付', 'completed', '两位客人，情侣约会'),
('ORD202401010002', 'A2', 4, 458.00, 20.00, 438.00, '支付宝', 'completed', '朋友聚会'),
('ORD202401010003', 'B1', 1, 158.00, 0.00, 158.00, '现金', 'completed', '老顾客'),
('ORD202401020001', 'A3', 3, 392.00, 0.00, 392.00, '微信支付', 'completed', '商务会谈'),
('ORD202401020002', 'B2', 2, 286.00, 10.00, 276.00, '信用卡', 'completed', '生日庆祝'),
('ORD202401030001', 'A1', 2, 324.00, 0.00, 324.00, '微信支付', 'completed', ''),
('ORD202401030002', 'C1', 6, 886.00, 50.00, 836.00, '支付宝', 'completed', '公司团建'),
('ORD202401030003', 'A2', 2, 246.00, 0.00, 246.00, '微信支付', 'completed', '');

INSERT INTO order_items (order_id, recipe_id, recipe_name, quantity, unit_price, subtotal, remark) VALUES
(1, 1, '金汤力', 2, 68.00, 136.00, ''),
(1, 2, '莫吉托', 1, 78.00, 78.00, '少糖'),
(1, 7, '自由古巴', 1, 68.00, 68.00, '多冰'),
(2, 3, '长岛冰茶', 2, 98.00, 196.00, ''),
(2, 4, '威士忌酸', 2, 88.00, 176.00, '不要蛋清'),
(2, 8, '玛格丽特', 1, 88.00, 88.00, '盐边'),
(3, 10, '特调-烟熏古典', 1, 158.00, 158.00, '重烟熏'),
(4, 5, '曼哈顿', 2, 98.00, 196.00, ''),
(4, 6, '马天尼', 1, 108.00, 108.00, '干马天尼'),
(4, 9, '特调-东方茉莉', 1, 128.00, 128.00, ''),
(5, 9, '特调-东方茉莉', 2, 128.00, 256.00, '少糖'),
(5, 1, '金汤力', 1, 68.00, 68.00, ''),
(6, 10, '特调-烟熏古典', 2, 158.00, 316.00, ''),
(6, 2, '莫吉托', 1, 78.00, 78.00, ''),
(7, 3, '长岛冰茶', 3, 98.00, 294.00, ''),
(7, 6, '马天尼', 2, 108.00, 216.00, ''),
(7, 5, '曼哈顿', 2, 98.00, 196.00, ''),
(7, 8, '玛格丽特', 2, 88.00, 176.00, ''),
(8, 4, '威士忌酸', 2, 88.00, 176.00, ''),
(8, 7, '自由古巴', 1, 68.00, 68.00, '');

INSERT INTO waste_records (ingredient_type, ingredient_id, ingredient_name, amount, unit, reason, cost, operator, remark) VALUES
('ingredient', 1, '新鲜柠檬', 5.00, '个', '腐烂变质', 17.50, '张三', '周末未使用完'),
('ingredient', 20, '薄荷叶', 10.00, '枝', '枯萎', 5.00, '张三', '保存不当'),
('spirit', 15, 'Angostura Bitters', 10.00, 'ml', '操作失误', 4.00, '李四', '调酒时倒多了'),
('ingredient', 13, '蛋清', 3.00, '个', '变质', 4.50, '李四', '温度控制不当'),
('ingredient', 3, '西柚汁', 200.00, 'ml', '过期', 16.00, '王五', '进货太多');

INSERT INTO special_creations (name, creator, inspiration, taste_profile, glass_type, serving_ml, price, preparation_method, garnish, ingredients_text, status, tasting_notes) VALUES
('秋日桂香', '李大调酒师', '秋天的桂花香气，结合威士忌的醇厚', '桂花香气、圆润、微甜', 'Rock', 200, 138.00, '1. 桂花浸泡波本威士忌60ml\n2. 加入枫糖浆15ml，比特酒2滴\n3. 加冰搅拌\n4. 桂花装饰', '干桂花、橙皮', '波本威士忌60ml、枫糖浆15ml、桂花、苦精2滴', 'draft', '第一轮测试：桂花香气突出，但甜度需要调整'),
('热带风暴', '王创意调酒师', '东南亚热带风情，多种水果碰撞', '热带水果、清爽、酸甜', 'Tiki', 350, 118.00, '1. 白朗姆50ml，黑朗姆15ml\n2. 菠萝汁60ml，西柚汁30ml\n3. 百香果糖浆20ml\n4. 加冰摇匀，倒入碎冰', '菠萝叶、热带水果串', '白朗姆50ml、黑朗姆15ml、菠萝汁60ml、西柚汁30ml、百香果糖浆20ml', 'approved', '顾客反馈很好，准备加入正式菜单'),
('玫瑰人生', '张首席调酒师', '浪漫的玫瑰花香，适合女士饮用', '玫瑰香、细腻、甜美', 'Coupe', 180, 148.00, '1. 玫瑰金酒50ml\n2. 圣哲曼利口酒20ml\n3. 柠檬汁15ml，玫瑰糖浆10ml\n4. 加冰摇匀，倒入冰镇杯', '玫瑰花瓣', '玫瑰金酒50ml、圣哲曼利口酒20ml、柠檬汁15ml、玫瑰糖浆10ml', 'testing', '第三轮测试，口感已经很完美，下周推出');

INSERT INTO purchases (purchase_no, supplier, total_amount, purchase_date, operator, remark) VALUES
('PUR20240101001', '洋酒进口商', 5600.00, '2024-01-01', '库管小李', '月度常规采购'),
('PUR20240101002', '果蔬供应商', 1250.00, '2024-01-02', '库管小李', '每周水果采购'),
('PUR20240105001', '调酒原料供应商', 890.00, '2024-01-05', '库管小李', '糖浆和苦精补充'),
('PUR20240108001', '洋酒进口商', 3400.00, '2024-01-08', '库管小李', '威士忌补货');

INSERT INTO purchase_items (purchase_id, ingredient_type, ingredient_id, ingredient_name, quantity, unit, unit_price, subtotal, batch_no, expiry_date) VALUES
(1, 'spirit', 1, 'Tanqueray Gin', 5.00, '瓶', 280.00, 1400.00, 'GIN20240101', '2026-01-01'),
(1, 'spirit', 5, 'Johnnie Walker Black', 4.00, '瓶', 320.00, 1280.00, 'WHISKY2024001', '2027-01-01'),
(1, 'spirit', 7, 'Havana Club 7 Year', 5.00, '瓶', 180.00, 900.00, 'RUM2024001', '2026-06-01'),
(1, 'spirit', 10, 'Jose Cuervo Gold', 6.00, '瓶', 150.00, 900.00, 'TEQ2024001', '2026-03-01'),
(1, 'spirit', 11, 'Cointreau', 4.00, '瓶', 210.00, 840.00, 'LIQ2024001', '2026-12-01'),
(1, 'spirit', 3, 'Absolut Vodka', 4.00, '瓶', 120.00, 480.00, 'VOD2024001', '2026-09-01'),
(2, 'ingredient', 1, '新鲜柠檬', 60.00, '个', 3.50, 210.00, NULL, '2024-01-10'),
(2, 'ingredient', 2, '新鲜青柠', 50.00, '个', 4.00, 200.00, NULL, '2024-01-10'),
(2, 'ingredient', 20, '薄荷', 80.00, '枝', 0.50, 40.00, NULL, '2024-01-08'),
(2, 'ingredient', 13, '蛋清', 30.00, '个', 1.50, 45.00, NULL, '2024-01-05'),
(2, 'ingredient', 3, '西柚汁', 3000.00, 'ml', 0.08, 240.00, 'JUICE2024001', '2024-01-15'),
(2, 'ingredient', 4, '橙汁', 3000.00, 'ml', 0.06, 180.00, 'JUICE2024002', '2024-01-15'),
(2, 'ingredient', 5, '蔓越莓汁', 2000.00, 'ml', 0.10, 200.00, 'JUICE2024003', '2024-01-20'),
(2, 'ingredient', 6, '菠萝汁', 2000.00, 'ml', 0.07, 140.00, 'JUICE2024004', '2024-01-20'),
(3, 'ingredient', 7, '红石榴糖浆', 1000.00, 'ml', 0.12, 120.00, 'SYR2024001', '2024-06-01'),
(3, 'ingredient', 8, '单糖浆', 2000.00, 'ml', 0.03, 60.00, NULL, '2024-02-01'),
(3, 'ingredient', 14, '安格斯特拉苦精', 500.00, 'ml', 0.40, 200.00, 'BIT2024001', '2026-01-01'),
(3, 'ingredient', 15, '橙味苦精', 500.00, 'ml', 0.35, 175.00, 'BIT2024002', '2026-01-01'),
(3, 'ingredient', 16, '苏打水', 10000.00, 'ml', 0.01, 100.00, NULL, '2024-03-01'),
(3, 'ingredient', 17, '汤力水', 6000.00, 'ml', 0.02, 120.00, NULL, '2024-03-01'),
(3, 'ingredient', 18, '干姜水', 5000.00, 'ml', 0.02, 100.00, NULL, '2024-03-01'),
(3, 'ingredient', 19, '可乐', 12000.00, 'ml', 0.01, 120.00, NULL, '2024-03-01'),
(4, 'spirit', 6, 'Macallan 12 Year', 2.00, '瓶', 680.00, 1360.00, 'MAC2024001', '2030-01-01'),
(4, 'spirit', 9, 'Patron Silver', 2.00, '瓶', 450.00, 900.00, 'PAT2024001', '2028-01-01'),
(4, 'spirit', 4, 'Grey Goose', 2.00, '瓶', 380.00, 760.00, 'GOOSE2024001', '2027-01-01'),
(4, 'spirit', 12, 'Baileys Irish Cream', 2.00, '瓶', 180.00, 360.00, 'BAI2024001', '2024-12-01');

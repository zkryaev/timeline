CREATE TABLE IF NOT EXISTS users(
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    passwd_hash TEXT NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    telephone VARCHAR(255) NOT NULL,
    social TEXT,
    about TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    verified BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS user_verify(
    user_verify_id SERIAL PRIMARY KEY,
    user_id INT,  -- добавлен столбец для связи с таблицей users
    code VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS orgs(
    org_id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    passwd_hash TEXT NOT NULL,
    org_name VARCHAR(255) NOT NULL,
    org_address VARCHAR(300) NOT NULL,
    telephone VARCHAR(255) NOT NULL,
    social VARCHAR(255),
    about TEXT,
    lat FLOAT NOT NULL,
    long FLOAT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    verified BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS city(
    city_id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS orgs_verify(
    org_verify_id SERIAL PRIMARY KEY,
    org_id INT NOT NULL,
    code VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (org_id) REFERENCES orgs(org_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS orgs_city(
    city_id INT,
    org_id INT,
    FOREIGN KEY (city_id) REFERENCES city(city_id) ON DELETE CASCADE,
    FOREIGN KEY (org_id) REFERENCES orgs(org_id) ON DELETE CASCADE 
);




CREATE TABLE IF NOT EXISTS users(
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    passwd_hash TEXT NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    telephone VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    about TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    verified BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS user_verify(
    user_verify_id SERIAL PRIMARY KEY,
    user_id INT,
    code VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS orgs(
    org_id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    passwd_hash TEXT NOT NULL,
    name VARCHAR(255) NOT NULL,
    rating NUMERIC(3, 2) DEFAULT 0.00,
    type VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    address VARCHAR(300) NOT NULL,
    telephone VARCHAR(255),
    lat FLOAT NOT NULL,
    long FLOAT NOT NULL,
    about TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    verified BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS orgs_verify(
    org_verify_id SERIAL PRIMARY KEY,
    org_id INT NOT NULL,
    code VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (org_id) REFERENCES orgs(org_id) ON DELETE CASCADE
);




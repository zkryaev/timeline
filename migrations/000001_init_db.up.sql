CREATE TABLE IF NOT EXISTS users (
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

CREATE TABLE IF NOT EXISTS users_verify (
    user_verify_id SERIAL PRIMARY KEY,
    user_id INT,
    code VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS orgs (
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

CREATE TABLE IF NOT EXISTS timetables (
    timetable_id SERIAL PRIMARY KEY,
    org_id INT NOT NULL,
    weekday INT NOT NULL,
    open TIMESTAMP NOT NULL,
    close TIMESTAMP NOT NULL,
    break_start TIMESTAMP NOT NULL,
    break_end TIMESTAMP NOT NULL,
    FOREIGN KEY (org_id) REFERENCES orgs(org_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS orgs_verify (
    org_verify_id SERIAL PRIMARY KEY,
    org_id INT NOT NULL,
    code VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (org_id) REFERENCES orgs(org_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS services (
    service_id SERIAL PRIMARY KEY,
    org_id INT,
    name VARCHAR(300) NOT NULL,
    cost NUMERIC(15,2) NOT NULL,
    description VARCHAR(400),
    FOREIGN KEY (org_id) REFERENCES orgs(org_id)
);
CREATE TABLE IF NOT EXISTS workers (
    worker_id SERIAL PRIMARY KEY,
    org_id INT,  -- Запятая добавлена
    service_id INT,  -- Запятая добавлена
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255),
    position VARCHAR(300),
    degree VARCHAR(300),
    FOREIGN KEY (service_id) REFERENCES services(service_id),
    FOREIGN KEY (org_id) REFERENCES orgs(org_id)
);

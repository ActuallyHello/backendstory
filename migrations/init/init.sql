-- Создание базы данных для приложения
CREATE DATABASE IF NOT EXISTS backendstory;

-- Создание базы данных для Keycloak
CREATE DATABASE IF NOT EXISTS backendstory_keycloak;

-- Создание пользователя для приложения (пароль должен совпадать с MYSQL_PASSWORD в .env)
CREATE USER IF NOT EXISTS 'app_user'@'%' IDENTIFIED BY 'app_password';
GRANT ALL PRIVILEGES ON backendstory.* TO 'app_user'@'%';

-- Создание пользователя для Keycloak (пароль должен совпадать с KEYCLOAK_DB_PASSWORD в .env)
CREATE USER IF NOT EXISTS 'kc_user'@'%' IDENTIFIED BY 'kc_password';
GRANT ALL PRIVILEGES ON backendstory_keycloak.* TO 'kc_user'@'%';

-- Обновление прав
FLUSH PRIVILEGES;
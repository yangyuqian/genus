CREATE USER 'scauser'@'172.17.0.1' IDENTIFIED BY 'scauser';
CREATE USER 'scauser'@'127.0.0.1' IDENTIFIED BY 'scauser';
GRANT ALL PRIVILEGES ON *.* TO 'scauser'@'172.17.0.1';
GRANT ALL PRIVILEGES ON *.* TO 'scauser'@'127.0.0.1';

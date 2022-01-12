import os


dsl = {
    'database': os.getenv('DBNAME', 'postgres'),
    'user': os.getenv('DBUSER', 'postgres'),
    'password': os.getenv('DBPASS', 'mysecretpassword'),
    'host': os.getenv('DBHOST', '127.0.0.1'),
    'port': os.getenv('DBPORT', 5432),
}

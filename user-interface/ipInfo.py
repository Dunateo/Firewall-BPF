#pip3 install ipinfo
import ipinfo
import asyncio
import sqlite3
import logging
from contextlib import closing

DATABASE_NAME="connexion.db"
access_token = '7fcbfbb488fd6c'
handler = ipinfo.getHandler(access_token)



#connect to the data base 
def dataConnexion(name):
    connection = sqlite3.connect(name)
    if(connection.total_changes!=0):
        return False
    return connection

#create table
def creationTable(cursor):
    cursor.execute("CREATE TABLE connex (id INTEGER, lon REAL, lat REAL, name TEXT)")

#insert a new data
def insertData(cursor, value):
    try:
        text = value.city + " "+ value.ip + " " +value.hostname
    except AttributeError:
        text = value.city + " "+ value.ip 
    
    prepare = [value.ip, value.longitude, value.latitude, text]
    print(prepare)
    cursor.execute('INSERT INTO connex VALUES(?,?,?,?)', prepare)
    logging.info('New element'+text)

#read data from db
def read_from_db(c):
    c.execute ('SELECT * FROM connex ORDER BY ROWID DESC ') 
    for row in c.fetchall():
        print(row)

#read data from db
def get_from_db(c):
    c.execute ('SELECT * FROM connex ORDER BY ROWID DESC LIMIT 10') 
    return c.fetchall()

#proper close of db
def closeConnect(connection, cursor, name):
    with closing(sqlite3.connect(name)) as connection:
        with closing(connection.cursor()) as cursor:
            rows = cursor.execute("SELECT 1").fetchall()
            print(rows)

#asynchrone request 
def do_req(ip_address):
    details = handler.getDetails(ip_address)
    return details


#launch the async details
def get_ip_details(ip_address):
    #creation sqllite
    try:
        with open(DATABASE_NAME): pass
    except IOError:
        logging.error('No database')
        logging.info('Database creation')
        connect = dataConnexion(DATABASE_NAME)
        cursor = connect.cursor()
        creationTable(cursor)
    
    #get info IPInfo
    details = do_req(ip_address)
    #print(details.all)

    #new data in 
    connect = dataConnexion(DATABASE_NAME)
    cursor = connect.cursor()
    insertData(cursor, details)

     #commit the change and close sql 
    connect.commit()
    closeConnect(connect, cursor,DATABASE_NAME)
    


ip_address = '25.112.82.140'
get_ip_details(ip_address)

#connect = dataConnexion(DATABASE_NAME)
#cursor = connect.cursor()
#read_from_db(cursor)
import ipinfo
import asyncio
import sqlite3
import logging
from contextlib import closing

DATABASE_NAME="connexion.db"
access_token = '7fcbfbb488fd6c'
handler = ipinfo.getHandlerAsync(access_token)



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
    text = value.city + " "+ value.ip + " " +value.hostname
    prepare = [value.ip, value.longitude, value.latitude, text]
    print(prepare)
    cursor.execute('INSERT INTO connex VALUES(?,?,?,?)', prepare)
    logging.info('New element'+text)

#read data from db
def read_from_db(c):
    c.execute ('SELECT * FROM connex') 
    for row in c.fetchall():
        print(row)

#proper close of db
def closeConnect(connection, cursor, name):
    with closing(sqlite3.connect(name)) as connection:
        with closing(connection.cursor()) as cursor:
            rows = cursor.execute("SELECT 1").fetchall()
            print(rows)

#asynchrone request 
async def do_req(ip_address):
    details = await handler.getDetails(ip_address)
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
    loop = asyncio.get_event_loop()
    details = loop.run_until_complete(do_req(ip_address))
    print(details.all)

    #new data in 
    connect = dataConnexion(DATABASE_NAME)
    cursor = connect.cursor()
    insertData(cursor, details)

     #commit the change and close sql 
    connect.commit()
    closeConnect(connect, cursor,DATABASE_NAME)

ip_address = '109.13.195.67'
get_ip_details(ip_address)

connect = dataConnexion(DATABASE_NAME)
cursor = connect.cursor()
read_from_db(cursor)
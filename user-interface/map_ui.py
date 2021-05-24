# import the library
#python3 -m pip install folium
#python3 -m pip install pandas
import folium
import pandas as pd
import logging
import ipInfo as inf

DATABASE_NAME="connexion.db"


def generate_Map(data):

    # Make an empty map
    m = folium.Map(location=[20,0], tiles="OpenStreetMap", zoom_start=2)

    for i in range(0,len(data)):
        folium.Marker(location=[data.iloc[i]['lat'], data.iloc[i]['lon']],popup=data.iloc[i]['name'],).add_to(m)

    # Show the map
    m

    m.save('./map.html')

def createMap():
    logging.info("Map creation")
    connect = inf.dataConnexion(DATABASE_NAME)
    cursor = connect.cursor()
    
    latitude=[]
    longitude=[]
    name=[]
    value=[]

    result = inf.get_from_db(cursor)
    for row in result:
        latitude.append(row[2])
        longitude.append(row[1])
        value.append(row[0])
        name.append(row[3])

    #Make a data frame with dots to show on the map
    data = pd.DataFrame({
    'lon':longitude,
    'lat':latitude,
    'name':name,
    'value':value
    }, dtype=str)

    generate_Map(data)

    inf.closeConnect(connect, cursor,DATABASE_NAME)
    logging.info("Map created")

createMap()
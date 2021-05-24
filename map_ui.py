# import the library
#python3 -m pip install folium
#python3 -m pip install pandas
import folium
import pandas as pd
import logging

DATABASE_NAME="connexion.db"

#Make a data frame with dots to show on the map
data = pd.DataFrame({
   'lon':[-58, 2, 145, 30.32, -4.03, -73.57, 36.82, -38.5],
   'lat':[-34, 49, -38, 59.93, 5.33, 45.52, -1.29, -12.97],
   'name':['Buenos Aires', 'Paris', 'melbourne', 'St Petersbourg', 'Abidjan', 'Montreal', 'Nairobi', 'Salvador'],
   'value':[10, 12, 40, 70, 23, 43, 100, 43]
}, dtype=str)

# Make an empty map
m = folium.Map(location=[20,0], tiles="OpenStreetMap", zoom_start=2)

for i in range(0,len(data)):
    folium.Marker(location=[data.iloc[i]['lat'], data.iloc[i]['lon']],popup=data.iloc[i]['name'],).add_to(m)

# Show the map
m

m.save('./map.html')

def createMap():

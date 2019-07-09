import json
import requests

req = requests.get('https://api-v3.mbta.com/stops')
jsonString = req.text
data = json.loads(jsonString)
stops = data['data']

IDList = {}

for stop in stops:
    relationships = stop['relationships']
    parent_station = relationships['parent_station']
    if parent_station['data'] == None:
        attributes = stop['attributes']
        IDList[attributes['name']] = stop['id']

print('Ready:')
while True:
    line = input()
    if line == 'q':
        break
    try:
        print(IDList[line])
    except:
        print('No ID Found')
    print()

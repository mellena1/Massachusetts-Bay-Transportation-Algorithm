import json
import requests
import pprint

req = requests.get("https://api-v3.mbta.com/stops?filter%5Broute_type%5D=0,1")
jsonString = req.text
data = json.loads(jsonString)
stops = data['data']

stopnames = []
for stop in stops:
    stopnames.append(stop['attributes']['name'])

stopnames = list(dict.fromkeys(stopnames))

for k in stopnames:
    print(k)

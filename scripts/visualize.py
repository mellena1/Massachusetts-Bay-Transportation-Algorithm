import networkx as nx
import json
import graphviz
import matplotlib.pyplot as plt

with open("../graph.json", 'r') as json_file:
    data = json.load(json_file)
    station_map = {}
    graph = nx.Graph()
    node_colors = []
    edge_colors = []

    for station in data:
        station_id = station['id']
        station_name = station['name']
        
        if station_name in ['Heath Street', 'Cleveland Circle', 'Boston College', 'Riverside', 'Lechmere', 'Copley', 'Kenmore']:
            node_colors.append('green')
        elif station_name in ['Oak Grove', 'Forest Hills']:
            node_colors.append('orange')
        elif station_name in ['Braintree', 'Mattapan', 'Alewife', 'JFK/UMass']:
            node_colors.append('red')
        elif station_name in ['Bowdoin', 'Wonderland']:
            node_colors.append('dodgerblue')
        elif len(station['edges']) > 1:
            node_colors.append('lightgrey')

        station_map[station_id] = station_name
        graph.add_node(station_name)
    
    for station in data:
        u_name = station['name']
        for edge in station['edges']:
            v_name = station_map[edge]
            uv = set([u_name, v_name])

            if graph.has_edge(u_name, v_name):
                continue
            else:
                graph.add_edge(u_name, v_name)
                if uv.issubset(['North Station', 'Haymarket']):
                    edge_colors.append('gray')
                elif uv.issubset(['Heath Street', 'Cleveland Circle', 'Boston College', 'Riverside', 'Lechmere', 'Copley', 'Kenmore', 'Park Street', 'North Station', 'Government Center', 'Haymarket']):
                    edge_colors.append('green')
                elif uv.issubset(['Oak Grove', 'Forest Hills', 'North Station', 'Downtown Crossing', 'State', 'Haymarket']):
                    edge_colors.append('orange')
                elif uv.issubset(['Braintree', 'Mattapan', 'Alewife', 'JFK/UMass', 'Park Street', 'Downtown Crossing']):
                    edge_colors.append('red')
                elif uv.issubset(['Bowdoin', 'Wonderland', 'State', 'Government Center']):
                    edge_colors.append('dodgerblue')
    
    pos = nx.spring_layout(graph, k=0.25, iterations=100)
    nx.draw(graph, with_labels=True, pos=pos, node_color=node_colors, edge_color=edge_colors)
    plt.show()
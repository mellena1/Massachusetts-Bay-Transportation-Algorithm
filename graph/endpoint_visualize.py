import networkx as nx
import json
import graphviz
import matplotlib.pyplot as plt

with open("../endpoint_graph.json", 'r') as json_file:
    data = json.load(json_file)
    graph = nx.Graph()
    node_colors = []
    edge_colors = []

    for station in data:
        station_id = station['id']
        station_name = station['name']
        
        if station_name in ['Heath Street', 'Cleveland Circle', 'Boston College', 'Riverside', 'Lechmere']:
            node_colors.append('green')
        elif station_name in ['Oak Grove', 'Forest Hills']:
            node_colors.append('orange')
        elif station_name in ['Braintree', 'Mattapan', 'Alewife']:
            node_colors.append('red')
        elif station_name in ['Bowdoin', 'Wonderland']:
            node_colors.append('dodgerblue')

        graph.add_node(station_name)
    
    for station in data:
        u_name = station['name']
        for station in data:
            v_name = station['name']
            if u_name == v_name:
                continue
            if graph.has_edge(u_name, v_name):
                continue
            else:
                graph.add_edge(u_name, v_name)
                edge_colors.append('gray')
    
    nx.draw(graph, with_labels=True, node_color=node_colors, edge_color=edge_colors)
    plt.show()
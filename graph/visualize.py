import networkx as nx
import json
import graphviz
import matplotlib.pyplot as plt

with open("../graph.json", 'r') as json_file:
    data = json.load(json_file)
    station_map = {}
    graph = nx.Graph()

    for station in data:
        station_id = station['id']
        station_name = station['name']
        station_map[station_id] = station_name

        graph.add_node(station_name)
    
    for station in data:
        u_name = station['name']
        for edge in station['edges']:
            v_name = station_map[edge]
            graph.add_edge(u_name, v_name)
    
    print("Nodes of graph: ")
    print(graph.nodes)
    print("Edges of graph: ")
    print(graph.edges)
    
    nx.draw(graph, with_labels=True)
    plt.savefig("graph.png")  # save as png
    plt.show()
        
    

# -*- coding: utf-8 -*-

def breadth_walk(graph, root):
	if not graph.has_key(root):
		return 'Invalid root'

	child_nodes = [root]

	while len(child_nodes):
		elem = child_nodes.pop(0)

		print(elem)

		child_nodes += graph[elem]


graph = {'A': ['B', 'C'],
             'B': ['D', 'E'],
             'C': ['F', 'G'],
             'D': [],
             'E': ['H'],
             'F': [],
             'G': [],
              'H': []}

breadth_walk(graph, 'B')

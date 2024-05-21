ANT COLONY SIMULATION
The program simulates an ant colony navigating through a graph from a start node to an end node. It reads the graph configuration from a file, processes the graph to find paths, and then simulates the ants moving through these paths while avoiding collisions.

FEATURES
Reads graph data from a file and initializes nodes and edges.
Finds all paths from the start node to the end node using BFS.
Filters paths to avoid node collisions.
Simulates multiple ants moving through the graph.
Handles alternative paths if the main path is blocked.
Prints the movement of ants at each step.
Measures and prints the total time taken for the simulation.

HOW TO RUN İT
You can follow the steps below to run the project:
First, clone the project or download the files.
In the terminal or command client, navigate to the directory where the project is located.
Use the following command to run the project, providing the graph file as an argument.

USAGE
go run main.go [filename]

Example:
go run main.go graph.txt

AUTHOR
GÜLAY YÜCEER
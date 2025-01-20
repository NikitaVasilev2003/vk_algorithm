from collections import deque, defaultdict
import sys

MAX_EDGE_WEIGHT = 9
NULL_PARENT = -1

# Check if the start and end positions are the same
def is_direct_answer(start_row, start_col, end_row, end_col):
    return start_row == end_row and start_col == end_col

# Print the path from start to end
def print_path(path):
    for row, col in reversed(path):
        print(f"{row} {col}")
    print('.')

# Build the graph based on grid weights
def build_graph(rows, cols, weights):
    graph = defaultdict(list)
    for r in range(rows):
        for c in range(cols):
            current_index = r * cols + c
            if weights[current_index] == 0:  # Skip impassable cells
                continue

            # Add edges to valid neighbors
            if r > 0 and weights[current_index - cols] != 0:  # Neighbor above
                graph[current_index].append(current_index - cols)
            if r < rows - 1 and weights[current_index + cols] != 0:  # Neighbor below
                graph[current_index].append(current_index + cols)
            if c > 0 and weights[current_index - 1] != 0:  # Neighbor to the left
                graph[current_index].append(current_index - 1)
            if c < cols - 1 and weights[current_index + 1] != 0:  # Neighbor to the right
                graph[current_index].append(current_index + 1)
    return graph

# Reconstruct the path from the end position to the start
def reconstruct_path(cols, start_row, start_col, end_row, end_col, parents):
    path = []
    while (end_row, end_col) != (start_row, start_col):
        path.append((end_row, end_col))
        index = parents[end_row * cols + end_col]
        if index == NULL_PARENT:  # No path exists
            return None
        end_row, end_col = divmod(index, cols)
    path.append((start_row, start_col))
    return path

# Perform BFS with multi-level queues to account for weighted edges
def bfs(graph, distances, queues, visited, weights, parents):
    level = 0
    active_nodes = 1
    while active_nodes > 0:
        # Find the next non-empty queue
        while not queues[level % (MAX_EDGE_WEIGHT + 1)]:
            level += 1

        current_node = queues[level % (MAX_EDGE_WEIGHT + 1)].popleft()
        active_nodes -= 1

        if visited[current_node]:  # Skip already visited nodes
            continue

        visited[current_node] = True
        for neighbor in graph[current_node]:
            weight = weights[neighbor]
            if distances[current_node] + weight < distances[neighbor]:
                distances[neighbor] = distances[current_node] + weight
                parents[neighbor] = current_node
                queues[distances[neighbor] % (MAX_EDGE_WEIGHT + 1)].append(neighbor)
                active_nodes += 1

def main():
    try:
        # Read the grid dimensions
        rows, cols = map(int, input().strip().split())

        # Read the grid weights
        weights = []
        for _ in range(rows):
            weights.extend(map(int, input().strip().split()))

        # Read the start and end positions
        start_row, start_col, end_row, end_col = map(int, input().strip().split())

        # Validate the input coordinates
        if not (0 <= start_row < rows and 0 <= end_row < rows and
                0 <= start_col < cols and 0 <= end_col < cols):
            print("Incorrect input", file=sys.stderr)
            return

        # If the start and end positions are the same, print the result directly
        if is_direct_answer(start_row, start_col, end_row, end_col):
            print_path([(start_row, start_col)])
            return

        # Build the graph representation
        graph = build_graph(rows, cols, weights)

        # Initialize BFS-related structures
        start_index = start_row * cols + start_col
        distances = [float('inf')] * (rows * cols)
        parents = [NULL_PARENT] * (rows * cols)
        visited = [False] * (rows * cols)
        queues = [deque() for _ in range(MAX_EDGE_WEIGHT + 1)]

        # Set up the starting point
        distances[start_index] = 0
        parents[start_index] = start_index
        queues[0].append(start_index)

        # Perform BFS
        bfs(graph, distances, queues, visited, weights, parents)

        # Reconstruct and print the path
        path = reconstruct_path(cols, start_row, start_col, end_row, end_col, parent